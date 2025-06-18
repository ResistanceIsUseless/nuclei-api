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

// ScanProgress represents the progress information for a scan
type ScanProgress struct {
	Percentage     float64   `json:"percentage"`
	TemplatesRun   int       `json:"templates_run"`
	TotalTemplates int       `json:"total_templates"`
	VulnsFound     int       `json:"vulns_found"`
	CurrentTarget  string    `json:"current_target,omitempty"`
	LastUpdate     time.Time `json:"last_update"`
}

// ScanRequest represents a nuclei scan request
type ScanRequest struct {
	ID        uuid.UUID    `json:"id"`
	Target    string       `json:"target"`
	Status    ScanStatus   `json:"status"`
	Progress  ScanProgress `json:"progress"`
	CreatedAt time.Time    `json:"created_at"`
	StartedAt *time.Time   `json:"started_at,omitempty"`
	EndedAt   *time.Time   `json:"ended_at,omitempty"`
	Error     string       `json:"error,omitempty"`
	Results   string       `json:"results,omitempty"`
}

// NewScanRequest creates a new scan request with a generated UUID
func NewScanRequest(target string) *ScanRequest {
	return &ScanRequest{
		ID:     uuid.New(),
		Target: target,
		Status: StatusQueued,
		Progress: ScanProgress{
			Percentage:     0,
			TemplatesRun:   0,
			TotalTemplates: 0,
			VulnsFound:     0,
			LastUpdate:     time.Now(),
		},
		CreatedAt: time.Now(),
	}
}

// UpdateProgress updates the progress information for a scan
func (s *ScanRequest) UpdateProgress(percentage float64, templatesRun, totalTemplates, vulnsFound int, currentTarget string) {
	s.Progress = ScanProgress{
		Percentage:     percentage,
		TemplatesRun:   templatesRun,
		TotalTemplates: totalTemplates,
		VulnsFound:     vulnsFound,
		CurrentTarget:  currentTarget,
		LastUpdate:     time.Now(),
	}
}
