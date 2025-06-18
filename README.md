# Nuclei API

A Go-based API wrapper around [ProjectDiscovery's Nuclei](https://github.com/projectdiscovery/nuclei) vulnerability scanner.

## Setup

1. Install Go 1.21 or later
2. Clone this repository
3. Install dependencies:
   ```bash
   go mod download
   ```
4. Run the server:
   ```bash
   go run main.go
   ```

## API Endpoints

### Health Check
- `GET /health`
- Returns the API status

### Scan
- `POST /scan`
- Request body:
  ```json
  {
    "target": "example.com"
  }
  ```
- Initiates a Nuclei scan against the specified target

## Development

This is a work in progress. More features and endpoints will be added soon. 