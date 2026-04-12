package blog

import (
	"strings"
	"testing"
)

func TestProcessContent_ImageProxy(t *testing.T) {
	base := "http://imgproxy.example.com"

	svc := &Service{
		imageProxyEnabled: true,
		imageProxyBaseUrl: base,
		domain:            "https://blog.example.com",
	}

	svcDisabled := &Service{
		imageProxyEnabled: false,
		imageProxyBaseUrl: base,
		domain:            "https://blog.example.com",
	}

	svcInternal := &Service{
		imageProxyEnabled: true,
		imageProxyBaseUrl: "",
		domain:            "https://blog.example.com",
	}

	tests := []struct {
		name      string
		svc       *Service
		html      string
		checkURL  string
		expectURL string
	}{
		{
			name:      "Rule1 attachment goes to /api/assets/",
			svc:       svc,
			html:      `<img src="/attachments/abc123.png"/>`,
			checkURL:  "/api/assets/abc123.png",
			expectURL: "/api/assets/abc123.png",
		},
		{
			name:      "Rule2 external image proxied when enabled",
			svc:       svc,
			html:      `<img src="https://external.com/photo.jpg"/>`,
			checkURL:  "img src=\"",
			expectURL: base + "?url=",
		},
		{
			name:      "Rule2 external http image proxied when enabled",
			svc:       svc,
			html:      `<img src="http://insecure.com/img.png"/>`,
			checkURL:  "img src=\"",
			expectURL: base + "?url=",
		},
		{
			name:      "Rule2 external image unchanged when disabled",
			svc:       svcDisabled,
			html:      `<img src="https://external.com/photo.jpg"/>`,
			checkURL:  `src="https://external.com/photo.jpg"`,
			expectURL: `src="https://external.com/photo.jpg"`,
		},
		{
			name:      "Rule3 no double rewrite for already proxied URL",
			svc:       svc,
			html:      `<img src="` + base + `?url=https%3A%2F%2Fother.com%2Fimg.jpg"/>`,
			checkURL:  base + `?url=`,
			expectURL: base + `?url=`,
		},
		{
			name:      "Rule4 relative URL unchanged",
			svc:       svc,
			html:      `<img src="/local/img.png"/>`,
			checkURL:  `src="/local/img.png"`,
			expectURL: `src="/local/img.png"`,
		},
		{
			name:      "Rule4 internal /api/assets/ URL unchanged when proxy enabled",
			svc:       svc,
			html:      `<img src="https://blog.example.com/api/assets/xyz.png"/>`,
			checkURL:  `src="https://blog.example.com/api/assets/xyz.png"`,
			expectURL: `src="https://blog.example.com/api/assets/xyz.png"`,
		},
		{
			name:      "Rule4 proxy base URL unchanged to avoid loop",
			svc:       svc,
			html:      `<img src="` + base + `/some/path.jpg"/>`,
			checkURL:  `src="` + base + `/some/path.jpg"`,
			expectURL: `src="` + base + `/some/path.jpg"`,
		},
		{
			name:      "Rule4 external CDN with /assets/ path still proxied (not treated as internal)",
			svc:       svc,
			html:      `<img src="https://cdn.example.com/assets/foo.png"/>`,
			checkURL:  base + "?url=",
			expectURL: base + "?url=",
		},
		{
			name:      "Rule4 external CDN with /api/assets/ path still proxied",
			svc:       svc,
			html:      `<img src="https://cdn.example.com/api/assets/123"/>`,
			checkURL:  base + "?url=",
			expectURL: base + "?url=",
		},
		{
			name:      "Rule4 external host with /api/imageproxy path still proxied",
			svc:       svc,
			html:      `<img src="https://cdn.example.com/api/imageproxy?url=https://evil.com/img.jpg"/>`,
			checkURL:  base + "?url=",
			expectURL: base + "?url=",
		},
		{
			name:      "Rule4 relative /assets/ path is not proxied",
			svc:       svc,
			html:      `<img src="/assets/logo.png"/>`,
			checkURL:  `src="/assets/logo.png"`,
			expectURL: `src="/assets/logo.png"`,
		},
		{
			name:      "Rule2 external with query string proxied correctly",
			svc:       svc,
			html:      `<img src="https://example.com/img.jpg?size=large"/>`,
			checkURL:  base + "?url=",
			expectURL: base + "?url=",
		},
		{
			name:      "internal fallback: external image uses /api/imageproxy when enabled but no baseUrl",
			svc:       svcInternal,
			html:      `<img src="https://external.com/photo.jpg"/>`,
			checkURL:  `/api/imageproxy?url=`,
			expectURL: `/api/imageproxy?url=`,
		},
		{
			name:      "internal fallback: relative URL unchanged when proxy enabled with no baseUrl",
			svc:       svcInternal,
			html:      `<img src="/local/img.png"/>`,
			checkURL:  `src="/local/img.png"`,
			expectURL: `src="/local/img.png"`,
		},
		{
			name:      "internal fallback: attachment still goes to /api/assets/",
			svc:       svcInternal,
			html:      `<img src="/attachments/abc123.png"/>`,
			checkURL:  `/api/assets/abc123.png`,
			expectURL: `/api/assets/abc123.png`,
		},
		{
			name:      "internal fallback: internal /api/assets/ URL unchanged",
			svc:       svcInternal,
			html:      `<img src="https://blog.example.com/api/assets/xyz.png"/>`,
			checkURL:  `src="https://blog.example.com/api/assets/xyz.png"`,
			expectURL: `src="https://blog.example.com/api/assets/xyz.png"`,
		},
		{
			name:      "internal fallback: proxy base URL unchanged to avoid loop",
			svc:       svcInternal,
			html:      `<img src="/api/imageproxy?url=https://other.com/img.jpg"/>`,
			checkURL:  `/api/imageproxy?url=`,
			expectURL: `/api/imageproxy?url=`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.svc.processContent(tt.html)
			if !strings.Contains(result, tt.expectURL) {
				t.Errorf("\ninput:  %s\noutput: %s\nexpected to contain: %s", tt.html, result, tt.expectURL)
			}
		})
	}
}

