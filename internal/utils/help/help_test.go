package help

import (
	"github.com/number571/hidden-lake/internal/utils/flag"
)

func ExamplePrintln() {
	Println(
		"Hidden Lake Service (HLS)",
		"anonymizes traffic using the QB-problem",
		flag.NewFlagsBuilder(
			flag.NewFlagBuilder("v", "version").
				WithDescription("print information about service"),
			flag.NewFlagBuilder("h", "help").
				WithDescription("print version of service"),
			flag.NewFlagBuilder("p", "path").
				WithDescription("set path to config, database files").
				WithDefaultValue("."),
			flag.NewFlagBuilder("n", "network").
				WithDescription("set network key for connections").
				WithDefaultValue(""),
			flag.NewFlagBuilder("t", "threads").
				WithDescription("set num of parallel functions to calculate PoW").
				WithDefaultValue("1"),
		).Build(),
	)
	// Output:
	// Hidden Lake Service (HLS)
	// anonymizes traffic using the QB-problem
	// [-v, --version] - print information about service
	// [-h, --help] - print version of service
	// [-p, --path] - set path to config, database files
	// [-n, --network] - set network key for connections
	// [-t, --threads] - set num of parallel functions to calculate PoW
}
