package controller

import (
	"net"
	"net/http"
	"time"

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
		ExpectContinueTimeout: 2 * time.Second,
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
		if err := storage.Clear(); err != nil {
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
