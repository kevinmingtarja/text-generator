package main

import (
	"fmt"
	"strings"

	"github.com/hypermodeinc/modus/sdk/go/pkg/models"
	"github.com/hypermodeinc/modus/sdk/go/pkg/models/openai"
)

func GenerateText(modelName, prompt string) (string, error) {
	model, err := models.GetModel[openai.ChatModel](modelName)
	if err != nil {
		return "", err
	}
	model.Debug = true

	messages := []openai.Message{
		openai.NewSystemMessage(`When you receive a tool call response, use the output to format an answer to the orginal user question.

You are a helpful assistant with tool calling capabilities.
		`),
		openai.NewUserMessage(prompt),
	}

	input, err := model.CreateInput(
		messages...,
	)
	if err != nil {
		return "", err
	}

	input.Temperature = 1
	input.Tools = []openai.Tool{
		{
			Type: "function",
			Function: openai.FunctionDefinition{
				Name:        "get_current_conditions",
				Description: "Get the current weather conditions for a specific location",
				Parameters: `{
        "type": "object",
        "properties": {
			"location": {
				"type": "string",
				"description": "The city and state, e.g., San Francisco, CA"
			},
			"unit": {
				"type": "string",
				"enum": ["Celsius", "Fahrenheit"],
				"description": "The temperature unit to use. Infer this from the user's location."
			}
        },
        "required": ["location", "unit"]
    }`,
			},
		},
	}
	input.ToolChoice = openai.ToolChoiceAuto
	input.MaxTokens = 500

	output, err := model.Invoke(input)
	if err != nil {
		return "", err
	}

	// TODO: parse this better
	toolCall := strings.TrimSpace(output.Choices[0].Message.Content)
	fmt.Println(toolCall)

	messages = append(messages, openai.NewAssistantMessage(toolCall))
	// Mock response, imagine this output to be the result of an API call
	messages = append(messages, openai.NewToolMessage(`{"output": "The weather in that city is currently 76°F with a low of 56°F tonight."}`, ""))

	input, err = model.CreateInput(
		messages...,
	)
	if err != nil {
		return "", err
	}

	output, err = model.Invoke(input)
	if err != nil {
		return "", err
	}

	return output.Choices[0].Message.Content, nil
}
