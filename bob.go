package main

import (
	"fmt"
	"os/exec"
	"strings"
)

var commandMap = map[string]map[string]string{
	"dictate_text": {
		"type":    "cli",
		"command": "shortcuts run dictate_text",
	},
	"open_chatGPT": {
		"type":    "cli",
		"command": "open -a 'Google Chrome' https://chat.openai.com/",
	},
	"paste_text": {
		"type":    "key_press",
		"command": "cmd+v",
	},
	"press_enter": {
		"type":    "key_press",
		"command": "enter",
	},
	"query_chatGPT": {
		"type":    "chain",
		"command": "dictate_text | open_chatGPT | paste_text | press_enter",
	},
	"dictate_text_and_paste": {
		"type":    "chain",
		"command": "dictate_text | paste_text",
	},
}

func main() {
	var cmd = "dictate_text_and_paste"
	handleCommand(cmd)

}

func handleCommand(cmd string) {
	// switch based on the command type
	switch commandMap[cmd]["type"] {
	case "cli":
		executeCLICommand(commandMap[cmd]["command"])
	case "key_press":
		executeKeyPressCommand(commandMap[cmd]["command"])
	case "chain":
		commands := strings.Split(commandMap[cmd]["command"], " | ")
		for _, c := range commands {
			handleCommand(c)
		}
	}
}

func executeCLICommand(s string) {
	var cmdStrings = strings.Fields(s)
	cmd := exec.Command(cmdStrings[0], cmdStrings[1:]...)
	err := cmd.Run()
	if err != nil {
		fmt.Println("Error: ", err)
		return
	}
}

func executeKeyPressCommand(s string) {
	var cmd *exec.Cmd

	// Handling commands with "cmd" modifier key
	if strings.HasPrefix(s, "cmd+") {
		key := strings.TrimPrefix(s, "cmd+")
		script := fmt.Sprintf("tell application \"System Events\" to keystroke \"%s\" using command down", key)
		cmd = exec.Command("osascript", "-e", script)
	} else {
		// Handling other types of keystrokes like "enter"
		switch s {
		case "enter":
			cmd = exec.Command("osascript", "-e", "tell application \"System Events\" to keystroke return")
		default:
			fmt.Println("Unsupported key press command:", s)
			return
		}
	}

	// Execute the command
	err := cmd.Run()
	if err != nil {
		fmt.Println("Error executing key press command:", err)
	}
}
