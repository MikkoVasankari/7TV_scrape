package main

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
)

type Emote struct {
	Href  string `json:"href"`
	Title string `json:"title"`
}

func main() {

	// When clicking Href it should copy item to user clipboard

	args := os.Args[1:]

	if len(args) < 1 {
		fmt.Printf("You need to add a name for emote")
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
	for _, emote := range emotes {
		fmt.Printf("Href: https://7tv.app%s  Title: %s\n", emote.Href, emote.Title)
	}
}
