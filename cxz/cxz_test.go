package cxz

import (
	"fmt"
	"testing"
)

func TestCxz(t *testing.T) {
	t.Log("cxz")
}

func TestReadCSVFile(t *testing.T) {
	// t.Log("cxz2")
	filePath := "data.csv" // 替换为你的 CSV 文件路径
	fmt.Println(" >>> Start to test script ...")
	// 使用第一个版本
	records, err := ReadCSVFile(filePath)
	if err != nil {
		fmt.Printf("错误: %v\n", err)
		return
	}
	fmt.Println("CSV 文件内容:")
	for i, record := range records {
		fmt.Printf("%d: %v\n", i+1, record)
	}
}

func TestReadCSVFileEfficient(t *testing.T) {
	// t.Log("cxz2")
	filePath := "data.csv" // 替换为你的 CSV 文件路径
	fmt.Println(" >>> Start to test script ...")
	// 使用第一个版本
	records, err := ReadCSVFileEfficient(filePath)
	if err != nil {
		fmt.Printf("错误: %v\n", err)
		return
	}
	fmt.Println("CSV 文件内容:")
	for i, record := range records {
		fmt.Printf("%d: %v\n", i+1, record)
	}
}
