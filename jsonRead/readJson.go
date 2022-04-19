package jsonRead

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
)

//json建立结构体
type CmsFeature struct {
	RuleID         string `json:"rule_id"`
	Level          string `json:"level"`
	Softhard       string `json:"softhard"`
	Product        string `json:"product"`
	Company        string `json:"company"`
	Category       string `json:"category"`
	ParentCategory string `json:"parent_category"`
	Rules          [][]struct {
		Match   string `json:"match"`
		Content string `json:"content"`
	} `json:"rules"`
}

//读取文件
func ReadJson(filename string) []CmsFeature {
	jsonFile, err := ioutil.ReadFile(filename)
	if err != nil {
		panic(err.Error())
	}
	fmt.Println("Have Opened fofa.json")
	var CMSList []CmsFeature
	err = json.Unmarshal(jsonFile, &CMSList)
	if err != nil {
		panic(err.Error())
	}
	return CMSList
}
