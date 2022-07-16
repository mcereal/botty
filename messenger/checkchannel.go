package messenger

import (
	"errors"
	"os"

	"github.com/mcereal/botty/config"
)

func (c ChannelType) checkChannelType() (string, error) {
	url := c.Channel
	env := os.Getenv("ENVIRONMENT")
	var channelURL string
	if env == "development" || config.AppConfig.Application.Environment == "development" {

		if c.ChannelType == "discord" {
			channelURL = os.Getenv("DEV_DISCORD_CHANNEL_WEBHOOK_URL")
		} else {
			channelURL = os.Getenv("DEV_SLACK_CHANNEL_WEBHOOK_URL")
		}
		url = channelURL
	}
	if env == "staging" || config.AppConfig.Application.Environment == "development" {
		if c.ChannelType == "discord" {
			channelURL = os.Getenv("STAGING_DISCORD_CHANNEL_WEBHOOK_URL")
		} else {
			channelURL = os.Getenv("STAGING_SLACK_CHANNEL_WEBHOOK_URL")
		}
		url = channelURL
	}
	if env == "production" || config.AppConfig.Application.Environment == "development" {
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

func CheckChannel(o IOpenPRs) (string, error) {
	url, err := o.checkChannelType()
	return url, err
}
