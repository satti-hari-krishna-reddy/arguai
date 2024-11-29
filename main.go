package main

import (
    "fmt"
	"encoding/json"
)

func GetCurrentConversation(debateID string) (string) {
    result := GetDebateHistory(debateID)
	jsonBytes, err := json.Marshal(result)
    if err != nil {
        return "[]"
    }

    return string(jsonBytes)
}


type DebateResponse struct {
    Success bool   `json:"success"`
    Message string `json:"message"`
}

func StartDebate(debateID, model1, model2, persona1, persona2, debateTopic string) DebateResponse {
    _= RegisterDebate(debateID, model1, model2, persona1, persona2, debateTopic)

    _, err := runDebate(debateID, model1, model2, persona1, persona2, debateTopic)
    if err != nil {
        return DebateResponse{
            Success: false,
            Message: fmt.Sprintf("Error running debate: %v", err),
        }
    }

    return DebateResponse{
        Success: true,
        Message: "Debate completed successfully",
    }
}
