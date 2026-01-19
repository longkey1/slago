package slack

import (
	"github.com/slack-go/slack"
)

// Client wraps the Slack API client
type Client struct {
	api *slack.Client
}

// NewClient creates a new Slack client
func NewClient(token string) *Client {
	return &Client{
		api: slack.New(token),
	}
}

// API returns the underlying Slack API client
func (c *Client) API() *slack.Client {
	return c.api
}
