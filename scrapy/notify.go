package scrapy

import (
	"net/http"
	"strings"
	"fmt"
	"io/ioutil"
	"log"
	"upper.io/db.v3"
	"encoding/json"
)

const NOTIFY_URL string = "https://oapi.dingtalk.com/robot/send?access_token=431c7b9470d4d20b88560f63c27e27d0b8475b308df2ddcd0f6fd1f3fe335bbf"

type Response struct{
	Errcode int `json:"errcode"`
	Errmsg string `json:"errmsg"`
}

func Notify(channel *Channel, storage db.Collection) error {
	for i := 0; i < len(channel.Items) ; i = i + 5 {
		slice := channel.Items[i:i+5]
		title, notifyText := buildMarkDownText(channel, slice)

		resp, err := sendHtmlNotify(title, notifyText)
		if err != nil || resp.StatusCode != 200 {
			//消息发送失败
			return http.ErrUseLastResponse
		}

		bytes, _ := ioutil.ReadAll(resp.Body)

		stb := &Response{}
		err = json.Unmarshal(bytes, &stb)

		if stb.Errmsg != "ok" || err != nil {
			//消息发送失败
			log.Printf("reponse code : %d, response msg : %s.", resp.StatusCode, string(bytes))
			continue
		}

		for _, item := range slice{
			//入库
			err := InsertRecord(storage, item)
			if err != nil {
				log.Println("[ERROR]Failed to insert item : " + err.Error())
			}
		}
	}
	return nil
}

func sendHtmlNotify(title string, notifyText string) (resp *http.Response, err error) {
	msg := fmt.Sprintf("{'%s':{'%s':'%s','%s':'%s'},'%s':'%s'}", "markdown", "title", title, "text", notifyText, "msgtype", "markdown")
	fmt.Println(msg)
	return http.Post(NOTIFY_URL, "application/json", strings.NewReader(msg))
}

type Message struct {
	Title string `json:"title"`
	Text string `json:"text"`
}

func buildMarkDownText(channel *Channel, slice []Item) (string, string) {
	text := ""
	for _, item := range slice {
		item.Title = strings.Replace(item.Title,"'", "",-1)
		item.PubDate = strings.Replace(item.PubDate,"'", "",-1)
		channel.Title = strings.Replace(channel.Title,"'", "",-1)
		text = fmt.Sprintf("%s  ### %s | %s \n #### [%s](%s) \n ###### 发布时间：%s \n", text, channel.Title, channel.Description, item.Title, item.Link, item.PubDate)
		if item.Img != "" {
			text = fmt.Sprintf("%s ![%s](%s) \n", text, item.Title, item.Img)
		}
	}
	return channel.Title, text
}
