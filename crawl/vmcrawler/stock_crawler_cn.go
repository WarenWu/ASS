package wmcrawler

import (
	"context"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/chromedp/cdproto/network"
	"github.com/chromedp/cdproto/page"
	"github.com/chromedp/chromedp"
	"github.com/go-resty/resty/v2"
	"github.com/sirupsen/logrus"
	"github.com/tidwall/gjson"

	"ASS/config"
	"ASS/crawl"
	"ASS/db"
	"ASS/utils"
)

type WMCrawlerCN struct {
	firstCondition string
	duration       int
	commonInfos    map[string]*db.StockCommonInfo
	stockCodes     []string
	stockCodesMtx  sync.Mutex
	crawlTimeout   int
	pe_sh          float64
	pe_sz          float64
	yield          float64
	isRunning      bool
}

func NewCNCrawl(firstCondition string, duration int, crawlTimeout int) (c *WMCrawlerCN) {
	c = &WMCrawlerCN{
		firstCondition: firstCondition,
		duration:       duration,
		commonInfos:    make(map[string]*db.StockCommonInfo, 0),
		stockCodes:     make([]string, 0),
		crawlTimeout:   crawlTimeout,
		pe_sh:          -1,
		pe_sz:          -1,
		yield:          -1,
		isRunning:      false,
	}
	return
}

func (crawler *WMCrawlerCN) Start() {
	crawler.isRunning = true
	go func() {
		logrus.Infoln("WMCrawlerCN start crawling...")
		for crawler.isRunning {
			crawler.crawlStockCodes()
			crawler.crawlPE()
			crawler.crawlYield()
			crawler.crawlStockInfos(crawler.GetFilter(nil))
			time.Sleep(time.Second * 600)
		}
		logrus.Infoln("WMCrawlerCN stop crawling...")
	}()
}

func (crawler *WMCrawlerCN) Stop() {
	crawler.isRunning = false
}

func (crawler *WMCrawlerCN) PutStockCode(stockCode string) {
	crawler.stockCodesMtx.Lock()
	crawler.stockCodes = append(crawler.stockCodes, stockCode)
	crawler.stockCodesMtx.Unlock()
}

func (crawler *WMCrawlerCN) DelStockCode(stockCode string) {
	crawler.stockCodesMtx.Lock()
	j := 0
	for _, v := range crawler.stockCodes {
		if v != stockCode {
			crawler.stockCodes[j] = v
		}
	}
	crawler.stockCodes = crawler.stockCodes[:j]
	crawler.stockCodesMtx.Unlock()
}

func (crawler *WMCrawlerCN) GetFilter(filter crawl.Filter) crawl.Filter {
	return func() []string {
		return crawler.filter(filter)
	}
}

func (crawler *WMCrawlerCN) filter(filter crawl.Filter) []string {
	options := make([]string, 0)
	if filter != nil {
		options = append(options, filter()...)
	}
	options = append(options,
		crawl.PE,
		crawl.PRICE,
		crawl.ROE,
		crawl.CASH_RATIO,
		crawl.ASSET_LIABILITY_RATIO,
		crawl.GROSS_PROFIT_RATIO,
		crawl.DIVIDEND_RATIO,
		crawl.INTEREST_RATIO)
	return options
}

func (crawler *WMCrawlerCN) GetStockCodes() []string {
	crawler.stockCodesMtx.Lock()
	stockCodes := make([]string, len(crawler.stockCodes))
	copy(stockCodes, crawler.stockCodes)
	crawler.stockCodesMtx.Unlock()
	return stockCodes
}

func (crawler *WMCrawlerCN) crawlStockCodes() {
	jsonResp := crawler.CrawlFromIndex(crawler.firstCondition)
	if jsonResp == "" {
		return
	}
	datasObject := gjson.Get(jsonResp, "data.answer.0.txt.0.content.components.0.data.datas")
	for _, data := range datasObject.Array() {
		stockInfo := data.Value().(map[string]interface{})
		crawler.stockCodes = append(crawler.stockCodes, stockInfo["code"].(string))
	}
}

