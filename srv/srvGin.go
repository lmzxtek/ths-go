package srv

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/lmzxtek/ths-go/gm"
)

var (
	localRand *rand.Rand
	gmapi     string = "http://localhost:5000"
	gmcsv     string = "http://localhost:5002"
)

// init 函数在程序启动时自动执行，设置随机数种子
func init() {
	// 创建一个新的本地随机数生成器
	localRand = rand.New(rand.NewSource(time.Now().UnixNano()))
}

func SetURL(gmAPI, gmCSV string) {
	if gmAPI != "" {
		gmapi = SmartURLHandler(gmAPI, false)
	}
	if gmCSV != "" {
		gmcsv = SmartURLHandler(gmCSV, false)
	}
}

var serverTag string = "/api"

func SetServerTag(srvtag string) {
	if srvtag != "" {
		serverTag = srvtag
	}
}

// RandomStock 从预定义的股票代码数组中随机返回一个元素
func randomStock() string {
	stocks := []string{"AAPL", "GOOG", "AMZN", "MSFT", "TSLA"}
	// 使用本地的随机数生成器
	return stocks[localRand.Intn(len(stocks))]
}

// randomIntWithRange 返回一个在 [a-b, a+b] 区间内的随机整数
func randomIntWithRange(a, b int) int {
	// 使用本地的随机数生成器生成一个在 [a-b, a+b] 区间内的随机整数
	return localRand.Intn(2*b+1) + a - b
}

// randomDateInLastYear 返回过去一年之内的随机日期
func randomDateInLastYear() string {
	// 获取当前时间
	now := time.Now()
	// 计算一年前的时间
	oneYearAgo := now.AddDate(-1, 0, 0)
	// 计算一年前和现在的时间差（以秒为单位）
	diff := now.Unix() - oneYearAgo.Unix()
	// 生成一个在过去一年内的随机时间点（以秒为单位）
	randomUnix := oneYearAgo.Unix() + int64(localRand.Intn(int(diff)))
	// 将随机时间点转换为时间对象
	randomTime := time.Unix(randomUnix, 0)
	// 格式化随机时间为 "YYYY-MM-DD" 格式
	return randomTime.Format("2006-01-02")
}

// 定义配置结构体
type HTMLConfig struct {
	ServerTag string
	HostURL   string
	Symbol    string
	Sididx    string
}

// type Config struct {
// 	Host      string `json:"host"`
// 	Port      int    `json:"port"`
// 	Debug     bool   `json:"debug"`
// 	FldData   string `json:"fld_data"`
// 	ServerTag string `json:"servertag"`
// }
// var cfg Config
// var baseDir string

// 生成HTML的构造函数
func BuildHTML(cfg HTMLConfig) string {
	url := cfg.HostURL
	// 获取当前时间
	now := time.Now()
	// 格式化当前日期为 "YYYY-MM-DD" 格式
	today := now.Format("2006-01-02")
	strCurYear := now.Format("2006")
	// strMonth := now.Format("01")
	// strDay := now.Format("02")
	// strDate := strYear + "-" + strMonth + "-" + strDay

	strPreMonth := now.AddDate(0, -1, 0).Format("2006-01-02")
	preYear := now.AddDate(-1, 0, 0)
	strYear1 := preYear.Format("2006")
	strPreYear := preYear.Format("2006-01-02")

	syms := "SHSE.601088,SZSE.300917"
	sym := cfg.Symbol
	idx := cfg.Sididx

	//===================================================================
	strDatePrevN1 := "prevn?date=" + today + "&count=10"
	strDatePrevN2 := "prevn?date=" + today + "&count=10&include=false"
	strDatePrevN3 := "prevn?date=" + today + "&count=365"
	strDateNextN1 := "nextn?date=" + today + "&count=10"
	strDateNextN2 := "nextn?date=" + today + "&count=10&include=false"

	// strTradeCalendar1 := "calendar"
	strTradeCalendar1 := "calendar"
	strTradeCalendar2 := "calendar?eyear=" + strYear1
	strTradeCalendar3 := "calendar?syear=" + strYear1 + "&eyear=" + strCurYear
	strTradeCalendar4 := "calendar?syear=" + "2005"

	kbGM1 := "gm1m?symbol=" + sym
	kbGM2 := "gm1m?symbol=" + sym + "&sdate=" + strPreMonth + "&edate=" + today
	kbGM3 := "gm1m?symbol=" + sym + "&sdate=" + strPreYear + "&edate=" + today

	strMarketInfo1 := "market_info?sec=stock"
	strMarketInfo2 := "market_info?sec=index"
	strMarketInfo3 := "market_info?sec=fund"

	strSymbolInfo1 := "symbols_info?sec=stock"
	strSymbolInfo2 := "symbols_info?sec=stock&symbols=" + sym
	strSymbolInfo3 := "symbols_info?sec=stock&symbols=" + syms
	strSymbolInfo4 := "symbols_info?sec=index"
	strSymbolInfo5 := "symbols_info?sec=index&symbols=" + idx

	strHistoryInfo1 := "history_info?symbol=" + sym
	strHistoryInfo2 := "history_info?symbol=" + sym + "&sdate=" + strPreYear
	strHistoryInfo3 := "history_info?symbol=" + sym + "&sdate=" + strPreMonth + "&edate=" + today

	//===================================================================
	return fmt.Sprintf(`<!DOCTYPE html>
<html>
<head><meta charset="UTF-8"><title>GM-API [go] -> [ %s ] </title></head>
<body>
    <h1>GM-API (go语言版本)</h1>
    <h2>服务器 : %s </h2>
    <h3>当前日期 : %s </h3>
    <ul>
        <li>说明: <a href="http://%s/usage" target="_blank">http://%s/usage</a></li>
        <li>测试1: <a href="http://%s/test" target="_blank">http://%s/test</a></li>
        <li>测试2: <a href="http://%s/test2" target="_blank">http://%s/test2</a></li>
        <li>测试3: <a href="http://%s/test3" target="_blank">http://%s/test3</a></li>
    </ul>

    <h3>交易日历</h3>
    <ul>
        <li>前N: <a href="http://%s/%s" target="_blank">http://%s/%s</a></li>
        <li>前N: <a href="http://%s/%s" target="_blank">http://%s/%s</a></li>
        <li>前N: <a href="http://%s/%s" target="_blank">http://%s/%s</a></li>
        <li>后N: <a href="http://%s/%s" target="_blank">http://%s/%s</a></li>
        <li>后N: <a href="http://%s/%s" target="_blank">http://%s/%s</a></li>
        <li>日历(当年): <a href="http://%s/%s" target="_blank">http://%s/%s</a></li>
        <li>日历(去年): <a href="http://%s/%s" target="_blank">http://%s/%s</a></li>
        <li>日历(两年): <a href="http://%s/%s" target="_blank">http://%s/%s</a></li>
        <li>日历(2005~): <a href="http://%s/%s" target="_blank">http://%s/%s</a></li>
    </ul>

    <h3>历史行情</h3>
    <ul>
		<li>综合接口 : <br>
			<a href="http://%s/%s" target="_blank">http://%s/%s</a><br>
			<a href="http://%s/%s" target="_blank">http://%s/%s</a><br>
			<a href="http://%s/%s" target="_blank">http://%s/%s</a><br>
		</li>
		<li>API接口 : <br>
		</li>
		<li>CSV接口 : <br>
		</li>
    </ul>

    <h3>基本资料</h3>
    <ul>
		<li>市场信息 : <br>
			<a href="http://%s/%s" target="_blank">http://%s/%s</a><br>
			<a href="http://%s/%s" target="_blank">http://%s/%s</a><br>
			<a href="http://%s/%s" target="_blank">http://%s/%s</a><br>
		</li>
		<li>个股信息 : <br>
			<a href="http://%s/%s" target="_blank">http://%s/%s</a><br>
			<a href="http://%s/%s" target="_blank">http://%s/%s</a><br>
			<a href="http://%s/%s" target="_blank">http://%s/%s</a><br>
			<a href="http://%s/%s" target="_blank">http://%s/%s</a><br>
			<a href="http://%s/%s" target="_blank">http://%s/%s</a><br>
		</li>
		<li>个股历史信息 : <br>
			<a href="http://%s/%s" target="_blank">http://%s/%s</a><br>
			<a href="http://%s/%s" target="_blank">http://%s/%s</a><br>
			<a href="http://%s/%s" target="_blank">http://%s/%s</a><br>
		</li>
    </ul>

    <h3>财务数据</h3>
    <ul>
        <li>vv: <a href="http://%s/%s" target="_blank">http://%s/%s</a></li>
    </ul>

</body>
</html>`,
		cfg.ServerTag, cfg.ServerTag, today,

		url, url,
		url, url,
		url, url,
		url, url,

		url, strDatePrevN1, url, strDatePrevN1,
		url, strDatePrevN2, url, strDatePrevN2,
		url, strDatePrevN3, url, strDatePrevN3,
		url, strDateNextN1, url, strDateNextN1,
		url, strDateNextN2, url, strDateNextN2,

		url, strTradeCalendar1, url, strTradeCalendar1,
		url, strTradeCalendar2, url, strTradeCalendar2,
		url, strTradeCalendar3, url, strTradeCalendar3,
		url, strTradeCalendar4, url, strTradeCalendar4,

		url, kbGM1, url, kbGM1,
		url, kbGM2, url, kbGM2,
		url, kbGM3, url, kbGM3,

		url, strMarketInfo1, url, strMarketInfo1,
		url, strMarketInfo2, url, strMarketInfo2,
		url, strMarketInfo3, url, strMarketInfo3,

		url, strSymbolInfo1, url, strSymbolInfo1,
		url, strSymbolInfo2, url, strSymbolInfo2,
		url, strSymbolInfo3, url, strSymbolInfo3,
		url, strSymbolInfo4, url, strSymbolInfo4,
		url, strSymbolInfo5, url, strSymbolInfo5,

		url, strHistoryInfo1, url, strHistoryInfo1,
		url, strHistoryInfo2, url, strHistoryInfo2,
		url, strHistoryInfo3, url, strHistoryInfo3,

		url, strSymbolInfo1, url, strSymbolInfo1,
		// url, fpathMonth1, url, fpathMonth1,
		// url, fpathYear2, url, fpathYear2,
	)
}

