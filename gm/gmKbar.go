package gm

import (
	"encoding/json"
	"fmt"
	"sort"
	"time"
)

type KBar struct {
	// Timestamp string  `json:"timestamp"`
	Open   float64 `json:"open"`
	High   float64 `json:"high"`
	Low    float64 `json:"low"`
	Close  float64 `json:"close"`
	Volume int64   `json:"volume"`
}

// KLineData 结构体( klinechart 图表数据)
type KLineData struct {
	Timestamp int64   `json:"timestamp"` // 时间戳
	Open      float64 `json:"open"`
	High      float64 `json:"high"`
	Low       float64 `json:"low"`
	Close     float64 `json:"close"`
	Volume    int64   `json:"volume"`
	Turnover  float64 `json:"turnover"` // 成交额, 非必须字段，如果需要展示技术指标'EMV'和'AVP'，则需要为该字段填充数据。
}

type KBarArray []KBar              // KBar 数组
type KBarMapInt64 map[int64]KBar   // KBar 字典
type KBarMapString map[string]KBar // KBar 字典

type KbarMapAny struct {
	Datalist map[any]KBar `json:"timestamp"`
}

func (kb *KBar) ToList(ts any) []any {
	var records []any
	records = append(records, ts)
	records = append(records, kb.Open)
	records = append(records, kb.High)
	records = append(records, kb.Low)
	records = append(records, kb.Close)
	records = append(records, kb.Volume)
	return records
}

// records 格式是一个 map 数组，每个 map 的键是列名，值是对应的数据
func (kb *KBar) ToRecords(ts any) map[any]any {
	// var records map[any]any
	records := make(map[any]any, 6)
	records["timestamp"] = ts
	records["open"] = kb.Open
	records["high"] = kb.High
	records["low"] = kb.Low
	records["close"] = kb.Close
	records["volume"] = kb.Volume
	return records
}

// records 格式是一个 map 数组，每个 map 的键是列名，值是对应的数据
func (kb *KBar) ToKLineDataFromString(ts string) KLineData {
	chinaLocation, _ := time.LoadLocation("Asia/Shanghai")
	tt, err := time.ParseInLocation("2006-01-02 15:04:05", ts, chinaLocation)
	if err != nil {
		fmt.Printf("解析时间失败: %v", err)
	}
	kl := KLineData{
		Timestamp: tt.UnixMilli(),
		Open:      kb.Open,
		High:      kb.High,
		Low:       kb.Low,
		Close:     kb.Close,
		Volume:    kb.Volume,
	}
	return kl
}

// records 格式是一个 map 数组，每个 map 的键是列名，值是对应的数据
func (kb *KBar) ToKLineDataFromInt(ti int64) KLineData {
	kl := KLineData{
		Timestamp: ti,
		Open:      kb.Open,
		High:      kb.High,
		Low:       kb.Low,
		Close:     kb.Close,
		Volume:    kb.Volume,
	}
	return kl
}

// TickRecord 分时行情数据结构
type TickRecord struct {
	Timestamp time.Time // 时间戳
	Price     float64   // 价格
	Volume    int64     // 成交量
}

// DailyKLine 日K线数据结构
type DailyKLine struct {
	Date   time.Time // 日期
	Open   float64   // 开盘价
	High   float64   // 最高价
	Low    float64   // 最低价
	Close  float64   // 收盘价
	Volume int64     // 成交量
}