func (crawler *WMCrawlerCN) GetStockInfos(filter crawl.Filter) []map[string]string {

	stockInfos := make([]map[string]string, crawler.duration)
	//从数据库读出来
	return stockInfos
}

func (crawler *WMCrawlerCN) crawlStockInfos(filter crawl.Filter) {
	stockCodes := crawler.GetStockCodes()
	for _, code := range stockCodes {
		// go func(code string) {
		// 	crawler.crawlStockInfo(code, filter)
		// }(code)
		crawler.crawlStockInfo(code, filter)
	}
}

func (crawler *WMCrawlerCN) GetStockInfo(stockCode string, filter crawl.Filter) []map[string]string {
	if stockCode == "" {
		return nil
	}
	stockInfos := make([]map[string]string, crawler.duration)
	//从数据库读出来
	return stockInfos
}

func (crawler *WMCrawlerCN) crawlStockInfo(stockCode string, filter crawl.Filter) {
	if stockCode == "" {
		return
	}
	stockCommonInfo := &db.StockCommonInfo{Code: stockCode}
	crawler.updateStockCommonInfo(stockCommonInfo)

	flags := filter()
	//更新数据库数据
	crawler.crawlStockName(stockCode)
	for _, flag := range flags {
		switch flag {
		case crawl.PE:
			crawler.crawlStockPE(stockCode)
		case crawl.PRICE:
			crawler.crawlStockPrice(stockCode)
		case crawl.ROE:
			crawler.crawlStockROE(stockCode)
		case crawl.CASH_RATIO:
			crawler.crawlStockCashRatio(stockCode)
		case crawl.ASSET_LIABILITY_RATIO:
			crawler.crawlStockAssetLiabilityRatio(stockCode)
		case crawl.GROSS_PROFIT_RATIO:
			crawler.crawlStockGrossProfitRatio(stockCode)
		case crawl.DIVIDEND_RATIO:
			crawler.crawlStockDividendRatio(stockCode)
		case crawl.INTEREST_RATIO:
			crawler.crawlStockInterestRatio(stockCode)
		}
	}
}

func (crawler *WMCrawlerCN) GetPE(t int) (pe float64) {
	if t == crawl.SH_PE {
		pe = crawler.pe_sh
	}
	if t == crawl.SZ_PE {
		pe = crawler.pe_sz
	}
	return
}

func (crawler *WMCrawlerCN) crawlPE() {
	opts := append(chromedp.DefaultExecAllocatorOptions[:],
		chromedp.Flag("headless", config.Headless),
		chromedp.UserAgent(`Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/91.0.4472.164 Safari/537.36`),
	)
	timeoutCtx, cancel := context.WithTimeout(context.Background(), time.Duration(crawler.crawlTimeout)*time.Second)
	defer cancel()
	allocCtx, cancel := chromedp.NewExecAllocator(timeoutCtx, opts...)
	defer cancel()

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
				logrus.Errorln(err)
				return err
			}
			return nil
		}),
		chromedp.Navigate(`http://value500.com/PE.asp`),
		chromedp.Sleep(1000*time.Millisecond),
		chromedp.Evaluate(`
		shA = document.querySelector('div:nth-child(3)>table:nth-child(4)>tbody>tr:nth-child(1)>td:nth-child(2)>table:nth-child(1)>tbody>tr:nth-child(1)>td:nth-child(1)>table:nth-child(4)>tbody>tr:nth-child(2)>td:nth-child(2)').innerText;
		szA = document.querySelector('div:nth-child(3)>table:nth-child(4)>tbody>tr:nth-child(1)>td:nth-child(2)>table:nth-child(1)>tbody>tr:nth-child(1)>td:nth-child(1)>table:nth-child(4)>tbody>tr:nth-child(2)>td:nth-child(3)').innerText;
		shA +":"+ szA;
		`, &text),
	)
	if err != nil {
		logrus.Errorln(err)
		return
	}
	pes := strings.Split(text, ":")
	if len(pes) == 2 {
		crawler.pe_sh, err = strconv.ParseFloat(pes[0], 64)
		if err != nil {
			logrus.Errorln(err)
			return
		}
		crawler.pe_sz, err = strconv.ParseFloat(pes[1], 64)
		if err != nil {
			logrus.Errorln(err)
			return
		}
	}
}

