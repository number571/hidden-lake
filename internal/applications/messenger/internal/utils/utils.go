package utils

import (
	"encoding/json"
	"fmt"
	"html"
	"net/url"
	"strings"

	"github.com/number571/hidden-lake/internal/webui"
)

type sEmojis struct {
	Emojis []struct {
		Emoji     string `json:"emoji"`
		Shortname string `json:"shortname"`
	} `json:"emojis"`
}

var (
	gEmojiReplacer map[string]string
)

func init() {
	emojiSimple := new(sEmojis)
	if err := json.Unmarshal(webui.GEmojiSimpleJSON, emojiSimple); err != nil {
		panic(err)
	}

	emoji := new(sEmojis)
	if err := json.Unmarshal(webui.GEmojiJSON, emoji); err != nil {
		panic(err)
	}

	gEmojiReplacer = make(map[string]string, len(emojiSimple.Emojis)+len(emoji.Emojis))

	for _, emoji := range emojiSimple.Emojis {
		gEmojiReplacer[emoji.Shortname] = emoji.Emoji
	}
	for _, emoji := range emoji.Emojis {
		gEmojiReplacer[emoji.Shortname] = emoji.Emoji
	}
}

func ReplaceTextToEmoji(pS string) string {
	splited := strings.Split(pS, " ")
	for i, s := range splited {
		v, ok := gEmojiReplacer[s]
		if !ok {
			continue
		}
		splited[i] = v
	}
	return strings.Join(splited, " ")
}

func ReplaceTextToURLs(pS string) string {
	tagTemplate := "<a style='background-color:#b9cdcf;color:black;' target='_blank' href='%[1]s'>%[2]s</a>"
	splited := strings.Split(pS, " ")
	for i, s := range splited {
		if _, err := url.ParseRequestURI(s); err != nil {
			continue
		}
		u, err := url.Parse(s)
		if err != nil {
			continue
		}
		splited[i] = fmt.Sprintf(tagTemplate, u.String(), html.EscapeString(s))
	}
	return strings.Join(splited, " ")
}
