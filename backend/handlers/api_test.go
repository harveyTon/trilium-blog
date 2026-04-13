package handlers

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/gin-gonic/gin"
)

func TestImageProxy_Security(t *testing.T) {
	gin.SetMode(gin.TestMode)
	h := &APIHandler{}

	tests := []struct {
		name       string
		query      string
		wantStatus int
	}{
		{
			name:       "empty url returns 400",
			query:      "",
			wantStatus: http.StatusBadRequest,
		},
		{
			name:       "invalid url returns 400",
			query:      "?url=no-scheme",
			wantStatus: http.StatusBadRequest,
		},
		{
			name:       "ftp scheme returns 400",
			query:      "?url=ftp://example.com/img.jpg",
			wantStatus: http.StatusBadRequest,
		},
		{
			name:       "localhost returns 403",
			query:      "?url=" + url.QueryEscape("http://localhost/img.jpg"),
			wantStatus: http.StatusForbidden,
		},
		{
			name:       "127.0.0.1 returns 403",
			query:      "?url=" + url.QueryEscape("http://127.0.0.1/img.jpg"),
			wantStatus: http.StatusForbidden,
		},
		{
			name:       "localhost with port returns 403",
			query:      "?url=" + url.QueryEscape("http://localhost:8080/img.jpg"),
			wantStatus: http.StatusForbidden,
		},
		{
			name:       "127.0.0.1 with port returns 403",
			query:      "?url=" + url.QueryEscape("http://127.0.0.1:8080/img.jpg"),
			wantStatus: http.StatusForbidden,
		},
		{
			name:       "private network 10.x returns 403",
			query:      "?url=" + url.QueryEscape("http://10.0.0.1/img.jpg"),
			wantStatus: http.StatusForbidden,
		},
		{
			name:       "private network 192.168.x returns 403",
			query:      "?url=" + url.QueryEscape("http://192.168.1.1/img.jpg"),
			wantStatus: http.StatusForbidden,
		},
		{
			name:       "private network 172.16.x returns 403",
			query:      "?url=" + url.QueryEscape("http://172.16.0.1/img.jpg"),
			wantStatus: http.StatusForbidden,
		},
		{
			name:       "link local 169.254.x returns 403",
			query:      "?url=" + url.QueryEscape("http://169.254.169.254/img.jpg"),
			wantStatus: http.StatusForbidden,
		},
		{
			name:       "ipv6 loopback returns 403",
			query:      "?url=" + url.QueryEscape("http://[::1]/img.jpg"),
			wantStatus: http.StatusForbidden,
		},
		{
			name:       "ipv6 unique local returns 403",
			query:      "?url=" + url.QueryEscape("http://[fd00::1]/img.jpg"),
			wantStatus: http.StatusForbidden,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)
			c.Request = httptest.NewRequest("GET", "/api/imageproxy"+tt.query, nil)

			h.ImageProxy(c)

			if w.Code != tt.wantStatus {
				t.Errorf("got status %d, want %d, body: %s", w.Code, tt.wantStatus, w.Body.String())
			}
		})
	}
}

func TestImageProxy_InvalidURL_Errors(t *testing.T) {
	gin.SetMode(gin.TestMode)
	h := &APIHandler{}

	invalidURLs := []string{
		"?url=" + url.QueryEscape("javascript:alert(1)"),
		"?url=" + url.QueryEscape("data:image/png;base64,abc"),
		"?url=/local/path",
		"?url=////invalid",
	}

	for _, rawURL := range invalidURLs {
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = httptest.NewRequest("GET", "/api/imageproxy"+rawURL, nil)

		h.ImageProxy(c)

		if w.Code == http.StatusOK {
			t.Errorf("URL %q should not return 200, got %d", rawURL, w.Code)
		}
	}
}

func TestIsBlockedHost(t *testing.T) {
	blocked := []string{
		"localhost",
		"127.0.0.1",
		"127.0.0.1:8080",
		"localhost:8080",
		"::1",
		"example.local",
		"10.0.0.1",
		"192.168.1.1",
		"172.16.0.1",
		"169.254.169.254",
		"100.64.0.1",
		"fd00::1",
		"fe80::1",
	}
	for _, h := range blocked {
		if !isBlockedHost(h) {
			t.Errorf("isBlockedHost(%q) = false, want true", h)
		}
	}

	allowed := []string{
		"example.com",
		"api.example.com",
		"185.199.108.153",
	}
	for _, h := range allowed {
		if isBlockedHost(h) {
			t.Errorf("isBlockedHost(%q) = true, want false", h)
		}
	}
}

func TestImageProxy_ErrorResponseFormat(t *testing.T) {
	gin.SetMode(gin.TestMode)
	h := &APIHandler{}

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest("GET", "/api/imageproxy", nil)

	h.ImageProxy(c)

	if w.Code != http.StatusBadRequest {
		t.Errorf("expected 400, got %d", w.Code)
	}

	var body map[string]interface{}
	if err := json.Unmarshal(w.Body.Bytes(), &body); err != nil {
		t.Fatalf("response is not valid JSON: %v", err)
	}

	if _, ok := body["error"]; !ok {
		t.Errorf("expected 'error' field in response body")
	}
}
