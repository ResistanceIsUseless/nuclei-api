package queue

import (
	"sync"

	"github.com/resistanceisuseless/nuclei-api/internal/models"
)

// Queue manages the scan request queue
type Queue struct {
	mu       sync.RWMutex
	requests map[string]*models.ScanRequest
	maxSize  int
}

// NewQueue creates a new queue with the specified maximum size
func NewQueue(maxSize int) *Queue {
	return &Queue{
		requests: make(map[string]*models.ScanRequest),
		maxSize:  maxSize,
	}
}

// Add adds a scan request to the queue
func (q *Queue) Add(request *models.ScanRequest) bool {
	q.mu.Lock()
	defer q.mu.Unlock()

	if len(q.requests) >= q.maxSize {
		return false
	}

	q.requests[request.ID.String()] = request
	return true
}

// Get retrieves a scan request by ID
func (q *Queue) Get(id string) (*models.ScanRequest, bool) {
	q.mu.RLock()
	defer q.mu.RUnlock()

	request, exists := q.requests[id]
	return request, exists
}

// Update updates a scan request in the queue
func (q *Queue) Update(request *models.ScanRequest) {
	q.mu.Lock()
	defer q.mu.Unlock()

	q.requests[request.ID.String()] = request
}

// Remove removes a scan request from the queue
func (q *Queue) Remove(id string) {
	q.mu.Lock()
	defer q.mu.Unlock()

	delete(q.requests, id)
}

// List returns all scan requests in the queue
func (q *Queue) List() []*models.ScanRequest {
	q.mu.RLock()
	defer q.mu.RUnlock()

	requests := make([]*models.ScanRequest, 0, len(q.requests))
	for _, request := range q.requests {
		requests = append(requests, request)
	}
	return requests
}
