package scan

import (
	"strings"
	"unicode"
)

// SanitizeDisplayText removes terminal-affecting and display-spoofing runes from
// untrusted text while keeping layout whitespace readable as ordinary spaces.
func SanitizeDisplayText(s string) string {
	return strings.Map(func(r rune) rune {
		switch r {
		case '\n', '\r', '\t', '\v', '\f':
			return ' '
		}
		if unicode.IsControl(r) || unicode.In(r, unicode.Cf, unicode.Zl, unicode.Zp) {
			return -1
		}
		return r
	}, s)
}