// 生成HTML的构造函数
func BuildHTML2(cfg HTMLConfig) string {
	url := cfg.HostURL
	// 获取当前时间
	now := time.Now()
	// 格式化当前日期为 "YYYY-MM-DD" 格式
	today := now.Format("2006-01-02")
	strCurYear := now.Format("2006")
	strCurMonth := now.Format("01")
	// strDay := now.Format("02")
	// strDate := strYear + "-" + strMonth + "-" + strDay

	strPreMonth := now.AddDate(0, -1, 0).Format("2006-01-02")
	strPreDay := now.AddDate(0, 0, -1).Format("2006-01-02")
	preYear := now.AddDate(-1, 0, 0)
	strYear1 := preYear.Format("2006")
	strPreYear := preYear.Format("2006-01-02")

	syms := "SHSE.601088,SZSE.300917"
	sym := cfg.Symbol
	idx := cfg.Sididx
	//===================================================================

	//=============================================================
	builder := strings.Builder{}
	strHead := fmt.Sprintf(`<!DOCTYPE html>
<html>
<head><meta charset="UTF-8"><title>GM-API [go] -> [ %s ] </title></head>
<body>
    <h1>GM-API (go语言版本)</h1>
    <h2>服务器 : %s </h2>
    <h3>当前日期 : %s </h3>
`, cfg.ServerTag, cfg.ServerTag, today)
	builder.WriteString(strHead)

	//=============================================================
	strTest := fmt.Sprintf(`
    <ul>
        <li>说明: <a href="http://%s/usage" target="_blank">http://%s/usage</a></li>
        <li>测试1: <a href="http://%s/test" target="_blank">http://%s/test</a></li>
        <li>测试2: <a href="http://%s/test2" target="_blank">http://%s/test2</a></li>
        <li>测试3: <a href="http://%s/test3" target="_blank">http://%s/test3</a></li>
    </ul>
`,
		url, url,
		url, url,
		url, url,
		url, url,
	)
	builder.WriteString(strTest)

	//=============================================================
	strDatePrevN1 := "prevn?date=" + today + "&count=10"
	strDatePrevN2 := "prevn?date=" + today + "&count=10&include=false"
	strDatePrevN3 := "prevn?date=" + today + "&count=365"
	strDateNextN1 := "nextn?date=" + today + "&count=10"
	strDateNextN2 := "nextn?date=" + today + "&count=10&include=false"

	strDatesList1 := "dateslist?sdate=" + strPreMonth
	strDatesList2 := "dateslist?sdate=" + strPreYear + "&edate=" + strPreMonth

	// strTradeCalendar1 := "calendar"
	strTradeCalendar1 := "calendar"
	strTradeCalendar2 := "calendar?eyear=" + strYear1
	strTradeCalendar3 := "calendar?syear=" + strYear1 + "&eyear=" + strCurYear
	strTradeCalendar4 := "calendar?syear=" + "2005"

	strCalendar := fmt.Sprintf(`
    <h3>交易日历</h3>
    <ul>
        <li>前N: <a href="http://%s/%s" target="_blank">http://%s/%s</a></li>
        <li>前N: <a href="http://%s/%s" target="_blank">http://%s/%s</a></li>
        <li>前N: <a href="http://%s/%s" target="_blank">http://%s/%s</a></li><br>
        <li>后N: <a href="http://%s/%s" target="_blank">http://%s/%s</a></li>
        <li>后N: <a href="http://%s/%s" target="_blank">http://%s/%s</a></li><br>
        <li>交易列表: <a href="http://%s/%s" target="_blank">http://%s/%s</a></li>
        <li>交易列表: <a href="http://%s/%s" target="_blank">http://%s/%s</a></li><br>
        <li>日历(当年): <a href="http://%s/%s" target="_blank">http://%s/%s</a></li>
        <li>日历(去年): <a href="http://%s/%s" target="_blank">http://%s/%s</a></li>
        <li>日历(两年): <a href="http://%s/%s" target="_blank">http://%s/%s</a></li>
        <li>日历(2005~): <a href="http://%s/%s" target="_blank">http://%s/%s</a></li>
    </ul>
`,
		url, strDatePrevN1, url, strDatePrevN1,
		url, strDatePrevN2, url, strDatePrevN2,
		url, strDatePrevN3, url, strDatePrevN3,
		url, strDateNextN1, url, strDateNextN1,
		url, strDateNextN2, url, strDateNextN2,

		url, strDatesList1, url, strDatesList1,
		url, strDatesList2, url, strDatesList2,

		url, strTradeCalendar1, url, strTradeCalendar1,
		url, strTradeCalendar2, url, strTradeCalendar2,
		url, strTradeCalendar3, url, strTradeCalendar3,
		url, strTradeCalendar4, url, strTradeCalendar4,
	)
	builder.WriteString(strCalendar)

	//=============================================================
	kbGM1 := "gm1m?symbol=" + sym
	kbGM2 := "gm1m?symbol=" + sym + "&time_stamp=true"
	kbGM3 := "gm1m?symbol=" + sym + "&sdate=" + strPreMonth + "&edate=" + today
	kbGM4 := "gm1m?symbol=" + sym + "&sdate=" + strPreYear + "&edate=" + today
	kbGM5 := "gm1m?symbol=" + sym + "&sdate=" + strPreMonth + "&edate=" + today + "&include=false"
	kbGM6 := "api1m?symbol=" + sym
	kbGM7 := "api1m?symbol=" + sym + "&time_stamp=true"
	kbGM8 := "api1m?symbol=" + sym + "&sdate=" + strPreDay + "&edate=" + today
	kbGM9 := "api1m?symbol=" + sym + "&sdate=" + strPreMonth + "&edate=" + today

	kbCSV1 := "csvyear?symbol=" + sym + "&year=" + strCurYear
	kbCSV2 := "csvyear?symbol=" + sym + "&year=" + strCurYear + "&time_stamp=true"
	kbCSV3 := "csvmonth?symbol=" + sym + "&year=" + strCurYear + "&month=" + strCurMonth
	kbCSV4 := "csv1m?symbol=" + sym + "&sdate=" + strPreMonth + "&edate=" + strPreMonth
	kbCSV5 := "csvtag?symbol=" + sym + "&sdate=" + strPreMonth + "&edate=" + today
	kbCSV6 := "csvtag?symbol=" + sym + "&sdate=" + strPreMonth + "&edate=" + today + "&tag=pe"
	kbCSV7 := "csvtag?symbol=" + sym + "&sdate=" + strPreMonth + "&edate=" + today + "&tag=vv" + "&clip=false"

	kbAPI1 := "kbars?symbols=" + sym
	kbAPI2 := "kbars?symbols=" + syms
	kbAPI3 := "kbars?symbols=" + syms + "&time_stamp=true"
	kbAPI4 := "kbars?symbols=" + syms + "&sdate=" + strPreMonth + "&edate=" + today
	kbAPI5 := "kbars?symbols=" + syms + "&sdate=" + strPreYear + "&edate=" + today + "&tag=1d"
	kbAPI6 := "kbars?symbols=" + syms + "&sdate=" + strPreYear + "&edate=" + today + "&tag=1d" + "&time_stamp=true"
	kbAPI7 := "kbarsn?symbol=" + sym + "&count=90" + "&tag=1m"
	kbAPI8 := "kbarsn?symbol=" + sym + "&count=30" + "&edate=" + today + "&tag=1d"
	kbAPI9 := "kbarsn?symbol=" + sym + "&count=30" + "&edate=" + today + "&tag=1d" + `&time_stamp=true`

	strKBars := fmt.Sprintf(`
    <h3>历史行情</h3>
    <ul>
		<li>综合接口(1m)  <br>
			<a href="http://%s/%s" target="_blank">http://%s/%s</a><br>
			<a href="http://%s/%s" target="_blank">http://%s/%s</a><br>
			<a href="http://%s/%s" target="_blank">http://%s/%s</a><br>
			<a href="http://%s/%s" target="_blank">http://%s/%s</a><br>
			<a href="http://%s/%s" target="_blank">http://%s/%s</a><br><br>
			<a href="http://%s/%s" target="_blank">http://%s/%s</a><br>
			<a href="http://%s/%s" target="_blank">http://%s/%s</a><br>
			<a href="http://%s/%s" target="_blank">http://%s/%s</a><br>
			<a href="http://%s/%s" target="_blank">http://%s/%s</a><br>
		</li>
		<br>
		<li>CSV接口  <br>
			<a href="http://%s/%s" target="_blank">http://%s/%s</a><br>
			<a href="http://%s/%s" target="_blank">http://%s/%s</a><br>
			<a href="http://%s/%s" target="_blank">http://%s/%s</a><br><br>
			<a href="http://%s/%s" target="_blank">http://%s/%s</a><br>
			<a href="http://%s/%s" target="_blank">http://%s/%s</a><br>
			<a href="http://%s/%s" target="_blank">http://%s/%s</a><br>
			<a href="http://%s/%s" target="_blank">http://%s/%s</a><br>
		</li>
		<br>
		<li>API接口  <br>
			<a href="http://%s/%s" target="_blank">http://%s/%s</a><br>
			<a href="http://%s/%s" target="_blank">http://%s/%s</a><br>
			<a href="http://%s/%s" target="_blank">http://%s/%s</a><br><br>
			<a href="http://%s/%s" target="_blank">http://%s/%s</a><br>
			<a href="http://%s/%s" target="_blank">http://%s/%s</a><br>
			<a href="http://%s/%s" target="_blank">http://%s/%s</a><br><br>
			<a href="http://%s/%s" target="_blank">http://%s/%s</a><br>
			<a href="http://%s/%s" target="_blank">http://%s/%s</a><br>
			<a href="http://%s/%s" target="_blank">http://%s/%s</a><br>
		</li>
    </ul>
`,
		url, kbGM1, url, kbGM1,
		url, kbGM2, url, kbGM2,
		url, kbGM3, url, kbGM3,
		url, kbGM4, url, kbGM4,
		url, kbGM5, url, kbGM5,
		url, kbGM6, url, kbGM6,
		url, kbGM7, url, kbGM7,
		url, kbGM8, url, kbGM8,
		url, kbGM9, url, kbGM9,

		url, kbCSV1, url, kbCSV1,
		url, kbCSV2, url, kbCSV2,
		url, kbCSV3, url, kbCSV3,
		url, kbCSV4, url, kbCSV4,
		url, kbCSV5, url, kbCSV5,
		url, kbCSV6, url, kbCSV6,
		url, kbCSV7, url, kbCSV7,

		url, kbAPI1, url, kbAPI1,
		url, kbAPI2, url, kbAPI2,
		url, kbAPI3, url, kbAPI3,
		url, kbAPI4, url, kbAPI4,
		url, kbAPI5, url, kbAPI5,
		url, kbAPI6, url, kbAPI6,
		url, kbAPI7, url, kbAPI7,
		url, kbAPI8, url, kbAPI8,
		url, kbAPI9, url, kbAPI9,
	)
	builder.WriteString(strKBars)

	//=============================================================
	strMarketInfo1 := "market_info?sec=stock"
	strMarketInfo2 := "market_info?sec=index"
	strMarketInfo3 := "market_info?sec=fund"

	strSymbolInfo1 := "symbols_info?sec=stock"
	strSymbolInfo2 := "symbols_info?sec=stock&symbols=" + sym
	strSymbolInfo3 := "symbols_info?sec=stock&symbols=" + syms
	strSymbolInfo4 := "symbols_info?sec=index"
	strSymbolInfo5 := "symbols_info?sec=index&symbols=" + idx

	strHistoryInfo1 := "history_info?symbol=" + sym
	strHistoryInfo2 := "history_info?symbol=" + sym + "&sdate=" + strPreYear
	strHistoryInfo3 := "history_info?symbol=" + sym + "&sdate=" + strPreMonth + "&edate=" + today

	strDaily1 := "daily_valuation?symbol=" + sym + "&sdate=" + strPreMonth
	strDaily2 := "daily_mktvalue?symbol=" + sym + "&sdate=" + strPreMonth
	strDaily3 := "daily_basic?symbol=" + sym + "&sdate=" + strPreMonth + "&edate=" + strPreDay
	strDaily4 := "daily_valuation_pt?symbols=" + syms + "&date=" + today
	strDaily5 := "daily_mktvalue_pt?symbols=" + syms + "&date=" + today
	strDaily6 := "daily_basic_pt?symbols=" + syms + "&date=" + today

	strShareHolder1 := "share_change?symbol=" + sym + "&sdate=" + "2020-01-01"
	strShareHolder2 := "shareholder_num?symbol=" + sym + "&sdate=" + strPreYear + "&tradable_holder=1"
	strShareHolder3 := "top_shareholder?symbol=" + sym + "&sdate=" + strPreYear
	strShareHolder4 := "ration?symbol=" + "SZSE.000728" + "&sdate=" + "2005-07-01"
	strShareHolder5 := "dvidend?symbol=" + sym + "&sdate=" + strPreYear
	strShareHolder6 := "adj_factor?symbol=" + sym + "&sdate=" + strPreYear
	strShareHolder7 := "trading_sessions?symbols=" + syms

	strInfo := fmt.Sprintf(`
    <h3>基本资料</h3>
    <ul>
		<li>市场信息  <br>
			股票：<a href="http://%s/%s" target="_blank">http://%s/%s</a><br>
			指数：<a href="http://%s/%s" target="_blank">http://%s/%s</a><br>
			基金：<a href="http://%s/%s" target="_blank">http://%s/%s</a><br>
		</li>
		<br>
		<li>个股信息  <br>
			股票：<a href="http://%s/%s" target="_blank">http://%s/%s</a><br>
			单股：<a href="http://%s/%s" target="_blank">http://%s/%s</a><br>
			多股：<a href="http://%s/%s" target="_blank">http://%s/%s</a><br>
			指数：<a href="http://%s/%s" target="_blank">http://%s/%s</a><br>
			沪指：<a href="http://%s/%s" target="_blank">http://%s/%s</a><br>
		</li>
		<br>
		<li>个股历史信息  <br>
			<a href="http://%s/%s" target="_blank">http://%s/%s</a><br>
			<a href="http://%s/%s" target="_blank">http://%s/%s</a><br>
			<a href="http://%s/%s" target="_blank">http://%s/%s</a><br>
		</li>
		<br>
		<li>个股历史数据  <br>
			<a href="http://%s/%s" target="_blank">http://%s/%s</a><br>
			<a href="http://%s/%s" target="_blank">http://%s/%s</a><br>
			<a href="http://%s/%s" target="_blank">http://%s/%s</a><br>
		</li>
		<br>
		<li>多股截面数据  <br>			
			<a href="http://%s/%s" target="_blank">http://%s/%s</a><br>
			<a href="http://%s/%s" target="_blank">http://%s/%s</a><br>
			<a href="http://%s/%s" target="_blank">http://%s/%s</a><br>
		</li>
		<br>
		<li>股东数据  <br>			
			股东变动：<a href="http://%s/%s" target="_blank">http://%s/%s</a><br>
			股东户数：<a href="http://%s/%s" target="_blank">http://%s/%s</a><br>
			十大股东：<a href="http://%s/%s" target="_blank">http://%s/%s</a><br>
			配股信息：<a href="http://%s/%s" target="_blank">http://%s/%s</a><br>
			分红送股：<a href="http://%s/%s" target="_blank">http://%s/%s</a><br>
		</li>
		<br>
		<li>其他  <br>			
			复权因子：<a href="http://%s/%s" target="_blank">http://%s/%s</a><br>
			交易时间段：<a href="http://%s/%s" target="_blank">http://%s/%s</a><br>
		</li>
    </ul>
`,
		url, strMarketInfo1, url, strMarketInfo1,
		url, strMarketInfo2, url, strMarketInfo2,
		url, strMarketInfo3, url, strMarketInfo3,

		url, strSymbolInfo1, url, strSymbolInfo1,
		url, strSymbolInfo2, url, strSymbolInfo2,
		url, strSymbolInfo3, url, strSymbolInfo3,
		url, strSymbolInfo4, url, strSymbolInfo4,
		url, strSymbolInfo5, url, strSymbolInfo5,

		url, strHistoryInfo1, url, strHistoryInfo1,
		url, strHistoryInfo2, url, strHistoryInfo2,
		url, strHistoryInfo3, url, strHistoryInfo3,

		url, strDaily1, url, strDaily1,
		url, strDaily2, url, strDaily2,
		url, strDaily3, url, strDaily3,
		url, strDaily4, url, strDaily4,
		url, strDaily5, url, strDaily5,
		url, strDaily6, url, strDaily6,

		url, strShareHolder1, url, strShareHolder1,
		url, strShareHolder2, url, strShareHolder2,
		url, strShareHolder3, url, strShareHolder3,
		url, strShareHolder4, url, strShareHolder4,
		url, strShareHolder5, url, strShareHolder5,
		url, strShareHolder6, url, strShareHolder6,
		url, strShareHolder7, url, strShareHolder7,
	)
	builder.WriteString(strInfo)

	//=============================================================
	strFinance1 := "finance_prime?symbol=" + sym + "&sdate=" + strPreYear
	strFinance2 := "finance_deriv?symbol=" + sym + "&sdate=" + strPreYear
	strFinance3 := "finance_prime_pt?symbols=" + syms + "&date=" + today
	strFinance4 := "finance_deriv_pt?symbols=" + syms + "&date=" + today

	strFundamental1 := "fundamentals_cashflow?symbol=" + sym + "&sdate=" + strPreYear
	strFundamental2 := "fundamentals_balance?symbol=" + sym + "&sdate=" + strPreYear
	strFundamental3 := "fundamentals_income?symbol=" + sym + "&sdate=" + strPreYear + "&edate=" + strPreDay
	strFundamental4 := "fundamentals_cashflow_pt?symbols=" + syms + "&date=" + today
	strFundamental5 := "fundamentals_balance_pt?symbols=" + syms + "&date=" + today
	strFundamental6 := "fundamentals_income_pt?symbols=" + syms + "&date=" + today

	strFF := fmt.Sprintf(`
    <h3>财务数据</h3>
    <ul>
		<li>财务数据  <br>
			主要指标：<a href="http://%s/%s" target="_blank">http://%s/%s</a><br>
			衍生指标：<a href="http://%s/%s" target="_blank">http://%s/%s</a><br><br>
			主要指标(pt)：<a href="http://%s/%s" target="_blank">http://%s/%s</a><br>
			衍生指标(pt)：<a href="http://%s/%s" target="_blank">http://%s/%s</a><br>
		</li>
		<br>
		<li>财务报表  <br>
			现金流量表：<a href="http://%s/%s" target="_blank">http://%s/%s</a><br>
			资产负债表：<a href="http://%s/%s" target="_blank">http://%s/%s</a><br>
			利润表：<a href="http://%s/%s" target="_blank">http://%s/%s</a><br><br>
			现金流量表(pt)：<a href="http://%s/%s" target="_blank">http://%s/%s</a><br>
			资产负债表(pt)：<a href="http://%s/%s" target="_blank">http://%s/%s</a><br>
			利润表(pt)：<a href="http://%s/%s" target="_blank">http://%s/%s</a><br>
		</li>
    </ul>
`,
		url, strFinance1, url, strFinance1,
		url, strFinance2, url, strFinance2,
		url, strFinance3, url, strFinance3,
		url, strFinance4, url, strFinance4,

		url, strFundamental1, url, strFundamental1,
		url, strFundamental2, url, strFundamental2,
		url, strFundamental3, url, strFundamental3,
		url, strFundamental4, url, strFundamental4,
		url, strFundamental5, url, strFundamental5,
		url, strFundamental6, url, strFundamental6,
	)
	builder.WriteString(strFF)

	//=============================================================
	strSecBk1 := "industry_category"
	strSecBk2 := "secotr_category" + "?sector_type=" + "1003"
	strSecBk3 := "industry_constituents" + "?industry_code=A" + "&date=" + today
	strSecBk4 := "sector_constituents" + "?sector_code=" + "007089"
	strSecBk5 := "index_constituents" + "?index=" + "SHSE.000300"
	strSecBk6 := "symbols_industry" + "?symbols=" + syms + "&date=" + today
	strSecBk7 := "symbols_sector" + "?symbols=" + syms + "&sector_type=1002"

	strAbnor1 := "abnor_change_stocks"
	strAbnor2 := "abnor_change_detail?" + "&trade_date=" + today

	strFund1 := "fnd_constituents?fund=SHSE.510880"
	strFund2 := "fnd_portfolio?fund=SHSE.510880" + "&sdate=" + strPreYear + "&report_type=1&portfolio_tpe=stk"
	strFund3 := "fnd_split?fund=SZSE.161725" + "&sdate=" + "2022-01-01" + "&edate=" + "2022-10-01"
	strFund4 := "fnd_dividend?fund=SHSE.510880" + "&sdate=" + strPreYear
	strFund5 := "fnd_netvalue?fund=SHSE.510880" + "&sdate=" + strPreYear
	strFund6 := "fnd_adj_factor?fund=SHSE.510880" + "&sdate=" + strPreYear

	strHK1 := "shszhk_quota_info?" + "&sdate=" + strPreYear
	strHK2 := "shszhk_active_stock_top10_info"
	strHK3 := "hk_inst_holding_info?symbols=" + sym + ",SZSE.001696"
	strHK4 := "hk_inst_holding_detail_info?symbols=" + sym + ",SZSE.001696"

	strOthers := fmt.Sprintf(`
    <h3>其他数据</h3>
    <ul>
		<li>行业与板块  <br>
			行业分类：<a href="http://%s/%s" target="_blank">http://%s/%s</a><br>
			板块分类：<a href="http://%s/%s" target="_blank">http://%s/%s</a><br>
			行业成分股：<a href="http://%s/%s" target="_blank">http://%s/%s</a><br>
			板块成分股：<a href="http://%s/%s" target="_blank">http://%s/%s</a><br>
			指数成分股：<a href="http://%s/%s" target="_blank">http://%s/%s</a><br>
			股票所属行业：<a href="http://%s/%s" target="_blank">http://%s/%s</a><br>
			股票所属板块：<a href="http://%s/%s" target="_blank">http://%s/%s</a><br>
		</li>
		<br>
		<li>龙虎榜  <br>
			股票：<a href="http://%s/%s" target="_blank">http://%s/%s</a><br>
			营业部：<a href="http://%s/%s" target="_blank">http://%s/%s</a><br>
		</li>
		<br>
		<li>基金数据  <br>
			成分股：<a href="http://%s/%s" target="_blank">http://%s/%s</a><br>
			资产组合：<a href="http://%s/%s" target="_blank">http://%s/%s</a><br>
			拆分折算：<a href="http://%s/%s" target="_blank">http://%s/%s</a><br><br>
			基金分红：<a href="http://%s/%s" target="_blank">http://%s/%s</a><br>
			净值数据：<a href="http://%s/%s" target="_blank">http://%s/%s</a><br>
			复权因子：<a href="http://%s/%s" target="_blank">http://%s/%s</a><br>
		</li>
		<br>
		<li>港市数据  <br>
			沪深港通额度数据：<a href="http://%s/%s" target="_blank">http://%s/%s</a><br>
			沪深港通十大活跃：<a href="http://%s/%s" target="_blank">http://%s/%s</a><br><br>
			港股机构持股数据：<a href="http://%s/%s" target="_blank">http://%s/%s</a><br>
			港股机构持股明细：<a href="http://%s/%s" target="_blank">http://%s/%s</a><br>
		</li>
    </ul>
`,
		url, strSecBk1, url, strSecBk1,
		url, strSecBk2, url, strSecBk2,
		url, strSecBk3, url, strSecBk3,
		url, strSecBk4, url, strSecBk4,
		url, strSecBk5, url, strSecBk5,
		url, strSecBk6, url, strSecBk6,
		url, strSecBk7, url, strSecBk7,

		url, strAbnor1, url, strAbnor1,
		url, strAbnor2, url, strAbnor2,

		url, strFund1, url, strFund1,
		url, strFund2, url, strFund2,
		url, strFund3, url, strFund3,
		url, strFund4, url, strFund4,
		url, strFund5, url, strFund5,
		url, strFund6, url, strFund6,

		url, strHK1, url, strHK1,
		url, strHK2, url, strHK2,
		url, strHK3, url, strHK3,
		url, strHK4, url, strHK4,
	)
	builder.WriteString(strOthers)

	//=============================================================
	strTail := `
	<br><br>
</body>
</html>`
	builder.WriteString(strTail)

	//===================================================================
	// fmt.Println(builder.String())
	return builder.String()
}
func RouteUsage(c *gin.Context) {
	hostURL := c.Request.Host

	html := BuildHTML2(HTMLConfig{
		ServerTag: serverTag,
		HostURL:   hostURL,
		Symbol:    "SHSE.601088",
		Sididx:    "SHSE.000001",
	})

	c.Data(http.StatusOK, "text/html; charset=utf-8", []byte(html))
}

