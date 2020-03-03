package models

import (
	"sync"
	"log"
	"encoding/json"
	"os"
	"io/ioutil"
	"strings"
)

type HttpConfig struct {
	Port  string  `json:"port"`
	Token  string  `json:"token"`
}

type MailConfig struct {
	Enable  bool  `json:"enable"`
	Sender  []*SenderConfig  `json:"sender"`
}

type SenderConfig struct {
	Protocol  string  `json:"protocol"`
	Name  string  `json:"name"`
	User  string  `json:"user"`
	Password  string  `json:"password"`
	Server  string  `json:"server"`
	Token  string  `json:"token"`
}

type GlobalConfig struct {
	Http  *HttpConfig  `json:"http"`
	Mail  *MailConfig  `json:"mail"`
}

var (
	config     *GlobalConfig
	lock       = new(sync.RWMutex)
)

func Config() *GlobalConfig {
	lock.RLock()
	defer lock.RUnlock()
	return config
}

func InitConfig(cfg string) error {
	if cfg == "" {
		log.Println("use -c to specify configuration file")
	}
	_, err := os.Stat(cfg)
	if os.IsExist(err) {
		log.Println("config file not found")
		return err
	}
	b,err := ioutil.ReadFile(cfg)
	if err != nil {
		log.Printf("read file %s error %v \n", cfg, err)
		return err
	}
	configContent := strings.TrimSpace(string(b))
	var c GlobalConfig
	err = json.Unmarshal([]byte(configContent), &c)
	if err != nil {
		log.Println("parse config file:", cfg, "fail:", err)
		return err
	}
	lock.Lock()
	defer lock.Unlock()
	config = &c
	log.Println("read config file:", cfg, "successfully")
	return nil
}
