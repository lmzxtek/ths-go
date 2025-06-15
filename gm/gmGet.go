package gm

import (
	"context"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/ulikunitz/xz"
)

// makeRequest 发起一个HTTP请求并打印响应状态和响应体
// func makeRequest(url string) {
// 	resp, err := http.Get(url)
// 	if err != nil {
// 		fmt.Println("Error fetching URL: ", err)
// 		return
// 	}
// 	defer resp.Body.Close()

// 	body, err := io.ReadAll(resp.Body)
// 	if err != nil {
// 		fmt.Println("Error reading response body: ", err)
// 		return
// 	}

// 	fmt.Println("Response Status: ", resp.Status)
// 	fmt.Println("Response Body: ", string(body))
// }

// fetchData 从指定URL获取数据
func FetchData(url string) ([]byte, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("请求失败: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("HTTP 请求失败，状态码: %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("读取响应体失败: %v", err)
	}

	return body, nil
}

// 获取URL数据
func fetchURLData(url string, timeout time.Duration, params map[string]string) ([]byte, error) {
	client := &http.Client{
		Timeout: timeout,
	}

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	q := req.URL.Query()
	for k, v := range params {
		q.Add(k, v)
	}
	req.URL.RawQuery = q.Encode()

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("请求失败: %s", resp.Status)
	}

	return io.ReadAll(resp.Body)
}

// 发送带重试的 GET 请求
// 参数：
//   - url: 请求的 URL
//   - timeout: 单次请求超时时间（time.Duration 类型，如 5*time.Second）
//   - maxRetries: 最大重试次数（如 3 次）
//
// 返回值：
//   - 响应内容或错误信息
func fetchWithRetry(url string, timeout time.Duration, maxRetries int, params map[string]string) ([]byte, error) {
	var lastErr error

	// 重试循环（总尝试次数 = maxRetries + 1）
	for i := 0; i <= maxRetries; i++ {
		// 创建带超时的 Context
		ctx, cancel := context.WithTimeout(context.Background(), timeout)
		defer cancel()

		// 创建 Request 并绑定 Context
		req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
		if err != nil {
			lastErr = fmt.Errorf("创建请求失败: %v", err)
			continue
		}
		q := req.URL.Query()
		for k, v := range params {
			q.Add(k, v)
		}
		req.URL.RawQuery = q.Encode()

		// 发送请求
		client := &http.Client{}
		resp, err := client.Do(req)
		if err != nil {
			// 记录错误，准备重试
			lastErr = fmt.Errorf("请求失败 (尝试 %d/%d): %v", i+1, maxRetries+1, err)

			// 如果达到最大重试次数，返回错误
			if i == maxRetries {
				return nil, lastErr
			}

			// 等待一段时间后重试（指数退避）
			sleepTime := time.Duration(i*i) * time.Second // 示例：二次方退避
			fmt.Printf("等待 %v 后重试...\n", sleepTime)
			time.Sleep(sleepTime)
			continue
		}
		defer resp.Body.Close()

		// 检查 HTTP 状态码
		if resp.StatusCode != http.StatusOK {
			lastErr = fmt.Errorf("状态码异常 (尝试 %d/%d): %d", i+1, maxRetries+1, resp.StatusCode)

			if i == maxRetries {
				return nil, lastErr
			}

			// 等待后重试
			sleepTime := time.Duration(i*i) * time.Second
			fmt.Printf("等待 %v 后重试...\n", sleepTime)
			time.Sleep(sleepTime)
			continue
		}

		// 读取响应内容（此处简化处理）
		return io.ReadAll(resp.Body)
	}

	return nil, fmt.Errorf("所有尝试均失败: %w", lastErr)
}

// 测试数据1
func GetTest(url string, timeoutSeconds int) ([]byte, error) {
	urlTar := fmt.Sprintf("%s/test", url)

	// 获取历史K线数据
	resp, err := fetchURLData(urlTar, time.Duration(timeoutSeconds)*time.Second, map[string]string{})
	if err != nil {
		fmt.Printf("获取数据失败: %s\n", err)
		return nil, err
	}

	return resp, nil
}

// 测试数据1
func GetTest2(url string, timeoutSeconds int) ([]byte, error) {
	urlTar := fmt.Sprintf("%s/test2", url)

	// 获取历史K线数据
	resp, err := fetchURLData(urlTar, time.Duration(timeoutSeconds)*time.Second, map[string]string{})
	if err != nil {
		fmt.Printf("获取数据失败: %s\n", err)
		return nil, err
	}

	return resp, nil
}

// 获取交易日历
func GetCalendar(gmapi string, syear string, eyear string, exchange string, timeoutSeconds int) ([]byte, error) {
	url := fmt.Sprintf("%s/get_dates_by_year", gmapi)
	params := map[string]string{
		"syear": syear,
		"eyear": eyear,
	}
	if exchange != "" {
		params["exchange"] = exchange
	}

	// 获取历史K线数据
	resp, err := fetchURLData(url, time.Duration(timeoutSeconds)*time.Second, params)
	if err != nil {
		// fmt.Printf(" 获取数据失败(gm.GetCalendar): %v\n", err)
		return nil, err
	}

	return resp, nil
}

// 查询指定日期的前n个交易日
func GetPrevNByte(gmapi string, date string, count int, timeoutSeconds int, include bool) ([]byte, error) {
	url := fmt.Sprintf("%s/get_dates_prev_n", gmapi)

	cdate := date
	if include {
		nxd, _ := time.Parse("2006-01-02", cdate)
		cdate = nxd.AddDate(0, 0, 1).Format("2006-01-02")
	}

	params := map[string]string{
		"date":  cdate,
		"count": fmt.Sprintf("%d", count),
	}

	resp, err := fetchURLData(url, time.Duration(timeoutSeconds)*time.Second, params)
	if err != nil {
		// fmt.Printf(" 获取数据失败(gm.GetCalendar): %v\n", err)
		return nil, err
	}

	return resp, nil
}

// 查询指定日期的前n个交易日
func GetPrevN(gmapi string, date string, count int, timeoutSeconds int, include bool) ([]any, error) {

	rawData, err := GetPrevNByte(gmapi, date, count, timeoutSeconds, include)
	if err != nil {
		return nil, err
	}

	// 将获取到的字符串数据解析为 JSON 格式
	var data []any
	if err = json.Unmarshal(rawData, &data); err != nil {
		return nil, fmt.Errorf("解析 JSON 数据失败: %v", err)
	}

	return data, nil
}

