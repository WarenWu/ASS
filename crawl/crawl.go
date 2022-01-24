package crawl

import (
	"context"
	"math"
	"strconv"
	"time"

	"github.com/chromedp/cdproto/network"
	"github.com/chromedp/cdproto/page"
	"github.com/chromedp/chromedp"
	"github.com/go-resty/resty/v2"
	"github.com/sirupsen/logrus"
	"github.com/tidwall/gjson"

	"ASS/utils"
)

const (
	CODE                  = "股票简称"
	NAME                  = "股票代码"
	DATE                  = "时间区间"
	PE                    = "市盈率"
	ROE                   = "加权净资产收益率"
	CASH_RATIO            = "现金含量占比"
	ASSET_LIABILITY_RATIO = "资产负债率"
	GROSS_PROFIT_RATIO    = "毛利率"
	DIVIDEND_RATIO        = "派息率"
	DIVIDEND              = "每股派息"
)

type Filter func() []string //返回股票需要查询的财务指标
type BuyJudge func() bool   //判断买入时机
type SellJudge func() bool  //判断卖出时机

type Crawler interface {
	GetStockInfo(string, Filter) []map[string]string //指定股票代码和筛选器，爬取股票信息
	GetPE() float64                                  //爬取市盈率
	GetYield() float64                               //爬取整体收益率（国债）
}

type Processor interface {
	BuyStatus(BuyJudge) bool
	SellStatus(SellJudge) bool
}

type Strategy interface {
	GetFilter() Filter
	GetBuyJudge() BuyJudge
	GetSellJudge() SellJudge
}

type WMCrawler struct {
	FirstCondition string
	Duration int
}

func (crawler *WMCrawler)CrawlFromIndex(condition string) string{
	opts := append(chromedp.DefaultExecAllocatorOptions[:],
		chromedp.Flag("headless", false),
		chromedp.UserAgent(`Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/91.0.4472.164 Safari/537.36`),
	)
	allocCtx, cancel := chromedp.NewExecAllocator(context.Background(), opts...)
	defer cancel()

	// create chrome instance
	ctx, cancel := chromedp.NewContext(
		allocCtx,
		chromedp.WithLogf(logrus.Printf),
	)
	defer cancel()

	var hexinV string //问财网cookie
	const expr = `delete navigator.__proto__.webdriver;`//绕过爬虫检测

	err := chromedp.Run(ctx,
		chromedp.ActionFunc(func(ctx context.Context) error {
			_, err := page.AddScriptToEvaluateOnNewDocument(expr).Do(ctx)
			if err != nil {
				logrus.Error(err)
				return err
			}
			return nil
		}),
		chromedp.Navigate(`http://www.iwencai.com/unifiedwap/home/index`),
		chromedp.Sleep(1*time.Second),
		//chromedp.Evaluate(`document.cookie`, &hexinV),
		chromedp.ActionFunc(func(ctx context.Context) error {
			cookies, err := network.GetAllCookies().Do(ctx)
			if err != nil {
				logrus.Error(err)
				return err
			}

			for i, cookie := range cookies {
				if cookie.Name == "v" {
					hexinV = cookie.Value
				}
				logrus.Trace("chrome cookie %d: %+v", i, cookie)
			}
			return nil
		}),
	)
	if err != nil {
		logrus.Error(err)
		return ""
	}

	client := resty.New()

	resp, err := client.R().SetFormData(map[string]string{
		"question":         condition,
		"perpage":          "1000",
		"page":             "1",
		"secondary_intent": "",
		"log_info":         `{"input_type":"click"}`,
		"source":           "Ths_iwencai_Xuangu",
		"version":          "2.0",
		"query_area":       "",
		"block_list":       "",
		"add_info":         `{"urp":{"scene":1,"company":1,"business":1},"contentType":"json"}`,
	}).SetHeaders(map[string]string{
		"Referer":    "http://www.iwencai.com/unifiedwap/result?w=%E8%BF%9E%E7%BB%AD%205%20%E5%B9%B4%E7%9A%84%20ROE%20%E5%A4%A7%E4%BA%8E%2015%25",
		"User-Agent": "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/91.0.4472.124 Safari/537.36 Edg/91.0.864.67",
		"Origin":     "http://www.iwencai.com",
		"hexin-v":    hexinV,
		"Host":       "www.iwencai.com",
	}).Post("http://www.iwencai.com/unifiedwap/unified-wap/v2/result/get-robot-data")

	if err != nil {
		logrus.Error(err)
		return ""
	}

	//去掉转义
	jsonResp, err := utils.UnescapeUnicode(resp.Body())
	if err != nil {
		logrus.Error(err)
		return ""
	}
	return jsonResp
}

