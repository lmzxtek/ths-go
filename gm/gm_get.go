package gm

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"
)

// ConvertToDuration 转换秒数为时间段
func ConvertToDuration(seconds int) time.Duration {
	return time.Duration(seconds) * time.Second
}

func ParseURL(urlStr string) (*url.URL, error) {
	return url.Parse(urlStr)
}

// makeRequest 发起一个HTTP请求并打印响应状态和响应体
func makeRequest(url string) {
	resp, err := http.Get(url)
	if err != nil {
		fmt.Println("Error fetching URL: ", err)
		return
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Error reading response body: ", err)
		return
	}

	fmt.Println("Response Status: ", resp.Status)
	fmt.Println("Response Body: ", string(body))
}

// fetchData 从指定URL获取数据
func fetchData(url string) ([]byte, error) {
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

	// 获取历史K线数据
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

	// 获取历史K线数据
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

	// 获取历史K线数据
	resp, err := fetchURLData(url, time.Duration(timeoutSeconds)*time.Second, params)
	if err != nil {
		// fmt.Printf("获取数据失败: %s\n", err)
		return nil, err
	}

	return resp, nil
}

// 获取K线行情数据
func GetCSVMonth(
	symbols string, tag string,
	count int, edate string,
	gmapi string, timeoutSeconds int) ([]byte, error) {

	url := fmt.Sprintf("%s/get_his_n", gmapi)
	params := map[string]string{
		"symbols": symbols,
		"tag":     tag,
		"count":   fmt.Sprintf("%d", count),
		"edate":   edate,
	}

	// 获取历史K线数据
	resp, err := fetchURLData(url, time.Duration(timeoutSeconds)*time.Second, params)
	if err != nil {
		// fmt.Printf("获取数据失败: %s\n", err)
		return nil, err
	}

	return resp, nil
}
