package messenger

import (
	"bytes"
	"errors"
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/mcereal/go-api-server-example/config"
	"golang.org/x/exp/slices"
)

func checkChannelType(channel, channelType string) (string, error) {
	url := channel
	env := os.Getenv("ENVIRONMENT")
	if env == "development" || config.AppConfig.Application.Environment == "development" {
		var channelURL string
		if channelType == "discord" {
			channelURL = os.Getenv("DEV_DISCORD_CHANNEL_WEBHOOK_URL")
		} else {
			channelURL = os.Getenv("DEV_SLACK_CHANNEL_WEBHOOK_URL")
		}
		url = channelURL
	}
	if url == "" {
		return "", errors.New("no webhook found")
	}
	log.Println(url)
	return url, nil
}

// CheckPayload filters payload info to determine what needs to be sent to channel
func CheckPayload(response []byte, c *gin.Context) (*bytes.Buffer, string) {
	openPRs := NewOpenPRs()
	openPRs.AddJSONData(response)
	payloadData := openPRs.Data
	// check to see if the the repo in payload is acceptable
	for _, v := range config.AppConfig.Team {
		if slices.Contains(v.Repos, payloadData.Repository.Name) {

			channelURL := os.Getenv(v.Channel)
			channelType := v.ChannelType

			url, err := checkChannelType(channelURL, channelType)
			if err != nil {
				log.Println("No Webhook found")
				return nil, "No Webhook found"
			}

			// Dont send a message if the user is ignored
			if slices.Contains(v.IgnoreUsers, payloadData.Sender.Login) {
				log.Println("Ignoring Message Post: User Ignored")
				return nil, "Ignoring Message Post: User Ignored"
			}

			// Don't send messages on reviews or reopened
			if payloadData.Action == "review_requested" || payloadData.Action == "synchronize" {
				log.Println("This is a dup event")
				return nil, "This is a dup event"
			}

			// Don't send messages when closed or merged
			if payloadData.Action == "closed" && payloadData.PullRequest.Merged {
				log.Println("Merged PR - not reporting")
				return nil, "Merged PR - not reporting"
			} else if payloadData.Action == "closed" {
				log.Println("Closed PR - not reporting")
				return nil, "Closed PR - not reporting"
			}

			// Don't send message if its a draft
			if payloadData.PullRequest.Draft {
				log.Println("PR is draft - not reporting")
				return nil, "PR is draft  - not reporting"
			}

			// Create the message
			messageContent := &TextInfo{
				Type:        "NewPR",
				ChannelType: v.ChannelType,
				Action:      payloadData.Action,
				URL:         payloadData.PullRequest.HTMLURL,
				MessageBody: payloadData.PullRequest.Title,
				Repo:        payloadData.Repository.Name,
				Pull:        payloadData.Number,
				Merged:      payloadData.PullRequest.Merged,
				AvatarURL:   payloadData.Sender.AvatarURL,
				Login:       payloadData.Sender.Login,
				AuthorURL:   payloadData.Sender.AuthorURL,
				Body:        payloadData.PullRequest.Body,
			}

			body := messageContent.CreateMessage()
			return body, url
		}
	}
	return nil, "Repo not in config"
}
