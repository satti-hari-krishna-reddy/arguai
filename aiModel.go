package main

import (

    "strings"

    "github.com/hypermodeinc/modus/sdk/go/pkg/models"
    "github.com/hypermodeinc/modus/sdk/go/pkg/models/openai"
)


func InvokeModel(modelName, instruction, prompt string) (string, error) {

	// if strings.Contains(modelName, "gpt") || strings.Contains(modelName, "Meta") {
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

    input.Temperature = 0.7

    output, err := model.Invoke(input)
    if err != nil {
        return "", err
    }

    return strings.TrimSpace(output.Choices[0].Message.Content), nil
}
//}