package help

import (
	"testing"

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

func TestPanicToFormatAppName(t *testing.T) {
	t.Parallel()

	defer func() {
		if r := recover(); r == nil {
			t.Fatal("nothing panics")
		}
	}()

	toFormatAppName("a=b=c")
}

func TestPanicToShortAppName(t *testing.T) {
	t.Parallel()

	defer func() {
		if r := recover(); r == nil {
			t.Fatal("nothing panics")
		}
	}()

	toShortAppName("a=b=c")
}

func TestToFormatAppName(t *testing.T) {
	formatName := toFormatAppName("hidden-lake-kernel")
	if formatName != "Hidden Lake Kernel" {
		t.Log(formatName)
		t.Fatal("invalid full name")
	}
	formatName = toFormatAppName("hidden-lake-adapters=common")
	if formatName != "Hidden Lake Adapters = Common" {
		t.Log(formatName)
		t.Fatal("invalid full name (with subname)")
	}
}

func TestToShortAppName(t *testing.T) {
	shortName := toShortAppName("hidden-lake-kernel")
	if shortName != "HLK" {
		t.Log(shortName)
		t.Fatal("invalid short name")
	}
	shortName = toShortAppName("hidden-lake-adapters=common")
	if shortName != "HLA=common" {
		t.Log(shortName)
		t.Fatal("invalid short name (with subname)")
	}
}
