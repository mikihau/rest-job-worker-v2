package main

import (
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/mikihau/rest-job-worker-v2/handlers"
)

func getStuff(w http.ResponseWriter, r *http.Request, l *log.Logger) {
	w.WriteHeader(http.StatusOK)
}

func TestAuth(t *testing.T) {

	cases := []struct {
		headerName, headerValue string
		code                    int
	}{
		{"", "", http.StatusUnauthorized},
		{"Authorization", "wow", http.StatusForbidden},
		{"Authorization", "123456", http.StatusOK},
	}

	req, err := http.NewRequest("GET", "/stuff", nil)
	if err != nil {
		t.Fatal(err)
	}

	// this test handler requires Role ReadWrite
	authorizedFunc := handlers.VerifyAuth(getStuff, log.New(os.Stdout, "", log.LstdFlags))
	handler := http.HandlerFunc(authorizedFunc)

	for _, c := range cases {
		rr := httptest.NewRecorder()
		req.Header.Set(c.headerName, c.headerValue)
		handler.ServeHTTP(rr, req)

		if status := rr.Code; status != c.code {
			t.Errorf("Wrong status code: with header %v:%v, expecting %v, but got %v",
				c.headerName, c.headerValue, c.code, status)
		}
	}
}
