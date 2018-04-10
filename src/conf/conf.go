package conf

import (
	"encoding/json"
	"fmt"
	"os"
)

type BaseConfig struct {
	MySQL struct {
		Master string   `json:"master"`
		Slaves []string `json:"slaves"`
	} `json:"mysql"`
	Redis_url string `json: "redis_url"`
	Log_path  string `json: "log_path"`
	Img_dir   string `json: "img_dir"`
	Collector struct {
		Max_depth    int    `json: "max_depth"`
		Cache_dir    string `json: "cache_dir"`
		Parallelism  int    `json: "parallelism"`
		Random_delay int    `json: "random_delay"`
	} `json: "collector"`
}

func load_json(filepath string, f interface{}) {

	fl, err := os.Open(filepath)
	if err != nil {
		fmt.Println("%s Open failed: %s.", filepath, err)
		return
	}
	defer fl.Close()

	fi, err := os.Stat(filepath)
	buf := make([]byte, fi.Size())

	for {
		n, _ := fl.Read(buf)
		if 0 == n {
			break
		}
	}

	err = json.Unmarshal(buf, &f)
	return
}

var (
	Config = &BaseConfig{}
)

func init() {
	gopath := os.Getenv("GOPATH")
	//	ENV := os.Getenv("DEPLOY_ENV")
	load_json(gopath+"/src/img-crawler/src/conf/crawler.conf", Config)
}
