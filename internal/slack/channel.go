package slack

import (
	"fmt"

	"github.com/slack-go/slack"
)

// GetChannelInfo gets information about a channel
func (c *Client) GetChannelInfo(channelID string) (*slack.Channel, error) {
	channel, err := c.api.GetConversationInfo(&slack.GetConversationInfoInput{
		ChannelID: channelID,
	})
	if err != nil {
		return nil, fmt.Errorf("conversations.info API error: %w", err)
	}
	return channel, nil
}

// GetChannelName gets the name of a channel, falling back to ID if not accessible
func (c *Client) GetChannelName(channelID string) string {
	channel, err := c.GetChannelInfo(channelID)
	if err != nil {
		return channelID
	}
	return channel.Name
}
