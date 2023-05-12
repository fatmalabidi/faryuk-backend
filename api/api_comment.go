package api

import (
	"FaRyuk/internal/db"
	"FaRyuk/internal/types"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

// TODO move to new file: API mapper
func AddCommentEndpoints(secure *mux.Router) {
	secure.HandleFunc("/api/v1/comments", ListComments).Methods(http.MethodGet)
	secure.HandleFunc("/api/v1/comments/{id}", RemoveCommentByID).Methods(http.MethodDelete)
	secure.HandleFunc("/api/v1/comments/{id}", GetCommentsByResultID).Methods(http.MethodGet)
}

func ListComments(w http.ResponseWriter, _ *http.Request) {
	dbHandler, err := db.CreateDbHandler(db.Config{CommentDbType: "mongo"})
	if err != nil {
		log.Fatal(err)
	}
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

	dbHandler, err := db.CreateDbHandler(db.Config{CommentDbType: "mongo"})
	if err != nil {
		log.Fatal(err)
	}
	defer dbHandler.CloseConnection()

	// Create a channel to receive the result of the database operation
	done := make(chan error)

	// Call the RemoveCommentByID method asynchronously
	go dbHandler.RemoveCommentByID(commentID, done)

	// Wait for the database operation to complete and check for errors
	err = <-done

	if err != nil {
		// If an error occurred, return an appropriate HTTP response
		writeInternalError(&w, err.Error())
		return
	}

	// If the operation was successful, return a success response
	returnSuccess(&w, nil)
}

// GetCommentByResultIDget the comments list for a result
func GetCommentsByResultID(w http.ResponseWriter, r *http.Request) {
	// Get the result ID from the request URL parameters
	vars := mux.Vars(r)
	resultID := vars["id"]
	dbHandler, err := db.CreateDbHandler(db.Config{CommentDbType: "mongo"})
	if err != nil {
		log.Fatal(err)
	}
	defer dbHandler.CloseConnection()

	// Create a channel to receive the result of the database operation
	commentsChan := make(chan types.CommentsWithErrorType)

	// Call the GetCommentsByResult method asynchronously
	go dbHandler.GetCommentsByResultID(resultID, commentsChan)

	// Wait for the database operation to complete and check for errors
	result := <-commentsChan
	if result.Err != nil {
		// If an error occurred, return an appropriate HTTP response
		writeInternalError(&w, result.Err.Error())
		return
	}

	// If the operation was successful, return a success response with the list of comments
	returnSuccess(&w, result.Comments)
}