// 查询指定日期的后n个交易日
func GetNextNByte(gmapi string, date string, count int, timeoutSeconds int, include bool) ([]byte, error) {
	url := fmt.Sprintf("%s/get_dates_next_n", gmapi)

	cdate := date
	if include {
		prd, _ := time.Parse("2006-01-02", cdate)
		cdate = prd.AddDate(0, 0, -1).Format("2006-01-02")
	}

	params := map[string]string{
		"date":  cdate,
		"count": fmt.Sprintf("%d", count),
	}

	resp, err := fetchURLData(url, time.Duration(timeoutSeconds)*time.Second, params)
	if err != nil {
		// fmt.Printf(" 获取数据失败(gm.GetCalendar): %v\n", err)
		return nil, err
	}

	return resp, nil
}

// 查询指定日期的前n个交易日
func GetNextN(gmapi string, date string, count int, timeoutSeconds int, include bool) ([]any, error) {

	rawData, err := GetNextNByte(gmapi, date, count, timeoutSeconds, include)
	if err != nil {
		return nil, err
	}

	// 将获取到的字符串数据解析为 JSON 格式
	var data []any
	if err = json.Unmarshal(rawData, &data); err != nil {
		return nil, fmt.Errorf("解析 JSON 数据失败: %v", err)
	}

	return data, nil
}

// 获取行情快照数据
func GetCurrentByte(gmapi string, symbols string, timeoutSeconds int, split bool) ([]byte, error) {
	url := fmt.Sprintf("%s/get_current", gmapi)
	params := map[string]string{
		"symbols": symbols,
	}
	if split {
		params["split"] = "true"
	} else {
		params["split"] = "false"
	}

	resp, err := fetchURLData(url, time.Duration(timeoutSeconds)*time.Second, params)
	if err != nil {
		// fmt.Printf("获取数据失败: %s\n", err)
		return nil, err
	}

	return resp, nil
}

// 获取行情快照数据
func GetCurrent(gmapi string, symbols string, timeoutSeconds int, split bool) ([]any, error) {
	rawData, err := GetCurrentByte(gmapi, symbols, timeoutSeconds, split)
	if err != nil {
		return nil, fmt.Errorf("获取数据失败(GetCurrentByte()): %v", err)
	}

	// 将获取到的字符串数据解析为 JSON 格式
	var data []any
	if err = json.Unmarshal(rawData, &data); err != nil {
		return nil, fmt.Errorf("解析 JSON 数据失败(GetCurrent()): %v", err)
	}
	return data, nil
}

// 请求市场个股列表(gm-api)
// 输入参数：
//   - exchange: 交易所代码: SHSE,SZSE,CFFEX,DCE,CZCE,SHFE,INE
//   - sec: 证券类型代码, 如 "stock"、"fund"、"index"
//
// 返回值：
//   - 市场个股列表
func GetMarketInfo(gmapi string,
	symbols string, sec string, exchange string,
	timeoutSeconds int) ([]map[string]any, error) {

	url := fmt.Sprintf("%s/get_infos", gmapi)
	params := map[string]string{
		"split": "true",
	}

	if exchange != "" {
		params["exchange"] = exchange
	}
	if symbols != "" {
		params["symbols"] = symbols
	}
	if sec != "" {
		params["sec"] = sec
	}

	resp, err := fetchURLData(url, time.Duration(timeoutSeconds)*time.Second, params)
	if err != nil {
		return nil, err
	}

	var rcd RawColData
	if unmarshalErr := json.Unmarshal(resp, &rcd); unmarshalErr != nil {
		return nil, fmt.Errorf("解析 JSON 数据失败(GetMarketInfo()): %v", unmarshalErr)
	}

	records, transformErr := rcd.ToRecords()
	if transformErr != nil {
		return nil, fmt.Errorf("转换数据失败(GetMarketInfo()): %v", transformErr)
	}

	return ConvertEob2Timestamp(records, false), nil
}

// 请求市场股票列表某交易日的交易数据(gm-api)
// 输入参数：
//   - exchange: 交易所代码: SHSE,SZSE,CFFEX,DCE,CZCE,SHFE,INE
//   - sec: 证券类型代码, 如 "stock"、"fund"、"index"
//
// 返回值：
//   - 市场个股列表
func GetSymbolsInfo(gmapi string,
	symbols string, sec string, exchange string, trade_date string,
	timeoutSeconds int) ([]map[string]any, error) {

	url := fmt.Sprintf("%s/get_symbols", gmapi)
	params := map[string]string{
		"split": "true",
	}

	if exchange != "" {
		params["exchange"] = exchange
	}
	if symbols != "" {
		params["symbols"] = symbols
	}
	if sec != "" {
		params["sec"] = sec
	}
	if trade_date != "" {
		params["trade_date"] = trade_date
	}

	resp, err := fetchURLData(url, time.Duration(timeoutSeconds)*time.Second, params)
	if err != nil {
		return nil, err
	}

	var rcd RawColData
	if unmarshalErr := json.Unmarshal(resp, &rcd); unmarshalErr != nil {
		return nil, fmt.Errorf("解析 JSON 数据失败(GetSymbolsInfo()): %v", unmarshalErr)
	}

	records, transformErr := rcd.ToRecords()
	if transformErr != nil {
		return nil, fmt.Errorf("转换数据失败(GetSymbolsInfo()): %v", transformErr)
	}

	return ConvertEob2Timestamp(records, false), nil
}

// 查询个股估值指标每日数据
// 输入参数：
//   - sdate: 开始日期, 格式: "2021-01-01"
//   - edate: 结束日期, 格式: "2021-01-01"
//   - symbol: 股票代码, 如 "SHSE.601088"
func GetDailyValuation(gmapi string,
	symbol string, sdate string, edate string, fields string,
	timeoutSeconds int) ([]map[string]any, error) {
	if symbol == "" {
		return nil, fmt.Errorf("Symbol为必须字段")
	}

	url := fmt.Sprintf("%s/get_daily_valuation", gmapi)
	params := map[string]string{
		"symbol": symbol,
		"split":  "true",
	}

	if sdate != "" {
		params["sdate"] = sdate
	}
	if edate != "" {
		params["edate"] = edate
	}
	if fields != "" {
		params["fields"] = fields
	}
	resp, err := fetchURLData(url, time.Duration(timeoutSeconds)*time.Second, params)
	if err != nil {
		return nil, err
	}

	var rcd RawColData
	if unmarshalErr := json.Unmarshal(resp, &rcd); unmarshalErr != nil {
		return nil, fmt.Errorf("解析 JSON 数据失败(get_daily_valuation()): %v", unmarshalErr)
	}

	records, transformErr := rcd.ToRecords()
	if transformErr != nil {
		return nil, fmt.Errorf("转换数据失败(get_daily_valuation()): %v", transformErr)
	}

	return ConvertEob2Timestamp(records, false), nil
}

