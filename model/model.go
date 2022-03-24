package model

type StockInfo struct {
	Code                string `xorm:"code" json:"code"`
	Period              string `xorm:"period" json:"period"`
	Name                string `xorm:"name" json:"name"`
	Price               string `xorm:"price" json:"price"`
	PE                  string `xorm:"pe" json:"pe"`
	ROE                 string `xorm:"roe" json:"roe"`
	CashRatio           string `xorm:"cash_ratio" json:"cash_ratio"`
	AssetLiabilityRatio string `xorm:"asset_liability_ratio" json:"asset_liability_ratio"`
	GrossProfitRatio    string `xorm:"gross_profit_ratio" json:"gross_profit_ratio"`
	DividendRatio       string `xorm:"dividend_ratio" json:"dividend_ratio"`
	InterestRatio       string `xorm:"interest_ratio" json:"interest_ratio"`
}
type StockCommonInfo struct {
	Code  string `xorm:"code"`
	Name  string `xorm:"name"`
	Price string `xorm:"price"`
	PE    string `xorm:"pe"`
}
type StockJudgeResult struct {
	CanBuy  bool `json:"can_buy"`
	CanSell bool `json:"can_sell"`
}

type StrategyCN struct {
	PEType        int     `json:"pe_type"`
	AimMinPE      float64 `json:"aim_min_pe"`       //买入目标市场市盈率
	AimMaxPE      float64 `json:"aim_max_pe"`       //卖出目标市场市盈率
	AimStockMinPE float64 `json:"aim_min_stock_pe"` //股票买入市盈率
	AimStockMaxPE float64 `json:"aim_max_stock_pe"` //股票卖出市盈率
	UseYield      bool    `json:"use_yield"`
}
