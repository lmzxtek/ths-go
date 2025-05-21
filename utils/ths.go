package utils

import (
	"fmt"
	"os"
	"time"

	"github.com/BurntSushi/toml"
)

var (
	gmIP      string
	gmPort    int
	csvIP     string
	csvPort   int
	thsDir    string
	dataFld   string
	count     int
	symbols   []string
	indexList []string
	debug     bool
)

// 配置结构体
type Config struct {
	THS struct {
		Debug     bool     `toml:"debug"`
		DataDir   string   `toml:"data_dir"`
		ThsDir    string   `toml:"ths_dir"`
		CSVIP     string   `toml:"csv_ip"`
		CSVPor    int      `toml:"csv_port"`
		GMIP      string   `toml:"gm_ip"`
		GMPort    int      `toml:"gm_port"`
		Count     int      `toml:"count"`
		Symbols   []string `toml:"symbols"`
		IndexList []string `toml:"index_list"`
	} `toml:"ths"`
}

func ToDropUTC(ttUTC string) time.Time {
	tt, err := time.Parse(time.RFC3339, ttUTC)
	if err != nil {
		fmt.Println("Error parsing time:", err)
		return time.Time{}
	}
	return tt.UTC().Truncate(24 * time.Hour)
}

func main() {
	args := os.Args[1:] // os.Args[0] 是脚本名，后面是参数

	if len(args) > 0 && args[0] == "test" {
		fmt.Println("\n >>> Start to test script ...")
		// year := 2024

		symbol := "SHSE.601088"
		ipo := "2007-10-09"

		// 读取配置文件
		if _, err := toml.DecodeFile("cfg.toml", &Config{}); err != nil {
			fmt.Println("Error loading config file:", err)
			return
		}
		fmt.Println(` -=> Loading params from: cfg.toml`)

		// 调用读取本地CSV的函数
		df1m, err := readLocalCSV1M(symbol, count, dataFld, ipo, true, gmIP, gmPort, csvIP, csvPort, debug)
		if err != nil {
			fmt.Println("Error reading local CSV 1M:", err)
			return
		}
		fmt.Println("df1m:\n", df1m)

		dfpe, err := readLocalCSVPe(symbol, count, dataFld, ipo, gmIP, gmPort, csvIP, csvPort, debug)
		if err != nil {
			fmt.Println("Error reading local CSV Pe:", err)
			return
		}
		fmt.Println("dfpe:\n", dfpe)

		dfvv, err := readLocalCSVVv(symbol, count, dataFld, ipo, gmIP, gmPort, csvIP, csvPort, debug)
		if err != nil {
			fmt.Println("Error reading local CSV Vv:", err)
			return
		}
		fmt.Println("dfvv:\n", dfvv)

	} else if len(args) > 0 && args[0] == "update" {
		fmt.Println("\n >>> Start update THS daily indications ...")

		// 读取配置文件
		if _, err := toml.DecodeFile("cfg.toml", &Config{}); err != nil {
			fmt.Println("Error loading config file:", err)
			return
		}
		fmt.Println(` -=> Loading params from: cfg.toml`)

		// 调用更新THS每日指标的函数
		if err := thsProcSymbols(symbols, count, true, csvIP, csvPort, gmIP, gmPort, dataFld, thsDir, debug); err != nil {
			fmt.Println("Error processing symbols:", err)
			return
		}

	} else {
		fmt.Println("\n >>> Opps: cmd options or params error:", args, "...")
		fmt.Println("\n         > go run ths_v15.go test ")
		fmt.Println("\n         > go run ths_v15.go update ")
	}

	// 打印完成时间
	fmt.Printf("\n !!! Nice, All mission finished at %s. \n", time.Now().Format("2006-01-02 15:04:05"))
}

// 假设这些函数已经实现
func readLocalCSV1M(symbol string, count int, dataFld string, ipo string, saveCSV bool, gmIP string, gmPort int, csvIP string, csvPort int, debug bool) (interface{}, error) {
	// 实现读取本地CSV 1M的逻辑
	return nil, nil
}

func readLocalCSVPe(symbol string, count int, dataFld string, ipo string, gmIP string, gmPort int, csvIP string, csvPort int, debug bool) (interface{}, error) {
	// 实现读取本地CSV Pe的逻辑
	return nil, nil
}

func readLocalCSVVv(symbol string, count int, dataFld string, ipo string, gmIP string, gmPort int, csvIP string, csvPort int, debug bool) (interface{}, error) {
	// 实现读取本地CSV Vv的逻辑
	return nil, nil
}

func thsProcSymbols(symbols []string, count int, toCSV bool, csvIP string, csvPort int, gmIP string, gmPort int, dataFld string, thsFld string, debug bool) error {
	// 实现更新THS每日指标的逻辑
	return nil
}
