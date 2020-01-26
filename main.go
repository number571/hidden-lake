package main

/*
   github.com/mattn/go-sqlite3
   github.com/number571/gopeer
   golang.org/x/net/websocket
*/

import (
    "./api"
    "./db"
    "./models"
    "./settings"
    "./utils"
    "./ws"
    "github.com/number571/gopeer"
    "golang.org/x/net/websocket"
    "html/template"
    "net/http"
    "crypto/tls"
    "os"
)

func init() {
    gopeer.Set(gopeer.SettingsType{
        "NETWORK": "[HIDDEN-LAKE]",
        "VERSION": "[1.0.0s]",
        "HMACKEY": "9163571392708145",
        "GENESIS": "[GENESIS-LAKE]",
        "NOISE":   "h19dlI#L9dkc8JA]1s-zSp,Nl/qs4;qf",
    })
    settings.InitializeDB(settings.DB_NAME)
    settings.InitializeCFG(settings.CFG_NAME)
    go settings.ClearUnusedTokens()
}

func main() {
    mux := http.NewServeMux()
    mux.Handle("/static/", http.StripPrefix(
        "/static/",
        handleFileServer(http.Dir(settings.PATH_STATIC))),
    )

    mux.HandleFunc("/", indexPage) // GET
    mux.HandleFunc("/api/login", api.Login) // POST
    mux.HandleFunc("/api/logout", api.Logout) // POST
    mux.HandleFunc("/api/signup", api.Signup) // POST
    mux.HandleFunc("/api/account", api.Account) // GET, POST, DELETE
    mux.HandleFunc("/api/network/", api.Network) // GET, POST, DELETE
    mux.HandleFunc("/api/network/client/", api.Client) // GET, POST, DELETE

    mux.Handle("/ws/network", websocket.Handler(ws.Network))

    handleServerTCP(&settings.CFG.Host.Tcp)
    handleServerHTTP(&settings.CFG.Host.Http, mux)
}

func handleServerTCP(model *models.Tcp) {
    settings.Listener = gopeer.NewListener(model.Ipv4 + model.Port)
    settings.Listener.Open().Run(handleActions)
}

// $ openssl req -new -x509 -nodes -newkey ec:<(openssl ecparam -name secp521r1) -keyout tls/cert.key -out tls/cert.crt -days 3650
func handleServerHTTP(model *models.Http, mux *http.ServeMux) {
    srv := &http.Server{
        Addr:         model.Ipv4+model.Port,
        Handler:      mux,
        TLSConfig:    &tls.Config{
            MinVersion:               tls.VersionTLS12,
            CurvePreferences:         []tls.CurveID{tls.CurveP521},
            PreferServerCipherSuites: true,
            CipherSuites: []uint16{
                tls.TLS_ECDHE_ECDSA_WITH_AES_256_CBC_SHA,
            },
        },
        TLSNextProto: make(map[string]func(*http.Server, *tls.Conn, http.Handler), 0),
    }
    srv.ListenAndServeTLS(model.Tls.Crt, model.Tls.Key)
}

func handleActions(client *gopeer.Client, pack *gopeer.Package) {
    client.HandleAction(settings.TITLE_MESSAGE, pack, getMessage, setMessage)
}

func setMessage(client *gopeer.Client, pack *gopeer.Package) {
    // if package delivered
}

func getMessage(client *gopeer.Client, pack *gopeer.Package) (set string) {
    var (
        hashname = pack.From.Sender.Hashname
        token    = settings.Tokens[client.Hashname]
        user     = settings.Users[token]
        time     = utils.CurrentTime()
    )

    if !db.InClients(user, hashname) {
        db.SetClient(user, &models.Client{
            Hashname: hashname,
            Address:  pack.From.Address,
            Public:   client.Connections[hashname].Public,
        })
    }

    if user.Hashname == hashname {
        return
    }

    db.SetChat(user, &models.Chat{
        Companion: hashname,
        Messages: []models.Message{
            models.Message{
                Name: hashname,
                Text: pack.Body.Data,
                Time: time,
            },
        },
    })

    var wsdata = struct {
        Comp struct {
            From string `json:"from"`
            To   string `json:"to"`
        } `json:"comp"`
        Text string `json:"text"`
        Time string `json:"time"`
    }{
        Comp: struct {
            From string `json:"from"`
            To   string `json:"to"`
        }{
            From: hashname,
            To:   user.Hashname,
        },
        Text: pack.Body.Data,
        Time: time,
    }

    if user.Session.Socket != nil {
        websocket.JSON.Send(user.Session.Socket, wsdata)
    }
    
    return
}

func handleFileServer(fs http.FileSystem) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        if _, err := fs.Open(r.URL.Path); os.IsNotExist(err) {
            return
        }
        http.FileServer(fs).ServeHTTP(w, r)
    })
}

func indexPage(w http.ResponseWriter, r *http.Request) {
    t, err := template.ParseFiles(
        settings.PATH_VIEWS+"index.html",
        settings.PATH_VIEWS+"home.html",
        settings.PATH_VIEWS+"about.html",
        settings.PATH_VIEWS+"login.html",
        settings.PATH_VIEWS+"signup.html",
        settings.PATH_VIEWS+"navbar.html",
        settings.PATH_VIEWS+"footer.html",
        settings.PATH_VIEWS+"account.html",
        settings.PATH_VIEWS+"network.html",
        settings.PATH_VIEWS+"settings.html",
        settings.PATH_VIEWS+"client.html",
        settings.PATH_VIEWS+"notfound.html",
    )
    if err != nil {
        panic("can't load hmtl files")
    }
    t.Execute(w, struct {
        WS   string
        HTTP string
        HOST string
    }{
        WS:   "wss://",
        HTTP: "https://",
        HOST: settings.CFG.Host.Http.Ipv4 + settings.CFG.Host.Http.Port,
    })
}
