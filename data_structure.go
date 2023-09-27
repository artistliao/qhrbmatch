package main

// 每个品种成交利润和手数
type MoneyHandData struct {
	Breedname  string  `json:"breedname"`
	Handingfee string  `json:"handingfee"`
	Lhandsp1   string  `json:"lhandsp1"`
	Lmoneyp1   string  `json:"lmoneyp1"`
	Phandsp1   string  `json:"phandsp1"`
	Playerid   string  `json:"playerid"`
	Pmoneyp1   string  `json:"pmoneyp1"`
	Pmoney     float64 `json:"pmoney"`
}

// 每天的利润
type ProfitStatData struct {
	Cumulativenet string  `json:"cumulativenet"`
	Grossprofit   string  `json:"grossprofit"`
	Handingfee    string  `json:"handingfee"`
	Netprofit     string  `json:"netprofit"`
	Playerid      string  `json:"playerid"`
	Profitrate    string  `json:"profitrate"`
	Spprofitrate  string  `json:"spprofitrate"`
	Tradedate     string  `json:"tradedate"`
	Volatility    float64 `json:"volatility"`
}
