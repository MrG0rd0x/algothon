package routes

import (
	"net/http"
)

func (rtr *router) logoutGetHandler(w http.ResponseWriter, r *http.Request) {
	session, err := rtr.cookie.Get(r, "session")
	if err != nil {
		log.Print("cannot decrypt session cookie: ", err)
		http.Redirect(w, r, "/", 302)
		return
	}
	session.Options.MaxAge = -1
	err = rtr.cookie.Save(r, w, session)
	if err != nil {
		log.Print("error saving session: ", err)
		return
	}
	http.Redirect(w, r, "/", http.StatusFound)
}
