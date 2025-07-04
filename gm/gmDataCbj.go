package gm

import (
	"fmt"
	"sort"
	"strings"
	"time"
)

type CbjData struct {
	// TS string `json:"timestamp"`
	TS time.Time `json:"timestamp"`

	Vmed int64   `json:"vmed"` // median volume
	Cbj  float64 `json:"cbj"`  // Price for volume avarged Hjj over Vmed
	Cb1  float64 `json:"cb1"`  // Average Price before 10:00
	Cb2  float64 `json:"cb2"`  // Average Price after 10:00

	Nup   int64 `json:"nup"`   // number of up ticks
	Ndown int64 `json:"ndown"` // number of down ticks
}
type CbjList []CbjData

// 显示
func (k *CbjData) Print(head bool) {
	if head {
		fmt.Println("Date\t\tvmed\tcbj\tcb1\tcb2\tnup\tndown\t")
		fmt.Println("----\t\t----\t---\t---\t---\t---\t-----\t")
	}
	fmt.Printf("%s\t%d\t%.2f\t%.2f\t%.2f\t%d\t%d\n",
		k.TS.Local().Format("2006-01-02"),
		k.Vmed, k.Cbj, k.Cb1, k.Cb2, k.Nup, k.Ndown)
}

func (k *CbjList) Head(n int) {
	num := min(len(*k), n)
	for i, kb := range (*k)[:num] {
		isHead := false
		if i == 0 {
			isHead = true
			fmt.Printf("CbjList Head(%d) of %d CbjData:\n", num, len(*k))
		}
		kb.Print(isHead)

		if i == num-1 {
			fmt.Println("=" + strings.Repeat("=", 64))
		}
	}
}

func (k *CbjList) Tail(n int) {
	num := min(len(*k), n)
	isHead := false
	for i, kb := range (*k)[(len(*k) - num):] {
		isHead = false
		if i == 0 {
			isHead = true
			fmt.Printf("CbjList Tail(%d) of %d CbjData:\n", num, len(*k))
		}

		kb.Print(isHead)

		if i == num-1 {
			fmt.Println("=" + strings.Repeat("=", 64))
		}
	}
}

func (k *CbjList) Sort(descend bool) {
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
