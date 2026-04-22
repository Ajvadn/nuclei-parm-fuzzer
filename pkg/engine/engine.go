package engine

import (
	"context"
	"sync"
)

// Provider defines the interface for all discovery sources.
type Provider interface {
	Name() string
	Execute(ctx context.Context, target string, isFile bool) ([]string, error)
}

// Result represents the output from a provider.
type Result struct {
	ProviderName string
	URLs         []string
	Err          error
}

// Engine coordinates the discovery process.
type Engine struct {
	Providers   []Provider
	Concurrency int
}

// NewEngine creates a new discovery engine.
func NewEngine(concurrency int) *Engine {
	if concurrency <= 0 {
		concurrency = 5
	}
	return &Engine{
		Concurrency: concurrency,
	}
}

// AddProvider adds a discovery source to the engine.
func (e *Engine) AddProvider(p Provider) {
	e.Providers = append(e.Providers, p)
}

// Run executes all registered providers in parallel with concurrency limit.
func (e *Engine) Run(ctx context.Context, target string, isFile bool) <-chan Result {
	results := make(chan Result, len(e.Providers))
	var wg sync.WaitGroup
	sem := make(chan struct{}, e.Concurrency)

	for _, p := range e.Providers {
		wg.Add(1)
		go func(p Provider) {
			defer wg.Done()
			sem <- struct{}{}        // Acquire
			defer func() { <-sem }() // Release

			urls, err := p.Execute(ctx, target, isFile)
			results <- Result{
				ProviderName: p.Name(),
				URLs:         urls,
				Err:          err,
			}
		}(p)
	}

	go func() {
		wg.Wait()
		close(results)
	}()

	return results
}
