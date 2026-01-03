package main

import (
	"fmt"
)

const (
	Red    = "\033[91m"
	Green  = "\033[92m"
	Reset  = "\033[0m"
	Yellow = "\033[93m"
	Cyan   = "\033[96m"
)

func PrintBanner() {
	fmt.Printf("%s", Red)
	fmt.Println(" _  _ _   _  ___ _    ___ ___   ___  _   ___ __  __   ___ _   _ ___________ ___ ")
	fmt.Println("| \\| | | | |/ __| |  | __|_ _| | _ \\/_\\ | _ \\  \\/  | | __| | | |_  /_  / __| _ \\")
	fmt.Println("| .` | |_| | (__| |__| _| | |  |  _/ _ \\|   / |\\/| | | _|| |_| |/ / / /| _||   /")
	fmt.Println("|_|\\_|\\___/ \\___|____|___|___| |_|/_/ \\_\\_|\\_\\_|  |_| |_|  \\___//___/___|___|_|_\\")
	fmt.Println("                                                                                ")
	fmt.Println("                                               by Ajvad-N")
	fmt.Printf("%s\n", Reset)
}

func LogInfo(message string) {
	fmt.Printf("%s[INFO] %s%s\n", Green, message, Reset)
}

func LogError(message string) {
	fmt.Printf("%s[ERROR] %s%s\n", Red, message, Reset)
}

func LogStatus(message string) {
	fmt.Printf("%s[+] %s%s\n", Green, message, Reset)
}
