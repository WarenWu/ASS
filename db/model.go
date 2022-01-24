package db

type StockInfo struct {
	Code                string  `xorm:"code"`
	Date                string  `xorm:"date"`
	Name                string  `xorm:"name"`
	Price               string `xorm:"price"`
	PE                  string `xorm:"pe"`
	ROE                 string `xorm:"roe"`
	CashRatio           string `xorm:"cash_ratio"`
	AssetLiabilityRatio string `xorm:"asset_liability_ratio"`
	GrossProfitRatio    string `xorm:"gross_profit_ratio"`
	DvidendRatio        string `xorm:"dvidend_ratio"`
	Dvidend             string `xorm:"dvidend"`
}
