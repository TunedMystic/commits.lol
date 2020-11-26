package utils

import "testing"

func Test_Sentry(t *testing.T) {
	flushSentry := SetupSentry()
	defer flushSentry()
}
