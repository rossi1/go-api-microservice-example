package internal

import (
	"time"

	"github.com/getsentry/sentry-go"
)

func NewSentry(cfg Config) error {
	err := sentry.Init(sentry.ClientOptions{
		Dsn:         cfg.DSN,
		Environment: cfg.ENVIRONMENT,
		Release:     "my-project-name@1.0.0",
		Debug:       cfg.DEBUG,
	})
	if err != nil {
		return err
	}

	defer sentry.Flush(2 * time.Second)

	return nil
}
