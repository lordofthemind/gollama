package services

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type GenerateRequest struct {
	Model       string                 `json:"model"`
	Prompt      string                 `json:"prompt"`
	Stream      bool                   `json:"stream"`
	Temperature float64                `json:"temperature,omitempty"`
	Options     map[string]interface{} `json:"options,omitempty"`
}

type GenerateResponse struct {
	Model    string `json:"model"`
	Response string `json:"response"`
	Done     bool   `json:"done"`
	Context  []int  `json:"context,omitempty"`
}

// GenerateCompletion sends a non-streaming request
func GenerateCompletion(model, prompt string, temperature float64) error {
	request := GenerateRequest{
		Model:       model,
		Prompt:      prompt,
		Temperature: temperature,
	}

	jsonData, err := json.Marshal(request)
	if err != nil {
		return fmt.Errorf("error marshaling request: %v", err)
	}

	resp, err := http.Post("http://localhost:11434/api/generate", "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		return fmt.Errorf("error making request: %v", err)
	}
	defer resp.Body.Close()

	var response GenerateResponse
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return fmt.Errorf("error decoding response: %v", err)
	}

	fmt.Printf("Model: %s\nResponse: %s\n", response.Model, response.Response)
	return nil
}

// GenerateStreamingCompletion sends a streaming request
func GenerateStreamingCompletion(model, prompt string, temperature float64) error {
	request := GenerateRequest{
		Model:       model,
		Prompt:      prompt,
		Stream:      true,
		Temperature: temperature,
	}

	jsonData, err := json.Marshal(request)
	if err != nil {
		return fmt.Errorf("error marshaling request: %v", err)
	}

	resp, err := http.Post("http://localhost:11434/api/generate", "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		return fmt.Errorf("error making request: %v", err)
	}
	defer resp.Body.Close()

	decoder := json.NewDecoder(resp.Body)
	for {
		var response GenerateResponse
		if err := decoder.Decode(&response); err != nil {
			if err == io.EOF {
				break
			}
			return fmt.Errorf("error decoding stream: %v", err)
		}

		fmt.Print(response.Response)

		if response.Done {
			fmt.Println("\nGeneration completed!")
			break
		}
	}
	return nil
}
