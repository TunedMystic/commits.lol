package utils

import (
	"sort"
	"strings"
)

// Grouper defines behavior for a content classifier.
type Grouper interface {
	Group(text string) string
}

// CommitGrouper ...
type CommitGrouper struct {
	GroupedTerms map[string]string
	SortedKeys   []string
}

// NewCommitGrouper ...
func NewCommitGrouper(groupedTerms map[string]string) *CommitGrouper {
	keys := make([]string, 0, len(groupedTerms))

	for key := range groupedTerms {
		keys = append(keys, key)
	}
	sort.Strings(keys)

	return &CommitGrouper{
		GroupedTerms: groupedTerms,
		SortedKeys:   keys,
	}
}

// Group ...
func (g *CommitGrouper) Group(text string) string {
	for _, term := range g.SortedKeys {
		if strings.Contains(strings.ToLower(text), term) {
			return g.GroupedTerms[term]
		}
	}
	return ""
}

var _ Grouper = &CommitGrouper{}
