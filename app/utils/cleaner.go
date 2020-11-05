package utils

import (
	"fmt"
	"strings"
)

// Cleaner defines behavior for a content cleaner.
type Cleaner interface {
	Clean(text string) (string, int)
}

// MessageCleaner ...
type MessageCleaner struct {
	BadWords     []string
	CensorTokens []string
}

// NewMessageCleaner ...
func NewMessageCleaner(badWords []string) *MessageCleaner {
	return &MessageCleaner{
		BadWords:     badWords,
		CensorTokens: []string{"#", "%", "@", "$", "!"},
	}
}

// makeReplacement masks the given word with censor tokens.
// Example:  badword  ->  b#%@$!#
func (t *MessageCleaner) makeReplacement(word string) string {
	if len(word) <= 1 {
		return word
	}

	fragment := ""
	for i := 0; i < len(word[1:]); i++ {
		fragment += t.CensorTokens[i%len(t.CensorTokens)]
	}
	return fmt.Sprintf(`<span class="censored">%v<span class="word">%v</span></span>`, word[:1], fragment)
}

// Clean ...
func (t *MessageCleaner) Clean(text string) (string, int) {
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

// Ensure the MessageCleaner type satisfies the Cleaner interface.
var _ Cleaner = &MessageCleaner{}
