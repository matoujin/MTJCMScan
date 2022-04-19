package main

import (
	"bufio"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"myCMStest/core"
	"myCMStest/jsonRead"
	"os"
	"strings"
	"time"
)

func banner() {
	banner :=
		`
 __  __ _____   _     ____ __  __ ____                  
|  \/  |_   _| | |  / ___|  \/  / ___|  ___ __ _ _ __  
| |\/| | | |_  | | | |   | |\/| \___ \ / __/ _  | '_ \ 
| |  | | | | |_| | | |___| |  | |___) | (_| (_| | | | |
|_|  |_| |_|\___/  |\____|_|  |_|____/ \___\__,_|_| |_|

Usage of MTJCmscan:
  -h string
        set one host, 127.0.0.1
  -hosts string
        set one file name ,hosts divided by \n
  -json string
        json of fingerprint
`
	print(banner)
}

func Use() {

}

func main() {
	banner()
	var (
		Info    core.ArgsInfo
		file    io.ReadCloser
		err     error
		hosts   []string
		results []string
		//output为生成日志名
		outputFileName string
	)

	t := time.Now()
	Info.Flag()
	if Info.Host != "" {
		//单个扫描
		fmt.Println("GET host: " + Info.Host)
		//命名时去除http的干扰
		outputFileName = "--" + Info.Host[strings.Index(Info.Host, "/")+2:] + ".txt"
		//Reader只是一个读取字符串
		file = ioutil.NopCloser(strings.NewReader(Info.Host))
	} else {
		//多个扫描:文件
		fmt.Println("GET hosts : " + Info.Hosts)
		outputFileName = "--" + Info.Hosts + ".txt"
		//打开文件
		file, err = os.Open(Info.Hosts)
		if err != nil {
			log.Fatalf("err ./Cannot open the hosts file %s: %s\n", Info.Hosts, err)
		}
	}
	//延迟处理
	defer file.Close()
	//读取流
	scan := bufio.NewScanner(file)

	//逐行读取
	for scan.Scan() {
		log.Printf("LOAD host: %s\n", scan.Text())
		hosts = append(hosts, scan.Text())
	}

	if Info.CmsJson != "" {
		//跳转到对json进行解析
		cmsList := jsonRead.ReadJson(Info.CmsJson)
		//fmt.Println(cmsList, cmsSortList)
		log.Println("Successfully parsed the json file: " + Info.CmsJson)
		log.Println("Start testing")
		//进入扫描
		results = jsonRead.HostWorker(hosts, cmsList)
	}

	//创建输出文件
	outfile, err := os.Create(outputFileName)
	if err != nil {
		log.Fatalf("err ./Error creating output file: %s", err)
		return
	}
	//延迟处理
	defer outfile.Close()
	//return resultsList []string
	for _, result := range results {
		fmt.Println(result)
		outfile.WriteString(result + "\n")
	}
	//计算运行时间
	tim := time.Since(t)
	fmt.Printf("time is: %s", tim)
}
