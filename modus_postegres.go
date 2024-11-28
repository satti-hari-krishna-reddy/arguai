package main

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/hypermodeinc/modus/sdk/go/pkg/postgresql"
)

const connection = "convdb" 

type DebateExchange struct {
	PartitionKey string    `json:"debate_id"`
	MessageId       int     `json:"message_id"`
	Timestamp    time.Time `json:"timestamp"`
	Speaker      string    `json:"speaker"`
	Message      string    `json:"message"`
}

func StoreExchange(debateID string, exchange DebateExchange) error {
	const query = `
		INSERT INTO debate_messages (debate_id, speaker, message, timestamp)
		VALUES ($1, $2, $3, $4)
	`

	_, err := postgresql.Execute(
		connection,
		query,
		debateID,          
		exchange.Speaker,  
		exchange.Message, 
		exchange.Timestamp, 
	)
	if err != nil {
		return fmt.Errorf("failed to store exchange: %w", err)
	}

	return nil
}


func GetDebateHistory(debateID string) ([]DebateExchange, error) {
	const query = `
		SELECT debate_id, message_id, speaker, message, timestamp
		FROM debate_messages
		WHERE debate_id = $1
		ORDER BY timestamp
	`

	rows, _, err := postgresql.Query[DebateExchange](connection, query, debateID)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve history: %w", err)
	}

	return rows, nil
}

func GetCurrentConversation(debateID string) (string, error) {
	debateHistory, err := GetDebateHistory(debateID)
	if err != nil {
		return "", err
	}

	if len(debateHistory) == 0 {
		return fmt.Sprintf("No conversation found for debate ID: %s", debateID), nil
	}

	// Serialize the data as a JSON-formatted string
	serialized, err := json.MarshalIndent(debateHistory, "", "  ")
	if err != nil {
		return "", fmt.Errorf("failed to format conversation: %w", err)
	}

	return string(serialized), nil
}

func RegisterDebate(debateID, model1, model2, persona1, persona2, debateTopic string) error {
	const query = `
		INSERT INTO debates (debate_id, model1, model2, persona1, persona2, topic)
		VALUES ($1, $2, $3, $4, $5, $6)
	`
	_, err := postgresql.Execute(
		connection,
		query,
		debateID,
		model1,
		model2,
		persona1,
		persona2,
		debateTopic,
	)
	if err != nil {
		return fmt.Errorf("failed to register debate: %w", err)
	}

	return nil
}