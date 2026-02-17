package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/chili-network/scoville/internal/config"
	"github.com/chili-network/scoville/internal/engine"
	"github.com/chili-network/scoville/internal/ui"
)

func main() {
	cfg := config.Load()

	if !cfg.PaperMode {
		fmt.Println("====================================================================")
		fmt.Println("  WARNING: LIVE TRADING MODE IS ENABLED")
		fmt.Println("====================================================================")
		fmt.Println("This will execute REAL trades on the blockchain with REAL funds.")
		fmt.Println("Ensure your configuration in the .env file is correct before proceeding.")
		fmt.Println("")
		fmt.Print("Type 'YES' to confirm and start live trading: ")

		reader := bufio.NewReader(os.Stdin)
		input, _ := reader.ReadString('\n')
		input = strings.TrimSpace(input)

		if input != "YES" {
			fmt.Println("Confirmation not received. Exiting.")
			os.Exit(0)
		}
		fmt.Println("Confirmation received. Starting Scoville...")
	}

	// Initialize Bot
	bot := engine.New(cfg)

	// Initialize TUI
	p := tea.NewProgram(ui.InitialModel(bot), tea.WithAltScreen())

	if _, err := p.Run(); err != nil {
		fmt.Printf("Alas, there's been an error: %v", err)
		os.Exit(1)
	}
}
