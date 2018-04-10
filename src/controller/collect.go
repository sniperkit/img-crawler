package controller

import (
	"net"
	"net/http"
	"time"

	"github.com/gocolly/colly"
	"github.com/gocolly/colly/extensions"
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
		MaxIdleConns:          100,
		IdleConnTimeout:       90 * time.Second,
		TLSHandshakeTimeout:   10 * time.Second,
		ExpectContinueTimeout: 1 * time.Second,
	})

	c.Async = true
	c.Limit(&colly.LimitRule{
		Parallelism: conf.Config.Collector.Parallelism,
		RandomDelay: time.Duration(conf.Config.Collector.Random_delay) * time.Second,
	})

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
