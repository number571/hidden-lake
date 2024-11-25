package name

import (
	"bytes"
	"unicode"
)

var (
	_ IServiceName = &sServiceName{}
)

type sServiceName struct {
	fShortName  string
	fFormatName string
}

func LoadServiceName(pFullName string) IServiceName {
	// Example: hidden-lake-adapters=common -> [Hidden Lake Adapters=common], [HLA]
	nameSplited := bytes.Split([]byte(pFullName), []byte("-"))
	shortName := make([]byte, 0, len(nameSplited))
	for i := range nameSplited {
		nameSplited[i][0] = byte(unicode.ToUpper(rune(nameSplited[i][0])))
		shortName = append(shortName, nameSplited[i][0])
	}
	nameJoined := bytes.Join(nameSplited, []byte(" "))

	// Example: Hidden Lake Adapters=common -> [Hidden Lake Adapters = Common], [HLA=common]
	nameSplited = bytes.Split(nameJoined, []byte("="))
	if len(nameSplited) > 1 {
		shortName = append(shortName, '=')
		shortName = append(shortName, bytes.Join(nameSplited[1:], []byte("="))...)
		nameSplited[1][0] = byte(unicode.ToUpper(rune(nameSplited[1][0])))
	}
	nameJoined = bytes.Join(nameSplited, []byte(" = "))

	return &sServiceName{
		fShortName:  string(shortName),
		fFormatName: string(nameJoined),
	}
}

func (p *sServiceName) Short() string {
	return p.fShortName
}

func (p *sServiceName) Format() string {
	return p.fFormatName
}
