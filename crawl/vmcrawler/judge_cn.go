package wmcrawler

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
	UseYield      bool    //是否使用股利
}

func (Judge *WMJudgeCN) BuyJudge() bool {
	if Judge.UseYield {
		return Judge.PE <= Judge.AimMinPE &&
			Judge.StockPE <= Judge.AimStockMinPE &&
			Judge.StockYield >= Judge.Yield // ||
		//Judge.StockPrice <= Judge.aimPrice()
	} else {
		return Judge.PE <= Judge.AimMinPE &&
			Judge.StockPE <= Judge.AimStockMinPE // ||
		//Judge.StockPrice <= Judge.aimPrice()
	}

}

//卖出和价格无关。需长期持有
func (Judge *WMJudgeCN) SellJudge() bool {
	if Judge.UseYield {
		return Judge.PE >= Judge.AimMaxPE ||
			Judge.StockPE >= Judge.AimStockMaxPE ||
			Judge.StockYield <= Judge.Yield
	} else {
		return Judge.PE >= Judge.AimMaxPE ||
			Judge.StockPE >= Judge.AimStockMaxPE
	}

}