// 查询个股基础指标每日数据
// 输入参数：
//   - sdate: 开始日期, 格式: "2021-01-01"
//   - edate: 结束日期, 格式: "2021-01-01"
//   - symbol: 股票代码, 如 "SHSE.601088"
func GetDailyBasic(gmapi string,
	symbol string, sdate string, edate string, fields string,
	timeoutSeconds int) ([]map[string]any, error) {
	if symbol == "" {
		return nil, fmt.Errorf("Symbol为必须字段")
	}

	url := fmt.Sprintf("%s/get_daily_basic", gmapi)
	params := map[string]string{
		"symbol": symbol,
		"split":  "true",
	}

	if sdate != "" {
		params["sdate"] = sdate
	}
	if edate != "" {
		params["edate"] = edate
	}
	if fields != "" {
		params["fields"] = fields
	}
	resp, err := fetchURLData(url, time.Duration(timeoutSeconds)*time.Second, params)
	if err != nil {
		return nil, err
	}

	var rcd RawColData
	if unmarshalErr := json.Unmarshal(resp, &rcd); unmarshalErr != nil {
		return nil, fmt.Errorf("解析 JSON 数据失败(get_daily_basic()): %v", unmarshalErr)
	}

	records, transformErr := rcd.ToRecords()
	if transformErr != nil {
		return nil, fmt.Errorf("转换数据失败(get_daily_basic()): %v", transformErr)
	}

	return ConvertEob2Timestamp(records, false), nil
}

// 查询个股市值指标每日数据
// 输入参数：
//   - sdate: 开始日期, 格式: "2021-01-01"
//   - edate: 结束日期, 格式: "2021-01-01"
//   - symbol: 股票代码, 如 "SHSE.601088"
func GetDailyMktvalue(gmapi string,
	symbol string, sdate string, edate string, fields string,
	timeoutSeconds int) ([]map[string]any, error) {
	if symbol == "" {
		return nil, fmt.Errorf("Symbol为必须字段")
	}

	url := fmt.Sprintf("%s/get_daily_mktvalue", gmapi)
	params := map[string]string{
		"symbol": symbol,
		"split":  "true",
	}

	if sdate != "" {
		params["sdate"] = sdate
	}
	if edate != "" {
		params["edate"] = edate
	}
	if fields != "" {
		params["fields"] = fields
	}
	resp, err := fetchURLData(url, time.Duration(timeoutSeconds)*time.Second, params)
	if err != nil {
		return nil, err
	}

	var rcd RawColData
	if unmarshalErr := json.Unmarshal(resp, &rcd); unmarshalErr != nil {
		return nil, fmt.Errorf("解析 JSON 数据失败(get_daily_mktvalue()): %v", unmarshalErr)
	}

	records, transformErr := rcd.ToRecords()
	if transformErr != nil {
		return nil, fmt.Errorf("转换数据失败(get_daily_mktvalue()): %v", transformErr)
	}

	return ConvertEob2Timestamp(records, false), nil
}

// 查询财务主要指标
// 输入参数：
//   - sdate: 开始日期, 格式: "2021-01-01"
//   - edate: 结束日期, 格式: "2021-01-01"
//   - symbol: 股票代码, 如 "SHSE.601088"
//   - fields: fields不能超过20个字段
//   - rpt_type: 按报告期查询可指定以下报表类型： 1-一季度报; 6-中报; 9-前三季报; 12-年报 默认None为不限
//   - data_type: 在发布原始财务报告以后，上市公司可能会对数据进行修正。 101-合并原始; 102-合并调整; 201-母公司原始; 202-母公司调整 默认None返回当期合并调整，如果没有调整返回合并原始
func GetFinancePrime(gmapi string,
	symbol string, sdate string, edate string, fields string,
	rpt_type string, data_type string,
	timeoutSeconds int) ([]map[string]any, error) {
	if symbol == "" {
		return nil, fmt.Errorf("Symbol为必须字段")
	}

	url := fmt.Sprintf("%s/get_finance_prime", gmapi)
	params := map[string]string{
		"symbol": symbol,
		"split":  "true",
	}

	if sdate != "" {
		params["sdate"] = sdate
	}
	if edate != "" {
		params["edate"] = edate
	}
	if rpt_type != "" {
		params["rpt_type"] = rpt_type
	}
	if data_type != "" {
		params["data_type"] = data_type
	}
	if fields != "" {
		params["fields"] = fields
	}
	resp, err := fetchURLData(url, time.Duration(timeoutSeconds)*time.Second, params)
	if err != nil {
		return nil, err
	}

	var rcd RawColData
	if unmarshalErr := json.Unmarshal(resp, &rcd); unmarshalErr != nil {
		return nil, fmt.Errorf("解析 JSON 数据失败(get_finance_prime()): %v", unmarshalErr)
	}

	records, transformErr := rcd.ToRecords()
	if transformErr != nil {
		return nil, fmt.Errorf("转换数据失败(get_finance_prime()): %v", transformErr)
	}

	return ConvertEob2Timestamp(records, false), nil
}

// 查询财务衍生指标
// 输入参数：
//   - sdate: 开始日期, 格式: "2021-01-01"
//   - edate: 结束日期, 格式: "2021-01-01"
//   - symbol: 股票代码, 如 "SHSE.601088"
//   - fields: fields不能超过20个字段
//   - rpt_type: 按报告期查询可指定以下报表类型： 1-一季度报; 6-中报; 9-前三季报; 12-年报 默认None为不限
//   - data_type: 在发布原始财务报告以后，上市公司可能会对数据进行修正。 101-合并原始; 102-合并调整; 201-母公司原始; 202-母公司调整 默认None返回当期合并调整，如果没有调整返回合并原始
func GetFinanceDeriv(gmapi string,
	symbol string, sdate string, edate string, fields string,
	rpt_type string, data_type string,
	timeoutSeconds int) ([]map[string]any, error) {
	if symbol == "" {
		return nil, fmt.Errorf("Symbol为必须字段")
	}

	url := fmt.Sprintf("%s/get_finance_deriv", gmapi)
	params := map[string]string{
		"symbol": symbol,
		"split":  "true",
	}

	if sdate != "" {
		params["sdate"] = sdate
	}
	if edate != "" {
		params["edate"] = edate
	}
	if rpt_type != "" {
		params["rpt_type"] = rpt_type
	}
	if data_type != "" {
		params["data_type"] = data_type
	}
	if fields != "" {
		params["fields"] = fields
	}
	resp, err := fetchURLData(url, time.Duration(timeoutSeconds)*time.Second, params)
	if err != nil {
		return nil, err
	}

	var rcd RawColData
	if unmarshalErr := json.Unmarshal(resp, &rcd); unmarshalErr != nil {
		return nil, fmt.Errorf("解析 JSON 数据失败(get_finance_deriv()): %v", unmarshalErr)
	}

	records, transformErr := rcd.ToRecords()
	if transformErr != nil {
		return nil, fmt.Errorf("转换数据失败(get_finance_deriv()): %v", transformErr)
	}

	return ConvertEob2Timestamp(records, false), nil
}

