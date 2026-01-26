package main

import (
	"bufio"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"os"
	"strings"

	"github.com/number571/go-peer/pkg/crypto/asymmetric"
	"github.com/number571/hidden-lake/build"
	"github.com/number571/hidden-lake/internal/kernel/pkg/settings"
	"github.com/number571/hidden-lake/internal/utils/flag"
	"github.com/number571/hidden-lake/internal/utils/help"
	hlk_client "github.com/number571/hidden-lake/pkg/api/kernel/client"
	"github.com/number571/hidden-lake/pkg/api/kernel/client/proc"
	"github.com/number571/hidden-lake/pkg/network/request"
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
		flag.NewFlagBuilder("-d", "--do").
			WithDescription("set runner [send|fetch|pubkey|get-onlines|del-online|get-connections|add-connection|del-connection|get-friends|add-friend|del-friend]").
			WithDefinedValue(""),
		flag.NewFlagBuilder("-a", "--arg").
			WithDescription("set argument for runner <send|fetch|del-online|add-connection|del-connection|add-friend|del-friend>").
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
		fmt.Println(err)
		os.Exit(1)
	}
}

func runFunction(pCtx context.Context, pArgs []string) error {
	hlkClient := hlk_client.NewClient(
		hlk_client.NewBuilder(),
		hlk_client.NewRequester(
			gFlags.Get("-k").GetStringValue(pArgs),
			&http.Client{Timeout: build.GetSettings().GetHttpCallbackTimeout()},
		),
	)

	reader := bufio.NewReader(os.Stdin)
	do := gFlags.Get("-d").GetStringValue(pArgs)
	switch do {

	// NETWORK
	case "send":
		friend := gFlags.Get("-a").GetStringValue(pArgs)
		req, err := request.LoadRequest(inputString(reader))
		if err != nil {
			return err
		}
		if err := hlkClient.SendRequest(pCtx, friend, req); err != nil {
			return err
		}
		fmt.Println("done!")
	case "fetch":
		friend := gFlags.Get("-a").GetStringValue(pArgs)
		req, err := request.LoadRequest(inputString(reader))
		if err != nil {
			return err
		}
		rsp, err := hlkClient.FetchRequest(pCtx, friend, req)
		if err != nil {
			return err
		}
		fmt.Println(rsp.ToString())

	// PROFILE
	case "pubkey":
		pubKey, err := hlkClient.GetPubKey(pCtx)
		if err != nil {
			return err
		}
		fmt.Println(pubKey.ToString())

	// ONLINES
	case "get-onlines":
		onlines, err := hlkClient.GetOnlines(pCtx)
		if err != nil {
			return err
		}
		fmt.Println(serializeJSON(onlines))
	case "del-online":
		online := gFlags.Get("-a").GetStringValue(pArgs)
		if err := hlkClient.DelOnline(pCtx, online); err != nil {
			return err
		}
		fmt.Println("done!")

	// CONNECTIONS
	case "get-connections":
		connections, err := hlkClient.GetConnections(pCtx)
		if err != nil {
			return err
		}
		fmt.Println(serializeJSON(connections))
	case "del-connection":
		connection := gFlags.Get("-a").GetStringValue(pArgs)
		if err := hlkClient.DelConnection(pCtx, connection); err != nil {
			return err
		}
		fmt.Println("done!")
	case "add-connection":
		connection := gFlags.Get("-a").GetStringValue(pArgs)
		if err := hlkClient.AddConnection(pCtx, connection); err != nil {
			return err
		}
		fmt.Println("done!")

	// FRIENDS
	case "get-friends":
		friends, err := hlkClient.GetFriends(pCtx)
		if err != nil {
			return err
		}
		fmt.Println(serializeJSON(proc.FriendsMapToList(friends)))
	case "del-friend":
		friend := gFlags.Get("-a").GetStringValue(pArgs)
		if err := hlkClient.DelFriend(pCtx, friend); err != nil {
			return err
		}
		fmt.Println("done!")
	case "add-friend":
		friend := gFlags.Get("-a").GetStringValue(pArgs)
		pubKey := asymmetric.LoadPubKey(inputString(reader))
		if pubKey == nil {
			return errors.New("load public key") // nolint: err113
		}
		if err := hlkClient.AddFriend(pCtx, friend, pubKey); err != nil {
			return err
		}
		fmt.Println("done!")

	default:
		return errors.Join(ErrUnknownAction, errors.New(do)) // nolint:err113
	}

	return nil
}

func serializeJSON(pData interface{}) string {
	res, _ := json.MarshalIndent(pData, "", "\t")
	return string(res)
}

func inputString(reader *bufio.Reader) string {
	rawInput, err := reader.ReadString('\n')
	if err != nil {
		panic(err)
	}
	return strings.TrimSpace(rawInput)
}
