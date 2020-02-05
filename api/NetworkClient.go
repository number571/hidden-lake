package api

import (
	"bytes"
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
	case "DELETE":
		clientDELETE(w, r)
	default:
		data.State = "Method should be GET"
		json.NewEncoder(w).Encode(data)
	}
}

// Disconnect from client.
func clientDELETE(w http.ResponseWriter, r *http.Request) {
	var data struct {
		State string `json:"state"`
	}

	var read struct {
		Hashname string `json:"hashname"`
	}

	err := json.NewDecoder(r.Body).Decode(&read)
	if err != nil {
		data.State = "Error decode json format"
		json.NewEncoder(w).Encode(data)
		return
	}

	token := r.Header.Get("Authorization")
	token = strings.Replace(token, "Bearer ", "", 1)
	if _, ok := settings.Users[token]; !ok {
		data.State = "Tokened user undefined"
		json.NewEncoder(w).Encode(data)
		return
	}

	err = settings.CheckLifetimeToken(token)
	if err != nil {
		data.State = "Token lifetime is over"
		json.NewEncoder(w).Encode(data)
		return
	} else {
		settings.Users[token].Session.Time = utils.CurrentTime()
	}

	hash := settings.Users[token].Hashname
	client, ok := settings.Listener.Clients[hash]
	if !ok {
		data.State = "Current client is not exist"
		json.NewEncoder(w).Encode(data)
		return
	}

	if !client.InConnections(read.Hashname) {
		data.State = "User is not connected"
		json.NewEncoder(w).Encode(data)
		return
	}

	dest := &gopeer.Destination{
		Address: client.Connections[read.Hashname].Address,
		Public:  client.Connections[read.Hashname].Public,
	}

	message := "connection closed"
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

	db.SetChat(settings.Users[token], &models.Chat{
		Companion: read.Hashname,
		Messages: []models.Message{
			models.Message{
				Name: hash,
				Text: message,
				Time: utils.CurrentTime(),
			},
		},
	})
	client.Disconnect(dest)

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

	var read struct {
		Address   string `json:"address"`
		PublicKey string `json:"public_key"`
	}

	err := json.NewDecoder(r.Body).Decode(&read)
	if err != nil {
		data.State = "Error decode json format"
		json.NewEncoder(w).Encode(data)
		return
	}

	token := r.Header.Get("Authorization")
	token = strings.Replace(token, "Bearer ", "", 1)
	if _, ok := settings.Users[token]; !ok {
		data.State = "Tokened user undefined"
		json.NewEncoder(w).Encode(data)
		return
	}

	err = settings.CheckLifetimeToken(token)
	if err != nil {
		data.State = "Token lifetime is over"
		json.NewEncoder(w).Encode(data)
		return
	} else {
		settings.Users[token].Session.Time = utils.CurrentTime()
	}

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

	user := settings.Users[token]
	client, ok := settings.Listener.Clients[user.Hashname]
	if !ok {
		data.State = "Current client is not exist"
		json.NewEncoder(w).Encode(data)
		return
	}

	dest := &gopeer.Destination{
		Address: read.Address,
		Public:  public,
	}
	err = client.Connect(dest)
	if err != nil {
		data.State = "Connect error"
		json.NewEncoder(w).Encode(data)
		return
	}

	hash := gopeer.HashPublic(public)
	err = db.SetClient(user, &models.Client{
		Hashname: hash,
		Address:  read.Address,
		Public:   public,
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

	err := json.NewDecoder(r.Body).Decode(&read)
	if err != nil {
		data.State = "Error decode json format"
		json.NewEncoder(w).Encode(data)
		return
	}

	token := r.Header.Get("Authorization")
	token = strings.Replace(token, "Bearer ", "", 1)
	if _, ok := settings.Users[token]; !ok {
		data.State = "Tokened user undefined"
		json.NewEncoder(w).Encode(data)
		return
	}

	err = settings.CheckLifetimeToken(token)
	if err != nil {
		data.State = "Token lifetime is over"
		json.NewEncoder(w).Encode(data)
		return
	} else {
		settings.Users[token].Session.Time = utils.CurrentTime()
	}

	user := settings.Users[token]
	client, ok := settings.Listener.Clients[user.Hashname]
	if !ok {
		data.State = "Current client is not exist"
		json.NewEncoder(w).Encode(data)
		return
	}

	if !client.InConnections(hashname) {
		data.State = "User is not connected"
		json.NewEncoder(w).Encode(data)
		return
	}

	dest := &gopeer.Destination{
		Address: client.Connections[hashname].Address,
		Public:  client.Connections[hashname].Public,
	}

	_, err = client.SendTo(dest, &gopeer.Package{
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

	if len(user.FileList) == 0 {
		data.State = "File not found"
		json.NewEncoder(w).Encode(data)
		return
	}

	hash, err := hex.DecodeString(user.FileList[0].Hash)
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

	file := db.GetFile(user, user.FileList[0].Hash)
	if file != nil {
		data.State = "This file already exist"
		json.NewEncoder(w).Encode(data)
		return
	}

	client.LoadFile(dest, user.FileList[0].Path, settings.PATH_ARCHIVE+pathhash)

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

	if hex.EncodeToString(checkhash) != user.FileList[0].Hash {
		os.Remove(settings.PATH_ARCHIVE + pathhash)
		data.State = "Hashes not equal"
		json.NewEncoder(w).Encode(data)
		return
	}

	db.SetFile(user, &models.File{
		Name: user.FileList[0].Name,
		Hash: user.FileList[0].Hash,
		Path: pathhash,
		Size: user.FileList[0].Size,
	})

	json.NewEncoder(w).Encode(data)
}

// Get client public information.
func clientGET(w http.ResponseWriter, r *http.Request) {
	var data struct {
		Connected bool   `json:"connected"`
		Address   string `json:"address"`
		Hashname  string `json:"hashname"`
		PublicKey string `json:"public_key"`
		State     string `json:"state"`
	}

	var read struct {
		Hashname string `json:"hashname"`
	}

	read.Hashname = strings.Replace(r.URL.Path, "/api/network/client/", "", 1)
	splited := strings.Split(read.Hashname, "/archive/")

	hashname := splited[0]

	token := r.Header.Get("Authorization")
	token = strings.Replace(token, "Bearer ", "", 1)
	if _, ok := settings.Users[token]; !ok {
		data.State = "Tokened user undefined"
		json.NewEncoder(w).Encode(data)
		return
	}

	err := settings.CheckLifetimeToken(token)
	if err != nil {
		data.State = "Token lifetime is over"
		json.NewEncoder(w).Encode(data)
		return
	} else {
		settings.Users[token].Session.Time = utils.CurrentTime()
	}

	user := settings.Users[token]
	clientData := db.GetClient(user, hashname)
	if clientData == nil {
		data.State = "Client undefined"
		json.NewEncoder(w).Encode(data)
		return
	}

	client, ok := settings.Listener.Clients[user.Hashname]
	if !ok {
		data.State = "Current client is not exist"
		json.NewEncoder(w).Encode(data)
		return
	}

	if len(splited) == 2 {
		clientArchiveGET(w, r, user, client, splited)
		return
	}

	if client.InConnections(hashname) {
		data.Connected = true
	}

	data.Address = clientData.Address
	data.Hashname = hashname
	data.PublicKey = gopeer.StringPublic(clientData.Public)

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

	if !client.InConnections(hashname) {
		data.State = "User is not connected"
		json.NewEncoder(w).Encode(data)
		return
	}

	dest := &gopeer.Destination{
		Address: client.Connections[hashname].Address,
		Public:  client.Connections[hashname].Public,
	}

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

		data.Files = user.FileList
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

	data.Files = user.FileList
	json.NewEncoder(w).Encode(data)
}