// 查询财务主要指标单日截面数据(point-in-time, 多标的)
// 输入参数：
//   - date: 交易日期, 格式: "2021-01-01"
//   - fields: fields不能超过20个字段
//   - rpt_type: 按报告期查询可指定以下报表类型： 1-一季度报; 6-中报; 9-前三季报; 12-年报 默认None为不限
//   - data_type: 在发布原始财务报告以后，上市公司可能会对数据进行修正。 101-合并原始; 102-合并调整; 201-母公司原始; 202-母公司调整 默认None返回当期合并调整，如果没有调整返回合并原始
//   - symbols: 股票列表, 如 "SHSE.601088,SZSE.000001"
func GetFinancePrimePt(gmapi string,
	symbols string, date string, fields string,
	rpt_type string, data_type string,
	timeoutSeconds int) ([]map[string]any, error) {
	if symbols == "" {
		return nil, fmt.Errorf("Symbols为必须字段")
	}

	url := fmt.Sprintf("%s/get_finance_prime_pt", gmapi)
	params := map[string]string{
		"symbols": symbols,
		"split":   "true",
	}

	if date != "" {
		params["date"] = date
	}
	if rpt_type != "" {
		params["rpt_type"] = rpt_type
	}
	if data_type != "" {
		params["data_type"] = data_type
	}
	if fields != "" {
		params["fields"] = fields
	}
	resp, err := fetchURLData(url, time.Duration(timeoutSeconds)*time.Second, params)
	if err != nil {
		return nil, err
	}

	var rcd RawColData
	if unmarshalErr := json.Unmarshal(resp, &rcd); unmarshalErr != nil {
		return nil, fmt.Errorf("解析 JSON 数据失败(get_finance_prime_pt()): %v", unmarshalErr)
	}

	records, transformErr := rcd.ToRecords()
	if transformErr != nil {
		return nil, fmt.Errorf("转换数据失败(get_finance_prime_pt()): %v", transformErr)
	}

	return ConvertEob2Timestamp(records, false), nil
}

// 查询财务主要指标单日截面数据(point-in-time, 多标的)
// 输入参数：
//   - date: 交易日期, 格式: "2021-01-01"
//   - fields: fields不能超过20个字段
//   - rpt_type: 按报告期查询可指定以下报表类型： 1-一季度报; 6-中报; 9-前三季报; 12-年报 默认None为不限
//   - data_type: 在发布原始财务报告以后，上市公司可能会对数据进行修正。 101-合并原始; 102-合并调整; 201-母公司原始; 202-母公司调整 默认None返回当期合并调整，如果没有调整返回合并原始
//   - symbols: 股票列表, 如 "SHSE.601088,SZSE.000001"
func GetFinanceDerivPt(gmapi string,
	symbols string, date string, fields string,
	rpt_type string, data_type string,
	timeoutSeconds int) ([]map[string]any, error) {
	if symbols == "" {
		return nil, fmt.Errorf("Symbols为必须字段")
	}

	url := fmt.Sprintf("%s/get_finance_deriv_pt", gmapi)
	params := map[string]string{
		"symbols": symbols,
		"split":   "true",
	}

	if date != "" {
		params["date"] = date
	}
	if rpt_type != "" {
		params["rpt_type"] = rpt_type
	}
	if data_type != "" {
		params["data_type"] = data_type
	}
	if fields != "" {
		params["fields"] = fields
	}
	resp, err := fetchURLData(url, time.Duration(timeoutSeconds)*time.Second, params)
	if err != nil {
		return nil, err
	}

	var rcd RawColData
	if unmarshalErr := json.Unmarshal(resp, &rcd); unmarshalErr != nil {
		return nil, fmt.Errorf("解析 JSON 数据失败(get_finance_deriv_pt()): %v", unmarshalErr)
	}

	records, transformErr := rcd.ToRecords()
	if transformErr != nil {
		return nil, fmt.Errorf("转换数据失败(get_finance_deriv_pt()): %v", transformErr)
	}

	return ConvertEob2Timestamp(records, false), nil
}

// 查询估值指标单日截面数据(point-in-time, 多标的)
// 输入参数：
//   - date: 交易日期, 格式: "2021-01-01"
//   - symbols: 股票列表, 如 "SHSE.601088,SZSE.000001"
func GetDailyBasicPt(gmapi string,
	symbols string, date string, fields string,
	timeoutSeconds int) ([]map[string]any, error) {
	if symbols == "" {
		return nil, fmt.Errorf("Symbols为必须字段")
	}

	url := fmt.Sprintf("%s/get_daily_basic_pt", gmapi)
	params := map[string]string{
		"symbols": symbols,
		"split":   "true",
	}

	if date != "" {
		params["date"] = date
	}
	if fields != "" {
		params["fields"] = fields
	}
	resp, err := fetchURLData(url, time.Duration(timeoutSeconds)*time.Second, params)
	if err != nil {
		return nil, err
	}

	var rcd RawColData
	if unmarshalErr := json.Unmarshal(resp, &rcd); unmarshalErr != nil {
		return nil, fmt.Errorf("解析 JSON 数据失败(get_daily_basic_pt()): %v", unmarshalErr)
	}

	records, transformErr := rcd.ToRecords()
	if transformErr != nil {
		return nil, fmt.Errorf("转换数据失败(get_daily_basic_pt()): %v", transformErr)
	}

	return ConvertEob2Timestamp(records, false), nil
}

