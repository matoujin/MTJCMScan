package jsonRead

import (
	"fmt"
	"myCMStest/reqHost"
	"net"
	"net/url"
	"regexp"
	"strconv"
	"strings"
	"sync"
)

func HostWorker(hosts []string, cmsList []CmsFeature) []string {
	hostsChan := make(chan string)
	//预留的result
	resultChan := make(chan string)
	var resultList []string

	//等待组:保证在并发环境完成指定任务
	var wg sync.WaitGroup
	for i := 0; i < len(hosts); i++ {
		//并发处理
		wg.Add(1)
		go cmsWorker(&wg, hostsChan, cmsList, resultChan)
	}
	//遍历添加到并发任务中
	for _, host := range hosts {
		hostsChan <- host
	}

	for i := 0; i < len(hosts); i++ {
		//跟据host数量返回最后结果数
		result := <-resultChan
		resultList = append(resultList, result)
	}

	//获得结果后，结束并发任务
	close(hostsChan)
	close(resultChan)
	return resultList
}

func cmsWorker(wg *sync.WaitGroup, hosts chan string, cmsList []CmsFeature, resultChan chan string) {

	//遍历多个host
	for host := range hosts {

		//初始化，缓冲区为10的map[键string], 值类型为[]CmsFeature
		//cmsListChan := make(chan map[string][]CmsFeature, 10)
		//提前进行内容提取
		//body，title，banner
		content, err := reqHost.BodyReq(host)
		if err != nil {
			panic(err.Error())
		}
		//header，server
		headersStruct, err := reqHost.HeadersReq(host)
		if err != nil {
			panic(err.Error())
		}
		//等待组:保证在并发环境完成指定任务
		//var wg sync.WaitGroup

		//根据cmsChan数量设定循环次数,cap返回slice容量

		featureWorker(host, cmsList, wg, resultChan, content, headersStruct)

		//resultChan <- fmt.Sprintf("The host: %s has no matching results", host)
	}
}

