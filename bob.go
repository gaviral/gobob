package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"os/exec"
	"strings"
)

// CommandInfo defines the structure for each command based on the JSON structure
type CommandInfo struct {
	Type    string   `json:"type"`
	Command string   `json:"command"`
	Phrases []string `json:"phrases"`
}

var commandMap map[string]CommandInfo

func main() {
	println("\n\nStarting Bob...")
	loadCommands("commands.json")
	// Continuously listen for commands
	for {
		var cmd string
		fmt.Print("Enter command: ")
		fmt.Scanln(&cmd)
		handleCommand(cmd)
	}
}

func loadCommands(filename string) {
	println("Loading commands from file:", filename)
	data, err := os.ReadFile(filename)
	if err != nil {
		log.Fatal("Error loading commands:", err)
	}
	if err := json.Unmarshal(data, &commandMap); err != nil {
		log.Fatalf("Error parsing commands: %v", err)
	}
}

func handleCommand(cmd string) {
	println("Handling command:", cmd)
	commandDetails := commandMap[cmd]
	switch commandDetails.Type {
	case "cli":
		executeCLICommand(commandDetails.Command)
	case "key_press":
		executeKeyPressCommand(commandDetails.Command)
	case "chain":
		commands := strings.Split(commandDetails.Command, " | ")
		for _, c := range commands {
			handleCommand(c)
		}
	}
}

func executeCLICommand(s string) {
	println("Executing CLI command:", s)
	cmd := exec.Command(strings.Fields(s)[0], strings.Fields(s)[1:]...)
	if err := cmd.Run(); err != nil {
		fmt.Printf("Error: %v\n", err)
	}
}

func executeKeyPressCommand(s string) {
	println("Executing key press command:", s)
	script := "tell application \"System Events\" to keystroke "
	switch {
	case strings.HasPrefix(s, "cmd+"):
		key := strings.TrimPrefix(s, "cmd+")
		script += fmt.Sprintf("\"%s\" using command down", key)
	case s == "enter":
		script += "return"
	default:
		fmt.Printf("Unsupported key press command: %s\n", s)
		return
	}
	if err := exec.Command("osascript", "-e", script).Run(); err != nil {
		fmt.Printf("Error executing key press command: %v\n", err)
	}
}
