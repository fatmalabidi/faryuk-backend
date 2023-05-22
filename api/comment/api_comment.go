package comment

import (
	"FaRyuk/api/utils"
	"FaRyuk/config"
	"FaRyuk/database"
	"FaRyuk/internal/types"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

// TODO move to new file: API mapper
func AddCommentEndpoints(secure *mux.Router) {
	secure.HandleFunc("/api/v1/comments", ListComments).Methods(http.MethodGet)
	secure.HandleFunc("/api/v1/comments/{id}", RemoveCommentByID).Methods(http.MethodDelete)
}

type CommentFilter struct {
	ResultID   string `schema:"resultID"`
	SearchText string `schema:"searchText"`
	Owner      string `schema:"owner"`
}

func ListComments(w http.ResponseWriter, rr *http.Request) {
	cfg, confErr := config.MakeConfig()
	if confErr != nil {
		log.Fatal(confErr)
	}

	dbHandler, err := database.CreateDbHandler(cfg)
	if err != nil {
		log.Fatal(err)
	}
	defer dbHandler.CloseConnection()
	resultID := rr.URL.Query().Get("resultID")
	searchText := rr.URL.Query().Get("searchText")

	// FIXME : should use the same driver methor (it's uo to the driver to checkout all possible filters
	// The driver will take a filter object as paramater that should be validated in the service layer and appliyed in the db
	commentsChan := make(chan types.CommentsWithErrorType)
	if resultID != "" {
		go dbHandler.GetCommentsByResultID(resultID, commentsChan)
	} else if searchText != "" {
		owner := rr.URL.Query().Get("owner")
		if owner != "" {
			go dbHandler.GetCommentsByTextAndOwner(searchText, owner, commentsChan)
		} else {
			go dbHandler.GetCommentsByText(searchText, commentsChan)
		}
	} else {
		go dbHandler.GetComments(commentsChan)
	}

	result := <-commentsChan
	if result.Err != nil {
		utils.WriteInternalError(&w, result.Err.Error())
		return
	}

	utils.ReturnSuccess(&w, result.Comments)
}

func RemoveCommentByID(w http.ResponseWriter, r *http.Request) {
	// Get the comment ID from the request URL parameters
	vars := mux.Vars(r)
	commentID := vars["id"]

	cfg, confErr := config.MakeConfig()
	if confErr != nil {
		log.Fatal(confErr)
	}
	dbHandler, err := database.CreateDbHandler(cfg)
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
		utils.WriteInternalError(&w, err.Error())
		return
	}

	// If the operation was successful, return a success response
	utils.ReturnSuccess(&w, nil)
}
