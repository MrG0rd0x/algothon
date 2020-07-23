package routes

import (
	"context"
	"net/http"
)

func (rtr *router) registerGetHandler(w http.ResponseWriter, r *http.Request) {
	rtr.tpl.ExecuteTemplate(w, "register.html", nil)
}

func (rtr *router) registerPostHandler(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	username := r.PostForm.Get("username")
	password := r.PostForm.Get("password")
	ctx := context.Background()
	defer ctx.Done()
	if err := rtr.db.Register(ctx, username, password); err != nil {
		log.Error("Could not create new user: ", err)
	}
	http.Redirect(w, r, "/login", 302)
}
