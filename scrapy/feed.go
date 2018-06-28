package scrapy

import "encoding/xml"

type Feed struct {
	XMLName xml.Name `xml:"rss"`
	Channel Channel `xml:"channel"`
}

type Channel struct {
	Title string `xml:"title"`
	Link string `xml:"link"`
	Description string `xml:"description"`
	LastBuildDate string `xml:"lastBuildDate"`
	Language string `xml:"language"`
	Items []Item `xml:"item"`
}

type Item struct {
	Title string `xml:"title"`
	Link string `xml:"link"`
	PubDate string `xml:"pubDate"`
	Img string `xml:"src,addr"`
	Description string `xml:"description"`
	Creator string `xml:"creator"`
	Category []string `xml:"category"`
}


