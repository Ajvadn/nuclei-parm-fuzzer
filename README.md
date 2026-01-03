# Nuclei Parameter Fuzzer 🚀

A high-performance DAST scanning tool written in Go. It automates URL discovery, parameter filtering, and vulnerability scanning using Nuclei. Optimized for speed, reliability, and security research.

## Features ✨

- **Parallel URL Discovery**: Simultaneously fetches URLs from `gau`, `waybackurls`, `katana`, `paramspider`, `hakrawler`, and `waymore`.
- **Smart Filtering**: Uses `uro` to filter for unique URLs with query parameters, reducing noise and scan time.
- **Liveness Checking**: Integrated with `httpx` to ensure only live targets are scanned.
- **DAST Scanning**: Leverages `nuclei` for powerful, template-based vulnerability detection.
- **Speed Optimized**: Fine-tuned concurrency and parallel processing for rapid results.
- **Self-Healing**: Automatically installs missing dependencies like `gau`, `katana`, etc.

## Prerequisites 🛠️

Ensure you have the following installed:
- [Go](https://golang.org/doc/install) (1.21+)
- [Python 3](https://www.python.org/downloads/) & `pip3`

## Installation 📥

1. Clone the repository:
   ```bash
   git clone https://github.com/YOUR_USERNAME/nuclei-parmter-fuzz.git
   cd nuclei-parmter-fuzz
   ```

2. Install the tool:
   ```bash
   go install .
   ```

3. Run it!
   ```bash
   nuclei-parm-fuzzer -h
   ```

   (Ensure your `$GOPATH/bin` is in your `$PATH`)

## Usage 🚀

### Scan a single domain
```bash
nuclei-parm-fuzzer -d example.com
```

### Scan multiple domains from a file
```bash
nuclei-parm-fuzzer -f domains.txt
```

### Update all tools and templates
```bash
nuclei-parm-fuzzer --update
```

## Options ⚙️

| Option | Description |
| :--- | :--- |
| `-d`, `--domain` | Target a single domain |
| `-f`, `--file` | File containing a list of domains |
| `-u`, `--update` | Update all backend tools and Nuclei templates |
| `-h`, `--help` | Show the help message |

## Disclaimer ⚠️

This tool is for educational and authorized security testing purposes only. Running it against targets without explicit permission is illegal.

---
**Author**: Ajvad-N
