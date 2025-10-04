package main

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
)

type sResponse struct {
	FEcho   string `json:"echo"`
	FReturn int    `json:"return"`
}

func main() {
	http.HandleFunc("/echo", echoPage)
	http.ListenAndServe(":8080", nil)
}

func echoPage(w http.ResponseWriter, r *http.Request) {
	log.Printf("[%s] %s\n", r.Method, r.Header.Get("Hlk-Sender-Name"))
	if r.Method != http.MethodPost {
		response(w, 2, "failed: incorrect method")
		return
	}
	res, err := io.ReadAll(r.Body)
	if err != nil {
		response(w, 3, "failed: read body")
		return
	}
	response(w, 1, string(res))
}

func response(w http.ResponseWriter, ret int, res string) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(&sResponse{
		FEcho:   res,
		FReturn: ret,
	})
}
