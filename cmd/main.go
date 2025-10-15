package main

import (
	"cibo/internal/pipelines"
	"cibo/internal/statistics/api"
	"cibo/internal/statistics/config"
	"cibo/internal/statistics/io"
	"cibo/internal/tui"
	"log"
	"os"

	tea "github.com/charmbracelet/bubbletea"
)

func main() {
	configPath := os.Getenv("API_KEYS_CONFIG_PATH")
	if configPath == "" {
		log.Fatal("Error: API_KEYS_CONFIG_PATH environment variable not set.")
	}

	cfg, err := config.LoadConfig(configPath)
	if err != nil {
		log.Fatalf("Error loading configuration: %v", err)
	}

	apiClient := api.NewClient(cfg.AlphaVantageAPIKey)
	parquetWriter := io.NewParquetIOAdapter()

	pipelines := pipelines.NewPipelines(apiClient, parquetWriter)

	p := tea.NewProgram(tui.NewModel(pipelines))

	if _, err := p.Run(); err != nil {
		log.Fatalf("Alas, there's been an error: %v", err)
	}
}
