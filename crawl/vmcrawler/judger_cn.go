package wmcrawler

import (
	"ASS/config"
	"ASS/crawl"
	"ASS/model"
	"context"
	"strconv"
	"sync"
	"time"

	"github.com/sirupsen/logrus"
)

var Judger_cn = NewCNJudger()

const interval = 60

type result struct {
	canBuy  bool
	canSell bool
	time    int
}

type WMJudgerCN struct {
	results   map[string]result
	strategys map[string]model.StrategyCN
	ctx       context.Context
	cancel    context.CancelFunc
	w         sync.WaitGroup
}

func NewCNJudger() (c *WMJudgerCN) {
	ctx, cancel := context.WithCancel(context.Background())
	return &WMJudgerCN{
		results:   make(map[string]result, 0),
		strategys: make(map[string]model.StrategyCN, 0),
		ctx:       ctx,
		cancel:    cancel,
	}
}

func (judger *WMJudgerCN) Start() {
	judger.w.Add(1)
	go func() {
		defer judger.w.Done()
		logrus.Infoln("WMJudgerCN start(1)...")
		for {
			codes := StockCrawler_cn.GetStockCodes()
			for _, v := range codes {
				judge := WMJudgeCN{}
				var err error
				stockInfos := StockCrawler_cn.GetStockInfo(v, nil)
				if len(stockInfos) == 0 {
					continue
				}
				//找到最近动态股利
				for _, v := range stockInfos {
					judge.StockYield, err = strconv.ParseFloat(v.InterestRatio, 64)
					if err == nil {
						break
					}
				}
				if err != nil {
					continue
				}
				judge.StockPrice, err = strconv.ParseFloat(stockInfos[0].Price, 64)
				if err != nil {
					continue
				}
				judge.StockPE, err = strconv.ParseFloat(stockInfos[0].PE, 64)
				if err != nil {
					continue
				}
				v, ok := judger.strategys[stockInfos[0].Code]
				strategy := model.StrategyCN{
					AimMaxPE:      float64(config.MaxPe),
					AimMinPE:      float64(config.MinPe),
					AimStockMaxPE: float64(config.MaxStockPe),
					AimStockMinPE: float64(config.MinStockPe),
					UseYield:      true,
					PEType:        crawl.SZ_PE,
				}
				if ok {
					strategy = v
				}
				judge.AimMaxPE = strategy.AimMaxPE
				judge.AimMinPE = strategy.AimMinPE
				judge.AimStockMinPE = strategy.AimStockMinPE
				judge.AimStockMaxPE = strategy.AimStockMaxPE
				judge.UseYield = strategy.UseYield

				judge.PE = StockCrawler_cn.GetPE(strategy.PEType)
				judge.Yield = StockCrawler_cn.GetYield()
				if judge.UseYield && (judge.PE == 0 || judge.Yield == 0) {
					continue
				}
				r := result{}
				r.time = time.Now().Second()
				r.canBuy = judge.BuyJudge()
				r.canSell = judge.SellJudge()
				judger.results[stockInfos[0].Code] = r
			}

			select {
			case <-time.After(time.Second * interval):
			case <-judger.ctx.Done():
				return
			}
		}

	}()
	judger.w.Add(1)
	go func() {
		defer judger.w.Done()
		logrus.Infoln("WMJudgerCN start(2)...")
		for k, v := range judger.results {
			if time.Now().Second()-v.time > interval {
				delete(judger.results, k)
			}
		}
		select {
		case <-time.After(time.Second * 1):
		case <-judger.ctx.Done():
			return
		}
	}()
}

func (judger *WMJudgerCN) Stop() {
	judger.cancel()
	judger.w.Wait()
	logrus.Infoln("WMJudgerCN stop...")
}

func (judger *WMJudgerCN) SetStrategy(code string, strategy model.StrategyCN) {
	judger.strategys[code] = strategy
}

func (judger *WMJudgerCN) GetJudgeResult(code string) (result model.StockJudgeResult) {
	//先不加锁
	result.CanBuy = judger.results[code].canBuy
	result.CanSell = judger.results[code].canSell
	return
}
