package srv

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"regexp"
	"time"
)

// SmartURLHandler 智能URL处理器，可以根据端口号判断协议
// 参数:
//
//	url: 输入的URL
//	preferHTTPS: 是否优先使用HTTPS
//
// 返回:
//
//	处理后的完整URL
func SmartURLHandler(url string, preferHTTPS bool) string {
	if url == "" {
		return url
	}

	// 如果已有协议，直接返回
	protocolRegex := regexp.MustCompile(`^[a-zA-Z][a-zA-Z0-9+.-]*://`)
	if protocolRegex.MatchString(url) {
		return url
	}

	// 检查常见的HTTPS端口
	httpsports := []string{"443", "8443"}
	httpPorts := []string{"80", "8080", "3000", "5000", "8000"}

	// 提取端口号
	portRegex := regexp.MustCompile(`:(\d+)`)
	portMatch := portRegex.FindStringSubmatch(url)

	if len(portMatch) > 1 {
		port := portMatch[1]

		// 检查是否为HTTPS端口
		for _, httpsPort := range httpsports {
			if port == httpsPort {
				return "https://" + url
			}
		}

		// 检查是否为HTTP端口或不优先HTTPS
		for _, httpPort := range httpPorts {
			if port == httpPort {
				if !preferHTTPS {
					return "http://" + url
				}
				break
			}
		}

		// 如果端口不在预定义列表中且优先HTTPS
		if preferHTTPS {
			return "https://" + url
		} else {
			return "http://" + url
		}
	}

	// 默认协议
	if preferHTTPS {
		return "https://" + url
	}
	return "http://" + url
}

// // 辅助函数：检查字符串是否在切片中
// func contains(slice []string, item string) bool {
// 	return slices.Contains(slice, item)
// }

// getURLWithoutRetry sends an HTTP GET request without retries.
// It returns the JSON unmarshaled response or an error if the request fails
// or the status code is not 2xx.
//
// connectTimeout: Connection timeout (default 3 seconds if 0).
// dataTimeout: Data transfer timeout (defaults to connectTimeout if 0).
func GetURLWithoutRetry(url string, params map[string]string, connectTimeout time.Duration, dataTimeout time.Duration) (map[string]any, error) {
	if connectTimeout == 0 {
		connectTimeout = 3 * time.Second
	}
	if dataTimeout == 0 {
		dataTimeout = connectTimeout
	}

	client := &http.Client{
		Timeout: connectTimeout + dataTimeout, // Total timeout for the request
	}

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	q := req.URL.Query()
	for key, val := range params {
		q.Add(key, val)
	}
	req.URL.RawQuery = q.Encode()

	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return nil, fmt.Errorf("HTTP error: status code %d", resp.StatusCode)
	}

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	var result map[string]any
	err = json.Unmarshal(bodyBytes, &result)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal JSON response: %w", err)
	}

	return result, nil
}

// getURLWithRetry sends an HTTP GET request with retry mechanism.
// It retries up to 5 times with a 3-second delay between attempts.
// It returns the JSON unmarshaled response or an error if all retries fail.
//
// connectTimeout: Connection timeout (default 3 seconds if 0).
// dataTimeout: Data transfer timeout (defaults to connectTimeout if 0).
func GetURLWithRetry(url string, params map[string]string, connectTimeout time.Duration, dataTimeout time.Duration) (map[string]any, error) {
	const (
		maxRetries = 5
		delay      = 3 * time.Second
	)

	if connectTimeout == 0 {
		connectTimeout = 3 * time.Second
	}
	if dataTimeout == 0 {
		dataTimeout = connectTimeout
	}

	client := &http.Client{
		Timeout: connectTimeout + dataTimeout, // Total timeout for the request
	}

	var lastErr error
	for range maxRetries {
		req, err := http.NewRequest("GET", url, nil)
		if err != nil {
			lastErr = fmt.Errorf("failed to create request: %w", err)
			time.Sleep(delay) // Still sleep to respect delay even if request creation fails
			continue
		}

		q := req.URL.Query()
		for key, val := range params {
			q.Add(key, val)
		}
		req.URL.RawQuery = q.Encode()

		resp, err := client.Do(req)
		if err != nil {
			lastErr = fmt.Errorf("failed to send request: %w", err)
			time.Sleep(delay)
			continue
		}
		defer resp.Body.Close() // Ensure body is closed after each attempt

		if resp.StatusCode < 200 || resp.StatusCode >= 300 {
			lastErr = fmt.Errorf("HTTP error: status code %d", resp.StatusCode)
			time.Sleep(delay)
			continue
		}

		bodyBytes, err := io.ReadAll(resp.Body)
		if err != nil {
			lastErr = fmt.Errorf("failed to read response body: %w", err)
			time.Sleep(delay)
			continue
		}

		var result map[string]any
		err = json.Unmarshal(bodyBytes, &result)
		if err != nil {
			lastErr = fmt.Errorf("failed to unmarshal JSON response: %w", err)
			time.Sleep(delay)
			continue
		}

		return result, nil // Success!
	}

	return nil, fmt.Errorf("all %d retries failed: %w", maxRetries, lastErr)
}

// RawColData 结构体用于匹配从 URL 获取的原始 JSON 格式
type RawColData struct {
	Columns []string `json:"columns"`
	Data    [][]any  `json:"data"` // 使用 interface{} 来处理不同类型的数据
}

// transformToRecords 将 InputData 格式转换为 records 格式
// records 格式是一个 map 数组，每个 map 的键是列名，值是对应的数据
func (rcd *RawColData) TransformToRecords() ([]map[string]any, error) {
	var records []map[string]any

	if len(rcd.Columns) == 0 && len(rcd.Data) > 0 {
		return nil, fmt.Errorf("当存在数据时，列名不能为空")
	}

	for _, row := range rcd.Data {
		// 检查数据行长度是否与列名长度匹配
		if len(row) != len(rcd.Columns) {
			return nil, fmt.Errorf("数据行长度 (%d) 与列名长度 (%d) 不匹配", len(row), len(rcd.Columns))
		}
		record := make(map[string]any)
		for i, colName := range rcd.Columns {
			record[colName] = row[i]
		}
		records = append(records, record)
	}
	return records, nil
}

// InputData 结构体用于匹配从 URL 获取的原始 JSON 格式
type StCalendar struct {
	Date          []string `json:"date"`
	NextTradeDate []string `json:"next_trade_date"`
	PrevTradeDate []string `json:"pre_trade_date"`
	TradeDate     []string `json:"trade_date"`
}

func TransformToStCalendar(input RawColData) ([]map[string]any, error) {
	var records []map[string]any

	if len(input.Columns) == 0 && len(input.Data) > 0 {
		return nil, fmt.Errorf("当存在数据时，列名不能为空")
	}

	for _, row := range input.Data {
		// 检查数据行长度是否与列名长度匹配
		if len(row) != len(input.Columns) {
			return nil, fmt.Errorf("数据行长度 (%d) 与列名长度 (%d) 不匹配", len(row), len(input.Columns))
		}
		record := make(map[string]any)
		for i, colName := range input.Columns {
			record[colName] = row[i]
		}
		records = append(records, record)
	}
	return records, nil
}
