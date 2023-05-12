package comment_test

// TODO refactor unit tests to use mocks and add other usecases
// TODO refactor test to add setup/cleanup tests
import (
	"testing"
	"time"

	"FaRyuk/internal/db/mongodb/comment"
	"FaRyuk/internal/types"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

var (
	db *comment.MongoCommentRepository
)

func TestMain(m *testing.M) {
	db = comment.NewMongoCommentRepository()
	comment := &types.Comment{
		ID:       "some-id",
		Content:  "some-search-text",
		IDResult: "some-result-id",
		Owner:    "owner-1",
	}
	db.InsertComment(comment)
	m.Run()
}

func TestInsertComment(t *testing.T) {
	comment := &types.Comment{
		ID:          uuid.NewString(),
		Content:     "new test comment",
		Owner:       "unit test",
		CreatedDate: time.Now(),
		UpdatedDate: time.Now(),
		IDResult:    uuid.NewString(),
	}
	err := db.InsertComment(comment)
	assert.NoError(t, err)
}

func TestGetComments(t *testing.T) {
	commentsChan := make(chan types.CommentsWithErrorType)
	go db.GetComments(commentsChan)
	result := <-commentsChan
	assert.NoError(t, result.Err)
	assert.NotEmpty(t, result.Comments)
}

func TestRemoveCommentByID(t *testing.T) {
	done := make(chan error)
	id := "some-id"
	go db.RemoveCommentByID(id, done)
	err := <-done
	assert.NoError(t, err)
}

func TestUpdateComment(t *testing.T) {
	commentToUpdate := &types.Comment{
		ID:          uuid.NewString(),
		Content:     "test comment to update",
		Owner:       "unit test",
		CreatedDate: time.Now(),
		UpdatedDate: time.Now(),
		IDResult:    uuid.NewString(),
	}
	err := db.InsertComment(commentToUpdate)
	assert.NoError(t, err)

	done := make(chan bool)
	comment := &types.Comment{
		ID:          uuid.NewString(),
		Content:     "updated comment",
		Owner:       "unit test",
		UpdatedDate: time.Now(),
		IDResult:    commentToUpdate.IDResult,
	}

	go db.UpdateComment(comment, done)
	result := <-done
	assert.True(t, result)
	// TODO get commlent by id and check if it is correctly updated
}

func TestGetCommentByID(t *testing.T) {

	result := make(chan types.CommentWithErrorType)
	id := "some-id"
	go db.GetCommentByID(id, result)
	ch := <-result
	assert.NotNil(t, ch.Comment)
	assert.NoError(t, ch.Err)
}

func TestGetCommentsByText(t *testing.T) {
	result := make(chan types.CommentsWithErrorType)
	search := "some-search-text"
	go db.GetCommentsByText(search, result)
	ch := <-result
	assert.NotNil(t, ch.Comments)
	assert.NoError(t, ch.Err)
}

func TestGetCommentsByTextAndOwner(t *testing.T) {
	search := "some-search-text"
	idUser := "owner-1"
	commentsChan := make(chan types.CommentsWithErrorType)
	go db.GetCommentsByTextAndOwner(search, idUser, commentsChan)
	result := <-commentsChan
	assert.NoError(t, result.Err)
	assert.NotEmpty(t, result.Comments)
}

func TestGetCommentsByResult(t *testing.T) {
	commentsChan := make(chan types.CommentsWithErrorType)
	idResult := "some-result-id"
	go db.GetCommentsByResultID(idResult, commentsChan)
	result := <-commentsChan
	assert.NoError(t, result.Err)
	assert.NotEmpty(t, result.Comments)
}
