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
		var userInput string
		fmt.Print("Enter command phrase: ")
		fmt.Scanln(&userInput)
		cmd := findCommandByPhrase(userInput)
		if cmd != "" {
			handleCommand(cmd)
		} else {
			fmt.Printf("No command found for phrase: %s\n", userInput)
		}
	}
}

// findCommandByPhrase searches for a command by matching the user input with phrases.
func findCommandByPhrase(phrase string) string {
	for cmd, cmdInfo := range commandMap {
		for _, p := range cmdInfo.Phrases {
			if strings.EqualFold(phrase, p) {
				return cmd
			}
		}
	}
	return ""
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
