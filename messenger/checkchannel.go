package messenger

import (
	"errors"
	"os"

	"github.com/mcereal/botty/config"
)

func (c ChannelType) checkChannelType() (string, error) {
	url := c.Channel
	env := os.Getenv("ENVIRONMENT")
	if env == "development" || config.AppConfig.Application.Environment == "development" {
		var channelURL string
		if c.ChannelType == "discord" {
			channelURL = os.Getenv("DEV_DISCORD_CHANNEL_WEBHOOK_URL")
		} else {
			channelURL = os.Getenv("DEV_SLACK_CHANNEL_WEBHOOK_URL")
		}
		url = channelURL
	}
	if url == "" {
		return "", errors.New("no webhook found")
	}
	return url, nil
}

func CheckChannel(o IOpenPRs) (string, error) {
	url, err := o.checkChannelType()
	return url, err
}
