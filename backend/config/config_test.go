package config

import "testing"

func TestAISummaryConfigNormalization(t *testing.T) {
	t.Setenv("TRILIUM_API_URL", "https://trilium.example.com")
	t.Setenv("TRILIUM_TOKEN", "token")
	t.Setenv("AI_SUMMARY_ENABLED", "true")
	t.Setenv("AI_SUMMARY_PROVIDER", "openai")
	t.Setenv("AI_SUMMARY_MODE", "AI")
	t.Setenv("AI_SUMMARY_TIMEOUT_MS", "45000")
	t.Setenv("AI_SUMMARY_MAX_INPUT_CHARS", "9000")

	LoadConfig()

	if Config.AISummary.Provider != "openai-compatible" {
		t.Fatalf("expected normalized provider, got %q", Config.AISummary.Provider)
	}
	if Config.AISummary.Mode != "ai" {
		t.Fatalf("expected normalized mode, got %q", Config.AISummary.Mode)
	}
	if !Config.AISummary.AIRequestsEnabled() {
		t.Fatalf("expected AI requests to be enabled")
	}
	if Config.AISummary.TimeoutMs != 45000 {
		t.Fatalf("expected timeout to be loaded, got %d", Config.AISummary.TimeoutMs)
	}
	if Config.AISummary.MaxInputChars != 9000 {
		t.Fatalf("expected max input chars to be loaded, got %d", Config.AISummary.MaxInputChars)
	}
}

func TestAISummaryModeDefaultsToCode(t *testing.T) {
	t.Setenv("TRILIUM_API_URL", "https://trilium.example.com")
	t.Setenv("TRILIUM_TOKEN", "token")
	t.Setenv("AI_SUMMARY_ENABLED", "true")
	t.Setenv("AI_SUMMARY_MODE", "unexpected")

	LoadConfig()

	if Config.AISummary.Mode != "code" {
		t.Fatalf("expected invalid mode to fall back to code, got %q", Config.AISummary.Mode)
	}
	if Config.AISummary.AIRequestsEnabled() {
		t.Fatalf("expected AI requests to stay disabled in code mode")
	}
}