// 查询估值指标单日截面数据(point-in-time, 多标的)
// 输入参数：
//   - date: 交易日期, 格式: "2021-01-01"
//   - symbols: 股票列表, 如 "SHSE.601088,SZSE.000001"
func GetDailyMktvaluePt(gmapi string,
	symbols string, date string, fields string,
	timeoutSeconds int) ([]map[string]any, error) {
	if symbols == "" {
		return nil, fmt.Errorf("Symbols为必须字段")
	}

	url := fmt.Sprintf("%s/get_daily_mktvalue_pt", gmapi)
	params := map[string]string{
		"symbols": symbols,
		"split":   "true",
	}

	if date != "" {
		params["date"] = date
	}
	if fields != "" {
		params["fields"] = fields
	}
	resp, err := fetchURLData(url, time.Duration(timeoutSeconds)*time.Second, params)
	if err != nil {
		return nil, err
	}

	var rcd RawColData
	if unmarshalErr := json.Unmarshal(resp, &rcd); unmarshalErr != nil {
		return nil, fmt.Errorf("解析 JSON 数据失败(get_daily_mktvalue_pt()): %v", unmarshalErr)
	}

	records, transformErr := rcd.ToRecords()
	if transformErr != nil {
		return nil, fmt.Errorf("转换数据失败(get_daily_mktvalue_pt()): %v", transformErr)
	}

	return ConvertEob2Timestamp(records, false), nil
}

// 查询估值指标单日截面数据(point-in-time, 多标的)
// 输入参数：
//   - date: 交易日期, 格式: "2021-01-01"
//   - symbols: 股票列表, 如 "SHSE.601088,SZSE.000001"
func GetDailyValuationPt(gmapi string,
	symbols string, date string, fields string,
	timeoutSeconds int) ([]map[string]any, error) {
	if symbols == "" {
		return nil, fmt.Errorf("Symbols为必须字段")
	}

	url := fmt.Sprintf("%s/get_daily_valuation_pt", gmapi)
	params := map[string]string{
		"symbols": symbols,
		"split":   "true",
	}

	if date != "" {
		params["date"] = date
	}
	if fields != "" {
		params["fields"] = fields
	}
	resp, err := fetchURLData(url, time.Duration(timeoutSeconds)*time.Second, params)
	if err != nil {
		return nil, err
	}

	var rcd RawColData
	if unmarshalErr := json.Unmarshal(resp, &rcd); unmarshalErr != nil {
		return nil, fmt.Errorf("解析 JSON 数据失败(get_daily_valuation_pt()): %v", unmarshalErr)
	}

	records, transformErr := rcd.ToRecords()
	if transformErr != nil {
		return nil, fmt.Errorf("转换数据失败(get_daily_valuation_pt()): %v", transformErr)
	}

	return ConvertEob2Timestamp(records, false), nil
}

// 请求个股历史交易段的个股信息(gm-api)
// 输入参数：
//   - sdate: 开始日期, 格式: "2021-01-01"
//   - edate: 结束日期, 格式: "2021-01-01"
//   - symbol: 股票代码, 如 "SHSE.601088"
func GetHistoryInfo(gmapi string,
	symbol string, sdate string, edate string,
	timeoutSeconds int) ([]map[string]any, error) {
	if symbol == "" {
		return nil, fmt.Errorf("Symbol为必须字段")
	}

	url := fmt.Sprintf("%s/get_his_symbol", gmapi)
	params := map[string]string{
		"symbol": symbol,
		"split":  "true",
	}

	if sdate != "" {
		params["sdate"] = sdate
	}
	if edate != "" {
		params["edate"] = edate
	}
	resp, err := fetchURLData(url, time.Duration(timeoutSeconds)*time.Second, params)
	if err != nil {
		return nil, err
	}

	var rcd RawColData
	if unmarshalErr := json.Unmarshal(resp, &rcd); unmarshalErr != nil {
		return nil, fmt.Errorf("解析 JSON 数据失败(GetHistoryInfo()): %v", unmarshalErr)
	}

	records, transformErr := rcd.ToRecords()
	if transformErr != nil {
		return nil, fmt.Errorf("转换数据失败(GetKbarsHisByte()): %v", transformErr)
	}

	return ConvertEob2Timestamp(records, false), nil
}

// 获取K线行情数据
func GetKbarsHisByte(gmapi string,
	symbols string, tag string,
	sdate string, edate string,
	timeoutSeconds int) ([]byte, error) {

	url := fmt.Sprintf("%s/get_his", gmapi)
	params := map[string]string{
		"symbols": symbols,
		"tag":     tag,
		"sdate":   sdate,
		"edate":   edate,
	}

	resp, err := fetchURLData(url, time.Duration(timeoutSeconds)*time.Second, params)
	if err != nil {
		// fmt.Printf("获取数据失败: %s\n", err)
		return nil, err
	}

	return resp, nil
}

// 获取K线行情数据
func GetKbarsHis(gmapi string,
	symbols string, tag string,
	sdate string, edate string, istimestamp bool,
	timeoutSeconds int) ([]map[string]any, error) {

	rawData, err := GetKbarsHisByte(gmapi, symbols, tag, sdate, edate, timeoutSeconds)
	if err != nil {
		return nil, fmt.Errorf("获取数据失败(GetKbarsHisByte()): %v", err)
	}

	var rcd RawColData
	if unmarshalErr := json.Unmarshal(rawData, &rcd); unmarshalErr != nil {
		return nil, fmt.Errorf("解析 JSON 数据失败(GetKbarsHisByte()): %v", unmarshalErr)
	}

	records, transformErr := rcd.ToRecords()
	if transformErr != nil {
		return nil, fmt.Errorf("转换数据失败(GetKbarsHisByte()): %v", transformErr)
	}

	return ConvertEob2Timestamp(records, istimestamp), nil
}

// 获取K线行情数据
func GetKbarsHisNByte(gmapi string,
	symbol string, tag string,
	count string, edate string,
	timeoutSeconds int) ([]byte, error) {

	url := fmt.Sprintf("%s/get_his_n", gmapi)
	params := map[string]string{
		"symbol": symbol,
		"tag":    tag,
		"edate":  edate,
	}
	if count != "" {
		params["count"] = count
	}

	resp, err := fetchURLData(url, time.Duration(timeoutSeconds)*time.Second, params)
	if err != nil {
		// fmt.Printf("获取数据失败: %s\n", err)
		return nil, err
	}

	return resp, nil
}

// 获取K线行情数据
func GetKbarsHisN(gmapi string,
	symbol string, tag string,
	count string, edate string, istimestamp bool,
	timeoutSeconds int) ([]map[string]any, error) {

	rawData, err := GetKbarsHisNByte(gmapi, symbol, tag, count, edate, timeoutSeconds)
	if err != nil {
		return nil, fmt.Errorf("获取数据失败(GetKbarsHisNByte()): %v", err)
	}

	var rcd RawColData
	if unmarshalErr := json.Unmarshal(rawData, &rcd); unmarshalErr != nil {
		return nil, fmt.Errorf("解析 JSON 数据失败(GetKbarsHisN()): %v", unmarshalErr)
	}

	records, transformErr := rcd.ToRecords()
	if transformErr != nil {
		return nil, fmt.Errorf("转换数据失败(GetKbarsHisN()): %v", transformErr)
	}

	return ConvertEob2Timestamp(records, istimestamp), nil
}

