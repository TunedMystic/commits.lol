package utils

import (
	"fmt"
	"strings"
)

// TextCleaner ...
type TextCleaner struct {
	BadWords     []string
	CensorTokens []string
}

// NewTextCleaner ...
func NewTextCleaner(badWords []string) *TextCleaner {
	return &TextCleaner{
		BadWords:     badWords,
		CensorTokens: []string{"#", "%", "@", "$", "!"},
	}
}

// makeReplacement masks the given word with censor tokens.
// Example:  badword  ->  b#%@$!#
func (t *TextCleaner) makeReplacement(word string) string {
	if len(word) <= 1 {
		return word
	}

	fragment := ""
	for i := 0; i < len(word[1:]); i++ {
		fragment += t.CensorTokens[i%len(t.CensorTokens)]
	}
	return fmt.Sprintf(`<span class="censored">%v<span class="word">%v</span></span>`, word[:1], fragment)
}

// CensorText ...
func (t *TextCleaner) CensorText(text string) (string, int) {
	newText := []string{}
	censoredWords := 0
	for _, word := range strings.Fields(text) {
		replacement := ""

		for _, badWord := range t.BadWords {
			if strings.Contains(strings.ToLower(word), badWord) {
				replacement = t.makeReplacement(word)
				censoredWords++
			}
		}
		if replacement == "" {
			replacement = word
		}
		newText = append(newText, replacement)
	}
	return strings.Join(newText, " "), censoredWords
}
