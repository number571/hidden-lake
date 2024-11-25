package help

import (
	"fmt"
	"strings"

	"github.com/number571/go-peer/pkg/encoding"
)

type sHelp struct {
	Name string `yaml:"name"`
	Desc string `yaml:"desc"`
	Args string `yaml:"args"`
}

func Println(yamlBytes []byte) {
	help := &sHelp{}
	if err := encoding.DeserializeYAML(yamlBytes, help); err != nil {
		panic(err)
	}
	fmt.Printf(
		"%s\n%s\n%s\n",
		strings.TrimSpace(help.Name),
		strings.TrimSpace(help.Desc),
		strings.TrimSpace(help.Args),
	)
}
