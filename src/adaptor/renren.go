package adaptor

import (
	"encoding/json"
	"fmt"
	"github.com/gocolly/colly"
	"img-crawler/src/controller"
	"img-crawler/src/log"
	"regexp"
	"strings"
)

func RenRenLogin(c *colly.Collector) (err error) {

	c.OnResponse(func(r *colly.Response) {
		//     log.Info("logon resp=", string(r.Body))
	})

	//vcode := "http://icode.renren.com/getcode.do?t=login&rnd=Math.random()"
	url := "http://www.renren.com/PLogin.do"
	err = c.Post(url, map[string]string{
		"email":     "chenny7@163.com",
		"password":  "01de442283a34e82bccaa036b91f775d70398fe90478d0c41eedb873009784ab",
		"autoLogin": "true",
		"key_id":    "1",
		"rkey":      "39b392090c635431e86ef76d46f31f40",
		"f":         "http%3A%2F%2Fwww.renren.com%2F222386426",
		"domain":    "renren.com"})

	return
}

func RenRen() *controller.Task {

	task := controller.NewTaskController(
		"人人网",
		"社交网络好友照片",
		[]string{"http://friend.renren.com/GetFriendList.do?curpage=0&id=221940758"},
		4,
		&controller.Login{Action: RenRenLogin})

	c := task.C[0]
	albumlist := task.C[1]
	album := task.C[2]

	c.URLFilters = []*regexp.Regexp{
		//regexp.MustCompile("^http://mat1\\.gtimg\\.com"),
	}

	c.OnRequest(func(r *colly.Request) {
		r.Headers.Set("Host", "friend.renren.com")
		r.Headers.Set("Upgrade-Insecure-Requests", "1")
		r.Headers.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; WOW64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/66.0.3359.170 Safari/537.36")
		r.Headers.Set("Cache-Control", "0")
	})

	c.OnResponse(func(r *colly.Response) {
		//log.Info("seed resp=", string(r.Body))
	})

	c.OnHTML(`#topPage>a[title="下一页"]`,
		func(e *colly.HTMLElement) {
			href := e.Attr("href")
			c.Visit("http://friend.renren.com" + href)
		})

	c.OnHTML(`dd >a[href*="profile.do"]`,
		func(e *colly.HTMLElement) {
			name := e.Text
			href := e.Attr("href")
			re := regexp.MustCompile(`profile\.do\?id=(\d+)`)
			ret := re.FindStringSubmatch(href)
			if len(ret) < 2 {
				log.Warn("capture friends uid failed")
				return
			}

			uid := ret[1]
			log.Infof("capture uid=%s, name=%s", uid, name)

			friend := fmt.Sprintf("http://friend.renren.com/GetFriendList.do?curpage=0&id=%s", uid)
			c.Visit(friend)

			/*
			   link := fmt.Sprintf("http://photo.renren.com/photo/%s/albumlist/v7?offset=0&limit=20", uid)
			   albumlist.Visit(link)
			*/
		})

	albumlist.OnResponse(func(r *colly.Response) {
		if r.StatusCode != 200 {
			return
		}

		var resp string = string(r.Body)
		resp = resp[13+strings.Index(resp, `'albumList': [`):]
		jsonData := resp[:2+strings.Index(resp, `}],`)]

		data := struct {
			AlbumList []struct {
				Cover         string `json: "cover"`
				AlbumName     string `json: "albumName"`
				AlbumId       string `json: "albumId"`
				OwnerId       uint64 `json: "ownerId"`
				SourceControl int    `json: "sourceControl"`
				PhotoCount    int    `json: "photoCount"`
				Type          int    `json: "type"`
			} `json: "albumList"`
		}{}

		err := json.Unmarshal([]byte(jsonData), &data.AlbumList)
		if err != nil {
			log.Warn("albumList unmarshal error: ", err)
			log.Warn(resp)
			return
		}

		for _, v := range data.AlbumList {
			link := fmt.Sprintf("http://photo.renren.com/photo/%d/album-%s/v7", v.OwnerId, v.AlbumId)
			album.Visit(link)
		}
	})

	album.OnResponse(func(r *colly.Response) {
		if r.StatusCode != 200 {
			return
		}

		var resp string = string(r.Body)
		resp = resp[13+strings.Index(resp, `{"photoList":{`):]
		jsonData := resp[:1+strings.Index(resp, `}};`)]

		rb := strings.NewReplacer(
			"\n", "",
		)
		jsonData = rb.Replace(jsonData)

		re := regexp.MustCompile(`'(\w+)' ?:`)
		jsonData = re.ReplaceAllString(jsonData, `"$1":`)

		re = regexp.MustCompile(": ?'([0-9a-zA-Z\u0391-\uFFE5]*)'")
		jsonData = re.ReplaceAllString(jsonData, `:"$1"`)

		log.Infof("json=%s", jsonData)

		data := struct {
			AlbumId   uint64 `json: "albumId"`
			AlbumName string `json: "albumName"`
			PhotoList []struct {
				Height  int    `json: "height"`
				Width   int    `json: "width"`
				URL     string `json: "url"`
				PhotoId string `json: "photoId"`
			} `json: "photoList"`
		}{}

		err := json.Unmarshal([]byte(jsonData), &data)
		if err != nil {
			log.Warn("photoList unmarshal error: ", err)
			return
		}

		for _, v := range data.PhotoList {
			title := data.AlbumName
			//task.CreateTaskItem(title, v.URL)
			log.Infof("got one image %s %s", title, v.URL)
		}

	})

	return task
}
