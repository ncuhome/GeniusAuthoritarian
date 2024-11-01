package cronAgent

import (
	"context"
	"github.com/Mmx233/BackoffCli/backoff"
	"github.com/ncuhome/GeniusAuthoritarian/internal/pkg/singleton"
	"github.com/robfig/cron/v3"
)

var Parser = cron.NewParser(
	cron.Minute | cron.Hour | cron.Dom | cron.Month | cron.Dow,
)

func New() *cron.Cron {
	return cron.New(cron.WithParser(Parser))
}

func FuncJobWithSingleton(bkInstance backoff.Backoff) cron.FuncJob {
	single := singleton.New(bkInstance.Run)
	return func() {
		_ = single.Run(context.Background())
	}
}
