package models

import (
	"testing"

	u "github.com/tunedmystic/commits.lol/app/utils"
)

type MockCleaner struct {
	MockClean func(text string) (string, int)
}

func (c *MockCleaner) Clean(text string) (string, int) {
	return c.MockClean(text)
}

func Test_SetCensoredMessage_updated(t *testing.T) {
	c := &MockCleaner{
		MockClean: func(text string) (string, int) {
			return text + " - updated", 1
		},
	}

	commit := GitCommit{Message: "Fixed a bug"}
	u.AssertEqual(t, commit.MessageCensored, "")

	// Censor the message, is possible.
	updated := commit.SetCensoredMessage(c)
	u.AssertEqual(t, updated, true)

	u.AssertEqual(t, commit.MessageCensored, "Fixed a bug - updated")
}

func Test_SetCensoredMessage_not_updated(t *testing.T) {
	c := &MockCleaner{
		MockClean: func(text string) (string, int) {
			return text, 0
		},
	}

	commit := GitCommit{Message: "Fixed a bug"}
	u.AssertEqual(t, commit.MessageCensored, "")

	// Censor the message, if possible.
	updated := commit.SetCensoredMessage(c)
	u.AssertEqual(t, updated, false)

	u.AssertEqual(t, commit.MessageCensored, "")
}

type MockGrouper struct {
	MockGroup func(text string) string
}

func (g *MockGrouper) Group(text string) string {
	return g.MockGroup(text)
}

func Test_SetGroup_updated(t *testing.T) {
	g := &MockGrouper{
		MockGroup: func(text string) string {
			return "funny"
		},
	}

	commit := GitCommit{Message: "Fixed a bug"}
	u.AssertEqual(t, commit.Group, "")

	// Group the message, if possible.
	updated := commit.SetGroup(g)
	u.AssertEqual(t, updated, true)

	u.AssertEqual(t, commit.Group, "funny")
}

func Test_SetGroup_not_updated(t *testing.T) {
	g := &MockGrouper{
		MockGroup: func(text string) string {
			return ""
		},
	}

	commit := GitCommit{Message: "Fixed a bug"}
	u.AssertEqual(t, commit.Group, "")

	// Group the message, if possible.
	updated := commit.SetGroup(g)
	u.AssertEqual(t, updated, false)

	u.AssertEqual(t, commit.Group, "")
}

func Test_SetColorTheme_updated(t *testing.T) {
	commit := GitCommit{
		Message: "Fixed a bug",
		Author: GitUser{
			Username: "Alice",
		},
	}
	u.AssertEqual(t, commit.ColorBackground, "")
	u.AssertEqual(t, commit.ColorForeground, "")

	// Calculate the color theme.
	commit.SetColorTheme()
	u.AssertEqual(t, commit.ColorBackground, "#ffd300")
	u.AssertEqual(t, commit.ColorForeground, "#000000")
}
