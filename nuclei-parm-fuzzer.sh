#!/bin/bash

# ANSI color codes
RED='\033[91m'
GREEN='\033[92m'
RESET='\033[0m'

# ASCII art banner
echo -e "${RED}"
cat << "EOF"
  _   _ ___ ___ _    ___ ___   ___  ___ ___ __  __   ___ _   _ _______________ 
 | | | | __| __| |  | __|_ _| | _ \/ __| _ \  \/  | | __| | | |_  /_  / __| _ \
 | |_| | _|| _|| |__| _| | |  |  _/ (__|   / |\/| | | _|| |_| |/ / / /| _||   /
  \___/|___|___|____|___|___| |_|  \___|_|_\_|  |_| |_|  \___//___/___|___|_|_\
                                                                               
                                               by Ajvad-N
EOF
echo -e "${RESET}"

# Add Go bin and local bin to PATH
export PATH=$PATH:$(go env GOPATH)/bin:$HOME/.local/bin

# Function to display help
usage() {
    echo -e "Usage: $0 [options]"
    echo -e ""
    echo -e "Options:"
    echo -e "  -d, --domain <domain>    Target single domain"
    echo -e "  -f, --file   <file>      File containing list of domains"
    echo -e "  -u, --update             Update all tools and Nuclei templates"
    echo -e "  -h, --help               Show this help message"
    echo -e ""
    echo -e "Example:"
    echo -e "  $0 -d example.com"
    echo -e "  $0 --file domains.txt"
    echo -e "  $0 --update"
    exit 1
}

# Function to check and install tools
check_and_install() {
    local tool=$1
    local install_cmd=$2

    if ! command -v "$tool" &>/dev/null; then
        echo -e "${RED}[ERROR] $tool is not installed.${RESET}"
        read -p "Do you want to install $tool? (y/n): " choice
        if [[ "$choice" == "y" || "$choice" == "Y" ]]; then
            echo -e "${GREEN}[INFO] Installing $tool...${RESET}"
            eval "$install_cmd"
            if command -v "$tool" &>/dev/null; then
                echo -e "${GREEN}[INFO] $tool installed successfully.${RESET}"
            else
                echo -e "${RED}[ERROR] Failed to install $tool. Please install it manually.${RESET}"
                exit 1
            fi
        else
            echo -e "${RED}[ERROR] $tool is required. Exiting.${RESET}"
            exit 1
        fi
    fi
}

# Function to update all tools
update_tools() {
    echo -e "${GREEN}[INFO] Updating all tools and Nuclei templates...${RESET}"
    
    echo -e "${GREEN}[+] Updating gau...${RESET}"
    go install github.com/lc/gau/v2/cmd/gau@latest
    
    echo -e "${GREEN}[+] Updating waybackurls...${RESET}"
    go install github.com/tomnomnom/waybackurls@latest
    
    echo -e "${GREEN}[+] Updating katana...${RESET}"
    go install github.com/projectdiscovery/katana/cmd/katana@latest
    
    echo -e "${GREEN}[+] Updating httpx...${RESET}"
    go install github.com/projectdiscovery/httpx/cmd/httpx@latest
    
    echo -e "${GREEN}[+] Updating nuclei...${RESET}"
    go install github.com/projectdiscovery/nuclei/v3/cmd/nuclei@latest
    
    echo -e "${GREEN}[+] Updating nuclei templates...${RESET}"
    nuclei -ut
    
    echo -e "${GREEN}[+] Updating uro...${RESET}"
    pip3 install --upgrade uro --break-system-packages
    
    echo -e "${GREEN}[+] Updating paramspider...${RESET}"
    pip3 install --upgrade git+https://github.com/devanshbatham/ParamSpider --break-system-packages
    
    echo -e "${GREEN}[INFO] All tools and templates updated successfully!${RESET}"
    exit 0
}

# Parse arguments
TARGET_DOMAIN=""
TARGET_FILE=""

while [[ "$#" -gt 0 ]]; do
    case $1 in
        -d|--domain) TARGET_DOMAIN="$2"; shift ;;
        -f|--file) TARGET_FILE="$2"; shift ;;
        -u|--update) update_tools ;;
        -h|--help) usage ;;
        *) echo -e "${RED}[ERROR] Unknown parameter: $1${RESET}"; usage ;;
    esac
    shift
done

# Validate input
if [[ -z "$TARGET_DOMAIN" && -z "$TARGET_FILE" ]]; then
    echo -e "${RED}[ERROR] You must specify a target domain (-d) or a file (--file).${RESET}"
    usage
fi

if [[ -n "$TARGET_DOMAIN" && -n "$TARGET_FILE" ]]; then
    echo -e "${RED}[ERROR] Please specify either (-d) or (--file), not both.${RESET}"
    usage
