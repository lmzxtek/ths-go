package gm

import (
	"slices"
	"strings"
	"time"
)

type VVData struct {
	// TS time.Time `json:"timestamp"`
	TS int64 `json:"timestamp"`
	// TS string `json:"timestamp"`

	Open   float64 `json:"open"`
	High   float64 `json:"high"`
	Low    float64 `json:"low"`
	Close  float64 `json:"close"`
	Volume int64   `json:"volume"`

	V931 int64 `json:"v931"` // volume at 9:31
	V932 int64 `json:"v932"` // volume at 9:32
	V935 int64 `json:"v935"` // volume at 9:35
	V940 int64 `json:"v940"` // volume at 9:40
	V150 int64 `json:"v150"` // volume at 15:00

	Hjj float64 `json:"hjj"` // Huangjing Price: 4*Close + 2*Open + High + Low / 8
	Pvj float64 `json:"pvj"` // volume adjusted average price in the last 30 minutes

	Vmed int64   `json:"vmed"` // median volume
	Cbj  float64 `json:"cbj"`  // Price for volume avarged Hjj over Vmed
	Cb1  float64 `json:"cb1"`  // Average Price before 10:00
	Cb2  float64 `json:"cb2"`  // Average Price after 10:00

	Nup   int64 `json:"nup"`   // number of up ticks
	Ndown int64 `json:"ndown"` // number of down ticks

	// PEttm float64 `json:"pettm"` // price to earnings to price ratio in the last 30 minutes
}
type VVList []VVData

// 计算单日成本价+成交量中值
func (k *VVData) Init(ohlcv OHLCVList, isOHLC bool, isV123 bool, isCbj bool) {
	if isOHLC {
		ohv := ohlcv.ToOHLCVData(true)
		k.TS = ohv.Timestamp.UnixMilli()
		// k.TS = ohv.Timestamp.Format("2006-01-02")
		k.Open = ohv.Open
		k.High = ohv.High
		k.Low = ohv.Low
		k.Close = ohv.Close
		k.Volume = ohv.Volume

		k.Hjj = (4*k.Close + 2*k.Open + k.High + k.Low) / 8
		k.Pvj = ohlcv.GetPvj(4, 2, 1, 1)
	}

	if isV123 {
		v123 := ohlcv.ToV123Data()
		k.V931 = v123.V931
		k.V932 = v123.V932
		k.V935 = v123.V935
		k.V940 = v123.V940
		k.V150 = v123.V150
	}

	if isCbj {
		cbj := ohlcv.ToCbjData("10:00:00")
		k.Cbj = cbj.Cbj
		k.Cb1 = cbj.Cb1
		k.Cb2 = cbj.Cb2

		k.Vmed = cbj.Vmed
		k.Nup = cbj.Nup
		k.Ndown = cbj.Ndown
	}
}

func CheckIndicators(indicators string) (isOHLC bool, isV123 bool, isCbj bool) {
	isOHLC = false
	isV123 = false
	isCbj = false
	if indicators == "" {
		return true, true, true
	}
	indicatorList := strings.Split(indicators, ",")
	indicatorOHLC := []string{"hjj", "pvj"}
	indicatorV123 := []string{"v931", "v932", "v935", "v940", "v150"}
	indicatorCbj := []string{"vmed", "cbj", "cb1", "cb2", "nup", "ndown"}
	for _, ind := range indicatorList {
		if slices.Contains(indicatorOHLC, ind) {
			isOHLC = true
		}
		if slices.Contains(indicatorV123, ind) {
			isV123 = true
		}
		if slices.Contains(indicatorCbj, ind) {
			isCbj = true
		}
	}

	return isOHLC, isV123, isCbj
}

