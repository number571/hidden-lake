package appname

import (
	"bytes"
	"unicode"
)

var (
	_ IFmtAppName = &sFmtAppName{}
)

type sFmtAppName struct {
	fShortName string
	fFullName  string
}

func LoadAppName(pFullName string) IFmtAppName {
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

	return &sFmtAppName{
		fShortName: string(shortName),
		fFullName:  string(nameJoined),
	}
}

func (p *sFmtAppName) Short() string {
	return p.fShortName
}

func (p *sFmtAppName) Full() string {
	return p.fFullName
}
