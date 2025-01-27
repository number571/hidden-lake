package main

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"io"
	"net/http"
	"time"
)

const (
	cRequestTemplate = `{
        "receiver":"%s",
        "req_data":{
			"method":"POST",
			"host":"hidden-lake-notifier",
			"path":"/notify",
			"head": {
				"Hl-Notifier-Public-Key-Hash": "999a5e25460dea77bec1113a54155576d32e718ba3dae8f80d2ac3ca0de6e5e6e6e812fa942d7ea192534a4ffb76281b"
			},
			"body":"%s"
		}
	}`
)

func main() {
	receiver := "Bob"
	message := "hello, world!"

	sendMessage(receiver, message)
}

func sendMessage(pReceiver string, pMessage string) {
	httpClient := http.Client{Timeout: time.Minute / 2}

	requestData := fmt.Sprintf(
		cRequestTemplate,
		pReceiver,
		base64.StdEncoding.EncodeToString([]byte(pMessage)),
	)

	req, err := http.NewRequest(
		http.MethodPut,
		"http://localhost:7572/api/network/request",
		bytes.NewBufferString(requestData),
	)
	if err != nil {
		panic(err)
	}

	resp, err := httpClient.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	res, err := io.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}

	fmt.Println(string(res))
}
