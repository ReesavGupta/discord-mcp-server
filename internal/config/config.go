package config

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

type Config struct {
	Server  ServerConfig  `yaml:"server"`
	Discord DiscordConfig `yaml:"discord"`
	Auth    AuthConfig    `yaml:"auth"`
	Logging LoggingConfig `yaml:"logging"`
	MCP     MCPConfig     `yaml:"mcp"`
}

type ServerConfig struct {
	Name        string `yaml:"name"`
	Version     string `yaml:"version"`
	Environment string `yaml:"environment"`
}

type DiscordConfig struct {
	BotToken     string   `yaml:"bot_token"`
	GuildID      string   `yaml:"guild_id"`
	AllowedRoles []string `yaml:"allowed_roles"`
}

type AuthConfig struct {
	JWTSecret    string   `yaml:"jwt_secret"`
	APIKeys      []string `yaml:"api_keys"`
	EnableAudit  bool     `yaml:"enable_audit"`
	AuditLogPath string   `yaml:"audit_log_path"`
}

type LoggingConfig struct {
	Level    string `yaml:"level"`
	Format   string `yaml:"format"`
	FilePath string `yaml:"file_path"`
}

type MCPConfig struct {
	ProtocolVersion string `yaml:"protocol_version"`
	Transport       string `yaml:"transport"`
	Debug           bool   `yaml:"debug"`
}

func LoadConfig(path string) (*Config, error) {
	config := &Config{}

	// Set defaults
	config.Server.Name = "discord-mcp-server"
	config.Server.Version = "1.0.0"
	config.Server.Environment = "development"
	config.MCP.ProtocolVersion = "2024-11-05"
	config.MCP.Transport = "stdio"
	config.Logging.Level = "info"
	config.Logging.Format = "json"

	// Load from file if exists
	if _, err := os.Stat(path); err == nil {
		file, err := os.ReadFile(path)
		if err != nil {
			return nil, fmt.Errorf("failed to read config file: %w", err)
		}

		if err := yaml.Unmarshal(file, config); err != nil {
			return nil, fmt.Errorf("failed to parse config file: %w", err)
		}
	}

	// Override with environment variables
	if token := os.Getenv("DISCORD_BOT_TOKEN"); token != "" {
		config.Discord.BotToken = token
	}

	if secret := os.Getenv("JWT_SECRET"); secret != "" {
		config.Auth.JWTSecret = secret
	}

	if keys := os.Getenv("API_KEYS"); keys != "" {
		config.Auth.APIKeys = append(config.Auth.APIKeys, keys)
	}

	return config, nil
}
