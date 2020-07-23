package routes

import (
	"html/template"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/gorilla/sessions"

	"github.com/double-nibble/algothon/userdb"
)

// Config supplies needed information about the router
type Config struct {
	DB         userdb.UserDB
	SessionKey []byte
	Salt       string
}

type router struct {
	*mux.Router
	db     userdb.UserDB
	cookie *sessions.CookieStore
	tpl    *template.Template
	salt   string
}

// NewRouter returns a new http.Handler for serving
func NewRouter(config *Config) http.Handler {
	rtr := &router{
		Router: mux.NewRouter(),
		db:     config.DB,
		cookie: sessions.NewCookieStore(config.SessionKey),
		tpl:    template.Must(template.ParseGlob("templates/*.html")),
		salt:   config.Salt,
	}
	rtr.generateRoutes()
	fs := http.FileServer(http.Dir("./static/"))
	rtr.PathPrefix("/static/").Handler(http.StripPrefix("/static/", blockIndex(fs)))
	rtr.Use(loggingWrapper)
	return rtr
}
