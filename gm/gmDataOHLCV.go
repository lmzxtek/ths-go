package gm

import (
	"encoding/json"
	"fmt"
	"math"
	"sort"
	"strings"
	"time"
)

type OHLCVData struct {
	// Timestamp int64   `json:"timestamp"`
	Timestamp time.Time `json:"timestamp"`
	Open      float64   `json:"open"`
	High      float64   `json:"high"`
	Low       float64   `json:"low"`
	Close     float64   `json:"close"`
	Volume    int64     `json:"volume"`
}
type OHLCVList []OHLCVData

func (k *OHLCVData) ReadJson(data []byte) {
	json.Unmarshal(data, k)
}

func (k *OHLCVData) ReadMap(data map[string]any) {
	ts, _ := ParseTimestamp(data["timestamp"])
	// chinaLocation, _ := time.LoadLocation("Asia/Shanghai")
	// k.Timestamp = ts.UnixMilli() - 8*3600*1000 // 北京时间转UTC时间
	tz := time.FixedZone("CST", 8*3600) // 北京时区
	k.Timestamp = ts.In(tz)
	// k.Timestamp = ts.Add(-8 * time.Hour).In(tz)
	k.Open = data["open"].(float64)
	k.High = data["high"].(float64)
	k.Low = data["low"].(float64)
	k.Close = data["close"].(float64)
	k.Volume = int64(data["volume"].(float64))
	// k.Volume = data["volume"].(int64)
}

// 计算黄金价：(4*收盘价 + 2*开盘价 + 最高价 + 最低价) / 8
func (k *OHLCVData) GetHjj(rC float64, rO float64, rH float64, rL float64) float64 {
	return (rH*k.High + rL*k.Low + rC*k.Close + rO*k.Open) / (rC + rH + rL + rO)
}

// 显示
func (k *OHLCVData) Print(head bool) {
	if head {
		fmt.Println("Timestamp\t\tOpen\tHigh\tLow\tClose\tVolume\t")
		fmt.Println("---------\t\t----\t----\t---\t-----\t------\t")
	}
	fmt.Printf("%s\t%.2f\t%.2f\t%.2f\t%.2f\t%d\t\n",
		k.Timestamp.Local().Format("2006-01-02 15:04:05"),
		k.Open, k.High, k.Low, k.Close, k.Volume)
}

// 添加KBar到KBMinuteList
func (k *OHLCVList) Add(kbar OHLCVData) {
	*k = append(*k, kbar)
}

// 从map列表中读取KBar列表
func (k *OHLCVList) ReadMapList(kbList []map[string]any) {
	for _, data := range kbList {
		kbar := OHLCVData{}
		kbar.ReadMap(data)
		k.Add(kbar)
	}
}

func (k *OHLCVList) Head(n int) {
	num := min(len(*k), n)
	for i, kb := range (*k)[:num] {
		isHead := false
		if i == 0 {
			isHead = true
			fmt.Printf("OHLCVList Head(%d) of %d OHLCVData:\n", num, len(*k))
		}
		kb.Print(isHead)

		if i == num-1 {
			fmt.Println("=" + strings.Repeat("=", 64))
		}
	}
}

func (k *OHLCVList) Tail(n int) {
	num := min(len(*k), n)
	isHead := false
	for i, kb := range (*k)[(len(*k) - num):] {
		isHead = false
		if i == 0 {
			isHead = true
			fmt.Printf("OHLCVList Tail(%d) of %d OHLCVData:\n", num, len(*k))
		}

		kb.Print(isHead)

		if i == num-1 {
			fmt.Println("=" + strings.Repeat("=", 64))
		}
	}
}

// 对KBMinuteList进行排序
func (k *OHLCVList) Sort(descend bool) {
	// 按时间排序
	if descend {
		sort.Slice(*k, func(i, j int) bool {
			return (*k)[i].Timestamp.After((*k)[j].Timestamp)
		})
	} else {
		sort.Slice(*k, func(i, j int) bool {
			return (*k)[i].Timestamp.Before((*k)[j].Timestamp)
		})
	}
}

