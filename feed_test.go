package feed

import (
	"net/http/httptest"
	"strings"
	"testing"
)

func TestFeed(t *testing.T) {
	payload := strings.NewReader("")
	req := httptest.NewRequest("GET", "/?view=json", payload)
	rr := httptest.NewRecorder()
	RSSFeed(rr, req)
	t.Log(rr.Body.String())
}
