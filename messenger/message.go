package messenger

import (
	"bytes"
	"encoding/json"
	"fmt"
)

// Message holds the text to be sent posted with the webhook
type Message struct {
	Content string `json:"content"`
}

// IMessage is interface that holds the methods for creating messages
type IMessage interface {
	CreateMessage() *bytes.Buffer
}

// TextInfo is the different contents of the message
type TextInfo struct {
	Type        string
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

// NewMessage is a constructor to create a new instance of TextInfo
func NewMessage() *Message {
	return &Message{
		Content: "",
	}
}

// CreateMessage builds the message string
func (s *TextInfo) CreateMessage() *bytes.Buffer {
	if s.Type == "NewPR" {
		s.Emoji = ":white_check_mark:"

		if s.Action == "closed" && s.Merged {
			s.Emoji = ":negative_squared_check_mark:"
		} else if s.Action == "closed" {
			s.Emoji = ":negative_squared_check_mark:"
		}
	}

	if s.Type == "Stale" {
		s.Emoji = ":warning:"
	}

	// https://cdn.discordapp.com/avatars/997248910991048874/df91181b3f1cf0ef1592fbe18e0962d7.webp?size=160
	// build the text string from the github url and description
	content := fmt.Sprintf("%s [%s:%v](%s) %s", s.Emoji, s.Repo, s.Pull, s.URL, s.MessageBody)

	// create the  text based off of the SlackText struct
	messageText := &Message{
		Content: content,
	}

	// json encode the text and create a buffer that can be used by the Rest client
	data, _ := json.Marshal(messageText)
	requestBytes := bytes.NewBuffer(data)
	return requestBytes
}
