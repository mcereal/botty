package config

import (
	"errors"
	"os"
)

// ICheckChannel is an interface for checking channel details
type ICheckChannel interface {
	checkChannelType() (string, error)
	// CheckPayload() (*bytes.Buffer, string)
}

// ChannelType holds the type of channel and webhhok value
type ChannelType struct {
	Channel, ChannelType string
}

func (c ChannelType) checkChannelType() (string, error) {
	url := c.Channel
	env := os.Getenv("ENVIRONMENT")
	var channelURL string
	if env == "development" || AppConfig.Application.Environment == "development" {

		if c.ChannelType == "discord" {
			channelURL = os.Getenv("DEV_DISCORD_CHANNEL_WEBHOOK_URL")
		} else {
			channelURL = os.Getenv("DEV_SLACK_CHANNEL_WEBHOOK_URL")
		}
		url = channelURL
	}
	if env == "staging" || AppConfig.Application.Environment == "development" {
		if c.ChannelType == "discord" {
			channelURL = os.Getenv("STAGING_DISCORD_CHANNEL_WEBHOOK_URL")
		} else {
			channelURL = os.Getenv("STAGING_SLACK_CHANNEL_WEBHOOK_URL")
		}
		url = channelURL
	}
	if env == "production" || AppConfig.Application.Environment == "development" {
		if c.ChannelType == "discord" {
			channelURL = os.Getenv("PRODUCTION_DISCORD_CHANNEL_WEBHOOK_URL")
		} else {
			channelURL = os.Getenv("PRODUCTION_SLACK_CHANNEL_WEBHOOK_URL")
		}
		url = channelURL
	}
	if url == "" {
		defaultURL := os.Getenv("DEFAULT_CHANNEL_WEBHOOK_URL")
		if len(defaultURL) != 0 {
			url = defaultURL
		} else {
			return "", errors.New("no webhook found")
		}
	}
	return url, nil
}

// CheckChannel checks for the correct webhook for the environment
func CheckChannel(o ICheckChannel) (string, error) {
	url, err := o.checkChannelType()
	return url, err
}
