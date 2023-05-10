package e2e_tests

import (
	"FaRyuk/api"
	"net/http"
	"net/http/httptest"
	"testing"
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

   // TODO override test dbHandler
	// TODO: Add more assertions to validate the response body or other aspects of the test
