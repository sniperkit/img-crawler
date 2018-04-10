package controller

import (
	"img-crawler/src/log"

	"github.com/gocolly/colly"
	//"img-crawler/src/dao"
)

type Task struct {
	name  string
	seeds []string
	desc  string
	C     *colly.Collector
}

func NewTask(name, desc string, seeds []string) *Task {

	return &Task{
		name:  name,
		seeds: seeds,
		desc:  desc,
		C:     CreateCollector()}
}

func (task *Task) Init() {
	c := task.C

	c.OnRequest(func(r *colly.Request) {
		log.Infoln("Visiting", r.URL.String())
	})

	c.OnError(func(r *colly.Response, err error) {
		log.Infoln("Error Request:", r.Request.URL)
	})

}

func (task *Task) Do() {

	log.Infof("Job %s Begin", task.name)
	task.Init()
	task.createTask()
	for _, url := range task.seeds {
		task.C.Visit(url)
	}

	// wait all requests(threads) return
	task.C.Wait()
	log.Infof("Job %s Done!", task.name)
}

// insert into task
func (task *Task) createTask() {

}
