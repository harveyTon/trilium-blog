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
