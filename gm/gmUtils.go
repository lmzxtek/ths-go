package gm

import (
	"fmt"
	"net/url"
	"strings"
	"time"
)

// 返回字符串的最后n个字符
func LastNChars(s string, n int) string {
	runes := []rune(s)
	num := len(runes)
	if num <= n {
		return s // 长度不足 n，返回原字符串
	}
	start := num - n
	return string(runes[start:])
}

// ConvertToDuration 转换秒数为时间段
func ConvertToDuration(seconds int) time.Duration {
	return time.Duration(seconds) * time.Second
}

// 将字符串转换为时间戳
func ConvertString2Time(s string) time.Time {
	ts := strings.TrimSpace(s)

	ts = strings.TrimSuffix(ts, "+08:00")
	ts = strings.Replace(ts, "T", " ", 1)

	tz, _ := time.LoadLocation("Asia/Shanghai")
	tt, _ := time.ParseInLocation("2006-01-02 15:04:05", ts, tz)

	// fmt.Println("=" + strings.Repeat("=", 50))
	return tt
}

// TimeStringToTimestamp 将时间字符串转换为Unix时间戳（秒）
func TimeStringToTimestamp(timeStr string) (int64, error) {
	// 解析时间字符串，使用本地时区
	t, err := time.Parse("2006-01-02 15:04:05", timeStr)
	if err != nil {
		return 0, fmt.Errorf("解析时间失败: %v", err)
	}

	// 返回Unix时间戳（秒）
	return t.Unix(), nil
}

// TimeStringToTimestampWithLocation 将时间字符串转换为指定时区的Unix时间戳（秒）
func TimeStringToTimestampWithLocation(timeStr string, locationName string) (int64, error) {
	// 加载指定时区
	location, err := time.LoadLocation(locationName)
	if err != nil {
		return 0, fmt.Errorf("加载时区失败: %v", err)
	}

	// 在指定时区解析时间字符串
	t, err := time.ParseInLocation("2006-01-02 15:04:05", timeStr, location)
	if err != nil {
		return 0, fmt.Errorf("解析时间失败: %v", err)
	}

	// 返回Unix时间戳（秒）
	return t.Unix(), nil
}

// TimeStringToTimestampMillis 将时间字符串转换为Unix时间戳（毫秒）
func TimeStringToTimestampMillis(timeStr string) (int64, error) {
	t, err := time.Parse("2006-01-02 15:04:05", timeStr)
	if err != nil {
		return 0, fmt.Errorf("解析时间失败: %v", err)
	}

	// 返回Unix时间戳（毫秒）
	return t.UnixMilli(), nil
}

// TimeStringToTimestampNano 将时间字符串转换为Unix时间戳（纳秒）
func TimeStringToTimestampNano(timeStr string) (int64, error) {
	t, err := time.Parse("2006-01-02 15:04:05", timeStr)
	if err != nil {
		return 0, fmt.Errorf("解析时间失败: %v", err)
	}

	// 返回Unix时间戳（纳秒）
	return t.UnixNano(), nil
}

// TimestampToTimeString 将Unix时间戳（秒）转换回时间字符串
func TimestampToTimeString(timestamp int64) string {
	t := time.Unix(timestamp, 0)
	return t.Format("2006-01-02 15:04:05")
}

// TimestampToTimeStringWithLocation 将Unix时间戳转换为指定时区的时间字符串
func TimestampToTimeStringWithLocation(timestamp int64, locationName string) (string, error) {
	location, err := time.LoadLocation(locationName)
	if err != nil {
		return "", fmt.Errorf("加载时区失败: %v", err)
	}

	t := time.Unix(timestamp, 0).In(location)
	return t.Format("2006-01-02 15:04:05"), nil
}

// MillisToTime 将毫秒时间戳转换为time.Time类型
func MillisToTime(millis int64) time.Time {
	// 将毫秒时间戳转换为秒和纳秒
	seconds := millis / 1000
	nanoseconds := (millis % 1000) * 1000000

	// 创建time.Time对象
	return time.Unix(seconds, nanoseconds)
}

// MillisToTimeInLocation 将毫秒时间戳转换为指定时区的time.Time类型
func MillisToTimeInLocation(millis int64, locationName string) (time.Time, error) {
	// 加载指定时区
	location, err := time.LoadLocation(locationName)
	if err != nil {
		return time.Time{}, fmt.Errorf("加载时区失败: %v", err)
	}

	// 先转换为UTC时间，再转换到指定时区
	utcTime := MillisToTime(millis)
	return utcTime.In(location), nil
}

// TimeToMillis 将time.Time类型转换为毫秒时间戳（反向转换）
func TimeToMillis(t time.Time) int64 {
	return t.UnixMilli()
}

