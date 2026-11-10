package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"strconv"

	"github.com/number571/hidden-lake/build"
	"github.com/number571/hidden-lake/internal/services/publisher/pkg/settings"
	"github.com/number571/hidden-lake/internal/utils/flag"
	"github.com/number571/hidden-lake/internal/utils/help"

	hls_client "github.com/number571/hidden-lake/pkg/api/services/publisher/client"
	"github.com/number571/hidden-lake/pkg/api/services/publisher/client/dto"
)

var (
	gFlags = flag.NewFlagsBuilder(
		flag.NewFlagBuilder("-v", "--version").
			WithDescription("print version of application"),
		flag.NewFlagBuilder("-h", "--help").
			WithDescription("print information about application"),
		flag.NewFlagBuilder("-s", "--service").
			WithDescription("set internal address of the HLS").
			WithDefinedValue("localhost:9551"),
		flag.NewFlagBuilder("-f", "--friend").
			WithDescription("set alias name of the friend").
			WithDefinedValue(""),
		flag.NewFlagBuilder("-r", "--remote").
			WithDescription("use remote alg option"),
		flag.NewFlagBuilder("-p", "--public").
			WithDescription("use public alg option"),
		flag.NewFlagBuilder("-d", "--do").
			WithDescription("set runner [size|push|load-by-uid|load-by-nid|listen]").
			WithDefinedValue(""),
		flag.NewFlagBuilder("-a", "--arg").
			WithDescription("set argument for runner <data|uid|nid>").
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
	friend := gFlags.Get("-f").GetStringValue(pArgs)
	do := gFlags.Get("-d").GetStringValue(pArgs)

	var (
		isRemote = gFlags.Get("-r").GetBoolValue(pArgs)
		isPublic = gFlags.Get("-p").GetBoolValue(pArgs)
	)

	hlsClient := hls_client.NewClient(
		gFlags.Get("-s").GetStringValue(pArgs),
		&http.Client{Timeout: build.GetSettings().GetHttpCallbackTimeout()},
		hls_client.WithRemote(isRemote),
		hls_client.WithPublic(isPublic),
	)

	switch do {
	case "size":
		size, err := hlsClient.GetBoxSize(pCtx, friend)
		if err != nil {
			return err
		}
		fmt.Println(size)
	case "push":
		if isRemote && isPublic {
			fmt.Println("warn: public is disabled with remote option (push runner)")
		}
		data := gFlags.Get("-a").GetStringValue(pArgs)
		uid := dto.NewMessageUID()
		nid, _, err := hlsClient.PushMessageToBox(pCtx, friend, uid, []byte(data))
		if err != nil {
			return err
		}
		fmt.Println(nid, uid)
	case "load-by-uid":
		uid := gFlags.Get("-a").GetStringValue(pArgs)
		nid, msg, err := hlsClient.LoadMessageFromBoxByUID(pCtx, friend, uid)
		if err != nil {
			return err
		}
		fmt.Println(nid)
		fmt.Println(msg.ToString())
	case "load-by-nid":
		nidStr := gFlags.Get("-a").GetStringValue(pArgs)
		nid, err := strconv.ParseUint(nidStr, 10, 64)
		if err != nil {
			return err
		}
		msg, err := hlsClient.LoadMessageFromBoxByEnum(pCtx, friend, nid)
		if err != nil {
			return err
		}
		fmt.Println(msg.ToString())
	case "listen":
		if isRemote {
			fmt.Println("warn: remote is always used (listen runner)")
		}
		if isPublic {
			// public message can't be pushed to another node
			return ErrAvailableOnlyForTypePersonal
		}
		for {
			sender, _, msg, err := hlsClient.ListenMessage(pCtx, "hls-publisher")
			if err != nil {
				fmt.Printf("warn: %s\n", err.Error())
				continue
			}
			if sender != friend {
				continue
			}
			fmt.Println(msg.ToString())
		}
	default:
		return ErrUnknownAction
	}

	return nil
}
