package services

import (
	"context"
	"log"
	"sync/atomic"
	"time"

	nuclei "github.com/projectdiscovery/nuclei/v3/lib"
	"github.com/projectdiscovery/nuclei/v3/pkg/output"
	"github.com/resistanceisuseless/nuclei-api/internal/models"
	"github.com/resistanceisuseless/nuclei-api/internal/queue"
)

// ScanService handles scan operations
type ScanService struct {
	queue *queue.Queue
}

// NewScanService creates a new scan service
func NewScanService(q *queue.Queue) *ScanService {
	return &ScanService{
		queue: q,
	}
}

// StartScan starts a new scan
func (s *ScanService) StartScan(ctx context.Context, request *models.ScanRequest) {
	request.Status = models.StatusRunning
	now := time.Now()
	request.StartedAt = &now
	s.queue.Update(request)

	err := s.executeScan(ctx, request)
	if err != nil {
		log.Printf("Error executing scan %s: %v", request.ID, err)
		request.Status = models.StatusFailed
		request.Error = err.Error()
	} else {
		request.Status = models.StatusCompleted
	}
	now = time.Now()
	request.EndedAt = &now
	s.queue.Update(request)
}

func (s *ScanService) executeScan(ctx context.Context, request *models.ScanRequest) error {
	// Create engine with any needed options (add more as needed)
	engine, err := nuclei.NewNucleiEngine()
	if err != nil {
		return err
	}
	defer engine.Close()

	// Load targets
	engine.LoadTargets([]string{request.Target}, false)

	var resultCount int32
	var vulnsFound int32
	var maxTemplatesSeen int32

	callback := func(result *output.ResultEvent) {
		if result == nil {
			return
		}
		atomic.AddInt32(&resultCount, 1)
		tags := result.Info.Tags.ToSlice()
		for _, tag := range tags {
			if tag == "vuln" {
				atomic.AddInt32(&vulnsFound, 1)
				break
			}
		}

		// Update max templates seen
		currentCount := atomic.LoadInt32(&resultCount)
		if currentCount > atomic.LoadInt32(&maxTemplatesSeen) {
			atomic.StoreInt32(&maxTemplatesSeen, currentCount)
		}

		// Calculate progress percentage based on templates run vs max seen
		templatesRun := int(currentCount)
		maxTemplates := int(atomic.LoadInt32(&maxTemplatesSeen))
		percentage := 0.0
		if maxTemplates > 0 {
			percentage = float64(templatesRun) / float64(maxTemplates) * 100
		}

		request.Progress.TemplatesRun = templatesRun
		request.Progress.TotalTemplates = maxTemplates
		request.Progress.VulnsFound = int(vulnsFound)
		request.Progress.Percentage = percentage
		request.Progress.LastUpdate = time.Now()
		s.queue.Update(request)
	}

	err = engine.ExecuteCallbackWithCtx(ctx, callback)

	// Final update
	request.Progress.TemplatesRun = int(resultCount)
	request.Progress.TotalTemplates = int(maxTemplatesSeen)
	request.Progress.VulnsFound = int(vulnsFound)
	request.Progress.Percentage = 100.0 // Set to 100% when complete
	request.Progress.LastUpdate = time.Now()
	s.queue.Update(request)
	return err
}
