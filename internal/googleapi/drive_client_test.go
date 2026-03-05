package googleapi

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/Alechan/gdrivectl/internal/fail"
)

func TestDriveSearchSuccessFixture(t *testing.T) {
	payload := readFixture(t, "drive/search_success.json")
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/drive/v3/files" {
			t.Fatalf("path = %q, want /drive/v3/files", r.URL.Path)
		}
		_, _ = w.Write(payload)
	}))
	defer srv.Close()

	c := NewDriveClient(testHTTPClient(srv.URL))
	got, err := c.Search(context.Background(), "tok", SearchRequest{
		Query: "name contains 'RFC'",
	})
	if err != nil {
		t.Fatalf("Search() error = %v", err)
	}
	files, ok := got["files"].([]any)
	if !ok || len(files) != 1 {
		t.Fatalf("files = %#v, want single element", got["files"])
	}
}

func TestDriveFileMetaSuccessFixture(t *testing.T) {
	payload := readFixture(t, "drive/file_meta_success.json")
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/drive/v3/files/file-123" {
			t.Fatalf("path = %q, want /drive/v3/files/file-123", r.URL.Path)
		}
		_, _ = w.Write(payload)
	}))
	defer srv.Close()

	c := NewDriveClient(testHTTPClient(srv.URL))
	got, err := c.FileMeta(context.Background(), "tok", FileMetaRequest{ID: "file-123"})
	if err != nil {
		t.Fatalf("FileMeta() error = %v", err)
	}
	if got["id"] != "file-123" {
		t.Fatalf("id = %#v, want file-123", got["id"])
	}
}

func TestDriveExportSuccess(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/drive/v3/files/doc-1/export" {
			t.Fatalf("path = %q, want /drive/v3/files/doc-1/export", r.URL.Path)
		}
		_, _ = w.Write([]byte("exported body"))
	}))
	defer srv.Close()

	c := NewDriveClient(testHTTPClient(srv.URL))
	got, err := c.ExportDoc(context.Background(), "tok", ExportRequest{ID: "doc-1", MIME: "text/plain"})
	if err != nil {
		t.Fatalf("ExportDoc() error = %v", err)
	}
	if string(got) != "exported body" {
		t.Fatalf("ExportDoc() body = %q, want %q", string(got), "exported body")
	}
}

func TestDriveStatusErrorMapping(t *testing.T) {
	tests := []struct {
		name       string
		statusCode int
		fixture    string
		wantCat    string
	}{
		{name: "unauthorized auth", statusCode: http.StatusUnauthorized, fixture: "drive/error_401.json", wantCat: "auth"},
		{name: "forbidden scope", statusCode: http.StatusForbidden, fixture: "drive/error_403.json", wantCat: "scope"},
		{name: "server api", statusCode: http.StatusInternalServerError, fixture: "drive/error_500.json", wantCat: "api"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			payload := readFixture(t, tt.fixture)
			srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(tt.statusCode)
				_, _ = w.Write(payload)
			}))
			defer srv.Close()

			c := NewDriveClient(testHTTPClient(srv.URL))
			_, err := c.Search(context.Background(), "tok", SearchRequest{Query: "x"})
			assertFailCategoryStrict(t, err, tt.wantCat)
		})
	}
}

func TestDriveSearchMalformedJSONMapsToAPI(t *testing.T) {
	payload := readFixture(t, "drive/malformed.json")
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write(payload)
	}))
	defer srv.Close()

	c := NewDriveClient(testHTTPClient(srv.URL))
	_, err := c.Search(context.Background(), "tok", SearchRequest{Query: "x"})
	assertFailCategoryStrict(t, err, "api")
}

func TestDriveNetworkErrorMapsToNetwork(t *testing.T) {
	c := NewDriveClient(testHTTPClientWithErr(errors.New("dial tcp: no such host")))
	_, err := c.Search(context.Background(), "tok", SearchRequest{Query: "x"})
	assertFailCategoryStrict(t, err, "network")
}

func assertFailCategoryStrict(t *testing.T, err error, want string) {
	t.Helper()
	var fe *fail.Error
	if !errors.As(err, &fe) {
		t.Fatalf("error type = %T, want *fail.Error (err=%v)", err, err)
	}
	if fe.Category != want {
		t.Fatalf("category = %q, want %q", fe.Category, want)
	}
}
