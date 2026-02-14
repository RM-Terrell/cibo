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
	"strings"
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

type spaHandler struct {
	staticFS  fs.FS
	indexPath string
}

func (h spaHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	path := r.URL.Path
	if path == "/" {
		path = h.indexPath
	}
	path = strings.TrimPrefix(path, "/")

	_, err := fs.Stat(h.staticFS, path)
	if os.IsNotExist(err) {
		http.ServeFileFS(w, r, h.staticFS, h.indexPath)
		return
	} else if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	http.ServeFileFS(w, r, h.staticFS, path)
}

// Method to handle shared server setup code
func newServerHandler(filePath string) http.Handler {
	staticFS, err := fs.Sub(EmbeddedFiles, "dist")
	if err != nil {
		log.Panicf("Failed to create sub-filesystem: %v", err)
	}

	parquetClient := io.NewParquetClient()

	mux := http.NewServeMux()
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

	mux.Handle("/", spaHandler{staticFS: staticFS, indexPath: "index.html"})

	return mux
}

// StartNonBlocking and StartServer remain the same...
func StartNonBlocking(listener net.Listener, filePath string) {
	handler := newServerHandler(filePath)
	server := &http.Server{Handler: handler}

	go func() {
		if err := server.Serve(listener); err != nil && err != http.ErrServerClosed {
			log.Printf("ERROR: Non-blocking web server failed: %v", err)
		}
	}()
}

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
