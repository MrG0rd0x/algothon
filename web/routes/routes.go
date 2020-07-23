package routes

import (
	"net/http"
)

type routeHandler func(*router, http.ResponseWriter, *http.Request)

func (rtr *router) newHandlerFunc(rh routeHandler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		rh(rtr, w, r)
	}
}

// Route represents a single web route
type Route struct {
	Path    string
	Methods []string
	Handler routeHandler
	Auth    bool
}

var routes = []Route{
	{
		Path:    "/",
		Methods: []string{"GET"},
		Handler: (*router).indexGetHandler,
		Auth:    true,
	},
	{
		Path:    "/",
		Methods: []string{"POST"},
		Handler: (*router).indexPostHandler,
		Auth:    true,
	},
	{
		Path:    "/login",
		Methods: []string{"GET"},
		Handler: (*router).loginGetHandler,
	},
	{
		Path:    "/login",
		Methods: []string{"POST"},
		Handler: (*router).loginPostHandler,
	},
	{
		Path:    "/logout",
		Methods: []string{"GET"},
		Handler: (*router).logoutGetHandler,
	},
	{
		Path:    "/register",
		Methods: []string{"GET"},
		Handler: (*router).registerGetHandler,
	},
	{
		Path:    "/register",
		Methods: []string{"POST"},
		Handler: (*router).registerPostHandler,
	},
}

func (rtr *router) generateRoutes() {
	for _, route := range routes {
		log.Debugf("adding route '%s' (%s)", route.Path, route.Methods)
		if route.Auth {
			rtr.Handle(route.Path, rtr.requireAuth(route.Handler)).Methods(route.Methods...)
		}
		rtr.Handle(route.Path, rtr.newHandlerFunc(route.Handler)).Methods(route.Methods...)
	}
}