func (crawler *WMCrawlerCN) GetYield() float64 {
	return crawler.yield
}

func (crawler *WMCrawlerCN) crawlYield() {
	opts := append(chromedp.DefaultExecAllocatorOptions[:],
		chromedp.Flag("headless", config.Headless),
		chromedp.UserAgent(`Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/91.0.4472.164 Safari/537.36`),
	)
	timeoutCtx, cancel := context.WithTimeout(context.Background(), time.Duration(crawler.crawlTimeout)*time.Second)
	defer cancel()
	allocCtx, cancel := chromedp.NewExecAllocator(timeoutCtx, opts...)
	defer cancel()

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
				logrus.Errorln(err)
				return err
			}
			return nil
		}),
		chromedp.Navigate(`https://yield.chinabond.com.cn/cbweb-mn/yield_main?locale=zh_CN`),
		chromedp.Sleep(1000*time.Millisecond),
		chromedp.Evaluate(`
		document.querySelector('#detailFrame>table.tablelist>tbody>tr:nth-child(13)>td:nth-child(2)').innerText;
		`, &text),
	)
	if err != nil {
		logrus.Errorln(err)
		return
	}
	crawler.yield, err = strconv.ParseFloat(text, 64)
	if err != nil {
		logrus.Errorln(err)
		return
	}
}

func (crawler *WMCrawlerCN) CrawlFromIndex(condition string) string {
	opts := append(chromedp.DefaultExecAllocatorOptions[:],
		chromedp.Flag("headless", config.Headless),
		chromedp.UserAgent(`Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/91.0.4472.164 Safari/537.36`),
	)
	timeoutCtx, cancel := context.WithTimeout(context.Background(), time.Duration(crawler.crawlTimeout)*time.Second)
	defer cancel()
	allocCtx, cancel := chromedp.NewExecAllocator(timeoutCtx, opts...)
	defer cancel()
	ctx, cancel := chromedp.NewContext(
		allocCtx,
		chromedp.WithLogf(logrus.Printf),
	)
	defer cancel()

	var hexinV string
	const expr = `delete navigator.__proto__.webdriver;`
	err := chromedp.Run(ctx,
		chromedp.ActionFunc(func(ctx context.Context) error {
			_, err := page.AddScriptToEvaluateOnNewDocument(expr).Do(ctx)
			if err != nil {
				logrus.Errorln(err)
				return err
			}
			return nil
		}),
		chromedp.Navigate(`http://www.iwencai.com/unifiedwap/home/index`),
		chromedp.ActionFunc(func(ctx context.Context) error {
			cookies, err := network.GetAllCookies().Do(ctx)
			if err != nil {
				logrus.Errorln(err)
				return err
			}

			for i, cookie := range cookies {
				if cookie.Name == "v" {
					hexinV = cookie.Value
				}
				logrus.Tracef("chrome cookie %d: %+v\n", i, cookie)
			}
			return nil
		}),
	)
	if err != nil {
		logrus.Errorln(err)
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
		logrus.Errorln(err)
		return ""
	}

	jsonResp, err := utils.UnescapeUnicode(resp.Body())
	if err != nil {
		logrus.Errorln(err)
		return ""
	}
	return jsonResp
}

func (crawler *WMCrawlerCN) crawlStockPE(code string) {
	opts := append(chromedp.DefaultExecAllocatorOptions[:],
		chromedp.Flag("headless", config.Headless),
		chromedp.UserAgent(`Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/91.0.4472.164 Safari/537.36`),
	)
	timeoutCtx, cancel := context.WithTimeout(context.Background(), time.Duration(crawler.crawlTimeout)*time.Second)
	defer cancel()
	allocCtx, cancel := chromedp.NewExecAllocator(timeoutCtx, opts...)
	defer cancel()

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
				logrus.Errorln(err)
				return err
			}
			return nil
		}),
		chromedp.Navigate(`http://stockpage.10jqka.com.cn/realHead_v2.html#hs_`+code),
		chromedp.Sleep(1000*time.Millisecond),
		chromedp.Evaluate(`document.getElementById('fvaluep').innerText`, &text),
	)
	if err != nil {
		logrus.Errorln(err)
		return
	}
	stockCommonInfo := &db.StockCommonInfo{Code: code, PE: text}
	crawler.updateStockCommonInfo(stockCommonInfo)

	stockInfos := make([]db.StockInfo, 0)
	db.DbEngine().Where("code = ?", code).Find(&stockInfos)
	for i, _ := range stockInfos {
		stockInfos[i].PE = text
		crawler.updateStockInfo(code, stockInfos[i].Period, &stockInfos[i])
	}

	_, ok := crawler.commonInfos[code]
	if ok {
		crawler.commonInfos[code].PE = text
	} else {
		crawler.commonInfos[code] = stockCommonInfo
	}
}

