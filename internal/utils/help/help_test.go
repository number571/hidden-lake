package help

import (
	"github.com/number571/hidden-lake/internal/utils/flag"
	"github.com/number571/hidden-lake/internal/utils/name"
)

func ExamplePrintln() {
	Println(
		name.LoadServiceName("hidden-lake-adapter=common"),
		"adapts HL traffic to a custom HTTP server",
		flag.NewFlagsBuilder(
			flag.NewFlagBuilder("-v", "--version").
				WithDescription("print information about service"),
			flag.NewFlagBuilder("-h", "--help").
				WithDescription("print version of service"),
			flag.NewFlagBuilder("-p", "--path").
				WithDescription("set path to config, database files").
				WithDefinedValue("."),
			flag.NewFlagBuilder("-n", "--network").
				WithDescription("set network key for connections").
				WithDefinedValue(""),
			flag.NewFlagBuilder("-t", "--threads").
				WithDescription("set num of parallel functions to calculate PoW").
				WithDefinedValue("1"),
		).Build(),
	)
	// Output:
	// <Hidden Lake Adapter = Common (HLA=common)>
	// Description: adapts HL traffic to a custom HTTP server
	// Arguments:
	// [ -v, --version ] = print information about service
	// [ -h, --help ] = print version of service
	// [ -p, --path ] = set path to config, database files
	// [ -n, --network ] = set network key for connections
	// [ -t, --threads ] = set num of parallel functions to calculate PoW
}
