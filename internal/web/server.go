package web

import (
	"cibo/internal/types"
	"encoding/json"
	"fmt"
	"io/fs"
	"log"
	"net"
	"net/http"
	"os/exec"
	"runtime"
)

func Start(chartData []types.CombinedPriceRecord) {
	// Use a listener on port 0 to get a random free port from the OS to try to avoid port conflicts
	listener, err := net.Listen("tcp", "localhost:0")
	if err != nil {
		log.Fatalf("Failed to find a free port: %v", err)
	}
	port := listener.Addr().(*net.TCPAddr).Port
	url := fmt.Sprintf("http://localhost:%d", port)

	// Create a sub-filesystem from our embedded files
	// This makes it so the server's root ("/") corresponds to the "static" directory.
	staticFS, err := fs.Sub(EmbeddedFiles, "static")
	if err != nil {
		log.Fatalf("Failed to create sub-filesystem: %v", err)
	}

	// HANDLER ROOT: Serve static files from the embedded filesystem.
	http.Handle("/", http.FileServer(http.FS(staticFS)))

	// HANDLER DATA: The data API endpoint.
	http.HandleFunc("/api/data", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(chartData)
	})

	log.Printf("Chart visualization server starting at: %s", url)
	log.Println("If your browser does not open automatically, please navigate to the URL above.")

	openBrowser(url)

	if err := http.Serve(listener, nil); err != nil {
		log.Fatalf("HTTP server failed: %v", err)
	}
}

func openBrowser(url string) {
	var cmd string
	var args []string

	switch runtime.GOOS {
	case "windows":
		cmd = "cmd"
		args = []string{"/c", "start", url}
	case "darwin":
		cmd = "open"
		args = []string{url}
	default: // "linux", "freebsd", "openbsd", "netbsd"
		cmd = "xdg-open"
		args = []string{url}
	}

	err := exec.Command(cmd, args...).Start()
	if err != nil {
		log.Printf("Error opening browser: %v", err)
	}
}
