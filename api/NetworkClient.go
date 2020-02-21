package api

import (
	"time"
	"bytes"
	"crypto/x509"
	"encoding/pem"
	"encoding/hex"
	"encoding/json"
	"github.com/number571/gopeer"
	"github.com/number571/hiddenlake/db"
	"github.com/number571/hiddenlake/models"
	"github.com/number571/hiddenlake/settings"
	"github.com/number571/hiddenlake/utils"
	"net/http"
	"os"
	"strings"
)

func NetworkClient(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var data struct {
		State string `json:"state"`
	}

	switch r.Method {
	case "GET":
		clientGET(w, r)
	case "POST":
		clientPOST(w, r)
	case "PATCH":
		clientPATCH(w, r)
	case "DELETE":
		clientDELETE(w, r)
	default:
		data.State = "Method should be GET, POST, PATCH or DELETE"
		json.NewEncoder(w).Encode(data)
	}
}

// Get client public information.
func clientGET(w http.ResponseWriter, r *http.Request) {
	var data struct {
		State     string `json:"state"`
		Info models.Connect `json:"info"`
	}

	var read struct {
		Hashname string `json:"hashname"`
	}

	read.Hashname = strings.Replace(r.URL.Path, "/api/network/client/", "", 1)
	hashname := strings.Split(read.Hashname, "/archive/")[0]
	hashname  = strings.Split(hashname, "/connects")[0]

	var (
		token string
		client = new(gopeer.Client)
	)
	switch {
	case isTokenAuthError(w, r, &token): return
	case isLifeTokenError(w, r, token): return
	case isGetClientError(w, r, client, token): return
	}

	user := settings.Users[token]
	clientData := db.GetClient(user, hashname)
	if clientData == nil {
		data.State = "Client undefined"
		json.NewEncoder(w).Encode(data)
		return
	}

	if strings.Contains(read.Hashname, "/archive/") {
		splited := strings.Split(read.Hashname, "/archive/")
		clientArchiveGET(w, r, user, client, splited)
		return
	}

	data.Info = models.Connect{
		Hidden: gopeer.HashPublic(clientData.Public) != gopeer.HashPublic(clientData.ThrowClient),
		Connected:   client.InConnections(hashname),
		Address:     clientData.Address,
		Hashname:    hashname,
		Public:      gopeer.StringPublic(clientData.Public),
		ThrowClient: gopeer.HashPublic(clientData.ThrowClient),
		Certificate: clientData.Certificate,
	}
	json.NewEncoder(w).Encode(data)
}

// Get list info of files / file from another node.
func clientArchiveGET(w http.ResponseWriter, r *http.Request, user *models.User, client *gopeer.Client, splited []string) {
	var data struct {
		State string        `json:"state"`
		Files []models.File `json:"files"`
	}

	hashname := splited[0]
	filehash := splited[1]

	switch {
	case isNotInConnectionsError(w, r, client, hashname): return
	}

	dest := client.Destination(hashname)

	if filehash == "" {
		_, err := client.SendTo(dest, &gopeer.Package{
			Head: gopeer.Head{
				Title:  settings.TITLE_ARCHIVE,
				Option: settings.OPTION_GET,
			},
		})
		if err != nil {
			data.State = "User can't receive message"
			json.NewEncoder(w).Encode(data)
			return
		}

		select {
		case <-client.Connections[hashname].Chans.Action:
			// pass
		case <-time.After(time.Duration(settings.WAITING_TIME) * time.Second):
			data.State = "Files not loaded"
			json.NewEncoder(w).Encode(data)
			return
		}
		data.Files = user.Temp.FileList
		json.NewEncoder(w).Encode(data)
		return
	}

	_, err := client.SendTo(dest, &gopeer.Package{
		Head: gopeer.Head{
			Title:  settings.TITLE_ARCHIVE,
			Option: settings.OPTION_GET,
		},
		Body: gopeer.Body{
			Data: filehash,
		},
	})
	if err != nil {
		data.State = "User can't receive message"
		json.NewEncoder(w).Encode(data)
		return
	}

	select {
	case <-client.Connections[hashname].Chans.Action:
		// pass
	case <-time.After(time.Duration(settings.WAITING_TIME) * time.Second):
		data.State = "Files not loaded"
		json.NewEncoder(w).Encode(data)
		return
	}
	data.Files = user.Temp.FileList
	json.NewEncoder(w).Encode(data)
}

