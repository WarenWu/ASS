package db

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
