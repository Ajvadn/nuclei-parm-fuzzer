#!/bin/bash

# parm-fuzzer v2.0 Enterprise Setup 🚀
# Redesigned for maximum performance and security.

set -e

# Colors
CYAN='\033[0;36m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m'

echo -e "${CYAN}Installing Enterprise Discovery Stack...${NC}"

# Install Go Tools
go_tools=(
    "github.com/lc/gau/v2/cmd/gau@latest"
    "github.com/tomnomnom/waybackurls@latest"
    "github.com/projectdiscovery/katana/cmd/katana@latest"
    "github.com/projectdiscovery/httpx/cmd/httpx@latest"
    "github.com/hakluke/hakrawler@latest"
    "github.com/lc/subjs@latest"
    "github.com/jaeles-project/gospider@latest"
)

for tool in "${go_tools[@]}"; do
    echo -e "${GREEN}[+] Installing $tool...${NC}"
    go install "$tool"
done

# Install Python Tools
python_tools=("uro" "waymore")
for tool in "${python_tools[@]}"; do
    echo -e "${GREEN}[+] Installing $tool...${NC}"
    pip3 install "$tool" --break-system-packages --upgrade
done

# Build the Enterprise binary
echo -e "${CYAN}Building parm-fuzzer v2.0 Enterprise...${NC}"
go build -o parm-fuzzer cmd/parm-fuzzer/main.go

echo -e "${GREEN}Done! Run ./parm-fuzzer -d example.com to start.${NC}"
