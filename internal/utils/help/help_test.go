package help

import "testing"

const (
	yamlString = `name: Hidden Lake Service (HLS)
desc: anonymizes traffic using the QB-problem
args: |
  [ -h, --help    ] - print information about service
  [ -v, --version ] - print version of service
  [ -p, --path    ] - set path to config, database files
  [ -n, --network ] - set network key for connections
  [ -t, --threads ] - set num of parallel functions to calculate PoW`
)

func TestPanicPrintln(t *testing.T) {
	t.Parallel()

	defer func() {
		if r := recover(); r == nil {
			t.Error("nothing panics")
			return
		}
	}()

	Println([]byte("="))
}

func ExamplePrintln() {
	Println([]byte(yamlString))
	// Output:
	// Hidden Lake Service (HLS)
	// anonymizes traffic using the QB-problem
	// [ -h, --help    ] - print information about service
	// [ -v, --version ] - print version of service
	// [ -p, --path    ] - set path to config, database files
	// [ -n, --network ] - set network key for connections
	// [ -t, --threads ] - set num of parallel functions to calculate PoW
}
