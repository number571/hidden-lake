package main

import (
	"bufio"
	"context"
	"fmt"
	"net/http"
	"os"
	"strings"

	"github.com/number571/hidden-lake/build"
	"github.com/number571/hidden-lake/internal/services/filesharer/pkg/settings"
	hls_client "github.com/number571/hidden-lake/internal/services/messenger/pkg/client"
	"github.com/number571/hidden-lake/internal/utils/flag"
	"github.com/number571/hidden-lake/internal/utils/help"
)

var (
	gFlags = flag.NewFlagsBuilder(
		flag.NewFlagBuilder("-v", "--version").
			WithDescription("print version of application"),
		flag.NewFlagBuilder("-h", "--help").
			WithDescription("print information about application"),
		flag.NewFlagBuilder("-s", "--service").
			WithDescription("set internal address of the HLS").
			WithDefinedValue("localhost:9591"),
		flag.NewFlagBuilder("-f", "--friend").
			WithDescription("set alias name of the friend").
			WithDefinedValue(""),
	).Build()
)

func main() {
	args := os.Args[1:]
	if ok := gFlags.Validate(args); !ok {
		fmt.Println(args)
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
	hlsClient := hls_client.NewClient(
		hls_client.NewRequester(
			gFlags.Get("-s").GetStringValue(pArgs),
			&http.Client{Timeout: build.GetSettings().GetHttpCallbackTimeout()},
		),
	)

	friend := gFlags.Get("-f").GetStringValue(pArgs)
	limit, err := hlsClient.GetMessageLimit(pCtx)
	if err != nil {
		fmt.Printf("error: %s\n", err.Error())
		return
	}

	fmt.Printf("{\n\t\"friend_name\": \"%s\",\n\t\"payload_limit\": %d\n}\n\n", friend, limit)

	msgs, err := hlsClient.LoadMessages(pCtx, friend, 256, 256, true)
	if err != nil {
		fmt.Printf("error: %s\n", err.Error())
		return
	}

	iam := "<iam>"
	for _, m := range msgs {
		sender := iam
		if m.IsIncoming() {
			sender = friend
		}
		fmt.Printf("%s: %s [%s]\n", sender, m.GetMessage(), m.GetTimestamp())
	}

	go func() {
		for {
			m, err := hlsClient.ListenChat(pCtx, friend, "hls-messenger-cli")
			if err != nil {
				fmt.Printf("error: %s\n", err.Error())
				continue
			}
			fmt.Printf("%s: %s [%s]\n", friend, m.GetMessage(), m.GetTimestamp())
		}
	}()

	reader := bufio.NewReader(os.Stdin)
	for {
		_, err := hlsClient.PushMessage(pCtx, friend, inputString(reader, ""))
		if err != nil {
			fmt.Printf("error: %s\n", err.Error())
			continue
		}
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
