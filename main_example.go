package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"time"

	"github.com/chromedp/cdproto/cdp"
	"github.com/chromedp/cdproto/network"
	"github.com/chromedp/cdproto/page"
	"github.com/chromedp/cdproto/runtime"
	"github.com/chromedp/chromedp"
	"github.com/go-resty/resty/v2"
	_ "github.com/go-sql-driver/mysql"
	"github.com/sirupsen/logrus"
	"github.com/tidwall/gjson"

	"ASS/utils"
)

func writeHTML(content string) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html")
		io.WriteString(w, strings.TrimSpace(content))
	})
}

func main1() {
	ctx, cancel := chromedp.NewContext(context.Background())
	defer cancel()

	ts := httptest.NewServer(writeHTML(`<!doctype html>
<html>
<body>
  <div id="content">the content</div>
</body>
</html>`))
	defer ts.Close()

	const expr = `(function(d, id, v) {
		var b = d.querySelector('body');
		var el = d.createElement('div');
		el.id = id;
		el.innerText = v;
		b.insertBefore(el, b.childNodes[0]);
	})(document, %q, %q);`

	var nodes []*cdp.Node
	if err := chromedp.Run(ctx,
		chromedp.Navigate(ts.URL),
		chromedp.Nodes(`document`, &nodes, chromedp.ByJSPath),
		chromedp.WaitVisible(`#content`),
		chromedp.ActionFunc(func(ctx context.Context) error {
			s := fmt.Sprintf(expr, "thing", "a new thing!")
			_, exp, err := runtime.Evaluate(s).Do(ctx)
			if err != nil {
				return err
			}
			if exp != nil {
				return exp
			}
			return nil
		}),
		chromedp.WaitVisible(`#thing`),
	); err != nil {
		logrus.Fatal(err)
	}

	logrus.Println("Document tree:")
	logrus.Print(nodes[0].Dump("  ", "  ", false))

}

func init() {
	logrus.SetLevel(logrus.DebugLevel)
	logrus.SetReportCaller(true)
}

