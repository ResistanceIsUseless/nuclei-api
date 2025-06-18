package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/resistanceisuseless/nuclei-api/internal/models"
	"github.com/resistanceisuseless/nuclei-api/internal/version"
)

func printStatus() error {
	// Fetch scans from API
	resp, err := http.Get("http://localhost:8080/scans")
	if err != nil {
		return fmt.Errorf("failed to connect to API: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("API returned status: %d", resp.StatusCode)
	}

	var scans []models.ScanRequest
	if err := json.NewDecoder(resp.Body).Decode(&scans); err != nil {
		return fmt.Errorf("failed to decode response: %v", err)
	}

	queued := 0
	running := 0
	completed := 0
	failed := 0

	for _, scan := range scans {
		switch scan.Status {
		case models.StatusQueued:
			queued++
		case models.StatusRunning:
			running++
		case models.StatusCompleted:
			completed++
		case models.StatusFailed:
			failed++
		}
	}

	fmt.Printf("\033[2J\033[H") // Clear screen and move cursor to top
	fmt.Println(version.String())
	fmt.Println("\nConnected to API at http://localhost:8080")
	fmt.Println("\nAPI Endpoints:")
	fmt.Println("  POST /scan     - Start a new scan")
	fmt.Println("  GET  /scan/:id - Get scan status")
	fmt.Println("  GET  /scans    - List all scans")
	fmt.Println("  GET  /version  - Get version info")
	fmt.Println("\nScan Status:")
	fmt.Printf("  Queued:    %d\n", queued)
	fmt.Printf("  Running:   %d\n", running)
	fmt.Printf("  Completed: %d\n", completed)
	fmt.Printf("  Failed:    %d\n", failed)
	fmt.Printf("  Total:     %d\n", len(scans))
	fmt.Println("\nPress Ctrl+C to stop")

	return nil
}

func main() {
	fmt.Println("Starting Nuclei API TUI Dashboard...")
	fmt.Println("Make sure the API server is running on http://localhost:8080")

	for {
		if err := printStatus(); err != nil {
			log.Printf("Error: %v", err)
		}
		time.Sleep(2 * time.Second)
	}
}