func getFilePathYear(symbol string, tag string, year int) string {
	// 构造行情数据文件路径
	key := fmt.Sprintf("%s-%s", symbol[:2], symbol[5:7])
	var subfld string
	if tag == "vv" || tag == "pe" {
		subfld = fmt.Sprintf("kbars-%s/%s-%d/%s-%d--%s/", tag, tag, year, tag, year, key)
	} else {
		subfld = fmt.Sprintf("kbars-year/year-%d/year-%d--%s/", year, year, key)
	}
	fname := fmt.Sprintf("kbars-%s--%s--%d-.csv.xz", tag, symbol, year)
	fpath := fmt.Sprintf("%s%s", subfld, fname)
	// fpath := filepath.Join(subfld, fname)
	return fpath
}

func getFilePathMonth(symbol string, year int, month int) string {
	// 构造分时行情文件路径
	tag := "1m"
	key := fmt.Sprintf("%s-%s", symbol[:2], symbol[5:7])
	subfld := fmt.Sprintf("kbars-month/month-%d/month-%d-%02d--%s/", year, year, month, key)
	fname := fmt.Sprintf("kbars-%s--%s--%d-%02d-.csv.xz", tag, symbol, year, month)
	fpath := fmt.Sprintf("%s%s", subfld, fname)
	// fpath := filepath.Join(subfld, fname)
	return fpath
}

// // downloadAndReadData 从指定URL下载数据并返回内容
// func downloadAndReadData(url string) ([]byte, error) {
// 	// 发送HTTP GET请求
// 	resp, err := http.Get(url)
// 	if err != nil {
// 		return nil, fmt.Errorf("请求失败: %w", err)
// 	}
// 	defer resp.Body.Close()

// 	// 检查HTTP状态码
// 	if resp.StatusCode != http.StatusOK {
// 		return nil, fmt.Errorf("HTTP请求失败，状态码: %d", resp.StatusCode)
// 	}

// 	// 由于文件是.xz格式（LZMA压缩），需要解压缩
// 	reader, err := lzma.NewReader(resp.Body)
// 	if err != nil {
// 		return nil, fmt.Errorf("创建LZMA解压器失败: %w", err)
// 	}

// 	// 读取解压后的数据
// 	data, err := io.ReadAll(reader)
// 	if err != nil {
// 		return nil, fmt.Errorf("读取数据失败: %w", err)
// 	}

// 	return data, nil
// }

// downloadAndReadData 从指定URL下载数据并返回内容
func downloadAndReadData(url string) ([]byte, error) {
	// 发送HTTP GET请求
	resp, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("请求失败: %w", err)
	}
	defer resp.Body.Close()

	// 检查HTTP状态码
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("HTTP请求失败，状态码: %d", resp.StatusCode)
	}

	// 使用github.com/ulikunitz/xz库创建XZ解压器
	reader, err := xz.NewReader(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("创建XZ解压器失败: %w", err)
	}

	// 读取解压后的数据
	data, err := io.ReadAll(reader)
	if err != nil {
		return nil, fmt.Errorf("读取数据失败: %w", err)
	}

	return data, nil
}

// parseValue attempts to convert a string value to its appropriate Go type.
func parseValue(s string) any {
	// Try parsing as integer
	if i, err := strconv.ParseInt(s, 10, 64); err == nil {
		return i
	}

	// Try parsing as float
	if f, err := strconv.ParseFloat(s, 64); err == nil {
		return f
	}

	// Try parsing as boolean (case-insensitive)
	lowerS := strings.ToLower(s)
	if lowerS == "true" {
		return true
	}
	if lowerS == "false" {
		return false
	}

	// If none of the above, return as string
	return s
}

// csvToJSON 将CSV数据转换为JSON格式(不进行类型判断)
func CsvToJSON(csvData []byte) ([]byte, error) {
	// 创建CSV读取器
	reader := csv.NewReader(strings.NewReader(string(csvData)))

	// 读取所有CSV记录
	records, err := reader.ReadAll()
	if err != nil {
		return nil, fmt.Errorf("解析CSV失败: %w", err)
	}

	if len(records) == 0 {
		return nil, fmt.Errorf("CSV文件为空")
	}

	// 第一行作为表头
	headers := records[0]

	// 将每一行数据转换为map
	var result []map[string]any
	for i := 1; i < len(records); i++ {
		row := make(map[string]any)
		for j, value := range records[i] {
			if j < len(headers) {
				row[headers[j]] = value
			}
		}
		result = append(result, row)
	}

	// 转换为JSON
	// jsonData, err := json.MarshalIndent(result, "", "  ")
	jsonData, err := json.Marshal(result)
	if err != nil {
		return nil, fmt.Errorf("转换JSON失败: %w", err)
	}

	return jsonData, nil
}

func CSVToRecords(csvData []byte, istimestamp bool, tskey string) ([]map[string]any, error) {

	reader := csv.NewReader(strings.NewReader(string(csvData)))

	// 读取所有CSV记录
	records, err := reader.ReadAll()
	if err != nil {
		return nil, fmt.Errorf("解析CSV失败: %w", err)
	}

	if len(records) == 0 {
		return nil, fmt.Errorf("CSV文件为空")
	}

	// 第一行作为表头
	headers := records[0]
	// 将每一行数据转换为map
	var result []map[string]any
	for i := 1; i < len(records); i++ {
		row := make(map[string]any)
		for j, value := range records[i] {
			if j < len(headers) {
				// row[headers[j]] = value
				parsedValue := parseValue(value)
				head := headers[j]
				if istimestamp && (head == tskey || head == "timestamp") {
					// 转换为时间戳
					// parsedValue,err = string.
					// tt := ConvertString2Time(parsedValue.(string))
					tt, err := ParseTimestamp(parsedValue.(string))
					if err != nil {
						// return nil, fmt.Errorf("转换时间戳失败: %w", err)
						continue // 跳过错误数据
					}
					row[head] = tt.UnixMilli()
				} else {
					row[head] = parsedValue
				}
			}
		}
		result = append(result, row)
	}

	return result, nil
}

// CSVToJson reads a CSV file and converts its content to a JSON array of objects.
// It attempts to infer data types for each field (int, float, bool, string).
func CSVToJson(csvData []byte, istimestamp bool, tskey string) ([]byte, error) {

	result, err := CSVToRecords(csvData, istimestamp, tskey)
	if err != nil {
		return nil, fmt.Errorf("解析CSV失败: %w", err)
	}

	// jsonData, err := json.MarshalIndent(result, "", "  ")
	jsonData, err := json.Marshal(result)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal JSON: %w", err)
	}

	return jsonData, nil
}

