package gm

import (
	"context"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/ulikunitz/xz"
)

// ConvertToDuration 转换秒数为时间段
func ConvertToDuration(seconds int) time.Duration {
	return time.Duration(seconds) * time.Second
}

func ParseURL(urlStr string) (*url.URL, error) {
	return url.Parse(urlStr)
}

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

// GetTradingDays 返回交易日历
// startDate和endDate可以为nil，表示使用默认值
// region: "us"(美国), "cn"(中国), ""(无节假日)
// func GetTradingDays(startDate, endDate *time.Time, region string) []string {
// 	now := time.Now()

// 	// 设置默认开始日期为当年1月1日
// 	var start time.Time
// 	if startDate == nil {
// 		start = time.Date(now.Year(), 1, 1, 0, 0, 0, 0, time.Local)
// 	} else {
// 		start = *startDate
// 	}

// 	// 设置默认结束日期为当年12月31日
// 	var end time.Time
// 	if endDate == nil {
// 		end = time.Date(now.Year(), 12, 31, 0, 0, 0, 0, time.Local)
// 	} else {
// 		end = *endDate
// 	}

// 	// 创建日历实例
// 	c := cal.NewCalendar()

// 	// 根据地区添加节假日
// 	switch region {
// 	case "us":
// 		cal.UnitedStates.AddTo(c)
// 	case "cn":
// 		cal.China.AddTo(c)
// 	}

// 	var tradingDays []string
// 	current := start

// 	// 遍历日期范围
// 	for !current.After(end) {
// 		// 排除周末和节假日
// 		if c.IsWorkday(current) {
// 			tradingDays = append(tradingDays, current.Format("2006-01-02"))
// 		}
// 		current = current.AddDate(0, 0, 1)
// 	}

// 	return tradingDays
// }

// 可用的地区代码 (部分)
// us - 美国节假日
// cn - 中国节假日
// ca - 加拿大节假日
// gb - 英国节假日
// de - 德国节假日
// jp - 日本节假日
// au - 澳大利亚节假日
// 更多地区请参考: https://github.com/rickar/cal

// 辅助函数：从字符串解析日期
func ParseDate(dateStr string) (*time.Time, error) {
	t, err := time.Parse("2006-01-02", dateStr)
	if err != nil {
		return nil, err
	}
	return &t, nil
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
func GetPrevN(gmapi string, date string, count int, timeoutSeconds int) ([]byte, error) {
	url := fmt.Sprintf("%s/get_dates_prev_n", gmapi)
	params := map[string]string{
		"date":  date,
		"count": fmt.Sprintf("%d", count),
	}

	resp, err := fetchURLData(url, time.Duration(timeoutSeconds)*time.Second, params)
	if err != nil {
		// fmt.Printf(" 获取数据失败(gm.GetCalendar): %v\n", err)
		return nil, err
	}

	return resp, nil
}

// 查询指定日期的后n个交易日
func GetNextN(gmapi string, date string, count int, timeoutSeconds int) ([]byte, error) {
	url := fmt.Sprintf("%s/get_dates_next_n", gmapi)
	params := map[string]string{
		"date":  date,
		"count": fmt.Sprintf("%d", count),
	}

	resp, err := fetchURLData(url, time.Duration(timeoutSeconds)*time.Second, params)
	if err != nil {
		// fmt.Printf(" 获取数据失败(gm.GetCalendar): %v\n", err)
		return nil, err
	}

	return resp, nil
}

// 获取行情快照数据
func GetCurrent(gmapi string, symbols string, timeoutSeconds int, split bool) ([]byte, error) {
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

// 获取K线行情数据
func GetKbarsHis(gmapi string,
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
func GetKbarsHisN(gmapi string,
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

// csvToJSON 将CSV数据转换为JSON格式
func csvToJSON(csvData []byte) ([]byte, error) {
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
	var result []map[string]string
	for i := 1; i < len(records); i++ {
		row := make(map[string]string)
		for j, value := range records[i] {
			if j < len(headers) {
				row[headers[j]] = value
			}
		}
		result = append(result, row)
	}

	// 转换为JSON
	jsonData, err := json.MarshalIndent(result, "", "  ")
	if err != nil {
		return nil, fmt.Errorf("转换JSON失败: %w", err)
	}

	return jsonData, nil
}

// downloadAndConvertToJSON 从URL下载CSV数据并转换为JSON
func DownloadAndConvertToJSON(url string) ([]byte, error) {
	// 下载并读取数据
	csvData, err := downloadAndReadData(url)
	if err != nil {
		return nil, err
	}

	// 转换为JSON
	jsonData, err := csvToJSON(csvData)
	if err != nil {
		return nil, err
	}

	return jsonData, nil
}

// 获取K线行情数据
func GetCSVMonth(gmcsv string,
	symbol string, tag string,
	month int, year int,
	timeoutSeconds int) ([]byte, error) {

	url := fmt.Sprintf("%s/download/", gmcsv)

	// now := time.Now()
	// year := now.Year()
	// month := int(now.Month())
	fpath := getFilePathMonth(symbol, year, month)
	fmt.Println(fpath)

	resp, err := DownloadAndConvertToJSON(url + fpath)
	if err != nil {
		// fmt.Printf("获取数据失败: %s\n", err)
		return nil, err
	}
	return resp, nil
}