// MillisToFormattedString 将毫秒时间戳直接转换为格式化字符串
func MillisToFormattedString(millis int64, layout string) string {
	t := MillisToTime(millis)
	return t.Format(layout)
}

// MillisToFormattedStringInLocation 将毫秒时间戳转换为指定时区的格式化字符串
func MillisToFormattedStringInLocation(millis int64, layout string, locationName string) (string, error) {
	t, err := MillisToTimeInLocation(millis, locationName)
	if err != nil {
		return "", err
	}
	return t.Format(layout), nil
}

// GetTimestampInfo 获取时间戳的详细信息
func GetTimestampInfo(millis int64) map[string]interface{} {
	t := MillisToTime(millis)

	info := map[string]interface{}{
		"timestamp_millis":  millis,
		"timestamp_seconds": millis / 1000,
		"time_utc":          t.UTC().Format("2006-01-02 15:04:05"),
		"time_local":        t.Local().Format("2006-01-02 15:04:05"),
		"year":              t.Year(),
		"month":             int(t.Month()),
		"day":               t.Day(),
		"hour":              t.Hour(),
		"minute":            t.Minute(),
		"second":            t.Second(),
		"millisecond":       t.Nanosecond() / 1000000,
		"weekday":           t.Weekday().String(),
		"unix_timestamp":    t.Unix(),
		"is_leap_year":      isLeapYear(t.Year()),
	}

	return info
}

// isLeapYear 判断是否为闰年
func isLeapYear(year int) bool {
	return year%4 == 0 && (year%100 != 0 || year%400 == 0)
}

// IsChineseStockMarketOpen 判断当前时间是否为中国股市开市时间
func IsChineseStockMarketOpen() bool {
	// 获取中国时区
	chinaLocation, err := time.LoadLocation("Asia/Shanghai")
	if err != nil {
		fmt.Printf("加载中国时区失败: %v\n", err)
		return false
	}

	// 获取当前中国时间
	now := time.Now().In(chinaLocation)

	return IsChineseStockMarketOpenAt(now)
}

// IsChineseStockMarketOpenAt 判断指定时间是否为中国股市开市时间
func IsChineseStockMarketOpenAt(t time.Time) bool {
	// 检查是否为工作日（周一到周五）
	weekday := t.Weekday()
	if weekday == time.Saturday || weekday == time.Sunday {
		return false
	}

	// 获取当前时间的小时和分钟
	hour := t.Hour()
	minute := t.Minute()
	currentTime := hour*100 + minute // 转换为HHMM格式便于比较

	// 上午交易时间：9:30-11:30
	morningStart := 930 // 9:30
	morningEnd := 1130  // 11:30

	// 下午交易时间：13:00-15:00
	afternoonStart := 1300 // 13:00
	afternoonEnd := 1500   // 15:00

	// 判断是否在交易时间内
	isInMorningSession := currentTime >= morningStart && currentTime <= morningEnd
	isInAfternoonSession := currentTime >= afternoonStart && currentTime <= afternoonEnd

	return isInMorningSession || isInAfternoonSession
}

// GetNextTradingTime 获取下一个交易时间
func GetNextTradingTime() time.Time {
	chinaLocation, _ := time.LoadLocation("Asia/Shanghai")
	now := time.Now().In(chinaLocation)

	// 如果当前是工作日
	if now.Weekday() >= time.Monday && now.Weekday() <= time.Friday {
		hour := now.Hour()
		minute := now.Minute()
		currentTime := hour*100 + minute

		// 如果在上午开市前
		if currentTime < 930 {
			return time.Date(now.Year(), now.Month(), now.Day(), 9, 30, 0, 0, chinaLocation)
		}
		// 如果在午休时间
		if currentTime > 1130 && currentTime < 1300 {
			return time.Date(now.Year(), now.Month(), now.Day(), 13, 0, 0, 0, chinaLocation)
		}
		// 如果在收市后，返回下一个交易日上午开市时间
		if currentTime > 1500 {
			nextDay := now.AddDate(0, 0, 1)
			for nextDay.Weekday() == time.Saturday || nextDay.Weekday() == time.Sunday {
				nextDay = nextDay.AddDate(0, 0, 1)
			}
			return time.Date(nextDay.Year(), nextDay.Month(), nextDay.Day(), 9, 30, 0, 0, chinaLocation)
		}
	}

	// 如果是周末，找到下周一
	nextDay := now
	for nextDay.Weekday() != time.Monday {
		nextDay = nextDay.AddDate(0, 0, 1)
	}
	return time.Date(nextDay.Year(), nextDay.Month(), nextDay.Day(), 9, 30, 0, 0, chinaLocation)
}

func ParseURL(urlStr string) (*url.URL, error) {
	return url.Parse(urlStr)
}
