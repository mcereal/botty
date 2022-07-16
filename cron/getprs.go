package cron

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/mcereal/botty/client"
	"github.com/mcereal/botty/config"
	"github.com/mcereal/botty/messenger"
)

// GetOpenPrs gets Github PRs
func GetOpenPrs() {
	ctx := context.Background()

	gitHubBaseURL := os.Getenv("GITHUB_URL")
	if gitHubBaseURL == "" {
		gitHubBaseURL = config.AppConfig.Github.GitHubURL
	}
	token := os.Getenv("GITHUB_SECRET_TOKEN")
	if token == "" {
		token = config.AppConfig.Github.GitHubToken
	}

	newHeaders := client.NewHeader()
	newHeaders.AddDefaultHeaders()
	newHeaders.AddNewHeader("Authorization", token)

	body := bytes.NewBuffer([]byte("PRs Please"))

	for _, v := range config.AppConfig.Team {
		repoList := v.Repos
		org := v.Owner
		channelURL := os.Getenv(v.Channel)
		channelType := v.ChannelType

		checkChannelType := messenger.ChannelType{Channel: channelURL, ChannelType: channelType}
		channelURL, err := messenger.CheckChannel(checkChannelType)
		if err != nil {
			log.Println("No Webhook found")
			// return nil, "No Webhook found"
		}
		elapsedDuration := v.CronElapsedDuration

		if v.EnableCron {
			for r := range repoList {
				url := fmt.Sprintf("%s/repos/%s/%s/pulls", gitHubBaseURL, org, repoList[r])
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

				listPulls := messenger.ListPulls{}
				error := json.Unmarshal(responseBytes, &listPulls)
				if error != nil {
					log.Println(error)
					log.Println("Failed to Unmarshal JSON")
				}
				for v := range listPulls {
					htmlURL := listPulls[v].HTMLURL
					repoName := listPulls[v].Head.Repo.Name
					pullNumber := listPulls[v].Number
					elapsedtime, elapsedMessage := ElapsedTime(listPulls[v].CreatedAt, elapsedDuration)
					if elapsedtime && !listPulls[v].Draft {
						// fmt.Println("DRAFT", listPulls[v].Draft)
						// fmt.Println("TIME", listPulls[v].CreatedAt)
						messageContent := &messenger.TextInfo{
							Type:        "Stale",
							ChannelType: channelType,
							URL:         htmlURL,
							MessageBody: elapsedMessage,
							Repo:        repoName,
							Pull:        pullNumber,
						}
						body := messageContent.CreateMessage()

						addHeaders := client.NewHeader()
						addHeaders.AddDefaultHeaders()

						// Create a REST client and then make the request using the message body
						restClient := &client.RestClient{
							Ctx:               ctx,
							BaseURL:           channelURL,
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