// ConvertTicksToDaily 将分时行情数据转换为日K线数据
func ConvertTicksToDaily(records []TickRecord) []DailyKLine {
	if len(records) == 0 {
		return []DailyKLine{}
	}

	// 按时间排序
	sort.Slice(records, func(i, j int) bool {
		return records[i].Timestamp.Before(records[j].Timestamp)
	})

	// 按日期分组
	dailyGroups := make(map[string][]TickRecord)
	for _, record := range records {
		dateKey := record.Timestamp.Format("2006-01-02")
		dailyGroups[dateKey] = append(dailyGroups[dateKey], record)
	}

	// 转换为日K线数据
	var dailyKLines []DailyKLine
	for dateStr, dayRecords := range dailyGroups {
		if len(dayRecords) == 0 {
			continue
		}

		// 解析日期
		date, err := time.Parse("2006-01-02", dateStr)
		if err != nil {
			continue
		}

		// 按时间排序当日数据
		sort.Slice(dayRecords, func(i, j int) bool {
			return dayRecords[i].Timestamp.Before(dayRecords[j].Timestamp)
		})

		// 计算OHLCV
		open := dayRecords[0].Price
		close := dayRecords[len(dayRecords)-1].Price
		high := dayRecords[0].Price
		low := dayRecords[0].Price
		var totalVolume int64

		for _, record := range dayRecords {
			if record.Price > high {
				high = record.Price
			}
			if record.Price < low {
				low = record.Price
			}
			totalVolume += record.Volume
		}

		dailyKLine := DailyKLine{
			Date:   date,
			Open:   open,
			High:   high,
			Low:    low,
			Close:  close,
			Volume: totalVolume,
		}

		dailyKLines = append(dailyKLines, dailyKLine)
	}

	// 按日期排序
	sort.Slice(dailyKLines, func(i, j int) bool {
		return dailyKLines[i].Date.Before(dailyKLines[j].Date)
	})

	return dailyKLines
}

// ConvertTicksToDailyWithValidation 带数据验证的转换函数
func ConvertTicksToDailyWithValidation(records []TickRecord) ([]DailyKLine, error) {
	if len(records) == 0 {
		return []DailyKLine{}, fmt.Errorf("输入数据为空")
	}

	// 数据验证
	validRecords := make([]TickRecord, 0, len(records))
	for _, record := range records {
		if record.Price <= 0 {
			continue // 跳过无效价格
		}
		if record.Volume < 0 {
			continue // 跳过负成交量
		}
		validRecords = append(validRecords, record)
	}

	if len(validRecords) == 0 {
		return []DailyKLine{}, fmt.Errorf("没有有效的分时数据")
	}

	return ConvertTicksToDaily(validRecords), nil
}

// ===================================================
type KBarData struct {
	Timestamp time.Time `json:"timestamp"` // 时间戳
	Open      float64   `json:"open"`      // 开盘价
	High      float64   `json:"high"`      // 最高价
	Low       float64   `json:"low"`       // 最低价
	Close     float64   `json:"close"`     // 收盘价
	Volume    int64     `json:"volume"`    // 成交量
}

// 自定义时间解析（处理多种时间戳格式）
func (k *KBarData) UnmarshalJSON(data []byte) error {
	// 定义一个临时结构体来处理原始JSON数据
	type Alias KBarData
	aux := &struct {
		Timestamp any `json:"timestamp"`
		*Alias
	}{
		Alias: (*Alias)(k),
	}

	if err := json.Unmarshal(data, &aux); err != nil {
		return err
	}

	// 处理不同类型的时间戳
	timestamp, err := ParseTimestamp(aux.Timestamp)
	if err != nil {
		return fmt.Errorf("解析时间戳失败: %v", err)
	}
	k.Timestamp = timestamp

	return nil
}

// ParseKBarFromJSON 从JSON字符串解析单个K线数据
func ParseKBarFromJSON(jsonStr string) (*KBarData, error) {
	var kbar KBarData
	err := json.Unmarshal([]byte(jsonStr), &kbar)
	if err != nil {
		return nil, fmt.Errorf("JSON解析失败: %v", err)
	}
	return &kbar, nil
}

// ParseKBarArrayFromJSON 从JSON字符串解析K线数组
func ParseKBarArrayFromJSON(jsonStr string) ([]KBarData, error) {
	var kbars []KBarData
	err := json.Unmarshal([]byte(jsonStr), &kbars)
	if err != nil {
		return nil, fmt.Errorf("JSON数组解析失败: %v", err)
	}
	return kbars, nil
}

