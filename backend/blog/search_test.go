package blog

import (
	"strings"
	"testing"
)

func TestSearchScore(t *testing.T) {
	score := searchScore(
		normalizeSearchText("AI 在社群中的应用"),
		normalizeSearchText("这里讨论 AI 内容生成 与 社群机器人"),
		normalizeSearchText("AI 社群"),
	)
	if score <= 0 {
		t.Fatalf("expected positive score, got %d", score)
	}
}

func TestContainsSubsequence(t *testing.T) {
	if !containsSubsequence("aisummary", "as") {
		t.Fatalf("expected subsequence match")
	}
	if containsSubsequence("blog", "ga") {
		t.Fatalf("did not expect subsequence match")
	}
}

func TestExtractSnippet(t *testing.T) {
	text := strings.Repeat("前置内容。", 20) + "这里专门讨论 AI summary 的生成过程，以及它和 code summary 的关系。" + strings.Repeat("后置内容。", 20)
	snippet := extractSnippet(text, "AI summary")
	if snippet == "" {
		t.Fatalf("expected snippet to be generated")
	}
	if snippet == text {
		t.Fatalf("expected shortened snippet around query")
	}
}
