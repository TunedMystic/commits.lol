package utils

import (
	"fmt"
	"strings"
)

// TextCleaner ...
type TextCleaner struct {
	BadWords []string
	Token    []string
}

// NewTextCleaner ...
func NewTextCleaner() *TextCleaner {
	return &TextCleaner{
		// TODO: Put bad words here...
		BadWords: []string{},
		// Token: "x",
		Token: []string{"#", "%", "@", "$", "!"},
	}
}

func (t *TextCleaner) makeReplacement(word string) string {
	/*
		shitty
		<span class="censored">sh<span>xxxx</span></span>
	*/
	if len(word) <= 2 {
		return `<span class="censored"><span class="word">XX</span></span>`
	}

	fragment := ""
	for i := 0; i < len(word[2:]); i++ {
		fragment += t.Token[i%len(t.Token)]
	}
	return fmt.Sprintf(`<span class="censored">%v<span class="word">%v</span></span>`, word[:2], fragment)
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