// Find hidden connection and connect.
func clientPATCH(w http.ResponseWriter, r *http.Request) {
	var data struct {
		State string `json:"state"`
	}

	var read struct {
		PublicKey string `json:"public_key"`
	}

	var (
		token string
		client = new(gopeer.Client)
	)
	switch {
	case isTokenAuthError(w, r, &token): return
	case isLifeTokenError(w, r, token): return
	case isDecodeError(w, r, &read): return
	case isGetClientError(w, r, client, token): return
	}

	user := settings.Users[token]
	public := gopeer.ParsePublic(read.PublicKey)
	if public == nil {
		data.State = "Error decode public key"
		json.NewEncoder(w).Encode(data)
		return
	}

	dest := &gopeer.Destination{
        Receiver: public,
    }

	err := client.Connect(dest)
	if err != nil {
		data.State = "Connect error"
		json.NewEncoder(w).Encode(data)
		return
	}

	hash := gopeer.HashPublic(public)
	err = db.SetClient(user, &models.Client{
		Hashname:    hash,
		Address:     client.Connections[hash].Address,
		Public:      public,
		ThrowClient: client.Connections[hash].ThrowClient,
		Certificate: string(client.Connections[hash].Certificate),
	})
	if err != nil {
		data.State = "Set client error"
		json.NewEncoder(w).Encode(data)
		return
	}

	message := "connection created"
	_, err = client.SendTo(dest, &gopeer.Package{
		Head: gopeer.Head{
			Title:  settings.TITLE_MESSAGE,
			Option: settings.OPTION_GET,
		},
		Body: gopeer.Body{
			Data: message,
		},
	})
	if err != nil {
		data.State = "User can't receive message"
		json.NewEncoder(w).Encode(data)
		return
	}

	err = db.SetChat(user, &models.Chat{
		Companion: hash,
		Messages: []models.Message{
			models.Message{
				Name: hash,
				Text: message,
				Time: utils.CurrentTime(),
			},
		},
	})
	if err != nil {
		data.State = "Set chat error"
		json.NewEncoder(w).Encode(data)
		return
	}

	json.NewEncoder(w).Encode(data)
}

// Connect to another client.
func clientPOST(w http.ResponseWriter, r *http.Request) {
	var data struct {
		State string `json:"state"`
	}

	hashname := strings.Replace(r.URL.Path, "/api/network/client/", "", 1)

	if strings.Contains(hashname, "/archive/") {
		hashname = strings.Split(hashname, "/archive/")[0]
		clientArchivePOST(w, r, hashname)
		return
	}

	if strings.Contains(hashname, "/connects") {
		hashname = strings.Split(hashname, "/connects")[0]
		clientConnectsPOST(w, r, hashname)
		return
	}

	var read struct {
		Address     string `json:"address"`
		Certificate string `json:"certificate"`
		PublicKey   string `json:"public_key"`
	}

	var (
		token string
		client = new(gopeer.Client)
	)
	switch {
	case isTokenAuthError(w, r, &token): return
	case isLifeTokenError(w, r, token): return
	case isDecodeError(w, r, &read): return
	case isGetClientError(w, r, client, token): return
	}

	user := settings.Users[token]
	if len(strings.Split(read.Address, ":")) != 2 {
		data.State = "Address is not corrected"
		json.NewEncoder(w).Encode(data)
		return
	}

	public := gopeer.ParsePublic(read.PublicKey)
	if public == nil {
		data.State = "Error decode public key"
		json.NewEncoder(w).Encode(data)
		return
	}

	block, _ := pem.Decode([]byte(read.Certificate))
	if block == nil {
		data.State = "Failed to parse certificate PEM"
		json.NewEncoder(w).Encode(data)
		return
	}

	_, err := x509.ParseCertificate(block.Bytes)
	if err != nil {
		data.State = "Failed to parse certificate"
		json.NewEncoder(w).Encode(data)
		return
	}

	dest := &gopeer.Destination{
        Address:     read.Address,
        Certificate: []byte(read.Certificate),
        Public:      public,
    }

	err = client.Connect(dest)
	if err != nil {
		data.State = "Connect error"
		json.NewEncoder(w).Encode(data)
		return
	}

	hash := gopeer.HashPublic(public)
	err = db.SetClient(user, &models.Client{
		Hashname:    hash,
		Address:     read.Address,
		Public:      public,
		Certificate: read.Certificate,
	})
	if err != nil {
		data.State = "Set client error"
		json.NewEncoder(w).Encode(data)
		return
	}

	message := "connection created"
	_, err = client.SendTo(dest, &gopeer.Package{
		Head: gopeer.Head{
			Title:  settings.TITLE_MESSAGE,
			Option: settings.OPTION_GET,
		},
		Body: gopeer.Body{
			Data: message,
		},
	})
	if err != nil {
		data.State = "User can't receive message"
		json.NewEncoder(w).Encode(data)
		return
	}

	err = db.SetChat(user, &models.Chat{
		Companion: hash,
		Messages: []models.Message{
			models.Message{
				Name: hash,
				Text: message,
				Time: utils.CurrentTime(),
			},
		},
	})
	if err != nil {
		data.State = "Set chat error"
		json.NewEncoder(w).Encode(data)
		return
	}

	json.NewEncoder(w).Encode(data)
}

