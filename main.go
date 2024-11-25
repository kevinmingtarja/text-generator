package main

import (
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

	input, err := model.CreateInput(
		openai.NewSystemMessage(`When you receive a tool call response, use the output to format an answer to the orginal user question.

You are a helpful assistant with tool calling capabilities.
		`),
		openai.NewUserMessage(prompt),
		// openai.NewAssistantMessage(`{"name": "get_current_conditions", "parameters": {"location": "San Francisco, CA", "unit": "Fahrenheit"}}`),
		// openai.NewUserMessage(`{"output": "Clouds giving way to sun Hi: 76° Tonight: Mainly clear early, then areas of low clouds forming Lo: 56°"}`),
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

	return strings.TrimSpace(output.Choices[0].Message.Content), nil
}
