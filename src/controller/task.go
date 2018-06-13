package controller

import (
	"database/sql"
	"img-crawler/src/dao"
	"img-crawler/src/log"
	"strings"
	"sync"
    "time"
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
    download *colly.Collector
	Logon *Login
	retry map[string]uint8
	lock  *sync.Mutex
}

func NewTaskController(name, desc string, seeds []string, 
        num_cc int, download_pic bool, login *Login) *Task {

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
        download: C[0],
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

func (task *Task) createTask() (err error) {
	t := new(dao.Task)
	t.Name = task.name
	t.Seeds = strings.Join(task.seeds, ",")
	if len(task.desc) > 0 {
		t.Desc = sql.NullString{task.desc, true}
	}

	task.Id, err = taskDAO.CreateTask(t)
	if err != nil {
		result, err := taskDAO.Get(map[string]interface{}{"name": t.Name})
		if err == nil {
			task.Id = result.ID
		} else {
			log.Fatalf("CreateTask failed %s", t.Name)
		}
	}

    taskDAO.CreateItemTable(task.Id)
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

func (task *Task) UpdateTaskItem(name, url, desc, digest, filepath string, status int) {

    // where
    c := map[string]interface{}{
        "name": name,
        "url": url,
    }

    // set
    v := map[string]interface{}{
        "status": status,
        "filepath": filepath,
        "digest": digest,
    }

	_, err := taskDAO.Update(true, c, v)
	if err != nil {
		log.Errorf("UpdateTaskItem failed %s,%s,%s", name, url, desc)
	}
}

func (task *Task) DownloadImg() {
 //   taskDAO.ListItems(Download_INIT)
//    taskDAO.ListItems(Download_SAVEFAIL)

	task.createTask()
    var num uint64 = 100
//    task.download.Async = false
    Download(task.download)
    for {
        items, err := taskDAO.ListItems(Download_DownFAIL, num)
        if err != nil {
            log.Errorf("ListItems error")
            return;
        }

        for _,item := range items {

            ctx := colly.NewContext()
            ctx.Put("name", item.Name)
            if item.Desc.Valid {
                ctx.Put("desc", item.Desc.String)
            }
            ctx.Put("task", task)
            task.download.Request("GET", item.Url, nil, ctx, nil)
            time.Sleep(time.Duration(10) * time.Millisecond)
        }

        task.download.Wait()
    }


}

var (
	taskDAO *dao.TaskDAOImpl
)

func init() {
	taskDAO = dao.NewTaskDAO(dao.Mpool)
}
