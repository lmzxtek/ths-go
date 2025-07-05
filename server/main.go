package main

import (
	"fmt"
	"time"

	"github.com/BurntSushi/toml"
	"github.com/gin-gonic/gin"
	"github.com/lmzxtek/ths-go/srv"
)

// 配置结构体
type Config struct {
	API struct {
		Port      int    `toml:"port"`
		Gmapi     string `toml:"gmapi"`
		Gmcsv     string `toml:"gmcsv"`
		ServerTag string `toml:"server_tag"`
	} `toml:"api"`
}

var cfg Config

func main() {
	// 读取配置文件
	if _, err := toml.DecodeFile("cfg.toml", &cfg); err != nil {
		fmt.Println("Error loading config file:", err)
		return
	}
	fmt.Println(` >>> Load cfg from: cfg.toml`)

	if cfg.API.ServerTag != "" {
		// srv.ServerTag = cfg.API.ServerTag
		srv.SetServerTag(cfg.API.ServerTag)
	}
	if cfg.API.Gmapi == "" {
		fmt.Println("Error: gmapi is empty in config file")
		return
	}
	if cfg.API.Gmcsv == "" {
		fmt.Println("Error: gmcsv is empty in config file")
		return
	}

	srv.SetURL(cfg.API.Gmapi, cfg.API.Gmcsv)

	fmt.Println("")
	now := time.Now()
	// 格式化当前日期为 "YYYY-MM-DD" 格式
	today := now.Format("2006-01-02")
	fmt.Println(" Today-> ", today)

	// fmt.Printf(" config file data:\n %v \n", cfg)
	fmt.Println("")
	fmt.Println(" port -> ", cfg.API.Port)
	fmt.Println(" server_tag -> " + cfg.API.ServerTag)
	fmt.Println(" gmapi -> " + cfg.API.Gmapi)
	fmt.Println(" gmcsv -> " + cfg.API.Gmcsv)
	fmt.Println("")

	r := gin.Default()
	//====================================================
	r.GET("/usage", srv.RouteUsage)

	r.GET("/test", srv.RouteTest)
	r.GET("/test2", srv.RouteTest2)
	r.GET("/test3", srv.RouteTest3)

	//====================================================
	r.GET("/gm1d", srv.RouteGM1d)
	r.GET("/gmpe", srv.RouteGMpe)
	r.GET("/gmvv", srv.RouteGMvv)
	r.GET("/gm1m", srv.RouteGM1m)
	r.GET("/api1m", srv.RouteGMApi1m)

	r.GET("/csv1m", srv.RouteCSVxz1m)
	r.GET("/csvtag", srv.RouteCSVxzTag)
	r.GET("/csvyear", srv.RouteCSVxzYear)
	r.GET("/csvmonth", srv.RouteCSVxzMonth)

	//====================================================
	r.GET("/kbdictts", srv.RouteKBDictTS)
	r.GET("/kbdict", srv.RouteKBDict)
	r.GET("/kbars", srv.RouteKbars)
	r.GET("/kbars2", srv.RouteKbars2)
	r.GET("/kbarsn", srv.RouteKbarsN)
	r.GET("/kbars2n", srv.RouteKbars2N)

	r.GET("/current", srv.RouteCurrent)
	//====================================================
	r.GET("/prevn", srv.RouteDatesPrevN)
	r.GET("/nextn", srv.RouteDatesNextN)
	r.GET("/dateslist", srv.RouteDatesList)
	r.GET("/calendar", srv.RouteCalendar)
	r.GET("/calendar2", srv.RouteCalendar2)

	//====================================================
	r.GET("/market_info", srv.RouteMarketInfo)
	r.GET("/symbols_info", srv.RouteSymbolsInfo)
	r.GET("/history_info", srv.RouteHistoryInfo)

	r.GET("/index_constituents", srv.RouteIndexConstituents)
	r.GET("/industry_constituents", srv.RouteIndustryConstituents)
	r.GET("/industry_category", srv.RouteIndustryCategory)
	r.GET("/symbols_industry", srv.RouteSymbolsIndustry)

	r.GET("/symbols_sector", srv.RouteSymbolsSector)
	r.GET("/secotr_category", srv.RouteSectorCategory)
	r.GET("/sector_constituents", srv.RouteSectorConstituents)

	r.GET("/ration", srv.RouteRation)
	r.GET("/dvidend", srv.RouteDvidend)
	r.GET("/adj_factor", srv.RouteAdjFactor)
	r.GET("/trading_sessions", srv.RouteTradingSessions)

	r.GET("/share_change", srv.RouteShareChange)
	r.GET("/shareholder_num", srv.RouteShareholderNum)
	r.GET("/top_shareholder", srv.RouteTopShareholder)

	//====================================================
	r.GET("/daily_valuation", srv.RouteDailyValuation)
	r.GET("/daily_mktvalue", srv.RouteDailyMktvalue)
	r.GET("/daily_basic", srv.RouteDailyBasic)
	r.GET("/daily_valuation_pt", srv.RouteDailyValuationPt)
	r.GET("/daily_mktvalue_pt", srv.RouteDailyMktvaluePt)
	r.GET("/daily_basic_pt", srv.RouteDailyBasicPt)

	r.GET("/finance_prime", srv.RouteFinancePrime)
	r.GET("/finance_deriv", srv.RouteFinanceDeriv)
	r.GET("/finance_prime_pt", srv.RouteFinancePrimePt)
	r.GET("/finance_deriv_pt", srv.RouteFinanceDerivPt)

	r.GET("/fundamentals_cashflow", srv.RouteFundamentalsCashflow)
	r.GET("/fundamentals_balance", srv.RouteFundamentalsBalance)
	r.GET("/fundamentals_income", srv.RouteFundamentalsIncome)
	r.GET("/fundamentals_cashflow_pt", srv.RouteFundamentalsCashflowPt)
	r.GET("/fundamentals_balance_pt", srv.RouteFundamentalsBalancePt)
	r.GET("/fundamentals_income_pt", srv.RouteFundamentalsIncomePt)

	//====================================================
	r.GET("/abnor_change_stocks", srv.RouteAbnorChangeStocks)
	r.GET("/abnor_change_detail", srv.RouteAbnorChangeDetail)

	r.GET("/fnd_constituents", srv.RouteFndConstituents)
	r.GET("/fnd_portfolio", srv.RouteFndPortfolio)
	r.GET("/fnd_split", srv.RouteFndSplit)
	r.GET("/fnd_dividend", srv.RouteFndDividend)
	r.GET("/fnd_netvalue", srv.RouteFndNetValue)
	r.GET("/fnd_adj_factor", srv.RouteFndAdjFactor)

	r.GET("/hk_inst_holding_info", srv.RouteHKInstHoldingInfo)
	r.GET("/hk_inst_holding_detail_info", srv.RouteHKInstHoldingDetailInfo)

	r.GET("/shszhk_quota_info", srv.RouteSHZSZHKQuotaInfo)
	r.GET("/shszhk_active_stock_top10_info", srv.RouteSHZSZHKActiveStockTop10Info)
	//====================================================

	addr := fmt.Sprintf(":%d", cfg.API.Port)
	fmt.Printf("\nServer running at http://*%s\n\n", addr)
	r.Run(addr)
}
