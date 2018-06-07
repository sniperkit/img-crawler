package controller

import (
	"database/sql"
	"img-crawler/src/dao"
	"img-crawler/src/log"
	"strings"
	"sync"

	"github.com/gocolly/colly"
)

type Login struct {
	Action func(c *colly.Collector) error
}

type Task struct {
	Id    uint64
	name  string
	seeds []string
	desc  string
	C     []*colly.Collector
	Logon *Login
	retry map[string]uint8
	lock  *sync.Mutex
}

func NewTaskController(name, desc string, seeds []string, num_cc int, login *Login) *Task {

	if num_cc < 1 {
		num_cc = 1
	}

	C := make([]*colly.Collector, num_cc)
	l := CreateCollector()

	// login
	if login != nil {
		//login.Action(l)
		l.Wait()
	}

	// clone
	for k := 0; k < len(C); k++ {
		C[k] = l.Clone()
	}

	return &Task{
		name:  name,
		seeds: seeds,
		desc:  desc,
		C:     C,
		retry: make(map[string]uint8),
		lock:  &sync.Mutex{},
	}
}

// general call back
func (task *Task) GeneralCB(cs ...*colly.Collector) {

	for _, c := range cs {

		c.OnRequest(func(r *colly.Request) {
			log.Info("Visiting:", r.URL.String())
			/*
				log.Info("Host:", r.Headers.Get("Host"))
				log.Info("Cookie:", r.Headers.Get("Cookie"))
				log.Info("Referer:", r.Headers.Get("Referer"))
				log.Info("User-Agent:", r.Headers.Get("User-Agent"))
			*/
		})

		c.OnResponse(func(r *colly.Response) {
			/*
				log.Info("Response:", r.Request.URL.String())
				log.Info("Status:", r.StatusCode)
				log.Info("Set-Cookie:", r.Headers.Get("Set-Cookie"))
			*/
		})

		c.OnError(func(r *colly.Response, err error) {
			log.Warnf("Error Request %s [%d]%s",
				r.Request.URL.String(), r.StatusCode, err)
			if r.StatusCode == 503 {
				task.Retry(r.Request, 3)
			}
		})
	}
}

func (task *Task) getRetryCnt(url string, max uint8) uint8 {
	task.lock.Lock()
	defer task.lock.Unlock()

	if _, ok := task.retry[url]; ok {
		task.retry[url]++
	} else {
		task.retry[url] = 1
	}

	cnt := task.retry[url]
	if cnt > max {
		//delete(task.retry, url)
	}
	return cnt
}

func (task *Task) Retry(r *colly.Request, max uint8) {
	url := r.URL.String()
	cnt := task.getRetryCnt(url, max)
	if cnt <= max {
		r.Retry()
	}
}

func (task *Task) Do() {

	log.Infof("Job %s Begin, seeds=%s", task.name, task.seeds)

	task.GeneralCB(task.C...)
	task.createTask()

	for _, url := range task.seeds {
		task.C[0].Visit(url)
	}

	// wait all requests(threads) finished
	for _, v := range task.C {
		v.Wait()
	}

	log.Infof("Job %s Done!", task.name)
}

func (task *Task) createTask() (id uint64, err error) {
	t := new(dao.Task)
	t.Name = task.name
	t.Seeds = strings.Join(task.seeds, ",")
	if len(task.desc) > 0 {
		t.Desc = sql.NullString{task.desc, true}
	}

	id, err = taskDAO.CreateTask(t)
	if err == nil {
		task.Id = id
		taskDAO.CreateItemTable(id)

	} else {
		result, err := taskDAO.Get(map[string]interface{}{"name": t.Name}, false)
		log.Infof("name=%s", t.Name)
		if err == nil {
			task.Id = result.ID
		} else {
			log.Fatalf("CreateTask got no id")
		}
	}
	return
}

func (task *Task) CreateTaskItem(name, url, desc, digest, filepath string, status int) {
	item := new(dao.TaskItem)
	item.TaskID = task.Id
	item.Name = name
	item.Url = url
	item.Status = status
	if len(desc) > 0 {
		item.Desc = sql.NullString{desc, true}
	}
	if len(digest) > 0 {
		item.Digest = sql.NullString{digest, true}
	}
	if len(filepath) > 0 {
		item.FilePath = sql.NullString{filepath, true}
	}

	_, err := taskDAO.CreateTaskItem(item, task.Id)
	if err != nil {
		log.Errorf("CreateTaskItem failed %s,%s,%s", name, url, desc)
	}
}

var (
	taskDAO *dao.TaskDAOImpl
)

func init() {
	taskDAO = dao.NewTaskDAO(dao.Mpool)
}