func (crawler *WMCrawlerCN) crawlStockName(code string) {
	opts := append(chromedp.DefaultExecAllocatorOptions[:],
		chromedp.Flag("headless", config.Headless),
		chromedp.UserAgent(`Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/91.0.4472.164 Safari/537.36`),
	)
	timeoutCtx, cancel := context.WithTimeout(context.Background(), time.Duration(crawler.crawlTimeout)*time.Second)
	defer cancel()
	allocCtx, cancel := chromedp.NewExecAllocator(timeoutCtx, opts...)
	defer cancel()

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
				logrus.Errorln(err)
				return err
			}
			return nil
		}),
		chromedp.Navigate(`http://stockpage.10jqka.com.cn/`+code),
		chromedp.WaitVisible(`#stockNamePlace`, chromedp.ByID),
		chromedp.Evaluate(`document.querySelector('#stockNamePlace').getAttribute('stockname');`, &text),
	)
	if err != nil {
		logrus.Errorln(err)
		return
	}
	stockCommonInfo := &db.StockCommonInfo{Code: code, Name: text}
	crawler.updateStockCommonInfo(stockCommonInfo)

	stockInfos := make([]db.StockInfo, 0)
	db.DbEngine().Where("code = ?", code).Find(&stockInfos)
	for i, _ := range stockInfos {
		stockInfos[i].Name = text
		crawler.updateStockInfo(code, stockInfos[i].Period, &stockInfos[i])
	}

	_, ok := crawler.commonInfos[code]
	if ok {
		crawler.commonInfos[code].Name = text
	} else {
		crawler.commonInfos[code] = stockCommonInfo
	}
}

func (crawler *WMCrawlerCN) crawlStockPrice(code string) {
	opts := append(chromedp.DefaultExecAllocatorOptions[:],
		chromedp.Flag("headless", config.Headless),
		chromedp.UserAgent(`Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/91.0.4472.164 Safari/537.36`),
	)
	timeoutCtx, cancel := context.WithTimeout(context.Background(), time.Duration(crawler.crawlTimeout)*time.Second)
	defer cancel()
	allocCtx, cancel := chromedp.NewExecAllocator(timeoutCtx, opts...)
	defer cancel()

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
				logrus.Errorln(err)
				return err
			}
			return nil
		}),
		chromedp.Navigate(`http://stockpage.10jqka.com.cn/realHead_v2.html#hs_`+code),
		chromedp.WaitVisible(`#topenprice`, chromedp.ByID),
		chromedp.Sleep(1000*time.Millisecond),
		chromedp.Text(`#hexm_curPrice`, &text),
	)
	if err != nil {
		logrus.Errorln(err)
		return
	}
	stockCommonInfo := &db.StockCommonInfo{Code: code, Price: text}
	crawler.updateStockCommonInfo(stockCommonInfo)

	stockInfos := make([]db.StockInfo, 0)
	db.DbEngine().Where("code = ?", code).Find(&stockInfos)
	for i, _ := range stockInfos {
		stockInfos[i].Price = text
		crawler.updateStockInfo(code, stockInfos[i].Period, &stockInfos[i])
	}

	_, ok := crawler.commonInfos[code]
	if ok {
		crawler.commonInfos[code].Price = text
	} else {
		crawler.commonInfos[code] = stockCommonInfo
	}
}

