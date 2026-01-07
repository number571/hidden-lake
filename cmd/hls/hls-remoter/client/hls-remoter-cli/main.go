package main

import (
	"bufio"
	"context"
	"fmt"
	"net/http"
	"os"
	"os/exec"
	"runtime"
	"strings"

	"github.com/number571/hidden-lake/build"
	hlk_client "github.com/number571/hidden-lake/internal/kernel/pkg/client"
	"github.com/number571/hidden-lake/internal/services/filesharer/pkg/settings"
	hls_remoter_client "github.com/number571/hidden-lake/internal/services/remoter/pkg/client"
	"github.com/number571/hidden-lake/internal/utils/flag"
	"github.com/number571/hidden-lake/internal/utils/help"
)

var (
	gFlags = flag.NewFlagsBuilder(
		flag.NewFlagBuilder("-v", "--version").
			WithDescription("print version of application"),
		flag.NewFlagBuilder("-h", "--help").
			WithDescription("print information about application"),
		flag.NewFlagBuilder("-k", "--kernel").
			WithDescription("set internal address of the HLK").
			WithDefinedValue("localhost:9572"),
		flag.NewFlagBuilder("-f", "--friend").
			WithDescription("set alias name of the friend").
			WithDefinedValue(""),
	).Build()
)

func main() {
	args := os.Args[1:]
	if ok := gFlags.Validate(args); !ok {
		fmt.Println("args invalid")
		os.Exit(1)
	}

	if gFlags.Get("-v").GetBoolValue(args) {
		fmt.Println(build.GetVersion())
		return
	}

	if gFlags.Get("-h").GetBoolValue(args) {
		help.Println(settings.CAppFullName, settings.CAppDescription, gFlags)
		return
	}

	ctx := context.Background()
	runFunction(ctx, args)
}

func runFunction(pCtx context.Context, pArgs []string) {
	reader := bufio.NewReader(os.Stdin)

	password := inputString(reader, "password: ")
	clearTerminal()

	hlkClient := hlk_client.NewClient(
		hlk_client.NewBuilder(),
		hlk_client.NewRequester(
			gFlags.Get("-k").GetStringValue(pArgs),
			&http.Client{Timeout: build.GetSettings().GetHttpCallbackTimeout()},
		),
	)
	hlrClient := hls_remoter_client.NewClient(
		hls_remoter_client.NewBuilder(password),
		hls_remoter_client.NewRequester(hlkClient),
	)

	friend := gFlags.Get("-f").GetStringValue(pArgs)
	for {
		result, err := hlrClient.Exec(
			pCtx,
			friend,
			inputString(reader, "> "),
		)
		if err != nil {
			fmt.Printf("error: %s\n", err.Error())
			continue
		}
		fmt.Println(string(result))
	}
}

func inputString(reader *bufio.Reader, prefix string) string {
	fmt.Print(prefix)
	rawInput, err := reader.ReadString('\n')
	if err != nil {
		panic(err)
	}
	return strings.TrimSpace(rawInput)
}

func clearTerminal() {
	var cmd *exec.Cmd
	if runtime.GOOS == "windows" {
		cmd = exec.Command("cmd", "/c", "cls")
	} else {
		cmd = exec.Command("clear")
	}
	cmd.Stdout = os.Stdout
	_ = cmd.Run()
}