func main2() {

	client := resty.New()

	resp, err := client.R().SetFormData(map[string]string{
		"question":         "连续 5 年的 ROE 大于 15%",
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
		"hexin-v":    "A-P9s3FRhz-dtEsLQNzc7vZEciyO2HcasWy7ThVAP8K5VA3anagHasE8S58m",
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

	datasJson := gjson.Get(jsonResp, "data.answer.0.txt.0.content.components.0.data.datas").String()

	datas := make([]map[string]interface{}, 0)

	err = json.Unmarshal([]byte(datasJson), &datas)

	logrus.Trace(datas)
}

func main3() {
	body, _ := utils.Request(utils.GET, "http://www.iwencai.com/unifiedwap/result", map[string]string{
		"w": "连续 5 年的 ROE 大于 15%",
	}, nil, nil)

	bodyResp := string(body)
	logrus.Println(gjson.Get(bodyResp, "data.answer.1.txt.1.disambiguationed"))
}

func main4() {

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
		//"question":         "连续 10 年 ROE 大于 15%， 连续 5 年净利润现金含量大于 80% ，连续 5 年毛利率大于 30%",
		"question":         "600519 连续 5 年 毛利率",
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

	logrus.Println(jsonResp)

	datasJson := gjson.Get(jsonResp, "data.answer.0.txt.0.content.components.1.data.datas")
	test1 := datasJson.Value()
	logrus.Print(test1)

	for _, name := range datasJson.Array() {
		data := name.Value().(map[string]interface{})
		fmt.Printf("%s", data["code"])
	}

	datas := make([]map[string]interface{}, 0)

	//err = json.Unmarshal([]byte(datasJson), &datas)

	logrus.Traceln(datas, len(datas))

	// var text string
	// err = chromedp.Run(ctx,
	// 	chromedp.ActionFunc(func(ctx context.Context) error {
	// 		_, err := page.AddScriptToEvaluateOnNewDocument(expr).Do(ctx)
	// 		if err != nil {
	// 			logrus.Error(err)
	// 			return err
	// 		}
	// 		return nil
	// 	}),
	// 	//chromedp.Navigate(`http://stockpage.10jqka.com.cn/realHead_v2.html#hs_` + "600519"),
	// 	chromedp.Navigate(`http://stockpage.10jqka.com.cn/600519/bonus/#bonuslist`),
	// 	//chromedp.WaitVisible(`#fvaluep`, chromedp.ByID),
	// 	chromedp.Evaluate(`bt = function getBt(){
	//         return document.getElementById('dataifm').contentWindow.document.querySelector('#bonus_table>tbody').children[1].children[9].innerText;
	//     }();   `, &text),
	// )
	// if err != nil {
	// 	logrus.Error(err)
	// 	return
	// }
}

func main5() {
	// 禁用chrome headless
	opts := append(chromedp.DefaultExecAllocatorOptions[:],
		chromedp.Flag("headless", false),
		chromedp.UserAgent(`Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/91.0.4472.164 Safari/537.36`),
	)
	timeoutCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	allocCtx, cancel := chromedp.NewExecAllocator(timeoutCtx, opts...)
	defer cancel()
	// create chrome instance
	ctx, cancel := chromedp.NewContext(
		allocCtx,
		chromedp.WithLogf(logrus.Printf),
	)
	defer cancel()

	const expr = `delete navigator.__proto__.webdriver;`

	//price :="120.0"
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
		//chromedp.Navigate(`http://stockpage.10jqka.com.cn/600519/bonus/#bonuslist`),
		//chromedp.Sleep(1000*time.Millisecond),
		// chromedp.Evaluate(`
		// bt = function getIR(){
		// 	var price = `+price+`
		//     let ret = "";
		//     var trs = document.getElementById('dataifm').contentWindow.document.querySelector('#bonus_table>tbody').children;
		//     for (var i = 0; i < trs.length; i++) {
		//         if (trs[i].className != "J_pageritem ")
		//         {
		//             continue;
		//         }
		//         t = trs[i].children[0].innerText.substr(0,4);
		//         arr = trs[i].children[4].innerText.split("派");
		// 		if(!arr.length)
		// 		{
		// 			continue;
		// 		}
		// 		n = parseFloat(arr[0]);
		// 		b = parseFloat(arr[1]);
		// 		r = 100*b/(n*parseFloat(price))
		// 		if(Object.is(r,NaN))
		//         {
		//             r="0"
		//         }
		//         ret +=t+"-"+r+":";
		//     }
		//     return ret
		// }()`, &text),

		// chromedp.Navigate(`http://stockpage.10jqka.com.cn/realHead_v2.html#hs_` + "600519"),
		// chromedp.Sleep(1000*time.Millisecond),
		// chromedp.Evaluate(`document.getElementById('fvaluep').innerText`,&text),

		// chromedp.Navigate(`http://stockpage.10jqka.com.cn/realHead_v2.html#hs_` + "600519"),
		// chromedp.Sleep(1000*time.Millisecond),
		// chromedp.Text(`#hexm_curPrice`, &text),

		// chromedp.Navigate(`http://value500.com/PE.asp`),
		// chromedp.Sleep(1000*time.Millisecond),
		// chromedp.Evaluate(`
		// shA = document.querySelector('div:nth-child(3)>table:nth-child(4)>tbody>tr:nth-child(1)>td:nth-child(2)>table:nth-child(1)>tbody>tr:nth-child(1)>td:nth-child(1)>table:nth-child(4)>tbody>tr:nth-child(2)>td:nth-child(2)').innerText;
		// szA = document.querySelector('div:nth-child(3)>table:nth-child(4)>tbody>tr:nth-child(1)>td:nth-child(2)>table:nth-child(1)>tbody>tr:nth-child(1)>td:nth-child(1)>table:nth-child(4)>tbody>tr:nth-child(2)>td:nth-child(3)').innerText;
		// shA +":"+ szA;
		// `, &text),

		chromedp.Navigate(`https://yield.chinabond.com.cn/cbweb-mn/yield_main?locale=zh_CN`),
		chromedp.Sleep(1000*time.Millisecond),
		chromedp.Evaluate(`
		document.querySelector('#detailFrame>table.tablelist>tbody>tr:nth-child(13)>td:nth-child(2)').innerText;
		`, &text),
	)
	if err != nil {
		logrus.Error(err)
		return
	}
}
