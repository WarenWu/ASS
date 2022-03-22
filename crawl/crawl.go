package crawl

const (
	CODE                  = "股票代码"
	NAME                  = "股票简称"
	DATE                  = "时间区间"
	PRICE                 = "价格"
	PE                    = "市盈率"
	ROE                   = "加权净资产收益率"
	CASH_RATIO            = "现金含量占比"
	ASSET_LIABILITY_RATIO = "资产负债率"
	GROSS_PROFIT_RATIO    = "毛利率"
	DIVIDEND_RATIO        = "派息率"
	INTEREST_RATIO        = "动态股利"
)

const (
	SH_PE    = iota //深交市盈率
	SZ_PE           //上交市盈率
	CN_YIELD        //10年期国债收益率
	AM_YIELD        //十年期美国国债收益率
)

type Filter func() []string //返回股票需要查询的财务指标

//定义爬虫接口：
type StockCrawler interface {
	Start()
	Stop()
	PutStockCode(string)        //向爬取池增加自定义股票
	DelStockCode(string)        //向爬取池删除自定义股票
	GetStockCodes() []string    //返回爬取池所有股票
	GetStockInfo(string) string //指定股票，获取股票json信息
	GetStockInfos() string      //获取所有股票json信息
	GetPE(int) float64          //爬取市盈率
	GetYield() float64          //爬取整体收益率（国债）
}

type Judge interface {
	BuyJudge() bool
	SellJudge() bool
}
