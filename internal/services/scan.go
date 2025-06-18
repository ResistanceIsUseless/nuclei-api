package services

import (
	"context"
	"os/exec"
	"time"

	"github.com/resistanceisuseless/nuclei-api/internal/models"
	"github.com/resistanceisuseless/nuclei-api/internal/queue"
)

// ScanService handles the execution of nuclei scans
type ScanService struct {
	queue *queue.Queue
}

// NewScanService creates a new scan service
func NewScanService(queue *queue.Queue) *ScanService {
	return &ScanService{
		queue: queue,
	}
}

// StartScan starts a new scan in the background
func (s *ScanService) StartScan(ctx context.Context, request *models.ScanRequest) error {
	// Update request status to running
	now := time.Now()
	request.Status = models.StatusRunning
	request.StartedAt = &now
	s.queue.Update(request)

	// Start scan in background
	go s.executeScan(ctx, request)
	return nil
}

// executeScan runs the nuclei scan and updates the request status
func (s *ScanService) executeScan(ctx context.Context, request *models.ScanRequest) {
	defer func() {
		now := time.Now()
		request.EndedAt = &now
		s.queue.Update(request)
	}()

	// Execute nuclei command
	cmd := exec.CommandContext(ctx, "nuclei", "-u", request.Target)
	output, err := cmd.CombinedOutput()

	if err != nil {
		request.Status = models.StatusFailed
		request.Error = err.Error()
		return
	}

	request.Status = models.StatusCompleted
	request.Results = string(output)
}
