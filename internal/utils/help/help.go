package help

import (
	"fmt"
	"strings"

	"github.com/number571/hidden-lake/internal/utils/flag"
	"github.com/number571/hidden-lake/internal/utils/name"
)

func Println(pServiceName name.IServiceName, pDescription string, pArgs flag.IFlags) {
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
		"<%s (%s)>\nDescription: %s\nArguments:\n%s\n",
		pServiceName.Format(),
		pServiceName.Short(),
		pDescription,
		strings.TrimSpace(args.String()),
	)
}
