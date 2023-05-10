package e2e_tests

import (
	"FaRyuk/api"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
)

func TestListComments(t *testing.T) {
	// TODO add token validation
	// Create a mock request
	req, err := http.NewRequest("GET", "/comments", nil)
	if err != nil {
		t.Fatal(err)
	}

	// Create a ResponseRecorder to record the response
	rr := httptest.NewRecorder()

	// Call the handler function directly
	api.ListComments(rr, req)

	// Check the response status code
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}
}


func TestRemoveCommentByID(t *testing.T) {
	// Create a new router
	router := mux.NewRouter()
	router.HandleFunc("/api/v1/comments/{id}", api.RemoveCommentByID).Methods(http.MethodDelete)

	// Create a test request with a comment ID
	commentID := "1"
	req, err := http.NewRequest(http.MethodDelete,fmt.Sprintf( "/api/v1/comments/%s",commentID), nil)
	assert.NoError(t, err)

	// Create a new response recorder to capture the response
	rr := httptest.NewRecorder()

	// Serve the request using the router
	router.ServeHTTP(rr, req)

	// Check the response status code
	assert.Equal(t, http.StatusOK, rr.Code)
}

  // TODO override test dbHandler
	// TODO: Add more assertions to validate the response body or other aspects of the test
