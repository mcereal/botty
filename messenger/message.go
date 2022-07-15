package messenger

import (
	"bytes"
	"encoding/json"
	"fmt"
)

// Message holds the text to be sent posted with the webhook
type Message struct {
	Username     string   `json:"username"`
	BotAvatarURL string   `json:"avatar_url"`
	Content      string   `json:"content"`
	Embeds       []Embeds `json:"embeds"`
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
	AvatarURL   string
	Login       string
	AuthorURL   string
}

// Embeds is the message embeded content
type Embeds struct {
	Author Author `json:"author"`
}

// Author is the author information
type Author struct {
	Name      string `json:"name"`
	AuthorURL string `json:"url"`
	IconURL   string `json:"icon_url"`
}

// NewMessage is a constructor to create a new instance of TextInfo
func NewMessage() *Message {
	return &Message{}
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

	// build the text string from the github url and description
	content := fmt.Sprintf("%s [%s:%v](%s) %s", s.Emoji, s.Repo, s.Pull, s.URL, s.MessageBody)

	embeds := []Embeds{{Author{s.Login, s.AuthorURL, s.AvatarURL}}}

	// create the  text based off of the SlackText struct
	messageText := &Message{
		Username:     "GitHub",
		BotAvatarURL: "https://cdn.discordapp.com/avatars/997248910991048874/df91181b3f1cf0ef1592fbe18e0962d7.webp?size=160",
		Content:      content,
		Embeds:       embeds,
	}

	// json encode the text and create a buffer that can be used by the Rest client
	data, _ := json.Marshal(messageText)
	requestBytes := bytes.NewBuffer(data)
	return requestBytes
}
