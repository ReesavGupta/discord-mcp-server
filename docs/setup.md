# Setup Guide

## Prerequisites

1. **Go 1.21+** - Install from [golang.org](https://golang.org)
2. **Discord Bot Token** - Create a bot on [Discord Developer Portal](https://discord.com/developers/applications)
3. **Discord Guild ID** - Enable Developer Mode in Discord and copy your server ID

## Discord Bot Setup

### 1. Create Discord Application

1. Go to [Discord Developer Portal](https://discord.com/developers/applications)
2. Click "New Application"
3. Name your application (e.g., "MCP Discord Bot")
4. Save the Application ID

### 2. Create Bot

1. Go to the "Bot" section
2. Click "Add Bot"
3. Copy the Bot Token (keep this secret!)
4. Enable these Privileged Gateway Intents:
   - Server Members Intent
   - Message Content Intent

### 3. Set Bot Permissions

Required permissions:
- Read Messages/View Channels
- Send Messages
- Read Message History
- Manage Messages (for moderation)
- Kick Members (for moderation)
- Ban Members (for moderation)

### 4. Invite Bot to Server

1. Go to "OAuth2" → "URL Generator"
2. Select "bot" scope
3. Select the permissions above
4. Copy the generated URL
5. Open URL in browser and invite bot to your server

## Installation

### Method 1: Build from Source

```bash
# Clone repository
git clone https://github.com/yourusername/discord-mcp-server.git
cd discord-mcp-server

# Install dependencies
make deps

# Build binary
make build

# The binary will be in bin/discord-mcp-server
```

### Method 2: Install from Releases

```bash
# Download the latest release for your platform
wget https://github.com/yourusername/discord-mcp-server/releases/latest/download/discord-mcp-server-linux-amd64

# Make executable
chmod +x discord-mcp-server-linux-amd64

# Move to PATH
sudo mv discord-mcp-server-linux-amd64 /usr/local/bin/discord-mcp-server
```

## Configuration

### 1. Environment Variables

Create a .env file or set these environment variables:

```bash
export DISCORD_BOT_TOKEN="your_bot_token_here"
export DISCORD_GUILD_ID="your_guild_id_here"
export JWT_SECRET="your_jwt_secret_here"
export API_KEY_1="your_api_key_1_here"
```

### 2. Configuration File

Copy configs/config.yaml and modify as needed:

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
```

### 3. Claude Desktop Configuration

Update your claude_desktop_config.json:

```json
{
  "mcpServers": {
    "discord-mcp-server": {
      "command": "/usr/local/bin/discord-mcp-server",
      "args": ["-config", "/path/to/config.yaml"],
      "env": {
        "DISCORD_BOT_TOKEN": "your_bot_token",
        "DISCORD_GUILD_ID": "your_guild_id",
        "JWT_SECRET": "your_jwt_secret",
        "API_KEY_1": "your_api_key"
      }
    }
  }
}
```

## Testing

```bash
# Run tests
make test

# Run tests with coverage
make test-coverage

# View coverage report
open coverage.html
```

## Troubleshooting

### Common Issues

- **Bot not responding**: Check bot token and permissions
- **Permission denied**: Ensure bot has required permissions in Discord
- **Connection failed**: Verify guild ID and bot is in the server
- **Authentication failed**: Check API keys and JWT secret

### Debug Mode

Enable debug logging:

```yaml
logging:
  level: "debug"
  
mcp:
  debug: true
```

### Logs

Check logs for errors:

```bash
# View server logs
tail -f logs/server.log

# View audit logs
tail -f logs/audit.log
```

---

# Configuration Reference

## Configuration File Structure

The server uses YAML configuration with environment variable substitution.

### Server Configuration

```yaml
server:
  name: "discord-mcp-server"      # Server name
  version: "1.0.0"                # Server version
  environment: "production"       # Environment (development/production)
```

### Discord Configuration

```yaml
discord:
  bot_token: "${DISCORD_BOT_TOKEN}"     # Discord bot token
  guild_id: "${DISCORD_GUILD_ID}"       # Discord server ID
  allowed_roles:                         # Roles allowed to use the bot
    - "Admin"
    - "Moderator"
```

### Authentication Configuration

```yaml
auth:
  jwt_secret: "${JWT_SECRET}"           # JWT signing secret
  api_keys:                             # Valid API keys
    - "${API_KEY_1}"
    - "${API_KEY_2}"
  enable_audit: true                    # Enable audit logging
  audit_log_path: "logs/audit.log"     # Audit log file path
```

### Logging Configuration

```yaml
logging:
  level: "info"                         # Log level (debug/info/warn/error)
  format: "json"                        # Log format (json/text)
  file_path: "logs/server.log"         # Log file path
```

### MCP Configuration

```yaml
mcp:
  protocol_version: "2024-11-05"       # MCP protocol version
  transport: "stdio"                   # Transport method
  debug: false                         # Enable debug mode
```

## Environment Variables

All configuration values can be overridden with environment variables:

| Variable | Description | Example |
|----------|-------------|---------|
| DISCORD_BOT_TOKEN | Discord bot token | MTA5NzE2NTI3... |
| DISCORD_GUILD_ID | Discord server ID | 123456789012345678 |
| JWT_SECRET | JWT signing secret | your-secret-key |
| API_KEY_1 | API key for authentication | sk-1234567890abcdef |

## Security Best Practices

- Never commit secrets to version control
- Use environment variables for sensitive data
- Rotate API keys regularly
- Enable audit logging in production
- Use strong JWT secrets (32+ characters)
- Limit bot permissions to minimum required
- Monitor logs for suspicious activity

## Multi-tenancy Support

The server supports multiple Discord bots:

```yaml
# Multiple bot configurations
discord:
  bots:
    - name: "bot1"
      token: "${BOT1_TOKEN}"
      guild_id: "${BOT1_GUILD_ID}"
    - name: "bot2"
      token: "${BOT2_TOKEN}"
      guild_id: "${BOT2_GUILD_ID}"
```

## Rate Limiting

Configure rate limits to prevent abuse:

```yaml
rate_limiting:
  enabled: true
  requests_per_minute: 60
  burst_size: 10
  per_user: true
```

## Monitoring

Enable monitoring and metrics:

```yaml
monitoring:
  enabled: true
  metrics_port: 8080
  health_check_port: 8081
```