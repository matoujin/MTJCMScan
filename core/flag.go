package core

import (
	"flag"
	"log"
)

type ArgsInfo struct {
	Host    string
	Hosts   string
	CmsJson string
}

func (Info *ArgsInfo) Flag() {

	flag.StringVar(&Info.Host, "h", "", "set one host, 127.0.0.1")
	flag.StringVar(&Info.Hosts, "hosts", "", "set one file name ,hosts divided by \\n")
	flag.StringVar(&Info.CmsJson, "json", "", "json of fingerprint")
	flag.Parse()

	if Info.Host == "" && Info.Hosts == "" {
		log.Fatalln("err:./no host parameter")
		return
	}

	if Info.Host != "" && Info.Hosts != "" {
		log.Fatalln("err:./only one host parameter")
		return
	}
}
