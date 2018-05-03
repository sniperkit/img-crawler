package adaptor

import (
	"fmt"
	"img-crawler/src/controller"
	"img-crawler/src/log"
	"regexp"
	"strings"

	"github.com/gocolly/colly"
)

func Pic_699() *controller.Task {

	task := controller.NewTaskController(
		"摄图网",
		"http://699pic.com/photo/",
		[]string{"http://699pic.com/photo/"},
        2)

	c, detailCollector := task.C[0], task.C[1]
	c.URLFilters = []*regexp.Regexp{
		//regexp.MustCompile("^https?://.*\\.699pic\\.com/.*"),
	}


	// callback
	// seed html
	c.OnHTML(`div[class="img-show"] a[href]:first-of-type`,
		func(e *colly.HTMLElement) {
			link := e.Attr("href")
			title := e.ChildAttr("img", "alt")

			ctx := colly.NewContext()
			ctx.Put("title", title)
			c.Request("GET", link, nil, ctx, nil)
		})

	// home html
	c.OnHTML(`div[class="list"] >a[href*="tupian"]:first-of-type`,
		func(e *colly.HTMLElement) {
			link := e.Attr("href")
			ctx := e.Request.Ctx

			parseURL(ctx.Get("title"), link)

			detailCollector.Request("GET", link, nil, ctx, nil)

		})

	// photo html
	detailCollector.OnHTML(`div[class="list"]>a[href]`,
		func(e *colly.HTMLElement) {
			link := e.Attr("href")
			ctx := e.Request.Ctx
			parseURL(ctx.Get("title"), link)
		})

	return task
}

func parseURL(title, link string) {

	reg := regexp.MustCompile(`\-[0-9]+\.`)
	picid := string(reg.Find([]byte(link)))
	if picid == "" {
		log.Warnf("parse %s error", link)
	}

	pid := picid[1 : len(picid)-1]
	if len(pid) < 9 {
		pid = strings.Repeat("0", 9-len(pid)) + pid
	}

	img := fmt.Sprintf("http://seopic.699pic.com/photo/%s/%s.jpg_wh1200.jpg", pid[:5], pid[5:])
	log.Infof("[%s] got one image, src=%s", title, img)

	controller.DownloadPic(title, img)
}
