package controller

import (
	"net"
	"net/http"
	"time"

	"github.com/gocolly/colly"
	"github.com/gocolly/colly/extensions"
	"github.com/gocolly/redisstorage"
	"github.com/gocolly/colly/proxy"

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

    /* global LimitRule */
	c.Limit(&colly.LimitRule{
        DomainGlob: "*",
		Parallelism: conf.Config.Collector.Parallelism,
		RandomDelay: time.Duration(conf.Config.Collector.Random_delay) * time.Millisecond,
	})

    /* redis storage backend */
	if conf.Config.Redis_url.Switch {
        storage := &redisstorage.Storage{
            Address:  conf.Config.Redis_url.URL,
            Password: "",
            DB:       conf.Config.Redis_url.DB,
            Prefix:   conf.Config.Redis_url.Prefix,
        }
		err := c.SetStorage(storage)
		if err != nil {
		    panic(err)
		}
	}

	/* Random UA & Refer */
	extensions.RandomUserAgent(c)
	extensions.Referrer(c)

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
