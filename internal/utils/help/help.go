package help

import (
	"fmt"
	"strings"

	"github.com/number571/hidden-lake/internal/utils/appname"
	"github.com/number571/hidden-lake/pkg/utils/flag"
)

func Println(pAppName appname.IAppName, pDescription string, pArgs flag.IFlags) {
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
		pAppName.Format(),
		pAppName.Short(),
		pDescription,
		strings.TrimSpace(args.String()),
	)
}
