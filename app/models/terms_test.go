package models

import (
	"testing"

	u "github.com/tunedmystic/commits.lol/app/utils"
)

func SliceContains(items []string, value string) bool {
	for _, item := range items {
		if item == value {
			return true
		}
	}
	return false
}

func Test_BadWords_ToStrings(t *testing.T) {
	bw := BadWords{
		{1, "crappy"},
		{2, "boring"},
		{3, "sleepy"},
	}

	words := bw.ToStrings()

	u.AssertEqual(t, len(words), 3)
	u.AssertEqual(t, SliceContains(words, "crappy"), true)
	u.AssertEqual(t, SliceContains(words, "boring"), true)
	u.AssertEqual(t, SliceContains(words, "sleepy"), true)
	u.AssertEqual(t, SliceContains(words, "not-here"), false)
}

func Test_SearchTerms_ToStrings(t *testing.T) {
	st := SearchTerms{
		{1, "great", 4},
		{2, "awesome", 1},
		{3, "nice", 2},
	}

	terms := st.ToStrings()

	u.AssertEqual(t, len(terms), 3)
	u.AssertEqual(t, SliceContains(terms, "great"), true)
	u.AssertEqual(t, SliceContains(terms, "awesome"), true)
	u.AssertEqual(t, SliceContains(terms, "nice"), true)
	u.AssertEqual(t, SliceContains(terms, "not-here"), false)
}

func Test_GroupTerms_ToMap(t *testing.T) {
	gt := GroupTerms{
		{1, "lol", "funny"},
		{2, "wtf", "funnny"},
		{3, "omg", "angry"},
	}

	terms := gt.ToMap()

	u.AssertEqual(t, len(terms), 3)

	var ok bool
	_, ok = terms["lol"]
	u.AssertEqual(t, ok, true)

	_, ok = terms["wtf"]
	u.AssertEqual(t, ok, true)

	_, ok = terms["omg"]
	u.AssertEqual(t, ok, true)

	_, ok = terms["not-here"]
	u.AssertEqual(t, ok, false)
}
