package state

import (
	"encoding/json"
	"os"
	"sync"
)

const StateFile = "scoville_progress.json"

type Progress struct {
	CurrentPhase  int     `json:"current_phase"`  // 1=Whale, 2=Retail, 3=LP
	WhaleIndex    int     `json:"whale_index"`    // Which whale are we on?
	RetailIndex   int     `json:"retail_index"`   // Which retail wallet?
	TotalSpent    float64 `json:"total_spent"`
	IsPaused      bool    `json:"is_paused"`
}

var (
	mu           sync.Mutex
	CurrentState Progress
)

// Load reads the JSON file or creates a default state
func Load() Progress {
	mu.Lock()
	defer mu.Unlock()

	file, err := os.ReadFile(StateFile)
	if err != nil {
		// Default State
		return Progress{
			CurrentPhase: 1,
			WhaleIndex:   0,
			RetailIndex:  0,
			TotalSpent:   0.0,
			IsPaused:     false,
		}
	}

	var p Progress
	json.Unmarshal(file, &p)
	CurrentState = p
	return p
}

// Save writes the current state to JSON
func Save(p Progress) error {
	mu.Lock()
	defer mu.Unlock()

	CurrentState = p
	data, _ := json.MarshalIndent(p, "", "  ")
	return os.WriteFile(StateFile, data, 0644)
}
