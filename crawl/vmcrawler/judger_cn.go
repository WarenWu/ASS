package wmcrawler

import (
	"ASS/config"
	"ASS/crawl"
	"ASS/db"
	"strconv"
	"time"

	"github.com/sirupsen/logrus"
)

var Judger = &WMJudgerCN{
	results:   make(map[string]result, 0),
	strageys:  make(map[string]Stragey),
	isRunning: make(chan struct{}),
}

const interval = 60

type result struct {
	canBuy  bool
	canSell bool
	time    int
}

type Stragey struct {
	PEType        int
	AimMinPE      float64 //买入目标市场市盈率
	AimMaxPE      float64 //卖出目标市场市盈率
	AimStockMinPE float64 //股票买入市盈率
	AimStockMaxPE float64 //股票卖出市盈率
}

type WMJudgerCN struct {
	results   map[string]result
	strageys  map[string]Stragey
	isRunning chan struct{}
}

func (judger *WMJudgerCN) Start() {
	go func() {
		logrus.Infoln("WMJudgerCN start...")
		for {

			stockInfos := make([]db.StockInfo, 0)
			codes := StockCrawler_cn.GetStockCodes()
			for _, v := range codes {
				judge := WMJudgeCN{}
				var err error
				stockInfos = StockCrawler_cn.GetStockInfo(v, nil)
				judge.StockPrice, err = strconv.ParseFloat(stockInfos[0].Price, 64)
				if err != nil {
					continue
				}
				judge.StockPE, err = strconv.ParseFloat(stockInfos[0].PE, 64)
				if err != nil {
					continue
				}
				judge.StockYield, err = strconv.ParseFloat(stockInfos[0].InterestRatio, 64)
				if err != nil {
					continue
				}
				v, ok := judger.strageys[stockInfos[0].Code]
				if ok {
					judge.AimMaxPE = v.AimMaxPE
					judge.AimMinPE = v.AimMinPE
					judge.AimStockMinPE = v.AimStockMinPE
					judge.AimStockMaxPE = v.AimStockMaxPE
					judge.PE = StockCrawler_cn.GetPE(v.PEType)
				} else {
					judge.AimMaxPE = float64(config.MaxPe)
					judge.AimMinPE = float64(config.MinPe)
					judge.AimStockMinPE = float64(config.MinStockPe)
					judge.AimStockMaxPE = float64(config.MaxStockPe)
					judge.PE = StockCrawler_cn.GetPE(crawl.SH_PE)
				}

				judge.Yield = StockCrawler_cn.GetYield()
				r := result{}
				r.time = time.Now().Second()
				r.canBuy = judge.BuyJudge()
				r.canSell = judge.SellJudge()
				judger.results[stockInfos[0].Code] = r
			}

			select {
			case <-time.After(time.Second * interval):
			case <-judger.isRunning:
				logrus.Infoln("WMJudgerCN stop...")
				return
			}
		}
	}()
	go func() {
		for k, v := range judger.results {
			if time.Now().Second()-v.time > interval {
				delete(judger.results, k)
			}
		}
		select {
		case <-time.After(time.Second * 1):
		case <-judger.isRunning:
			logrus.Infoln("WMJudgerCN stop...")
			return
		}
	}()
}

func (judger *WMJudgerCN) Stop() {
	judger.isRunning <- struct{}{}
}

func (judger *WMJudgerCN) SetStragey(code string, stragey Stragey) {
	judger.strageys[code] = stragey
}
