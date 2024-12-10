package ai

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/bwmarrin/discordgo"
)

func FinalMessage(dg *discordgo.Session, channelID string, defcon_countdown_days int) {

	fmt.Println("Generating AI message")

	message, err := funnyAIMessage()
	if err != nil {
		panic(err)
	}

	fmt.Println("Generated message!")

	dg.ChannelMessageSend(channelID, fmt.Sprintf("Good morning! You're just %d days away from Defcon!\n\n%s", defcon_countdown_days, message))
}

// Response represents the structure of each line of JSON response
type Response struct {
	Model     string `json:"model"`
	CreatedAt string `json:"created_at"`
	Response  string `json:"response"`
	Done      bool   `json:"done"`
}

func funnyAIMessage() (string, error) {
	prompt := "Write a short one liner joke about cyber security."

	response, err := postRequest("http://localhost:11434/api/generate", fmt.Sprintf(`{"model": "llama2", "prompt":"%s"}`, prompt))
	if err != nil {
		fmt.Println("Error:", err)
		return "", err
	}
	arrayOfResponses := strings.Split(response, "\n")

	actualResponse, err := extractResponseText(arrayOfResponses)
	if err != nil {
		fmt.Println(err)
	}

	return actualResponse, nil
}

func postRequest(url string, data string) (string, error) {
	// Create a new buffer with the data for the POST request
	buffer := bytes.NewBufferString(data)

	// Create a POST request
	req, err := http.NewRequest("POST", url, buffer)
	if err != nil {
		return "", err
	}

	// Set the appropriate headers (Content-Type in this case)
	req.Header.Set("Content-Type", "application/json")

	// Make the request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	// Read the response body
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	return string(body), nil
}

// ExtractResponseText extracts the response text from a series of JSON strings.
func extractResponseText(jsonResponses []string) (string, error) {
	var responseText strings.Builder

	for _, jsonResponse := range jsonResponses {
		var resp Response
		if err := json.Unmarshal([]byte(jsonResponse), &resp); err != nil {
			return "", err
		}

		responseText.WriteString(resp.Response)

		if resp.Done {
			break
		}
	}

	return responseText.String(), nil
}
