package math

import (
	"errors"
)

// 加法
func Add(a, b float64) float64 {
	return a + b
}

// 减法
func Subtract(a, b float64) float64 {
	return a - b
}

// 乘法
func Multiply(a, b float64) float64 {
	return a * b
}

// 除法（处理除以零的错误）
func Divide(a, b float64) (float64, error) {
	if b == 0 {
		return 0, errors.New("除数不能为零")
	}
	return a / b, nil
}
