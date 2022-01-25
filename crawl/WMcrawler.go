package crawl

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/chromedp/cdproto/network"
	"github.com/chromedp/cdproto/page"
	"github.com/chromedp/chromedp"
	"github.com/go-resty/resty/v2"
	"github.com/sirupsen/logrus"
	"github.com/tidwall/gjson"

	"ASS/db"
	"ASS/utils"
)

type StockCommonInfo struct {
	Code  string `xorm:"code"`
	Name  string `xorm:"name"`
	Price string `xorm:"price"`
	PE    string `xorm:"pe"`
}

type WMCrawler struct {
	FirstCondition string
	Duration       int
	CommonInfos    map[string]*StockCommonInfo
}


func (crawler *WMCrawler) CrawlFromIndexOfA(condition string) string {
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

	var hexinV string                                    //问财网cookie
	const expr = `delete navigator.__proto__.webdriver;` //绕过爬虫检测

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

func (crawler *WMCrawler) GetStockCodesOfA() []string {
	stocks := make([]string, 0)
	jsonResp := crawler.CrawlFromIndexOfA(crawler.FirstCondition)
	if jsonResp == "" {
		return nil
	}
	datasObject := gjson.Get(jsonResp, "data.answer.0.txt.0.content.components.0.data.datas")
	for _, data := range datasObject.Array() {
		stockInfo := data.Value().(map[string]interface{})
		stocks = append(stocks, stockInfo["code"].(string))
	}
	return stocks
}

func (crawler *WMCrawler) GetStockInfosOfA(stockCodes []string, filter Filter) []map[string]string {
	codes := stockCodes
	if len(codes) == 0 {
		codes = crawler.GetStockCodesOfA()
		if codes == nil {
			return nil
		}
	}
	stockInfos := make([]map[string]string, 0)
	for _, code := range codes {
		stockInfo := crawler.GetStockInfoOfA(code, filter)
		stockInfos = append(stockInfos, stockInfo...)
	}
	return stockInfos
}

func (crawler *WMCrawler) GetStockInfoOfA(stockCode string, filter Filter) []map[string]string {
	if stockCode == "" {
		return nil
	}
	stockCommonInfo := &StockCommonInfo{Code: stockCode}
	crawler.updateStockCommonInfo(stockCommonInfo)

	flags := filter()
	stockInfos := make([]map[string]string, crawler.Duration)

	//更新数据库数据
	for _, flag := range flags {
		switch flag {
		case PE:
			crawler.crawlStockPEOfA(stockCode)
		case ROE:
			crawler.crawlStockROEOfA(stockCode)
		case CASH_RATIO:
			crawler.crawlStockCashRatioOfA(stockCode)
		case ASSET_LIABILITY_RATIO:
			crawler.crawlStockAssetLiabilityRatioOfA(stockCode)
		case GROSS_PROFIT_RATIO:
			crawler.crawlStockGrossProfitRatioOfA(stockCode)
		case DIVIDEND_RATIO:
			crawler.crawlStockDividendRatioOfA(stockCode)
		case DIVIDEND:
			crawler.crawlStockDividendOfA(stockCode)
		}
	}

	//更新数据后从数据库读出来
	return stockInfos
}

//func (crawler *WMCrawler)GetPE() float64 {
//}
//
//func (crawler *WMCrawler)GetPE() float64 {
//}

func (crawler *WMCrawler) crawlStockPEOfA(code string) {
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

	const expr = `delete navigator.__proto__.webdriver;` //绕过爬虫检测

	var text string
	err := chromedp.Run(ctx,
		chromedp.ActionFunc(func(ctx context.Context) error {
			_, err := page.AddScriptToEvaluateOnNewDocument(expr).Do(ctx)
			if err != nil {
				logrus.Error(err)
				return err
			}
			return nil
		}),
		chromedp.Navigate(`http://stockpage.10jqka.com.cn/`+code),
		chromedp.WaitVisible(`#fvaluep`, chromedp.ByID),
		chromedp.Text(`#fvaluep`, &text),
	)
	if err != nil {
		logrus.Error(err)
		return
	}
	if err != nil {
		logrus.Error(err)
		return
	}
	stockCommonInfo := &StockCommonInfo{Code: code, PE: text}
	crawler.updateStockCommonInfo(stockCommonInfo)

	stockInfos := make([]db.StockInfo, 0)
	db.DbEngine().Where("code = ?", code).Find(&stockInfos)
	for i, _ := range stockInfos {
		stockInfos[i].PE = text
		crawler.updateStockInfo(code, stockInfos[i].Date, &stockInfos[i])
	}

	crawler.CommonInfos[code].PE = text
}

func (crawler *WMCrawler) crawlStockNameOfA(code string) {
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

	const expr = `delete navigator.__proto__.webdriver;` //绕过爬虫检测

	var text string
	err := chromedp.Run(ctx,
		chromedp.ActionFunc(func(ctx context.Context) error {
			_, err := page.AddScriptToEvaluateOnNewDocument(expr).Do(ctx)
			if err != nil {
				logrus.Error(err)
				return err
			}
			return nil
		}),
		chromedp.Navigate(`http://stockpage.10jqka.com.cn/`+code),
		chromedp.WaitVisible(`#stockNamePlace`, chromedp.ByID),
		chromedp.Evaluate(`document.querySelector('#stockNamePlace').getAttribute('stockname');`, &text),
	)
	if err != nil {
		logrus.Error(err)
		return
	}
	if err != nil {
		logrus.Error(err)
		return
	}
	stockCommonInfo := &StockCommonInfo{Code: code, Name: text}
	crawler.updateStockCommonInfo(stockCommonInfo)

	stockInfos := make([]db.StockInfo, 0)
	db.DbEngine().Where("code = ?", code).Find(&stockInfos)
	for i, _ := range stockInfos {
		stockInfos[i].Name = text
		crawler.updateStockInfo(code, stockInfos[i].Date, &stockInfos[i])
	}

	crawler.CommonInfos[code].Name = text
}

func (crawler *WMCrawler) crawlStockPriceOfA(code string) {
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

	const expr = `delete navigator.__proto__.webdriver;` //绕过爬虫检测

	var text string
	err := chromedp.Run(ctx,
		chromedp.ActionFunc(func(ctx context.Context) error {
			_, err := page.AddScriptToEvaluateOnNewDocument(expr).Do(ctx)
			if err != nil {
				logrus.Error(err)
				return err
			}
			return nil
		}),
		chromedp.Navigate(`http://stockpage.10jqka.com.cn/`+"600519"),
		chromedp.WaitVisible(`#hexm_curPrice`, chromedp.ByID),
		chromedp.Text(`#hexm_curPrice`, &text),
	)
	if err != nil {
		logrus.Error(err)
		return
	}
	if err != nil {
		logrus.Error(err)
		return
	}
	stockCommonInfo := &StockCommonInfo{Code: code, Price: text}
	crawler.updateStockCommonInfo(stockCommonInfo)

	stockInfos := make([]db.StockInfo, 0)
	db.DbEngine().Where("code = ?", code).Find(&stockInfos)
	for i, _ := range stockInfos {
		stockInfos[i].Price = text
		crawler.updateStockInfo(code, stockInfos[i].Date, &stockInfos[i])
	}

	crawler.CommonInfos[code].Price = text
}

func (crawler *WMCrawler) crawlStockROEOfA(code string) {
	condition := code + ` 连续 ` + strconv.Itoa(crawler.Duration) + ` 年 ROE`
	jsonResp := crawler.CrawlFromIndexOfA(condition)
	if jsonResp == "" {
		return
	}

	datasObject := gjson.Get(jsonResp, "data.answer.0.txt.0.content.components.0.data.datas")
	for _, data := range datasObject.Array() {
		dataInfo := data.Value().(map[string]interface{})
		roe := strconv.FormatFloat(dataInfo["ROE"].(float64), 'f', 6, 64)
		date := strconv.FormatUint(dataInfo["时间"].(uint64), 10)[0:3]
		stockInfo := db.StockInfo{
			Code:  code,
			Name: crawler.CommonInfos[code].Name,
			PE:    crawler.CommonInfos[code].PE,
			Price: crawler.CommonInfos[code].Price,
			ROE:   roe,
		}
		crawler.updateStockInfo(code, date, &stockInfo)
	}
}

func (crawler *WMCrawler) crawlStockCashRatioOfA(code string) {
	condition := code + ` 连续 ` + strconv.Itoa(crawler.Duration) + ` 年 净利润现金含量`
	jsonResp := crawler.CrawlFromIndexOfA(condition)
	if jsonResp == "" {
		return
	}

	datasObject := gjson.Get(jsonResp, "data.answer.0.txt.0.content.components.1.data.datas")
	for _, data := range datasObject.Array() {
		dataInfo := data.Value().(map[string]interface{})
		cashRatio := strconv.FormatFloat(dataInfo["净利润现金含量占比"].(float64), 'f', 6, 64)
		date := dataInfo["时间"].(string)[0:3]
		stockInfo := db.StockInfo{
			Code:      code,
			Name: crawler.CommonInfos[code].Name,
			PE:        crawler.CommonInfos[code].PE,
			Price:     crawler.CommonInfos[code].Price,
			CashRatio: cashRatio,
		}
		crawler.updateStockInfo(code, date, &stockInfo)
	}
}

func (crawler *WMCrawler) crawlStockAssetLiabilityRatioOfA(code string) {
	condition := code + ` 连续 ` + strconv.Itoa(crawler.Duration) + ` 年 资产负债率`
	jsonResp := crawler.CrawlFromIndexOfA(condition)
	if jsonResp == "" {
		return
	}

	datasObject := gjson.Get(jsonResp, "data.answer.0.txt.0.content.components.0.data.datas")
	for _, data := range datasObject.Array() {
		dataInfo := data.Value().(map[string]interface{})
		assetLiabilityRatio := strconv.FormatFloat(dataInfo["资产负债率(%)"].(float64), 'f', 6, 64)
		date := strconv.FormatUint(dataInfo["时间区间"].(uint64), 10)[0:3]
		stockInfo := db.StockInfo{
			Code:                code,
			Name: crawler.CommonInfos[code].Name,
			PE:                  crawler.CommonInfos[code].PE,
			Price:               crawler.CommonInfos[code].Price,
			AssetLiabilityRatio: assetLiabilityRatio,
		}
		crawler.updateStockInfo(code, date, &stockInfo)
	}
}

func (crawler *WMCrawler) crawlStockGrossProfitRatioOfA(code string) {
	condition := code + ` 连续 ` + strconv.Itoa(crawler.Duration) + ` 年 毛利率`
	jsonResp := crawler.CrawlFromIndexOfA(condition)
	if jsonResp == "" {
		return
	}

	datasObject := gjson.Get(jsonResp, "data.answer.0.txt.0.content.components.1.tab_list.0.list.0.data.datas")
	for _, data := range datasObject.Array() {
		dataInfo := data.Value().(map[string]interface{})
		grossProfitRatio := strconv.FormatFloat(dataInfo["销售毛利率"].(float64), 'f', 6, 64)
		date := dataInfo["年度"].(string)
		stockInfo := db.StockInfo{
			Code:             code,
			Name: crawler.CommonInfos[code].Name,
			PE:               crawler.CommonInfos[code].PE,
			Price:            crawler.CommonInfos[code].Price,
			GrossProfitRatio: grossProfitRatio,
		}
		crawler.updateStockInfo(code, date, &stockInfo)
	}
}

func (crawler *WMCrawler) crawlStockDividendRatioOfA(code string) {
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

	const expr = `delete navigator.__proto__.webdriver;` //绕过爬虫检测

	var text string
	//var nodes []*cdp.Node
	err := chromedp.Run(ctx,
		chromedp.ActionFunc(func(ctx context.Context) error {
			_, err := page.AddScriptToEvaluateOnNewDocument(expr).Do(ctx)
			if err != nil {
				logrus.Error(err)
				return err
			}
			return nil
		}),
		chromedp.Navigate(`http://stockpage.10jqka.com.cn/`+code+`/bonus/#bonuslist`),
		chromedp.WaitVisible(`#bonus_table`, chromedp.ByID),
		//chromedp.Evaluate(`document.querySelector('#bonus_table tbody').getAttribute('stockname');`, &text),
		//chromedp.Nodes(`#bonus_table tbody`,&nodes)
		//fmt.Println(nodes)
	)
	if err != nil {
		logrus.Error(err)
		return
	}
	if err != nil {
		logrus.Error(err)
		return
	}
	stockCommonInfo := &StockCommonInfo{Code: code, Name: text}
	crawler.updateStockCommonInfo(stockCommonInfo)

	stockInfos := make([]db.StockInfo, 0)
	db.DbEngine().Where("code = ?", code).Find(&stockInfos)
	for i, _ := range stockInfos {
		stockInfos[i].Name = text
		crawler.updateStockInfo(code, stockInfos[i].Date, &stockInfos[i])
	}

	crawler.CommonInfos[code].Name = text

}

func (crawler *WMCrawler) crawlStockDividendOfA(code string) {

}

func (crawler *WMCrawler) updateStockCommonInfo(stockCommonInfo *StockCommonInfo) {

	has, _ := db.DbEngine().Where("code = ? ", stockCommonInfo.Code).Get(&StockCommonInfo{})
	if has {
		db.DbEngine().Where("code = ? ", stockCommonInfo.Code).Update(stockCommonInfo)
	} else {
		db.DbEngine().Insert(stockCommonInfo)
	}
}

func (crawler *WMCrawler) updateStockInfo(code string, date string, stockInfo *db.StockInfo) {

	has, _ := db.DbEngine().Where("code = ? and date = ?", code, date).Get(&db.StockInfo{})
	if has {
		_, err := db.DbEngine().Where("code = ? and date = ?", code, date).Update(stockInfo)
		if err != nil {
			logrus.Error(err)
		}
	} else {
		_, err := db.DbEngine().Insert(stockInfo)
		if err != nil {
			logrus.Error(err)
		}
	}
}
