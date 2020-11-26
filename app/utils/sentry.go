package utils

import (
	"log"
	"time"

	"github.com/getsentry/sentry-go"
	"github.com/tunedmystic/commits.lol/app/config"
)

// SetupSentry ...
func SetupSentry() func() {
	err := sentry.Init(sentry.ClientOptions{
		Dsn: config.App.SentryDSN,
	})

	if err != nil {
		log.Fatalf("sentry.Init: %s", err)
	}

	return func() {
		sentry.Flush(2 * time.Second)
	}
}
