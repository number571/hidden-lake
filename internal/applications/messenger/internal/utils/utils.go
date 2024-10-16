package utils

import (
	"encoding/json"
	"strings"

	"github.com/number571/hidden-lake/internal/applications/messenger/web"
)

type sEmojis struct {
	Emojis []struct {
		Emoji     string `json:"emoji"`
		Shortname string `json:"shortname"`
	} `json:"emojis"`
}

var (
	gEmojiReplacer *strings.Replacer
)

func init() {
	emojiSimple := new(sEmojis)
	if err := json.Unmarshal(web.GEmojiSimpleJSON, emojiSimple); err != nil {
		panic(err)
	}

	emoji := new(sEmojis)
	if err := json.Unmarshal(web.GEmojiJSON, emoji); err != nil {
		panic(err)
	}

	replacerList := make([]string, 0, len(emojiSimple.Emojis)+len(emoji.Emojis))

	for _, emoji := range emojiSimple.Emojis {
		replacerList = append(replacerList, emoji.Shortname, emoji.Emoji)
	}
	for _, emoji := range emoji.Emojis {
		replacerList = append(replacerList, emoji.Shortname, emoji.Emoji)
	}

	gEmojiReplacer = strings.NewReplacer(replacerList...)
}

func ReplaceTextToEmoji(pS string) string {
	return gEmojiReplacer.Replace(pS)
}