// downloadAndConvertToJSON 从URL下载CSV数据并转换为JSON
func DownloadAndConvertToJSON(url string, istimestamp bool, tskey string) ([]byte, error) {
	// 下载并读取数据
	csvData, err := downloadAndReadData(url)
	if err != nil {
		return nil, err
	}

	// 转换为JSON
	// jsonData, err := csvToJSON(csvData)
	jsonData, err := CSVToJson(csvData, istimestamp, tskey)
	if err != nil {
		return nil, err
	}

	return jsonData, nil
}

// 获取Csv.xz按月行情数据
func GetCSVMonthJson(gmcsv string,
	symbol string,
	month int, year int, istimestamp bool,
	timeoutSeconds int) ([]byte, error) {

	url := fmt.Sprintf("%s/download/", gmcsv)

	fpath := getFilePathMonth(symbol, year, month)
	// fmt.Println(fpath)

	resp, err := DownloadAndConvertToJSON(url+fpath, istimestamp, "timestamp")
	if err != nil {
		// fmt.Printf("获取数据失败: %s\n", err)
		return nil, err
	}
	return resp, nil
}

// 获取Csv.xz按年行情数据
func GetCSVYearJson(gmcsv string,
	symbol string, tag string, year int, istimestamp bool, tskey string,
	timeoutSeconds int) ([]byte, error) {

	url := fmt.Sprintf("%s/download/", gmcsv)

	fpath := getFilePathYear(symbol, tag, year)
	// fmt.Println(fpath)

	resp, err := DownloadAndConvertToJSON(url+fpath, istimestamp, tskey)
	if err != nil {
		// fmt.Printf("获取数据失败: %s\n", err)
		return nil, err
	}
	return resp, nil
}

// 获取Csv.xz按月行情数据
func GetCSVMonth(gmcsv string,
	symbol string,
	month int, year int, istimestamp bool,
	timeoutSeconds int) ([]map[string]any, error) {

	url := fmt.Sprintf("%s/download/", gmcsv)
	fpath := getFilePathMonth(symbol, year, month)
	// fmt.Println(fpath)

	// 下载并读取数据
	csvData, err := downloadAndReadData(url + fpath)
	if err != nil {
		return nil, err
	}
	result, err := CSVToRecords(csvData, istimestamp, "timestamp")
	if err != nil {
		return nil, fmt.Errorf("解析CSV失败: %w", err)
	}
	return result, nil
}

// 获取Csv.xz按年行情数据
func GetCSVYear(gmcsv string,
	symbol string, tag string, year int, istimestamp bool, tskey string,
	timeoutSeconds int) ([]map[string]any, error) {

	url := fmt.Sprintf("%s/download/", gmcsv)
	fpath := getFilePathYear(symbol, tag, year)
	// fmt.Println(url + fpath)

	// 下载并读取数据
	csvData, err := downloadAndReadData(url + fpath)
	if err != nil {
		fmt.Println(url + fpath)
		return nil, err
	}
	result, err := CSVToRecords(csvData, istimestamp, tskey)
	if err != nil {
		fmt.Println("解析CSV失败: %w", err)
		return nil, fmt.Errorf("解析CSV失败: %w", err)
	}
	return result, nil
}

// filterDataByDate 根据sdate和edate筛选数据，日期列为timestamp
// 参数:
//
//	data: 原始数据，[]map[string]any
//	dateKey: 数据中表示日期的键名，例如 "timestamp"
//	sdateStr: 开始日期字符串，例如 "2023-01-01"
//	edateStr: 结束日期字符串，例如 "2023-12-31"
//
// 返回:
//
//	筛选后的数据，[]map[string]any
//	错误信息，如果日期解析失败
func filterDataByDate(data []map[string]any, dateKey, sdateStr, edateStr string) ([]map[string]any, error) {
	// 定义日期解析格式 (根据实际情况调整，例如 "2006-01-02 15:04:05" 或 "2006/01/02")
	const dateFormat = "2006-01-02"

	// 解析开始日期
	stime, err := time.Parse(dateFormat, sdateStr)
	if err != nil {
		return nil, fmt.Errorf("解析开始日期失败: %w", err)
	}
	stime = GetDayStart(stime)

	// 解析结束日期
	etime, err := time.Parse(dateFormat, edateStr)
	if err != nil {
		return nil, fmt.Errorf("解析结束日期失败: %w", err)
	}
	etime = GetDayEnd(etime)

	var filteredData []map[string]any

	for _, item := range data {
		timestampVal, ok := item[dateKey]
		if !ok {
			// 如果没有日期键，跳过或根据需求处理
			continue
		}

		itemTime, parseErr := ParseTimestamp(timestampVal)
		if parseErr != nil {
			// 如果解析失败，根据需求跳过或报错
			fmt.Printf("警告: 无法解析日期 '%s': %v\n", timestampVal, parseErr)
			continue
		}
		// 判断是否在范围内
		if (itemTime.Equal(stime) || itemTime.After(stime)) &&
			(itemTime.Equal(etime) || itemTime.Before(etime)) {
			filteredData = append(filteredData, item)
		}
		// // 判断是否在范围内
		// if (itemTime.Equal(sdateTruncated) || itemTime.After(sdateTruncated)) &&
		// 	(itemTime.Equal(edateTruncated) || itemTime.Before(edateTruncated)) {
		// 	filteredData = append(filteredData, item)
		// }
	}

	return filteredData, nil
}

