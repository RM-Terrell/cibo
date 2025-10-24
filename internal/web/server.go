package web

import (
	"cibo/internal/statistics/io"

	"context"
	"encoding/json"
	"fmt"
	"io/fs"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func PrepareListener() (net.Listener, string, error) {
	listener, err := net.Listen("tcp", "localhost:0")
	if err != nil {
		return nil, "", fmt.Errorf("failed to find a free port: %w", err)
	}
	port := listener.Addr().(*net.TCPAddr).Port
	url := fmt.Sprintf("http://localhost:%d", port)
	return listener, url, nil
}

// Method to handle shared server setup code
func newServerHandler(filePath string) http.Handler {
	staticFS, err := fs.Sub(EmbeddedFiles, "static")
	if err != nil {
		// Using panic here because if embedded files are broken, the app can't run.
		log.Panicf("Failed to create sub-filesystem: %v", err)
	}

	parquetClient := io.NewParquetClient()

	mux := http.NewServeMux()
	mux.Handle("/", http.FileServer(http.FS(staticFS)))
	mux.HandleFunc("/api/data", func(w http.ResponseWriter, r *http.Request) {
		records, err := parquetClient.ReadCombinedPriceDataFromParquet(filePath)
		if err != nil {
			log.Printf("API ERROR: %v", err)
			http.Error(w, "Could not read data file", http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(records)
	})

	return mux
}

// StartNonBlocking starts the web server in a goroutine for use with the TUI.
func StartNonBlocking(listener net.Listener, filePath string) {
	handler := newServerHandler(filePath)
	server := &http.Server{Handler: handler}

	go func() {
		if err := server.Serve(listener); err != nil && err != http.ErrServerClosed {
			// Since this is a background task, we log the error instead of fatally exiting.
			log.Printf("ERROR: Non-blocking web server failed: %v", err)
		}
	}()
}

// StartServer starts the web server and blocks, handling graceful shutdown.
func StartServer(listener net.Listener, filePath string) {
	handler := newServerHandler(filePath)
	server := &http.Server{Handler: handler}

	go func() {
		if err := server.Serve(listener); err != nil && err != http.ErrServerClosed {
			log.Fatalf("HTTP server failed: %v", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Shutdown signal received, shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		log.Fatalf("Server forced to shutdown: %v", err)
	}
}
