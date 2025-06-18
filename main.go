package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/resistanceisuseless/nuclei-api/internal/models"
	"github.com/resistanceisuseless/nuclei-api/internal/queue"
	"github.com/resistanceisuseless/nuclei-api/internal/services"
)

func main() {
	// Initialize queue with max 10 concurrent scans
	scanQueue := queue.NewQueue(10)
	scanService := services.NewScanService(scanQueue)

	r := gin.Default()

	// Start a new scan
	r.POST("/scan", func(c *gin.Context) {
		var request struct {
			Target string `json:"target" binding:"required"`
		}

		if err := c.ShouldBindJSON(&request); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		scanRequest := models.NewScanRequest(request.Target)
		if !scanQueue.Add(scanRequest) {
			c.JSON(http.StatusServiceUnavailable, gin.H{"error": "scan queue is full"})
			return
		}

		if err := scanService.StartScan(c.Request.Context(), scanRequest); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusAccepted, scanRequest)
	})

	// Get scan status
	r.GET("/scan/:id", func(c *gin.Context) {
		id := c.Param("id")
		request, exists := scanQueue.Get(id)
		if !exists {
			c.JSON(http.StatusNotFound, gin.H{"error": "scan not found"})
			return
		}

		c.JSON(http.StatusOK, request)
	})

	// List all scans
	r.GET("/scans", func(c *gin.Context) {
		requests := scanQueue.List()
		c.JSON(http.StatusOK, requests)
	})

	r.Run(":8080")
}
