package clients

import (
	"net/http"
	"../settings"
)

func ClientHTTP() {
	settings.ServerListenHTTP = &http.Server{
		Addr: settings.IPV4_HTTP + settings.PORT_HTTP,
	}

	if err := settings.ServerListenHTTP.ListenAndServe(); err != nil {
		settings.ServerListenHTTP = nil
		return
	}
}
