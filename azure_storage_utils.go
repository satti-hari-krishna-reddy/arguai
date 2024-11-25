package main

import (
	// "context"
	// "encoding/json"
	// "fmt"
	// "strings"
	"time"

	// "github.com/Azure/azure-sdk-for-go/sdk/data/aztables"
)

type DebateExchange struct {
	PartitionKey string
	RowKey       string
	Timestamp    time.Time
	Speaker      string
	Message      string
}

// type DebateStore struct {
// 	client *aztables.Client
// }

// func NewDebateStore(accountName, accountKey, tableName string) (*DebateStore, error) {
// 	serviceURL := fmt.Sprintf("https://%s.table.core.windows.net", accountName)
// 	cred, err := aztables.NewSharedKeyCredential(accountName, accountKey)
// 	if err != nil {
// 		return nil, err
// 	}

// 	client, err := aztables.NewClientWithSharedKey(serviceURL, cred, nil)
// 	if err != nil {
// 		return nil, err
// 	}

// 	_, err = client.CreateTable(context.Background(), &aztables.CreateTableOptions{})
// 	if err != nil {
// 		return nil, err
// 	}

// 	return &DebateStore{client: client}, nil
// }

// func (ds *DebateStore) StoreExchange(debateID, speaker, message string) error {
// 	entity := aztables.EDMEntity{
// 		Entity: aztables.Entity{
// 			PartitionKey: debateID,
// 			RowKey:       fmt.Sprintf("%d", time.Now().UnixNano()),
// 		},
// 		Properties: map[string]interface{}{
// 			"Speaker":   speaker,
// 			"Message":   message,
// 			"Timestamp": time.Now().Format(time.RFC3339),
// 		},
// 	}
// 	entityBytes, err := json.Marshal(entity)
// 	if err != nil {
// 		return err
// 	}
// 	_, err = ds.client.AddEntity(context.Background(), entityBytes, &aztables.AddEntityOptions{})
// 	return err
// }

// func (ds *DebateStore) GetDebateHistory(debateID string) ([]DebateExchange, error) {
// 	var exchanges []DebateExchange

// 	filter := fmt.Sprintf("PartitionKey eq '%s'", debateID)
// 	options := &aztables.ListEntitiesOptions{
// 		Filter: &filter,
// 	}

// 	pager := ds.client.NewListEntitiesPager(options)
// 	for pager.More() {
// 		page, err := pager.NextPage(context.Background())
// 		if err != nil {
// 			return nil, err
// 		}
// 		for _, entity := range page.Entities {
// 			var rawEntity map[string]interface{}
// 			if err := json.Unmarshal(entity, &rawEntity); err != nil {
// 				return nil, err
// 			}

// 			exchange := DebateExchange{
// 				PartitionKey: rawEntity["PartitionKey"].(string),
// 				RowKey:       rawEntity["RowKey"].(string),
// 				Speaker:      rawEntity["Speaker"].(string),
// 				Message:      rawEntity["Message"].(string),
// 				Timestamp:    parseTime(rawEntity["Timestamp"].(string)),
// 			}
// 			exchanges = append(exchanges, exchange)
// 		}
// 	}
// 	return exchanges, nil
// }

// func (ds *DebateStore) GetCurrentConversation(debateID string) (string, error) {
	
// 	debateHistory, err := ds.GetDebateHistory(debateID)
// 	if err != nil {
// 		return "", err
// 	}

// 	if len(debateHistory) == 0 {
// 		return fmt.Sprintf("No conversation found for debate ID: %s", debateID), nil
// 	}

// 	conversationStr := "[\n"
// 	for _, exchange := range debateHistory {
// 		conversationStr += fmt.Sprintf(
// 			"  {\n    \"PartitionKey\": \"%s\",\n    \"RowKey\": \"%s\",\n    \"Timestamp\": \"%s\",\n    \"Speaker\": \"%s\",\n    \"Message\": \"%s\"\n  },\n",
// 			exchange.PartitionKey, exchange.RowKey, exchange.Timestamp.Format(time.RFC3339), exchange.Speaker, exchange.Message,
// 		)
// 	}
// 	conversationStr = strings.TrimSuffix(conversationStr, ",\n") + "\n]"

// 	return conversationStr, nil
// }


// func parseTime(timeStr string) time.Time {
// 	t, err := time.Parse(time.RFC3339, timeStr)
// 	if err != nil {
// 		return time.Now() 
// 	}
// 	return t
// }
