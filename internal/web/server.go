package web

import (
	"cibo/internal/types"
	"encoding/json"
	"fmt"
	"io/fs"
	"log"
	"net"
	"net/http"

	"github.com/xitongsys/parquet-go-source/local"
	"github.com/xitongsys/parquet-go/reader"
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

func Start(listener net.Listener, filePath string) {
	staticFS, err := fs.Sub(EmbeddedFiles, "static")
	if err != nil {
		log.Fatalf("Failed to create sub-filesystem: %v", err)
	}

	mux := http.NewServeMux()
	mux.Handle("/", http.FileServer(http.FS(staticFS)))
	mux.HandleFunc("/api/data", func(w http.ResponseWriter, r *http.Request) {
		fr, err := local.NewLocalFileReader(filePath)
		if err != nil {
			http.Error(w, "Could not read data file", http.StatusInternalServerError)
			return
		}
		defer fr.Close()

		pr, err := reader.NewParquetReader(fr, new(types.CombinedPriceRecord), 4)
		if err != nil {
			http.Error(w, "Could not read data file", http.StatusInternalServerError)
			return
		}
		defer pr.ReadStop()

		numRecords := int(pr.GetNumRows())
		records := make([]types.CombinedPriceRecord, numRecords)
		if err := pr.Read(&records); err != nil {
			http.Error(w, "Could not read data file", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(records)
	})

	if err := http.Serve(listener, mux); err != nil {
		log.Fatalf("HTTP server failed: %v", err)
	}
}
