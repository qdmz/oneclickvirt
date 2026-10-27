package utils

import (
	"strings"
	"testing"
)

func TestSanitizeHTMLRemovesActiveContent(t *testing.T) {
	input := `<p onclick="alert(1)">ok<script>alert(1)</script><a href="javascript:alert(1)" onmouseover="x()">bad</a><a href="https://example.com" target="_blank">good</a></p>`
	got := SanitizeHTML(input)

	for _, blocked := range []string{"script", "onclick", "onmouseover", "javascript:"} {
		if strings.Contains(strings.ToLower(got), blocked) {
			t.Fatalf("sanitized HTML still contains %q: %s", blocked, got)
		}
	}
	if !strings.Contains(got, `<a href="https://example.com" target="_blank" rel="noopener noreferrer">good</a>`) {
		t.Fatalf("safe link was not preserved correctly: %s", got)
	}
}

func TestSanitizeHTMLRemovesIframeAndDangerousURLs(t *testing.T) {
	input := `<div><iframe src="https://evil.example/x"></iframe><img src="https://example.com/a.png" alt="a"><a href="data:text/html;base64,AAAA">d</a><a href="/relative/path">rel</a></div>`
	got := SanitizeHTML(input)

	for _, blocked := range []string{"iframe", "data:text/html", "evil.example"} {
		if strings.Contains(strings.ToLower(got), blocked) {
			t.Fatalf("sanitized HTML still contains %q: %s", blocked, got)
		}
	}
	// Note: the first-pass whitelist only allows absolute http/https/mailto
	// hrefs, so the relative link keeps its text but loses the href.
	for _, expected := range []string{`<img src="https://example.com/a.png" alt="a">`, `<a rel="noopener noreferrer">rel</a>`} {
		if !strings.Contains(got, expected) {
			t.Fatalf("expected sanitized HTML to contain %q, got: %s", expected, got)
		}
	}
}

func TestMarkdownToSafeHTMLStripsRawScript(t *testing.T) {
	got := MarkdownToSafeHTML("hello <script>alert(1)</script> **bold**")
	if strings.Contains(strings.ToLower(got), "script") || strings.Contains(got, "alert(1)") {
		t.Fatalf("markdown HTML still contains script content: %s", got)
	}
	if !strings.Contains(got, "<strong>bold</strong>") {
		t.Fatalf("markdown formatting lost: %s", got)
	}
}

func TestSanitizeHTMLPreservesSafeEditorFormatting(t *testing.T) {
	input := `<p class="ql-align-center evil" style="color: rgb(230, 0, 0); background-image: url(javascript:alert(1)); text-align: center"><span style="background-color: #fff">hi</span><sub>2</sub><sup>3</sup><img src="data:image/png;base64,AAAA" width="120" onerror="alert(1)"></p>`
	got := SanitizeHTML(input)

	for _, expected := range []string{
		`class="ql-align-center"`,
		`color: rgb(230, 0, 0)`,
		`text-align: center`,
		`background-color: #fff`,
		`<sub>2</sub>`,
		`<sup>3</sup>`,
		`<img src="data:image/png;base64,AAAA" width="120">`,
	} {
		if !strings.Contains(got, expected) {
			t.Fatalf("expected sanitized HTML to contain %q, got: %s", expected, got)
		}
	}

	for _, blocked := range []string{"evil", "background-image", "javascript:", "onerror"} {
		if strings.Contains(strings.ToLower(got), blocked) {
			t.Fatalf("sanitized HTML still contains %q: %s", blocked, got)
		}
	}
}
