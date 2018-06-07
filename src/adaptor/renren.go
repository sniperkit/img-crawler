package adaptor

import (
	"encoding/json"
	"fmt"
	"img-crawler/src/controller"
	"img-crawler/src/log"
	"img-crawler/src/utils"
	"regexp"
	"strings"

	"github.com/gocolly/colly"
)

func RenRenLogin(c *colly.Collector) (err error) {

	c.OnResponse(func(r *colly.Response) {
		log.Info("logon resp=", string(r.Body))
	})

	//vcode := "http://icode.renren.com/getcode.do?t=login&rnd=Math.random()"
	url := "http://www.renren.com/PLogin.do"
	err = c.Post(url, map[string]string{
		"email":        "chenny7@163.com",
		"autoLogin":    "true",
		"captcha_type": "web_login",
		"origURL":      "http://www.renren.com/home",
		"key_id":       "1",
		"rkey":         "5b625fd2e7bb436f20f4c7905430b4e1",
		"password":     "060bf26c485a659818f9da48415291d4440d15071f40eeb8e4a2948fdcebecaf",
		"f":            "http%3A%2F%2Fwww.renren.com%2F222386426",
		"domain":       "renren.com"})

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
	photolist := task.C[2]

	c.URLFilters = []*regexp.Regexp{
		//regexp.MustCompile("^http://mat1\\.gtimg\\.com"),
	}

	c.OnRequest(func(r *colly.Request) {
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

			link := fmt.Sprintf("http://photo.renren.com/photo/%s/albumlist/v7?offset=0&limit=20", uid)
			albumlist.Visit(link)
		})

	albumlist.OnRequest(func(r *colly.Request) {
		r.Headers.Set("Pragma", "no-cache")
		r.Headers.Set("Cache-Control", "no-cache")
		r.Headers.Set("Connection", "keep-alive")
		r.Headers.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/webp,image/apng,*/*;q=0.8")
		r.Headers.Set("Accept-Encoding", "gzip, deflate")
		r.Headers.Set("Accept-Language", "zh-CN,zh;q=0.9,en;q=0.8")
		r.Headers.Set("Upgrade-Insecure-Requests", "1")
		r.Headers.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; WOW64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/66.0.3359.170 Safari/537.36")
	})

	albumlist.OnResponse(func(r *colly.Response) {
		if r.StatusCode != 200 {
			return
		}

		var (
			reader []byte
			err    error
			resp   string
		)

		switch r.Headers.Get("Content-Encoding") {
		case "gzip":
			reader, err = utils.ParseGzip(r.Body)
			if err != nil {
				return
			}
		default:
			reader = r.Body
		}

		resp = string(reader)
		// skip if there is no permission or album is empty
		if strings.Contains(resp, `您没有操作本资源的权限`) {
			return
		} else if strings.Contains(resp, `'albumCount': 0`) {
			return
		} else if strings.Contains(resp, `抱歉，出错了`) {
			task.Retry(r.Request, 3)
			return
		}

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

		err = json.Unmarshal([]byte(jsonData), &data.AlbumList)
		if err != nil {
			log.Warnf("albumList unmarshal %s error %s: %s",
				err, r.Request.URL.String(), resp)
			return
		}

		for _, v := range data.AlbumList {
			link := fmt.Sprintf("http://photo.renren.com/photo/%d/album-%s/v7", v.OwnerId, v.AlbumId)
			fmt.Println(link)
			photolist.Visit(link)
		}
	})

	photolist.OnRequest(func(r *colly.Request) {
		r.Headers.Set("Upgrade-Insecure-Requests", "1")
		r.Headers.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; WOW64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/66.0.3359.170 Safari/537.36")
		r.Headers.Set("Cache-Control", "0")
	})

	photolist.OnResponse(func(r *colly.Response) {
		if r.StatusCode != 200 {
			return
		}

		var resp string = string(r.Body)
		// skip if there is no permission
		if strings.Contains(resp, `您没有操作本资源的权限`) {
			return
		} else if strings.Contains(resp, `该相册设置了密码保护，输入密码后才允许访问`) {
			return
		} else if strings.Contains(resp, `抱歉，出错了`) {
			task.Retry(r.Request, 3)
			return
		}

		controller.HTMLPreview(r,
			`.fd-nav-item > a[href $="profile"]`,
			func(e *colly.HTMLElement) {
				e.Response.Ctx.Put("name", e.Text)
			})

		resp = resp[13+strings.Index(resp, `{"photoList":{`):]
		jsonData := resp[:1+strings.Index(resp, `}};`)]

		rb := strings.NewReplacer(
			"\n", "",
		)
		jsonData = rb.Replace(jsonData)

		re := regexp.MustCompile(`'(\w+)' ?:`)
		jsonData = re.ReplaceAllString(jsonData, `"$1":`)

		re = regexp.MustCompile(": ?'(.*?)'")
		jsonData = re.ReplaceAllString(jsonData, `:"$1"`)

		reg := regexp.MustCompile(`photo\/(\d+)\/album-`)
		ret := reg.FindStringSubmatch(r.Request.URL.String())
		uid := "0"
		if len(ret) < 2 {
			log.Warn("capture photo uid failed")
		}
		uid = ret[1]

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
			log.Warnf("photoList unmarshal %s error %s: %s",
				err, r.Request.URL.String(), resp)
			return
		}

		for _, v := range data.PhotoList {
			name := uid + "_" + r.Ctx.Get("name")
			desc := data.AlbumName
			log.Infof("got one image %s %s %s", name, desc, v.URL)
			status := controller.Download_INIT
			imgContent := controller.Download(v.URL)
			if imgContent == nil {
				status = controller.Download_DownFAIL
				continue
			}
			digest, filepath, err := controller.Save(name, desc, imgContent)
			if err == nil {
				status = controller.Download_NORMAL
			} else {
				status = controller.Download_SAVEFAIL
			}
			task.CreateTaskItem(name, v.URL, desc, digest, filepath, status)
		}

	})

	return task
}
