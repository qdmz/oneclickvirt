package product

import (
	"strings"
	"testing"
)

func TestSanitizeSiteConfigValueCleansHeaderFooter(t *testing.T) {
	raw := `<div onclick="alert(1)"><script>alert(2)</script><p>公告</p></div>`
	for _, key := range []string{"custom_header", "custom_footer"} {
		got, ok := sanitizeSiteConfigValue(key, raw).(string)
		if !ok {
			t.Fatalf("%s: expected string result", key)
		}
		for _, blocked := range []string{"script", "onclick", "alert"} {
			if strings.Contains(strings.ToLower(got), blocked) {
				t.Fatalf("%s: sanitized value still contains %q: %s", key, blocked, got)
			}
		}
		if !strings.Contains(got, "<p>公告</p>") {
			t.Fatalf("%s: safe content lost: %s", key, got)
		}
	}
}

func TestSanitizeSiteConfigValueLeavesOtherFieldsUntouched(t *testing.T) {
	if got := sanitizeSiteConfigValue("site_name", "<b>站名</b>"); got != "<b>站名</b>" {
		t.Fatalf("unexpected change for non-HTML field: %v", got)
	}
	if got := sanitizeSiteConfigValue("custom_header", 42); got != 42 {
		t.Fatalf("non-string value should pass through, got: %v", got)
	}
}