func (crawler *WMCrawlerCN) crawlStockROE(code string) {
	condition := code + ` 连续 ` + strconv.Itoa(crawler.duration) + ` 年 ROE`
	jsonResp := crawler.CrawlFromIndex(condition)
	if jsonResp == "" {
		logrus.Errorln("***" + code + "***" + `爬取ROE失败`)
		return
	}
	datasObject := gjson.Get(jsonResp, "data.answer.0.txt.0.content.components.1.data.datas")
	for _, data := range datasObject.Array() {
		dataInfo := data.Value().(map[string]interface{})
		var roe string
		var period string
		if v, ok := dataInfo["ROE"]; ok {
			roe = strconv.FormatFloat(v.(float64), 'f', 12, 64)[:6]
		} else {
			logrus.Warnln("***" + code + "***")
			continue
		}
		if v, ok := dataInfo["时间区间"]; ok {
			period = (strconv.FormatFloat(v.(float64), 'f', 12, 64))[0:4]
		} else {
			logrus.Warnln("***" + code + "***")
			continue
		}

		stockInfo := db.StockInfo{
			Code:   code,
			Period: period,
			Name:   crawler.commonInfos[code].Name,
			PE:     crawler.commonInfos[code].PE,
			Price:  crawler.commonInfos[code].Price,
			ROE:    roe,
		}
		crawler.updateStockInfo(code, period, &stockInfo)
	}
}

func (crawler *WMCrawlerCN) crawlStockCashRatio(code string) {
	condition := code + ` 连续 ` + strconv.Itoa(crawler.duration) + ` 年 净利润现金含量`
	jsonResp := crawler.CrawlFromIndex(condition)
	if jsonResp == "" {
		logrus.Errorln("***" + code + "***" + `爬取现金含量失败`)
		return
	}
	datasObject := gjson.Get(jsonResp, "data.answer.0.txt.0.content.components.0.data.datas")
	for _, data := range datasObject.Array() {
		dataInfo := data.Value().(map[string]interface{})
		var cashRatio string
		var period string
		if v, ok := dataInfo["净利润现金含量占比"]; ok {
			cashRatio = strconv.FormatFloat(v.(float64), 'f', 12, 64)[:6]
		} else {
			logrus.Warnln("***" + code + "***")
			continue
		}
		if v, ok := dataInfo["时间区间"]; ok {
			data := v.(string)
			period = "20" + data[0:2]
		} else {
			logrus.Warnln("***" + code + "***")
			continue
		}

		stockInfo := db.StockInfo{
			Code:      code,
			Period:    period,
			Name:      crawler.commonInfos[code].Name,
			PE:        crawler.commonInfos[code].PE,
			Price:     crawler.commonInfos[code].Price,
			CashRatio: cashRatio,
		}
		//logrus.Infoln(`---------------------------` + code + `:` + cashRatio + `---------------------------------`)
		crawler.updateStockInfo(code, period, &stockInfo)
	}
}

func (crawler *WMCrawlerCN) crawlStockAssetLiabilityRatio(code string) {
	condition := code + ` 连续 ` + strconv.Itoa(crawler.duration) + ` 年 资产负债率`
	jsonResp := crawler.CrawlFromIndex(condition)
	if jsonResp == "" {
		logrus.Errorln("***" + code + "***" + `爬取资产负债率失败`)
		return
	}

	datasObject := gjson.Get(jsonResp, "data.answer.0.txt.0.content.components.0.data.datas")
	for _, data := range datasObject.Array() {
		dataInfo := data.Value().(map[string]interface{})

		var assetLiabilityRatio string
		var period string
		if v, ok := dataInfo["资产负债率(%)"]; ok {
			assetLiabilityRatio = strconv.FormatFloat(v.(float64), 'f', 12, 64)[:6]
		} else {
			logrus.Warnln("***" + code + "***")
			continue
		}
		if v, ok := dataInfo["时间区间"]; ok {
			period = (strconv.FormatFloat(v.(float64), 'f', 12, 64))[0:4]
		} else {
			logrus.Warnln("***" + code + "***")
			continue
		}

		stockInfo := db.StockInfo{
			Code:                code,
			Period:              period,
			Name:                crawler.commonInfos[code].Name,
			PE:                  crawler.commonInfos[code].PE,
			Price:               crawler.commonInfos[code].Price,
			AssetLiabilityRatio: assetLiabilityRatio,
		}
		crawler.updateStockInfo(code, period, &stockInfo)
	}
}

