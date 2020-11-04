package utils

import (
	"testing"
)

func Test_Cleaner(t *testing.T) {
	cleaner := NewTextCleaner([]string{"crappy", "bad"})
	text := "fixed a crappy bug"
	censoredText, censoredWords := cleaner.CensorText(text)
	expected := `fixed a <span class="censored">c<span class="word">#%@$!</span></span> bug`
	AssertEqual(t, censoredText, expected)
	AssertEqual(t, censoredWords, 1)
}

func Test_Cleaner_no_censored_words(t *testing.T) {
	cleaner := NewTextCleaner([]string{"crappy", "bad"})
	text := "fixed a bug"
	censoredText, censoredWords := cleaner.CensorText(text)
	AssertEqual(t, censoredText, text)
	AssertEqual(t, censoredWords, 0)
}

func Test_Cleaner_make_replacement(t *testing.T) {
	cleaner := NewTextCleaner([]string{})
	AssertEqual(t, cleaner.makeReplacement("crappy"), `<span class="censored">c<span class="word">#%@$!</span></span>`)
	AssertEqual(t, cleaner.makeReplacement("bad"), `<span class="censored">b<span class="word">#%</span></span>`)
	AssertEqual(t, cleaner.makeReplacement("a"), "a")
}
