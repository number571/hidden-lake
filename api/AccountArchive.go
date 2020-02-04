package api

import (
    "os"
    // "fmt"
    "bytes"
    "strings"
    "net/http"
    "encoding/hex"
    "encoding/json"
    "github.com/number571/gopeer"
    "../db"
    "../utils"
    "../models"
    "../settings"
)

func AccountArchive(w http.ResponseWriter, r *http.Request) {
    w.Header().Set("Content-Type", "application/json")

    var data struct {
        State string `json:"state"`
    }

    switch r.Method {
    case "GET":
        accountArchiveGET(w, r)
    case "PUT":
        accountArchivePUT(w, r)
    case "DELETE":
        accountArchiveDELETE(w, r)
    default:
        data.State = "Method should be GET, PUT, POST or DELETE"
        json.NewEncoder(w).Encode(data)
    }
}

// Delete file.
func accountArchiveDELETE(w http.ResponseWriter, r *http.Request) {
    var data struct {
        State   string  `json:"state"`
        Files []*models.File `json:"files"`
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
    err = db.DeleteFile(user, read.Filehash)
    if err != nil {
        data.State = "File already deleted"
        json.NewEncoder(w).Encode(data)
        return
    }

    json.NewEncoder(w).Encode(data)
}

// Upload file.
func accountArchivePUT(w http.ResponseWriter, r *http.Request) {
    var data struct {
        State string `json:"state"`
        Filehash string `json:"filehash"`
    }

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

    r.ParseMultipartForm(settings.PACKAGE_SIZE)

    input, handler, err := r.FormFile("uploadfile")
    if err != nil {
        data.State = "Error read upload file"
        json.NewEncoder(w).Encode(data)
        return
    }

    if len(handler.Filename) > 128 {
        data.State = "File length should be <= 128 chars"
        json.NewEncoder(w).Encode(data)
        return
    }

    tempname := utils.RandomString(16)
    output, err := os.OpenFile(
        settings.PATH_ARCHIVE + tempname, 
        os.O_WRONLY | os.O_CREATE, 
        0666,
    )
    if err != nil {
        data.State = "Error push file to archive"
        json.NewEncoder(w).Encode(data)
        return
    }

    var (
        size = uint64(0)
        hash = make([]byte, 32)
        buffer = make([]byte, settings.BUFFER_SIZE)
    )

    for {
        length, err := input.Read(buffer)
        if err != nil {
            break
        }
        size += uint64(length)
        hash = gopeer.HashSum(bytes.Join(
            [][]byte{hash, buffer[:length]},
            []byte{},
        ))
        output.Write(buffer[:length])
    }

    output.Close()
    input.Close()

    user := settings.Users[token]
    filehash := hex.EncodeToString(hash)

    file := db.GetFile(user, filehash)
    if file != nil {
        os.Remove(settings.PATH_ARCHIVE + tempname)
        data.State = "This file already exist"
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

    os.Rename(
        settings.PATH_ARCHIVE + tempname, 
        settings.PATH_ARCHIVE + pathhash,
    )

    db.SetFile(user, &models.File{
        Name: handler.Filename,
        Hash: filehash,
        Path: pathhash,
        Size: size,
    })

    data.Filehash = filehash
    json.NewEncoder(w).Encode(data)
}

// List of files / file information.
func accountArchiveGET(w http.ResponseWriter, r *http.Request) {
    var data struct {
        State   string  `json:"state"`
        Files []models.File `json:"files"`
    }

    var read struct {
        Filehash string `json:"filehash"`
    }

    read.Filehash = strings.Replace(r.URL.Path, "/api/account/archive/", "", 1)

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
    switch read.Filehash {
    case "", "null", "undefined":
        data.Files = db.GetAllFiles(user)
    default:
        file := db.GetFile(user, read.Filehash)
        if file == nil {
            data.State = "File undefined"
            json.NewEncoder(w).Encode(data)
            return
        }
        data.Files = append(data.Files, *file)
    }
    json.NewEncoder(w).Encode(data)
}
