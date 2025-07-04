package mcp

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/ReesavGupta/discord-mcp-server/internal/discord"
)

func (s *Server) handleToolsList(request JSONRPCRequest) error {
	tools := []Tool{
		{
			Name:        "send_message",
			Description: "Send a message to a Discord channel",
			InputSchema: json.RawMessage(`{
				"type": "object",
				"properties": {
					"channel_id": {
						"type": "string",
						"description": "The ID of the channel to send the message to"
					},
					"content": {
						"type": "string",
						"description": "The content of the message to send"
					}
				},
				"required": ["channel_id", "content"]
			}`),
		},
		{
			Name:        "get_messages",
			Description: "Retrieve message history from a Discord channel",
			InputSchema: json.RawMessage(`{
				"type": "object",
				"properties": {
					"channel_id": {
						"type": "string",
						"description": "The ID of the channel to get messages from"
					},
					"limit": {
						"type": "number",
						"description": "Number of messages to retrieve (default: 50, max: 100)",
						"default": 50
					}
				},
				"required": ["channel_id"]
			}`),
		},
		{
			Name:        "get_channel_info",
			Description: "Get information about a Discord channel",
			InputSchema: json.RawMessage(`{
				"type": "object",
				"properties": {
					"channel_id": {
						"type": "string",
						"description": "The ID of the channel to get info about"
					}
				},
				"required": ["channel_id"]
			}`),
		},
		{
			Name:        "search_messages",
			Description: "Search for messages in a Discord channel with filters",
			InputSchema: json.RawMessage(`{
				"type": "object",
				"properties": {
					"channel_id": {
						"type": "string",
						"description": "The ID of the channel to search in"
					},
					"content": {
						"type": "string",
						"description": "Content to search for (case-insensitive)"
					},
					"user_id": {
						"type": "string",
						"description": "Filter by user ID"
					},
					"before": {
						"type": "string",
						"description": "Search messages before this timestamp (ISO 8601)"
					},
					"after": {
						"type": "string",
						"description": "Search messages after this timestamp (ISO 8601)"
					},
					"limit": {
						"type": "number",
						"description": "Maximum number of messages to return (default: 50)",
						"default": 50
					}
				},
				"required": ["channel_id"]
			}`),
		},
		{
			Name:        "moderate_content",
			Description: "Perform moderation actions (delete messages, kick/ban users)",
			InputSchema: json.RawMessage(`{
				"type": "object",
				"properties": {
					"action": {
						"type": "string",
						"enum": ["delete_message", "kick_user", "ban_user"],
						"description": "The moderation action to perform"
					},
					"channel_id": {
						"type": "string",
						"description": "Channel ID (required for delete_message)"
					},
					"message_id": {
						"type": "string",
						"description": "Message ID (required for delete_message)"
					},
					"guild_id": {
						"type": "string",
						"description": "Guild ID (required for kick_user and ban_user)"
					},
					"user_id": {
						"type": "string",
						"description": "User ID (required for kick_user and ban_user)"
					},
					"reason": {
						"type": "string",
						"description": "Reason for the moderation action"
					},
					"delete_message_days": {
						"type": "number",
						"description": "Days of messages to delete when banning (0-7, default: 0)",
						"default": 0
					}
				},
				"required": ["action"]
			}`),
		},
	}

	response := JSONRPCResponse{
		JSONRPC: "2.0",
		ID:      request.ID,
		Result: ListToolsResult{
			Tools: tools,
		},
	}

	return s.sendResponse(response)
}

func (s *Server) handleToolsCall(request JSONRPCRequest) error {
	params, ok := request.Params.(map[string]interface{})
	if !ok {
		s.sendError(request.ID, InvalidParams, "Invalid parameters")
		return nil
	}

	toolName, ok := params["name"].(string)
	if !ok {
		s.sendError(request.ID, InvalidParams, "Invalid tool name")
		return nil
	}

	args, ok := params["arguments"].(map[string]interface{})
	if !ok {
		s.sendError(request.ID, InvalidParams, "Invalid arguments")
		return nil
	}

	var result CallToolResult
	var err error

	switch toolName {
	case "send_message":
		result, err = s.handleSendMessage(args)
	case "get_messages":
		result, err = s.handleGetMessages(args)
	case "get_channel_info":
		result, err = s.handleGetChannelInfo(args)
	case "search_messages":
		result, err = s.handleSearchMessages(args)
	case "moderate_content":
		result, err = s.handleModerateContent(args)
	default:
		s.sendError(request.ID, MethodNotFound, "Unknown tool")
		return nil
	}

	if err != nil {
		s.sendError(request.ID, InternalError, err.Error())
		return nil
	}

	response := JSONRPCResponse{
		JSONRPC: "2.0",
		ID:      request.ID,
		Result:  result,
	}

	return s.sendResponse(response)
}