func RouteTest(c *gin.Context) {
	// now := time.Now()
	// 格式化当前日期为 "YYYY-MM-DD" 格式
	// today := now.Format("2006-01-02")
	// fmt.Println(" -=> Today is", today)

	// symbols := ["AAPL", "GOOG", "AMZN", "MSFT", "TSLA"]
	// symbol := c.DefaultQuery("symbol", randomStock())
	// time := c.DefaultQuery("time", today)

	// 定义价格和成交量的范围
	pa, pr := 200, 50
	va, vr := 1000, 100

	data := gin.H{
		"Symbol": []string{randomStock(), randomStock(), randomStock(), randomStock(), randomStock()},
		"Time":   []string{randomDateInLastYear(), randomDateInLastYear(), randomDateInLastYear(), randomDateInLastYear(), randomDateInLastYear()},
		"Price":  []int{randomIntWithRange(pa, pr), randomIntWithRange(pa, pr), randomIntWithRange(pa, pr), randomIntWithRange(pa, pr), randomIntWithRange(pa, pr)},
		"Volume": []int{randomIntWithRange(va, vr), randomIntWithRange(va, vr), randomIntWithRange(va, vr), randomIntWithRange(va, vr), randomIntWithRange(va, vr)},
	}
	c.JSON(http.StatusOK, data)
}

