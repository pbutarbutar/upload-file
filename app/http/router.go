package router

import (
	"brks/app/handler"
	"net/http"

	"github.com/gorilla/mux"
)

type Route struct {
	Name        string
	Method      string
	Pattern     string
	HandlerFunc http.HandlerFunc
	IsExcluded  bool
}

var routes = []Route{}

func GetRouters(hup handler.UploadInt) *mux.Router {
	uploadimage := []Route{
		{
			"Upload Images",
			"POST",
			"/v1/upload",
			hup.UploadFile,
			false,
		},
	}
	routes = append(routes, uploadimage...)

	router := mux.NewRouter().StrictSlash(true)

	for _, route := range routes {
		router.
			Methods(route.Method).
			Path(route.Pattern).
			Name(route.Name).
			Handler(route.HandlerFunc)
	}
	return router
}
