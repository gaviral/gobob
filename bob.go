package main

import (
	"bufio"
	"fmt"
	"gopkg.in/yaml.v2"
	"log"
	"os"
	"os/exec"
	"strings"
)

type CommandInfo struct {
	Type    string   `yaml:"type"`
	Command string   `yaml:"command"`
	Phrases []string `yaml:"phrases"`
}

var commandMap map[string]CommandInfo

func main() {
	println("\nStarting Bob...")
	if err := loadCommands("commands.yaml"); err != nil {
		log.Fatalf("Failed to load commands: %v", err)
	}

	scanner := bufio.NewScanner(os.Stdin)
	for {
		fmt.Print("Enter command phrase: ")
		scanner.Scan()
		userInput := scanner.Text()
		if err := scanner.Err(); err != nil {
			fmt.Fprintln(os.Stderr, "Reading standard input:", err)
			continue
		}

		cmd, found := findCommandByPhrase(userInput)
		if !found {
			fmt.Printf("No command found for phrase: %s\n", userInput)
			continue
		}

		if cmd == "exit_program" {
			println("Exiting Bob...")
			break
		}

		executeCommand(cmd)
	}
}

func loadCommands(filename string) error {
	println("Loading commands from file:", filename)
	data, err := os.ReadFile(filename)
	if err != nil {
		return err
	}
	return yaml.Unmarshal(data, &commandMap)
}

func findCommandByPhrase(phrase string) (string, bool) {
	for cmd, cmdInfo := range commandMap {
		for _, p := range cmdInfo.Phrases {
			if strings.EqualFold(phrase, p) {
				return cmd, true
			}
		}
	}
	return "", false
}

func executeCommand(cmd string) {
	println("Handling command:", cmd)
	details := commandMap[cmd]

	switch details.Type {
	case "cli", "key_press":
		if err := executeSystemCommand(details.Command, details.Type); err != nil {
			fmt.Printf("Command execution error: %v\n", err)
		}
	case "chain":
		for _, part := range strings.Split(details.Command, " | ") {
			executeCommand(part)
		}
	default:
		fmt.Println("Unknown command type:", details.Type)
	}
}

func executeSystemCommand(commandStr, commandType string) error {
	var cmd *exec.Cmd
	if commandType == "cli" {
		println("Executing CLI command:", commandStr)
		cmd = exec.Command("/bin/sh", "-c", commandStr)
	} else if commandType == "key_press" {
		println("Executing key press command:", commandStr)
		cmd = exec.Command("osascript", "-e", commandStr)
	}
	if cmd == nil {
		return fmt.Errorf("invalid command type: %s", commandType)
	}
	return cmd.Run()
}
