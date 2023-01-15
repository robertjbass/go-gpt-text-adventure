package main

import (
	"bufio"
	"context"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/PullRequestInc/go-gpt3"
	"github.com/joho/godotenv"
)

type Client struct {
	apiKey     string
	history    []string
	gpt3Client gpt3.Client
}

func NewClient(apiKey string) *Client {
	client := gpt3.NewClient(apiKey)
	return &Client{apiKey: apiKey, gpt3Client: client}
}

func (c *Client) GetCompletion(prompt string) (string, error) {
	ctx := context.Background()

	if prompt == "" {

		c.history = []string{
			"CONTEXT:",
			"\n",
			"\n",
			"The following is a conversation with an AI dungeon master like in D&D. The bot is using GPT-3 to generate responses. You can talk to the bot freely and it should respond to anything you say. If the player is going along with the story in a reasonable way, the DM will allow it, the DM will push-back against outrageous requests, everything must be mildly realistic. The DM is trying to tell a coherent story. Unlike DND, this game takes place in the modern day.",
			"\n",
			"\n",
			"STORY LINE:",
			"\n",
			"\n",
			os.Getenv("STORY_LINE"),
			"\n",
			"\n",
			"WIN CONDITION:",
			"\n",
			"\n",
			os.Getenv("WIN_CONDITION"),
			"\n",
			"\n",
			"EXAMPLES:",
			"(this is an example of a back and forth conversation, this is not part of the story - this is for example only)",
			"\n",
			"\n",
			"Player: Go to the bathroom",
			"DM: You walk into the bathroom and see a mirror. You look in the mirror and realize that you're wearing lipstick.",
			"\n",
			"\n",
			"Player: Wipe off the lipstick",
			"DM: You wipe off the lipstick and realize that you're wearing a dress.",
			"\n",
			"\n",
			"DIALOGUE:",
			"(This is the actual story)",
			"\n",
			"\n",
			"DM: " + os.Getenv("OPENING_DIALOGUE")}

		fmt.Println(c.history[len(c.history)-1])

		return "", nil
	} else {
		c.history = append(c.history, prompt)

		promptWithHistoryString := strings.Join(c.history, "\n\n") + "DM:"

		completionRequest := gpt3.CompletionRequest{
			Prompt:    []string{promptWithHistoryString},
			MaxTokens: gpt3.IntPtr(1000),
			Stop:      []string{"Player: ", "DM: ", "\n\n"},
			Echo:      false,
		}

		resp, err := c.gpt3Client.Completion(ctx, completionRequest)

		responseText := resp.Choices[0].Text

		c.history = append(c.history, "", "", "", responseText)

		return responseText, err
	}

}

func readInput() string {
	reader := bufio.NewReader(os.Stdin)
	fmt.Print("What would you like to do?\n> ")
	input, _ := reader.ReadString('\n')
	return input
}

func playGame(client *Client) {

	if len(client.history) == 0 {
		client.GetCompletion("")
	} else {

		prompt := readInput()

		resp, err := client.GetCompletion(prompt)
		fmt.Println("DM: ", resp)

		if err != nil {
			log.Fatalln(err)
		}

	}

}

func main() {
	godotenv.Load()

	apiKey := os.Getenv("OPENAI_KEY")
	if apiKey == "" {
		log.Fatalln("Missing API KEY")
	}

	client := NewClient(apiKey)

	for {
		playGame(client)
	}

}
