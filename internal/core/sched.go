package core

var sched *Scheduler

type Scheduler struct {
}

func NewScheduler() *Scheduler {
	return &Scheduler{}
}

func init() {
	sched = NewScheduler()
}
