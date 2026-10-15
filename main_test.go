package main

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestRoutes(t *testing.T) {
	app, err := newApp()
	if err != nil {
		t.Fatal(err)
	}

	server := httptest.NewServer(app.routes())
	defer server.Close()

	tests := []struct {
		path string
		want int
	}{
		{path: "/", want: http.StatusOK},
		{path: "/about", want: http.StatusOK},
		{path: "/experience", want: http.StatusOK},
		{path: "/missing", want: http.StatusNotFound},
	}

	for _, tt := range tests {
		t.Run(tt.path, func(t *testing.T) {
			resp, err := http.Get(server.URL + tt.path)
			if err != nil {
				t.Fatal(err)
			}
			defer resp.Body.Close()

			if resp.StatusCode != tt.want {
				t.Errorf(
					"GET %s: got %d, want %d",
					tt.path,
					resp.StatusCode,
					tt.want,
				)
			}
		})
	}
}
