package controller

import (
	"crypto/md5"
	"encoding/hex"
	"img-crawler/src/conf"
    "img-crawler/src/log"
	"os"
	"strings"
	//	"github.com/gocolly/colly/queue"
	"path/filepath"
	"regexp"
	"github.com/gocolly/colly"
)

func genFilename(dir, desc, suffix, digest string) string {

	filename := digest
	if len(desc) > 0 {
		filename = desc + "_" + filename
	}
	if len(suffix) > 0 {
		filename += "." + suffix
	}

	filename = strings.Replace(filename, "/", "", -1)
	filename = filepath.Join(dir, filename)
	return filename
}

func Download(c *colly.Collector) {

	c.OnResponse(func(r *colly.Response) {

		var (
			digest   string
			suffix   string
			filename string
		)
		status := Download_INIT
		url := r.Request.URL.String()
		name := r.Ctx.Get("name")
		desc := r.Ctx.Get("desc")
        task := r.Ctx.GetAny("task").(*Task)

        log.Infof("Download %s", url)

		// download failed
		if r.StatusCode != 200 {
			status = Download_DownFAIL
            log.Errorf("download get failed %d %s", url, r.StatusCode)
            return
		}

		// image suffix
		ct := r.Headers.Get("Content-Type")
		reg := regexp.MustCompile(`image/(\w+)`)
		ret := reg.FindStringSubmatch(ct)
		if len(ret) >= 2 {
			suffix = ret[1]
		}

		// md5sum
		h := md5.New()
		h.Write(r.Body)
		digest = hex.EncodeToString(h.Sum(nil))

		// mkdir
		base := filepath.Join(conf.Config.Img_dir, name)
		if _, err := os.Stat(base); err != nil {
			if err := os.MkdirAll(base, 0750); err != nil {
                log.Errorf("download mkdir failed %s %s", base, err)
                return
			}
		}

		// write to disk
		filename = genFilename(base, desc, suffix, digest)
		if err := r.Save(filename); err != nil {
			status = Download_SAVEFAIL
		} else {
			status = Download_NORMAL
		}

		// insert into mysql
		task.UpdateTaskItem(name, url, desc, digest, filename, status)

	})

}
