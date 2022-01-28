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
	SH_PE     = iota //深交市盈率
	SZ_PE            //上交市盈率
	CN_YIELD        //10年期国债收益率
	AM_YIELD        //十年期美国国债收益率
)

type Filter func() []string //返回股票需要查询的财务指标
type BuyJudge func() bool   //判断买入时机
type SellJudge func() bool  //判断卖出时机

//定义爬虫接口：
type Crawler interface {
	GetFilter() Filter                               //返回默认信息项
	GetStockCodes() []string                         //返回默认股票代码
	GetStockInfo(string, Filter) []map[string]string //指定股票代码和筛选器，爬取股票信息
	GetStockInfos([]string, Filter) []map[string]string
	GetPE(int) float64    //爬取市盈率
	GetYield(int) float64 //爬取整体收益率（国债）
}

//买入卖出接口，返回是否可以买入或卖出
type Processor interface {
	BuyStatus(BuyJudge) bool
	SellStatus(SellJudge) bool
}

//策略接口，包括返回股票信息项、买入策略、卖出策略
type Strategy interface {
	GetBuyJudge() BuyJudge
	GetSellJudge() SellJudge
}