func RouteTest2(c *gin.Context) {
	// 定义价格和成交量的范围
	pa, pr := 200, 50
	va, vr := 1000, 100

	data := gin.H{
		"columns": []string{"Symbol", "Time", "Price", "Volume"},
		"data": [][]any{
			{randomStock(), randomDateInLastYear(), randomIntWithRange(pa, pr), randomIntWithRange(va, vr)},
			{randomStock(), randomDateInLastYear(), randomIntWithRange(pa, pr), randomIntWithRange(va, vr)},
			{randomStock(), randomDateInLastYear(), randomIntWithRange(pa, pr), randomIntWithRange(va, vr)},
			{randomStock(), randomDateInLastYear(), randomIntWithRange(pa, pr), randomIntWithRange(va, vr)},
			{randomStock(), randomDateInLastYear(), randomIntWithRange(pa, pr), randomIntWithRange(va, vr)},
		},
	}

	c.JSON(http.StatusOK, data)
}

// 定义一个record结构体，包含每个字段
type record struct {
	Symbol string `json:"Symbol"`
	Time   string `json:"Time"`
	Price  int    `json:"Price"`
	Volume int    `json:"Volume"`
}

func RouteTest3(c *gin.Context) {
	// 定义价格和成交量的范围
	pa, pr := 200, 50
	va, vr := 1000, 100

	// 创建一个包含5个record的切片
	records := make([]record, 5)
	for i := range records {
		records[i] = record{
			Symbol: randomStock(),
			Time:   randomDateInLastYear(),
			Price:  randomIntWithRange(pa, pr),
			Volume: randomIntWithRange(va, vr),
		}
	}

	// 将records切片转换为JSON格式
	// jsonData, err := json.MarshalIndent(records, "", "  ")
	// if err != nil {
	// 	fmt.Println("Error converting to JSON:", err)
	// 	// return c.JSON(http.StatusOK, err.Error())
	// }

	// 输出JSON数据
	// fmt.Println(string(jsonData))

	// c.JSON(http.StatusOK, string(jsonData))
	c.JSON(http.StatusOK, records)
}

func RouteCalendar2(c *gin.Context) {
	// yy := fmt.Sprintf("%d", time.Now().Year())
	yy := time.Now().Format("2006")
	syear := c.DefaultQuery("syear", yy)
	eyear := c.DefaultQuery("eyear", yy)
	exchange := c.DefaultQuery("exchange", "")
	// timeoutSeconds := 10

	url := fmt.Sprintf("%s/get_dates_by_year", gmapi)
	pars := map[string]string{
		"syear": syear,
		"eyear": eyear,
	}
	if exchange != "" {
		pars["exchange"] = exchange
	}
	rawData, err := GetURLWithoutRetry(url, pars, 30*time.Second, 0)
	if err != nil {
		c.JSON(http.StatusNotAcceptable, gin.H{"error": err.Error()})
		return
	}

	jsonBytes, marshalErr := json.Marshal(rawData)
	if marshalErr != nil {
		c.JSON(http.StatusNotFound, marshalErr)
		return
	}

	var rcd RawColData
	unmarshalErr := json.Unmarshal(jsonBytes, &rcd)
	if unmarshalErr != nil {
		c.JSON(http.StatusNotFound, unmarshalErr)
		return
	}

	// 处理获取到的 JSON 数据为 records 形式
	records, transformErr := rcd.TransformToRecords()
	if transformErr != nil {
		c.JSON(http.StatusNotFound, transformErr)
		return
	}

	// fmt.Println("\n--- 转换后的 records 格式 (getURLWithRetry) ---")
	// recordsJSON, _ := json.MarshalIndent(records[:5], "", "  ") // 格式化输出 JSON
	// fmt.Printf("%s\n", recordsJSON)
	c.JSON(http.StatusOK, records)

}

func RouteCalendar(c *gin.Context) {
	// yy := fmt.Sprintf("%d", time.Now().Year())
	yy := time.Now().Format("2006")
	syear := c.DefaultQuery("syear", yy)
	eyear := c.DefaultQuery("eyear", yy)
	if syear > eyear {
		syear = eyear
	}

	exchange := c.DefaultQuery("exchange", "")
	timeoutSeconds := 30

	rawData, err := gm.GetCalendar(gmapi, syear, eyear, exchange, timeoutSeconds)
	if err != nil {
		c.JSON(http.StatusNotAcceptable, gin.H{" Err(gm.Calendar)": err.Error()})
		return
	}

	var rcd RawColData
	unmarshalErr := json.Unmarshal(rawData, &rcd)
	if unmarshalErr != nil {
		c.JSON(http.StatusNotFound, unmarshalErr)
		return
	}

	// 处理获取到的 JSON 数据为 records 形式
	records, transformErr := rcd.TransformToRecords()
	if transformErr != nil {
		c.JSON(http.StatusNotFound, gin.H{" Err(TransformToRecords)": transformErr.Error()})
		return
	}

	c.JSON(http.StatusOK, records)
}

func RouteDatesList(c *gin.Context) {
	sdate := c.DefaultQuery("sdate", "")
	edate := c.DefaultQuery("edate", "")

	timeoutSeconds := 30
	rawData, err := gm.GetDatesList(gmapi, sdate, edate, timeoutSeconds)
	if err != nil {
		c.JSON(http.StatusNotAcceptable, gin.H{" Err(gm.GetDatesList)": err.Error()})
		return
	}
	c.JSON(http.StatusOK, rawData)
}

func RouteDatesPrevN(c *gin.Context) {
	today := time.Now().Format("2006-01-02")
	date := c.DefaultQuery("date", "")
	if date == "" {
		date = today
	}
	scount := c.DefaultQuery("count", "10")
	count, err := strconv.Atoi(scount)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{" Err(count)": err.Error()})
		return
	}

	include := c.DefaultQuery("include", "true")
	isinclude := true
	if include == "false" {
		isinclude = false
	}

	timeoutSeconds := 30
	// rawData, err := gm.GetPrevNByte(gmapi, date, count, timeoutSeconds, isinclude)
	rawData, err := gm.GetPrevN(gmapi, date, count, isinclude, timeoutSeconds)
	if err != nil {
		c.JSON(http.StatusNotAcceptable, gin.H{" Err(gm.GetPrevN)": err.Error()})
		return
	}
	c.JSON(http.StatusOK, rawData)
}

func RouteDatesNextN(c *gin.Context) {
	today := time.Now().Format("2006-01-02")
	date := c.DefaultQuery("date", "")

	if date == "" {
		date = today
	}
	scount := c.DefaultQuery("count", "10")
	count, err := strconv.Atoi(scount)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{" Err(count)": err.Error()})
		return
	}

	include := c.DefaultQuery("include", "true")
	isinclude := true
	if include == "false" {
		isinclude = false
	}

	timeoutSeconds := 30
	// rawData, err := gm.GetNextNByte(gmapi, date, count, timeoutSeconds, isinclude)
	rawData, err := gm.GetNextN(gmapi, date, count, isinclude, timeoutSeconds)
	if err != nil {
		c.JSON(http.StatusNotAcceptable, gin.H{" Err(gm.GetNextN)": err.Error()})
		return
	}
	c.JSON(http.StatusOK, rawData)
}

func RouteCurrent(c *gin.Context) {
	symbols := c.DefaultQuery("symbols", "")
	if symbols == "" {
		c.JSON(http.StatusBadRequest, fmt.Errorf("symbols 参数为必须"))
		return
	}

	split := c.DefaultQuery("split", "false")
	issplit := false
	if split == "true" {
		issplit = true
	}

	timeoutSeconds := 30

	// rawData, err := gm.GetCurrentByte(gmapi, symbols, timeoutSeconds, issplit)
	rawData, err := gm.GetCurrent(gmapi, symbols, timeoutSeconds, issplit)
	if err != nil {
		c.JSON(http.StatusNotAcceptable, gin.H{" Err(gm.GetCurrent)": err.Error()})
		return
	}
	c.JSON(http.StatusOK, rawData)

	// // 将获取到的字符串数据解析为 JSON 格式
	// var data any
	// if err = json.Unmarshal(rawData, &data); err != nil {
	// 	c.JSON(http.StatusNotAcceptable, gin.H{" Err(unmarshaling JSON)": err.Error()})
	// 	return
	// }
	// c.JSON(http.StatusOK, data)
	// c.JSON(http.StatusOK, records)
}

func RouteDailyValuation(c *gin.Context) {
	symbols := c.DefaultQuery("symbol", "")
	if symbols == "" {
		c.JSON(http.StatusBadRequest, fmt.Errorf("symbol 参数为必须"))
		return
	}

	timeoutSeconds := 30
	sdate := c.DefaultQuery("sdate", "")
	edate := c.DefaultQuery("edate", "")
	fields := c.DefaultQuery("fields", "")

	rawData, err := gm.GetDailyValuation(gmapi, symbols, sdate, edate, fields, timeoutSeconds)
	if err != nil {
		c.JSON(http.StatusNotAcceptable, gin.H{" Err(gm.GetDailyValuation)": err.Error()})
		return
	}
	c.JSON(http.StatusOK, rawData)
}

func RouteDailyBasic(c *gin.Context) {
	symbols := c.DefaultQuery("symbol", "")
	if symbols == "" {
		c.JSON(http.StatusBadRequest, fmt.Errorf("symbol 参数为必须"))
		return
	}

	timeoutSeconds := 30
	sdate := c.DefaultQuery("sdate", "")
	edate := c.DefaultQuery("edate", "")
	fields := c.DefaultQuery("fields", "")

	rawData, err := gm.GetDailyBasic(gmapi, symbols, sdate, edate, fields, timeoutSeconds)
	if err != nil {
		c.JSON(http.StatusNotAcceptable, gin.H{" Err(gm.GetDailyBasic)": err.Error()})
		return
	}
	c.JSON(http.StatusOK, rawData)
}

func RouteDailyMktvalue(c *gin.Context) {
	symbols := c.DefaultQuery("symbol", "")
	if symbols == "" {
		c.JSON(http.StatusBadRequest, fmt.Errorf("symbol 参数为必须"))
		return
	}

	timeoutSeconds := 30
	sdate := c.DefaultQuery("sdate", "")
	edate := c.DefaultQuery("edate", "")
	fields := c.DefaultQuery("fields", "")

	rawData, err := gm.GetDailyMktvalue(gmapi, symbols, sdate, edate, fields, timeoutSeconds)
	if err != nil {
		c.JSON(http.StatusNotAcceptable, gin.H{" Err(gm.GetDailyMktvalue)": err.Error()})
		return
	}
	c.JSON(http.StatusOK, rawData)
}

