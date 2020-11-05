package utils

import (
	"testing"
)

func Test_Grouper(t *testing.T) {
	type test struct {
		label    string
		text     string
		expected string
	}

	tests := []test{
		{
			"text_with_one_term",
			"lol that was a crazy bug",
			"funny",
		},
		{
			"text_with_two_terms",
			"lol that was a crazy haha bug",
			"funny",
		},
		{
			"text_with_mixed_terms",
			"i was all lol at first, and then all argh later",
			"angry",
		},
		{
			"text_with_no_terms",
			"that was a crazy bug",
			"",
		},
	}

	groupedTerms := map[string]string{
		"lol":  "funny",
		"haha": "funny",
		"argh": "angry",
	}
	grouper := NewCommitGrouper(groupedTerms)

	for _, test := range tests {
		t.Run(test.label, func(t *testing.T) {
			AssertEqual(t, grouper.Group(test.text), test.expected)
		})
	}
}
