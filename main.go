package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"github.com/go-resty/resty/v2"
	"os"
	"strings"
)

func main() {
	fmt.Println("Hello! I'm your chatbot. How can I assist you today?")

	// Create an HTTP client for making API requests
	client := resty.New()
	scanner := bufio.NewScanner(os.Stdin)

	for {
		fmt.Print("You: ")
		scanner.Scan()
		input := scanner.Text()

		// Convert the user input to lowercase for case-insensitive matching
		input = strings.ToLower(input)

		// Use Wit.ai for natural language understanding
		intent, entities, err := getWitIntent(client, input)

		// Based on the intent read by Wit get an answer hard coded in the function
		if err != nil {
			fmt.Println("Chatbot: I'm not sure how to respond to that.")
		} else {
			response := getResponse(intent, entities)
			fmt.Println("Chatbot:", response)
		}
	}
}

// getWitIntent sends the user input to Wit.ai and extracts the detected intent
func getWitIntent(client *resty.Client, input string) (string, string, error) {

	// Set your Wit.ai API token
	const witToken = "VCI3NXCT645IGUORIU3MHET2YLDNXQON"
	const witURL = "https://api.wit.ai/message"

	// Make a GET request to Wit.ai with the user input
	resp, err := client.R().
		SetHeader("Authorization", "Bearer "+witToken).
		SetQueryParams(map[string]string{
			"q": input,
		}).
		Get(witURL)

	if err != nil {
		return "", "", err
	}

	// Parse the JSON response from Wit.ai
	var witResponse map[string]interface{}
	err = json.Unmarshal(resp.Body(), &witResponse)
	if err != nil {
		return "", "", err
	}

	// Extract the detected intent
	intent, ok := witResponse["intents"].([]interface{})
	if !ok || len(intent) == 0 {
		return "", "", nil
	}

	firstIntent, ok := intent[0].(map[string]interface{})
	if !ok {
		return "", "", nil
	}

	// Extract the entities
	entities, ok := witResponse["entities"].(map[string]interface{})
	if ok {
		contact, ok := entities["contact:contact"].([]interface{})
		if ok && len(contact) > 0 {
			contactData, ok := contact[0].(map[string]interface{})
			if ok {
				value, ok := contactData["value"].(string)
				if ok {
					return firstIntent["name"].(string), value, nil
				}
			}
		}
	}

	//fmt.Printf("Wit Response: %+v\n", witResponse)

	return firstIntent["name"].(string), "", nil
}

// getResponse generates a response based on the detected intent
func getResponse(intent string, entities string) string {

	fmt.Printf("Intent: %s\n", intent)
	fmt.Printf("Entities: %s\n", entities)

	// Check for specific entities and customize responses based on them
	if intent == "greet" && entities != "" {
		return "Hello " + entities + "! How can I help you?"
	}

	// Define a response library with predefined responses for different intents
	responseLibrary := map[string]string{
		"greet":     "Hello! How can I help you?",
		"introduce": "I'm a chatbot written in Go.",
		"goodbye":   "Goodbye! Have a great day.",
		"support":   "Please, fill in this form and we'll contact you as soon as we can.",
		"default":   "I'm not sure how to respond to that.",
	}

	// Retrieve the appropriate response from the library or use a default response
	response, ok := responseLibrary[intent]
	if !ok {
		response = responseLibrary["default"]
	}
	return response
}
