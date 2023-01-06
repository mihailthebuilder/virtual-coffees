package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
)

func main() {
	botToken := os.Getenv("BOT_TOKEN")
	serverId := os.Getenv("SERVER_ID")

	// createCoffeeTables(botToken, serverId, 5)
	deleteCoffeeTables(botToken, serverId, 5)
}

func createCoffeeTables(botToken string, serverId string, numberOfTables int) {

	for i := 0; i < numberOfTables; i++ {
		body := []byte(`{"name":"coffee-table","type":2,"user_limit":2,"parent_id":"1035886233135620127"}`)

		url := fmt.Sprintf("https://discord.com/api/guilds/%s/channels", serverId)

		client := &http.Client{}

		req, err := http.NewRequest("POST", url, bytes.NewBuffer(body))
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		req.Header.Set("Authorization", fmt.Sprintf("Bot %s", botToken))
		req.Header.Set("Content-Type", "application/json")

		resp, err := client.Do(req)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		defer resp.Body.Close()

		_, err = io.Copy(os.Stdout, resp.Body)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	}
}

func deleteCoffeeTables(botToken string, serverId string, numberOfTables int) {
	tableIds := getListOfCoffeeTableIds(botToken, serverId)

	for i := 0; i < numberOfTables; i++ {
		url := fmt.Sprintf("https://discord.com/api/channels/%s", tableIds[i])

		client := &http.Client{}

		req, err := http.NewRequest("DELETE", url, nil)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		req.Header.Set("Authorization", fmt.Sprintf("Bot %s", botToken))
		req.Header.Set("Content-Type", "application/json")

		resp, err := client.Do(req)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		defer resp.Body.Close()

		_, err = io.Copy(os.Stdout, resp.Body)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	}
}

type ChannelOrCategory struct {
	Id       string `json:"id"`
	ParentId string `json:"parent_id"`
}

func getListOfCoffeeTableIds(botToken string, serverId string) []string {
	url := fmt.Sprintf("https://discord.com/api/guilds/%s/channels", serverId)

	client := &http.Client{}

	request, err := http.NewRequest("GET", url, nil)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	request.Header.Set("Authorization", fmt.Sprintf("Bot %s", botToken))
	request.Header.Set("Content-Type", "application/json")

	response, err := client.Do(request)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	defer response.Body.Close()

	var co []ChannelOrCategory

	// resp.Body to string
	buf := new(bytes.Buffer)
	buf.ReadFrom(response.Body)

	json.Unmarshal(buf.Bytes(), &co)

	var out []string

	for _, v := range co {
		if v.ParentId == "1035886233135620127" {
			out = append(out, v.Id)
		}
	}

	return out
}
