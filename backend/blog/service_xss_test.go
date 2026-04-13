package blog

import (
	"strings"
	"testing"
)

func TestSanitizeContent_RemovesDangerousMarkup(t *testing.T) {
	service := &Service{}
	input := `
		<h1>Hello</h1>
		<p onclick="alert(1)">safe text</p>
		<script>alert("xss")</script>
		<img src="x" onerror="alert(2)">
		<a href="javascript:alert(3)">click</a>
		<iframe src="https://evil.example"></iframe>
	`

	output := service.sanitizeContent(input)

	blockedFragments := []string{
		"<script",
		"onclick=",
		"onerror=",
		"javascript:alert",
		"<iframe",
	}
	for _, fragment := range blockedFragments {
		if strings.Contains(strings.ToLower(output), strings.ToLower(fragment)) {
			t.Fatalf("sanitized article HTML still contains blocked fragment %q: %s", fragment, output)
		}
	}

	if !strings.Contains(output, "<h1") || !strings.Contains(output, "safe text") {
		t.Fatalf("expected safe article content to be preserved, got %s", output)
	}
}

func TestSearchSanitization_StripsExecutableContentFromSnippetSource(t *testing.T) {
	input := `
		<p>Hello world</p>
		<script>alert("xss")</script>
		<a href="javascript:alert(1)">Click me</a>
		<svg><script>alert(2)</script></svg>
	`

	sanitized := sanitizeSearchContent(input)
	plainText := htmlToPlainText(sanitized)
	snippet := extractSnippet(plainText, "hello")

	blockedFragments := []string{
		"<script",
		"javascript:",
		"<svg",
		"alert(",
	}
	for _, fragment := range blockedFragments {
		if strings.Contains(strings.ToLower(sanitized), strings.ToLower(fragment)) ||
			strings.Contains(strings.ToLower(plainText), strings.ToLower(fragment)) ||
			strings.Contains(strings.ToLower(snippet), strings.ToLower(fragment)) {
			t.Fatalf("search sanitization leaked dangerous fragment %q", fragment)
		}
	}

	if !strings.Contains(strings.ToLower(snippet), "hello world") {
		t.Fatalf("expected safe snippet text to remain, got %q", snippet)
	}
}
