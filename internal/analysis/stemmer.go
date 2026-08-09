package analysis

import "github.com/kljensen/snowball"

func StemTokens(tokens []Token) []Token {
	for i := range tokens {
		stemmed, err := snowball.Stem(tokens[i].Text, "english", true)
		if err != nil {
			continue
		}
		tokens[i].Text = stemmed
	}

	return tokens
}
