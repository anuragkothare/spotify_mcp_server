package mcp

import (
	"encoding/json"
	"fmt"

	"github.com/anuragkothare/spotify_mcp_server/internal/spotify"
	"github.com/sirupsen/logrus"
)

type Server struct {
	spotifyClient *spotify.Client
	logger        *logrus.Logger
	tools         map[string]Tool
}

type MCPRequest struct {
	JSONRPC string          `json:"jsonrpc"`
	ID      interface{}     `json:"id"`
	Method  string          `json:"method"`
	Params  json.RawMessage `json:"params,omitempty"`
}

type MCPResponse struct {
	JSONRPC string      `json:"jsonrpc"`
	ID      interface{} `json:"id"`
	Result  interface{} `json:"result,omitempty"`
	Error   *MCPError   `json:"error,omitempty"`
}

type MCPError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

type Tool struct {
	Name        string                                            `json:"name"`
	Description string                                            `json:"description"`
	InputSchema interface{}                                       `json:"inputSchema"`
	Handler     func(params json.RawMessage) (interface{}, error) `json:"-"`
}

// ToolInfo represents tool information for JSON responses (without Handler)
type ToolInfo struct {
	Name        string      `json:"name"`
	Description string      `json:"description"`
	InputSchema interface{} `json:"inputSchema"`
}

func NewServer(spotifyClient *spotify.Client, logger *logrus.Logger) *Server {
	server := &Server{
		spotifyClient: spotifyClient,
		logger:        logger,
		tools:         make(map[string]Tool),
	}

	server.registerTools()
	return server
}

func (s *Server) registerTools() {
	s.tools["search_tracks"] = Tool{
		Name:        "search_tracks",
		Description: "Search for tracks on Spotify",
		InputSchema: map[string]interface{}{
			"type": "object",
			"properties": map[string]interface{}{
				"query": map[string]interface{}{
					"type":        "string",
					"description": "Search query for tracks",
				},
				"limit": map[string]interface{}{
					"type":        "integer",
					"description": "Maximum number of results (default: 10)",
					"minimum":     1,
					"maximum":     50,
				},
			},
			"required": []string{"query"},
		},
		Handler: s.handleSearchTracks,
	}

	s.tools["search_artists"] = Tool{
		Name:        "search_artists",
		Description: "Search for artists on Spotify",
		InputSchema: map[string]interface{}{
			"type": "object",
			"properties": map[string]interface{}{
				"query": map[string]interface{}{
					"type":        "string",
					"description": "Search query for artists",
				},
				"limit": map[string]interface{}{
					"type":        "integer",
					"description": "Maximum number of results (default: 10)",
					"minimum":     1,
					"maximum":     50,
				},
			},
			"required": []string{"query"},
		},
		Handler: s.handleSearchArtists,
	}

	s.tools["get_track"] = Tool{
		Name:        "get_track",
		Description: "Get detailed information about a specific track",
		InputSchema: map[string]interface{}{
			"type": "object",
			"properties": map[string]interface{}{
				"track_id": map[string]interface{}{
					"type":        "string",
					"description": "Spotify track ID",
				},
			},
			"required": []string{"track_id"},
		},
		Handler: s.handleGetTrack,
	}
}

func (s *Server) HandleRequest(req *MCPRequest) *MCPResponse {
	switch req.Method {
	case "tools/list":
		return s.handleListTools(req)
	case "tools/call":
		return s.handleToolCall(req)
	default:
		return &MCPResponse{
			JSONRPC: "2.0",
			ID:      req.ID,
			Error: &MCPError{
				Code:    -32601,
				Message: fmt.Sprintf("Method not found: %s", req.Method),
			},
		}
	}
}

func (s *Server) handleListTools(req *MCPRequest) *MCPResponse {
	tools := make([]ToolInfo, 0, len(s.tools))
	for _, tool := range s.tools {
		// Convert Tool to ToolInfo (without Handler function)
		tools = append(tools, ToolInfo{
			Name:        tool.Name,
			Description: tool.Description,
			InputSchema: tool.InputSchema,
		})
	}

	return &MCPResponse{
		JSONRPC: "2.0",
		ID:      req.ID,
		Result: map[string]interface{}{
			"tools": tools,
		},
	}
}

func (s *Server) handleToolCall(req *MCPRequest) *MCPResponse {
	var params struct {
		Name      string          `json:"name"`
		Arguments json.RawMessage `json:"arguments"`
	}

	if err := json.Unmarshal(req.Params, &params); err != nil {
		return &MCPResponse{
			JSONRPC: "2.0",
			ID:      req.ID,
			Error: &MCPError{
				Code:    -32602,
				Message: "Invalid parameters",
			},
		}
	}

	tool, exists := s.tools[params.Name]
	if !exists {
		return &MCPResponse{
			JSONRPC: "2.0",
			ID:      req.ID,
			Error: &MCPError{
				Code:    -32602,
				Message: fmt.Sprintf("Tool not found: %s", params.Name),
			},
		}
	}

	result, err := tool.Handler(params.Arguments)
	if err != nil {
		return &MCPResponse{
			JSONRPC: "2.0",
			ID:      req.ID,
			Error: &MCPError{
				Code:    -32603,
				Message: err.Error(),
			},
		}
	}

	return &MCPResponse{
		JSONRPC: "2.0",
		ID:      req.ID,
		Result: map[string]interface{}{
			"content": []map[string]interface{}{
				{
					"type": "text",
					"text": fmt.Sprintf("%+v", result),
				},
			},
		},
	}
}

// Tool handler methods
func (s *Server) handleSearchTracks(params json.RawMessage) (interface{}, error) {
	var args struct {
		Query string `json:"query"`
		Limit int    `json:"limit"`
	}

	if err := json.Unmarshal(params, &args); err != nil {
		return nil, err
	}

	if args.Limit == 0 {
		args.Limit = 10
	}

	return s.spotifyClient.SearchTracks(args.Query, args.Limit)
}

func (s *Server) handleSearchArtists(params json.RawMessage) (interface{}, error) {
	var args struct {
		Query string `json:"query"`
		Limit int    `json:"limit"`
	}

	if err := json.Unmarshal(params, &args); err != nil {
		return nil, err
	}

	if args.Limit == 0 {
		args.Limit = 10
	}

	return s.spotifyClient.SearchArtists(args.Query, args.Limit)
}

func (s *Server) handleGetTrack(params json.RawMessage) (interface{}, error) {
	var args struct {
		TrackID string `json:"track_id"`
	}

	if err := json.Unmarshal(params, &args); err != nil {
		return nil, err
	}

	return s.spotifyClient.GetTrack(args.TrackID)
}
