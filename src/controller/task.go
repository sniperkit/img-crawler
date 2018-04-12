package controller

import (
	"database/sql"
	"img-crawler/src/dao"
	"img-crawler/src/log"
	"strings"

	"github.com/gocolly/colly"
)

type Task struct {
	name  string
	seeds []string
	desc  string
	C     *colly.Collector
}

func NewTaskController(name, desc string, seeds []string) *Task {

	return &Task{
		name:  name,
		seeds: seeds,
		desc:  desc,
		C:     CreateCollector()}
}

// general call back
func (task *Task) GeneralCB(cs ...*colly.Collector) {

	for _, c := range cs {

		c.OnRequest(func(r *colly.Request) {
			log.Infoln("Visiting", r.URL.String())
		})

		c.OnError(func(r *colly.Response, err error) {
			log.Warnf("Error Request %s %s", r.Request.URL.String(), err)
		})
	}

}

func (task *Task) Do() {

	log.Infof("Job %s Begin, seeds=%s", task.name, task.seeds)
	task.createTask()
	for _, url := range task.seeds {
		task.C.Visit(url)
	}

	// wait all requests(threads) return
	task.C.Wait()
	log.Infof("Job %s Done!", task.name)
}

func (task *Task) createTask() (uint64, error) {
	t := new(dao.Task)
	t.Name = task.name
	t.Seeds = strings.Join(task.seeds, ",")
	if len(task.desc) > 0 {
		t.Desc = sql.NullString{task.desc, true}
	}

	return taskDAO.Create(t)
}

var (
	taskDAO *dao.TaskDAOImpl
)

func init() {
	taskDAO = dao.NewTaskDAO(dao.Mpool)
}
