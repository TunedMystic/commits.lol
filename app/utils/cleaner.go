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
func NewTextCleaner() *TextCleaner {
	return &TextCleaner{
		// TODO: Put bad words here...
		BadWords:     []string{},
		CensorTokens: []string{"#", "%", "@", "$", "!"},
	}
}

func (t *TextCleaner) makeReplacement(word string) string {
	/*
		shitty
		<span class="censored">s<span>xxxxx</span></span>
	*/
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
func (t *TextCleaner) CensorText(text string) string {
	newText := []string{}
	for _, word := range strings.Fields(text) {
		replacement := ""

		for _, badWord := range t.BadWords {
			if strings.Contains(strings.ToLower(word), badWord) {
				replacement = t.makeReplacement(word)
			}
		}
		if replacement == "" {
			replacement = word
		}
		newText = append(newText, replacement)
	}
	return strings.Join(newText, " ")
}
