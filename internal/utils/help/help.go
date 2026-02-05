package help

import (
	"bytes"
	"fmt"
	"strings"
	"unicode"

	"github.com/number571/hidden-lake/internal/utils/flag"
)

func Println(pAppName, pDescription string, pArgs flag.IFlags) {
	args := strings.Builder{}
	args.Grow(1 << 10)

	for _, arg := range pArgs.List() {
		aliases := arg.GetAliases()
		args.WriteString(fmt.Sprintf(
			"[ %s ] = %s\n",
			strings.Join(aliases, ", "),
			arg.GetDescription(),
		))
	}

	fmt.Printf(
		"<%s (%s)>\nDescription: %s\nArguments:\n%s\n",
		toFormatAppName(pAppName),
		toShortAppName(pAppName),
		pDescription,
		strings.TrimSpace(args.String()),
	)
}

func toFormatAppName(pFullName string) string {
	// Example: hidden-lake-adapters=common -> hidden-lake-adapters = Common
	splitedEq := bytes.Split([]byte(pFullName), []byte("="))
	if len(splitedEq) > 2 {
		panic("length of splited by '=' > 2")
	}
	if len(splitedEq) == 2 {
		splitedEq[1][0] = byte(unicode.ToUpper(rune(splitedEq[1][0])))
	}
	joinedEq := bytes.Join(splitedEq, []byte(" = "))

	// Example: hidden-lake-adapters=common -> Hidden Lake Adapters = Common
	splitedSb := bytes.Split(joinedEq, []byte("-"))
	for i := range splitedSb {
		splitedSb[i][0] = byte(unicode.ToUpper(rune(splitedSb[i][0])))
	}
	joinedSb := bytes.Join(splitedSb, []byte(" "))

	return string(joinedSb)
}

func toShortAppName(pFullName string) string {
	result := strings.Builder{}
	result.Grow(len(pFullName))

	splitedEq := strings.Split(pFullName, "=")
	if len(splitedEq) > 2 {
		panic("length of splited by '=' > 2")
	}

	splitedSb := strings.Split(splitedEq[0], "-")
	for _, v := range splitedSb {
		firstCh := unicode.ToUpper(rune(v[0]))
		_, _ = result.WriteRune(firstCh)
	}

	if len(splitedEq) == 2 {
		_ = result.WriteByte('=')
		_, _ = result.WriteString(splitedEq[1])
	}

	return result.String()
}
