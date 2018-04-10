# img-crawler

图片爬虫。

## 安装
从gitlab下载代码后放到$GOPATH/src/路径下,
即项目实际路径为：$GOPATH/src/img-crawler
到项目目录下执行下面的命令完成安装：

```
go get github.com/golang/dep/cmd/dep    # 安装dep工具
go get golang.org/x/tools/cmd/goimports # 安装goimports工具
make install  # 使用dep工具安装依赖的库
```
以上命令只有第一次部署需要执行，之后修改代码后只要执行下面的make编译命令即可：
```
make          # 编译go源码
```

## 代码目录
- *bin/*      # 可执行文件目录
- *cmd/*      # main包
- *src/*      # go源码
- *vendor/*   # 依赖的第三方库

## 配置
配置文件路径：
```
src/conf/crawler.conf
```


## 运行
```
bin/spider
```

## 参考
- [colly github](https://github.com/gocolly/colly)
- [colly Doc](http://go-colly.org/docs/)
- [colly API](https://godoc.org/github.com/gocolly/colly)
- [W3school](http://www.w3school.com.cn/h.asp)
- [CSS Selector](http://www.w3school.com.cn/cssref/css_selectors.asp)
- [goquery](https://godoc.org/github.com/PuerkitoBio/goquery)


[回到顶部](#readme)
