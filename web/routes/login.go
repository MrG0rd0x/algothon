package routes

import (
	"context"
	"net/http"
)

func (rtr *router) loginGetHandler(w http.ResponseWriter, r *http.Request) {
	session, err := rtr.cookie.Get(r, "session")
	if err != nil {
		log.Print("failed to get session: ", err)
		return
	}
	untyped, ok := session.Values["username"]
	if !ok || untyped == nil {
		rtr.tpl.ExecuteTemplate(w, "login.html", nil)
		return
	}
	username, ok := untyped.(string)
	if !ok {
		log.Print("failed parsing username: ", err)
		return
	}
	w.Write([]byte(username))
}

func (rtr *router) loginPostHandler(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	username := r.PostForm.Get("username")
	password := r.PostForm.Get("password")
	session, err := rtr.cookie.Get(r, "session")
	if err != nil {
		log.Print("Failed to create session: ", err)
		return
	}
	session.Values["username"] = username
	ctx := context.Background()
	defer ctx.Done()
	if !rtr.db.Login(ctx, username, password) {
		rtr.tpl.ExecuteTemplate(w, "login.html", "invalid credentials")
		return
	}
	err = session.Save(r, w)
	if err != nil {
		log.Error(err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	http.Redirect(w, r, "/", 302)
}
