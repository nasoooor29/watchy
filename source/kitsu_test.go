package source

import (
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
)

func TestFetchDetails(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var filepath string
		switch r.URL.Path {
		case "/anime":
			if r.URL.String() == "/anime?filter%5Btext%5D=Cowboy+Bebop" {
				filepath = "testdata/kitsu/anime/1/detail.json"
				break
			}
			fallthrough
		case "/anime/1/genres":
			filepath = "testdata/kitsu/anime/1/genres.json"
		default:
			t.Fatalf("unexpected route hit: %s", r.URL.Path)
		}
		payload, err := os.ReadFile(filepath)
		if err != nil {
			t.Fatalf("failed reading mock data: %v", err)
		}
		testBaseURL := "http://" + r.Host
		patchedPayload := strings.ReplaceAll(string(payload), "https://kitsu.io/api/edge", testBaseURL)

		w.Header().Set("Content-Type", "application/vnd.api+json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(patchedPayload))
	}))
	defer server.Close()

	src := NewKitsuSource(server.URL, server.Client())

	res, err := src.FetchDetails("Cowboy Bebop")
	if err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}
	if res == nil {
		t.Fatal("expected a ShowDetail struct, got nil")
	}
}