func RouteFinancePrime(c *gin.Context) {
	symbols := c.DefaultQuery("symbol", "")
	if symbols == "" {
		c.JSON(http.StatusBadRequest, fmt.Errorf("symbol 参数为必须"))
		return
	}

	timeoutSeconds := 30
	sdate := c.DefaultQuery("sdate", "")
	edate := c.DefaultQuery("edate", "")
	fields := c.DefaultQuery("fields", "")
	rpt_type := c.DefaultQuery("rpt_type", "")
	data_type := c.DefaultQuery("data_type", "")

	rawData, err := gm.GetFinancePrime(gmapi, symbols, sdate, edate, fields, rpt_type, data_type, timeoutSeconds)
	if err != nil {
		c.JSON(http.StatusNotAcceptable, gin.H{" Err(gm.GetFinancePrime)": err.Error()})
		return
	}
	c.JSON(http.StatusOK, rawData)
}

func RouteFinanceDeriv(c *gin.Context) {
	symbols := c.DefaultQuery("symbol", "")
	if symbols == "" {
		c.JSON(http.StatusBadRequest, fmt.Errorf("symbol 参数为必须"))
		return
	}

	timeoutSeconds := 30
	sdate := c.DefaultQuery("sdate", "")
	edate := c.DefaultQuery("edate", "")
	fields := c.DefaultQuery("fields", "")
	rpt_type := c.DefaultQuery("rpt_type", "")
	data_type := c.DefaultQuery("data_type", "")

	rawData, err := gm.GetFinanceDeriv(gmapi, symbols, sdate, edate, fields, rpt_type, data_type, timeoutSeconds)
	if err != nil {
		c.JSON(http.StatusNotAcceptable, gin.H{" Err(gm.GetFinanceDeriv)": err.Error()})
		return
	}
	c.JSON(http.StatusOK, rawData)
}
func RouteFundamentalsCashflow(c *gin.Context) {
	symbols := c.DefaultQuery("symbol", "")
	if symbols == "" {
		c.JSON(http.StatusBadRequest, fmt.Errorf("symbol 参数为必须"))
		return
	}

	timeoutSeconds := 30
	sdate := c.DefaultQuery("sdate", "")
	edate := c.DefaultQuery("edate", "")
	fields := c.DefaultQuery("fields", "")
	rpt_type := c.DefaultQuery("rpt_type", "")
	data_type := c.DefaultQuery("data_type", "")

	rawData, err := gm.GetFundamentalsCashflow(gmapi, symbols, sdate, edate, fields, rpt_type, data_type, timeoutSeconds)
	if err != nil {
		c.JSON(http.StatusNotAcceptable, gin.H{" Err(gm.GetFundamentalsCashflow)": err.Error()})
		return
	}
	c.JSON(http.StatusOK, rawData)
}

func RouteFundamentalsIncome(c *gin.Context) {
	symbols := c.DefaultQuery("symbol", "")
	if symbols == "" {
		c.JSON(http.StatusBadRequest, fmt.Errorf("symbol 参数为必须"))
		return
	}

	timeoutSeconds := 30
	sdate := c.DefaultQuery("sdate", "")
	edate := c.DefaultQuery("edate", "")
	fields := c.DefaultQuery("fields", "")
	rpt_type := c.DefaultQuery("rpt_type", "")
	data_type := c.DefaultQuery("data_type", "")

	rawData, err := gm.GetFundamentalsIncome(gmapi, symbols, sdate, edate, fields, rpt_type, data_type, timeoutSeconds)
	if err != nil {
		c.JSON(http.StatusNotAcceptable, gin.H{" Err(gm.GetFundamentalsIncome)": err.Error()})
		return
	}
	c.JSON(http.StatusOK, rawData)
}

func RouteFundamentalsBalance(c *gin.Context) {
	symbols := c.DefaultQuery("symbol", "")
	if symbols == "" {
		c.JSON(http.StatusBadRequest, fmt.Errorf("symbol 参数为必须"))
		return
	}

	timeoutSeconds := 30
	sdate := c.DefaultQuery("sdate", "")
	edate := c.DefaultQuery("edate", "")
	fields := c.DefaultQuery("fields", "")
	rpt_type := c.DefaultQuery("rpt_type", "")
	data_type := c.DefaultQuery("data_type", "")

	rawData, err := gm.GetFundamentalsBalance(gmapi, symbols, sdate, edate, fields, rpt_type, data_type, timeoutSeconds)
	if err != nil {
		c.JSON(http.StatusNotAcceptable, gin.H{" Err(gm.GetFundamentalsBalance)": err.Error()})
		return
	}
	c.JSON(http.StatusOK, rawData)
}

func RouteFundamentalsBalancePt(c *gin.Context) {
	symbols := c.DefaultQuery("symbols", "")
	if symbols == "" {
		c.JSON(http.StatusBadRequest, fmt.Errorf("symbols 参数为必须"))
		return
	}

	timeoutSeconds := 30
	edate := c.DefaultQuery("date", "")
	fields := c.DefaultQuery("fields", "")
	rpt_type := c.DefaultQuery("rpt_type", "")
	data_type := c.DefaultQuery("data_type", "")

	rawData, err := gm.GetFundamentalsBalancePt(gmapi, symbols, edate, fields, rpt_type, data_type, timeoutSeconds)
	if err != nil {
		c.JSON(http.StatusNotAcceptable, gin.H{" Err(gm.GetFundamentalsBalancePt)": err.Error()})
		return
	}
	c.JSON(http.StatusOK, rawData)
}

func RouteFundamentalsCashflowPt(c *gin.Context) {
	symbols := c.DefaultQuery("symbols", "")
	if symbols == "" {
		c.JSON(http.StatusBadRequest, fmt.Errorf("symbols 参数为必须"))
		return
	}

	timeoutSeconds := 30
	edate := c.DefaultQuery("date", "")
	fields := c.DefaultQuery("fields", "")
	rpt_type := c.DefaultQuery("rpt_type", "")
	data_type := c.DefaultQuery("data_type", "")

	rawData, err := gm.GetFundamentalsCashflowPt(gmapi, symbols, edate, fields, rpt_type, data_type, timeoutSeconds)
	if err != nil {
		c.JSON(http.StatusNotAcceptable, gin.H{" Err(gm.GetFundamentalsCashflowPt)": err.Error()})
		return
	}
	c.JSON(http.StatusOK, rawData)
}

func RouteFundamentalsIncomePt(c *gin.Context) {
	symbols := c.DefaultQuery("symbols", "")
	if symbols == "" {
		c.JSON(http.StatusBadRequest, fmt.Errorf("symbols 参数为必须"))
		return
	}

	timeoutSeconds := 30
	edate := c.DefaultQuery("date", "")
	fields := c.DefaultQuery("fields", "")
	rpt_type := c.DefaultQuery("rpt_type", "")
	data_type := c.DefaultQuery("data_type", "")

	rawData, err := gm.GetFundamentalsIncomePt(gmapi, symbols, edate, fields, rpt_type, data_type, timeoutSeconds)
	if err != nil {
		c.JSON(http.StatusNotAcceptable, gin.H{" Err(gm.GetFundamentalsIncomePt)": err.Error()})
		return
	}
	c.JSON(http.StatusOK, rawData)
}

func RouteFinancePrimePt(c *gin.Context) {
	symbols := c.DefaultQuery("symbols", "")
	if symbols == "" {
		c.JSON(http.StatusBadRequest, fmt.Errorf("symbols 参数为必须"))
		return
	}

	timeoutSeconds := 30
	edate := c.DefaultQuery("date", "")
	fields := c.DefaultQuery("fields", "")
	rpt_type := c.DefaultQuery("rpt_type", "")
	data_type := c.DefaultQuery("data_type", "")

	rawData, err := gm.GetFinancePrimePt(gmapi, symbols, edate, fields, rpt_type, data_type, timeoutSeconds)
	if err != nil {
		c.JSON(http.StatusNotAcceptable, gin.H{" Err(gm.GetFinancePrimePt)": err.Error()})
		return
	}
	c.JSON(http.StatusOK, rawData)
}

func RouteFinanceDerivPt(c *gin.Context) {
	symbols := c.DefaultQuery("symbols", "")
	if symbols == "" {
		c.JSON(http.StatusBadRequest, fmt.Errorf("symbols 参数为必须"))
		return
	}

	timeoutSeconds := 30
	edate := c.DefaultQuery("date", "")
	fields := c.DefaultQuery("fields", "")
	rpt_type := c.DefaultQuery("rpt_type", "")
	data_type := c.DefaultQuery("data_type", "")

	rawData, err := gm.GetFinanceDerivPt(gmapi, symbols, edate, fields, rpt_type, data_type, timeoutSeconds)
	if err != nil {
		c.JSON(http.StatusNotAcceptable, gin.H{" Err(gm.GetFinanceDerivPt)": err.Error()})
		return
	}
	c.JSON(http.StatusOK, rawData)
}

func RouteDailyValuationPt(c *gin.Context) {
	symbols := c.DefaultQuery("symbols", "")
	if symbols == "" {
		c.JSON(http.StatusBadRequest, fmt.Errorf("symbols 参数为必须"))
		return
	}

	timeoutSeconds := 30
	edate := c.DefaultQuery("date", "")
	fields := c.DefaultQuery("fields", "")

	rawData, err := gm.GetDailyValuationPt(gmapi, symbols, edate, fields, timeoutSeconds)
	if err != nil {
		c.JSON(http.StatusNotAcceptable, gin.H{" Err(gm.GetDailyValuationPt)": err.Error()})
		return
	}
	c.JSON(http.StatusOK, rawData)
}

func RouteDailyBasicPt(c *gin.Context) {
	symbols := c.DefaultQuery("symbols", "")
	if symbols == "" {
		c.JSON(http.StatusBadRequest, fmt.Errorf("symbols 参数为必须"))
		return
	}

	timeoutSeconds := 30
	edate := c.DefaultQuery("date", "")
	fields := c.DefaultQuery("fields", "")

	rawData, err := gm.GetDailyBasicPt(gmapi, symbols, edate, fields, timeoutSeconds)
	if err != nil {
		c.JSON(http.StatusNotAcceptable, gin.H{" Err(gm.GetDailyBasicPt)": err.Error()})
		return
	}
	c.JSON(http.StatusOK, rawData)
}

func RouteDailyMktvaluePt(c *gin.Context) {
	symbols := c.DefaultQuery("symbols", "")
	if symbols == "" {
		c.JSON(http.StatusBadRequest, fmt.Errorf("symbols 参数为必须"))
		return
	}

	timeoutSeconds := 30
	edate := c.DefaultQuery("date", "")
	fields := c.DefaultQuery("fields", "")

	rawData, err := gm.GetDailyMktvaluePt(gmapi, symbols, edate, fields, timeoutSeconds)
	if err != nil {
		c.JSON(http.StatusNotAcceptable, gin.H{" Err(gm.GetDailyMktvaluePt)": err.Error()})
		return
	}
	c.JSON(http.StatusOK, rawData)
}

func RouteSectorCategory(c *gin.Context) {
	symbols := c.DefaultQuery("sector_type", "")
	if symbols == "" {
		c.JSON(http.StatusBadRequest, fmt.Errorf("sector_type 参数为必须"))
		return
	}

	timeoutSeconds := 30
	rawData, err := gm.GetSectorCategory(gmapi, symbols, timeoutSeconds)
	if err != nil {
		c.JSON(http.StatusNotAcceptable, gin.H{" Err(gm.GetSectorCategory)": err.Error()})
		return
	}
	c.JSON(http.StatusOK, rawData)
}

func RouteSectorConstituents(c *gin.Context) {
	symbols := c.DefaultQuery("sector_code", "")
	if symbols == "" {
		c.JSON(http.StatusBadRequest, fmt.Errorf("sector_code 参数为必须"))
		return
	}

	timeoutSeconds := 30
	rawData, err := gm.GetSectorConstituents(gmapi, symbols, timeoutSeconds)
	if err != nil {
		c.JSON(http.StatusNotAcceptable, gin.H{" Err(gm.GetSectorConstituents)": err.Error()})
		return
	}
	c.JSON(http.StatusOK, rawData)
}

