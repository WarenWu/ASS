package main

import (
	"context"
	"fmt"
	"time"

	"github.com/chromedp/cdproto/network"
	"github.com/chromedp/cdproto/page"
	"github.com/chromedp/chromedp"
	"github.com/go-resty/resty/v2"
	"github.com/sirupsen/logrus"
	"github.com/tidwall/gjson"

	"ASS/utils"
)

func main() {

	// 禁用chrome headless
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

	var hexinV string
	const expr = `delete navigator.__proto__.webdriver;`

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
		// wait for footer element is visible (ie, page is loaded)
		//chromedp.WaitVisible(`.search-icon.base-icon.pointer`, chromedp.ByQuery),
		//chromedp.Click(`textarea.search-input`, chromedp.ByQuery), // find and click "Expand All"
		//chromedp.Sleep(1*time.Second),
		//chromedp.SendKeys(`.search-input`, `茅台`, chromedp.ByQuery),
		//chromedp.Sleep(5*time.Second),

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
	}

	client := resty.New()

	resp, err := client.R().SetFormData(map[string]string{
		"question":         "连续 10 年 ROE 大于 15%， 连续 5 年净利润现金含量大于 80% ，连续 5 年毛利率大于 30%",
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
		return
	}

	jsonResp, err := utils.UnescapeUnicode(resp.Body())
	if err != nil {
		logrus.Error(err)
		return
	}

	logrus.Trace(jsonResp)

	datasJson := gjson.Get(jsonResp, "data.answer.0.txt.0.content.components.0.data.datas")
	test1 := datasJson.Value()
	logrus.Print(test1)

	for _, name := range datasJson.Array() {
		data := name.Value().(map[string]interface{})
		fmt.Printf("%s",data["code"])
	}

	datas := make([]map[string]interface{}, 0)

	//err = json.Unmarshal([]byte(datasJson), &datas)

	logrus.Traceln(datas, len(datas))

	var test string
	err = chromedp.Run(ctx,
		chromedp.ActionFunc(func(ctx context.Context) error {
			_, err := page.AddScriptToEvaluateOnNewDocument(expr).Do(ctx)
			if err != nil {
				logrus.Error(err)
				return err
			}
			return nil
		}),
		chromedp.Navigate(`http://stockpage.10jqka.com.cn/600519/bonus/`),
		chromedp.WaitVisible(`#bonus_table > tbody > tr:nth-child(2) > td:nth-child(10)`),
		chromedp.Text(`#bonus_table > tbody > tr:nth-child(2) > td:nth-child(10)`, &test),
	)
	if err != nil {
		logrus.Error(err)
	}
}
