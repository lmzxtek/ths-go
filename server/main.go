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
		Port  int    `toml:"port"`
		gmapi string `toml:"gmapi"`
		gmcsv string `toml:"gmcsv"`
	} `toml:"api"`
}

var cfg Config

func main() {
	// 读取配置文件
	if _, err := toml.DecodeFile("cfg.toml", &cfg); err != nil {
		fmt.Println("Error loading config file:", err)
		return
	}
	fmt.Println(` -=> Loading params from: cfg.toml`)

	now := time.Now()
	// 格式化当前日期为 "YYYY-MM-DD" 格式
	today := now.Format("2006-01-02")
	fmt.Println(" -=> Today: ", today)

	fmt.Printf(" config file data:\n %v ", cfg)

	r := gin.Default()

	r.GET("/test", srv.RouteTest)
	r.GET("/test2", srv.RouteTest2)
	r.GET("/test3", srv.RouteTest3)

	r.GET("/gm1m", srv.RouteGM1m)
	r.GET("/csv1m", srv.RouteCSVxz1m)
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

	addr := fmt.Sprintf(":%d", cfg.API.Port)
	fmt.Printf("\nServer running at http://*%s\n\n", addr)
	r.Run(addr)
}
