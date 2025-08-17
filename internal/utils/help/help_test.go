package help

import (
	"github.com/number571/hidden-lake/internal/utils/flag"
)

func ExamplePrintln() {
	Println(
		"hidden-lake-adapter=common",
		"adapts HL traffic to a custom HTTP server",
		flag.NewFlagsBuilder(
			flag.NewFlagBuilder("-v", "--version").
				WithDescription("print version of application"),
			flag.NewFlagBuilder("-h", "--help").
				WithDescription("print information about application"),
			flag.NewFlagBuilder("-p", "--path").
				WithDescription("set path to config, database files").
				WithDefinedValue("."),
			flag.NewFlagBuilder("-n", "--network").
				WithDescription("set network key of connections from build").
				WithDefinedValue(""),
		).Build(),
	)
	// Output:
	// <Hidden Lake Adapter = Common (HLA=common)>
	// Description: adapts HL traffic to a custom HTTP server
	// Arguments:
	// [ -v, --version ] = print version of application
	// [ -h, --help ] = print information about application
	// [ -p, --path ] = set path to config, database files
	// [ -n, --network ] = set network key of connections from build
}
