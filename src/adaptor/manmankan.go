package adaptor

import (
	//"encoding/json"
	"img-crawler/src/controller"
	"img-crawler/src/log"
	"img-crawler/src/utils"
	"regexp"
	//"strings"
	"github.com/gocolly/colly"
)

func Manmankan() *controller.Task {
	download_pic := false

	task := controller.NewTaskController(
		"manmankan",
		"漫漫看明星库",
		[]string{
			"http://www.manmankan.com/dy2013/mingxing/yanyuan/neidi/",
			"http://www.manmankan.com/dy2013/mingxing/yanyuan/xianggang/",
			"http://www.manmankan.com/dy2013/mingxing/yanyuan/taiwan/",
			"http://www.manmankan.com/dy2013/mingxing/yanyuan/oumei/",
			"http://www.manmankan.com/dy2013/mingxing/yanyuan/riben/",
			"http://www.manmankan.com/dy2013/mingxing/yanyuan/hanguo/",
		},
		1,
		download_pic,
		nil)

	if download_pic {
		return task
	}

	c := task.C[0]

	c.URLFilters = []*regexp.Regexp{
		//regexp.MustCompile("^http://mat1\\.gtimg\\.com"),
	}

	c.OnHTML(`div>a[title][class="show"]`,
		func(e *colly.HTMLElement) {

			title := e.Attr("title")
            title = utils.ConvertToString(title, "gbk", "utf-8")
            log.Infof("got one name: %s ", title)
		})

	return task
}
