package adaptor

import (
	"img-crawler/src/controller"
	"img-crawler/src/log"
	"regexp"

	"github.com/gocolly/colly"
)

func Ent_qq() *controller.Task {

	task := controller.NewTaskController(
		"qq娱乐明星库",
		"test",
		[]string{"http://ent.qq.com/c/all_star.shtml"},
        2)

	c, detailCollector := task.C[0], task.C[1]

	c.URLFilters = []*regexp.Regexp{
		regexp.MustCompile("^https?://.*\\.qq\\.com/.*"),
		regexp.MustCompile("^http://mat1\\.gtimg\\.com"),
	}

	// callback
	// seed html
	c.OnHTML(`a[title][href$="index.shtml"]:not([title=''])`,
		func(e *colly.HTMLElement) {
			link := e.Attr("href")
			title := e.Attr("title")

			ctx := colly.NewContext()
			ctx.Put("title", title)
			detailCollector.Request("GET", link, nil, ctx, nil)
		})

	// star home html
	detailCollector.OnHTML(`#star_face > a[href]`,
		func(e *colly.HTMLElement) {

			link := "http://datalib.ent.qq.com"
			link += e.Attr("href")
			ctx := e.Request.Ctx
			ctx.Put("base_link", link)
			detailCollector.Request("GET", link, nil, ctx, nil)

		})

	// photo html
	detailCollector.OnHTML(`div[id="disp_right"]`,
		func(e *colly.HTMLElement) {

			ctx := e.Request.Ctx
			title := ctx.Get("title")

			log.Infof("chenqi %s", e.Response.Body)
			e.ForEach(`a[href]`, func(_ int, el *colly.HTMLElement) {
				log.Infof("test %s", title)

				img := el.ChildAttr("img", "src")
				if img == "" {
					return
				}

				href := el.Attr("href")
				var picid string
				if onclikc := el.Attr("onclick"); onclikc != "" {
					reg := regexp.MustCompile(`\([0-9]+\)`)
					if next := string(reg.Find([]byte(onclikc))); next != "" {
						picid = next[1 : len(next)-1]
					}
				}

				if href == "#" && picid == "" {
					return
				}

				// get one image
				log.Infof("[%s] got one image, src=%s", title, img)
				link := ctx.Get("base_link") + "&picid=" + picid
				detailCollector.Request("GET", link, nil, ctx, nil)

			})

			return
		})

	return task
}