// 提取为单个KBarDataMinute
func (k *OHLCVList) ToOHLCVData(isDaily bool) OHLCVData {
	// var kbar KBarDataMinute
	nlen := len(*k)
	if nlen == 0 {
		return OHLCVData{}
	}

	tz := time.FixedZone("CST", 8*3600) // 北京时区
	ts := (*k)[0].Timestamp             //.Truncate(24 * time.Hour)
	t1st := time.Date(ts.Year(), ts.Month(), ts.Day(), 0, 0, 0, 0, tz)
	if !isDaily {
		// 非日频数据采用最后一个KBar的时刻作为时间戳
		t1st = (*k)[nlen-1].Timestamp.In(tz)
	}

	oo := (*k)[0].Open
	hh := (*k)[0].High
	ll := (*k)[0].Low
	cc := (*k)[nlen-1].Close
	vv := int64(0)
	for _, kb := range *k {
		vv += kb.Volume
		hh = max(hh, kb.High)
		ll = min(ll, kb.Low)
	}
	return OHLCVData{
		Timestamp: t1st,
		Open:      oo,
		High:      hh,
		Low:       ll,
		Close:     cc,
		Volume:    vv,
	}
}

// 聚合为日KBarDataMinute
func (k *OHLCVList) ToDailyList() OHLCVList {
	var kbList OHLCVList

	nlen := len(*k)
	if nlen == 0 {
		return kbList
	}

	tz := time.FixedZone("CST", 8*3600) // 北京时区

	kDic := make(map[string]OHLCVList)
	for _, kb := range *k {
		ts := kb.Timestamp //.Truncate(24 * time.Hour)
		tDay := time.Date(ts.Year(), ts.Month(), ts.Day(), 0, 0, 0, 0, tz)
		tStr := tDay.Format("2006-01-02")
		if _, ok := kDic[tStr]; !ok {
			kDic[tStr] = OHLCVList{}
		}
		kDic[tStr] = append(kDic[tStr], kb)
	}

	for _, v := range kDic {
		kb := v.ToOHLCVData(true)
		kbList.Add(kb)
	}
	kbList.Sort(false)
	return kbList
}

// 聚合分时数据为5分钟KBarDataMinute
//
//	nm: 分钟数: 5, 15, 30, 60
func (k *OHLCVList) To5m(nm int) OHLCVList {
	var kbList OHLCVList

	nlen := len(*k)
	if nlen == 0 {
		return kbList
	}

	tz := time.FixedZone("CST", 8*3600) // 北京时区

	kDic := make(map[string]OHLCVList)
	for _, kb := range *k {
		ts := kb.Timestamp //.Truncate(24 * time.Hour)
		nms := ts.Minute() + 60*ts.Hour()
		isAm := ts.Format("15:04:05") <= "11:30:00"
		if isAm {
			// 如果是上午，minute加30分钟
			nms = nms + 30
		}
		nmin := int(math.Ceil(float64(nms)/float64(nm)) * float64(nm))
		if isAm {
			// 上午的话，minute减30分钟（恢复时间）
			nmin = nmin - 30
		}
		nhr := int(nmin / 60)
		nmin = nmin - nhr*60
		tDay := time.Date(ts.Year(), ts.Month(), ts.Day(), nhr, nmin, 0, 0, tz)
		tStr := tDay.Format("2006-01-02 15:03:04")
		if _, ok := kDic[tStr]; !ok {
			kDic[tStr] = OHLCVList{}
		}
		kDic[tStr] = append(kDic[tStr], kb)
	}

	for _, v := range kDic {
		kb := v.ToOHLCVData(false)
		kbList.Add(kb)
	}
	kbList.Sort(false)
	return kbList
}

// 统计相邻OHLCVData的涨跌数量
// 可以此数量用于判断交易是否活跃
func (k *OHLCVList) GetUpDownNum() (nup int, ndown int) {
	nup = 0
	ndown = 0
	for i, kb := range *k {
		if i == 0 {
			continue
		}
		cc := (*k)[i-1].Close // 前一个KBar的收盘价
		oo := kb.Open
		if cc > oo {
			nup += 1
		} else if cc < oo {
			ndown += 1
		}
	}
	return nup, ndown
}

// 提取成交量中值: pct为百分比，如50.0表示中位数, 12.5相当于30分钟的K线
func (k *OHLCVList) GetVmed(pct float64) int64 {
	nlen := len(*k)
	if nlen == 0 {
		return 0.0
	}

	vols := make([]float64, nlen)
	for i, kb := range *k {
		vv := float64(kb.Volume)
		vols[i] += vv

	}

	sort.Float64s(vols)
	vmed := CalcMedianPct(vols, 100.0-pct)
	return int64(vmed)
}

