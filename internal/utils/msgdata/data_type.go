package msgdata

import (
	"bytes"
	"encoding/base64"
	"html"
	"html/template"
	"strings"

	"github.com/number571/hidden-lake/internal/utils/chars"
)

const (
	cIsText = 0x01
	cIsFile = 0x02
)

func isText(pBytes []byte) bool {
	if len(pBytes) == 0 {
		return false
	}
	return pBytes[0] == cIsText
}

func isFile(pBytes []byte) bool {
	if len(pBytes) == 0 {
		return false
	}
	return pBytes[0] == cIsFile
}

func wrapText(pMsg string) []byte {
	return bytes.Join([][]byte{
		{cIsText},
		[]byte(pMsg),
	}, []byte{})
}

func wrapFile(filename string, pBytes []byte) []byte {
	return bytes.Join([][]byte{
		{cIsFile},
		[]byte(filename),
		{cIsFile},
		pBytes,
	}, []byte{})
}

func unwrapText(pBytes []byte) template.HTML {
	if !isText(pBytes) {
		return ""
	}
	text := string(pBytes[1:]) //nolint:gosec
	if chars.HasNotGraphicCharacters(text) {
		return ""
	}
	text = ReplaceTextToEmoji(strings.TrimSpace(text))
	text = ReplaceTextToURLs(html.EscapeString(text))
	return template.HTML(text) // nolint: gosec
}

func unwrapFile(pBytes []byte) (template.HTML, string) {
	if !isFile(pBytes) {
		return "", ""
	}
	splited := bytes.Split(pBytes[1:], []byte{cIsFile}) //nolint:gosec
	if len(splited) < 2 {
		return "", ""
	}
	filename := string(splited[0])
	if chars.HasNotGraphicCharacters(filename) {
		return "", ""
	}
	fileBytes := bytes.Join(splited[1:], []byte{cIsFile})
	if len(fileBytes) == 0 {
		return "", ""
	}
	escapedFilename := FilenameEscape(strings.TrimSpace(filename))
	base64FileBytes := base64.StdEncoding.EncodeToString(fileBytes)
	return template.HTML(escapedFilename), base64FileBytes // nolint: gosec
}
