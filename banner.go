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
	fmt.Println(" ___   _   ___ __  __   ___ _   _ __________ ___ ")
	fmt.Println("| _ \\ /_\\ | _ \\  \\/  | | __| | | |_  /_  / __| _ \\")
	fmt.Println("|  _/ _ \\|   / |\\/| | | _|| |_| |/ / / /| _||   /")
	fmt.Println("|_|/_/ \\_\\_|\\_\\_|  |_| |_|  \\___//___/___|___|_|_\\")
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
