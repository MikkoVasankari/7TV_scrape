package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

type Emote struct {
	Href  string `json:"href"`
	Title string `json:"title"`
}

func main() {

	args := os.Args[1:]

	if len(args) < 1 {
		fmt.Println("You need to add a name for emote")
		return
	}

	emote_name := args[0]

	cmd := exec.Command("node", "fetchEmote.js", emote_name)

	output, err := cmd.Output()
	if err != nil {
		fmt.Println("Error executing command:", err)
		return
	}

	var emotes []Emote
	err = json.Unmarshal(output, &emotes)
	if err != nil {
		fmt.Println("Error parsing JSON:", err)
		return
	}

	fmt.Println("Suggested Emotes:")
	for index, emote := range emotes {
		fmt.Printf("%v Link: https://7tv.app%s  Title: %s\n", index, emote.Href, emote.Title)
	}

	ReadUserInput("Select emote by giving it's index ", emotes, len(emotes)-1)
}

func ReadUserInput(label string, emotes []Emote, lengthOfItemList int) {
	inputReader := bufio.NewReader(os.Stdin)
	xclipCmd := exec.Command("xclip", "-selection", "clipboard")

	for {
		fmt.Fprintf(os.Stderr, label+"")
		input, err := inputReader.ReadString('\n')

		if err != nil {
			fmt.Println("Error reading user input:", err)
			return
		}

		input = strings.TrimSpace(input)
		userInput, err := strconv.Atoi(input)
		if err != nil {
			fmt.Println("Error converting input to integer:", err)
			return
		}

		if userInput >= 0 && userInput <= lengthOfItemList {
			xclipCmd.Stdin = bytes.NewReader([]byte("https://7tv.app" + emotes[userInput].Href))
			if err := xclipCmd.Run(); err != nil {
				fmt.Println("Error copying to user clipboard:", err)
			}
			fmt.Println("Copied emote " + emotes[userInput].Title + " to clipboard")
			break
		} else {
			fmt.Println("Input didn't match any emotes. Please try again.")
		}
	}
}
