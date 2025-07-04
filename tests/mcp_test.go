package tests

import (
	"encoding/json"
	"testing"

	"github.com/ReesavGupta/discord-mcp-server/internal/mcp"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestJSONRPCTypes(t *testing.T) {
	t.Run("JSONRPCRequest", func(t *testing.T) {
		jsonStr := `{
			"jsonrpc": "2.0",
			"id": 1,
			"method": "initialize",
			"params": {"test": "value"}
		}`

		var req mcp.JSONRPCRequest
		err := json.Unmarshal([]byte(jsonStr), &req)
		require.NoError(t, err)

		assert.Equal(t, "2.0", req.JSONRPC)
		assert.Equal(t, float64(1), req.ID)
		assert.Equal(t, "initialize", req.Method)
	})

	t.Run("JSONRPCResponse", func(t *testing.T) {
		resp := mcp.JSONRPCResponse{
			JSONRPC: "2.0",
			ID:      1,
			Result:  map[string]string{"status": "ok"},
		}

		data, err := json.Marshal(resp)
		require.NoError(t, err)

		assert.Contains(t, string(data), `"jsonrpc":"2.0"`)
		assert.Contains(t, string(data), `"id":1`)
		assert.Contains(t, string(data), `"result"`)
	})
}

func TestInitializeResult(t *testing.T) {
	result := mcp.InitializeResult{
		ProtocolVersion: "2024-11-05",
		ServerInfo: mcp.ServerInfo{
			Name:    "discord-mcp-server",
			Version: "1.0.0",
		},
		Capabilities: mcp.Capabilities{
			Tools: map[string]interface{}{},
		},
	}

	data, err := json.Marshal(result)
	require.NoError(t, err)

	assert.Contains(t, string(data), `"protocolVersion":"2024-11-05"`)
	assert.Contains(t, string(data), `"name":"discord-mcp-server"`)
}

func TestToolDefinitions(t *testing.T) {
	tool := mcp.Tool{
		Name:        "send_message",
		Description: "Send a message to a Discord channel",
		InputSchema: json.RawMessage(`{"type":"object"}`),
	}

	data, err := json.Marshal(tool)
	require.NoError(t, err)

	assert.Contains(t, string(data), `"name":"send_message"`)
	assert.Contains(t, string(data), `"description":"Send a message to a Discord channel"`)
	assert.Contains(t, string(data), `"inputSchema":{"type":"object"}`)
}