func (s *Server) handleSendMessage(args map[string]interface{}) (CallToolResult, error) {
	channelID, ok := args["channel_id"].(string)
	if !ok {
		return CallToolResult{}, fmt.Errorf("channel_id is required")
	}

	content, ok := args["content"].(string)
	if !ok {
		return CallToolResult{}, fmt.Errorf("content is required")
	}

	message, err := s.discordClient.SendMessage(channelID, content)
	if err != nil {
		return CallToolResult{}, fmt.Errorf("failed to send message: %w", err)
	}

	return CallToolResult{
		Content: []ToolContent{
			{
				Type: "text",
				Text: fmt.Sprintf("Message sent successfully. Message ID: %s", message.ID),
			},
		},
	}, nil
}

func (s *Server) handleGetMessages(args map[string]interface{}) (CallToolResult, error) {
	channelID, ok := args["channel_id"].(string)
	if !ok {
		return CallToolResult{}, fmt.Errorf("channel_id is required")
	}

	limit := 50
	if l, ok := args["limit"].(float64); ok {
		limit = int(l)
		if limit > 100 {
			limit = 100
		}
	}

	messages, err := s.discordClient.GetMessages(channelID, limit)
	if err != nil {
		return CallToolResult{}, fmt.Errorf("failed to get messages: %w", err)
	}

	var messageTexts []string
	for _, msg := range messages {
		messageTexts = append(messageTexts, fmt.Sprintf("[%s] %s: %s",
			msg.Timestamp.Format(time.RFC3339),
			msg.Author.Username,
			msg.Content))
	}

	resultText := fmt.Sprintf("Retrieved %d messages from channel %s:\n%s",
		len(messages), channelID, fmt.Sprintf("%v", messageTexts))

	return CallToolResult{
		Content: []ToolContent{
			{
				Type: "text",
				Text: resultText,
			},
		},
	}, nil
}

func (s *Server) handleGetChannelInfo(args map[string]interface{}) (CallToolResult, error) {
	channelID, ok := args["channel_id"].(string)
	if !ok {
		return CallToolResult{}, fmt.Errorf("channel_id is required")
	}

	channel, err := s.discordClient.GetChannelInfo(channelID)
	if err != nil {
		return CallToolResult{}, fmt.Errorf("failed to get channel info: %w", err)
	}

	resultText := fmt.Sprintf("Channel Info:\nID: %s\nName: %s\nType: %d\nGuild ID: %s\nPosition: %d",
		channel.ID, channel.Name, channel.Type, channel.GuildID, channel.Position)

	return CallToolResult{
		Content: []ToolContent{
			{
				Type: "text",
				Text: resultText,
			},
		},
	}, nil
}

func (s *Server) handleSearchMessages(args map[string]interface{}) (CallToolResult, error) {
	channelID, ok := args["channel_id"].(string)
	if !ok {
		return CallToolResult{}, fmt.Errorf("channel_id is required")
	}

	filter := discord.MessageFilter{
		ChannelID: channelID,
		Limit:     50,
	}

	if content, ok := args["content"].(string); ok {
		filter.Content = content
	}

	if userID, ok := args["user_id"].(string); ok {
		filter.UserID = userID
	}

	if limit, ok := args["limit"].(float64); ok {
		filter.Limit = int(limit)
	}

	if before, ok := args["before"].(string); ok {
		if t, err := time.Parse(time.RFC3339, before); err == nil {
			filter.Before = &t
		}
	}

	if after, ok := args["after"].(string); ok {
		if t, err := time.Parse(time.RFC3339, after); err == nil {
			filter.After = &t
		}
	}

	messages, err := s.discordClient.SearchMessages(filter)
	if err != nil {
		return CallToolResult{}, fmt.Errorf("failed to search messages: %w", err)
	}

	var messageTexts []string
	for _, msg := range messages {
		messageTexts = append(messageTexts, fmt.Sprintf("[%s] %s: %s",
			msg.Timestamp.Format(time.RFC3339),
			msg.Author.Username,
			msg.Content))
	}

	resultText := fmt.Sprintf("Found %d messages matching criteria:\n%s",
		len(messages), fmt.Sprintf("%v", messageTexts))

	return CallToolResult{
		Content: []ToolContent{
			{
				Type: "text",
				Text: resultText,
			},
		},
	}, nil
}

