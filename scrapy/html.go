package scrapy

type Image struct {
	Src string `json:"src"`
	Height int `json:"height"`
	Width int `json:width`
	Alt string `json:"alt"`
}