func RouteSymbolsSector(c *gin.Context) {
	symbols := c.DefaultQuery("symbols", "")
	if symbols == "" {
		c.JSON(http.StatusBadRequest, fmt.Errorf("symbols 参数为必须"))
		return
	}

	timeoutSeconds := 30
	// today := time.Now().Format("2006-01-02")
	sdate := c.DefaultQuery("sector_type", "")
	// edate := c.DefaultQuery("edate", "")
	// bdate := c.DefaultQuery("bdate", "")

	rawData, err := gm.GetSymbolsSector(gmapi, symbols, sdate, timeoutSeconds)
	if err != nil {
		c.JSON(http.StatusNotAcceptable, gin.H{" Err(gm.GetSymbolsSector)": err.Error()})
		return
	}
	c.JSON(http.StatusOK, rawData)
}

func RouteDvidend(c *gin.Context) {
	symbols := c.DefaultQuery("symbol", "")
	if symbols == "" {
		c.JSON(http.StatusBadRequest, fmt.Errorf("symbol 参数为必须"))
		return
	}

	timeoutSeconds := 30
	// today := time.Now().Format("2006-01-02")
	sdate := c.DefaultQuery("sdate", "")
	edate := c.DefaultQuery("edate", "")
	// bdate := c.DefaultQuery("bdate", "")

	rawData, err := gm.GetDividend(gmapi, symbols, sdate, edate, timeoutSeconds)
	if err != nil {
		c.JSON(http.StatusNotAcceptable, gin.H{" Err(gm.GetDividend)": err.Error()})
		return
	}
	c.JSON(http.StatusOK, rawData)
}

func RouteRation(c *gin.Context) {
	symbols := c.DefaultQuery("symbol", "")
	if symbols == "" {
		c.JSON(http.StatusBadRequest, fmt.Errorf("symbol 参数为必须"))
		return
	}

	timeoutSeconds := 30
	// today := time.Now().Format("2006-01-02")
	sdate := c.DefaultQuery("sdate", "")
	edate := c.DefaultQuery("edate", "")
	// bdate := c.DefaultQuery("bdate", "")

	rawData, err := gm.GetRation(gmapi, symbols, sdate, edate, timeoutSeconds)
	if err != nil {
		c.JSON(http.StatusNotAcceptable, gin.H{" Err(gm.GetRation)": err.Error()})
		return
	}
	c.JSON(http.StatusOK, rawData)
}

func RouteShareholderNum(c *gin.Context) {
	symbols := c.DefaultQuery("symbol", "")
	if symbols == "" {
		c.JSON(http.StatusBadRequest, fmt.Errorf("symbol 参数为必须"))
		return
	}

	timeoutSeconds := 30
	// today := time.Now().Format("2006-01-02")
	sdate := c.DefaultQuery("sdate", "")
	edate := c.DefaultQuery("edate", "")
	// bdate := c.DefaultQuery("bdate", "")

	rawData, err := gm.GetShareholderNum(gmapi, symbols, sdate, edate, timeoutSeconds)
	if err != nil {
		c.JSON(http.StatusNotAcceptable, gin.H{" Err(gm.GetShareholderNum)": err.Error()})
		return
	}
	c.JSON(http.StatusOK, rawData)
}

func RouteShareChange(c *gin.Context) {
	symbols := c.DefaultQuery("symbol", "")
	if symbols == "" {
		c.JSON(http.StatusBadRequest, fmt.Errorf("symbol 参数为必须"))
		return
	}

	timeoutSeconds := 30
	// today := time.Now().Format("2006-01-02")
	sdate := c.DefaultQuery("sdate", "")
	edate := c.DefaultQuery("edate", "")
	// bdate := c.DefaultQuery("bdate", "")

	rawData, err := gm.GetShareChange(gmapi, symbols, sdate, edate, timeoutSeconds)
	if err != nil {
		c.JSON(http.StatusNotAcceptable, gin.H{" Err(gm.GetShareChange)": err.Error()})
		return
	}
	c.JSON(http.StatusOK, rawData)
}

func RouteAdjFactor(c *gin.Context) {
	symbols := c.DefaultQuery("symbol", "")
	if symbols == "" {
		c.JSON(http.StatusBadRequest, fmt.Errorf("symbol 参数为必须"))
		return
	}

	timeoutSeconds := 30
	// today := time.Now().Format("2006-01-02")
	sdate := c.DefaultQuery("sdate", "")
	edate := c.DefaultQuery("edate", "")
	bdate := c.DefaultQuery("bdate", "")

	rawData, err := gm.GetAdjFactor(gmapi, symbols, sdate, edate, bdate, timeoutSeconds)
	if err != nil {
		c.JSON(http.StatusNotAcceptable, gin.H{" Err(gm.GetAdjFactor)": err.Error()})
		return
	}
	c.JSON(http.StatusOK, rawData)
}

func RouteTopShareholder(c *gin.Context) {
	symbols := c.DefaultQuery("symbol", "")
	if symbols == "" {
		c.JSON(http.StatusBadRequest, fmt.Errorf("symbol 参数为必须"))
		return
	}

	timeoutSeconds := 30
	// today := time.Now().Format("2006-01-02")
	sdate := c.DefaultQuery("sdate", "")
	tradable_holder := c.DefaultQuery("tradable_holder", "")
	edate := c.DefaultQuery("edate", "")

	rawData, err := gm.GetTopShareholder(gmapi, symbols, sdate, edate, tradable_holder, timeoutSeconds)
	if err != nil {
		c.JSON(http.StatusNotAcceptable, gin.H{" Err(gm.GetTopShareholder)": err.Error()})
		return
	}
	c.JSON(http.StatusOK, rawData)
}

func RouteAbnorChangeStocks(c *gin.Context) {
	symbols := c.DefaultQuery("symbols", "")
	// if symbols == "" {
	// 	c.JSON(http.StatusBadRequest, fmt.Errorf("symbols 参数为必须"))
	// 	return
	// }

	timeoutSeconds := 30
	// today := time.Now().Format("2006-01-02")
	sdate := c.DefaultQuery("trade_date", "")
	fields := c.DefaultQuery("fields", "")
	change_types := c.DefaultQuery("change_types", "")
	// edate := c.DefaultQuery("edate", "")

	rawData, err := gm.GetAbnorChangeStocks(gmapi, symbols, change_types, sdate, fields, timeoutSeconds)
	if err != nil {
		c.JSON(http.StatusNotAcceptable, gin.H{" Err(gm.GetAbnorChangeStocks)": err.Error()})
		return
	}
	c.JSON(http.StatusOK, rawData)
}

func RouteAbnorChangeDetail(c *gin.Context) {
	symbols := c.DefaultQuery("symbols", "")
	// if symbols == "" {
	// 	c.JSON(http.StatusBadRequest, fmt.Errorf("symbols 参数为必须"))
	// 	return
	// }

	timeoutSeconds := 30
	// today := time.Now().Format("2006-01-02")
	sdate := c.DefaultQuery("trade_date", "")
	fields := c.DefaultQuery("fields", "")
	change_types := c.DefaultQuery("change_types", "")
	// edate := c.DefaultQuery("edate", "")

	rawData, err := gm.GetAbnorChangeDetail(gmapi, symbols, change_types, sdate, fields, timeoutSeconds)
	if err != nil {
		c.JSON(http.StatusNotAcceptable, gin.H{" Err(gm.GetAbnorChangeDetail)": err.Error()})
		return
	}
	c.JSON(http.StatusOK, rawData)
}

func RouteHKInstHoldingInfo(c *gin.Context) {
	symbols := c.DefaultQuery("symbols", "")
	// if symbols == "" {
	// 	c.JSON(http.StatusBadRequest, fmt.Errorf("symbols 参数为必须"))
	// 	return
	// }

	// today := time.Now().Format("2006-01-02")
	sdate := c.DefaultQuery("trade_date", "")
	// edate := c.DefaultQuery("edate", "")
	// count := c.DefaultQuery("count", "")
	timeoutSeconds := 30

	rawData, err := gm.GetHKInstHoldingInfo(gmapi, symbols, sdate, timeoutSeconds)
	if err != nil {
		c.JSON(http.StatusNotAcceptable, gin.H{" Err(gm.GetHKInstHoldingDetailInfo)": err.Error()})
		return
	}
	c.JSON(http.StatusOK, rawData)
}

func RouteHKInstHoldingDetailInfo(c *gin.Context) {
	symbols := c.DefaultQuery("symbols", "")
	// if symbols == "" {
	// 	c.JSON(http.StatusBadRequest, fmt.Errorf("symbols 参数为必须"))
	// 	return
	// }

	// today := time.Now().Format("2006-01-02")
	sdate := c.DefaultQuery("trade_date", "")
	// edate := c.DefaultQuery("edate", "")
	// count := c.DefaultQuery("count", "")
	timeoutSeconds := 30

	rawData, err := gm.GetHKInstHoldingDetailInfo(gmapi, symbols, sdate, timeoutSeconds)
	if err != nil {
		c.JSON(http.StatusNotAcceptable, gin.H{" Err(gm.GetHKInstHoldingDetailInfo)": err.Error()})
		return
	}
	c.JSON(http.StatusOK, rawData)
}

func RouteSHZSZHKActiveStockTop10Info(c *gin.Context) {
	symbols := c.DefaultQuery("types", "")
	// if symbols == "" {
	// 	c.JSON(http.StatusBadRequest, fmt.Errorf("types 参数为必须"))
	// 	return
	// }

	// today := time.Now().Format("2006-01-02")
	sdate := c.DefaultQuery("trade_date", "")
	// edate := c.DefaultQuery("edate", "")
	// count := c.DefaultQuery("count", "")
	timeoutSeconds := 30

	rawData, err := gm.GetSHSZHKActiveStockTop10Info(gmapi, symbols, sdate, timeoutSeconds)
	if err != nil {
		c.JSON(http.StatusNotAcceptable, gin.H{" Err(gm.GetSHSZHKActiveStockTop10Info)": err.Error()})
		return
	}
	c.JSON(http.StatusOK, rawData)
}

func RouteSHZSZHKQuotaInfo(c *gin.Context) {
	symbols := c.DefaultQuery("types", "")
	// if symbols == "" {
	// 	c.JSON(http.StatusBadRequest, fmt.Errorf("types 参数为必须"))
	// 	return
	// }

	// today := time.Now().Format("2006-01-02")
	sdate := c.DefaultQuery("sdate", "")
	edate := c.DefaultQuery("edate", "")
	count := c.DefaultQuery("count", "")
	timeoutSeconds := 30

	rawData, err := gm.GetSHSZHKQuotaInfo(gmapi, symbols, sdate, edate, count, timeoutSeconds)
	if err != nil {
		c.JSON(http.StatusNotAcceptable, gin.H{" Err(gm.GetSHSZHKQuotaInfo)": err.Error()})
		return
	}
	c.JSON(http.StatusOK, rawData)
}

