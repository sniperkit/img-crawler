package controller

import (
	"bytes"
	"net"
	"net/http"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/gocolly/colly"
	"github.com/gocolly/colly/extensions"
	"github.com/gocolly/colly/proxy"
	"github.com/gocolly/redisstorage"

	"img-crawler/src/conf"
	"img-crawler/src/log"
)

func CreateCollector() *colly.Collector {

	/* initialize spider */
	c := colly.NewCollector()

	/* HTTP configuration */
	c.WithTransport(&http.Transport{
		Proxy: http.ProxyFromEnvironment,
		DialContext: (&net.Dialer{
			Timeout:   30 * time.Second,
			KeepAlive: 30 * time.Second,
			DualStack: true,
		}).DialContext,
		MaxIdleConns:          65535,
		IdleConnTimeout:       90 * time.Second,
		TLSHandshakeTimeout:   10 * time.Second,
		ExpectContinueTimeout: 30 * time.Second,
	})

	c.Async = true
	c.IgnoreRobotsTxt = true

	/* global LimitRule */
	c.Limit(&colly.LimitRule{
		DomainGlob:  "*",
		Parallelism: conf.Config.Collector.Parallelism,
		RandomDelay: time.Duration(conf.Config.Collector.Random_delay) * time.Millisecond,
	})

	/* redis storage backend */
	if conf.Config.Redis_backend.Switch {
		storage := &redisstorage.Storage{
			Address:  conf.Config.Redis_backend.URL,
			Password: "",
			DB:       conf.Config.Redis_backend.DB,
			Prefix:   conf.Config.Redis_backend.Prefix,
		}

		if err := c.SetStorage(storage); err != nil {
			log.Fatal(err)
		}

		// delete previous data from storage
		if err := storage.ClearURL(); err != nil {
			log.Fatal(err)
		}
	}

	/* Random UA & Refer */
	extensions.RandomUserAgent(c)
	extensions.Referrer(c)

	/* cookiejar */
	/*
			j, err := cookiejar.New(&cookiejar.Options{Filename: "cookie.db"})
			if err == nil {
				c.SetCookieJar(j)
			} else {
		        log.Fatal(err)
		    }
	*/

	/* Set proxy */
	pxy := conf.Config.Collector.Proxy
	if len(pxy) > 0 {
		rp, err := proxy.RoundRobinProxySwitcher(pxy...)
		if err != nil {
			log.Fatal(err)
		}

		c.SetProxyFunc(rp)
	}

	c.MaxDepth = conf.Config.Collector.Max_depth
	c.CacheDir = conf.Config.Collector.Cache_dir

	return c
}

func HTMLPreview(resp *colly.Response, goquerySelector string, f colly.HTMLCallback) error {
	if !strings.Contains(strings.ToLower(resp.Headers.Get("Content-Type")), "html") {
		return nil
	}
	doc, err := goquery.NewDocumentFromReader(bytes.NewBuffer(resp.Body))
	if err != nil {
		return err
	}

	doc.Find(goquerySelector).Each(func(i int, s *goquery.Selection) {
		for _, n := range s.Nodes {
			e := colly.NewHTMLElementFromSelectionNode(resp, s, n)
			f(e)
		}
	})
	return nil
}
