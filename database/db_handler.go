package database

import (
	"FaRyuk/config"
	"FaRyuk/database/mongodb/comment"
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
	InsertComment(r *types.Comment, done chan<- error)
	GetComments(comments chan<- types.CommentsWithErrorType)
	RemoveCommentByID(id string, done chan<- error)
	UpdateComment(r *types.Comment, done chan<- error)
	GetCommentByID(id string, result chan<- types.CommentWithErrorType)
	GetCommentsByText(search string, result chan types.CommentsWithErrorType)
	GetCommentsByTextAndOwner(search string, idUser string, result chan types.CommentsWithErrorType)
	GetCommentsByResultID(idResult string, result chan<- types.CommentsWithErrorType)
	CloseConnection()
}

func CreateDbHandler(cfg *config.Config) (DbHandler, error) {
	var dbHandler MainDb
	switch cfg.Database.DbType {
	case "mongo":
		dbHandler.CommentHandler = comment.NewMongoCommentRepository(cfg)
	default:
		return nil, errors.New("db type not supported")
	}
	return dbHandler, nil
}
