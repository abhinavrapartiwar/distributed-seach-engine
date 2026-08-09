package analysis

func Analyse(text string) []Token {
	tokens := Tokenize(text)
	normalisedToken := NormalizeAndLowercase(tokens)
	stopWordsRemovedTokens := StopWordRemoval(normalisedToken)
	stemmedTokens := StemTokens(stopWordsRemovedTokens)
	return stemmedTokens
}