func featureWorker(host string, cmsList []CmsFeature, wg *sync.WaitGroup, resultChan chan string, content []byte, headersStruct *reqHost.ResponseStruct) {

	var resultstr string
	var flag, globalflag bool

	//返回cms键值对
	for _, cms := range cmsList {
		//resultstr = ""
		//从cms中提取Rules
		rulesList := cms.Rules
		id := cms.RuleID
		id = id
		flag = false
		//对rule遍历
	ArrayRule:
		for _, rule := range rulesList {
			//默认rules内为未扫描到状态
			//flag = false
			//对match和content的提取
			//ArrayKey:
			for _, key := range rule {

				match := strings.Split(key.Match, "_")[0]

				switch match {
				//对body进行特征检查
				case "body":
					/*
						这里就不处理状态码了
						if resp.StatusCode != 200 {
							continue
						}
					*/
					//content, err := reqHost.BodyReq(host)
					//根据json格式，option为关键词时
					if strings.Contains(strings.ToLower(string(content)), strings.ToLower(key.Content)) {
						resultstr = resultstr + fmt.Sprintf("\n"+" %s matches Productname[id=%s] : **%s** -category **%s** -level **%s** -matchfeature [content]:[%s]", host, cms.RuleID, cms.Product, cms.Category, cms.Level, key.Content)
						//resultChan <- fmt.Sprintf(" \n %s matches Productname[id=%s] : %s -category %s -level %s -matchfeature [content]:[%s]", host, cms.RuleID, cms.Product, cms.Category, cms.Level, key.Content)
						//flag设置为扫描到
						flag = true
						//扫描到就结束
					} else {
						//未扫描到
						flag = false
						break
					}

				case "protocol":
					/*
						暂不支持
					*/
					flag = false
					break

				case "title":
					titleRe := regexp.MustCompile("(?im)<\\s*title.*>(.*?)<\\s*/\\s*title>")
					//FindStringSubMatch 返回的第一位是字符串本身
					title := titleRe.FindStringSubmatch(string(content))
					if len(title) > 0 {
						if strings.Contains(strings.ToLower(title[1]), strings.ToLower(key.Content)) {
							resultstr = resultstr + fmt.Sprintf("\n"+" %s matches Productname[id=%s] : **%s** -category **%s** -level **%s** -matchfeature [title]:[%s]", host, cms.RuleID, cms.Product, cms.Category, cms.Level, key.Content)
							//resultChan <- fmt.Sprintf(" %s matches Productname[id=%s] : %s -category %s -level %s -matchfeature [title]:[%s]", host, cms.RuleID, cms.Product, cms.Category, cms.Level, key.Match)
							//flag设置为扫描到
							flag = true
						} else {
							flag = false
							break
						}
					}

				case "header":
					//处理当content格式为 server：的情况
					if headersStruct.Headers.Get(strings.Split(key.Content, ":")[0]) == strings.ToLower(key.Content) {
						resultstr = resultstr + fmt.Sprintf("\n"+" %s matches Productname[id=%s]: **%s** -category **%s** -level **%s** -matchfeature [header]:[%s]", host, cms.RuleID, cms.Product, cms.Category, cms.Level, key.Content)
						//resultChan <- fmt.Sprintf(" %s matches Productname[id=%s]: %s -category %s -level %s -matchfeature [header]:[%s]", host, cms.RuleID, cms.Product, cms.Category, cms.Level, key.Match)
						flag = true
					} else {
						flag = false
						break
					}

				case "server":
					if strings.Contains(strings.ToLower(headersStruct.Headers.Get("Server")), strings.ToLower(key.Content)) {
						flag = true
						resultstr = resultstr + fmt.Sprintf("\n"+" %s matches Productname[id=%s] : **%s** -category **%s** -level **%s** -matchfeature [server]:[%s]", host, cms.RuleID, cms.Product, cms.Category, cms.Level, key.Content)
						//resultChan <- fmt.Sprintf(" %s matches Productname[id=%s] : %s -category %s -level %s -matchfeature [server]:[%s]", host, cms.RuleID, cms.Product, cms.Category, cms.Level, key.Match)
					} else {
						flag = false
						break
					}

				case "banner":
					compileRegex := regexp.MustCompile("(?im)<\\s*banner.*>(.*?)<\\s*/\\s*banner>")
					banner := compileRegex.FindStringSubmatch(string(content))
					if len(banner) > 0 {
						if strings.Contains(strings.ToLower(banner[1]), strings.ToLower(key.Content)) {
							resultstr = resultstr + fmt.Sprintf("\n"+" %s matches Productname[id=%s] : **%s** -category **%s** -level **%s** -matchfeature [banner]:[%s]", host, cms.RuleID, cms.Product, cms.Category, cms.Level, key.Content)
							//resultChan <- fmt.Sprintf(" %s matches Productname[id=%s] : %s -category %s -level %s -matchfeature [banner]:[%s]", host, cms.RuleID, cms.Product, cms.Category, cms.Level, key.Match)
							//flag设置为扫描到
							flag = true
						} else {
							flag = false
							break
						}
					} else {
						flag = false
						break
					}

				case "port":
					u, _ := url.Parse(host)
					host := u.Host
					address := net.ParseIP(host)
					if address == nil {
						flag = false
						break
					}
					ho := strings.Split(host, ":")
					if len(ho) > 1 {
						port := ho[1]
						if port == key.Content {
							flag = true
							resultstr = resultstr + fmt.Sprintf("\n"+" %s matches Productname[id=%s] : **%s** -category **%s** -level **%s** -matchfeature [Port]:[%s]", host, cms.RuleID, cms.Product, cms.Category, cms.Level, key.Content)
							//resultChan <- fmt.Sprintf(" %s matches Productname[id=%s] : %s -category %s -level %s -matchfeature [Port]:[%s]", host, cms.RuleID, cms.Product, cms.Category, cms.Level, key.Match)
						} else {
							flag = false
							break
						}
					} else {
						flag = false
						break
					}

				}

				//当flag为true表示满足and条件其一，则继续循环；不满足则直接换key-value
				if flag == false {
					break ArrayRule
				} else if flag == true {
					globalflag = true

				}

			}

		}
		nowid, _ := strconv.Atoi(cms.RuleID)
		if nowid > 759789 {
			if globalflag == true {
				resultChan <- resultstr
				wg.Done()
			} else {
				resultChan <- fmt.Sprintf("The host: %s has no matching results", host)
			}

		}

	}
	return

}
