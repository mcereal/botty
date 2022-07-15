package github

import (
	"bytes"
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/mcereal/go-api-server/config"
	"github.com/mcereal/go-api-server/slack"
	"golang.org/x/exp/slices"
)

// CheckPayload filters payload info to determine what needs to be sent to slack
func CheckPayload(response []byte, c *gin.Context) (*bytes.Buffer, string) {
	openPRs := NewOpenPRs()
	openPRs.AddJSONData(response)
	payloadData := openPRs.Data

	// check to see if the the repo in payload is acceptable
	for _, v := range config.AppConfig.Team {
		if slices.Contains(v.Repos, payloadData.Repository.Name) {
			url := os.Getenv(v.Channel)
			env := os.Getenv("ENVIRONMENT")
			if env == "development" || config.AppConfig.Application.Environment == "development" {
				url = os.Getenv("PR_BOT_TEST_PUBLIC_SLACK_WEBHOOK_URL")
			}
			if url == "" {
				log.Println("No Webhook found")
				return nil, "No Webhook found"
			}
			// Dont send a message if the user is ignored
			if slices.Contains(v.IgnoreUsers, payloadData.Sender.Login) {
				log.Println("Ignoring Slack Post: User Ignored")
				return nil, "Ignoring Slack Post: User Ignored"
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

			// Create the slack message
			slackMessage := &slack.TextInfo{
				Type:        "NewPR",
				SlackTeam:   v.SlackGroupID,
				Action:      payloadData.Action,
				URL:         payloadData.PullRequest.HTMLURL,
				MessageBody: payloadData.PullRequest.Title,
				Repo:        payloadData.Repository.Name,
				Pull:        payloadData.Number,
				Merged:      payloadData.PullRequest.Merged,
			}

			body := slackMessage.CreateMessage()
			return body, url
		}
	}
	return nil, "Repo not in config"
}
