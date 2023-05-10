package api

import (
	"FaRyuk/internal/db"
	"FaRyuk/internal/types"
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
)

// TODO move to new file: API mapper
func AddCommentEndpoints(secure *mux.Router) {
	secure.HandleFunc("/api/v1/comments", ListComments).Methods(http.MethodGet)
	secure.HandleFunc("/api/v1/comments/{id}", RemoveCommentByID).Methods(http.MethodDelete)
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


func RemoveCommentByID(w http.ResponseWriter, r *http.Request) {
	// Get the comment ID from the request URL parameters
	vars := mux.Vars(r)
	commentID := vars["id"]
  fmt.Println("id to delete", commentID)
	dbHandler := db.NewDBHandler()
	defer dbHandler.CloseConnection()

	// Create a channel to receive the result of the database operation
	done := make(chan error)

	// Call the RemoveCommentByID method asynchronously
	go dbHandler.RemoveCommentByID(commentID, done)

	// Wait for the database operation to complete and check for errors
	err := <-done
	if err != nil {
		// If an error occurred, return an appropriate HTTP response
		writeInternalError(&w, err.Error())
		return
	}

	// If the operation was successful, return a success response
	returnSuccess(&w, nil)
}
