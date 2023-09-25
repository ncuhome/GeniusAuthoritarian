package cronAgent

import (
	"github.com/robfig/cron/v3"
)

var Parser = cron.NewParser(
	cron.Minute | cron.Hour | cron.Dom | cron.Month | cron.Dow,
)

func New() *cron.Cron {
	return cron.New(cron.WithParser(Parser))
}