func RouteFndNetValue(c *gin.Context) {
	symbols := c.DefaultQuery("fund", "")
	if symbols == "" {
		c.JSON(http.StatusBadRequest, fmt.Errorf("fund 参数为必须"))
		return
	}

	today := time.Now().Format("2006-01-02")
	sdate := c.DefaultQuery("sdate", today)
	edate := c.DefaultQuery("edate", today)
	timeoutSeconds := 30

	rawData, err := gm.GetFndNetValue(gmapi, symbols, sdate, edate, timeoutSeconds)
	if err != nil {
		c.JSON(http.StatusNotAcceptable, gin.H{" Err(gm.GetFndNetValue)": err.Error()})
		return
	}
	c.JSON(http.StatusOK, rawData)
}
func RouteFndSplit(c *gin.Context) {
	symbols := c.DefaultQuery("fund", "")
	if symbols == "" {
		c.JSON(http.StatusBadRequest, fmt.Errorf("fund 参数为必须"))
		return
	}

	today := time.Now().Format("2006-01-02")
	sdate := c.DefaultQuery("sdate", today)
	edate := c.DefaultQuery("edate", today)
	timeoutSeconds := 30

	rawData, err := gm.GetFndSplit(gmapi, symbols, sdate, edate, timeoutSeconds)
	if err != nil {
		c.JSON(http.StatusNotAcceptable, gin.H{" Err(gm.GetFndSplit)": err.Error()})
		return
	}
	c.JSON(http.StatusOK, rawData)
}

func RouteFndPortfolio(c *gin.Context) {
	symbols := c.DefaultQuery("fund", "")
	if symbols == "" {
		c.JSON(http.StatusBadRequest, fmt.Errorf("fund 参数为必须"))
		return
	}

	today := time.Now().Format("2006-01-02")
	sdate := c.DefaultQuery("sdate", today)
	edate := c.DefaultQuery("edate", today)
	report_type := c.DefaultQuery("report_type", "")
	portfolio_type := c.DefaultQuery("portfolio_type", "")
	timeoutSeconds := 30

	rawData, err := gm.GetFndPortfolio(gmapi, symbols, report_type, portfolio_type, sdate, edate, timeoutSeconds)
	if err != nil {
		c.JSON(http.StatusNotAcceptable, gin.H{" Err(gm.GetFndPortfolio)": err.Error()})
		return
	}
	c.JSON(http.StatusOK, rawData)
}

func RouteFndConstituents(c *gin.Context) {
	symbols := c.DefaultQuery("fund", "")
	if symbols == "" {
		c.JSON(http.StatusBadRequest, fmt.Errorf("fund 参数为必须"))
		return
	}

	timeoutSeconds := 30
	rawData, err := gm.GetFndConstituents(gmapi, symbols, timeoutSeconds)
	if err != nil {
		c.JSON(http.StatusNotAcceptable, gin.H{" Err(gm.GetFndPortfolio)": err.Error()})
		return
	}
	c.JSON(http.StatusOK, rawData)
}

func RouteFndDividend(c *gin.Context) {
	symbols := c.DefaultQuery("fund", "")
	if symbols == "" {
		c.JSON(http.StatusBadRequest, fmt.Errorf("fund 参数为必须"))
		return
	}

	today := time.Now().Format("2006-01-02")
	sdate := c.DefaultQuery("sdate", today)
	edate := c.DefaultQuery("edate", today)
	timeoutSeconds := 30

	rawData, err := gm.GetFndDividend(gmapi, symbols, sdate, edate, timeoutSeconds)
	if err != nil {
		c.JSON(http.StatusNotAcceptable, gin.H{" Err(gm.GetFndDividend)": err.Error()})
		return
	}
	c.JSON(http.StatusOK, rawData)
}

func RouteFndAdjFactor(c *gin.Context) {
	symbols := c.DefaultQuery("fund", "")
	if symbols == "" {
		c.JSON(http.StatusBadRequest, fmt.Errorf("fund 参数为必须"))
		return
	}

	today := time.Now().Format("2006-01-02")
	sdate := c.DefaultQuery("sdate", "")
	edate := c.DefaultQuery("edate", today)
	bdate := c.DefaultQuery("bdate", "")
	timeoutSeconds := 30

	rawData, err := gm.GetFndAdjFactor(gmapi, symbols, sdate, edate, bdate, timeoutSeconds)
	if err != nil {
		c.JSON(http.StatusNotAcceptable, gin.H{" Err(gm.GetFndAdjFactor)": err.Error()})
		return
	}
	c.JSON(http.StatusOK, rawData)
}
func RouteIndustryCategory(c *gin.Context) {
	symbols := c.DefaultQuery("source", "zjh2012")
	if symbols == "" {
		c.JSON(http.StatusBadRequest, fmt.Errorf("source 参数为必须"))
		return
	}

	// today := time.Now().Format("2006-01-02")
	sdate := c.DefaultQuery("level", "1")
	timeoutSeconds := 30

	rawData, err := gm.GetIndustryCategory(gmapi, symbols, sdate, timeoutSeconds)
	if err != nil {
		c.JSON(http.StatusNotAcceptable, gin.H{" Err(gm.GetIndustryCategory)": err.Error()})
		return
	}
	c.JSON(http.StatusOK, rawData)
}
func RouteIndustryConstituents(c *gin.Context) {
	symbols := c.DefaultQuery("industry_code", "")
	if symbols == "" {
		c.JSON(http.StatusBadRequest, fmt.Errorf("industry_code 参数为必须"))
		return
	}

	// today := time.Now().Format("2006-01-02")
	sdate := c.DefaultQuery("date", "")
	timeoutSeconds := 30

	rawData, err := gm.GetIndustryConstituents(gmapi, symbols, sdate, timeoutSeconds)
	if err != nil {
		c.JSON(http.StatusNotAcceptable, gin.H{" Err(gm.GetIndustryConstituents)": err.Error()})
		return
	}
	c.JSON(http.StatusOK, rawData)
}

func RouteSymbolsIndustry(c *gin.Context) {
	symbols := c.DefaultQuery("symbols", "")
	if symbols == "" {
		c.JSON(http.StatusBadRequest, fmt.Errorf("symbols 参数为必须"))
		return
	}

	// today := time.Now().Format("2006-01-02")
	level := c.DefaultQuery("level", "")
	source := c.DefaultQuery("source", "")
	sdate := c.DefaultQuery("date", "")
	timeoutSeconds := 30

	rawData, err := gm.GetSymbolIndustry(gmapi, symbols, source, level, sdate, timeoutSeconds)
	if err != nil {
		c.JSON(http.StatusNotAcceptable, gin.H{" Err(gm.GetSymbolIndustry)": err.Error()})
		return
	}
	c.JSON(http.StatusOK, rawData)
}

func RouteIndexConstituents(c *gin.Context) {
	symbols := c.DefaultQuery("index", "")
	if symbols == "" {
		c.JSON(http.StatusBadRequest, fmt.Errorf("index 参数为必须"))
		return
	}

	today := time.Now().Format("2006-01-02")
	sdate := c.DefaultQuery("trade_date", today)
	// edate := c.DefaultQuery("edate", today)
	// sec := c.DefaultQuery("sec", "stock")
	timeoutSeconds := 30

	rawData, err := gm.GetIndexConstituents(gmapi, symbols, sdate, timeoutSeconds)
	if err != nil {
		c.JSON(http.StatusNotAcceptable, gin.H{" Err(gm.GetIndexConstituents)": err.Error()})
		return
	}
	c.JSON(http.StatusOK, rawData)
}

func RouteTradingSessions(c *gin.Context) {
	symbols := c.DefaultQuery("symbols", "")
	if symbols == "" {
		c.JSON(http.StatusBadRequest, fmt.Errorf("symbols 参数为必须"))
		return
	}

	// today := time.Now().Format("2006-01-02")
	// sdate := c.DefaultQuery("sdate", today)
	// edate := c.DefaultQuery("edate", today)
	// sec := c.DefaultQuery("sec", "stock")
	timeoutSeconds := 30

	rawData, err := gm.GetTradingSessions(gmapi, symbols, timeoutSeconds)
	if err != nil {
		c.JSON(http.StatusNotAcceptable, gin.H{" Err(gm.GetTradingSessions)": err.Error()})
		return
	}
	c.JSON(http.StatusOK, rawData)
}
func RouteMarketInfo(c *gin.Context) {
	symbols := c.DefaultQuery("symbols", "")
	// if symbols == "" {
	// 	c.JSON(http.StatusBadRequest, fmt.Errorf("symbols 参数为必须"))
	// 	return
	// }

	// today := time.Now().Format("2006-01-02")
	// sdate := c.DefaultQuery("sdate", today)
	// edate := c.DefaultQuery("edate", today)
	sec := c.DefaultQuery("sec", "stock")
	exchange := c.DefaultQuery("exchange", "")
	timeoutSeconds := 30

	rawData, err := gm.GetMarketInfo(gmapi, symbols, sec, exchange, timeoutSeconds)
	if err != nil {
		c.JSON(http.StatusNotAcceptable, gin.H{" Err(gm.GetMarketInfo)": err.Error()})
		return
	}
	c.JSON(http.StatusOK, rawData)
}

func RouteSymbolsInfo(c *gin.Context) {
	symbols := c.DefaultQuery("symbols", "")
	// if symbols == "" {
	// 	c.JSON(http.StatusBadRequest, fmt.Errorf("symbols 参数为必须"))
	// 	return
	// }

	// today := time.Now().Format("2006-01-02")
	// sdate := c.DefaultQuery("sdate", today)
	// edate := c.DefaultQuery("edate", today)
	sec := c.DefaultQuery("sec", "stock")
	exchange := c.DefaultQuery("exchange", "")
	trade_date := c.DefaultQuery("trade_date", "")
	timeoutSeconds := 30

	rawData, err := gm.GetSymbolsInfo(gmapi, symbols, sec, exchange, trade_date, timeoutSeconds)
	if err != nil {
		c.JSON(http.StatusNotAcceptable, gin.H{" Err(gm.GetSymbolsInfo)": err.Error()})
		return
	}
	c.JSON(http.StatusOK, rawData)
}

func RouteHistoryInfo(c *gin.Context) {
	symbol := c.DefaultQuery("symbol", "")
	if symbol == "" {
		c.JSON(http.StatusBadRequest, fmt.Errorf("symbol 参数为必须"))
		return
	}

	today := time.Now().Format("2006-01-02")
	sdate := c.DefaultQuery("sdate", today)
	edate := c.DefaultQuery("edate", today)
	timeoutSeconds := 30

	rawData, err := gm.GetHistoryInfo(gmapi, symbol, sdate, edate, timeoutSeconds)
	if err != nil {
		c.JSON(http.StatusNotAcceptable, gin.H{" Err(gm.GetHistoryInfo)": err.Error()})
		return
	}
	c.JSON(http.StatusOK, rawData)
}

func RouteGMApi1m(c *gin.Context) {
	symbols := c.DefaultQuery("symbol", "")
	if symbols == "" {
		c.JSON(http.StatusBadRequest, fmt.Errorf("symbols 参数为必须"))
		return
	}

	timeoutSeconds := 30
	today := time.Now().Format("2006-01-02")
	sdate := c.DefaultQuery("sdate", today)
	edate := c.DefaultQuery("edate", today)
	// tag := c.DefaultQuery("tag", "1m")

	timestamp := c.DefaultQuery("time_stamp", "false")
	istimestamp := false
	if timestamp == "true" {
		istimestamp = true
	}
	datesList, _ := gm.GetDatesList(gmapi, sdate, edate, timeoutSeconds)
	if len(datesList) == 0 {
		c.JSON(http.StatusNotAcceptable, gin.H{" Err(gm.GetDatesList)": "日期列表为空 " + sdate + "~" + edate})
		return
	}
	rawData, err := gm.Get1mByDatelist(gmapi, symbols, datesList, istimestamp, timeoutSeconds)
	if err != nil {
		c.JSON(http.StatusNotAcceptable, gin.H{" Err(gm.Get1mByDatelist)": err.Error()})
		return
	}
	c.JSON(http.StatusOK, rawData)
}

