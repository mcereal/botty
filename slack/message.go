package slack

import (
	"bytes"
	"encoding/json"
	"fmt"
)

// Message holds the text to be sent posted with the webhook
type Message struct {
	Text string `json:"text"`
}

// ISlackMessage is interface that holds the methods for creating messages
type ISlackMessage interface {
	CreateMessage() *bytes.Buffer
}

// TextInfo is the different contents of the message
type TextInfo struct {
	Type        string
	SlackTeam   string
	Action      string
	URL         string
	Title       string
	Repo        string
	Description string
	Emoji       string
	Pull        int
	Merged      bool
	MessageBody string
}

// NewMessage is a constructor to create a new instance of SlackText
func NewMessage() *Message {
	return &Message{
		Text: "",
	}
}

// CreateMessage builds the message string
func (s *TextInfo) CreateMessage() *bytes.Buffer {
	if s.Type == "NewPR" {
		s.Emoji = ":pr:"

		if s.Action == "closed" && s.Merged {
			s.Emoji = ":pr-merged:"
		} else if s.Action == "closed" {
			s.Emoji = ":pr-closed:"
		}
	}

	if s.Type == "Stale" {
		s.Emoji = ":ibm-warning-filled:"
	}
	// build the text string from the github url and description
	text := fmt.Sprintf("%s %s <%s|%s:%v> %s", s.SlackTeam, s.Emoji, s.URL, s.Repo, s.Pull, s.MessageBody)

	// create the  text based off of the SlackText struct
	slackText := &Message{
		Text: text,
	}

	// json encode the text and create a buffer that can be used by the Rest client
	data, _ := json.Marshal(slackText)
	requestBytes := bytes.NewBuffer(data)
	return requestBytes
}
