package wmcrawler

import (
	"math"

	"ASS/crawl"
)

type WMStrategyCN struct {
	//StockType     string
	StockPrice    float64
	StockPE       float64
	StockYield    float64 //派息率
	PE            float64 //市场PE
	Yield         float64 //国债收益率
	AimMinPE      float64 //买入目标市场市盈率
	AimMaxPE      float64 //卖出目标市场市盈率
	AimStockMinPE float64 //股票买入市盈率
	AimStockMaxPE float64 //股票卖出市盈率
}

func (strategy *WMStrategyCN) GetBuyJudge() crawl.BuyJudge {
	return func() bool {
		return strategy.PE <= strategy.AimMinPE &&
			strategy.StockPE <= strategy.AimStockMinPE &&
			strategy.StockYield >= strategy.Yield &&
			strategy.StockPrice <= strategy.aimPrice()
	}
}

//卖出和价格无关。需长期持有
func (strategy *WMStrategyCN) GetSellJudge() crawl.SellJudge {
	return func() bool {
		return strategy.PE >= strategy.AimMaxPE &&
			strategy.StockPE >= strategy.AimStockMaxPE &&
			strategy.StockYield <= strategy.Yield
	}
}

func (strategy *WMStrategyCN) aimPrice() float64 {
	priceFormPE := strategy.StockPrice * strategy.AimMinPE / strategy.AimStockMinPE
	priceFromYiled := strategy.StockYield / strategy.Yield
	return math.Min(priceFormPE, priceFromYiled)
}