func TestExtractAttachmentId(t *testing.T) {
	tests := []struct {
		src  string
		want string
	}{
		{"/attachments/abc123.png", "abc123.png"},
		{"/attachments/abc123/image/photo.png", "abc123"},
		{"attachments/abc123/image/photo.png", "abc123"},
		{"api/attachments/MDblUNpajcNC/image/this.png", "MDblUNpajcNC"},
		{"/api/attachments/xyrZy5SotOgq/image/image.png", "xyrZy5SotOgq"},
		{"api/attachments/abc123", "abc123"},
		{"/api/attachments/abc123", "abc123"},
		{"https://external.com/photo.jpg", ""},
		{"/local/img.png", ""},
		{"", ""},
	}
	for _, tt := range tests {
		got := extractAttachmentId(tt.src)
		if got != tt.want {
			t.Errorf("extractAttachmentId(%q) = %q, want %q", tt.src, got, tt.want)
		}
	}
}

func TestProcessContent_RelativeAttachments(t *testing.T) {
	svc := &Service{
		imageProxyEnabled: true,
		imageProxyBaseUrl: "http://imgproxy.example.com",
		domain:            "https://blog.example.com",
	}

	tests := []struct {
		name      string
		html      string
		expectURL string
	}{
		{
			name:      "relative api/attachments with subpath",
			html:      `<img src="api/attachments/MDblUNpajcNC/image/this.png"/>`,
			expectURL: "/api/assets/MDblUNpajcNC",
		},
		{
			name:      "absolute /api/attachments with subpath",
			html:      `<img src="/api/attachments/xyrZy5SotOgq/image/image.png"/>`,
			expectURL: "/api/assets/xyrZy5SotOgq",
		},
		{
			name:      "relative attachments without api prefix",
			html:      `<img src="attachments/abc123/image/photo.png"/>`,
			expectURL: "/api/assets/abc123",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := svc.processContent(tt.html)
			if !strings.Contains(result, tt.expectURL) {
				t.Errorf("\ninput:  %s\noutput: %s\nexpected to contain: %s", tt.html, result, tt.expectURL)
			}
		})
	}
}