func (crawler *WMCrawlerCN) crawlStockGrossProfitRatio(code string) {
	condition := code + ` 连续 ` + strconv.Itoa(crawler.duration) + ` 年 毛利率`
	jsonResp := crawler.CrawlFromIndex(condition)
	if jsonResp == "" {
		logrus.Errorln("***" + code + "***" + `爬取毛利率失败`)
		return
	}

	datasObject := gjson.Get(jsonResp, "data.answer.0.txt.0.content.components.1.data.datas")
	for _, data := range datasObject.Array() {
		dataInfo := data.Value().(map[string]interface{})
		var grossProfitRatio string
		var period string
		if v, ok := dataInfo["销售毛利率"]; ok {
			grossProfitRatio = v.(string)
		} else {
			logrus.Warnln("***" + code + "***")
			continue
		}
		if v, ok := dataInfo["报告期"]; ok {
			data := v.(string)
			period = "20" + data[0:2]
		} else {
			logrus.Warnln("***" + code + "***")
			continue
		}

		stockInfo := db.StockInfo{
			Code:             code,
			Period:           period,
			Name:             crawler.commonInfos[code].Name,
			PE:               crawler.commonInfos[code].PE,
			Price:            crawler.commonInfos[code].Price,
			GrossProfitRatio: grossProfitRatio,
		}
		crawler.updateStockInfo(code, period, &stockInfo)
	}
}

func (crawler *WMCrawlerCN) crawlStockDividendRatio(code string) {
	opts := append(chromedp.DefaultExecAllocatorOptions[:],
		chromedp.Flag("headless", config.Headless),
		chromedp.UserAgent(`Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/91.0.4472.164 Safari/537.36`),
	)
	timeoutCtx, cancel := context.WithTimeout(context.Background(), time.Duration(crawler.crawlTimeout)*time.Second)
	defer cancel()
	allocCtx, cancel := chromedp.NewExecAllocator(timeoutCtx, opts...)
	defer cancel()

	ctx, cancel := chromedp.NewContext(
		allocCtx,
		chromedp.WithLogf(logrus.Printf),
	)
	defer cancel()

	const expr = `delete navigator.__proto__.webdriver;`
	var drs string
	js := `
	bt = function getDR(){
		let ret = "";
		var trs
		if (window.frames["CC"] == undefined)
		{
			trs = document.querySelector('#bonus_table>tbody').children;
		}
		else
		{
			trs = window.frames["CC"].document.querySelector('#bonus_table>tbody').children;
		}
		for (var i = 0; i < trs.length; i++) {
			if (trs[i].className != "J_pageritem ")
			{
				continue;
			}
			t = trs[i].children[0].innerText.substr(0,4);
			b = parseFloat(trs[i].children[9].innerText);
			if(Object.is(b,NaN))
			{
				continue;
			}
			ret += t + "-"+ b + ":"; 
		}
		return ret
	}()`
	err := chromedp.Run(ctx,
		chromedp.ActionFunc(func(ctx context.Context) error {
			_, err := page.AddScriptToEvaluateOnNewDocument(expr).Do(ctx)
			if err != nil {
				logrus.Errorln(err)
				return err
			}
			return nil
		}),
		chromedp.Navigate(`http://basic.10jqka.com.cn/`+code+`/bonus.html`),
		chromedp.Sleep(1000*time.Millisecond),
		chromedp.Evaluate(js, &drs),
	)
	if err != nil {
		logrus.Errorln(err)
		return
	}
	for _, dr := range strings.Split(drs, ":") {
		v := strings.Split(dr, "-")
		if len(v) != 2 {
			continue
		}
		period := v[0]
		dividendRatio := v[1]
		stockInfo := db.StockInfo{
			Code:          code,
			Period:        period,
			Name:          crawler.commonInfos[code].Name,
			PE:            crawler.commonInfos[code].PE,
			Price:         crawler.commonInfos[code].Price,
			DividendRatio: dividendRatio,
		}
		crawler.updateStockInfo(code, period, &stockInfo)
	}
}

