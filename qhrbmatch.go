package main

import (
	"qhrbmatch/mlog"
)

func main() {
	//日志初始化
	mlog.InitLogger(&mlog.Params{
		Path:       "./log/qhrbmatch.log", //文件路径
		MaxSize:    2,                     //MB 单个日志文件最大
		MaxBackups: 3,                     //备份个数
		MaxAge:     10,                    //保存时间,天
		Level:      -1,                    //# 日志级别
	})

	mlog.Debugf("qhrbmatch start!")
	defer mlog.Debugf("qhrbmatch end!")

	//_ = MoneyHand("./data/2020_kly_moneyhand_data.txt")
	//_ = MoneyHand("./data/2022_kly_moneyhand_data.txt")
	//_ = Volatility("./data/2020_kly_profit.txt")
	_ = Volatility2023("./data/2023_wdd_profit.txt")
}