// Hidden connect throw node.
func clientConnectsPOST(w http.ResponseWriter, r *http.Request, hashname string) {
	var data struct {
		State string `json:"state"`
	}

	var read struct {
		PublicKey string `json:"public_key"`
	}

	var (
		token string
		client = new(gopeer.Client)
	)
	switch {
	case isTokenAuthError(w, r, &token): return
	case isLifeTokenError(w, r, token): return
	case isDecodeError(w, r, &read): return
	case isGetClientError(w, r, client, token): return
	case isNotInConnectionsError(w, r, client, hashname): return
	}

	user := settings.Users[token]
	public := gopeer.ParsePublic(read.PublicKey)
	if public == nil {
		data.State = "Error decode public key"
		json.NewEncoder(w).Encode(data)
		return
	}

	dest := &gopeer.Destination{
        Address:     client.Connections[hashname].Address,
        Certificate: client.Connections[hashname].Certificate,
        Public:      client.Connections[hashname].Public,
        Receiver:    public,
    }

    err := client.Connect(dest)
	if err != nil {
		data.State = "Connect error"
		json.NewEncoder(w).Encode(data)
		return
	}

	hash := gopeer.HashPublic(public)
	err = db.SetClient(user, &models.Client{
		Hashname: hash,
		Certificate: string(client.Connections[hashname].Certificate),
		Address: client.Connections[hashname].Address,
		Public: public,
		ThrowClient: client.Connections[hashname].Public,
	})
	if err != nil {
		data.State = "Set client error"
		json.NewEncoder(w).Encode(data)
		return
	}

	message := "connection created"
	_, err = client.SendTo(dest, &gopeer.Package{
		Head: gopeer.Head{
			Title:  settings.TITLE_MESSAGE,
			Option: settings.OPTION_GET,
		},
		Body: gopeer.Body{
			Data: message,
		},
	})
	if err != nil {
		data.State = "User can't receive message"
		json.NewEncoder(w).Encode(data)
		return
	}

	err = db.SetChat(user, &models.Chat{
		Companion: hash,
		Messages: []models.Message{
			models.Message{
				Name: hash,
				Text: message,
				Time: utils.CurrentTime(),
			},
		},
	})
	if err != nil {
		data.State = "Set chat error"
		json.NewEncoder(w).Encode(data)
		return
	}

	json.NewEncoder(w).Encode(data)
}

