# Nuclei API

A powerful API wrapper for [Nuclei](https://github.com/projectdiscovery/nuclei) that provides asynchronous scanning capabilities, progress tracking, and a real-time TUI dashboard.

## Features

- 🆔 **Request IDs**: Every scan gets a unique identifier for tracking
- ⚡ **Async Processing**: Queue scans for background execution
- 📊 **Progress Tracking**: Real-time scan progress and statistics
- 🖥️ **TUI Dashboard**: Interactive terminal UI for monitoring scans
- 🔄 **Queue Management**: View and manage queued scans

## Prerequisites

- Go 1.21 or later
- Nuclei installed and available in your PATH
- Git

## Installation

### Option 1: Using go install (Recommended)

Install the API server:
```bash
go install github.com/resistanceisuseless/nuclei-api/cmd/nuclei-api@latest
```

Install the TUI dashboard (optional):
```bash
go install github.com/resistanceisuseless/nuclei-api/cmd/nuclei-tui@latest
```

This will install both tools to your `$GOPATH/bin` directory.

### Option 2: Building from source

1. Clone the repository:
```bash
git clone https://github.com/ResistanceIsUseless/nuclei-api.git
cd nuclei-api
```

2. Install dependencies:
```bash
go mod download
```

3. Build the project:
```bash
# Build both components
go build ./cmd/nuclei-api
go build ./cmd/nuclei-tui

# Or build specific component
go build ./cmd/nuclei-api  # API server only
go build ./cmd/nuclei-tui  # TUI dashboard only
```

## Usage

### Starting the API Server

If installed via `go install`:
```bash
nuclei-api
```

If built from source:
```bash
./nuclei-api
```

The server will start on port 8080 by default.

### Using the TUI Dashboard

If installed via `go install`:
```bash
nuclei-tui
```

If built from source:
```bash
./nuclei-tui
```

The TUI dashboard will connect to the API server at http://localhost:8080 and display real-time scan status.

### API Endpoints

#### Start a Scan
```bash
curl -X POST http://localhost:8080/scan \
  -H "Content-Type: application/json" \
  -d '{
    "target": "example.com",
    "templates": ["cves/", "vulnerabilities/"]
  }'
```

#### Get Scan Status
```bash
curl http://localhost:8080/scan/{scan_id}
```

#### List All Scans
```bash
curl http://localhost:8080/scans
```

### TUI Dashboard Features

The TUI dashboard provides real-time monitoring of:
- Server status
- Active scans
- Queue status
- Scan progress
- Statistics

## Configuration

The server can be configured using environment variables:

- `PORT`: Server port (default: 8080)
- `MAX_CONCURRENT_SCANS`: Maximum number of concurrent scans (default: 10)
- `NUCLEI_PATH`: Path to nuclei executable (default: uses PATH)

## Contributing

1. Fork the repository
2. Create your feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add some amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## Acknowledgments

- [ProjectDiscovery](https://github.com/projectdiscovery) for the amazing Nuclei tool
- [Gin](https://github.com/gin-gonic/gin) for the web framework
- [Bubble Tea](https://github.com/charmbracelet/bubbletea) for the TUI framework 