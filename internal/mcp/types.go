package mcp

import (
	"encoding/json"
	"time"
)

// JSON-RPC 2.0 base types
type JSONRPCRequest struct {
	JSONRPC string      `json:"jsonrpc"`
	ID      interface{} `json:"id"`
	Method  string      `json:"method"`
	Params  interface{} `json:"params,omitempty"`
}

type JSONRPCResponse struct {
	JSONRPC string        `json:"jsonrpc"`
	ID      interface{}   `json:"id"`
	Result  interface{}   `json:"result,omitempty"`
	Error   *JSONRPCError `json:"error,omitempty"`
}

type JSONRPCError struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

// MCP Protocol specific types
type InitializeResult struct {
	ProtocolVersion string       `json:"protocolVersion"`
	ServerInfo      ServerInfo   `json:"serverInfo"`
	Capabilities    Capabilities `json:"capabilities"`
}

type ServerInfo struct {
	Name    string `json:"name"`
	Version string `json:"version"`
}

type Capabilities struct {
	Tools     map[string]interface{} `json:"tools"`
	Resources map[string]interface{} `json:"resources"`
	Prompts   map[string]interface{} `json:"prompts"`
}

type Tool struct {
	Name        string          `json:"name"`
	Description string          `json:"description"`
	InputSchema json.RawMessage `json:"inputSchema"`
}

type ListToolsResult struct {
	Tools []Tool `json:"tools"`
}

type CallToolResult struct {
	Content []ToolContent `json:"content"`
}

type ToolContent struct {
	Type string `json:"type"`
	Text string `json:"text"`
}

// Discord-specific types
type DiscordMessage struct {
	ID        string    `json:"id"`
	Content   string    `json:"content"`
	Author    string    `json:"author"`
	Timestamp time.Time `json:"timestamp"`
	ChannelID string    `json:"channel_id"`
}

type DiscordChannel struct {
	ID       string `json:"id"`
	Name     string `json:"name"`
	Type     int    `json:"type"`
	GuildID  string `json:"guild_id"`
	Position int    `json:"position"`
}

// Error codes
const (
	ParseError     = -32700
	InvalidRequest = -32600
	MethodNotFound = -32601
	InvalidParams  = -32602
	InternalError  = -32603
)
