package core

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// 创建 Map 对象
//
// 有三种方法可以创建 Map 实例:
// 1. 使用 `var` 关键字定义空 Map, 例如 `var m map[string]int`, 这种方式无法为 Map 增加 Key/Value;
// 2. 定义并初始化 Map, 例如 `m := map[string]int{"a": 100, "b": 200}`;
// 2. 使用 `make` 函数创建 Map, 例如 `m := make(map[string]int)`;
func TestMap_New(t *testing.T) {
	// 定义一个 Map 变量
	var m map[string]int

	// 变量此时为 nil
	assert.Nil(t, m)
	// 值为 nil 的 Map 的长度为 0
	assert.Equal(t, 0, len(m))

	// 定义 Map 并初始化
	m = map[string]int{
		"a": 100,
		"b": 200,
	}

	// Map 中包含 2 个 Key
	assert.Equal(t, 2, len(m))
	// 根据 Key 获取 Value
	assert.Equal(t, 100, m["a"])
	// 根据 Key 获取 Value
	assert.Equal(t, 200, m["b"])

	// 通过 make 函数初始化 map, 第二个参数为 Map 的初始容积, 默认为 0
	m = make(map[string]int, 100)
	// 此时 Map 长度为 0
	assert.Equal(t, 0, len(m))
}

// 从 Map 中添加或删除 Key
//
// 通过 `m[key] = value` 方式添加 Key/Value, 如果 Key 已经存在, 则会覆盖原有的 Value
//
// 通过 `delete(m, key)` 方式删除 Key/Value, 如果 Key 不存在, 则不会报错
//
// 从 Map 中根据 Key 获取 Value 有两种形式, 即:
//
//	m := map[string]int{"a": 100}
//	v := m["a"]
//	v, ok := m["a"]
//
// 第一种形式根据 Key 获取 Value, 第二种形式返回 Value 和一个布尔值, 表示 Key 是否存在
func TestMap_AddDeleteKey(t *testing.T) {
	// 声明一个空 map
	m := map[string]int{}

	// 设置 key "a"
	m["a"] = 100

	// 获取 key "a" 是否存在, 以及其值
	v, exist := m["a"]
	assert.True(t, exist)
	assert.Equal(t, 100, v)

	// 删除 key "a"
	delete(m, "a")

	// 判断 key "a" 是否存在
	_, exist = m["a"]
	assert.False(t, exist)
}

// 遍历 Map
//
// 可以通过 `range` 循环遍历 Map, 循环变量 `key` 表示 Map 的 Key, 循环变量 `value` 表示 Map 的 Value, 例如:
//
//	m := map[string]int{"a": 100, "b": 200}
//	for key, value := range m {
//	    fmt.Println(key, value)
//	}
func TestMap_Through(t *testing.T) {
	m := map[string]any{
		"a": 100,
		"b": "B",
		"c": []int{1, 2, 3},
	}

	// 遍历数组
	ks := make([]string, 0, len(m))
	vs := make([]any, 0, len(m))

	// 遍历 Key/Value
	for k, v := range m {
		ks = append(ks, k)
		vs = append(vs, v)
	}

	assert.ElementsMatch(t, []string{"a", "b", "c"}, ks)
	assert.ElementsMatch(t, []any{100, "B", []int{1, 2, 3}}, vs)

	// 遍历所有 Key
	ks = make([]string, 0, len(m))
	for k := range m {
		ks = append(ks, k)
	}
	assert.ElementsMatch(t, []string{"a", "b", "c"}, ks)

	// 遍历所有的 Value
	vs = make([]any, 0, len(m))
	for _, v := range m {
		vs = append(vs, v)
	}
	assert.ElementsMatch(t, []any{100, "B", []int{1, 2, 3}}, vs)
}