// ParseKBarFromFile 从文件读取JSON并解析
func ParseKBarFromFile(filename string) ([]KBarData, error) {
	// 这里需要导入io/ioutil或os包
	// data, err := io.ReadFile(filename)
	// if err != nil {
	// 	return nil, err
	// }
	// return ParseKBarArrayFromJSON(string(data))

	// 示例实现（实际使用时取消注释上面的代码）
	return nil, fmt.Errorf("请实现文件读取功能")
}

// 示例：处理嵌套JSON结构
type APIResponse struct {
	Code    int        `json:"code"`
	Message string     `json:"message"`
	Data    []KBarData `json:"data"`
}

// ParseKBarFromAPIResponse 从API响应中解析K线数据
func ParseKBarFromAPIResponse(jsonStr string) ([]KBarData, error) {
	var response APIResponse
	err := json.Unmarshal([]byte(jsonStr), &response)
	if err != nil {
		return nil, fmt.Errorf("API响应解析失败: %v", err)
	}

	if response.Code != 0 {
		return nil, fmt.Errorf("API返回错误: %s", response.Message)
	}

	return response.Data, nil
}

// OHLCV 不包含时间戳的价格数据
type OHLCV struct {
	Open   float64 `json:"open"`
	High   float64 `json:"high"`
	Low    float64 `json:"low"`
	Close  float64 `json:"close"`
	Volume int64   `json:"volume"`
}

// 方法1: 转换为以RFC3339格式时间戳为键的JSON
func ConvertToTimestampKeyedJSON(kbars []KBarData) (string, error) {
	result := make(map[string]OHLCV)

	for _, kbar := range kbars {
		timeKey := kbar.Timestamp.Format(time.RFC3339)
		result[timeKey] = OHLCV{
			Open:   kbar.Open,
			High:   kbar.High,
			Low:    kbar.Low,
			Close:  kbar.Close,
			Volume: kbar.Volume,
		}
	}

	jsonData, err := json.MarshalIndent(result, "", "  ")
	if err != nil {
		return "", fmt.Errorf("JSON序列化失败: %v", err)
	}

	return string(jsonData), nil
}

// 方法2: 转换为以Unix时间戳为键的JSON
func ConvertToUnixTimestampKeyedJSON(kbars []KBarData) (string, error) {
	result := make(map[string]OHLCV)

	for _, kbar := range kbars {
		timeKey := fmt.Sprintf("%d", kbar.Timestamp.Unix())
		result[timeKey] = OHLCV{
			Open:   kbar.Open,
			High:   kbar.High,
			Low:    kbar.Low,
			Close:  kbar.Close,
			Volume: kbar.Volume,
		}
	}

	jsonData, err := json.MarshalIndent(result, "", "  ")
	if err != nil {
		return "", fmt.Errorf("JSON序列化失败: %v", err)
	}

	return string(jsonData), nil
}

// 方法3: 转换为以自定义格式时间戳为键的JSON
func ConvertToCustomTimestampKeyedJSON(kbars []KBarData, timeFormat string) (string, error) {
	result := make(map[string]OHLCV)

	for _, kbar := range kbars {
		timeKey := kbar.Timestamp.Format(timeFormat)
		result[timeKey] = OHLCV{
			Open:   kbar.Open,
			High:   kbar.High,
			Low:    kbar.Low,
			Close:  kbar.Close,
			Volume: kbar.Volume,
		}
	}

	jsonData, err := json.MarshalIndent(result, "", "  ")
	if err != nil {
		return "", fmt.Errorf("JSON序列化失败: %v", err)
	}

	return string(jsonData), nil
}

// 方法4: 转换为以毫秒时间戳为键的JSON
func ConvertToMillisTimestampKeyedJSON(kbars []KBarData) (string, error) {
	result := make(map[string]OHLCV)

	for _, kbar := range kbars {
		timeKey := fmt.Sprintf("%d", kbar.Timestamp.UnixMilli())
		result[timeKey] = OHLCV{
			Open:   kbar.Open,
			High:   kbar.High,
			Low:    kbar.Low,
			Close:  kbar.Close,
			Volume: kbar.Volume,
		}
	}

	jsonData, err := json.MarshalIndent(result, "", "  ")
	if err != nil {
		return "", fmt.Errorf("JSON序列化失败: %v", err)
	}

	return string(jsonData), nil
}

