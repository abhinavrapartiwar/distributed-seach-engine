package analysis

import (
	"strings"

	"golang.org/x/text/unicode/norm"
)

func NormalizeAndLowercase(tokens []Token) []Token {
	for i := range tokens {
		text := norm.NFC.String(tokens[i].Text)
		text = strings.ToLower(text)

		tokens[i].Text = text
	}

	return tokens
}

