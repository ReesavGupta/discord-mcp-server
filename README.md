# Discord MCP Server

A blazingly fast, robust and extensible Model Context Protocol (MCP) server for Discord, written in Go. This project enables advanced automation, moderation, and integration with Discord servers via a standardized JSON-RPC 2.0 interface. It is designed for use with Claude Desktop and other MCP-compatible clients.

---

## Features

- **Discord Bot Integration**: Connects to Discord as a bot, supporting message sending, retrieval, moderation, and channel info.
- **MCP Protocol**: Implements the Model Context Protocol (MCP) for standardized tool, resource, and prompt management.
- **JSON-RPC 2.0**: Communicates via JSON-RPC 2.0 over stdio for easy integration.
- **Authentication**: Supports JWT and API key authentication, with optional audit logging.
- **Extensible Tools**: Built-in tools for sending messages, searching, moderation, and more.
- **Configurable**: YAML-based configuration with environment variable overrides.
- **Logging**: Structured logging with support for file and JSON/text formats.
- **Testing**: Includes unit tests and coverage reporting.
- **Multi-tenancy**: Supports multiple Discord bots and servers.
- **Rate Limiting & Monitoring**: Optional support for rate limiting and metrics.

---

## Table of Contents
- [Features](#features)
- [Setup](#setup)
- [Configuration](#configuration)
- [Usage](#usage)
- [Testing](#testing)
- [Troubleshooting](#troubleshooting)
- [Security Best Practices](#security-best-practices)
- [Contributing](#contributing)
- [License](#license)

---

## Setup

See [`docs/setup.md`](docs/setup.md) for a detailed step-by-step guide, including Discord bot creation, permissions, and environment setup.

### Quick Start

```bash
# Clone repository
git clone https://github.com/yourusername/discord-mcp-server.git
cd discord-mcp-server

# Install dependencies
make deps

# Build the server
make build

# Copy and edit config
cp configs/config.yaml.example configs/config.yaml
# Or edit configs/config.yaml directly

# Run the server
./bin/discord-mcp-server -config configs/config.yaml
```

---

## Configuration

- **Main config:** `configs/config.yaml` (YAML, supports env vars)
- **Environment:** `.env` file or system env vars
- **Claude Desktop:** See `configs/claude_desktop_config.json` for integration

Example config:
```yaml
server:
  name: "discord-mcp-server"
  version: "1.0.0"
  environment: "production"
discord:
  bot_token: "${DISCORD_BOT_TOKEN}"
  guild_id: "${DISCORD_GUILD_ID}"
auth:
  jwt_secret: "${JWT_SECRET}"
  api_keys:
    - "${API_KEY_1}"
  enable_audit: true
logging:
  level: "info"
  format: "json"
  file_path: "logs/server.log"
mcp:
  protocol_version: "2024-11-05"
  transport: "stdio"
  debug: false
```

---

## Usage

The server runs as a background process and communicates via stdio (for Claude Desktop) or can be extended for other transports.

### Supported Tools (via MCP)
- `send_message`: Send a message to a Discord channel
- `get_messages`: Retrieve message history
- `get_channel_info`: Get channel metadata
- `search_messages`: Search messages with filters (content, user, time)
- `moderate_content`: Delete messages, kick/ban users

See the MCP tool schemas in [`internal/mcp/handlers.go`](internal/mcp/handlers.go) for details.

---

## Testing

```bash
# Run all tests
make test

# Run with coverage
make test-coverage
# Coverage report: coverage.html
```

---

## Troubleshooting

- **Bot not responding**: Check bot token, permissions, and logs
- **Permission denied**: Ensure bot has required Discord permissions
- **Connection failed**: Verify guild ID and bot invite
- **Authentication failed**: Check API keys and JWT secret
- **Debug logs**: Set `logging.level: debug` in config

---

## Security Best Practices
- Never commit secrets to version control
- Use environment variables for sensitive data
- Rotate API keys regularly
- Enable audit logging in production
- Use strong JWT secrets (32+ chars)
- Limit bot permissions to minimum required
- Monitor logs for suspicious activity

---

## Contributing

Contributions are welcome! Please open issues or pull requests. See [`docs/setup.md`](docs/setup.md) for development setup.

---

## License

MIT License. See [LICENSE](LICENSE) for details.

---

## Acknowledgements
- [discordgo](https://github.com/bwmarrin/discordgo)
- [sirupsen/logrus](https://github.com/sirupsen/logrus)
- [golang-jwt/jwt](https://github.com/golang-jwt/jwt)
- [Claude Desktop](https://github.com/anthropics/claude-desktop)
