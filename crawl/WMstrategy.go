package crawl

import (
	"math"
)

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
		PE,
		ROE,
		CASH_RATIO,
		ASSET_LIABILITY_RATIO,
		GROSS_PROFIT_RATIO,
		DIVIDEND_RATIO,
		DIVIDEND)

	return filter
}
