package main

import "github.com/gorilla/mux"
import "net/http"

type Route struct {
	Name       string
	Method     string
	Pattern    string
	HandleFunc http.HandlerFunc
}

type Routes []Route

func newRouter() *mux.Router {
	router := mux.NewRouter().StrictSlash(true)

	for _, route := range routes {
		router.Methods(route.Method).
			Path(route.Pattern).
			Name(route.Name).
			Handler(route.HandleFunc)

	}

	return router
}

var routes = Routes{
	Route{
		"NovedadList",
		"GET",
		"/novedades",
		NovedadList,
	},
	Route{
		"NovedadShow",
		"GET",
		"/novedades/{id}",
		NovedadShow,
	},
	Route{
		"NovedadAdd",
		"POST",
		"/novedades",
		NovedadAdd,
	},
	Route{
		"NovedadUpdate",
		"PUT",
		"/novedades/{id}",
		NovedadUpdate,
	},
	Route{
		"NovedadRemove",
		"DELETE",
		"/novedades/{id}",
		NovedadRemove,
	},
}