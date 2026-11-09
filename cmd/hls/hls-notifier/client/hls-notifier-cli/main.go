package main

import (
	"bufio"
	"context"
	"fmt"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/number571/hidden-lake/build"
	"github.com/number571/hidden-lake/internal/services/notifier/pkg/settings"
	"github.com/number571/hidden-lake/internal/utils/chars"
	"github.com/number571/hidden-lake/internal/utils/flag"
	"github.com/number571/hidden-lake/internal/utils/help"
	hls_client "github.com/number571/hidden-lake/pkg/api/services/notifier/client"
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
		flag.NewFlagBuilder("-w", "--wait").
			WithDescription("wait response from node"),
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
		help.Println(settings.CAppFullName+"-cli", settings.CAppDescription, gFlags)
		return
	}

	ctx := context.Background()
	if err := runFunction(ctx, args); err != nil {
		fmt.Fprintln(os.Stderr, "\n", err)
		os.Exit(1)
	}
}

func runFunction(pCtx context.Context, pArgs []string) error {
	hlsClient := hls_client.NewClient(
		gFlags.Get("-s").GetStringValue(pArgs),
		&http.Client{Timeout: build.GetSettings().GetHttpCallbackTimeout()},
		hls_client.WithWaitResponse(gFlags.Get("-w").GetBoolValue(pArgs)),
	)

	limit, err := hlsClient.GetNotifyLimit(pCtx)
	if err != nil {
		return err
	}

	friend := gFlags.Get("-f").GetStringValue(pArgs)
	fmt.Printf("{\n\t\"friend_name\": \"%s\",\n\t\"payload_limit\": %d\n}\n\n", friend, limit)

	go func() {
		for {
			sender, msg, err := hlsClient.RecvNotification(pCtx, "hls-notifier-cli")
			if err != nil {
				fmt.Printf("warn: %s\n", err.Error())
				continue
			}
			if friend != sender {
				continue
			}
			strMsg := string(msg)
			if chars.HasNotGraphicCharacters(strMsg) {
				strMsg = "<binary data>"
			}
			fmt.Printf("%s: %s [%s]\n", friend, strMsg, time.Now().String())
		}
	}()

	reader := bufio.NewReader(os.Stdin)
	for {
		msg := inputString(reader, "")
		msgLen := uint64(len(msg))
		if msgLen == 0 {
			continue
		}
		if msgLen > limit {
			fmt.Printf("ERROR: msglen(%d) > limit(%d)\n", msgLen, limit)
			continue
		}
		go func() { _ = hlsClient.SendNotification(pCtx, friend, []byte(msg)) }()
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
