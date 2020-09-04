package server

import (
	"testing"

	"github.com/matryer/is"
)

func Test_SomeHandler(t *testing.T) {
	is := is.New(t)
	t.Log("Testing SomeHandler...")
	is.Equal(1+1, 2)
}
