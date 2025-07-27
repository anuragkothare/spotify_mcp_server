package mcp

import (
	"encoding/json"
)

// Resource represents an MCP resource
type Resource struct {
	URI         string                 `json:"uri"`
	Name        string                 `json:"name"`
	Description string                 `json:"description,omitempty"`
	MimeType    string                 `json:"mimeType,omitempty"`
	Metadata    map[string]interface{} `json:"metadata,omitempty"`
}

// Request represents a generic MCP request
type Request struct {
	Method string          `json:"method"`
	Params json.RawMessage `json:"params,omitempty"`
	ID     interface{}     `json:"id,omitempty"`
}

// Response represents a generic MCP response
type Response struct {
	Result interface{} `json:"result,omitempty"`
	Error  *Error      `json:"error,omitempty"`
	ID     interface{} `json:"id,omitempty"`
}

// Error represents an MCP error
type Error struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

// InitializeRequest represents an MCP initialize request
type InitializeRequest struct {
	ProtocolVersion string                 `json:"protocolVersion"`
	Capabilities    map[string]interface{} `json:"capabilities"`
	ClientInfo      ClientInfo             `json:"clientInfo"`
}

// InitializeResponse represents an MCP initialize response
type InitializeResponse struct {
	ProtocolVersion string                 `json:"protocolVersion"`
	Capabilities    map[string]interface{} `json:"capabilities"`
	ServerInfo      ServerInfo             `json:"serverInfo"`
}

// ClientInfo represents client information
type ClientInfo struct {
	Name    string `json:"name"`
	Version string `json:"version"`
}

// ServerInfo represents server information
type ServerInfo struct {
	Name    string `json:"name"`
	Version string `json:"version"`
}

// ListToolsRequest represents a list tools request
type ListToolsRequest struct {
	Cursor string `json:"cursor,omitempty"`
}

// ListToolsResponse represents a list tools response
type ListToolsResponse struct {
	Tools      []*Tool `json:"tools"`
	NextCursor string  `json:"nextCursor,omitempty"`
}

// CallToolRequest represents a call tool request
type CallToolRequest struct {
	Name      string          `json:"name"`
	Arguments json.RawMessage `json:"arguments,omitempty"`
}

// CallToolResponse represents a call tool response
type CallToolResponse struct {
	Content []Content `json:"content"`
	IsError bool      `json:"isError,omitempty"`
}

// Content represents MCP content
type Content struct {
	Type string      `json:"type"`
	Text string      `json:"text,omitempty"`
	Data interface{} `json:"data,omitempty"`
}

// ListResourcesRequest represents a list resources request
type ListResourcesRequest struct {
	Cursor string `json:"cursor,omitempty"`
}

// ListResourcesResponse represents a list resources response
type ListResourcesResponse struct {
	Resources  []*Resource `json:"resources"`
	NextCursor string      `json:"nextCursor,omitempty"`
}

// ReadResourceRequest represents a read resource request
type ReadResourceRequest struct {
	URI string `json:"uri"`
}

// ReadResourceResponse represents a read resource response
type ReadResourceResponse struct {
	Contents []ResourceContent `json:"contents"`
}

// ResourceContent represents resource content
type ResourceContent struct {
	URI      string `json:"uri"`
	MimeType string `json:"mimeType,omitempty"`
	Text     string `json:"text,omitempty"`
	Blob     []byte `json:"blob,omitempty"`
}

// Common MCP error codes
const (
	ErrorCodeInvalidRequest      = -32600
	ErrorCodeMethodNotFound      = -32601
	ErrorCodeInvalidParams       = -32602
	ErrorCodeInternalError       = -32603
	ErrorCodeParseError          = -32700
	ErrorCodeResourceNotFound    = -32001
	ErrorCodeResourceUnavailable = -32002
	ErrorCodeToolNotFound        = -32003
	ErrorCodeToolExecutionError  = -32004
)

// NewError creates a new MCP error
func NewError(code int, message string, data interface{}) *Error {
	return &Error{
		Code:    code,
		Message: message,
		Data:    data,
	}
}

// NewTextContent creates new text content
func NewTextContent(text string) Content {
	return Content{
		Type: "text",
		Text: text,
	}
}

// NewDataContent creates new data content
func NewDataContent(data interface{}) Content {
	return Content{
		Type: "data",
		Data: data,
	}
}

// NewSuccessResponse creates a successful response
func NewSuccessResponse(result interface{}, id interface{}) *Response {
	return &Response{
		Result: result,
		ID:     id,
	}
}

// NewErrorResponse creates an error response
func NewErrorResponse(err *Error, id interface{}) *Response {
	return &Response{
		Error: err,
		ID:    id,
	}
}
