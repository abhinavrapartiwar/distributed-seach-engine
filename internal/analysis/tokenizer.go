package analysis

import (
	"fmt"
	"unicode"
)


func spaceSeparatedWords(text string) []string {
	var words []string
	word := ""
	for _, ch := range text {
		if ch == ' ' && word == "" {
			continue
		}

		if ch == ' ' {
			words = append(words, word)
			word = ""
		} else {
			word += fmt.Sprintf("%c", ch)
		}
	}
	words = append(words, word)
	return words
}

func trimSpecialChars(word string) string {

	runes := []rune(word)
	start := 0
	end := len(runes) - 1
	// forward check
	for start < len(runes) {
		if unicode.IsDigit(runes[start]) || unicode.IsLetter(runes[start]) {
			break
		}
		start++
	}

	// trailling check
	for end >= 0 {
		if unicode.IsDigit(runes[end]) || unicode.IsLetter(runes[end]) {
			break
		}
		end--
	}
	if start > end {
		return ""
	}
	trimmed := runes[start : end+1]
	return string(trimmed)
}

func splitOnSpecialChars(text string) []string {
	var words []string
	var current []rune

	for _, ch := range text {
		if unicode.IsLetter(ch) || unicode.IsDigit(ch) {
			current = append(current, ch)
		} else {
			if len(current) > 0 {
				words = append(words, string(current))
				current = current[:0]
			}
		}
	}

	if len(current) > 0 {
		words = append(words, string(current))
	}

	return words
}

func Tokenize(text string) []Token {
	spaceSeparated := spaceSeparatedWords(text)
	var tokenized []string
	for _, word := range spaceSeparated {
		tokenized = append(tokenized, splitOnSpecialChars(word)...)
	}
	var removedEmptyTokenized []Token
	pos := 0
	for _, word := range tokenized {
		if word != "" {
			removedEmptyTokenized = append(removedEmptyTokenized, Token{Text: word, position: pos})
			pos++
		}
	}
	return removedEmptyTokenized
}

