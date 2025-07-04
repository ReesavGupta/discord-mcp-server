package mcp

import (
	"encoding/json"
	"fmt"
	"io"
	"os"

	"github.com/ReesavGupta/discord-mcp-server/internal/auth"
	"github.com/ReesavGupta/discord-mcp-server/internal/config"
	"github.com/ReesavGupta/discord-mcp-server/internal/discord"
	"github.com/sirupsen/logrus"
)

type Server struct {
	config        *config.Config
	logger        *logrus.Logger
	authManager   *auth.AuthManager
	discordClient *discord.Client
	decoder       *json.Decoder
	encoder       *json.Encoder
}

func NewServer(cfg *config.Config, logger *logrus.Logger) (*Server, error) {
	// Initialize auth manager
	authManager, err := auth.NewAuthManager(
		cfg.Auth.JWTSecret,
		cfg.Auth.APIKeys,
		logger,
		cfg.Auth.EnableAudit,
		cfg.Auth.AuditLogPath,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create auth manager: %w", err)
	}

	// Initialize Discord client
	discordClient, err := discord.NewClient(cfg.Discord.BotToken, logger)
	if err != nil {
		return nil, fmt.Errorf("failed to create Discord client: %w", err)
	}

	return &Server{
		config:        cfg,
		logger:        logger,
		authManager:   authManager,
		discordClient: discordClient,
		decoder:       json.NewDecoder(os.Stdin),
		encoder:       json.NewEncoder(os.Stdout),
	}, nil
}

func (s *Server) Start() error {
	// Connect to Discord
	if err := s.discordClient.Connect(); err != nil {
		return fmt.Errorf("failed to connect to Discord: %w", err)
	}

	s.logger.Info("Discord MCP Server started")

	// Main message loop
	for {
		var request JSONRPCRequest
		if err := s.decoder.Decode(&request); err != nil {
			if err == io.EOF {
				break
			}
			s.logger.WithError(err).Error("Failed to decode request")
			s.sendError(nil, ParseError, "Failed to parse JSON")
			continue
		}

		s.logger.WithFields(logrus.Fields{
			"method": request.Method,
			"id":     request.ID,
		}).Debug("Received request")

		if request.JSONRPC != "2.0" {
			s.sendError(request.ID, InvalidRequest, "Only JSON-RPC 2.0 is supported")
			continue
		}

		if err := s.handleRequest(request); err != nil {
			s.logger.WithError(err).Error("Failed to handle request")
			s.sendError(request.ID, InternalError, err.Error())
		}
	}

	return s.discordClient.Disconnect()
}

func (s *Server) handleRequest(request JSONRPCRequest) error {
	switch request.Method {
	case "initialize":
		return s.handleInitialize(request)
	case "notifications/initialized", "initialized":
		s.logger.Info("Server initialized successfully")
		return nil
	case "tools/list":
		return s.handleToolsList(request)
	case "tools/call":
		return s.handleToolsCall(request)
	case "resources/list":
		return s.handleResourcesList(request)
	case "prompts/list":
		return s.handlePromptsList(request)
	case "cancelled":
		return s.handleCancelled(request)
	default:
		s.sendError(request.ID, MethodNotFound, "Method not implemented")
		return nil
	}
}

func (s *Server) handleInitialize(request JSONRPCRequest) error {
	response := JSONRPCResponse{
		JSONRPC: "2.0",
		ID:      request.ID,
		Result: InitializeResult{
			ProtocolVersion: s.config.MCP.ProtocolVersion,
			ServerInfo: ServerInfo{
				Name:    s.config.Server.Name,
				Version: s.config.Server.Version,
			},
			Capabilities: Capabilities{
				Tools:     map[string]interface{}{},
				Resources: map[string]interface{}{},
				Prompts:   map[string]interface{}{},
			},
		},
	}

	return s.sendResponse(response)
}

func (s *Server) sendResponse(response interface{}) error {
	return s.encoder.Encode(response)
}

func (s *Server) sendError(id interface{}, code int, message string) error {
	response := JSONRPCResponse{
		JSONRPC: "2.0",
		ID:      id,
		Error: &JSONRPCError{
			Code:    code,
			Message: message,
		},
	}
	return s.encoder.Encode(response)
}
