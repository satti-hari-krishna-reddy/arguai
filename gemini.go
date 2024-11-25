package main

import (
	"github.com/hypermodeinc/modus/sdk/go/pkg/models"
)

type GeminiModelInput struct {
	Contents []struct {
		Parts []struct {
			Text string `json:"text"`
		} `json:"parts"`
	} `json:"contents"`
}

type GeminiModelOutput struct {
	Candidates []struct {
		Content struct {
			Parts []struct {
				Text string `json:"text"`
			} `json:"parts"`
		} `json:"content"`
	} `json:"candidates"`
}

type GeminiModel struct {
	geminiModelBase 
	
}

type geminiModelBase = models.ModelBase[GeminiModelInput, GeminiModelOutput]
