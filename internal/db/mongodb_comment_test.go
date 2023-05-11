package db

// TODO refactor unit tests to use mocks and add other usecases
// TODO refactor test to add setup/cleanup tests
import (
	"testing"

	"FaRyuk/internal/types"

	"github.com/stretchr/testify/assert"
)

func init() {
	db := NewDBHandler()
	comment := &types.Comment{
		ID:       "some-id",
		Content:  "some-search-text",
		IDResult: "some-result-id",
		Owner:    "owner-1",
	}
	db.InsertComment(comment)
}

func TestInsertComment(t *testing.T) {
	db := NewDBHandler()
	comment := &types.Comment{}
	err := db.InsertComment(comment)
	assert.NoError(t, err)
}

func TestGetComments(t *testing.T) {
	db := NewDBHandler()

	commentsChan := make(chan types.CommentsWithErrorType)
	go db.GetComments(commentsChan)
	comments := <-commentsChan
	assert.NoError(t, comments.Err)
}

func TestRemoveCommentByID(t *testing.T) {
	db := NewDBHandler()
	done := make(chan error)
	id := "some-id"
	go db.RemoveCommentByID(id, done)
	err := <-done
	assert.NoError(t, err)
}

func TestUpdateComment(t *testing.T) {
	db := NewDBHandler()
	done := make(chan bool)
	comment := &types.Comment{}
	go db.UpdateComment(comment, done)
	result := <-done
	assert.True(t, result)
}

func TestGetCommentByID(t *testing.T) {
	db := NewDBHandler()
	result := make(chan types.CommentWithErrorType)
	id := "some-id"
	go db.GetCommentByID(id, result)
	ch := <-result
	assert.NotNil(t, ch.Comment)
	assert.NoError(t, ch.Err)

}

func TestGetCommentsByText(t *testing.T) {
	db := NewDBHandler()
	result := make(chan types.CommentsWithErrorType)

	search := "some-search-text"
	go db.GetCommentsByText(search, result)
	ch := <-result
	assert.NotNil(t, ch.Comments)
	assert.NoError(t, ch.Err)
}

func TestGetCommentsByTextAndOwner(t *testing.T) {
	db := NewDBHandler()

	search := "some-search-text"
	idUser := "some-user-id"
	result := make(chan types.CommentsWithErrorType)
	go db.GetCommentsByTextAndOwner(search, idUser, result)
	comments := <-result
	assert.NoError(t, comments.Err)
}

func TestGetCommentsByResult(t *testing.T) {
	db := NewDBHandler()

	commentsChan := make(chan types.CommentsWithErrorType)
	idResult := "some-result-id"
	go db.GetCommentsByResultID(idResult, commentsChan)
	result := <-commentsChan
	assert.NoError(t, result.Err)
	assert.NotEmpty(t, result.Comments)
}
