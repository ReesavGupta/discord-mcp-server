package tests

import (
	"testing"
	"time"

	"github.com/ReesavGupta/discord-mcp-server/internal/discord"
	"github.com/stretchr/testify/assert"
)

func TestMessageFilter(t *testing.T) {
	filter := discord.MessageFilter{
		ChannelID: "123456789",
		UserID:    "987654321",
		Content:   "test message",
		Limit:     50,
	}

	assert.Equal(t, "123456789", filter.ChannelID)
	assert.Equal(t, "987654321", filter.UserID)
	assert.Equal(t, "test message", filter.Content)
	assert.Equal(t, 50, filter.Limit)
}

func TestMessageFilterWithTime(t *testing.T) {
	now := time.Now()
	before := now.Add(-time.Hour)
	after := now.Add(-2 * time.Hour)

	filter := discord.MessageFilter{
		ChannelID: "123456789",
		Before:    &before,
		After:     &after,
		Limit:     100,
	}

	assert.Equal(t, "123456789", filter.ChannelID)
	assert.Equal(t, before, *filter.Before)
	assert.Equal(t, after, *filter.After)
	assert.Equal(t, 100, filter.Limit)
}
