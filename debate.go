package main

import "fmt"

var availableModels = []string{
    "gpt-4o", "gpt-4o-mini", "gpt-3-5-turbo", "Meta-Llama-8B",
}

var availablePersonas = []string{
	"pragmatic_analyst", "visionary_idealist", "data_driven_skeptic", "empathetic_humanist",
}

func IsModelAvailable(modelName string) bool {
    for _, model := range availableModels {
        if model == modelName {
            return true
        }
    }
    return false
}

func IsPersonaAvailable(personaName string) bool {
	for _, persona := range availablePersonas {
		if persona == personaName {
			return true
		}
	}
	return false
}

func RunDebate(model1, model2, persona1, persona2, debateTopic string) (string, error) {
   
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
	responsePrompt := "Stick to the debate style you were assigned in the instructions, Respond using that style and address the opponent's argument."
    prompt := ""

     prompt, err := ps.GetPrompt(persona1, "opening", debateTopic)
	 if err != nil {
        return "", err
    }
	// Opening statement
	var result1 string
	var result2 string
	result1, err = InvokeModel(model1, prompt, "Begin your opening argument")
	if err != nil {
        return "", err
    }
    for i := 2; i <= 8; i++ {
        if i%2 == 0 {
            if i == 6 || i == 8 {
				// Model2 escalates
				prompt, err = ps.GetPrompt(persona2, "escalation", result1)
			} else {
				// Model2 responds
				prompt, err = ps.GetPrompt(persona2, "rebuttal", result1)
			}
					
			if err != nil {
				return "", err
			}
			result2, err = InvokeModel(model2, prompt, responsePrompt)
        } else {
          // Model1 escalates
			if i == 7 {
				prompt, err = ps.GetPrompt(persona1, "escalation", result2)
			} else {
				// Model1 responds
			prompt, err = ps.GetPrompt(persona1, "rebuttal", result2)
			}
			if err != nil {
				return "", err
			}
			result1, err = InvokeModel(model1, prompt, responsePrompt)
        }
    }

    // Closing statements to the judge
    fmt.Println("Model1 Closing:", GetPrompt("Model1", "closing", debateTopic))
    fmt.Println("Model2 Closing:", GetPrompt("Model2", "closing", debateTopic))

    // Judge final decision
    winner := JudgeDecision(model1, model2, debateTopic)
    return winner, nil
}
