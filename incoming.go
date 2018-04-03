package sakuya

import (
	"bytes"
	"encoding/json"
	"fmt"
	"image/color"
	"net/http"

	"github.com/nlopes/slack"
)

type IncomingWriter struct {
	apiURL                string
	baseColor             color.RGBA
	postMessageParameters slack.PostMessageParameters
}

func NewIncomingWriter(apiURL string, userName string) *IncomingWriter {
	return &IncomingWriter{
		postMessageParameters: slack.PostMessageParameters{
			Username: defString(userName),
		},
		apiURL: apiURL,
		baseColor: color.RGBA{
			0x00,
			0xff,
			0x00,
			0xff,
		},
	}
}

func (i *IncomingWriter) reset() {
	i.postMessageParameters = slack.PostMessageParameters{}
}

func (i *IncomingWriter) ChangeChannel(channel string) {
	i.postMessageParameters.Channel = channel
}

func (i *IncomingWriter) SetBaseColor(c color.RGBA) {
	i.baseColor = c
}

func (i *IncomingWriter) AddInfo(txt string) {
	a := slack.Attachment{
		Text:  txt,
		Color: fmt.Sprintf("%02X%02X%02X", 0x00, 0xff, 0x00),
	}

	i.postMessageParameters.Attachments = append(i.postMessageParameters.Attachments, a)
}

func (i *IncomingWriter) AddWarn(txt string) {
	a := slack.Attachment{
		Text:  txt,
		Color: fmt.Sprintf("%02X%02X%02X", 0xff, 0xff, 0x00),
	}

	i.postMessageParameters.Attachments = append(i.postMessageParameters.Attachments, a)
}

func (i *IncomingWriter) AddError(txt string) {
	a := slack.Attachment{
		Text:  txt,
		Color: fmt.Sprintf("%02X%02X%02X", 0xff, 0x00, 0x00),
	}

	i.postMessageParameters.Attachments = append(i.postMessageParameters.Attachments, a)
}

func (i *IncomingWriter) AddUnknown(txt string) {
	a := slack.Attachment{
		Text:  txt,
		Color: fmt.Sprintf("%02X%02X%02X", 0xaa, 0xaa, 0xaa),
	}

	i.postMessageParameters.Attachments = append(i.postMessageParameters.Attachments, a)

}

func (i *IncomingWriter) AddText(txt string) {
	a := slack.Attachment{
		Text:  txt,
		Color: fmt.Sprintf("%02X%02X%02X", i.baseColor.R, i.baseColor.G, i.baseColor.B),
	}
	i.postMessageParameters.Attachments = append(i.postMessageParameters.Attachments, a)
}

func (i *IncomingWriter) AddTextAndColor(txt string, c color.RGBA) {
	a := slack.Attachment{
		Text:  txt,
		Color: fmt.Sprintf("%02X%02X%02X", c.R, c.G, c.B),
	}
	i.postMessageParameters.Attachments = append(i.postMessageParameters.Attachments, a)

}

func (i *IncomingWriter) AddAttachment(a slack.Attachment) {
	i.postMessageParameters.Attachments = append(i.postMessageParameters.Attachments, a)
}

func (i *IncomingWriter) Flush() error {
	defer i.reset()
	buffer := new(bytes.Buffer)
	if err := json.NewEncoder(buffer).Encode(i.postMessageParameters); err != nil {
		return err
	}
	req, err := http.NewRequest("POST", i.apiURL, buffer)
	if err != nil {
		return err
	}
	if _, err := http.DefaultClient.Do(req); err != nil {
		return err
	}
	return nil
}

func (i *IncomingWriter) Write(p []byte) (n int, err error) {
	attachment := slack.Attachment{
		Text:  string(p),
		Color: fmt.Sprintf("%02X%02X%02X", i.baseColor.R, i.baseColor.G, i.baseColor.B),
	}
	params := slack.PostMessageParameters{
		Username:    i.postMessageParameters.Username,
		Attachments: []slack.Attachment{attachment},
	}

	buffer := new(bytes.Buffer)
	if err := json.NewEncoder(buffer).Encode(params); err != nil {
		return 0, err
	}
	req, err := http.NewRequest("POST", i.apiURL, buffer)
	if err != nil {
		return 0, err
	}
	if _, err := http.DefaultClient.Do(req); err != nil {
		return 0, err
	}
	return len(p), nil
}