fi

# Check and install required tools
check_and_install "gau" "go install github.com/lc/gau/v2/cmd/gau@latest"
check_and_install "waybackurls" "go install github.com/tomnomnom/waybackurls@latest"
check_and_install "katana" "go install github.com/projectdiscovery/katana/cmd/katana@latest"
check_and_install "httpx" "go install github.com/projectdiscovery/httpx/cmd/httpx@latest"
check_and_install "nuclei" "go install github.com/projectdiscovery/nuclei/v3/cmd/nuclei@latest"
check_and_install "uro" "pip3 install uro --break-system-packages"
check_and_install "paramspider" "pip3 install git+https://github.com/devanshbatham/ParamSpider --break-system-packages"

# Create temporary files
ALL_URLS_FILE=$(mktemp)
FILTERED_URLS_FILE="filtered_urls.txt"
NUCLEI_RESULTS="nuclei_results.txt"

# Step 1: Fetch URLs using multiple tools
echo -e "${GREEN}[INFO] Fetching URLs using multiple tools (gau, waybackurls, katana, paramspider) in parallel...${RESET}"

# Create temporary files for parallel output
GAU_OUT=$(mktemp)
WAYBACK_OUT=$(mktemp)
KATANA_OUT=$(mktemp)
PARAMSPIDER_OUT=$(mktemp)

if [[ -n "$TARGET_FILE" ]]; then
    echo -e "${GREEN}[INFO] Processing list from file: $TARGET_FILE${RESET}"
    
    gau --subs < "$TARGET_FILE" > "$GAU_OUT" &
    waybackurls < "$TARGET_FILE" > "$WAYBACK_OUT" &
    katana -list "$TARGET_FILE" -d 5 -silent -jc -concurrency 50 -timeout 10 > "$KATANA_OUT" &
    paramspider -l "$TARGET_FILE" -s > "$PARAMSPIDER_OUT" &
else
    echo -e "${GREEN}[INFO] Processing single domain: $TARGET_DOMAIN${RESET}"
    
    gau "$TARGET_DOMAIN" --subs > "$GAU_OUT" &
    echo "$TARGET_DOMAIN" | waybackurls > "$WAYBACK_OUT" &
    katana -u "$TARGET_DOMAIN" -d 5 -silent -jc -concurrency 50 -timeout 10 > "$KATANA_OUT" &
    paramspider -d "$TARGET_DOMAIN" -s > "$PARAMSPIDER_OUT" &
fi

# Wait for all background processes to finish
wait

# Combine results
cat "$GAU_OUT" "$WAYBACK_OUT" "$KATANA_OUT" "$PARAMSPIDER_OUT" >> "$ALL_URLS_FILE"
rm "$GAU_OUT" "$WAYBACK_OUT" "$KATANA_OUT" "$PARAMSPIDER_OUT"

# Step 2: Filter URLs with query parameters
echo -e "${GREEN}[INFO] Filtering URLs with query parameters...${RESET}"
grep -E '\?[^=]+=.+$' "$ALL_URLS_FILE" | uro | sort -u > "$FILTERED_URLS_FILE"
rm "$ALL_URLS_FILE"

# Step 3: Check live URLs using httpx
echo -e "${GREEN}[INFO] Checking for live URLs using httpx (optimized)...${RESET}"
# Increased threads and reduced timeout for speed
httpx -silent -threads 500 -rl 300 -timeout 5 < "$FILTERED_URLS_FILE" > "$FILTERED_URLS_FILE.tmp"
mv "$FILTERED_URLS_FILE.tmp" "$FILTERED_URLS_FILE"

# Step 4: Run nuclei for DAST scanning
echo -e "${GREEN}[INFO] Running nuclei for DAST scanning (optimized)...${RESET}"
# Increased concurrency and rate-limiting
nuclei -dast -retries 2 -silent -concurrency 50 -rate-limit 100 -o "$NUCLEI_RESULTS" < "$FILTERED_URLS_FILE"

# Step 5: Show saved results
echo -e "${GREEN}[INFO] Nuclei results saved to $NUCLEI_RESULTS${RESET}"
echo -e "${GREEN}[INFO] Filtered URLs saved to $FILTERED_URLS_FILE for manual testing.${RESET}"
echo -e "${GREEN}[INFO] Automation completed successfully!${RESET}"

# Check if Nuclei found any vulnerabilities
if [ ! -s "$NUCLEI_RESULTS" ]; then
    echo -e "${GREEN}[INFO] No vulnerable URLs found.${RESET}"
else
    echo -e "${GREEN}[INFO] Vulnerabilities were detected. Check $NUCLEI_RESULTS for details.${RESET}"
fi
