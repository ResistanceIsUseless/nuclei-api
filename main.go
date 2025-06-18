package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/resistanceisuseless/nuclei-api/internal/models"
	"github.com/resistanceisuseless/nuclei-api/internal/queue"
	"github.com/resistanceisuseless/nuclei-api/internal/services"
	"github.com/resistanceisuseless/nuclei-api/internal/version"
)

func printStatus(q *queue.Queue) {
	scans := q.List()
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
	fmt.Println("\nServer is running on http://localhost:8080")
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
}

func main() {
	// Create a new queue
	q := queue.NewQueue()

	// Create a new scan service
	scanService := services.NewScanService(q)

	// Create a new Gin router
	r := gin.Default()

	// Add version information to the status UI
	r.GET("/version", func(c *gin.Context) {
		c.JSON(http.StatusOK, version.Get())
	})

	// Start a new scan
	r.POST("/scan", func(c *gin.Context) {
		var request struct {
			Target string `json:"target" binding:"required"`
		}
		if err := c.ShouldBindJSON(&request); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		// Create a new scan request
		scan := models.NewScanRequest(request.Target)

		// Add the scan to the queue
		if err := q.Add(scan); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		// Start the scan in a goroutine
		go scanService.StartScan(context.Background(), scan)

		c.JSON(http.StatusAccepted, gin.H{
			"id":     scan.ID,
			"status": scan.Status,
		})
	})

	// Get scan status
	r.GET("/scan/:id", func(c *gin.Context) {
		id := c.Param("id")
		scan, err := q.Get(id)
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "Scan not found"})
			return
		}

		c.JSON(http.StatusOK, scan)
	})

	// List all scans
	r.GET("/scans", func(c *gin.Context) {
		scans := q.List()
		c.JSON(http.StatusOK, scans)
	})

	// Start the server in a goroutine
	go func() {
		if err := r.Run(":8080"); err != nil {
			log.Fatal(err)
		}
	}()

	// Start status display in a goroutine
	go func() {
		for {
			printStatus(q)
			time.Sleep(2 * time.Second)
		}
	}()

	// Wait for interrupt signal
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	fmt.Println("\nShutting down server...")
}