// 计算单日成本价+成交量中值
func (k *VVData) ToRecord(isOHLC, isV123, isCbj, istimestamp bool) map[string]any {
	rec := make(map[string]any)
	// isOHLC, isV123, isCbj := CheckIndicators(indicators)
	if isOHLC || isV123 || isCbj {
		if !istimestamp {
			tz := time.FixedZone("CST", 8*3600)
			// ts, _ := time.ParseInLocation("2006-01-02", k.TS, tz)
			ts := MillisToTime(k.TS).In(tz)
			rec["timestamp"] = ts.Format("2006-01-02")
		} else {
			rec["timestamp"] = k.TS
		}
	}
	if isOHLC {
		rec["open"] = k.Open
		rec["high"] = k.High
		rec["low"] = k.Low
		rec["close"] = k.Close
		rec["volume"] = k.Volume
		rec["hjj"] = k.Hjj
		rec["pvj"] = k.Pvj
	}

	if isV123 {
		rec["v931"] = k.V931
		rec["v932"] = k.V932
		rec["v935"] = k.V935
		rec["v940"] = k.V940
		rec["v150"] = k.V150
	}

	if isCbj {
		rec["vmed"] = k.Vmed
		rec["cbj"] = k.Cbj
		rec["cb1"] = k.Cb1
		rec["cb2"] = k.Cb2
		rec["nup"] = k.Nup
		rec["ndown"] = k.Ndown
	}

	return rec
}

func (k *VVData) ToOHLCVData() (dd OHLCVData) {
	// ohv := ohlcv.ToOHLCVData(true)
	tz := time.FixedZone("CST", 8*3600)
	// ts, _ := time.ParseInLocation("2006-01-02", k.TS, tz)
	ts := MillisToTime(k.TS).In(tz)

	dd.Timestamp = ts
	dd.Open = k.Open
	dd.High = k.High
	dd.Low = k.Low
	dd.Close = k.Close
	dd.Volume = k.Volume

	return dd
}

func (k *VVData) ToV123Data() (dd V123Data) {
	// ohv := ohlcv.ToOHLCVData(true)
	tz := time.FixedZone("CST", 8*3600)
	ts := MillisToTime(k.TS).In(tz)
	// ts, _ := time.ParseInLocation("2006-01-02", k.TS, tz)

	dd.TS = ts
	dd.V931 = k.V931
	dd.V932 = k.V932
	dd.V935 = k.V935
	dd.V940 = k.V940
	dd.V150 = k.V150

	return dd
}

func (k *VVData) ToCbjData() (dd CbjData) {
	// ohv := ohlcv.ToOHLCVData(true)
	tz := time.FixedZone("CST", 8*3600)
	ts := MillisToTime(k.TS).In(tz)
	// ts, _ := time.ParseInLocation("2006-01-02", k.TS, tz)

	dd.TS = ts
	dd.Vmed = k.Vmed
	dd.Cbj = k.Cbj
	dd.Cb1 = k.Cb1
	dd.Cb2 = k.Cb2
	dd.Nup = k.Nup
	dd.Ndown = k.Ndown

	return dd
}

func (k *VVList) ToCbjList() (ll CbjList) {
	for _, kk := range *k {
		ll = append(ll, kk.ToCbjData())
	}
	return ll
}

func (k *VVList) ToV123List() (ll V123List) {
	for _, kk := range *k {
		ll = append(ll, kk.ToV123Data())
	}
	return ll
}

func (k *VVList) ToOHLCVList() (ll OHLCVList) {
	for _, kk := range *k {
		ll = append(ll, kk.ToOHLCVData())
	}
	return ll
}

// 将数据转换为map[string]any切片
func (k *VVList) ToRecords(isOHLC, isV123, isCbj, istimestamp bool) []map[string]any {
	var dd []map[string]any
	for _, kk := range *k {
		dd = append(dd, kk.ToRecord(isOHLC, isV123, isCbj, istimestamp))
	}
	return dd
}

func (k *VVList) Head(n int, isohlcv bool, iscb bool, isv123 bool) {
	if isohlcv {
		ohv := (*k).ToOHLCVList()
		ohv.Head(n)
	}

	if isv123 {
		v123 := (*k).ToV123List()
		v123.Head(n)
	}

	if iscb {
		cbj := (*k).ToCbjList()
		cbj.Head(n)
	}

}

func (k *VVList) Tail(n int, isohlcv bool, isv123 bool, iscb bool) {
	if isohlcv {
		ohv := (*k).ToOHLCVList()
		ohv.Tail(n)
	}

	if isv123 {
		v123 := (*k).ToV123List()
		v123.Tail(n)
	}

	if iscb {
		cbj := (*k).ToCbjList()
		cbj.Tail(n)
	}
}
