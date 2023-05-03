package api

import (
	"FaRyuk/internal/db"
	"FaRyuk/internal/types"
	"net/http"

	"github.com/gorilla/mux"
)

// TODO move to new file: API mapper
func AddCommentEndpoints(secure *mux.Router) {
	secure.HandleFunc("/api/v1/comments", ListComments).Methods("GET")
}

func ListComments(w http.ResponseWriter, r *http.Request) {	
	dbHandler := db.NewDBHandler()
	defer dbHandler.CloseConnection()

	commentsChan := make(chan types.CommentsWithErrorType)
	go dbHandler.GetComments(commentsChan)

	result := <-commentsChan

	if result.Err != nil {
		writeInternalError(&w, result.Err.Error())
		return
	}

	returnSuccess(&w, result.Comments)
}
