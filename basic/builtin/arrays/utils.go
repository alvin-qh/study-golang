package arrays

import "reflect"

// 获取数组元素类型名称
//
// 参数:
// `array`: 数组对象引用
func GetTypeNameOfArrayElement(array []any) []string {
	// 创建一个空切片, 用于存储数组元素类型名称
	types := make([]string, 0, len(array))

	// 遍历数组元素, 获取每个元素的类型名称并添加到切片中
	for _, v := range array {
		types = append(types, reflect.TypeOf(v).Name())
	}

	// 返回类型名称切片
	return types
}
