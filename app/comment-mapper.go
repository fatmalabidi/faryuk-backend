package app

import (
	"FaRyuk/services/comment_service"

	"github.com/gorilla/mux"
)

func AddCommentEndpoints(secure *mux.Router) {
	secure.HandleFunc("/api/get-comments/{id}", comment_service.GetComments).Methods("GET")
	secure.HandleFunc("/api/comments-highlights", comment_service.GetCommentsHighlight).Methods("GET")
	secure.HandleFunc("/api/count-highlights", comment_service.GetCountHighlight).Methods("GET")
	secure.HandleFunc("/api/get-tags", comment_service.GetTags).Methods("GET")
	secure.HandleFunc("/api/comment", comment_service.InsertComment).Methods("POST")
	secure.HandleFunc("/api/comment/{id}", comment_service.UpdateComment).Methods("POST")
	secure.HandleFunc("/api/comment/{id}", comment_service.DeleteComment).Methods("DELETE")
}
