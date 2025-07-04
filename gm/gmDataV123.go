package gm

import (
	"fmt"
	"sort"
	"strings"
	"time"
)

type V123Data struct {
	// Timestamp int64   `json:"timestamp"`
	TS time.Time `json:"timestamp"`
	// V0 int64     `json:"v0"` // volume for 1d
	V931 int64 `json:"v931"` // volume at 9:31
	V932 int64 `json:"v932"` // volume at 9:32
	V935 int64 `json:"v935"` // volume at 9:35
	V940 int64 `json:"v940"` // volume at 9:40
	V150 int64 `json:"v150"` // volume at 15:00
}
type V123List []V123Data

// 显示
func (k *V123Data) Print(head bool) {
	if head {
		fmt.Println("Date\t\tv931\tv932\tv935\tv940\tv150\t")
		fmt.Println("----\t\t----\t----\t----\t----\t----\t")
	}
	fmt.Printf("%s\t%d\t%d\t%d\t%d\t%d\t\n",
		k.TS.Local().Format("2006-01-02"),
		k.V931, k.V932, k.V935, k.V940, k.V150)
}

func (k *V123List) Head(n int) {
	num := min(len(*k), n)
	for i, kb := range (*k)[:num] {
		isHead := false
		if i == 0 {
			isHead = true
			fmt.Printf("V123List Head(%d) of %d V123Data:\n", num, len(*k))
		}
		kb.Print(isHead)

		if i == num-1 {
			fmt.Println("=" + strings.Repeat("=", 64))
		}
	}
}

func (k *V123List) Tail(n int) {
	num := min(len(*k), n)
	isHead := false
	for i, kb := range (*k)[(len(*k) - num):] {
		isHead = false
		if i == 0 {
			isHead = true
			fmt.Printf("V123List Tail(%d) of %d V123Data:\n", num, len(*k))
		}

		kb.Print(isHead)

		if i == num-1 {
			fmt.Println("=" + strings.Repeat("=", 64))
		}
	}
}

func (k *V123List) Sort(descend bool) {
	// 按时间排序
	if descend {
		sort.Slice(*k, func(i, j int) bool {
			return (*k)[i].TS.After((*k)[j].TS)
		})
	} else {
		sort.Slice(*k, func(i, j int) bool {
			return (*k)[i].TS.Before((*k)[j].TS)
		})
	}
}
