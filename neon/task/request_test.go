package task

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	t "time"

	_build "github.com/c4s4/neon/neon/build"
)

func TestRequestTimeout(test *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		select {
		case <-r.Context().Done():
		case <-t.After(3 * t.Second):
		}
	}))
	defer server.Close()
	context := _build.NewContext(&_build.Build{})
	start := t.Now()
	err := request(context, requestArgs{Request: server.URL, Timeout: 1})
	if err == nil || !strings.Contains(err.Error(), "Timeout") {
		test.Errorf("request should time out, got: %v", err)
	}
	if t.Since(start) > 2*t.Second {
		test.Errorf("request timeout was not applied (took %v)", t.Since(start))
	}
}
