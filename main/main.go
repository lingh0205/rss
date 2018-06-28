package main

import (
	"log"
	"flag"
	"github.com/rss_subscribe/scrapy"
	"upper.io/db.v3"
)

var (
	configPath    = flag.String("d", "config.json", "Path to get rss.")
)

type Config struct {
	Name string `json:"name"`
	Url string `json:"url"`
	Parse string `json:"parse"`
	Img string `json:"img"`
	Enable bool `json:"enable"`
}

func main(){
	//获取命令行参数
	flag.Parse()

	//加载数据源配置文件
	configs := make([]Config,10)
	jsonParser := scrapy.NewJsonParse()
	e := jsonParser.Load(*configPath, &configs)
	if e != nil {
		log.Fatal("[ERROR]Failed to get config : " + e.Error())
		return
	}

	//数据库初始化
	sess, eor := scrapy.DbInit()
	if eor != nil {
		log.Fatal("[ERROR]Failed to connect db : " + eor.Error())
		return
	}
	defer sess.Close()
	storage := sess.Collection(scrapy.Subscribe)
	if !storage.Exists() {
		log.Fatal("Storage table does not exist")
	}

	helper := scrapy.ScrapyHelper{}

	for _, c := range configs {
		if c.Enable {
			scrapyFunc(&helper, c, &storage)
		}
	}

}

func scrapyFunc(helper *scrapy.ScrapyHelper, config Config, storage *db.Collection) error {
	parser, ex := helper.NewParse(config.Parse)
	if ex != nil {
		log.Println("[ERROR]Failed to get parser : " + ex.Error())
		return ex
	}
	feed := parser.Scrapy(config.Url)
	slice := make([]scrapy.Item, 0)

	for _,item := range feed.Channel.Items{
		if scrapy.IsAlreadySubcribe(*storage, item.Link) {
			// 已订阅
			continue
		}

		if item.Img == "" {
			image, err := scrapy.ScrapyOneImage(item.Link, config.Img, 300, 200)
			if err != nil {
				//获取图片失败
				continue
			}
			item.Img = image.Src
		}

		slice = append(slice, item)
	}

	//通知
	feed.Channel.Items = slice
	return scrapy.Notify(&feed.Channel, *storage)
}