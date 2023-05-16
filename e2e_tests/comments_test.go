package e2e_tests

import (
	"FaRyuk/api"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

// TODO override test dbHandler (use db test)
// TODO: Add more assertions to validate the response body or other aspects of the test
func TestMain(m *testing.M) {
  os.Setenv("CONFIGOR_ENV", "test")
	os.Exit(m.Run())
}

func TestListComments(t *testing.T) {
	// TODO add token validation
	// Create a mock request
	req, err := http.NewRequest("GET", "/comments", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	api.ListComments(rr, req)
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}
}

func TestRemoveCommentByID(t *testing.T) {

	commentID := "some-id"

	req, err := http.NewRequest("DELETE", fmt.Sprintf("/comments:%s", commentID), nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	api.GetCommentsByResultID(rr, req)
	assert.Equal(t, http.StatusOK, rr.Code)
}

func TestGetCommentsByResultID(t *testing.T) {
	resultID := "some-result-id"

	req, err := http.NewRequest("GET", fmt.Sprintf("/comments:%s", resultID), nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder() 
	api.GetCommentsByResultID(rr, req)
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	assert.Equal(t, http.StatusOK, rr.Code)
	assert.NotNil(t, rr.Body)
	assert.Empty(t, rr.Body)
}
