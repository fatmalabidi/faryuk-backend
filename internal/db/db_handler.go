package db

import (
	"FaRyuk/internal/db/mongodb/comment"
	"FaRyuk/internal/types"
	"errors"
)

type DbHandler interface {
	CommentHandler
	// TODO add the remaining interfaces 
}

type MainDb struct {
	CommentHandler
	// TODO add the remaining handlers
}

type CommentHandler interface {
	InsertComment(r *types.Comment) error
	GetComments(comments chan<- types.CommentsWithErrorType)
	RemoveCommentByID(id string, done chan<- error)
	UpdateComment(r *types.Comment, done chan<- bool)
	GetCommentByID(id string, result chan<- types.CommentWithErrorType)
	GetCommentsByText(search string, result chan types.CommentsWithErrorType)
	GetCommentsByTextAndOwner(search string, idUser string, result chan types.CommentsWithErrorType)
	GetCommentsByResultID(idResult string, result chan<- types.CommentsWithErrorType)
	CloseConnection()
}

func CreateDbHandler(config Config) (DbHandler, error) {
	var dbHandler MainDb
	switch config.CommentDbType {
	case "mongo":
		dbHandler.CommentHandler= comment.NewMongoCommentRepository()
	default:
		return nil, errors.New("db type not supported")
	}
	return dbHandler, nil
}


type Config struct{
	CommentDbType string
	// TODO add the remaining model's db type
	// each model can be implemented with different driver/db (MySql, Mongo, DynamoDb, redis...)
}