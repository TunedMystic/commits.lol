package utils

import (
	"time"

	"github.com/getsentry/sentry-go"
	"github.com/tunedmystic/commits.lol/app/config"
)

// SetupSentry ...
func SetupSentry() func() {
	sentry.Init(sentry.ClientOptions{
		Dsn:         config.App.SentryDSN,
		Environment: config.App.Environment,
	})
	return func() {
		sentry.Flush(2 * time.Second)
	}
}
