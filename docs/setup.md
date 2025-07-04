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
git clone https://github.com/ReesavGupta/discord-mcp-server.git
cd discord-mcp-server

# Install dependencies
make deps

# Build binary
make build

# The binary will be in bin/discord-mcp-server
```

# Method 2: Install from Releases

```bash
# Download the latest release for your platform
wget https://github.com/ReesavGupta/discord-mcp-server/releases/latest/download/discord-mcp-server-linux-amd64

# Make executable
chmod +x discord-mcp-server-linux-amd64

# Move to PATH
sudo mv discord-mcp-server-linux-amd64 /usr/local/bin/discord-mcp-server
```