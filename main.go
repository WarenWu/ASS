package main

import (
	"ASS/config"
	crawler "ASS/crawl/vmcrawler"
	"ASS/module"
	"ASS/router"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/sirupsen/logrus"
)

// @title 股票信息获取
// @version 1.0.0
// @description 提供股票信息相关接口

// @contact.name wuweiming
// @contact.email wuweimingoen@163.com

func main() {
	module.InitModule()
	router := router.Router()

	go func() {
		logrus.Infof("%s start listening at %s", config.AppName, config.Ip+":"+config.Port)
		err := router.Run(config.Ip + ":" + config.Port)
		if err != nil {
			logrus.Errorf("start %s error: %s", config.AppName, err.Error())
			os.Exit(1)
		}
	}()

	start := time.Now().Unix()

	crawler.StockCrawler_cn.Start()
	crawler.Judger_cn.Start()

	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM)
	logrus.Println("quit (%v)", <-sig)
	crawler.Judger_cn.Stop()
	crawler.StockCrawler_cn.Stop()
	duration := time.Now().Unix() - start
	logrus.Println("总耗时:", duration)
}
