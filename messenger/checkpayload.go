package messenger

import (
	"bytes"
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/mcereal/botty/config"
	"golang.org/x/exp/slices"
)

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

			checkChannelType := ChannelType{Channel: channelURL, ChannelType: channelType}
			url, err := CheckChannel(checkChannelType)
			if err != nil {
				log.Println("No Webhook found")
				return nil, "No Webhook found"
			}

			// Dont send a message if the user is ignored
			if slices.Contains(v.IgnoreUsers, payloadData.Sender.Login) {
				log.Println("Ignoring Message Post: User Ignored")
				return nil, "Ignoring Message Post: User Ignored"
			}

			// check the payload data for pull request events
			switch {
			case payloadData.Action == "review_requested":
				// Don't send messages on reviews
				log.Println("This is a dup event")
				return nil, "This is a dup event"
			case payloadData.Action == "synchronize":
				// Don't send messages on reopened
				log.Println("This is a dup event")
				return nil, "This is a dup event"
			case payloadData.Action == "closed" && payloadData.PullRequest.Merged:
				// Don't send messages when closed and merged
				log.Println("Merged PR - not reporting")
				return nil, "Merged PR - not reporting"
			case payloadData.Action == "closed":
				// Don't send messages when closed
				log.Println("Closed PR - not reporting")
				return nil, "Closed PR - not reporting"
			case payloadData.PullRequest.Draft:
				// Don't send message if its a draft
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
			log.Println(body)
			return body, url
		}
	}
	return nil, "Repo not in config"
}
