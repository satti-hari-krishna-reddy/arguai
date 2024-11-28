// package main

// import (
//     "strings"

//     "github.com/hypermodeinc/modus/sdk/go/pkg/models"
//     "github.com/hypermodeinc/modus/sdk/go/pkg/models/openai"
// )

// // func GenerateText(instruction, prompt, modelName string) (string, error) {
// //     model, err := models.GetModel[openai.ChatModel](modelName)
// //     if err != nil {
// //         return "", err
// //     }

// //     input, err := model.CreateInput(
// //         openai.NewSystemMessage(instruction),
// //         openai.NewUserMessage(prompt),
// //     )
// //     if err != nil {
// //         return "", err
// //     }

// //     input.Temperature = 0.7

// //     output, err := model.Invoke(input)
// //     if err != nil {
// //         return "", err
// //     }

// //     return strings.TrimSpace(output.Choices[0].Message.Content), nil
// // }

// func StartDebate(debateTopic, model1, persona1, model2, persona2, debateID string) (string, error) {
//     ps := NewPersonaStore()
//     ds, err := NewDebateStore("modus", "modus", "DebateHistory")
//     if err != nil {
//         return "", err
//     }

//     debateID := "debate_" + time.Now().Format("20060102150405")
//     debateHistory := []DebateExchange{}

//     model1Response, err := initiateDebateRound(ps, model1, persona1, debateTopic, "", "", "role")
//     if err != nil {
//         return "", err
//     }
//     debateHistory = append(debateHistory, DebateExchange{
//         PartitionKey: debateID,
//         RowKey:       "1_" + model1,
//         Timestamp:    time.Now(),
//         Speaker:      model1,
//         Message:      model1Response,
//     })
//     ds.StoreExchange(debateID, debateHistory[len(debateHistory)-1].RowKey, debateHistory[len(debateHistory)-1].Message) // Store role

//     model2Response, err := initiateDebateRound(ps, model2, persona2, debateTopic, "", "", "role")
//     if err != nil {
//         return "", err
//     }
//     debateHistory = append(debateHistory, DebateExchange{
//         PartitionKey: debateID,
//         RowKey:       "2_" + model2,
//         Timestamp:    time.Now(),
//         Speaker:      model2,
//         Message:      model2Response,
//     })
//     ds.StoreExchange(debateID, debateHistory[len(debateHistory)-1].RowKey, debateHistory[len(debateHistory)-1].Message) // Store role

//     for i := 0; i < 4; i++ {
//         currentModel := model1
//         if i%2 == 0 {
//             currentModel = model2
//         }

//         var model1Response, model2Response string
//         if currentModel == model1 {
//             model1Response, err = initiateDebateRound(ps, model1, persona1, debateTopic, model1Response, model2Response, "rebuttal")
//             if err != nil {
//                 return "", err
//             }
//             debateHistory = append(debateHistory, DebateExchange{
//                 PartitionKey: debateID,
//                 RowKey:       fmt.Sprintf("%d_%s", i+3, model1),
//                 Timestamp:    time.Now(),
//                 Speaker:      model1,
//                 Message:      model1Response,
//             })
//             ds.StoreExchange(debateID, debateHistory[len(debateHistory)-1].RowKey, debateHistory[len(debateHistory)-1].Message) // Store rebuttal
//         } else {
//             model2Response, err = initiateDebateRound(ps, model2, persona2, debateTopic, model1Response, model2Response, "rebuttal")
//             if err != nil {
//                 return "", err
//             }
//             debateHistory = append(debateHistory, DebateExchange{
//                 PartitionKey: debateID,
//                 RowKey:       fmt.Sprintf("%d_%s", i+3, model2),
//                 Timestamp:    time.Now(),
//                 Speaker:      model2,
//                 Message:      model2Response,
//             })
//             ds.StoreExchange(debateID, debateHistory[len(debateHistory)-1].RowKey, debateHistory[len(debateHistory)-1].Message) // Store rebuttal
//         }
//     }
package main

import (
    "fmt"
    "strings"
    "time"

)

// Available models and personas
var availableModels = []string{"gpt-4o", "gpt-4o-mini", "gpt-3-5-turbo", "Meta-Llama-8B", "gemini-1-5-flash", "gemini-pro"}
var availablePersonas = []string{"pragmatic_analyst", "visionary_idealist", "data_driven_skeptic", "empathetic_humanist"}

// Check if a model is available
func IsModelAvailable(modelName string) bool {
    for _, model := range availableModels {
        if model == modelName {
            return true
        }
    }
    return false
}

// Check if a persona is available
func IsPersonaAvailable(personaName string) bool {
    for _, persona := range availablePersonas {
        if persona == personaName {
            return true
        }
    }
    return false
}

func getModelArguments(history []DebateExchange, model string) []string {
    var arguments []string
    for _, exchange := range history {
        if exchange.Speaker == model {
            arguments = append(arguments, exchange.Message)
        }
    }
    return arguments
}


func initiateDebateRound(ps *PersonaStore, model, personaID, debateTopic, argument1, argument2, promptType string) (string, error) {

    persona, ok := ps.personas[personaID]
    if !ok {
        return "", fmt.Errorf("persona '%s' not found in PersonaStore", personaID)
    }

    var role, prompt string
    var err error

    
        role, err = GetPrompt(persona, "role", debateTopic, "", "")
        if err != nil {
            return "", err
        }
    

    if promptType != "role" && promptType != "" {
        prompt, err = GetPrompt(persona, promptType, debateTopic, argument1, argument2)
        if err != nil {
            return "", err
        }
    }

    return InvokeModel(model, role, prompt)
}





