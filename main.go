package main

import (
	"fmt"
	"io"
	"net/http"
	"os"
)

func main() {
	// Replace YOUR_TOKEN_HERE with your Discord bot token
	token := os.Getenv("BOT_TOKEN")
	// Replace SERVER_ID_HERE with the ID of the server you want to retrieve channels from
	serverID := os.Getenv("SERVER_ID")

	fmt.Println(token, serverID)
	url := fmt.Sprintf("https://discord.com/api/guilds/%s/channels", serverID)

	client := &http.Client{}

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	req.Header.Set("Authorization", fmt.Sprintf("Bot %s", token))

	resp, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	defer resp.Body.Close()

	// Use io.Copy to copy the response body to stdout
	_, err = io.Copy(os.Stdout, resp.Body)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
