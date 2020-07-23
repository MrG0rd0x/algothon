package routes

import (
	"net/http"
	"strings"
)

func blockIndex(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if strings.HasSuffix(r.URL.Path, "/") {
			http.NotFound(w, r)
			return
		}
		h.ServeHTTP(w, r)
	})
}

func (rtr *router) requireAuth(h routeHandler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		session, err := rtr.cookie.Get(r, "session")
		if !checkError("Cannot decrypt session cookie", err, w, r) {
			return
		}
		raw, ok := session.Values["username"]
		if !ok {
			log.Debug("User has no session, redirecting to login")
			http.Redirect(w, r, "/login", http.StatusFound)
			return
		}
		username, ok := raw.(string)
		if !checkOk("Type assertion on session username failed", ok, w, r) {
			return
		}
		log.Debugf("User logged in as '%s'", username)
		h(rtr, w, r)
	}
}

func checkOk(str string, ok bool, w http.ResponseWriter, r *http.Request) bool {
	if ok {
		return true
	}
	log.Error(str)
	http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	return false
}

func checkError(str string, err error, w http.ResponseWriter, r *http.Request) bool {
	if err == nil {
		return true
	}
	log.Errorf("%s: %s", str, err)
	http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	return false
}
