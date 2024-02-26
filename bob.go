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
}

func main() {
	var cmd = "dictate_text"
	// switch based on the command type
	switch commandMap[cmd]["type"] {
	case "cli":
		executeCLICommand(commandMap[cmd]["command"])
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
