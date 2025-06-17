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

	r.GET("/usage", srv.RouteUsage)

	r.GET("/test", srv.RouteTest)
	r.GET("/test2", srv.RouteTest2)
	r.GET("/test3", srv.RouteTest3)

	r.GET("/gm1m", srv.RouteGM1m)

	r.GET("/csv1m", srv.RouteCSVxz1m)
	r.GET("/csvtag", srv.RouteCSVxzTag)
	r.GET("/csvyear", srv.RouteCSVxzYear)
	r.GET("/csvmonth", srv.RouteCSVxzMonth)

	r.GET("/kbdictts", srv.RouteKBDictTS)
	r.GET("/kbdict", srv.RouteKBDict)
	r.GET("/kbars", srv.RouteKbars)
	r.GET("/kbarsn", srv.RouteKbarsN)

	r.GET("/current", srv.RouteCurrent)

	r.GET("/prevn", srv.RouteDatesPrevN)
	r.GET("/nextn", srv.RouteDatesNextN)
	r.GET("/calendar", srv.RouteCalendar)
	r.GET("/calendar2", srv.RouteCalendar2)

	r.GET("/marketinfo", srv.RouteMarketInfo)
	r.GET("/symbolsinfo", srv.RouteSymbolsInfo)
	r.GET("/historyinfo", srv.RouteHistoryInfo)

	r.GET("/trading_sessions", srv.RouteTradingSessions)
	r.GET("/index_constituents", srv.RouteIndexConstituents)
	r.GET("/symbol_industry", srv.RouteSymbolIndustry)
	r.GET("/industry_constituents", srv.RouteIndustryConstituents)
	r.GET("/industry_category", srv.RouteIndustryCategory)

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

	addr := fmt.Sprintf(":%d", cfg.API.Port)
	fmt.Printf("\nServer running at http://*%s\n\n", addr)
	r.Run(addr)
}
