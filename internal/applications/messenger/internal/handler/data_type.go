package handler

import (
	"bytes"
	"encoding/base64"
	"html"
	"html/template"
	"strings"

	"github.com/number571/hidden-lake/internal/applications/messenger/internal/utils"
	hlm_settings "github.com/number571/hidden-lake/internal/applications/messenger/pkg/settings"
	"github.com/number571/hidden-lake/internal/utils/chars"
)

func isText(pBytes []byte) bool {
	if len(pBytes) == 0 {
		return false
	}
	return pBytes[0] == hlm_settings.CIsText
}

func isFile(pBytes []byte) bool {
	if len(pBytes) == 0 {
		return false
	}
	return pBytes[0] == hlm_settings.CIsFile
}

func wrapText(pMsg string) []byte {
	return bytes.Join([][]byte{
		{hlm_settings.CIsText},
		[]byte(pMsg),
	}, []byte{})
}

func wrapFile(filename string, pBytes []byte) []byte {
	return bytes.Join([][]byte{
		{hlm_settings.CIsFile},
		[]byte(filename),
		{hlm_settings.CIsFile},
		pBytes,
	}, []byte{})
}

func unwrapText(pBytes []byte) template.HTML {
	if !isText(pBytes) {
		return ""
	}
	text := string(pBytes[1:])
	if chars.HasNotGraphicCharacters(text) {
		return ""
	}
	text = utils.ReplaceTextToEmoji(strings.TrimSpace(text))
	text = utils.ReplaceTextToURLs(html.EscapeString(text))
	return template.HTML(text) // nolint: gosec
}

func unwrapFile(pBytes []byte) (template.HTML, string) {
	if !isFile(pBytes) {
		return "", ""
	}
	splited := bytes.Split(pBytes[1:], []byte{hlm_settings.CIsFile})
	if len(splited) < 2 {
		return "", ""
	}
	filename := string(splited[0])
	if chars.HasNotGraphicCharacters(filename) {
		return "", ""
	}
	fileBytes := bytes.Join(splited[1:], []byte{hlm_settings.CIsFile})
	if len(fileBytes) == 0 {
		return "", ""
	}
	escapedFilename := utils.FilenameEscape(strings.TrimSpace(filename))
	base64FileBytes := base64.StdEncoding.EncodeToString(fileBytes)
	return template.HTML(escapedFilename), base64FileBytes // nolint: gosec
}
