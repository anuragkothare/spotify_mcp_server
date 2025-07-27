<img width="16" height="16" alt="image" src="https://github.com/user-attachments/assets/3b9106a6-6c97-4c37-8d65-9f10acc34912" />#  Spotify MCP Server

A **Model Context Protocol (MCP) server** that provides AI assistants with access to Spotify's music catalog.

[![Go Version](https://img.shields.io/badge/Go-1.23+-00ADD8?style=flat&logo=go)](https://golang.org/)
[![Docker](https://img.shields.io/badge/Docker-Ready-2496ED?style=flat&logo=docker)](https://www.docker.com/)

## 🎯 **What is this?**

This MCP server enables AI assistants to search music on Spotify through three tools:

- **search_tracks**: Find songs by name/artist
- **search_artists**: Find artists with popularity scores
- **get_track**: Get detailed track information

## 📋 **Prerequisites**

- [Spotify Developer Account](https://developer.spotify.com/) (for API credentials)
- [Docker & Docker Compose](https://www.docker.com/) or [Go 1.23+](https://golang.org/)

## 🚀 **Quick Start**

### **1. Get Spotify Credentials**

1. Go to [Spotify Developer Dashboard](https://developer.spotify.com/dashboard)
2. Create a new app
3. Copy your **Client ID** and **Client Secret**

### **2. Setup & Run**

```bash
# Clone and setup
git clone https://github.com/anuragkothare/spotify_mcp_server.git
cd spotify-mcp-server

# Configure environment
cp .env.example .env
# Edit .env with your Spotify credentials:
# SPOTIFY_CLIENT_ID=your_client_id_here
# SPOTIFY_CLIENT_SECRET=your_client_secret_here

# Run with Docker (recommended)
docker-compose up -d

# Or run locally
go mod tidy
go run ./cmd/server
```

### **3. Test**

```bash
# Health check
curl http://localhost:8080/health

# List available tools
curl -X POST http://localhost:8080/mcp \
  -H "Content-Type: application/json" \
  -d '{"jsonrpc":"2.0","id":1,"method":"tools/list"}'

# Search for tracks
curl -X POST http://localhost:8080/mcp \
  -H "Content-Type: application/json" \
  -d '{
    "jsonrpc": "2.0",
    "id": 2,
    "method": "tools/call",
    "params": {
      "name": "search_tracks",
      "arguments": {"query": "bohemian rhapsody", "limit": 5}
    }
  }'
```

## ⚙️ **Configuration**

**Required environment variables:**

```bash
SPOTIFY_CLIENT_ID=your_spotify_client_id
SPOTIFY_CLIENT_SECRET=your_spotify_client_secret
```

**Optional:**

```bash
SPOTIFY_REDIRECT_URI=http://localhost:8080/callback
SERVER_PORT=8080
LOG_LEVEL=info
```

## 🐳 **Docker Commands**

```bash
# Development
docker-compose up -d          # Start services
docker-compose logs -f        # View logs
docker-compose down           # Stop services

# Production
docker-compose -f docker-compose.yml -f docker-compose.prod.yml up -d

# Useful commands
docker-compose ps             # Check status
docker-compose exec spotify-mcp sh  # Access container
```

## 🎵 **Available MCP Tools**

### **search_tracks**

```json
{
  "name": "search_tracks",
  "arguments": {
    "query": "song name or artist",
    "limit": 10
  }
}
```

### **search_artists**

```json
{
  "name": "search_artists",
  "arguments": {
    "query": "artist name",
    "limit": 10
  }
}
```

### **get_track**

```json
{
  "name": "get_track",
  "arguments": {
    "track_id": "spotify_track_id"
  }
}
```

## 🛠️ **Project Structure**

```
spotify-mcp-server/
├── cmd/server/           # Main application
├── internal/
│   ├── config/          # Configuration
│   ├── mcp/            # MCP protocol implementation
│   ├── spotify/        # Spotify API client
│   └── handlers/       # HTTP handlers
├── configs/            # Config files
├── Dockerfile          # Container build
├── docker-compose.yml  # Service orchestration
├── .env.example        # Environment template
└── Makefile           # Build commands
```

## 🚨 **Troubleshooting**

### **"invalid_client" error**

- Check your `.env` file has actual Spotify credentials (not placeholders)
- Verify credentials work: Test at [Spotify Console](https://developer.spotify.com/console/)

### **".env file not found"**

- Ensure `.env` exists in project root
- Run from the directory containing `docker-compose.yml`

### **Docker issues**

```bash
# Clean up and restart
docker-compose down --volumes --remove-orphans
docker-compose up -d --build

# Check logs
docker-compose logs spotify-mcp
```

### **Port conflicts**

```bash
# Check what's using port 8080
lsof -i :8080

# Change port in .env
SERVER_PORT=8081
```

## 🔗 **Integration**

### **Python Example**

```python
import requests

def search_tracks(query, limit=5):
    response = requests.post('http://localhost:8080/mcp', json={
        "jsonrpc": "2.0",
        "id": 1,
        "method": "tools/call",
        "params": {
            "name": "search_tracks",
            "arguments": {"query": query, "limit": limit}
        }
    })
    return response.json()

# Usage
results = search_tracks("hello adele")
print(results)
```

### **Claude Desktop Integration**

Add to your Claude Desktop config:

```json
{
  "mcpServers": {
    "spotify": {
      "command": "docker",
      "args": [
        "run",
        "-p",
        "8080:8080",
        "--env-file",
        ".env",
        "spotify-mcp-server"
      ],
      "env": {
        "SPOTIFY_CLIENT_ID": "your_client_id",
        "SPOTIFY_CLIENT_SECRET": "your_client_secret"
      }
    }
  }
}
```

## 📄 **License**

MIT License - see LICENSE file for details.

---

**Quick Setup:** Copy your Spotify credentials to `.env` and run `docker-compose up -d` 🎵