func (crawler *WMCrawler)GetStockCodes() []string{
	stocks := make([]string,0)
	jsonResp := crawler.CrawlFromIndex(crawler.FirstCondition)
	if jsonResp == ""{
		return nil
	}
	datasObject := gjson.Get(jsonResp, "data.answer.0.txt.0.content.components.0.data.datas")
	for _, data := range datasObject.Array() {
		stockInfo := data.Value().(map[string]interface{})
		stocks=append(stocks,stockInfo["code"].(string))
	}
	return stocks
}

func (crawler *WMCrawler)GetStocksInfo(stockCodes []string, filter Filter) []map[string]string{
	codes := stockCodes
	if len(codes) == 0 {
		codes = crawler.GetStockCodes()
		if codes == nil {
			return nil
		}
	}

}

func (crawler *WMCrawler)CrawlStockName(code string)string{

}

func (crawler *WMCrawler)CrawlStockPE(code string)string{

}

func (crawler *WMCrawler)CrawlStockDate(code string)string{

}

func (crawler *WMCrawler)CrawlStockROE(code string)string{

}

func (crawler *WMCrawler) CrawlStockCashRatio(code string)string{

}

func (crawler *WMCrawler) CrawlStockAssetLiabilityRatio(code string)string{

}

func (crawler *WMCrawler) CrawlStockGrossProfitRatio(code string)string{

}

func (crawler *WMCrawler) CrawlStockDividendRatio(code string)string{

}

func (crawler *WMCrawler) CrawlStockDividend(code string)string{

}

func (crawler *WMCrawler)GetStockInfo(stockCode string, filter Filter) []map[string]string{
	flags := filter()
	condition := stockCode
	
	for _,flag:= range flags{
		switch flag {
		case NAME:

		}
		'
	}
	 
	jsonResp := crawler.Crawl(crawler.FirstCondition)
	if jsonResp == ""{
		return nil
	}
}

//func (crawler *WMCrawler)GetPE() float64 {
//}
//
//func (crawler *WMCrawler)GetPE() float64 {
//}

type WMStrategy struct {
	StockType     string
	StockPrice    float64
	StockPE       float64
	StockYield    float64 //派息率
	StockDividend float64 //每股派息
	PE            float64 //市场PE
	Yield         float64 //国债收益率
	AimPE         float64 //目标市场市盈率
	AimStockPE    float64 //目标股票市盈率
}

func (strategy *WMStrategy) GetFilter() Filter {
	return func() []string {
		return strategy.filter()
	}
}

func (strategy *WMStrategy) GetBuyJudge() BuyJudge {
	return func() bool {
		return strategy.PE < strategy.AimPE &&
			strategy.StockPE < strategy.AimStockPE &&
			strategy.StockYield < strategy.Yield &&
			strategy.StockPrice < strategy.aimPrice()
	}
}

func (strategy *WMStrategy) GetSellJudge() SellJudge {
	return func() bool {
		return true
	}
}

func (strategy *WMStrategy) aimPrice() float64 {
	priceFormPE := strategy.StockPrice * strategy.AimPE / strategy.StockPE
	priceFromYiled := strategy.StockDividend / strategy.Yield
	return math.Min(priceFormPE, priceFromYiled)
}

func (strategy *WMStrategy) filter() []string {
	filter := make([]string, 0)
	filter = append(filter,
		NAME,
		DATE,
		PE,
		ROE,
		CASH_RATIO,
		ASSET_LIABILITY_RATIO,
		GROSS_PROFIT_RATIO,
		DIVIDEND_RATIO,
		DIVIDEND)

	return filter
}
