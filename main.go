package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
)

type DiscordApi struct {
	botToken string
	serverId string
}

type ChannelOrCategory struct {
	Id       string `json:"id"`
	ParentId string `json:"parent_id"`
}

const TablesCategoryId = "1035886233135620127"

func main() {
	api := DiscordApi{
		botToken: os.Getenv("BOT_TOKEN"),
		serverId: os.Getenv("SERVER_ID"),
	}

	// createCoffeeTables(botToken, serverId, 5)
	api.deleteCoffeeTables(1)
}

func createCoffeeTables(botToken string, serverId string, numberOfTables int) {

	for i := 0; i < numberOfTables; i++ {
		body := []byte(fmt.Sprintf(`{"name":"coffee-table","type":2,"user_limit":2,"parent_id":"%s"}`, TablesCategoryId))

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

func (d *DiscordApi) deleteCoffeeTables(numberOfTables int) {
	tableIds := d.getListOfCoffeeTableIds()

	for i := 0; i < numberOfTables; i++ {
		url := fmt.Sprintf("https://discord.com/api/channels/%s", tableIds[i])

		response := d.sendRequest("DELETE", url, nil)
		fmt.Println(string(response))
	}
}

func (d *DiscordApi) getListOfCoffeeTableIds() []string {
	url := fmt.Sprintf("https://discord.com/api/guilds/%s/channels", d.serverId)

	response := d.sendRequest("GET", url, nil)

	var co []ChannelOrCategory

	err := json.Unmarshal(response, &co)

	if err != nil {
		panic(err)
	}

	var out []string

	for _, v := range co {
		if v.ParentId == TablesCategoryId {
			out = append(out, v.Id)
		}
	}

	return out
}

func (d *DiscordApi) sendRequest(method string, url string, requestBody []byte) []byte {
	client := &http.Client{}

	request, err := http.NewRequest(method, url, bytes.NewBuffer(requestBody))
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	request.Header.Set("Authorization", fmt.Sprintf("Bot %s", d.botToken))
	request.Header.Set("Content-Type", "application/json")

	response, err := client.Do(request)
	if err != nil {
		panic(err)
	}
	defer response.Body.Close()

	responseBody, err := io.ReadAll(response.Body)
	if err != nil {
		panic(err)
	}

	if response.StatusCode != http.StatusOK {
		panic(fmt.Sprintf("Status: %d, body: %b", response.StatusCode, responseBody))
	}

	return responseBody
}