// 方法5: 转换为紧凑格式（数组形式的OHLCV）
func ConvertToCompactTimestampKeyedJSON(kbars []KBarData) (string, error) {
	result := make(map[string][]interface{})

	for _, kbar := range kbars {
		timeKey := kbar.Timestamp.Format(time.RFC3339)
		result[timeKey] = []interface{}{
			kbar.Open,
			kbar.High,
			kbar.Low,
			kbar.Close,
			kbar.Volume,
		}
	}

	jsonData, err := json.MarshalIndent(result, "", "  ")
	if err != nil {
		return "", fmt.Errorf("JSON序列化失败: %v", err)
	}

	return string(jsonData), nil
}

// 方法6: 转换为嵌套结构（按日期分组）
func ConvertToNestedTimestampKeyedJSON(kbars []KBarData) (string, error) {
	result := make(map[string]map[string]OHLCV)

	for _, kbar := range kbars {
		dateKey := kbar.Timestamp.Format("2006-01-02")
		timeKey := kbar.Timestamp.Format("15:04:05")

		if result[dateKey] == nil {
			result[dateKey] = make(map[string]OHLCV)
		}

		result[dateKey][timeKey] = OHLCV{
			Open:   kbar.Open,
			High:   kbar.High,
			Low:    kbar.Low,
			Close:  kbar.Close,
			Volume: kbar.Volume,
		}
	}

	jsonData, err := json.MarshalIndent(result, "", "  ")
	if err != nil {
		return "", fmt.Errorf("JSON序列化失败: %v", err)
	}

	return string(jsonData), nil
}

// 通用转换函数，支持多种选项
type ConvertOptions struct {
	TimeFormat    string // 时间格式
	UseUnixTime   bool   // 使用Unix时间戳
	UseMillis     bool   // 使用毫秒时间戳
	CompactFormat bool   // 使用紧凑格式（数组）
	GroupByDate   bool   // 按日期分组
	PrettyPrint   bool   // 美化输出
}

func ConvertKBarToTimestampJSON(kbars []KBarData, options ConvertOptions) (string, error) {
	if len(kbars) == 0 {
		return "{}", nil
	}

	// 按日期分组
	if options.GroupByDate {
		return ConvertToNestedTimestampKeyedJSON(kbars)
	}

	// 紧凑格式
	if options.CompactFormat {
		return ConvertToCompactTimestampKeyedJSON(kbars)
	}

	// 普通格式
	result := make(map[string]OHLCV)

	for _, kbar := range kbars {
		var timeKey string

		if options.UseMillis {
			timeKey = fmt.Sprintf("%d", kbar.Timestamp.UnixMilli())
		} else if options.UseUnixTime {
			timeKey = fmt.Sprintf("%d", kbar.Timestamp.Unix())
		} else {
			format := options.TimeFormat
			if format == "" {
				format = time.RFC3339
			}
			timeKey = kbar.Timestamp.Format(format)
		}

		result[timeKey] = OHLCV{
			Open:   kbar.Open,
			High:   kbar.High,
			Low:    kbar.Low,
			Close:  kbar.Close,
			Volume: kbar.Volume,
		}
	}

	var jsonData []byte
	var err error

	if options.PrettyPrint {
		jsonData, err = json.MarshalIndent(result, "", "  ")
	} else {
		jsonData, err = json.Marshal(result)
	}

	if err != nil {
		return "", fmt.Errorf("JSON序列化失败: %v", err)
	}

	return string(jsonData), nil
}

// 辅助函数：从JSON字符串中根据时间戳键获取数据
func GetKBarByTimestamp(jsonStr, timestampKey string) (*OHLCV, error) {
	var data map[string]OHLCV
	err := json.Unmarshal([]byte(jsonStr), &data)
	if err != nil {
		return nil, fmt.Errorf("JSON解析失败: %v", err)
	}

	if ohlcv, exists := data[timestampKey]; exists {
		return &ohlcv, nil
	}

	return nil, fmt.Errorf("未找到时间戳 %s 对应的数据", timestampKey)
}
