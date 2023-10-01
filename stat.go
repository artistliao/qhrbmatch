package main

import (
	"encoding/json"
	"io"
	"math"
	"os"
	"qhrbmatch/mlog"
	"sort"
	"strconv"
)

func MoneyHand(filename string) error {
	mlog.Debugf("MoneyHand filename:%s", filename)
	defer mlog.Debugf("MoneyHand end!")

	var allMoneyHandData []MoneyHandData
	file, err := os.Open(filename)
	if file == nil || err != nil {
		mlog.Warnf("read file fail, err:%v", err)
		return err
	}
	defer func(file *os.File) {
		cerr := file.Close()
		if cerr != nil {
			mlog.Warnf("file Close err:%v", cerr)
		}
	}(file)

	content, err := io.ReadAll(file)
	if err != nil {
		mlog.Warnf("read to fd fail, err:%v", err)
		return err
	}

	err = json.Unmarshal(content, &allMoneyHandData)
	if err != nil {
		mlog.Warnf("json.Unmarshal content:%v, err:%v", string(content), err)
		return err
	}

	//{
	//"breedname": "铁矿石",
	//"handingfee": "",
	//"lhandsp1": "7097",
	//"lmoneyp1": "9196950.000",
	//"phandsp1": "5287",
	//"playerid": "",
	//"pmoneyp1": "18604650.000"
	//}

	for i, tsd := range allMoneyHandData {
		allMoneyHandData[i].Pmoney, _ = strconv.ParseFloat(tsd.Pmoneyp1, 64)
	}
	sort.Slice(allMoneyHandData, func(i, j int) bool {
		return allMoneyHandData[i].Pmoney > allMoneyHandData[j].Pmoney
	})

	for _, tsd := range allMoneyHandData {
		lhandsp1, _ := strconv.Atoi(tsd.Lhandsp1)
		lmoneyp1, _ := strconv.ParseFloat(tsd.Lmoneyp1, 64)
		phandsp1, _ := strconv.Atoi(tsd.Phandsp1)
		pmoneyp1, _ := strconv.ParseFloat(tsd.Pmoneyp1, 64)
		lavg := lmoneyp1 / float64(lhandsp1)
		pavg := pmoneyp1 / float64(phandsp1)
		mlog.Debugf("品种:%s, 每手盈利:%.2f, 每手亏损:%.2f, 成功率:%.2f%%, 盈亏比:%.2f",
			tsd.Breedname, pavg, lavg, float64(phandsp1)*100/float64(phandsp1+lhandsp1), pavg/lavg)
	}

	return nil
}

func Volatility(filename string) error {
	mlog.Debugf("Volatility filename:%s", filename)
	defer mlog.Debugf("Volatility end!")

	file, err := os.Open(filename)
	if file == nil || err != nil {
		mlog.Warnf("read file fail, err:%v", err)
		return err
	}
	defer func(file *os.File) {
		cerr := file.Close()
		if cerr != nil {
			mlog.Warnf("file Close err:%v", cerr)
		}
	}(file)

	content, err := io.ReadAll(file)
	if err != nil {
		mlog.Warnf("read to fd fail, err:%v", err)
		return err
	}

	var AllProfitStatData []ProfitStatData
	err = json.Unmarshal(content, &AllProfitStatData)
	if err != nil {
		mlog.Warnf("json.Unmarshal content:%v, err:%v", string(content), err)
		return err
	}

	//{
	//	"cumulativenet": "1.32705",
	//	"grossprofit": "23580.000",
	//	"handingfee": "299.490",
	//	"netprofit": "23280.51000",
	//	"playerid": "726878",
	//	"profitrate": "32.70517",
	//	"spprofitrate": "32.70517",
	//	"tradedate": "2020-03-30"
	//}

	var startMoney float64 = 71182.96
	var lastNetprofit float64
	for i, psd := range AllProfitStatData {
		Netprofit, _ := strconv.ParseFloat(psd.Netprofit, 64)
		lastNetprofit = 0.0
		if i > 0 {
			lastNetprofit, _ = strconv.ParseFloat(AllProfitStatData[i-1].Netprofit, 64)
		}
		AllProfitStatData[i].Volatility = (Netprofit - lastNetprofit) / (startMoney + lastNetprofit)
		//mlog.Debugf("日期:%s, 波动率:%.2f%%", psd.Tradedate, volatility*100.0)
	}

	sort.Slice(AllProfitStatData, func(i, j int) bool {
		return AllProfitStatData[i].Volatility > AllProfitStatData[j].Volatility
	})

	var allVolatility float64 = 0.0
	for _, psd := range AllProfitStatData {
		mlog.Debugf("日期:%s, 波动率:%.2f%%", psd.Tradedate, psd.Volatility*100.0)
		allVolatility += math.Abs(psd.Volatility)
	}

	mlog.Debugf("平均波动率:%.2f%%", (allVolatility/float64(len(AllProfitStatData)))*100.0)

	return nil
}

func Volatility2023(filename string) error {
	mlog.Debugf("Volatility filename:%s", filename)
	defer mlog.Debugf("Volatility end!")

	file, err := os.Open(filename)
	if file == nil || err != nil {
		mlog.Warnf("read file fail, err:%v", err)
		return err
	}
	defer func(file *os.File) {
		cerr := file.Close()
		if cerr != nil {
			mlog.Warnf("file Close err:%v", cerr)
		}
	}(file)

	content, err := io.ReadAll(file)
	if err != nil {
		mlog.Warnf("read to fd fail, err:%v", err)
		return err
	}

	var profitData2023 ProfitData2023
	err = json.Unmarshal(content, &profitData2023)
	if err != nil {
		mlog.Warnf("json.Unmarshal content:%v, err:%v", string(content), err)
		return err
	}

	var startMoney float64 = 2436325.41
	var lastNetprofit float64
	for i, dp := range profitData2023.DataPoints {
		lastNetprofit = 0.0
		if i > 0 {
			lastNetprofit = profitData2023.DataPoints[i-1].NetProfit
		}
		profitData2023.DataPoints[i].Volatility = (dp.NetProfit - lastNetprofit) / (startMoney + lastNetprofit)
		//mlog.Debugf("日期:%s, 波动率:%.2f%%", psd.Tradedate, volatility*100.0)
	}

	sort.Slice(profitData2023.DataPoints, func(i, j int) bool {
		return profitData2023.DataPoints[i].Volatility > profitData2023.DataPoints[j].Volatility
	})

	var allVolatility float64 = 0.0
	for _, dp := range profitData2023.DataPoints {
		mlog.Debugf("日期:%s, 波动率:%.2f%%", dp.TradeDate, dp.Volatility*100.0)
		allVolatility += math.Abs(dp.Volatility)
	}

	mlog.Debugf("平均波动率:%.2f%%", (allVolatility/float64(len(profitData2023.DataPoints)))*100.0)

	return nil
}
