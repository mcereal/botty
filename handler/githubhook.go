package handler

import (
	"io/ioutil"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/mcereal/go-api-server/client"
	"github.com/mcereal/go-api-server/github"
	log "github.com/sirupsen/logrus"
)

// GitHubWebhookHandler is the handler for Github webhooks
func (h *Handler) GitHubWebhookHandler(c *gin.Context) {
	//Load Environment variables
	err := godotenv.Load()
	if err != nil {
		log.Println("Error loading .env file")
	}

	// set headers and send a json response that the webhook was recieved
	c.Header("Access-Control-Allow-Origin", "*")
	c.Header("Content-Type", "application/json")
	c.Header("User-Agent", "cs-code-review-bot")
	c.JSON(http.StatusOK, gin.H{
		"Message": "Webhook recieved",
	})

	// read the request body
	jsonDataBytes, err := ioutil.ReadAll(c.Request.Body)
	if err != nil {
		log.Error(err)
	}

	// get the message and url
	body, url := github.CheckPayload(jsonDataBytes, c)
	if body == nil {
		return
	}
	addHeaders := client.NewHeader()
	addHeaders.AddDefaultHeaders()

	// Create a REST client and then make the request using the slack message body
	restClient := &client.RestClient{
		Ctx:               c,
		BaseURL:           url,
		Verb:              "POST",
		Body:              body,
		AdditionalHeaders: addHeaders.Header,
	}
	// make the rest call
	responseBytes, responseHeader, err := restClient.MakeRestCall()
	if err != nil {
		log.Println("Failed to make request")
	}
	_ = responseBytes
	_ = responseHeader
}
