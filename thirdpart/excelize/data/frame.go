package data

import (
	"encoding/json"
	"os"

	df "github.com/go-gota/gota/dataframe"
)

// 读取 csv 文件, 返回 `DataFrame` 实例
func LoadCSV(path string) df.DataFrame {
	r, err := os.Open(path)
	if err != nil {
		return df.DataFrame{Err: err}
	}
	defer r.Close()

	return df.ReadCSV(r, df.HasHeader(true), df.DetectTypes(true))
}

// 求 `map[string][]any` 类型 Map 中所有值的最大长度
func maxLenOfMapValues(m map[string][]any) int {
	max := 0
	for _, vs := range m {
		if len(vs) > max {
			max = len(vs)
		}
	}
	return max
}

// 将 `map[string][]any` 类型 Map 转换为 `[]map[string]any` 类型 Map
//
// 例如:
//
//	map[string][]any{"A": {1, 2, 3}, "B": {4, 5, 6}} => []map[string]any{{A: 1, B: 4}, {A: 2, B: 5}, {A: 3, B: 6}}
//
// 后者可作为产生 `DataFrame` 的参数使用
func mapConvert(m map[string][]any) []map[string]any {
	max := maxLenOfMapValues(m)

	r := make([]map[string]any, max)
	for k, vs := range m {
		for i, v := range vs {
			if r[i] == nil {
				r[i] = make(map[string]any)
			}
			r[i][k] = v
		}
	}
	return r
}

// 读取 json 文件, 返回 `DataFrame` 实例
func LoadJSON(path string) df.DataFrame {
	data, err := os.ReadFile(path)
	if err != nil {
		return df.DataFrame{Err: err}
	}

	// 读取 json 文件并转换为 `map[string][]any` 类型 Map
	m := make(map[string][]any)
	err = json.Unmarshal(data, &m)
	if err != nil {
		return df.DataFrame{Err: err}
	}

	return df.LoadMaps(mapConvert(m), df.HasHeader(true), df.DetectTypes(true))
}
