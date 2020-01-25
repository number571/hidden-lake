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
	http.Handle("/static/", http.StripPrefix(
		"/static/",
		handleFileServer(http.Dir(settings.PATH_STATIC))),
	)

	http.HandleFunc("/", indexPage) // GET
	http.HandleFunc("/api/login", api.Login) // POST
	http.HandleFunc("/api/logout", api.Logout) // POST
	http.HandleFunc("/api/signup", api.Signup) // POST
	http.HandleFunc("/api/account", api.Account) // GET, POST, DELETE
	http.HandleFunc("/api/network/", api.Network) // GET, POST, DELETE
	http.HandleFunc("/api/network/client/", api.Client) // GET, POST, DELETE

	http.Handle("/ws/network", websocket.Handler(ws.Network))

	handleServerTCP(&settings.CFG.Host.Tcp)
	handleServerHTTP(&settings.CFG.Host.Http)
}

func handleServerTCP(model *models.Tcp) {
	settings.Listener = gopeer.NewListener(model.Ipv4 + model.Port)
	settings.Listener.Open().Run(handleActions)
}

func handleServerHTTP(model *models.Http) {
	err := http.ListenAndServeTLS(
		model.Ipv4+model.Port,
		model.Tls.Crt,
		model.Tls.Key,
		nil,
	)
	if err != nil {
		http.ListenAndServe(model.Ipv4+model.Port, nil)
	}
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
	var ws = "ws://"
	var http = "http://"
	if settings.CFG.Host.Http.Tls.Crt != "" && settings.CFG.Host.Http.Tls.Key != "" {
		ws = "wss://"
		http = "https://"
	}
	t.Execute(w, struct {
		WS   string
		HTTP string
		HOST string
	}{
		WS:   ws,
		HTTP: http,
		HOST: settings.CFG.Host.Http.Ipv4 + settings.CFG.Host.Http.Port,
	})
}
