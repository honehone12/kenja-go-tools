package opensearch

import (
	"slices"
	"strings"
	"testing"
	"unicode"
)

func spaceIsSpace(r rune) rune {
	if unicode.IsSpace(r) {
		return ' '
	} else {
		return r
	}
}

func cleanKeywords(replacer *strings.Replacer, raw string) string {
	s := replacer.Replace(raw)
	s = strings.Map(spaceIsSpace, s)
	sep := strings.Split(s, " ")
	sep = slices.DeleteFunc(sep, func(s string) bool {
		return len(strings.TrimSpace(s)) == 0
	})
	return strings.Join(sep, " ")
}

func TestOpensearchKeywords(t *testing.T) {
	s := "_+-*%&|=<>!?^~{}[]():,./'`"
	clean := cleanKeywords(__OPENSEARCH_REPLACER, s)
	if clean != "" {
		panic("replacer fail")
	}

	s = `"\`
	clean = cleanKeywords(__OPENSEARCH_REPLACER, s)
	if clean != "" {
		panic("replacer fail")
	}

	s = "  a  b    c         de "
	clean = cleanKeywords(__OPENSEARCH_REPLACER, s)
	if clean != "a b c de" {
		panic("replacer fail")
	}

	s = "\ta\nb \vc\fd\r\ne"
	clean = cleanKeywords(__OPENSEARCH_REPLACER, s)
	if clean != "a b c d e" {
		panic("replacer fail")
	}

	s = "school music band club"
	clean = cleanKeywords(__OPENSEARCH_REPLACER, s)
	if clean != "school music band club" {
		panic("breaking keywords")
	}
}
