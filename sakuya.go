package sakuya

import (
	"github.com/nlopes/slack"
	"image/color"
	"fmt"
	"io"
	"os"
)

type SlackWriter struct {
	client *slack.Client
	channel string
	text string
	userName string
	c color.RGBA
	iconURL string
	writer io.Writer
}

func defString(s string) string {
	if s == "" {
		return "sakuya"
	}
	return s
}

func New(token, channel, text, name string)*SlackWriter {
	null, err := os.Open(os.DevNull)
	if err != nil {
		return nil
	}
	return &SlackWriter{
		slack.New(token),
		channel,
		text,
		defString(name),
		color.RGBA{0x00, 0xff, 0x00, 0xff},
		"",
		null,
	}
}

func (s *SlackWriter) UseStdout() {
	s.writer = os.Stdout
}

func (s *SlackWriter) UseStderr() {
	s.writer = os.Stderr
}

func (s *SlackWriter) SetColor(c color.RGBA) {
	s.c = c
}

func (s *SlackWriter) SetIconURL(u string) {
	s.iconURL = u
}

func (s *SlackWriter) Write(p []byte) (n int, err error) {
	attachment := slack.Attachment{
		Text: string(p),
		Color: fmt.Sprintf("%02X%02X%02X", s.c.R, s.c.G, s.c.B),
	}
	params := slack.PostMessageParameters{
		Username: s.userName,
		Attachments: []slack.Attachment{attachment},
		IconURL: s.iconURL,
	}
	if _, _, err := s.client.PostMessage(s.channel, s.text, params); err != nil {
		return 0, err
	}
	if l, err := s.writer.Write(p); err != nil {
		return l, err
	}
	return len(p), nil
}

