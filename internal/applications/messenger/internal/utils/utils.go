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
	gEmojiReplacer map[string]string
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
