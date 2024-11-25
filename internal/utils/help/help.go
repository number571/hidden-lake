package help

import (
	"bytes"
	"fmt"
	"strings"
	"unicode"

	"github.com/number571/hidden-lake/internal/utils/flag"
)

func Println(pFullName, pDescription string, pArgs flag.IFlags) {
	args := strings.Builder{}
	args.Grow(1 << 10)

	for _, arg := range pArgs.List() {
		aliases := arg.GetAliases()
		slice := make([]string, 0, len(aliases))
		for _, a := range aliases {
			slice = append(slice, "-"+a)
		}
		args.WriteString(fmt.Sprintf(
			"[ %s ] = %s\n",
			strings.Join(slice, ", "),
			arg.GetDescription(),
		))
	}

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

	fmt.Printf(
		"<%s (%s)>\nDescription: %s\nArguments:\n%s\n",
		string(nameJoined),
		string(shortName),
		pDescription,
		strings.TrimSpace(args.String()),
	)
}
