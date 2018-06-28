// Copyright (c) 2017-2018 LinGH, http://blog.lingh.vip
// Use of this source code to transfer data to json
// Or transfer json to map, struct
package scrapy

import (
	"io/ioutil"
	"encoding/json"
)

type JsonParse struct {}

func NewJsonParse() *JsonParse {
	return &JsonParse{}
}

func (parse *JsonParse) Load(filename string, v interface{}) error {
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		return err
	}

	err = json.Unmarshal(data, v)
	if err != nil {
		return err
	}
	return nil
}

func ParseMap(data string, obj map[string]interface{}) error {
	return json.Unmarshal([]byte(data), obj)
}

func ParseStruct(data string, obj interface{}) (error) {
	return json.Unmarshal([]byte(data), obj)
}
