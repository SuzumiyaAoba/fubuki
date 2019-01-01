package util

import (
	"unicode"
	"unicode/utf8"
)

func reverseStrings(xs []string) {
	for i, j := 0, len(xs)-1; i < j; i, j = i+1, j-1 {
		xs[i], xs[j] = xs[j], xs[i]
	}
}

func IsDigit(r rune) bool {
	return '0' <= r && r <= '9'
}

func IsLetter(r rune) bool {
	return 'a' <= r && r <= 'z' ||
		'A' <= r && r <= 'Z' ||
		r == '_' ||
		r >= utf8.RuneSelf && unicode.IsLetter(r)
}