// 提取成交量加权黄金价 PVJ: rC, rO, rH, rL 为权重
//
//	计算方法：(4*收盘价 + 2*开盘价 + 最高价 + 最低价) / 8
func (k *OHLCVList) GetPvj(rC, rO, rH, rL float64) float64 {
	sumPvj := 0.0
	sumVol := 0.0
	for _, kb := range *k {
		vv := float64(kb.Volume)
		sumVol += vv
		hjj := kb.GetHjj(rC, rO, rH, rL)
		sumPvj += vv * hjj
	}
	if sumVol == 0.0 {
		return 0.0
	}
	return sumPvj / sumVol
}

// 计算黄金价 Hjj: rC, rO, rH, rL 为权重
//
//	计算方法：(4*收盘价 + 2*开盘价 + 最高价 + 最低价) / 8
func (k *OHLCVList) GetHjjList(rC, rO, rH, rL float64) []float64 {
	var hjjList []float64
	for _, kb := range *k {
		hjj := kb.GetHjj(rC, rO, rH, rL)
		hjjList = append(hjjList, hjj)
	}
	return hjjList
}

// 提取单日特定成交量：v931, v932, v935, v940, v150
func (k *OHLCVList) GetV123Data() V123Data {
	// var kbar KBarDataMinute
	nlen := len(*k)
	if nlen == 0 {
		return V123Data{}
	}

	ts := (*k)[0].Timestamp             //.Truncate(24 * time.Hour)
	tz := time.FixedZone("CST", 8*3600) // 北京时区
	tDay := time.Date(ts.Year(), ts.Month(), ts.Day(), 0, 0, 0, 0, tz)

	v931 := int64(0)
	v932 := int64(0)
	v935 := int64(0)
	v940 := int64(0)
	v150 := int64(0)
	for _, kb := range *k {
		tt := kb.Timestamp //.Truncate(24 * time.Hour)
		vv := kb.Volume
		tStr := tt.Format("15:04:05")
		if tStr == "09:31:00" {
			v931 = vv
		}
		if tStr == "09:32:00" {
			v932 = vv
		}
		if tStr <= "09:35:00" {
			v935 += vv
		}
		if tStr > "09:35:00" && tStr <= "09:40:00" {
			v940 += vv
		}
		if tStr == "15:00:00" {
			v150 = vv
		}
	}

	return V123Data{
		TS:   tDay,
		V931: v931,
		V932: v932,
		V935: v935,
		V940: v940,
		V150: v150,
	}
}

// 提取日频成交量：v931, v932, v935, v940, v150
func (k *OHLCVList) GetV123List() V123List {
	var kbList V123List

	nlen := len(*k)
	if nlen == 0 {
		return kbList
	}

	tz := time.FixedZone("CST", 8*3600) // 北京时区

	kDic := make(map[string]OHLCVList)
	for _, kb := range *k {
		ts := kb.Timestamp //.Truncate(24 * time.Hour)
		tDay := time.Date(ts.Year(), ts.Month(), ts.Day(), 0, 0, 0, 0, tz)
		tStr := tDay.Format("2006-01-02")
		if _, ok := kDic[tStr]; !ok {
			kDic[tStr] = OHLCVList{}
		}
		kDic[tStr] = append(kDic[tStr], kb)
		// kd := kDic[tStr]
		// kd.Add(kb)
		// kDic[tStr] = kd
	}

	for _, v := range kDic {
		kb := v.GetV123Data()
		kbList = append(kbList, kb)
	}
	kbList.Sort(false)
	return kbList
}

// 计算百分比中值的函数
func CalcMedianPct(data []float64, pct float64) float64 {
	n := len(data)
	if n == 0 {
		return 0
	}
	// sort.Float64s(data)

	ratio := pct / 100.0
	n1 := min(n, int(math.Ceil(float64(n)*ratio)))
	n2 := max(0, int(math.Floor(float64(n)*ratio)))

	return (data[n1-1] + data[n2-1]) / 2
}

