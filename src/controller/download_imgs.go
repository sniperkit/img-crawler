package controller

import (
	"img-crawler/src/conf"
	"img-crawler/src/log"
	"img-crawler/src/utils"
	"io/ioutil"
	"net"
	"net/http"
	"os"
	"path/filepath"
	"time"
)

func DownloadPic(title, url string) {
	res, err := client.Get(url)
	if err != nil {
		log.Warnf("download get %s error %s", url, err)
		return
	}

	defer res.Body.Close()
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		log.Warnf("download read %s error %s", url, err)
		return
	}

	base := filepath.Join(conf.Config.Img_dir, title)
	err = os.MkdirAll(base, 0666)
	if err != nil {
		log.Warnf("download mkdir %s error %s", title, err)
		return
	}

	filename := utils.GenerateUuidV5(url)
	err = ioutil.WriteFile(filepath.Join(base, title, filename), body, 0666)
	if err != nil {
		log.Warnf("download write %s error %s", url, err)
	}

}

func NewClient() *http.Client {
	tr := &http.Transport{
		Proxy: http.ProxyFromEnvironment,
		DialContext: (&net.Dialer{
			Timeout:   30 * time.Second,
			KeepAlive: 30 * time.Second,
			DualStack: true,
		}).DialContext,
		MaxIdleConns:          128,
		IdleConnTimeout:       90 * time.Second,
		TLSHandshakeTimeout:   10 * time.Second,
		ExpectContinueTimeout: 1 * time.Second,
	}

	client := &http.Client{
		Transport: tr,
	}

	return client
}

var client *http.Client

func init() {
	client = NewClient()
}
