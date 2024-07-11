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
	args := os.Args

	if len(args) < 2 || len(args) > 2 {
		fmt.Println("You need to add a name for emote")
		fmt.Println("Like " + args[0] + " emotename")
		return
	}

	emoteName := args[1]

	var emotes []Emote
	var input string
	emotesFetched := false

	for {

		if !emotesFetched {
			if !getEmotes(emoteName, &emotes) {
				emoteName = getUserInput("Write new emote name to search for: ")
				continue
			} else {
				emotesFetched = true
			}
		}

		input = getUserInput("Select emote by giving its index or (q to quit | n new search) ")

		if input == "n" {
			emoteName = getUserInput("Write new emote name to search for: ")
			emotesFetched = false
			continue
		}

		if input >= "0" || input <= "9" {
			copyEmoteToClipboard(input, emotes, len(emotes)-1)
		}
	}
}

func getUserInput(infoText string) string {
	fmt.Fprint(os.Stderr, infoText)

	inputReader := bufio.NewReader(os.Stdin)
	input, err := inputReader.ReadString('\n')
	input = strings.TrimSpace(input)

	if err != nil {
		return "Error reading user input"
	}

	if input == "q" {
		fmt.Println("Exiting application ... ")
		os.Exit(0)
	}

	return input
}

func getEmotes(emote_name string, emotes *[]Emote) bool {
	gotAnyEmotes := false
	cmd := exec.Command("node", "fetchEmote.js", emote_name)

	output, err := cmd.Output()
	if err != nil {
		fmt.Println("Error executing command:", err)
		return gotAnyEmotes
	}

	if strings.TrimSpace(bytes.NewBuffer(output).String()) == "No emotes found" {
		fmt.Println("\nNo emotes found for " + emote_name + ", make a new search")
		return false
	}

	err = json.Unmarshal(output, &emotes)

	if err != nil {
		fmt.Println("Error parsing JSON:", err)
		return false
	}

	fmt.Println("Suggested Emotes:")
	for index, emote := range *emotes {
		fmt.Printf("Index: %v Link: https://7tv.app%s  Title: %s\n", index, emote.Href, emote.Title)
	}
	gotAnyEmotes = true

	return gotAnyEmotes
}

func copyEmoteToClipboard(input string, emotes []Emote, lengthOfItemList int) {
	xclipCmd := exec.Command("xclip", "-selection", "clipboard")

	if len(input) > 1 {
		return
	}

	if "a" <= input && input <= "z" || "A" <= input && input <= "Z" {
		return
	}

	userInput, err := strconv.Atoi(input)

	if err != nil {
		fmt.Println("Not a valid selection")
		return
	}

	if userInput >= 0 && userInput <= lengthOfItemList {
		xclipCmd.Stdin = bytes.NewReader([]byte("https://7tv.app" + emotes[userInput].Href))

		if err := xclipCmd.Run(); err != nil {
			fmt.Println("Error copying to user clipboard:", err)
		}

		fmt.Println("Copied emote " + emotes[userInput].Title + " to clipboard")
	} else {
		fmt.Println("Input didn't match any emotes. Please try again.")
	}
}