// Retry helper function with exponential backoff
func retryWithBackoff(attempts int, delay time.Duration, fn func() (string, error)) (string, error) {
	var result string
	var err error
	for i := 0; i < attempts; i++ {
		result, err = fn()
		if err == nil {
			return result, nil
		}
		// Wait before retrying
		time.Sleep(delay)
		delay *= 2 // Exponential backoff
	}
	return "", fmt.Errorf("max retries reached: %w", err)
}

func runDebate(debateID, model1, model2, persona1, persona2, debateTopic string) (string, error) {
	// Check if models and personas are available
	if !IsModelAvailable(model1) || !IsModelAvailable(model2) {
		return "", fmt.Errorf("one or more models not available")
	}
	if !IsPersonaAvailable(persona1) || !IsPersonaAvailable(persona2) {
		return "", fmt.Errorf("one or more personas not available")
	}
	if debateTopic == "" {
		return "", fmt.Errorf("debate topic cannot be empty")
	}

	ps := NewPersonaStore()
	var debateHistory []DebateExchange
	var debateArgument DebateExchange

	// Helper to invoke API calls with retry
	invokeWithRetry := func(model, persona, debateTopic, prevArg1, prevArg2, promptType string) (string, error) {
		return retryWithBackoff(5, 500*time.Millisecond, func() (string, error) {
			return initiateDebateRound(ps, model, persona, debateTopic, prevArg1, prevArg2, promptType)
		})
	}

	// Initial opening statements
	model1Response, err := invokeWithRetry(model1, persona1, debateTopic, "", "", "")
	if err != nil {
		return "", err
	}
	debateArgument = DebateExchange{
		PartitionKey: debateID,
		Timestamp:    time.Now(),
		Speaker:      persona1,
		Message:      model1Response,
	}
	debateHistory = append(debateHistory, debateArgument)
	StoreExchange(debateID, debateArgument)
	time.Sleep(500 * time.Millisecond) // Wait before next API call

	model2Response, err := invokeWithRetry(model2, persona2, debateTopic, "", model1Response, "rebuttal")
	if err != nil {
		return "", err
	}
	debateArgument = DebateExchange{
		PartitionKey: debateID,
		Timestamp:    time.Now(),
		Speaker:      persona2,
		Message:      model2Response,
	}
	debateHistory = append(debateHistory, debateArgument)
	StoreExchange(debateID, debateArgument)
	time.Sleep(500 * time.Millisecond) 

	// Debate loop for back-and-forth exchanges
	for i := 3; i <= 8; i++ {
		var currentModel, currentPersona, promptType string

		if i%2 != 0 { 
			currentModel = model1
			currentPersona = persona1
			promptType = "rebuttal"
			if i == 7 {
				promptType = "escalation"
			}
		} else { 
			currentModel = model2
			currentPersona = persona2
			promptType = "rebuttal"
			if i == 6 || i == 8 {
				promptType = "escalation"
			}
		}

		response, err := invokeWithRetry(currentModel, currentPersona, debateTopic, model1Response, model2Response, promptType)
		if err != nil {
			return "", err
		}
		debateArgument = DebateExchange{
			PartitionKey: debateID,
			Timestamp:    time.Now(),
			Speaker:      currentPersona,
			Message:      response,
		}
		debateHistory = append(debateHistory, debateArgument)
		StoreExchange(debateID, debateArgument)
		time.Sleep(500 * time.Millisecond) 

		if i%2 != 0 {
			model1Response = response
		} else {
			model2Response = response
		}
	}

	// Closing statements
	model1Response, err = invokeWithRetry(model1, persona1, debateTopic, strings.Join(getModelArguments(debateHistory, model1), "\n"), "", "closing")
	if err != nil {
		return "", err
	}
	debateArgument = DebateExchange{
		PartitionKey: debateID,
		Timestamp:    time.Now(),
		Speaker:      persona1,
		Message:      model1Response,
	}
	debateHistory = append(debateHistory, debateArgument)
	StoreExchange(debateID, debateArgument)
	time.Sleep(500 * time.Millisecond)

	model2Response, err = invokeWithRetry(model2, persona2, debateTopic, strings.Join(getModelArguments(debateHistory, model2), "\n"), "", "closing")
	if err != nil {
		return "", err
	}
	debateArgument = DebateExchange{
		PartitionKey: debateID,
		Timestamp:    time.Now(),
		Speaker:      persona2,
		Message:      model2Response,
	}
	debateHistory = append(debateHistory, debateArgument)
	StoreExchange(debateID, debateArgument)
	time.Sleep(500 * time.Millisecond)

	// Judge's decision
	conversationSummary := ""
	for i, exchange := range debateHistory {
		conversationSummary += fmt.Sprintf("Turn %d - %s: %s\n", i+1, exchange.Speaker, exchange.Message)
	}

	judgeResponse, err := invokeWithRetry("gemini-pro", "debate_judge", debateTopic, conversationSummary, "", "judge")
	if err != nil {
		return "", err
	}
	debateArgument = DebateExchange{
		PartitionKey: debateID,
		Timestamp:    time.Now(),
		Speaker:      "Judge",
		Message:      judgeResponse,
	}
	StoreExchange(debateID, debateArgument)

	return judgeResponse, nil
}


func StartDebate(hi string) string{
    // Test the RunDebate function
    debateID := "test_debate_006"
    model1 := "gemini-1-5-flash"
    model2 := "gemini-1-5-flash"
    persona1 := "pragmatic_analyst"
    persona2 := "visionary_idealist"
    debateTopic := "The impact of AI on society"

    judgeResponse, err := runDebate(debateID, model1, model2, persona1, persona2, debateTopic)
    if err != nil {
        return fmt.Sprintf("Error running debate: %v\n", err)
    }

    return fmt.Sprintf("Judge's decision: %s\n", judgeResponse)
}
