package adaptor

import (
	"encoding/json"
	"img-crawler/src/controller"
	"img-crawler/src/log"
	"regexp"
	"strings"

	"github.com/gocolly/colly"
)

func Ent_qq() *controller.Task {

	task := controller.NewTaskController(
		"qq娱乐明星库",
		"test",
		[]string{"http://ent.qq.com/c/all_star.shtml"},
		1)

	c := task.C[0]

	c.URLFilters = []*regexp.Regexp{
		regexp.MustCompile("^https?://.*\\.qq\\.com/.*"),
		regexp.MustCompile("^http://mat1\\.gtimg\\.com"),
	}

	c.OnHTML(`a[title][href$="index.shtml"]:not([title=''])`,
		func(e *colly.HTMLElement) {

			link := strings.Replace(e.Attr("href"), "index.shtml", "starpicslist.js", 1)
			title := e.Attr("title")
			pageProcess(title, link)
		})

	return task
}

func pageProcess(title, link string) {
	defer func() {
		if r := recover(); r != nil {
			err := r.(error)
			log.Errorf("pageProcess errror %s", err)
		}
	}()

	res := controller.Download(link)
	if res == nil {
		return
	}

	r := strings.NewReplacer(
		"\n", "",
		"\r", "",
		"\t", "",
		"arrPic", `"arrPic"`,
		"nID", `"nID"`,
		"nDataID", `"nDataID"`,
		"nTypeID", `"nTypeID"`,
		"sOriginalImgUrl", `"sOriginalImgUrl"`,
		"sZoomImgUrl", `"sZoomImgUrl"`,
		"sDesc", `"sDesc"`,
	)

	resp := strings.TrimSpace(string(res))
	resp = r.Replace(resp)
	jsonData := resp[strings.Index(resp, "{") : len(resp)-1]

	data := struct {
		ArrPic []struct {
			NID             string `json: "nID"`
			NDataID         string `json: "nDataID"`
			NTypeID         string `json: "nTypeID"`
			SOriginalImgUrl string `json: "sOriginalImgUrl"`
			SZoomImgUrl     string `json: "sZoomImgUrl"`
			SDesc           string `json: "sDesc"`
		} `json: "arrPic"`
	}{}

	err := json.Unmarshal([]byte(jsonData), &data)
	if err != nil {
		log.Warn("unmarshal error: ", err)
		return
	}

	base_url := "http://mat1.gtimg.com/datalib_img/star/"
	for _, img := range data.ArrPic {
		url := base_url + img.SOriginalImgUrl
		log.Infof("got one image %s %s", title, url)
	}
}
