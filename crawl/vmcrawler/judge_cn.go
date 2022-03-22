package wmcrawler

import (
	"math"
)

type WMJudgeCN struct {
	StockPrice    float64
	StockPE       float64
	StockYield    float64 //动态股利
	PE            float64 //市场PE
	Yield         float64 //国债收益率
	AimMinPE      float64 //买入目标市场市盈率
	AimMaxPE      float64 //卖出目标市场市盈率
	AimStockMinPE float64 //股票买入市盈率
	AimStockMaxPE float64 //股票卖出市盈率
}

func (Judge *WMJudgeCN) BuyJudge() bool {
	return Judge.PE <= Judge.AimMinPE &&
		Judge.StockPE <= Judge.AimStockMinPE &&
		Judge.StockYield >= Judge.Yield &&
		Judge.StockPrice <= Judge.aimPrice()
}

//卖出和价格无关。需长期持有
func (Judge *WMJudgeCN) SellJudge() bool {
	return Judge.PE >= Judge.AimMaxPE &&
		Judge.StockPE >= Judge.AimStockMaxPE &&
		Judge.StockYield <= Judge.Yield
}

func (Judge *WMJudgeCN) aimPrice() float64 {
	priceFormPE := Judge.StockPrice * Judge.AimMinPE / Judge.AimStockMinPE
	priceFromYiled := Judge.StockYield / Judge.Yield
	return math.Min(priceFormPE, priceFromYiled)
}
