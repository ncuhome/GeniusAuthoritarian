package agent

import (
	"github.com/Mmx233/tool"
	"github.com/robfig/cron/v3"
	"time"
)

var Parser cron.Parser

func init() {
	Parser = cron.NewParser(
		cron.Minute | cron.Hour | cron.Dom | cron.Month | cron.Dow,
	)

	c = cron.New(cron.WithParser(Parser))

	c.Start()
}

var c *cron.Cron

type Event struct {
	T string
	E func()
}

// AddRegular
//
// e.Spec: Minute Hour Dom Month Dow
func AddRegular(e *Event) (cron.EntryID, error) {
	return c.AddFunc(e.T, func() {
		defer tool.Recover()
		e.E()
	})
}

type offsetStruct struct {
	schedule cron.Schedule
	value    time.Duration
}

func (a offsetStruct) Next(t time.Time) time.Time {
	return a.schedule.Next(t).Add(a.value)
}

func AddWithOffset(e *Event, offset time.Duration) (cron.EntryID, error) {
	schedule, err := Parser.Parse(e.T)
	if err != nil {
		return 0, err
	}
	return c.Schedule(offsetStruct{
		schedule: schedule,
		value:    offset,
	}, cron.FuncJob(e.E)), nil
}

func CancelRegular(ID cron.EntryID) {
	c.Remove(ID)
}

func Start() {
	c.Start()
}

func Stop() {
	c.Stop()
}
