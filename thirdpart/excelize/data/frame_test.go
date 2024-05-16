package data

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

// 测试将 CSV 文件读取为 `DataFrame`
func TestLoadCSV(t *testing.T) {
	frame := LoadCSV("../files/simple_data.csv")
	assert.NoError(t, frame.Error())

	fmt.Println(frame)

	assert.Equal(t, [][]string{
		{"A", "B", "C", "D", "E"},
		{"0.100000", "0.200000", "0.300000", "0.400000", "0.500000"},
		{"1.100000", "1.200000", "1.300000", "1.400000", "1.500000"},
		{"2.100000", "2.200000", "2.300000", "2.400000", "2.500000"},
		{"3.100000", "3.200000", "3.300000", "3.400000", "3.500000"},
		{"4.100000", "4.200000", "4.300000", "4.400000", "4.500000"},
	}, frame.Records())
}

// 测试将 JSON 文件读取为 `DataFrame`
func TestLoadJSON(t *testing.T) {
	frame := LoadJSON("../files/simple_data.json")
	assert.NoError(t, frame.Error())

	fmt.Println(frame)

	assert.Equal(t, [][]string{
		{"A", "B", "C", "D", "E"},
		{"0.100000", "0.200000", "0.300000", "0.400000", "0.500000"},
		{"1.100000", "1.200000", "1.300000", "1.400000", "1.500000"},
		{"2.100000", "2.200000", "2.300000", "2.400000", "2.500000"},
		{"3.100000", "3.200000", "3.300000", "3.400000", "3.500000"},
		{"4.100000", "4.200000", "4.300000", "4.400000", "4.500000"},
	}, frame.Records())
}
