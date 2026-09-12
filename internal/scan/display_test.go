package scan

import "testing"

func TestSanitizeDisplayText(t *testing.T) {
	tests := map[string]struct {
		input string
		want  string
	}{
		"layout whitespace": {
			input: "a\nb\rc\td\ve\f",
			want:  "a b c d e ",
		},
		"terminal and spoofing controls": {
			input: "a\x1b\a\u0085\u009db\u202Ec\u200Bd\uFEFFe\u00ADf\u2028g\u2029h",
			want:  "abcdefgh",
		},
		"readable unicode survives": {
			input: "box─ · ⑂ CJK 日本 emoji 😀 café",
			want:  "box─ · ⑂ CJK 日本 emoji 😀 café",
		},
	}
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			if got := SanitizeDisplayText(tc.input); got != tc.want {
				t.Errorf("SanitizeDisplayText(%q) = %q, want %q", tc.input, got, tc.want)
			}
		})
	}
}
