package ui

import (
	"fmt"
	"os"
	"time"

	"github.com/schollz/progressbar/v3"
)

const (
	Banner = `
 ___   _   ___ __  __   ___ _   _ __________ ___ 
| _ \ /_\ | _ \  \/  | | __| | | |_  /_  / __| _ \
|  _/ _ \|   / |\/| | | _|| |_| |/ / / /| _||   /
|_|/_/ \_\_|\_\_|  |_| |_|  \___//___/___|___|_|_\
                                 v2.0 - Enterprise
`
	Cyan   = "\033[36m"
	Green  = "\033[32m"
	Yellow = "\033[33m"
	Red    = "\033[31m"
	Reset  = "\033[0m"
)

var GlobalBar *progressbar.ProgressBar

func PrintBanner() {
	fmt.Println(Cyan + Banner + Reset)
	fmt.Println(Yellow + "           Aggressive Parameter & JS Discovery" + Reset)
	fmt.Println()
}

func NewProgressBar(total int) {
	GlobalBar = progressbar.NewOptions(total,
		progressbar.OptionEnableColorCodes(true),
		progressbar.OptionShowBytes(false),
		progressbar.OptionSetWidth(40),
		progressbar.OptionSetDescription("[cyan][reset] Discovering URLs..."),
		progressbar.OptionSetTheme(progressbar.Theme{
			Saucer:        "[green]=[reset]",
			SaucerHead:    "[green]>[reset]",
			SaucerPadding: " ",
			BarStart:      "[",
			BarEnd:        "]",
		}))
}

func Log(level, msg string) {
	timestamp := time.Now().Format("15:04:05")
	color := ""
	switch level {
	case "INFO":
		color = Cyan
	case "SUCCESS":
		color = Green
	case "WARN":
		color = Yellow
	case "ERROR":
		color = Red
	}
	
	output := fmt.Sprintf("%s[%s] %s%s%s\n", color, timestamp, Reset, msg, Reset)
	
	if GlobalBar != nil {
		fmt.Fprint(os.Stderr, "\r\033[K") // Clear progress bar line
		fmt.Print(output)
		GlobalBar.RenderBlank()
	} else {
		fmt.Print(output)
	}
}
