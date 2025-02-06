package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/cookiejar"
	"os"
	"strings"

	"blackjack/api"
	"blackjack/client/config"
)

var helpMessage = "Available actions\nNew: Start a new game. Forfeits the current game if one is active.\nHit: Take another card from the deck.\nStand: End the current game and score the player's hand.\nGet: Returns the current game if one is active.\nDelete: End the current game session.\nHelp: Display this message again.\nExit: Close the game.\n\n"

func main() {
	config, err := config.LoadConfig()
	if err != nil {
		panic(err)
	}

	address := config.ServerAddress
	if config.Port != "" {
		address += ":" + config.Port
	}

	cookieJar, err := cookiejar.New(nil)
	if err != nil {
		panic(err)
	}

	client := http.Client{
		Jar: cookieJar,
	}

	reader := bufio.NewReader(os.Stdin)
	fmt.Print("Welcome to Blackjack!\n\n")
	fmt.Print(helpMessage)

	for {
		fmt.Print("Select action: ")
		action, err := reader.ReadString('\n')
		if err != nil {
			panic(err)
		}

		action = strings.ToLower(strings.TrimSpace(action))
		switch action {
		case "get":
			resp, err := get(&client, address)
			if err != nil {
				fmt.Printf("Failed to get: %s\n", err.Error())
				continue
			}
			fmt.Printf("Dealer hand: %s\nPlayer hand: %s\nWins: %d\nLosses: %d\nTies: %d\n\n", resp.DealerHand, resp.PlayerHand, resp.Wins, resp.Losses, resp.Ties)

		case "delete":
			err = delete(&client, address)
			if err != nil {
				fmt.Printf("Failed to delete: %s\n", err.Error())
				continue
			}
			fmt.Printf("The current game session was deleted\n\n")

		case "new":
			resp, err := upsert(&client, address, api.UpsertCreate)
			if err != nil {
				fmt.Printf("Failed to create new game: %s\n", err.Error())
				continue
			}
			fmt.Printf("Dealer hand: %s\nPlayer hand: %s\nWins: %d\nLosses: %d\nTies: %d\n\n", resp.DealerHand, resp.PlayerHand, resp.Wins, resp.Losses, resp.Ties)

		case "hit":
			resp, err := upsert(&client, address, api.UpsertHit)
			if err != nil {
				fmt.Printf("Failed to hit: %s\n", err.Error())
				continue
			}
			fmt.Printf("Dealer hand: %s\nPlayer hand: %s\nWins: %d\nLosses: %d\nTies: %d\n\n", resp.DealerHand, resp.PlayerHand, resp.Wins, resp.Losses, resp.Ties)

		case "stand":
			resp, err := upsert(&client, address, api.UpsertStand)
			if err != nil {
				fmt.Printf("Failed to stand: %s\n", err.Error())
				continue
			}
			fmt.Printf("Dealer hand: %s\nPlayer hand: %s\nWins: %d\nLosses: %d\nTies: %d\n\n", resp.DealerHand, resp.PlayerHand, resp.Wins, resp.Losses, resp.Ties)

		case "help":
			fmt.Print(helpMessage)

		case "exit":
			fmt.Println("Thanks for playing!")
			os.Exit(0)

		default:
			fmt.Printf("Unknown action '%s'\n\n", action)
		}
	}
}

func get(client *http.Client, serverAddress string) (*api.GameStateResponse, error) {
	return doRequestWithMarshalledResponse(client, http.MethodGet, serverAddress+"/game", nil, http.StatusOK)
}

func delete(client *http.Client, serverAddress string) error {
	_, err := doRequest(client, http.MethodDelete, serverAddress+"/game", nil, http.StatusNoContent)
	return err
}

func upsert(client *http.Client, serverAddress string, action api.UpsertAction) (*api.GameStateResponse, error) {
	body := api.UpsertGameRequestBody{Action: action}
	bodyBytes, err := json.Marshal(body)
	if err != nil {
		return nil, err
	}

	bodyReader := bytes.NewReader(bodyBytes)
	return doRequestWithMarshalledResponse(client, http.MethodPost, serverAddress+"/game", bodyReader, http.StatusOK)
}

func doRequest(client *http.Client, method, address string, body io.Reader, expectedStatus int) ([]byte, error) {
	req, err := http.NewRequest(method, address, body)
	if err != nil {
		return nil, err
	}

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	respBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != expectedStatus {
		err = fmt.Errorf(string(respBytes))
	}

	return respBytes, err
}

func doRequestWithMarshalledResponse(client *http.Client, method, address string, body io.Reader, expectedStatus int) (*api.GameStateResponse, error) {
	respBytes, err := doRequest(client, method, address, body, expectedStatus)
	if err != nil {
		return nil, err
	}

	var gameStateResponse api.GameStateResponse
	err = json.Unmarshal(respBytes, &gameStateResponse)
	if err != nil {
		return nil, err
	}

	return &gameStateResponse, nil
}
