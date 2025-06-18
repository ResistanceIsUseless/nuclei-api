package queue

import (
	"errors"
	"sync"

	"github.com/resistanceisuseless/nuclei-api/internal/models"
)

// Queue manages scan requests
type Queue struct {
	requests map[string]*models.ScanRequest
	mu       sync.RWMutex
}

// NewQueue creates a new queue
func NewQueue() *Queue {
	return &Queue{
		requests: make(map[string]*models.ScanRequest),
	}
}

// Add adds a scan request to the queue
func (q *Queue) Add(request *models.ScanRequest) error {
	q.mu.Lock()
	defer q.mu.Unlock()

	if _, exists := q.requests[request.ID.String()]; exists {
		return errors.New("scan request already exists")
	}

	q.requests[request.ID.String()] = request
	return nil
}

// Get retrieves a scan request by ID
func (q *Queue) Get(id string) (*models.ScanRequest, error) {
	q.mu.RLock()
	defer q.mu.RUnlock()

	request, exists := q.requests[id]
	if !exists {
		return nil, errors.New("scan request not found")
	}

	return request, nil
}

// Update updates a scan request
func (q *Queue) Update(request *models.ScanRequest) error {
	q.mu.Lock()
	defer q.mu.Unlock()

	if _, exists := q.requests[request.ID.String()]; !exists {
		return errors.New("scan request not found")
	}

	q.requests[request.ID.String()] = request
	return nil
}

// Remove removes a scan request from the queue
func (q *Queue) Remove(id string) error {
	q.mu.Lock()
	defer q.mu.Unlock()

	if _, exists := q.requests[id]; !exists {
		return errors.New("scan request not found")
	}

	delete(q.requests, id)
	return nil
}

// List returns all scan requests
func (q *Queue) List() []*models.ScanRequest {
	q.mu.RLock()
	defer q.mu.RUnlock()

	requests := make([]*models.ScanRequest, 0, len(q.requests))
	for _, request := range q.requests {
		requests = append(requests, request)
	}

	return requests
}