func (s *Server) handleModerateContent(args map[string]interface{}) (CallToolResult, error) {
	action, ok := args["action"].(string)
	if !ok {
		return CallToolResult{}, fmt.Errorf("action is required")
	}

	reason := "No reason provided"
	if r, ok := args["reason"].(string); ok {
		reason = r
	}

	switch action {
	case "delete_message":
		channelID, ok := args["channel_id"].(string)
		if !ok {
			return CallToolResult{}, fmt.Errorf("channel_id is required for delete_message")
		}

		messageID, ok := args["message_id"].(string)
		if !ok {
			return CallToolResult{}, fmt.Errorf("message_id is required for delete_message")
		}

		err := s.discordClient.DeleteMessage(channelID, messageID)
		if err != nil {
			return CallToolResult{}, fmt.Errorf("failed to delete message: %w", err)
		}

		return CallToolResult{
			Content: []ToolContent{
				{
					Type: "text",
					Text: fmt.Sprintf("Message %s deleted successfully", messageID),
				},
			},
		}, nil

	case "kick_user":
		guildID, ok := args["guild_id"].(string)
		if !ok {
			return CallToolResult{}, fmt.Errorf("guild_id is required for kick_user")
		}

		userID, ok := args["user_id"].(string)
		if !ok {
			return CallToolResult{}, fmt.Errorf("user_id is required for kick_user")
		}

		err := s.discordClient.KickUser(guildID, userID, reason)
		if err != nil {
			return CallToolResult{}, fmt.Errorf("failed to kick user: %w", err)
		}

		return CallToolResult{
			Content: []ToolContent{
				{
					Type: "text",
					Text: fmt.Sprintf("User %s kicked successfully. Reason: %s", userID, reason),
				},
			},
		}, nil

	case "ban_user":
		guildID, ok := args["guild_id"].(string)
		if !ok {
			return CallToolResult{}, fmt.Errorf("guild_id is required for ban_user")
		}

		userID, ok := args["user_id"].(string)
		if !ok {
			return CallToolResult{}, fmt.Errorf("user_id is required for ban_user")
		}

		deleteMessageDays := 0
		if d, ok := args["delete_message_days"].(float64); ok {
			deleteMessageDays = int(d)
			if deleteMessageDays > 7 {
				deleteMessageDays = 7
			}
		}

		err := s.discordClient.BanUser(guildID, userID, reason, deleteMessageDays)
		if err != nil {
			return CallToolResult{}, fmt.Errorf("failed to ban user: %w", err)
		}

		return CallToolResult{
			Content: []ToolContent{
				{
					Type: "text",
					Text: fmt.Sprintf("User %s banned successfully. Reason: %s", userID, reason),
				},
			},
		}, nil

	default:
		return CallToolResult{}, fmt.Errorf("unknown moderation action: %s", action)
	}
}

func (s *Server) handleResourcesList(request JSONRPCRequest) error {
	response := JSONRPCResponse{
		JSONRPC: "2.0",
		ID:      request.ID,
		Result: map[string]interface{}{
			"resources": []interface{}{},
		},
	}

	return s.sendResponse(response)
}

func (s *Server) handlePromptsList(request JSONRPCRequest) error {
	response := JSONRPCResponse{
		JSONRPC: "2.0",
		ID:      request.ID,
		Result: map[string]interface{}{
			"prompts": []interface{}{},
		},
	}

	return s.sendResponse(response)
}

func (s *Server) handleCancelled(request JSONRPCRequest) error {
	if params, ok := request.Params.(map[string]interface{}); ok {
		s.logger.WithFields(map[string]interface{}{
			"request_id": params["requestId"],
			"reason":     params["reason"],
		}).Info("Received cancellation notification")
	}
	return nil
}
