package internal

import (
	"FaRyuk/api"
	"FaRyuk/app"
	"github.com/gorilla/mux"
)

func MainServer() {
	myRouter := mux.NewRouter().StrictSlash(true)

	secure := myRouter.PathPrefix("/").Subrouter()
	api.HandleRequests()
	// TODO add mappers
	app.AddCommentEndpoints(secure)
}