func RouteKbars(c *gin.Context) {
	symbols := c.DefaultQuery("symbols", "")
	if symbols == "" {
		c.JSON(http.StatusBadRequest, fmt.Errorf("symbols 参数为必须"))
		return
	}

	today := time.Now().Format("2006-01-02")
	sdate := c.DefaultQuery("sdate", today)
	edate := c.DefaultQuery("edate", today)
	tag := c.DefaultQuery("tag", "1m")
	timestamp := c.DefaultQuery("time_stamp", "false")
	timeoutSeconds := 30

	istimestamp := false
	if timestamp == "true" {
		istimestamp = true
	}
	// rawData, err := gm.GetKbarsHisByte(gmapi, symbols, tag, sdate, edate, timeoutSeconds)
	rawData, err := gm.GetKbarsHis(gmapi, symbols, tag, sdate, edate, istimestamp, timeoutSeconds)
	if err != nil {
		c.JSON(http.StatusNotAcceptable, gin.H{" Err(gm.GetKbarsHis)": err.Error()})
		return
	}
	c.JSON(http.StatusOK, rawData)

	// var rcd RawColData
	// if unmarshalErr := json.Unmarshal(rawData, &rcd); unmarshalErr != nil {
	// 	c.JSON(http.StatusNotFound, gin.H{" Err(unmarshalErr)": unmarshalErr.Error()})
	// 	return
	// }

	// records, transformErr := rcd.TransformToRecords()
	// if transformErr != nil {
	// 	c.JSON(http.StatusInternalServerError, gin.H{" Err(transformErr)": transformErr.Error()})
	// 	return
	// }

	// c.JSON(http.StatusOK, records)
}

// 获取股票当日的 K 线数据，返回字典json格式，键为代码，值为K线数据
func RouteKBDict(c *gin.Context) {
	symbols := c.DefaultQuery("symbols", "")
	if symbols == "" {
		c.JSON(http.StatusBadRequest, fmt.Errorf("symbols 参数为必须"))
		return
	}

	today := time.Now().Format("2006-01-02")
	sdate := c.DefaultQuery("sdate", today)
	edate := c.DefaultQuery("edate", today)
	tag := c.DefaultQuery("tag", "1m")
	timestamp := c.DefaultQuery("time_stamp", "false")
	timeoutSeconds := 30

	istimestamp := false
	if timestamp == "true" {
		istimestamp = true
	}

	rawData, err := gm.GetKbarsHis(gmapi, symbols, tag, sdate, edate, istimestamp, timeoutSeconds)
	if err != nil {
		c.JSON(http.StatusNotAcceptable, gin.H{" Err(gm.GetKbarsHis)": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gm.ConvertRecords2Dict(rawData))
}

// 获取股票当日的 K 线数据，返回字典json格式，键为代码，值为K线数据
func RouteKBDictTS(c *gin.Context) {
	symbols := c.DefaultQuery("symbols", "")
	if symbols == "" {
		c.JSON(http.StatusBadRequest, fmt.Errorf("symbols 参数为必须"))
		return
	}

	today := time.Now().Format("2006-01-02")
	sdate := c.DefaultQuery("sdate", today)
	edate := c.DefaultQuery("edate", today)
	tag := c.DefaultQuery("tag", "1m")
	timestamp := c.DefaultQuery("time_stamp", "false")
	timeoutSeconds := 30

	istimestamp := false
	if timestamp == "true" {
		istimestamp = true
	}

	rawData, err := gm.GetKbarsHis(gmapi, symbols, tag, sdate, edate, istimestamp, timeoutSeconds)
	if err != nil {
		c.JSON(http.StatusNotAcceptable, gin.H{" Err(gm.GetKbarsHis)": err.Error()})
		return
	}
	if istimestamp {
		c.JSON(http.StatusOK, gm.ConvertRecords2DictTSInt(rawData))
	} else {
		c.JSON(http.StatusOK, gm.ConvertRecords2DictTSString(rawData))
	}
	// c.JSON(http.StatusOK, gm.ConvertRecords2DictTSString(rawData))
}

func RouteKbarsN(c *gin.Context) {
	symbol := c.DefaultQuery("symbol", "")
	if symbol == "" {
		c.JSON(http.StatusBadRequest, fmt.Errorf("symbols 参数为必须"))
		return
	}

	today := time.Now().Format("2006-01-02")
	edate := c.DefaultQuery("edate", today)
	count := c.DefaultQuery("count", "")
	tag := c.DefaultQuery("tag", "1d")
	timestamp := c.DefaultQuery("time_stamp", "false")
	timeoutSeconds := 30

	istimestamp := false
	if timestamp == "true" {
		istimestamp = true
	}

	// rawData, err := gm.GetKbarsHisNByte(gmapi, symbol, tag, count, edate, timeoutSeconds)
	rawData, err := gm.GetKbarsHisN(gmapi, symbol, tag, count, edate, istimestamp, timeoutSeconds)
	if err != nil {
		c.JSON(http.StatusNotAcceptable, gin.H{" Err(gm.GetKbarsHisN)": err.Error()})
		return
	}
	c.JSON(http.StatusOK, rawData)
}

func RouteCSVxzMonth(c *gin.Context) {
	symbol := c.DefaultQuery("symbol", "")
	if symbol == "" {
		c.JSON(http.StatusBadRequest, fmt.Errorf("symbol 参数为必须参数"))
		return
	}

	timeoutSeconds := 30
	now := time.Now()
	year, month, _ := now.Date()
	yearStr := fmt.Sprintf("%d", year)
	monthStr := fmt.Sprintf("%02d", month)

	cyear := c.DefaultQuery("year", yearStr)
	cmonth := c.DefaultQuery("month", monthStr)
	// tag := c.DefaultQuery("tag", "1m")
	imonth, _ := strconv.Atoi(cmonth)
	iyear, _ := strconv.Atoi(cyear)
	timestamp := c.DefaultQuery("time_stamp", "false")
	istimestamp := false
	if timestamp == "true" {
		istimestamp = true
	}

	// rawData, err := gm.GetCSVMonthJson(gmcsv, symbol, imonth, iyear, timeoutSeconds)
	rawData, err := gm.GetCSVMonth(gmcsv, symbol, imonth, iyear, istimestamp, timeoutSeconds)
	if err != nil {
		c.JSON(http.StatusNotAcceptable, gin.H{" Err(gm.GetCSVMonth)": err.Error()})
		return
	}
	// c.JSON(http.StatusOK, string(rawData))
	c.JSON(http.StatusOK, rawData)

	// 将获取到的字符串数据解析为 JSON 格式
	// var data any
	// if err = json.Unmarshal(rawData, &data); err != nil {
	// 	c.JSON(http.StatusNotAcceptable, gin.H{" Err(unmarshaling JSON)": err.Error()})
	// 	return
	// }
	// c.JSON(http.StatusOK, data)
}

func RouteCSVxzYear(c *gin.Context) {
	symbol := c.DefaultQuery("symbol", "")
	if symbol == "" {
		c.JSON(http.StatusBadRequest, fmt.Errorf("symbol 参数为必须参数"))
		return
	}

	timeoutSeconds := 30
	now := time.Now()
	year := now.Year()
	yearStr := fmt.Sprintf("%d", year)

	cyear := c.DefaultQuery("year", yearStr)
	tag := c.DefaultQuery("tag", "1m")
	timestamp := c.DefaultQuery("time_stamp", "false")
	istimestamp := false
	if timestamp == "true" {
		istimestamp = true
	}
	iyear, _ := strconv.Atoi(cyear)

	lookuptab := map[string]string{
		"1m": "timestamp",
		"vv": "timestamp",
		"pe": "trade_date",
	}

	rawData, err := gm.GetCSVYear(gmcsv, symbol, tag, iyear, istimestamp, lookuptab[tag], timeoutSeconds)
	if err != nil {
		c.JSON(http.StatusNotAcceptable, gin.H{" Err(gm.GetCSVYear)": err.Error()})
		return
	}
	// c.JSON(http.StatusOK, string(rawData))
	c.JSON(http.StatusOK, rawData)

	// 将获取到的字符串数据解析为 JSON 格式
	// var data any
	// if err = json.Unmarshal(rawData, &data); err != nil {
	// 	c.JSON(http.StatusNotAcceptable, gin.H{" Err(unmarshaling JSON)": err.Error()})
	// 	return
	// }
	// c.JSON(http.StatusOK, data)
}

func RouteCSVxz1m(c *gin.Context) {
	symbol := c.DefaultQuery("symbol", "")
	if symbol == "" {
		c.JSON(http.StatusBadRequest, fmt.Errorf("symbol 参数为必须参数"))
		return
	}

	timeoutSeconds := 30
	now := time.Now()
	today := now.Format("2006-01-02")

	sdate := c.DefaultQuery("sdate", today)
	edate := c.DefaultQuery("edate", today)
	// tag := c.DefaultQuery("tag", "1m")

	timestamp := c.DefaultQuery("time_stamp", "false")
	istimestamp := false
	if timestamp == "true" {
		istimestamp = true
	}

	clip := c.DefaultQuery("clip", "true")
	isclip := false
	if clip == "true" {
		isclip = true
	}

	rawData, err := gm.GetCSV1m(gmcsv, symbol, sdate, edate, istimestamp, isclip, timeoutSeconds)
	if err != nil {
		c.JSON(http.StatusNotAcceptable, gin.H{" Err(gm.GetCSV1m)": err.Error()})
		return
	}
	c.JSON(http.StatusOK, rawData)

	// // jsonData, err := json.MarshalIndent(result, "", "  ")
	// jsonData, err := json.Marshal(rawData)
	// if err != nil {
	// 	c.JSON(http.StatusNotAcceptable, gin.H{" Err(unmarshaling JSON)": err.Error()})
	// }

	// c.JSON(http.StatusOK, string(jsonData))
}

func RouteCSVxzTag(c *gin.Context) {
	symbol := c.DefaultQuery("symbol", "")
	if symbol == "" {
		c.JSON(http.StatusBadRequest, fmt.Errorf("symbol 参数为必须参数"))
		return
	}

	timeoutSeconds := 30
	now := time.Now()
	today := now.Format("2006-01-02")

	sdate := c.DefaultQuery("sdate", today)
	edate := c.DefaultQuery("edate", today)
	tag := c.DefaultQuery("tag", "vv")

	timestamp := c.DefaultQuery("time_stamp", "false")
	istimestamp := false
	if timestamp == "true" {
		istimestamp = true
	}

	clip := c.DefaultQuery("clip", "true")
	isclip := false
	if clip == "true" {
		isclip = true
	}

	rawData, err := gm.GetCSVTag(gmcsv, tag, symbol, sdate, edate, istimestamp, isclip, timeoutSeconds)
	if err != nil {
		c.JSON(http.StatusNotAcceptable, gin.H{" Err(gm.GetCSVTag)": err.Error()})
		return
	}
	c.JSON(http.StatusOK, rawData)
}

func RouteGM1m(c *gin.Context) {
	symbol := c.DefaultQuery("symbol", "")
	if symbol == "" {
		c.JSON(http.StatusBadRequest, fmt.Errorf("symbol 参数为必须参数"))
		return
	}

	timeoutSeconds := 30
	now := time.Now()
	today := now.Format("2006-01-02")
	prday := now.AddDate(0, 0, -1)
	yesterday := prday.Format("2006-01-02")

	// tag := c.DefaultQuery("tag", "1m")

	timestamp := c.DefaultQuery("time_stamp", "false")
	istimestamp := false
	if timestamp == "true" {
		istimestamp = true
	}
	include := c.DefaultQuery("include", "true")
	isinclude := true
	if include == "false" {
		isinclude = false
	}
	cday := yesterday
	if isinclude && gm.IsAOpen() {
		cday = today
	}
	sdate := c.DefaultQuery("sdate", cday)
	edate := c.DefaultQuery("edate", cday)

	rawData, err := gm.GetGM1m(gmcsv, gmapi, symbol, sdate, edate, istimestamp, isinclude, timeoutSeconds)
	if err != nil {
		c.JSON(http.StatusNotAcceptable, gin.H{" Err(gm.GetGM1m)": err.Error()})
		return
	}
	c.JSON(http.StatusOK, rawData)
}
