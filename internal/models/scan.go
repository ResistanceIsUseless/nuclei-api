package models

import (
	"time"

	"github.com/google/uuid"
)

// ScanStatus represents the current status of a scan
type ScanStatus string

const (
	StatusQueued    ScanStatus = "queued"
	StatusRunning   ScanStatus = "running"
	StatusCompleted ScanStatus = "completed"
	StatusFailed    ScanStatus = "failed"
)

// ScanRequest represents a nuclei scan request
type ScanRequest struct {
	ID        uuid.UUID  `json:"id"`
	Target    string     `json:"target"`
	Status    ScanStatus `json:"status"`
	CreatedAt time.Time  `json:"created_at"`
	StartedAt *time.Time `json:"started_at,omitempty"`
	EndedAt   *time.Time `json:"ended_at,omitempty"`
	Error     string     `json:"error,omitempty"`
	Results   string     `json:"results,omitempty"`
}

// NewScanRequest creates a new scan request with a generated UUID
func NewScanRequest(target string) *ScanRequest {
	return &ScanRequest{
		ID:        uuid.New(),
		Target:    target,
		Status:    StatusQueued,
		CreatedAt: time.Now(),
	}
}
