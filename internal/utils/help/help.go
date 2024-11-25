package help

import (
	"fmt"
	"strings"

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
	fmt.Printf(
		"%s\n%s\n%s\n",
		pFullName,
		pDescription,
		strings.TrimSpace(args.String()),
	)
}
