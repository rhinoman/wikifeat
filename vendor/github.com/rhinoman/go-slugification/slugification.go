// Package slugification provides methods for generating 'slugged' versions
// of strings suitable for use in URLs.
package slugification

import (
	"strings"
	"unicode"
)

// Returns a slugified string
// Example: "Page Title" becomes "page-title"
func Slugify(inputString string) string {
	replaceChar := func(r rune) rune {
		switch {
		case unicode.IsLetter(r):
			return unicode.ToLower(r)
		case unicode.IsNumber(r), r == '_', r == '-', r == '+':
			return r
		case unicode.IsSpace(r):
			return '-'
		default:
			return -1
		}
	}
	return strings.Map(replaceChar, inputString)
}
