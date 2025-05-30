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

// 配置结构体
// type Config struct {
// 	API struct {
// 		Port  int    `json:"port"`
// 		gmapi string `json:"gmapi"`
// 		gmcsv string `json:"gmcsv"`
// 	}
// }

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
	fmt.Println(" -=> Today is", today)

	// filename := "cfg.json"
	// file, err := os.Open(filename)
	// if err != nil {
	// 	fmt.Printf("配置文件打开失败: %v", err)
	// }
	// if err := json.NewDecoder(file).Decode(&cfg); err != nil {
	// 	fmt.Printf("配置解析失败: %v", err)
	// }

	fmt.Printf(" config file data:\n %v ", cfg)

	// http.HandleFunc("/", handler)
	// http.ListenAndServe(":8080", nil)

	r := gin.Default()
	r.GET("/test", srv.RouteTest)
	r.GET("/test2", srv.RouteTest2)
	r.GET("/test3", srv.RouteTest3)

	r.GET("/current", srv.RouteCurrent)
	r.GET("/kbars", srv.RouteKbars)
	r.GET("/kbarsn", srv.RouteKbarsN)

	r.GET("/calendar", srv.RouteCalendar)
	r.GET("/calendar2", srv.RouteCalendar2)

	addr := fmt.Sprintf(":%d", cfg.API.Port)
	fmt.Printf("\nServer running at http://*%s\n\n", addr)
	r.Run(addr)
}

// func handler(w http.ResponseWriter, r *http.Request) {
// 	fmt.Fprintf(w, "Hello, world!")
// }
