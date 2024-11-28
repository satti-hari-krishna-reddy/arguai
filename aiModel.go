package main

import (
    "fmt"
    "strings"

    "github.com/hypermodeinc/modus/sdk/go/pkg/models"
    "github.com/hypermodeinc/modus/sdk/go/pkg/models/openai"
)

func InvokeModel(modelName, instruction, prompt string) (string, error) {
    var err error
    var output any

    if strings.Contains(modelName, "gpt") || strings.Contains(modelName, "Meta") {
        output, err = invokeOpenAIModel(modelName, instruction, prompt)
    } else {
        output, err = invokeGeminiModel(modelName, instruction, prompt)
    }

    if err != nil {
        return "", err
    }

    if result, ok := output.(string); ok {
        return result, nil
    }

    return "", fmt.Errorf("unexpected output type")
}

func invokeOpenAIModel(modelName, instruction, prompt string) (string, error) {
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

func invokeGeminiModel(modelName, instruction, prompt string) (string, error) {
    model, err := models.GetModel[GeminiModel](modelName)
    if err != nil {
        return "", err
    }

    combinedText := instruction
    if prompt != "" {
        combinedText = fmt.Sprintf("%s and \n your task: %s", instruction, prompt)
    }

    input := &GeminiModelInput{
        Contents: []struct {
            Parts []struct {
                Text string `json:"text"`
            } `json:"parts"`
        }{
            {
                Parts: []struct {
                    Text string `json:"text"`
                }{
                    {Text: combinedText},
                },
            },
        },
    }

    output, err := model.Invoke(input)
    if err != nil {
        return "", err
    }

    if len(output.Candidates) > 0 && len(output.Candidates[0].Content.Parts) > 0 {
        return output.Candidates[0].Content.Parts[0].Text, nil
    }

    return "", fmt.Errorf("no valid text found in the response")
}
