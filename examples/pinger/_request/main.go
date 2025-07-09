package main

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"time"
)

const (
	cRequestTemplate = `{
        "receiver":"%s",
        "req_data":{
			"method":"GET",
			"host":"hidden-lake-service=pinger",
			"path":"/ping"
		}
	}`
)

func main() {
	t1 := time.Now()
	defer func() {
		t2 := time.Now()
		diff := t2.Sub(t1)
		fmt.Println("Request took", diff)
	}()

	receiver := "Bob"
	sendMessage(receiver)
}

func sendMessage(pReceiver string) {
	httpClient := http.Client{Timeout: time.Hour}

	requestData := fmt.Sprintf(cRequestTemplate, pReceiver)
	req, err := http.NewRequest(
		http.MethodPost,
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
