package webui

import (
	"embed"
	"io/fs"
)

var (
	//go:embed static/json/emoji.json
	GEmojiJSON []byte

	//go:embed static/json/emoji_simple.json
	GEmojiSimpleJSON []byte

	//go:embed static
	gEmbededStatic embed.FS

	//go:embed template
	gEmbededTemplate embed.FS
)

func GetStaticPath() fs.FS {
	fsys, err := fs.Sub(gEmbededStatic, "static")
	if err != nil {
		panic(err)
	}
	return fsys
}

func GetTemplatePath() fs.FS {
	fsys, err := fs.Sub(gEmbededTemplate, "template")
	if err != nil {
		panic(err)
	}
	return fsys
}
