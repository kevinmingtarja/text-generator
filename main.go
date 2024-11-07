package main

import (
	"strings"

	"github.com/hypermodeinc/modus/sdk/go/pkg/models"
	"github.com/hypermodeinc/modus/sdk/go/pkg/models/openai"
)

// this model name should match the one defined in the modus.json manifest file
const modelName = "text-generator"

func GenerateText(instruction, prompt string) (string, error) {
	model, err := models.GetModel[openai.ChatModel](modelName)
	if err != nil {
		return "", err
	}

	input, err := model.CreateInput(
		openai.NewSystemMessage(instruction),
		openai.NewUserMessage(prompt),
	)
	if err != nil {
		return "", err
	}

	// this is one of many optional parameters available for the OpenAI chat interface
	input.Temperature = 0.7

	output, err := model.Invoke(input)
	if err != nil {
		return "", err
	}

	return strings.TrimSpace(output.Choices[0].Message.Content), nil
}