func (crawler *WMCrawlerCN) crawlStockInterestRatio(code string) {
	opts := append(chromedp.DefaultExecAllocatorOptions[:],
		chromedp.Flag("headless", config.Headless),
		chromedp.UserAgent(`Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/91.0.4472.164 Safari/537.36`),
	)
	timeoutCtx, cancel := context.WithTimeout(context.Background(), time.Duration(crawler.crawlTimeout)*time.Second)
	defer cancel()
	allocCtx, cancel := chromedp.NewExecAllocator(timeoutCtx, opts...)
	defer cancel()

	ctx, cancel := chromedp.NewContext(
		allocCtx,
		chromedp.WithLogf(logrus.Printf),
	)
	defer cancel()

	const expr = `delete navigator.__proto__.webdriver;`
	var irs string
	js := `bt = function getIR(){
		var price = ` + crawler.commonInfos[code].Price + `
		let ret = "";
		var irs;
		if (window.frames["CC"] == undefined)
		{
			irs = document.querySelector('#bonus_table>tbody').children;
		}
		else
		{
			irs = window.frames["CC"].document.querySelector('#bonus_table>tbody').children;
		}
		for (var i = 0; i < irs.length; i++) {
			if (irs[i].className != "J_pageritem ")
			{
				continue;
			}
			t = irs[i].children[0].innerText.substr(0,4);
			arr = irs[i].children[4].innerText.split("派");
			if(arr.length != 2)
			{
				continue;
			}
			n = parseFloat(arr[0]);
			b = parseFloat(arr[1]);
			r = 100*b/(n*parseFloat(price))
			if(Object.is(r,NaN))
			{
				r="0"
			}
			ret +=t + "-" + r + ":";  
		}
		return ret
	}()`
	err := chromedp.Run(ctx,
		chromedp.ActionFunc(func(ctx context.Context) error {
			_, err := page.AddScriptToEvaluateOnNewDocument(expr).Do(ctx)
			if err != nil {
				logrus.Errorln(err)
				return err
			}
			return nil
		}),
		chromedp.Navigate(`http://basic.10jqka.com.cn/`+code+`/bonus.html`),
		chromedp.Sleep(1000*time.Millisecond),
		chromedp.Evaluate(js, &irs),
	)
	if err != nil {
		logrus.Errorln(err)
		return
	}
	for _, ir := range strings.Split(irs, ":") {
		v := strings.Split(ir, "-")
		if len(v) != 2 {
			continue
		}
		period := v[0]
		interestRatio := v[1][:6]
		stockInfo := db.StockInfo{
			Code:          code,
			Period:        period,
			Name:          crawler.commonInfos[code].Name,
			PE:            crawler.commonInfos[code].PE,
			Price:         crawler.commonInfos[code].Price,
			InterestRatio: interestRatio,
		}
		crawler.updateStockInfo(code, period, &stockInfo)
	}
}

func (crawler *WMCrawlerCN) updateStockCommonInfo(stockCommonInfo *db.StockCommonInfo) {

	has, _ := db.DbEngine().Where("code = ? ", stockCommonInfo.Code).Get(&db.StockCommonInfo{})
	if has {
		db.DbEngine().Where("code = ? ", stockCommonInfo.Code).Update(stockCommonInfo)
	} else {
		db.DbEngine().Insert(stockCommonInfo)
	}
}

func (crawler *WMCrawlerCN) updateStockInfo(code string, period string, stockInfo *db.StockInfo) {

	has, _ := db.DbEngine().Where("code = ? and period = ?", code, period).Get(&db.StockInfo{})
	if has {
		_, err := db.DbEngine().Where("code = ? and period = ?", code, period).Update(stockInfo)
		if err != nil {
			logrus.Errorln(err)
		}
	} else {
		_, err := db.DbEngine().Insert(stockInfo)
		if err != nil {
			logrus.Errorln(err)
		}
	}
}