// Install file from another node.
func clientArchivePOST(w http.ResponseWriter, r *http.Request, hashname string) {
	var data struct {
		State string `json:"state"`
	}

	var read struct {
		Filehash string `json:"filehash"`
	}

	var (
		client = new(gopeer.Client)
		token string
	)
	switch {
	case isTokenAuthError(w, r, &token): return
	case isLifeTokenError(w, r, token): return
	case isDecodeError(w, r, &read): return
	case isGetClientError(w, r, client, token): return
	case isNotInConnectionsError(w, r, client, hashname): return
	}

	user := settings.Users[token]
	dest := client.Destination(hashname)

	_, err := client.SendTo(dest, &gopeer.Package{
		Head: gopeer.Head{
			Title:  settings.TITLE_ARCHIVE,
			Option: settings.OPTION_GET,
		},
		Body: gopeer.Body{
			Data: read.Filehash,
		},
	})
	if err != nil {
		data.State = "User can't receive message"
		json.NewEncoder(w).Encode(data)
		return
	}

	select {
	case <-client.Connections[hashname].Chans.Action:
		// pass
	case <-time.After(time.Duration(settings.WAITING_TIME) * time.Second):
		data.State = "File not loaded"
		json.NewEncoder(w).Encode(data)
		return
	}

	if len(user.Temp.FileList) == 0 {
		data.State = "File not found"
		json.NewEncoder(w).Encode(data)
		return
	}

	hash, err := hex.DecodeString(user.Temp.FileList[0].Hash)
	if err != nil {
		data.State = "Error decode hex format"
		json.NewEncoder(w).Encode(data)
		return
	}

	pathhash := hex.EncodeToString(gopeer.HashSum(bytes.Join(
		[][]byte{
			hash,
			gopeer.HashSum(gopeer.GenerateRandomBytes(16)),
			gopeer.Base64Decode(user.Hashname),
		},
		[]byte{},
	)))

	file := db.GetFile(user, user.Temp.FileList[0].Hash)
	if file != nil {
		data.State = "This file already exist"
		json.NewEncoder(w).Encode(data)
		return
	}

	client.LoadFile(dest, user.Temp.FileList[0].Path, settings.PATH_ARCHIVE+pathhash)

	output, err := os.Open(settings.PATH_ARCHIVE + pathhash)
	if err != nil {
		data.State = "Installed file does not open"
		json.NewEncoder(w).Encode(data)
		return
	}
	var (
		checkhash = make([]byte, 32)
		buffer    = make([]byte, settings.BUFFER_SIZE)
	)
	for {
		length, err := output.Read(buffer)
		if err != nil {
			break
		}
		checkhash = gopeer.HashSum(bytes.Join(
			[][]byte{checkhash, buffer[:length]},
			[]byte{},
		))
	}
	output.Close()

	if hex.EncodeToString(checkhash) != user.Temp.FileList[0].Hash {
		os.Remove(settings.PATH_ARCHIVE + pathhash)
		data.State = "Hashes not equal"
		json.NewEncoder(w).Encode(data)
		return
	}

	db.SetFile(user, &models.File{
		Name: user.Temp.FileList[0].Name,
		Hash: user.Temp.FileList[0].Hash,
		Path: pathhash,
		Size: user.Temp.FileList[0].Size,
	})

	json.NewEncoder(w).Encode(data)
}

// Disconnect from client.
func clientDELETE(w http.ResponseWriter, r *http.Request) {
	var data struct {
		State string `json:"state"`
	}

	var read struct {
		Hashname string `json:"hashname"`
	}

	var (
		client = new(gopeer.Client)
		token string
	)
	switch {
	case isTokenAuthError(w, r, &token): return
	case isLifeTokenError(w, r, token): return
	case isDecodeError(w, r, &read): return
	case isGetClientError(w, r, client, token): return
	case isNotInConnectionsError(w, r, client, read.Hashname): return
	}

	dest := client.Destination(read.Hashname) 
	
	message := "connection closed"
	_, err := client.SendTo(dest, &gopeer.Package{
		Head: gopeer.Head{
			Title:  settings.TITLE_MESSAGE,
			Option: settings.OPTION_GET,
		},
		Body: gopeer.Body{
			Data: message,
		},
	})
	if err != nil {
		data.State = "User can't receive message"
		json.NewEncoder(w).Encode(data)
		return
	}

	db.SetChat(settings.Users[token], &models.Chat{
		Companion: read.Hashname,
		Messages: []models.Message{
			models.Message{
				Name: client.Hashname,
				Text: message,
				Time: utils.CurrentTime(),
			},
		},
	})
	client.Disconnect(dest)

	json.NewEncoder(w).Encode(data)
}