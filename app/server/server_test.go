package server

import (
	"testing"

	u "github.com/tunedmystic/commits.lol/app/utils"
)

func Test_SomeHandler(t *testing.T) {
	t.Log("Testing SomeHandler...")
	u.AssertEqual(t, 1+1, 2)
}
