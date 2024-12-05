package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"runtime"
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

	operatingSystem := runtime.GOOS
	var emotes []Emote
	var input string
	emotesFetched := false

	for {

		if !emotesFetched {
			if !getEmotes(emoteName, &emotes, operatingSystem) {
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
			copyEmoteToClipboard(input, emotes, len(emotes)-1, operatingSystem)
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

func getEmotes(emote_name string, emotes *[]Emote, operatingSystem string) bool {
	var pathSeparator string
	
	switch operatingSystem {
	case "linux":
		pathSeparator = "/"
	case "windows":
		pathSeparator = "\\"
	default:
		fmt.Println("unsupported platform")
	}

	wd, err := os.Executable()
	if err != nil {
		fmt.Print("Couldn't get os.Executable()")
	}
	
	wdParsed := wd[:strings.LastIndex(wd, pathSeparator)]

	gotAnyEmotes := false

	cmd := exec.Command("node", wdParsed+pathSeparator+"fetchEmote.js", emote_name)

	output, err := cmd.Output()
	if err != nil {
		fmt.Println("Error executing command:", err)
		return false
	}

	if strings.TrimSpace(bytes.NewBuffer(output).String()) == "No emotes found" {
		fmt.Println("\nNo emotes found for " + emote_name + ", make a new search")
		return gotAnyEmotes
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

func copyEmoteToClipboard(input string, emotes []Emote, lengthOfItemList int, operatingSystem string) {
	var clipboardCmd *exec.Cmd

	switch operatingSystem {
	case "linux":
		clipboardCmd = exec.Command("xclip", "-selection", "clipboard")
	case "windows":
		clipboardCmd = exec.Command("cmd", "/c", "clip")
	default:
		fmt.Println("unsupported platform")
	}

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
		clipboardCmd.Stdin = bytes.NewReader([]byte("https://7tv.app" + emotes[userInput].Href))

		if err := clipboardCmd.Run(); err != nil {
			fmt.Println("Error copying to user clipboard:", err)
		}

		fmt.Println("Copied emote " + emotes[userInput].Title + " to clipboard")
	} else {
		fmt.Println("Input didn't match any emotes. Please try again.")
	}
}
