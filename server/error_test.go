package server

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

// TestErrorResponse will test the method ErrorResponse()
func TestErrorResponse(t *testing.T) {
	t.Run("placeholder test", func(_ *testing.T) {
		w := httptest.NewRecorder()
		ErrorResponse(w, &http.Request{}, ErrorMethodNotFound, "test message", http.StatusBadRequest)

		// todo: actually test the error response
	})
}
