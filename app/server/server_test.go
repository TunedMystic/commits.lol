package server

import (
	"testing"

	"github.com/tunedmystic/commits.lol/app/db"
	"github.com/tunedmystic/commits.lol/app/models"
	u "github.com/tunedmystic/commits.lol/app/utils"
)

func Test_SomeHandler(t *testing.T) {
	t.Log("Testing SomeHandler...")
	u.AssertEqual(t, 1+1, 2)
}

func Test_SomeOtherHandler(t *testing.T) {
	mockDB := db.MockDB{
		AllCommitsMock: func() (models.GitCommits, error) {
			return models.GitCommits{{Message: "Fixed a bug"}}, nil
		},
	}
	commits, _ := mockDB.AllCommits()
	u.AssertEqual(t, len(commits), 1)
	u.AssertEqual(t, commits[0].Message, "Fixed a bug")
}
