package adaptor

import (
	"img-crawler/src/controller"
	"img-crawler/src/log"
	"regexp"

	"github.com/gocolly/colly"
)

func Ent_qq() *controller.Task {

	task := controller.NewTask(
		"qq娱乐明星库",
		"test",
		[]string{"http://ent.qq.com/c/all_star.shtml"})

	c := task.C
	c.URLFilters = []*regexp.Regexp{
		regexp.MustCompile(`^https?://.*\.qq\.com/.*`),
	}

	detailCollector := c.Clone()

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
	detailCollector.OnHTML(`#disp_right  a[onclick*="goPage"]`,
		func(e *colly.HTMLElement) {

			img := e.ChildAttr("img", "src")
			if img == "" {
				return
			}

			// get one image
			ctx := e.Request.Ctx
			title := ctx.Get("title")
			log.Infof("[%s] got one image, src=%s", title, img)

			// next image
			// onclick = oSerialPicInfoRight.goPage(21157);return false;
			onclikc := e.Attr("onclick")
			reg := regexp.MustCompile(`\([0-9]+\)`)
			next := string(reg.Find([]byte(onclikc)))
			if next == "" {
				return
			}

			picid := next[1 : len(next)-1]
			link := ctx.Get("base_link") + "&picid=" + picid
			detailCollector.Request("GET", link, nil, ctx, nil)
		})

	return task
}
