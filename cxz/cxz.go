package cxz

import (
	"encoding/csv"
	"fmt"
	"io"
	"log"
	"os"

	"github.com/go-gota/gota/dataframe"
	"github.com/ulikunitz/xz"
)

// ReadCSVFile 读取本地 CSV 文件并返回内容
func ReadCSVFile(filePath string) ([][]string, error) {
	// 打开 CSV 文件
	file, err := os.Open(filePath)
	if err != nil {
		return nil, fmt.Errorf("无法打开文件: %v", err)
	}
	defer file.Close()

	// 创建 CSV reader
	reader := csv.NewReader(file)

	// 读取所有记录
	var records [][]string

	for {
		record, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, fmt.Errorf("读取 CSV 数据时出错: %v", err)
		}
		records = append(records, record)
	}

	return records, nil
}

// 更高效的版本，使用 ReadAll 方法
func ReadCSVFileEfficient(filePath string) ([][]string, error) {
	// 打开 CSV 文件
	file, err := os.Open(filePath)
	if err != nil {
		return nil, fmt.Errorf("无法打开文件: %v", err)
	}
	defer file.Close()

	// 创建 CSV reader
	reader := csv.NewReader(file)

	// 使用 ReadAll 一次性读取所有记录
	records, err := reader.ReadAll()
	if err != nil {
		return nil, fmt.Errorf("读取 CSV 数据时出错: %v", err)
	}

	return records, nil
}

// 将dataframe数据保存为csv文件
func SaveDataframeToCSVxz(df *dataframe.DataFrame, filePath string) error {
	file, err := os.Create(filePath)
	if err != nil {
		return fmt.Errorf("无法创建文件: %v", err)
	}
	defer file.Close()

	// 创建 xz.Writer
	xzWriter, err := xz.NewWriter(file)
	if err != nil {
		log.Fatalf("创建 xz.Writer 失败: %v", err)
	}
	defer xzWriter.Close()

	if err := df.WriteCSV(xzWriter); err != nil {
		log.Fatalf("写入压缩 csv.xz 失败: %v", err)
		return fmt.Errorf("写入压缩 csv.xz 失败: %v", err)
	}

	fmt.Printf("DataFrame 已保存为 %s", filePath)
	return nil
}

// 将dataframe数据保存为csv文件
func SaveDataframeToCSV(df *dataframe.DataFrame, filePath string) error {
	file, err := os.Create(filePath)
	if err != nil {
		return fmt.Errorf("无法创建文件: %v", err)
	}
	defer file.Close()

	if err := df.WriteCSV(file); err != nil {
		log.Fatalf("写入 CSV 失败: %v", err)
		return fmt.Errorf("写入 CSV 失败: %v", err)
	}

	fmt.Printf("DataFrame 已保存为 %s", filePath)

	// writer := csv.NewWriter(file)
	// defer writer.Flush()

	// // 写入标题行
	// titleRow := make([]string, len(df.Names()))
	// for i, name := range df.Names() {
	// 	titleRow[i] = name
	// }
	// err = writer.Write(titleRow)
	// if err != nil {
	// 	return fmt.Errorf("写入标题行失败: %v", err)
	// }

	// // 写入数据行
	// for _, row := range df.Values() {
	// 	err = writer.Write(row)
	// 	if err != nil {
	// 		return fmt.Errorf("写入数据行失败: %v", err)
	// 	}
	// }

	return nil
}
