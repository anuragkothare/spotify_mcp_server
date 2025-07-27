#!/bin/bash

set -e

echo "Setting up Spotify MCP Server..."

# Check if Go is installed
if ! command -v go &> /dev/null; then
    echo "Go is not installed. Please install Go 1.21 or later."
    exit 1
fi

# Check Go version
GO_VERSION=$(go version | awk '{print $3}' | sed 's/go//')
REQUIRED_VERSION="1.21"

if [ "$(printf '%s\n' "$REQUIRED_VERSION" "$GO_VERSION" | sort -V | head -n1)" != "$REQUIRED_VERSION" ]; then
    echo "Go version $REQUIRED_VERSION or later is required. Current version: $GO_VERSION"
    exit 1
fi

# Initialize Go module if go.mod doesn't exist
if [ ! -f go.mod ]; then
    echo "Initializing Go module..."
    go mod init github.com/your-org/spotify-mcp-server
fi

# Create .env file if it doesn't exist
if [ ! -f .env ]; then
    echo "Creating .env file from template..."
    cp .env.example .env
    echo "Please edit .env file with your Spotify API credentials"
fi

# Download dependencies
echo "Downloading Go dependencies..."
go mod tidy
go mod download

# Verify dependencies
echo "Verifying dependencies..."
go mod verify

# Build the application
echo "Building the application..."
make build

echo "Setup complete!"
echo ""
echo "Next steps:"
echo "1. Edit the .env file with your Spotify API credentials"
echo "2. Run 'make run' to start the server"
echo "3. Or run 'make docker-run' to start with Docker"
```

## Setup Instructions

1. **Get Spotify API Credentials**:
   - Go to [Spotify Developer Dashboard](https://developer.spotify.com/dashboard)
   - Create a new app
   - Copy Client ID and Client Secret

2. **Setup Environment**:
   ```bash
   git clone <repository>
   cd spotify-mcp-server
   
   # If you don't have go.mod yet:
   make init
   
   # Install dependencies
   make deps
   
   chmod +x scripts/setup.sh
   ./scripts/setup.sh
   ```

3. **Configure Environment Variables**:
   ```bash
   cp .env.example .env
   # Edit .env with your Spotify credentials
   ```

4. **Run the Server**:
   ```bash
   # Local development
   make run
   
   # Or with Docker
   make docker-run
   ```

## Troubleshooting

### Dependency Issues

If you get import errors, try these steps:

```bash
# Clean and reinitialize
go clean -modcache
make init
make deps

# Or manually:
go mod init github.com/your-org/spotify-mcp-server
go mod tidy
go get github.com/sirupsen/logrus@v1.9.3
go get github.com/zmb3/spotify/v2@v2.4.1
go get github.com/spf13/viper@v1.18.2
go get golang.org/x/oauth2@v0.15.0
```

### Common Issues

1. **"no required module provides package"** - Run `make deps` or `go mod tidy`
2. **Type conversion errors** - Make sure you're using the latest spotify/v2 library
3. **Import cycle** - Check that internal packages don't import each other circularly
4. **"undefined method handleSearchTracks"** - All handler methods are now in `server.go`, no separate `tools.go` file needed

## API Endpoints

- `POST /mcp` - MCP protocol endpoint
- `GET /health` - Health check endpoint

## Available MCP Tools

1. **search_tracks** - Search for tracks on Spotify
2. **search_artists** - Search for artists on Spotify  
3. **get_track** - Get detailed track information

## Example MCP Requests

```json
{
  "jsonrpc": "2.0",
  "id": 1,
  "method": "tools/call",
  "params": {
    "name": "search_tracks",
    "arguments": {
      "query": "bohemian rhapsody",
      "limit": 5
    }
  }
}