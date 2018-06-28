package scrapy

import (
	"net/http"
	"io/ioutil"
	"github.com/gocolly/colly"
	"strconv"
	"encoding/xml"
	"log"
)

func scrapy2Byte(url string) ([]byte, error) {

	resp, err := http.Get(url)

	defer resp.Body.Close()

	if err!= nil{
		log.Println("[ERROR]Failed to fetch html pages with url : " + url)
		return nil, err
	}

	if resp.StatusCode == 200 {
		body,_ := ioutil.ReadAll(resp.Body)
		return body, nil
	}

	return nil, err
}

type Scrapy interface  {
	Scrapy(url string) *Feed
}

type HtmlScrapy struct {}

type RssScrapy struct {}

func (HtmlScrapy) Scrapy(url string)  *Feed {
	collector := colly.NewCollector()
	feed := Feed{}
	slice := make([]Item, 0)
	collector.OnHTML("head", func(e *colly.HTMLElement) {
		feed.Channel.Title = e.ChildText("title")
		feed.Channel.Link = url
	})

	collector.OnHTML("#archive", func(e *colly.HTMLElement) {
		e.ForEach("div.post.floated-thumb", func(_ int, e *colly.HTMLElement) {
			item := Item{
				Title: e.ChildAttr("div.post-meta > p:nth-child(1) > a.meta-title", "title"),
				Link: e.ChildAttr("div.post-meta > p:nth-child(1) > a.meta-title", "href"),
				PubDate: e.ChildText("div.post-meta > p:nth-child(1)"),
				Description: e.ChildText("div.post-meta > span > p"),
				Img: e.ChildAttr("div.post-thumb > a > img", "src"),
				Creator: "",
				Category: []string{e.ChildText("div.post-meta > p:nth-child(1) > a:nth-child(3)")},
			}
			slice = append(slice, item)
		})
	})

	collector.OnRequest(func(r *colly.Request) {
		log.Println("Visiting Protal URL : ", r.URL)
	})

	err := collector.Visit(url)

	if err != nil {
		log.Println("[ERROR] Failed to Fetch data with URL : " + url)
		log.Println(err)
		return nil
	}
	feed.Channel.Items = slice
	return &feed
}

func (RssScrapy) Scrapy(url string) *Feed {
	content, _ := scrapy2Byte(url)
	feed := Feed{}
	xml.Unmarshal(content, &feed)
	return &feed
}

func ScrapyOneImage(url string, imgPath string, cWidth int, cHeight int) (Image, error) {
	collector := colly.NewCollector()
	image := Image{}
	isFinish := false

	collector.OnHTML("img", func(e *colly.HTMLElement) {
		if isFinish {
			return
		}
		if imgPath == "" {
			imgPath = "src"
		}
		src := e.Attr(imgPath)
		width,_ := strconv.Atoi(e.Attr("width"))
		height,_ := strconv.Atoi(e.Attr("height"))
		if width >= cWidth && height >= cHeight {
			image.Src = src
			image.Width = width
			image.Height = height
			isFinish = true
		}else {
			log.Println("Skip url : " + src)
		}
	})

	collector.OnRequest(func(r *colly.Request) {
		log.Println("Visiting Protal URL : ", r.URL)
	})

	err := collector.Visit(url)

	if err != nil {
		log.Println("[ERROR] Failed to Fetch data with URL : " + url)
		log.Println(err)
		return image, err
	}

	return image, nil
}
