# Nuclei Parameter Fuzzer 🚀

A high-performance DAST scanning tool that automates URL discovery, parameter filtering, and vulnerability scanning using Nuclei. Optimized for speed and security research.

## Features ✨

- **Parallel URL Discovery**: Simultaneously fetches URLs from `gau`, `waybackurls`, `katana`, and `paramspider`.
- **Smart Filtering**: Uses `uro` to filter for unique URLs with query parameters, reducing noise and scan time.
- **Liveness Checking**: Integrated with `httpx` to ensure only live targets are scanned.
- **DAST Scanning**: Leverages `nuclei` for powerful, template-based vulnerability detection.
- **Speed Optimized**: Fine-tuned concurrency and parallel processing for rapid results.
- **One-Click Updates**: Keep all your tools and templates up to date with a single command.

## Prerequisites 🛠️

Ensure you have the following installed:
- [Go](https://golang.org/doc/install)
- [Python 3](https://www.python.org/downloads/)
- [pip3](https://pip.pypa.io/en/stable/installation/)

The script will offer to install any missing tools automatically on the first run.

## Installation 📥

1. Clone the repository:
   ```bash
   git clone https://github.com/YOUR_USERNAME/nuclei-parmter-fuzz.git
   cd nuclei-parmter-fuzz
   ```
2. Make the script executable:
   ```bash
   chmod +x nuclei-parm-fuzzer.sh
   ```

## Usage 🚀

### Scan a single domain
```bash
./nuclei-parm-fuzzer.sh -d example.com
```

### Scan multiple domains from a file
```bash
./nuclei-parm-fuzzer.sh -f domains.txt
```

### Update all tools and templates
```bash
./nuclei-parm-fuzzer.sh --update
```

## Options ⚙️

| Option | Description |
| :--- | :--- |
| `-d`, `--domain` | Target a single domain |
| `-f`, `--file` | File containing a list of domains |
| `-u`, `--update` | Update all backend tools and Nuclei templates |
| `-h`, `--help` | Show the help message |

## Performance Tuning ⚡

The tool is pre-configured with high-concurrency settings: "
- **Threads (httpx)**: 500
- **Concurrency (Nuclei)**: 50
- **Rate Limit (Nuclei)**: 100 req/s

## Disclaimer ⚠️

This tool is for educational and authorized security testing purposes only. Running it against targets without explicit permission is illegal.

---
**Author**: Ajvad-N
