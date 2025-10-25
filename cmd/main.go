package main

import (
	"cibo/internal/pipelines"
	"cibo/internal/statistics/api"
	"cibo/internal/statistics/config"
	"cibo/internal/statistics/io"
	"cibo/internal/tui"
	"cibo/internal/web"
	"flag"
	"fmt"
	"log"
	"os"

	tea "github.com/charmbracelet/bubbletea"
)

const (
	alphaVantageURL     = "https://www.alphavantage.co"
	mockAlphaVantageURL = "http://localhost:8080"
)

func main() {
	webModeFilePath := flag.String("webMode", "", "Path to a Parquet file to display in standalone web mode.")
	useMockAPI := flag.Bool("mock-api", false, "Use the mock API server instead of Alpha Vantage")
	flag.Parse()

	if *webModeFilePath != "" {
		fmt.Printf("Starting server in standalone mode for file: %s\n", *webModeFilePath)
		listener, url, err := web.PrepareListener()
		if err != nil {
			log.Fatalf("Failed to prepare web listener: %v", err)
		}

		fmt.Printf("Web server starting. Open this URL in your browser: %s\n", url)
		fmt.Println("Press Ctrl+C to shut down the server.")

		web.StartServer(listener, *webModeFilePath)
		fmt.Println("Server shutting down.")
		return
	}

	var baseURL string
	if *useMockAPI {
		baseURL = mockAlphaVantageURL
		log.Println("Using mock API server.")
	} else {
		baseURL = alphaVantageURL
		log.Println("Using live Alpha Vantage API.")
	}

	configPath := os.Getenv("API_KEYS_CONFIG_PATH")
	if configPath == "" {
		log.Fatal("Error: API_KEYS_CONFIG_PATH environment variable not set.")
	}

	cfg, err := config.LoadConfig(configPath)
	if err != nil {
		log.Fatalf("Error loading configuration: %v", err)
	}

	apiClient := api.NewClient(cfg.AlphaVantageAPIKey, baseURL)
	parquetWriter := io.NewParquetClient()

	pipelines := pipelines.NewPipelines(apiClient, parquetWriter)

	p := tea.NewProgram(tui.NewModel(pipelines))

	if _, err := p.Run(); err != nil {
		log.Fatalf("There's been an error: %v", err)
	}
}
