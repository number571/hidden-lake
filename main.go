package main

import (
	"crypto/tls"
	"fmt"
	"github.com/number571/gopeer"
	"github.com/number571/hiddenlake/api"
	"github.com/number571/hiddenlake/db"
	"github.com/number571/hiddenlake/handle"
	"github.com/number571/hiddenlake/models"
	"github.com/number571/hiddenlake/settings"
	"github.com/number571/hiddenlake/utils"
	"github.com/number571/hiddenlake/ws"
	"golang.org/x/net/websocket"
	"html/template"
	"net/http"
	"os"
	"strings"
)

func init() {
	gopeer.Set(gopeer.SettingsType{
		"NETWORK": "[HIDDEN-LAKE]",
		"VERSION": "[1.0.2s]",
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

	mux.HandleFunc("/", indexPage)                  // GET
	mux.HandleFunc("/static/archive/", archivePage) // GET

	mux.HandleFunc("/api/login", api.Login)                      // POST
	mux.HandleFunc("/api/logout", api.Logout)                    // POST
	mux.HandleFunc("/api/signup", api.Signup)                    // POST
	mux.HandleFunc("/api/account", api.Account)                  // GET, POST, DELETE
	mux.HandleFunc("/api/account/connects", api.AccountConnects) // GET, PATCH
	mux.HandleFunc("/api/account/archive/", api.AccountArchive)  // GET, PUT, DELETE
	mux.HandleFunc("/api/network/chat/", api.NetworkChat)        // GET, POST, DELETE
	mux.HandleFunc("/api/network/client/", api.NetworkClient)    // GET, POST, DELETE
	//             "/api/network/client/:id/archive/"            // GET, POST
	//             "/api/network/client/:id/connects"            // GET, POST

	mux.Handle("/ws/network", websocket.Handler(ws.Network))

	handleServerTCP(&settings.CFG.Host.Tcp)
	handleServerHTTP(&settings.CFG.Host.Http, mux)
}

func handleServerTCP(model *models.Tcp) {
	settings.Listener = gopeer.NewListener(model.Ipv4 + model.Port)
	settings.Listener.Open().Run(handle.Actions)
}

func handleServerHTTP(model *models.Http, mux *http.ServeMux) {
	srv := &http.Server{
		Addr:    model.Ipv4 + model.Port,
		Handler: mux,
		TLSConfig: &tls.Config{
			MinVersion: tls.VersionTLS12,
			CurvePreferences: []tls.CurveID{
				tls.CurveP521,
				tls.CurveP384,
				tls.CurveP256,
			},
			PreferServerCipherSuites: true,
			CipherSuites: []uint16{
				tls.TLS_ECDHE_RSA_WITH_AES_256_GCM_SHA384,
				tls.TLS_ECDHE_RSA_WITH_AES_256_CBC_SHA,
				tls.TLS_RSA_WITH_AES_256_GCM_SHA384,
				tls.TLS_RSA_WITH_AES_256_CBC_SHA,
			},
		},
		TLSNextProto: make(map[string]func(*http.Server, *tls.Conn, http.Handler), 0),
	}
	srv.ListenAndServeTLS(model.Tls.Crt, model.Tls.Key)
}

func handleFileServer(fs http.FileSystem) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if _, err := fs.Open(r.URL.Path); os.IsNotExist(err) {
			indexPage(w, r)
			return
		}
		http.FileServer(fs).ServeHTTP(w, r)
	})
}

func archivePage(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		indexPage(w, r)
		return
	}

	filehash := strings.TrimPrefix(r.URL.Path, "/static/archive/")
	token := r.URL.Query().Get("token")

	if _, ok := settings.Users[token]; !ok {
		indexPage(w, r)
		return
	}

	err := settings.CheckLifetimeToken(token)
	if err != nil {
		indexPage(w, r)
		return
	} else {
		settings.Users[token].Session.Time = utils.CurrentTime()
	}

	user := settings.Users[token]
	file := db.GetFile(user, filehash)
	if file == nil {
		indexPage(w, r)
		return
	}

	w.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=\"%s\"", file.Name))
	http.ServeFile(w, r, settings.PATH_ARCHIVE+file.Path)
}

func indexPage(w http.ResponseWriter, r *http.Request) {
	t, err := template.ParseFiles(
		settings.PATH_VIEWS+"index.html",
		settings.PATH_VIEWS+"home.html",
		settings.PATH_VIEWS+"about.html",
		settings.PATH_VIEWS+"login.html",
		settings.PATH_VIEWS+"signup.html",
		settings.PATH_VIEWS+"account.html",
		settings.PATH_VIEWS+"network.html",
		settings.PATH_VIEWS+"chat.html",
		settings.PATH_VIEWS+"settings.html",
		settings.PATH_VIEWS+"client.html",
		settings.PATH_VIEWS+"clientarchive.html",
		settings.PATH_VIEWS+"clientarchivefile.html",
		settings.PATH_VIEWS+"clientconnects.html",
		settings.PATH_VIEWS+"clients.html",
		settings.PATH_VIEWS+"archive.html",
		settings.PATH_VIEWS+"archivefile.html",
		settings.PATH_VIEWS+"notfound.html",
		settings.PATH_VIEWS+"message_part.html",
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
