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
		"Healthy",
		"GET",
		"/api/novedad/healthy",
		Healthy,
	},
	Route{
		"NovedadList",
		"GET",
		"/api/novedad/novedades",
		NovedadList,
	},
	Route{
		"NovedadShow",
		"GET",
		"/api/novedad/novedades/{id}",
		NovedadShow,
	},
	Route{
		"NovedadAdd",
		"POST",
		"/api/novedad/novedades",
		NovedadAdd,
	},
	Route{
		"NovedadUpdate",
		"PUT",
		"/api/novedad/novedades/{id}",
		NovedadUpdate,
	},
	Route{
		"NovedadRemove",
		"DELETE",
		"/api/novedad/novedades/{id}",
		NovedadRemove,
	},
	Route{
		"NovedadesRemoveMasivo",
		"DELETE",
		"/api/novedad/novedades",
		NovedadesRemoveMasivo,
	},
}
