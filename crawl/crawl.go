package crawl

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
