package comment

import (
	"FaRyuk/api/utils"
	"FaRyuk/config"
	"FaRyuk/internal/types"
)

type CommentInMemoryRepository struct {
	config *config.AppConfig
}

func NewInMemoryRepository(config *config.AppConfig) *CommentInMemoryRepository {
	return &CommentInMemoryRepository{}
}

func (repo *CommentInMemoryRepository) Create(comment *types.Comment, done chan<- error) {

}

func (repo *CommentInMemoryRepository) Delete(id string, done chan<- error) {

}

func (repo *CommentInMemoryRepository) Update(r *types.Comment, done chan<- error) {

}

func (repo *CommentInMemoryRepository) GetByID(id string, result chan<- *types.CommentWithErrorType) {

}

func (repo *CommentInMemoryRepository) List(listCommentsFilter utils.ListCommentsFilter, commentsChan chan<- *types.CommentsWithErrorType) {

}

// CloseConnection : closes connection with mongo db
func (repo *CommentInMemoryRepository) CloseCommentDBConnection() {

}

var data []*types.Comment = []*types.Comment{}
