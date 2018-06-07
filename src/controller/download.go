package controller

import (
	"crypto/md5"
	"encoding/hex"
	"img-crawler/src/conf"
	"img-crawler/src/log"
	"io/ioutil"
	"net"
	"net/http"
	"os"
	"strings"
	//    "github.com/gocolly/colly/queue"
	"path/filepath"
	"time"
)

func Save(dir, desc string, content []byte) (string, string, error) {

	base := filepath.Join(conf.Config.Img_dir, dir)
	err := os.MkdirAll(base, 0755)
	if err != nil {
		log.Warnf("Save mkdir %s error %s", dir, err)
		return "", "", err
	}

	h := md5.New()
	h.Write(content)
	digest := hex.EncodeToString(h.Sum(nil))
	filename := digest
	if len(desc) > 0 {
		filename = desc + "_" + digest
	}

	filename = filepath.Join(base, filename)
	err = ioutil.WriteFile(filename, content, 0744)
	if err != nil {
		log.Warnf("Save write %s error %s", dir, err)
		return "", "", err
	}
	return digest, filename, nil
}

func Download(url string) []byte {

	res, err := client.Get(url)
	if err != nil {
		log.Warnf("download get %s error %s", url, err)
		return nil
	}

	defer res.Body.Close()

	if res.StatusCode != 200 {
		log.Warnf("download bad res, status=%s", url, res.StatusCode)
		return nil
	}

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		log.Warnf("download read %s error %s", url, err)
		return nil
	}

    ct := res.Header.Get("Content-Type")
    if !strings.Contains(ct, "image") {
        log.Warnf("%s Content-Type not contains image", url)
    }

	return body
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
		ExpectContinueTimeout: 30 * time.Second,
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