// 按日期范围获取1m分时行情数据
func GetCSV1m(gmcsv string,
	symbol string, sdate string, edate string,
	istimestamp bool, clip bool,
	timeoutSeconds int) ([]map[string]any, error) {

	// url := fmt.Sprintf("%s/download/", gmcsv)
	url := gmcsv

	now := time.Now()
	today := now.Format("2006-01-02")

	sday := today
	eday := today
	if sdate != "" {
		sday = sdate
	}
	if edate != "" {
		eday = edate
	}
	if sday > eday {
		return nil, fmt.Errorf("开始日期大于结束日期")
	}

	// 解析开始和结束日期
	sdateTime, err := time.Parse("2006-01-02", sday)
	if err != nil {
		return nil, fmt.Errorf(" error parsing start date: %v", err)
	}

	edateTime, err := time.Parse("2006-01-02", eday)
	if err != nil {
		return nil, fmt.Errorf(" error parsing end date: %v", err)
	}

	syy := sdateTime.Year()
	smm := int(sdateTime.Month())
	eyy := edateTime.Year()
	emm := int(edateTime.Month())

	tag := "1m"
	var ddd []map[string]any
	for yy := syy; yy <= eyy; yy++ {
		if yy == syy || yy == eyy {
			var smonth, emonth int

			smonth = 1
			emonth = 12
			if yy == syy {
				smonth = smm
			}
			if yy == eyy {
				emonth = emm
			}
			if yy == now.Year() && emonth > int(now.Month()) {
				emonth = int(now.Month())
			}

			var ddm []map[string]any
			for im := smonth; im <= emonth; im++ {

				rsp, err := GetCSVMonth(url, symbol, im, yy, istimestamp, timeoutSeconds)
				if err != nil {
					fmt.Printf("获取月CSV数据失败: %s", err)
				}
				ddm = append(ddm, rsp...)
			}

			if ddm != nil {
				ddd = append(ddd, ddm...)
			}
		} else {
			// 调用年份数据获取函数
			rsp, err := GetCSVYear(url, symbol, tag, yy, istimestamp, "timestamp", timeoutSeconds)
			if err != nil {
				fmt.Printf("没有获取到数据: %d", yy)
			}
			if rsp != nil {
				ddd = append(ddd, rsp...)
			}
		}
	}

	if len(ddd) == 0 {
		return nil, fmt.Errorf("没有获取到数据: %d - %d", syy, eyy)
	}
	if clip {
		// 过滤数据
		ddd, err = filterDataByDate(ddd, "timestamp", sday, eday)
		if err != nil {
			return nil, fmt.Errorf("过滤数据失败: %w", err)
		}
	}
	return ddd, nil
}

// 按日期范围获取[vv,pe]日频行情数据
// tag string: vv,pe
func GetCSVTag(gmcsv string, tag string,
	symbol string, sdate string, edate string,
	istimestamp bool, clip bool,
	timeoutSeconds int) ([]map[string]any, error) {

	// url := fmt.Sprintf("%s/download/", gmcsv)
	url := gmcsv

	now := time.Now()
	today := now.Format("2006-01-02")

	sday := today
	eday := today
	if sdate != "" {
		sday = sdate
	}
	if edate != "" {
		eday = edate
	}
	if sday > eday {
		return nil, fmt.Errorf("开始日期大于结束日期")
	}

	// 解析开始和结束日期
	sdateTime, err := time.Parse("2006-01-02", sday)
	if err != nil {
		return nil, fmt.Errorf(" error parsing start date: %v", err)
	}

	edateTime, err := time.Parse("2006-01-02", eday)
	if err != nil {
		return nil, fmt.Errorf(" error parsing end date: %v", err)
	}

	syy := sdateTime.Year()
	eyy := edateTime.Year()

	lookuptab := map[string]string{
		"1m": "timestamp",
		"vv": "timestamp",
		"pe": "trade_date",
	}

	// istimestamp = false
	var ddd []map[string]any
	for yy := syy; yy <= eyy; yy++ {
		rsp, err := GetCSVYear(url, symbol, tag, yy, istimestamp, lookuptab[tag], timeoutSeconds)
		if err != nil {
			// fmt.Printf("获取年CSV数据失败: %s", err)
			fmt.Printf("没有获取到数据: %d\n", yy)
		}
		if rsp != nil {
			ddd = append(ddd, rsp...)
		}
	}

	if len(ddd) == 0 {
		return nil, fmt.Errorf("没有获取到数据(%s): %d - %d", tag, syy, eyy)
	}
	if clip {
		// 过滤数据
		ddd, err = filterDataByDate(ddd, lookuptab[tag], sday, eday)
		if err != nil {
			fmt.Printf("过滤数据失败: %s\n", err)
			return nil, fmt.Errorf("过滤数据失败: %w", err)
		}
		if len(ddd) == 0 {
			fmt.Printf("过滤数据后数据为空:(%s): %d - %d\n", tag, syy, eyy)
			return nil, fmt.Errorf("过滤数据后数据为空(%s): %d - %d", tag, syy, eyy)
		}
	}
	return ddd, nil
}

// 按日期范围获取1m分时行情数据
func GetGM1m(gmcsv string, gmapi string,
	symbol string, sdate string, edate string, istimestamp bool, include bool,
	timeoutSeconds int) ([]map[string]any, error) {

	var ddd []map[string]any

	sday := sdate
	eday := edate
	if !include {
		etime, _ := time.Parse("2006-01-02", eday)
		etime = etime.AddDate(0, 0, -1)
		eday = etime.Format("2006-01-02")
	}
	if sday > eday {
		return nil, fmt.Errorf("开始日期大于结束日期: sdate=%s, edate=%s", sday, eday)
	}

	isclip := true

	dcsv, _ := GetCSV1m(gmcsv, symbol, sday, eday, istimestamp, isclip, timeoutSeconds)
	if len(dcsv) > 0 {
		ddd = append(ddd, dcsv...)
		// fmt.Printf("获取CSV数据成功: %d条: %s - %s\n", len(dcsv), sday, eday)
		// etime, _ := time.Parse("2006-01-02", dcsv[len(dcsv)-1]["timestamp"].(string))
		tsStr := dcsv[len(dcsv)-1]["timestamp"]
		etime, _ := ParseTimestamp(tsStr)
		// var etime time.Time
		// if istimestamp {
		// 	etime = ConvertString2Time(tsStr.(string))
		// } else {
		// 	etime = time.UnixMilli(tsStr.(int64))
		// }
		etime = etime.AddDate(0, 0, 1)
		sday = etime.Format("2006-01-02") // 开始时间设置为下一天
		// fmt.Printf(" 下个开始日期: %s \n", sday)
	}

	if sday > eday {
		return ddd, nil
	}

	dapi, _ := GetKbarsHis(gmapi, symbol, "1m", sday, eday, istimestamp, timeoutSeconds)

	// jsonData, _ := json.Marshal(dapi[:5])
	// fmt.Println(string(jsonData))

	for i := range dapi {
		// 去掉API数据中的symbol字段
		dd1 := make(map[string]any, 6)

		dd1["timestamp"] = dapi[i]["timestamp"]
		dd1["open"] = dapi[i]["open"]
		dd1["high"] = dapi[i]["high"]
		dd1["low"] = dapi[i]["low"]
		dd1["close"] = dapi[i]["close"]
		dd1["volume"] = dapi[i]["volume"]
		ddd = append(ddd, dd1)
	}
	// ddd = append(ddd, dapi...)
	// if dapi != nil {
	// 	// fmt.Printf("获取API数据成功: %d条: %s - %s\n", len(dapi), sday, eday)
	// 	// 处理一下数据结构，使得CSV数据和API数据合并
	// }

	return ddd, nil
}
