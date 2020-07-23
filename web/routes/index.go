package routes

import (
	"net/http"
)

func (rtr *router) indexGetHandler(w http.ResponseWriter, r *http.Request) {
	err := rtr.tpl.ExecuteTemplate(w, "index.html", nil)
	if err != nil {
		log.Error("failed to render template: ", err)
	}
}

func (rtr *router) indexPostHandler(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	http.Redirect(w, r, "/", 302)
}
