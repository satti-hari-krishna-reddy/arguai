package main

import (
	"fmt"
	"time"
	"github.com/hypermodeinc/modus/sdk/go/pkg/postgresql"
)

const connection = "convdb" 

type DebateExchange struct {
	PartitionKey string    `json:"debate_id"`
	MessageId       int     `json:"message_id"`
	Speaker      string    `json:"speaker"`
	Message      string    `json:"message"`
	Timestamp    time.Time `json:"timestamp"`
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

func GetDebateHistory(debateID string) ([]DebateExchange) {
	const query = `
		SELECT debate_id, message_id, speaker, message, timestamp
		FROM debate_messages
		WHERE debate_id = $1
		ORDER BY timestamp
	`

	rows, _, err := postgresql.Query[DebateExchange](connection, query, debateID)
	if err != nil {
		return rows
	}

	return rows
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