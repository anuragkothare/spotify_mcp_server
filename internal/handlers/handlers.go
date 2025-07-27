package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/anuragkothare/spotify_mcp_server/internal/mcp"
	"github.com/sirupsen/logrus"
)

type Handler struct {
	mcpServer *mcp.Server
	logger    *logrus.Logger
}

func NewHandler(mcpServer *mcp.Server, logger *logrus.Logger) *Handler {
	return &Handler{
		mcpServer: mcpServer,
		logger:    logger,
	}
}

func (h *Handler) HandleMCP(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req mcp.MCPRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.logger.Errorf("Failed to decode request: %v", err)
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	response := h.mcpServer.HandleRequest(&req)

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(response); err != nil {
		h.logger.Errorf("Failed to encode response: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
}

func (h *Handler) HandleHealth(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{
		"status": "healthy",
	})
}
