package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/joho/godotenv"
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
	err := godotenv.Load()
	if err != nil {
		panic("Error loading .env file")
	}

	api := DiscordApi{
		botToken: os.Getenv("BOT_TOKEN"),
		serverId: os.Getenv("SERVER_ID"),
	}

	method := flag.String("method", "", "create or delete")
	numberOfTables := flag.Int("number", 1, "number of tables to create or delete")

	flag.Parse()

	if *method == "create" {
		api.createCoffeeTables(*numberOfTables)
	} else if *method == "delete" {
		api.deleteCoffeeTables(*numberOfTables)
	} else {
		panic("Invalid method")
	}
}

func (d *DiscordApi) createCoffeeTables(numberOfTables int) {

	for i := 0; i < numberOfTables; i++ {
		body := []byte(fmt.Sprintf(`{"name":"coffee-table","type":2,"user_limit":2,"parent_id":"%s"}`, TablesCategoryId))

		url := fmt.Sprintf("https://discord.com/api/guilds/%s/channels", d.serverId)

		statusCode, responseBody := d.sendRequest("POST", url, body)
		if statusCode != http.StatusCreated {
			panic("Error creating channel: " + string(responseBody))
		}

		fmt.Println(string(responseBody))
	}
}

func (d *DiscordApi) deleteCoffeeTables(numberOfTables int) {
	tableIds := d.getListOfCoffeeTableIds()

	if len(tableIds) < numberOfTables {
		panic("Not enough tables to delete")
	}

	for i := numberOfTables; i > 0; i-- {
		url := fmt.Sprintf("https://discord.com/api/channels/%s", tableIds[i])

		statusCode, responseBody := d.sendRequest("DELETE", url, nil)
		if statusCode != http.StatusOK {
			panic("Error deleting channel with id: " + tableIds[i] + "error: " + string(responseBody))
		}

		fmt.Println(string(responseBody))
	}
}

func (d *DiscordApi) getListOfCoffeeTableIds() []string {
	url := fmt.Sprintf("https://discord.com/api/guilds/%s/channels", d.serverId)

	statusCode, responseBody := d.sendRequest("GET", url, nil)
	if statusCode != http.StatusOK {
		panic("Error getting channels: " + string(responseBody))
	}

	var co []ChannelOrCategory

	err := json.Unmarshal(responseBody, &co)

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

func (d *DiscordApi) sendRequest(method string, url string, requestBody []byte) (int, []byte) {
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

	return response.StatusCode, responseBody
}
