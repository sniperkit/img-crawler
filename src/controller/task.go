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
	C     []*colly.Collector
}

func NewTaskController(name, desc string, seeds []string, num_cc int) *Task {

	if num_cc < 1 {
		num_cc = 1
	}

	C := make([]*colly.Collector, num_cc)
	C[0] = CreateCollector()

	for k := 1; k < len(C); k++ {
		C[k] = C[0].Clone()
	}

	return &Task{
		name:  name,
		seeds: seeds,
		desc:  desc,
		C:     C,
	}
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

	task.GeneralCB(task.C...)

	id, err := task.createTask()
	log.Infof("chenqi\t", id, err)

	for _, url := range task.seeds {
		task.C[0].Visit(url)
	}

	// wait all requests(threads) finished
	for _, v := range task.C {
		v.Wait()
	}

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
