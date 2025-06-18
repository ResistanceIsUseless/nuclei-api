# Nuclei API Enhancement Progress

This document tracks the progress of planned enhancements to the Nuclei API.

## Features

### Core Features
- [x] 🆔 Request IDs: Every scan gets a unique identifier (UUID) for tracking
  - Implemented UUID generation for each scan request
  - Added status tracking (queued, running, completed, failed)
  - Added timestamps for request lifecycle

- [x] ⚡ Async Processing: Queue scans for background execution
  - Implemented queue system with 10 concurrent scan limit
  - Added background processing of scans
  - Added queue management endpoints

### In Progress
- [ ] 📊 Progress Tracking: Real-time scan progress and statistics
  - [ ] Add progress percentage tracking
  - [ ] Add scan statistics (templates run, vulnerabilities found)
  - [ ] Add real-time progress updates via WebSocket

### Planned Features
- [ ] 🔄 Optional TUI to see the status of the server and progress of each scan
  - [ ] Design TUI layout
  - [ ] Implement server status dashboard
  - [ ] Add real-time scan progress visualization
  - [ ] Add queue management interface

- [ ] 🔄 Queue Management: View and manage queued scans
  - [ ] Add queue status endpoint
  - [ ] Add ability to cancel queued scans
  - [ ] Add queue statistics
  - [ ] Add queue management UI

- [ ] 🪝 Webhooks: Automatic result delivery via callbacks
  - [ ] Add webhook configuration
  - [ ] Implement webhook delivery system
  - [ ] Add webhook retry mechanism
  - [ ] Add webhook security (signatures)

- [ ] 📈 Priority Queuing: High-priority scans jump the queue
  - [ ] Add priority levels to scan requests
  - [ ] Implement priority queue system
  - [ ] Add priority management endpoints
  - [ ] Add priority visualization in UI

- [ ] 🔄 Uploading templates through the API
  - [ ] Add template upload endpoint
  - [ ] Add template validation
  - [ ] Add template management system
  - [ ] Add template versioning

## Implementation Notes

### Completed Features
1. Request IDs and Async Processing
   - Branch: `feature/request-id-and-async`
   - Status: ✅ Completed
   - Key Components:
     - UUID generation and tracking
     - Queue system with concurrent scan limit
     - Background processing
     - Status tracking and timestamps

### Next Steps
1. Progress Tracking
   - Create new branch: `feature/progress-tracking`
   - Implement real-time progress updates
   - Add statistics collection
   - Add WebSocket support for live updates

## Technical Debt
- Add proper error handling and logging
- Add comprehensive test coverage
- Add API documentation
- Add rate limiting
- Add authentication and authorization 