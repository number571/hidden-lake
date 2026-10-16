package main

import (
	"bufio"
	"context"
	"fmt"
	"net/http"
	"os"
	"strings"

	"github.com/number571/hidden-lake/build"
	"github.com/number571/hidden-lake/internal/services/notifier/pkg/settings"
	"github.com/number571/hidden-lake/internal/utils/flag"
	"github.com/number571/hidden-lake/internal/utils/help"
	hls_client "github.com/number571/hidden-lake/pkg/api/services/notifier/client"
	"github.com/number571/hidden-lake/pkg/api/services/notifier/client/dto"
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
		help.Println(settings.CAppFullName+"-cli", settings.CAppDescription, gFlags)
		return
	}

	ctx := context.Background()
	runFunction(ctx, args)
}

func runFunction(pCtx context.Context, pArgs []string) {
	const (
		load = 256
		iam  = "<iam>"
	)

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

	count, err := hlsClient.GetChatSize(pCtx, friend)
	if err != nil {
		fmt.Printf("error: %s\n", err.Error())
		return
	}

	index := uint64(0)
	if count >= load {
		index = count - load
	}

	for ; index < count; index++ {
		msg, err := hlsClient.LoadMessageByEnum(pCtx, friend, index)
		if err != nil {
			fmt.Printf("error: %s\n", err.Error())
			return
		}
		state := ""
		sender := iam
		switch msg.GetHead().GetState() {
		case dto.CIncomingState:
			sender = friend
			state = "<"
		case dto.CReceivedState:
			state = "v"
		case dto.CSendingState:
			go hlsClient.PushMessage(pCtx, friend, msg.GetHead().GetUID(), msg.GetBody().GetData()) // nolint:errcheck
			state = ">"
		}
		fmt.Printf("%s: %s [%s][%s]\n", sender, msg.GetBody().GetData(), state, msg.GetBody().GetTimestamp())
	}

	go func() {
		for {
			m, err := hlsClient.ListenChan(pCtx, friend, "hls-notifier-cli")
			if err != nil {
				fmt.Printf("error: %s\n", err.Error())
				continue
			}
			fmt.Printf("%s: %s [%s]\n", friend, m.GetData(), m.GetTimestamp())
		}
	}()

	reader := bufio.NewReader(os.Stdin)
	for {
		uid := dto.NewMessageUID()
		msg := inputString(reader, "")
		go hlsClient.PushMessage(pCtx, friend, uid, msg) // nolint:errcheck
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
