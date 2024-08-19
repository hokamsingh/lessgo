package scheduler

import (
	"github.com/robfig/cron/v3"
)

type Scheduler interface {
	AddJob(schedule string, job func()) error
	Start()
	Stop()
}

type CronScheduler struct {
	cron *cron.Cron
}

func NewCronScheduler() *CronScheduler {
	return &CronScheduler{
		cron: cron.New(),
	}
}

func (s *CronScheduler) AddJob(schedule string, job func()) error {
	_, err := s.cron.AddFunc(schedule, job)
	if err != nil {
		return err
	}
	return nil
}

func (s *CronScheduler) Start() {
	s.cron.Start()
}

func (s *CronScheduler) Stop() {
	s.cron.Stop()
}
