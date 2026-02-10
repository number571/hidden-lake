package main

import (
	"flag"
	"net/http"

	"github.com/number571/hidden-lake/build"
	hlk_client "github.com/number571/hidden-lake/pkg/api/kernel/client"
	hlf_client "github.com/number571/hidden-lake/pkg/api/services/filesharer/client"
	hlm_client "github.com/number571/hidden-lake/pkg/api/services/messenger/client"
	hlp_client "github.com/number571/hidden-lake/pkg/api/services/pinger/client"
)

var (
	hlkClient hlk_client.IClient
	hlmClient hlm_client.IClient
	hlpClient hlp_client.IClient
	hlfClient hlf_client.IClient
)

func init() {
	kernelAddr := flag.String("kernel", "localhost:9572", "address of the HLK")
	flag.StringVar(kernelAddr, "k", *kernelAddr, "alias for -kernel")

	messengerAddr := flag.String("messenger", "localhost:9591", "address of the HLS=messenger")
	flag.StringVar(messengerAddr, "m", *messengerAddr, "alias for -messenger")

	pingerAddr := flag.String("pinger", "localhost:9551", "address of the HLS=pinger")
	flag.StringVar(pingerAddr, "p", *pingerAddr, "alias for -pinger")

	filesharerAddr := flag.String("filesharer", "localhost:9541", "address of the HLS=filesharer")
	flag.StringVar(filesharerAddr, "f", *filesharerAddr, "alias for -filesharer")

	flag.Parse()
	initClients(*kernelAddr, *messengerAddr, *pingerAddr, *filesharerAddr)
}

func initClients(kernelAddr, messengerAddr, pingerAddr, filesharerAddr string) {
	hlkClient = hlk_client.NewClient(
		hlk_client.NewBuilder(),
		hlk_client.NewRequester(
			kernelAddr,
			&http.Client{Timeout: build.GetSettings().GetHttpCallbackTimeout()},
		),
	)

	hlmClient = hlm_client.NewClient(
		hlm_client.NewRequester(
			messengerAddr,
			&http.Client{Timeout: build.GetSettings().GetHttpCallbackTimeout()},
		),
	)

	hlpClient = hlp_client.NewClient(
		hlp_client.NewRequester(
			pingerAddr,
			&http.Client{Timeout: build.GetSettings().GetHttpCallbackTimeout()},
		),
	)

	hlfClient = hlf_client.NewClient(
		hlf_client.NewRequester(
			filesharerAddr,
			&http.Client{Timeout: build.GetSettings().GetHttpCallbackTimeout()},
		),
	)
}
