package crond

import (
	"time"

	"poison/slf4go"
)

type Crond struct {
	// 定时任务列表
	Jobs []*CrondJob

	// 定时任务执行时间
	When int64
}

func (this *Crond) AddJob(job *CrondJob) {
	this.Jobs = append(this.Jobs, job)
}

var cronds = []*Crond{}

func NewCrond(job *CrondJob, t time.Time) *Crond {
	return &Crond{
		Jobs: []*CrondJob{job},
		When: t.Unix(),
	}
}

func AddCrondJob(job *CrondJob) {
	tn := time.Now()
	loc, _ := time.LoadLocation("Local")

	t := GetExecTime(tn, loc, job.Timer)

	slf4go.Info("Add Crond Job: %v, Exec Time: %v,", job.Name, t)

	for _, crond := range cronds {
		if crond.When == t.Unix() {
			crond.AddJob(job)
			return
		}
	}

	c := NewCrond(job, t)
	cronds = append(cronds, c)

	go Start(t.Sub(tn), c)
}

func Start(d time.Duration, crond *Crond) {
	timer := time.NewTimer(d)

	<-timer.C

	for _, job := range crond.Jobs {
		slf4go.Info("job %v exec now.", job.Name)
		job.Func()
	}

	End(crond.When)
}

func End(when int64) {
	var idx int
	for index, v := range cronds {
		if v.When == when {
			idx = index
			break
		}
	}

	jobs := cronds[idx].Jobs

	cronds = append(cronds[:idx], cronds[idx+1:]...)

	slf4go.Debug("%v", cronds)

	for _, job := range jobs {
		AddCrondJob(job)
	}
}
