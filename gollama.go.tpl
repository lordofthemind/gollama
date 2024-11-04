package main

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
	Model              string `json:"model"`
	CreatedAt          string `json:"created_at"`
	Response           string `json:"response"`
	Done               bool   `json:"done"`
	Context            []int  `json:"context,omitempty"`
	TotalDuration      int64  `json:"total_duration"`
	LoadDuration       int64  `json:"load_duration"`
	PromptEvalCount    int    `json:"prompt_eval_count"`
	PromptEvalDuration int64  `json:"prompt_eval_duration"`
	EvalCount          int    `json:"eval_count"`
	EvalDuration       int64  `json:"eval_duration"`
}

// Simple completion request
func generateCompletion(model, prompt string) error {
	request := GenerateRequest{
		Model:  model,
		Prompt: prompt,
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

// Streaming completion request
func generateStreamingCompletion(model, prompt string) error {
	request := GenerateRequest{
		Model:  model,
		Prompt: prompt,
		Stream: true,
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

		// Print response chunk without newline
		fmt.Print(response.Response)

		if response.Done {
			fmt.Println("\n\nGeneration completed!")
			fmt.Printf("Total duration: %dms\n", response.TotalDuration/1e6)
			break
		}
	}
	return nil
}

// Chat with context
func chatWithContext(model string) error {
	var context []int

	for {
		// Read user input
		fmt.Print("\nYou: ")
		var input string
		fmt.Scanln(&input)

		if input == "exit" {
			break
		}

		request := GenerateRequest{
			Model:   model,
			Prompt:  input,
			Stream:  true,
			Options: map[string]interface{}{"context": context},
		}

		jsonData, err := json.Marshal(request)
		if err != nil {
			return fmt.Errorf("error marshaling request: %v", err)
		}

		resp, err := http.Post("http://localhost:11434/api/generate", "application/json", bytes.NewBuffer(jsonData))
		if err != nil {
			return fmt.Errorf("error making request: %v", err)
		}

		fmt.Print("Assistant: ")
		decoder := json.NewDecoder(resp.Body)
		for {
			var response GenerateResponse
			if err := decoder.Decode(&response); err != nil {
				if err == io.EOF {
					break
				}
				resp.Body.Close()
				return fmt.Errorf("error decoding stream: %v", err)
			}

			fmt.Print(response.Response)

			if response.Done {
				context = response.Context
				break
			}
		}
		resp.Body.Close()
	}
	return nil
}

func main() {
	// Example usage of simple completion
	fmt.Println("=== Simple Completion ===")
	err := generateCompletion("llama3.2:3b", "Write a haiku about programming")
	if err != nil {
		fmt.Printf("Error: %v\n", err)
	}

	fmt.Println("\n=== Streaming Completion ===")
	err = generateStreamingCompletion("llama3.2:3b", "Explain how garbage collection works in Go")
	if err != nil {
		fmt.Printf("Error: %v\n", err)
	}

	fmt.Println("\n=== Interactive Chat ===")
	fmt.Println("Type 'exit' to end the conversation")
	err = chatWithContext("llama3.2:3b")
	if err != nil {
		fmt.Printf("Error: %v\n", err)
	}
}
