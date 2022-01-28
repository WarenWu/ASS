package db

type StockInfo struct {
	Code                string `xorm:"code"`
	Period                string `xorm:"period"`
	Name                string `xorm:"name"`
	Price               string `xorm:"price"`
	PE                  string `xorm:"pe"`
	ROE                 string `xorm:"roe"`
	CashRatio           string `xorm:"cash_ratio"`
	AssetLiabilityRatio string `xorm:"asset_liability_ratio"`
	GrossProfitRatio    string `xorm:"gross_profit_ratio"`
	DividendRatio       string `xorm:"dividend_ratio"`
	InterestRatio       string `xorm:"interest_ratio"`
}

type StockCommonInfo struct {
	Code  string `xorm:"code"`
	Name  string `xorm:"name"`
	Price string `xorm:"price"`
	PE    string `xorm:"pe"`
}
