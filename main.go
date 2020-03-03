package main

import (
	"github.com/WeBankPartners/wecube-plugins-notifications/models"
	"flag"
	"log"
	"github.com/WeBankPartners/wecube-plugins-notifications/api"
)

func main() {
	cfgFile := flag.String("c", "conf/default.json", "config file")
	port := flag.Int("p", 0, "http port")
	flag.Parse()
	err := models.InitConfig(*cfgFile)
	if err != nil {
		log.Printf("init config fail : %v \n", err)
		return
	}
	api.InitHttpServer(*port)
}