// 计算单日成本价+成交量中值
func (k *OHLCVList) GetCbjData(pTime string) CbjData {

	nlen := len(*k)
	if nlen == 0 {
		return CbjData{}
	}

	ts := (*k)[0].Timestamp             //.Truncate(24 * time.Hour)
	tz := time.FixedZone("CST", 8*3600) // 北京时区
	tDay := time.Date(ts.Year(), ts.Month(), ts.Day(), 0, 0, 0, 0, tz)

	nup := int64(0)
	ndown := int64(0)
	vmed := k.GetVmed(12.5)

	// var kCbj OHLCVList
	// var kCb1 OHLCVList
	// var kCb2 OHLCVList

	sumCbj := float64(0.0)
	sumCb1 := float64(0.0)
	sumCb2 := float64(0.0)
	volCbj := float64(0.0)
	volCb1 := float64(0.0)
	volCb2 := float64(0.0)

	for i, kb := range *k {
		vv := kb.Volume
		tStr := kb.Timestamp.Format("15:04:05") //.Truncate(24 * time.Hour)
		hjj := kb.GetHjj(4, 2, 1, 1)

		if vv > vmed {
			// kCbj.Add(kb)
			sumCbj += hjj * float64(vv)
			volCbj += float64(vv)
		}

		if tStr <= pTime {
			// kCb1.Add(kb)
			sumCb1 += hjj * float64(vv)
			volCb1 += float64(vv)
		} else {
			// kCb2.Add(kb)
			sumCb2 += hjj * float64(vv)
			volCb2 += float64(vv)
		}

		if i > 0 {
			cc := (*k)[i-1].Close // 前一个KBar的收盘价
			oo := kb.Open
			if cc > oo {
				nup += 1
			} else if cc < oo {
				ndown += 1
			}
		}
	}
	// cbj := kCbj.GetPvj(4, 2, 1, 1)
	// cb1 := kCb1.GetPvj(4, 2, 1, 1)
	// cb2 := kCb2.GetPvj(4, 2, 1, 1)
	cbj := float64(0.0)
	cb1 := float64(0.0)
	cb2 := float64(0.0)
	if volCbj > 0.0 {
		cbj = sumCbj / volCbj
	}
	if volCb1 > 0.0 {
		cb1 = sumCb1 / volCb1
	}
	if volCb2 > 0.0 {
		cb2 = sumCb2 / volCb2
	}

	return CbjData{
		TS:    tDay,
		Vmed:  vmed,
		Cbj:   cbj,
		Cb1:   cb1,
		Cb2:   cb2,
		Nup:   nup,
		Ndown: ndown,
	}
}

// 计算日频成本价+成交量中值
func (k *OHLCVList) GetCbjList(pTime string) CbjList {
	var kbList CbjList

	nlen := len(*k)
	if nlen == 0 {
		return kbList
	}

	tz := time.FixedZone("CST", 8*3600) // 北京时区

	kDic := make(map[string]OHLCVList)
	for _, kb := range *k {
		ts := kb.Timestamp //.Truncate(24 * time.Hour)
		tDay := time.Date(ts.Year(), ts.Month(), ts.Day(), 0, 0, 0, 0, tz)
		tStr := tDay.Format("2006-01-02")
		if _, ok := kDic[tStr]; !ok {
			kDic[tStr] = OHLCVList{}
		}
		kDic[tStr] = append(kDic[tStr], kb)
	}

	for _, v := range kDic {
		kb := v.GetCbjData(pTime)
		kbList = append(kbList, kb)
	}
	kbList.Sort(false)
	return kbList
}

// 计算日频成本价+成交量中值
func (k *OHLCVList) ToVVList() VVList {
	var kbList VVList

	nlen := len(*k)
	if nlen == 0 {
		return kbList
	}

	tz := time.FixedZone("CST", 8*3600) // 北京时区

	kDic := make(map[string]OHLCVList)
	for _, kb := range *k {
		ts := kb.Timestamp //.Truncate(24 * time.Hour)
		tDay := time.Date(ts.Year(), ts.Month(), ts.Day(), 0, 0, 0, 0, tz)
		tStr := tDay.Format("2006-01-02")
		if _, ok := kDic[tStr]; !ok {
			kDic[tStr] = OHLCVList{}
		}
		kDic[tStr] = append(kDic[tStr], kb)
	}

	for _, v := range kDic {
		var kd VVData
		kd.Init(v)
		kbList = append(kbList, kd)
	}

	return kbList
}
