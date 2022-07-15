package github

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/mcereal/go-api-server/client"
	"github.com/mcereal/go-api-server/config"
	"github.com/mcereal/go-api-server/slack"
)

// GetOpenPrs gets Github PRs
func GetOpenPrs() {
	ctx := context.Background()

	gitHubBaseURL := os.Getenv("GITHUB_URL")
	if gitHubBaseURL == "" {
		gitHubBaseURL = config.AppConfig.Github.GitHubURL
	}
	token := os.Getenv("GITHUB_TOKEN")
	if token == "" {
		token = config.AppConfig.Github.GitHubToken
	}

	newHeaders := client.NewHeader()
	newHeaders.AddDefaultHeaders()
	newHeaders.AddNewHeader("Authorization", token)

	body := bytes.NewBuffer([]byte("PRs Please"))

	for _, v := range config.AppConfig.Team {
		slackGroupID := v.SlackGroupID
		repoList := v.Repos
		slackURL := os.Getenv(v.Channel)
		env := os.Getenv("ENVIRONMENT")
		if env == "development" || config.AppConfig.Application.Environment == "development" {
			slackURL = os.Getenv("PR_BOT_TEST_PUBLIC_SLACK_WEBHOOK_URL")
		}
		if slackURL == "" {
			log.Println("No Webhook found")
			return
		}
		if v.EnableCron {
			for r := range repoList {
				url := fmt.Sprintf("%s/repos/CIO-SETS/%s/pulls", gitHubBaseURL, repoList[r])
				fmt.Println("URLLLLLLLLL", url)
				newClient := &client.RestClient{
					Ctx:               ctx,
					BaseURL:           url,
					Verb:              "GET",
					Body:              body,
					AdditionalHeaders: newHeaders.Header,
				}

				responseBytes, responseHeader, err := newClient.MakeRestCall()
				if err != nil {
					log.Println("Failed to make request")
				}
				_ = responseHeader
				// fmt.Println(responseBytes)
				listPulls := ListPulls{}
				error := json.Unmarshal(responseBytes, &listPulls)
				if error != nil {
					log.Println("Failed to Unmarshal JSON")
				}
				for v := range listPulls {
					htmlURL := listPulls[v].HTMLURL
					repoName := listPulls[v].Head.Repo.Name
					pullNumber := listPulls[v].Number
					elapsedtime, elapsedMessage := ElapsedTime(listPulls[v].CreatedAt)
					if elapsedtime && !listPulls[v].Draft {
						fmt.Println("DRAFT", listPulls[v].Draft)
						fmt.Println("TIME", listPulls[v].CreatedAt)
						slackMessage := &slack.TextInfo{
							Type:        "Stale",
							SlackTeam:   slackGroupID,
							URL:         htmlURL,
							MessageBody: elapsedMessage,
							Repo:        repoName,
							Pull:        pullNumber,
						}
						body := slackMessage.CreateMessage()

						addHeaders := client.NewHeader()
						addHeaders.AddDefaultHeaders()

						// Create a REST client and then make the request using the slack message body
						restClient := &client.RestClient{
							Ctx:               ctx,
							BaseURL:           slackURL,
							Verb:              "POST",
							Body:              body,
							AdditionalHeaders: addHeaders.Header,
						}
						responseBytes, responseHeader, err := restClient.MakeRestCall()
						if err != nil {
							log.Println("Failed to make request")
						}
						_ = responseBytes
						_ = responseHeader
					}
				}
			}
		}
	}
}
