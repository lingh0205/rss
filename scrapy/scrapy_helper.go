package scrapy

type ScrapyHelper struct {
	rssScrapy RssScrapy
	htmlScrapy HtmlScrapy
}

var helper = ScrapyHelper{
	rssScrapy: RssScrapy{},
	htmlScrapy: HtmlScrapy{},
}

type ScrapyError struct {
	msg string
}

func (sh ScrapyError) Error() string {
	return "Failed to get parser with name : " + sh.msg
}

func (ScrapyHelper) NewParse(source string) (Scrapy, error){
	switch source {
		case "rss":
			return helper.rssScrapy, nil
		case "html":
			return helper.htmlScrapy, nil
		default:
			return nil, ScrapyError{msg: source}
	}
}

