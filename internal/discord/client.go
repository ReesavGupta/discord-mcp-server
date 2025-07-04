package discord

import (
	"fmt"
	"strings"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/sirupsen/logrus"
)

type Client struct {
	session *discordgo.Session
	logger  *logrus.Logger
	guildID string
}

type MessageFilter struct {
	ChannelID string
	UserID    string
	Content   string
	Before    *time.Time
	After     *time.Time
	Limit     int
}

func NewClient(token string, logger *logrus.Logger) (*Client, error) {
	session, err := discordgo.New("Bot " + token)
	if err != nil {
		return nil, fmt.Errorf("failed to create Discord session: %w", err)
	}

	return &Client{
		session: session,
		logger:  logger,
	}, nil
}

func (c *Client) Connect() error {
	c.session.AddHandler(func(s *discordgo.Session, r *discordgo.Ready) {
		c.logger.Info("Discord bot is ready")
	})

	return c.session.Open()
}

func (c *Client) Disconnect() error {
	return c.session.Close()
}

func (c *Client) SendMessage(channelID, content string) (*discordgo.Message, error) {
	c.logger.WithFields(logrus.Fields{
		"channel_id": channelID,
		"content":    content,
	}).Info("Sending message")

	return c.session.ChannelMessageSend(channelID, content)
}

func (c *Client) GetMessages(channelID string, limit int) ([]*discordgo.Message, error) {
	c.logger.WithFields(logrus.Fields{
		"channel_id": channelID,
		"limit":      limit,
	}).Info("Fetching messages")

	return c.session.ChannelMessages(channelID, limit, "", "", "")
}

func (c *Client) GetChannelInfo(channelID string) (*discordgo.Channel, error) {
	c.logger.WithFields(logrus.Fields{
		"channel_id": channelID,
	}).Info("Fetching channel info")

	return c.session.Channel(channelID)
}

func (c *Client) SearchMessages(filter MessageFilter) ([]*discordgo.Message, error) {
	messages, err := c.GetMessages(filter.ChannelID, filter.Limit)
	if err != nil {
		return nil, err
	}

	var filtered []*discordgo.Message
	for _, msg := range messages {
		if c.matchesFilter(msg, filter) {
			filtered = append(filtered, msg)
		}
	}

	return filtered, nil
}

func (c *Client) matchesFilter(msg *discordgo.Message, filter MessageFilter) bool {
	if filter.UserID != "" && msg.Author.ID != filter.UserID {
		return false
	}

	if filter.Content != "" && !strings.Contains(strings.ToLower(msg.Content), strings.ToLower(filter.Content)) {
		return false
	}

	msgTime, err := msg.Timestamp.Parse()
	if err != nil {
		return false
	}

	if filter.Before != nil && msgTime.After(*filter.Before) {
		return false
	}

	if filter.After != nil && msgTime.Before(*filter.After) {
		return false
	}

	return true
}

func (c *Client) DeleteMessage(channelID, messageID string) error {
	c.logger.WithFields(logrus.Fields{
		"channel_id": channelID,
		"message_id": messageID,
	}).Info("Deleting message")

	return c.session.ChannelMessageDelete(channelID, messageID)
}

func (c *Client) KickUser(guildID, userID, reason string) error {
	c.logger.WithFields(logrus.Fields{
		"guild_id": guildID,
		"user_id":  userID,
		"reason":   reason,
	}).Info("Kicking user")

	return c.session.GuildMemberDeleteWithReason(guildID, userID, reason)
}

func (c *Client) BanUser(guildID, userID, reason string, deleteMessageDays int) error {
	c.logger.WithFields(logrus.Fields{
		"guild_id":            guildID,
		"user_id":             userID,
		"reason":              reason,
		"delete_message_days": deleteMessageDays,
	}).Info("Banning user")

	return c.session.GuildBanCreateWithReason(guildID, userID, reason, deleteMessageDays)
}
