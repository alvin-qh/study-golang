package builtin

// GO 集合操作
// GO 集合包括：数组、切片和列表

import (
	"basic/builtin/set"
	"container/list"
	"fmt"
	"reflect"
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
)

// 测试数组
func TestArray(t *testing.T) {
	// 创建长度为 10 的数组
	var a1 [5]int
	length := len(a1)                            // 获取数组长度
	assert.Equal(t, 5, length)                   // len 函数用于测量数组的长度
	assert.Equal(t, [...]int{0, 0, 0, 0, 0}, a1) // 数组元素的初始值均为 0

	// 遍历数组
	for i := 0; i < len(a1); i++ {
		a1[i] = i + 1 // 使用下标访问数组元素
	}
	assert.Equal(t, [...]int{1, 2, 3, 4, 5}, a1) // 数组实际的结果值

	// 通过 range 方法遍历数组，其中 index 为数组下标，value 为对应的元素值
	for index, value := range a1 {
		assert.Equal(t, index+1, value) // 查看遍历的结果值
	}

	// 可以通过赋值的方式创建数组
	a2 := [3]int{1} // 初始化数组的前 1 个元素
	assert.Len(t, a2, 3)
	assert.Equal(t, [...]int{1, 0, 0}, a2) // 除了显式初始化的元素外，其余元素值为 0

	// 隐式声明数组长度，初始化全部数组元素
	// '...' 表示数组长度由初始化元素个数决定
	// 注意：如果省略了省略号，则表示一个 切片 而非 数组
	a3 := [...]int{1, 2, 3}
	assert.Len(t, a3, 3)
	assert.Equal(t, [...]int{1, 2, 3}, a3)

	// 初始化多维数组
	a4 := [9][9]int{}
	assert.Len(t, a4, 9)    // 数组的第 1 维长度
	assert.Len(t, a4[0], 9) // 数组的第 2 维长度

	// 通过循环给数组赋值
	for i := 0; i < len(a4); i++ {
		for j := 0; j < len(a4[i]); j++ {
			a4[i][j] = (i + 1) * (j + 1)
		}

		assert.Equal(t, [...]int{
			(i + 1) * 1,
			(i + 1) * 2,
			(i + 1) * 3,
			(i + 1) * 4,
			(i + 1) * 5,
			(i + 1) * 6,
			(i + 1) * 7,
			(i + 1) * 8,
			(i + 1) * 9, // 注意，这里要多一个 逗号，表示参数换行
		}, a4[i])
	}

	// 定义一个元素为任意类型的数组
	a5 := [...]interface{}{"Hello", 1, false}
	assert.Len(t, a5, 3)
	assert.Equal(t, "string", reflect.TypeOf(a5[0]).Name())
	assert.Equal(t, "int", reflect.TypeOf(a5[1]).Name())
	assert.Equal(t, "bool", reflect.TypeOf(a5[2]).Name())

	// 测试数组的指针
	a6 := [...]int{1, 2, 3}
	pa6 := &a6                               // var pArray7 *[3]Any = &array7，获取数组的指针
	assert.Equal(t, [...]int{1, 2, 3}, *pa6) // 解引指针，获取数组

	pa6[0] = 10 // 通过指针改变数组的元素值
	assert.Equal(t, [...]int{10, 2, 3}, a6)

	// 表示数组的拷贝
	a7 := a6
	a7[0] = 100
	assert.NotEqual(t, a6, a7) // array7 和 array8 是两个不同的数组
}

// 测试切片
// 切片表示一个 可变长度 的数组的 引用类型
// 切片有两个长度 len, cap，前者表示切片的元素个数，后者表示切片的实际长度，超过 cap 后，切片需要重新分配
func TestSlice(t *testing.T) {
	// 创建切片类型 空 变量
	var s []int // slice := []int(nil)，创建一个切片类型变量
	assert.Nil(t, s)
	assert.Equal(t, 0, len(s)) // nil 的长度为 0

	// 创建一个空的切片
	s = []int{}
	assert.Len(t, s, 0) // 为 nil 的切片长度为 0

	s = append(s, 0)       // 向切片中添加元素
	s = append(s, 1, 2, 3) //向切片中添加多个元素
	assert.Equal(t, []int{0, 1, 2, 3}, s)
	assert.Equal(t, 4, cap(s))

	// 创建长度为 5 的切片并初始化元素
	s = []int{1, 2, 3, 4, 5}
	assert.Equal(t, []int{1, 2, 3, 4, 5}, s)
	assert.Equal(t, 5, cap(s)) //  切片 cap 为 5

	// 通过 make 函数初始化切片，初始长度 3
	// 切片无法通过 [n]int{} 创建，这和创建数组的语法冲突
	s = make([]int, 3)
	assert.Equal(t, []int{0, 0, 0}, s) // 切片初始长度为 3
	assert.Equal(t, 3, cap(s))         // 数组初始容积为 3

	s[0], s[1], s[2] = 100, 200, 300 // 通过下标给切片赋值
	assert.Equal(t, []int{100, 200, 300}, s)

	s = append(s, 400)                            // 向切片中添加元素
	assert.Equal(t, s, []int{100, 200, 300, 400}) // 切片长度为 4
	assert.Equal(t, 6, cap(s))                    // 数组 cap 增长为 6

	s = append(s, 500, 600, 700) // 向切片中添加多个元素
	assert.Equal(t, len(s), 7)   // 切片的长度变为 7
	assert.Equal(t, cap(s), 12)  // 数组 cap 增长为 12

	s = make([]int, 0, 100) // 创建切片并指定 cap
	assert.Len(t, s, 0)
	assert.Equal(t, 100, cap(s))

	array := [...]int{1, 2, 3, 4, 5}

	// 从数组（或另一个切片）中获取切片
	s = array[:2]
	assert.Equal(t, []int{1, 2}, s) // 切片为数组的前两个元素

	s = array[2:]
	assert.Equal(t, []int{3, 4, 5}, s) // 切片为数组的后三个元素

	s = array[2:3]
	assert.Equal(t, []int{3}, s) // 切片为数组的低桑额元素

	s = append(s, 4) // 为切片增加新的元素
	assert.Equal(t, []int{3, 4}, s)

	// 删除第 3 个元素
	// ...运算符用于将切片展开成参数
	s = []int{1, 2, 3, 4, 5}
	s = append(s[:2], s[3:]...) // 先取前 2 个元素的切片，在其之上添加后 2 个元素，相当于删除第 3 个元素
	assert.Equal(t, []int{1, 2, 4, 5}, s)

	// 切片引用
	// 切片是引用类型，赋值只会传递切片的引用
	s1 := []int{1, 2, 3}
	s2 := s1
	assert.Equal(t, []int{1, 2, 3}, s2) // 赋值运算符会传递切片的引用

	s2[1] = 20
	assert.Equal(t, s1, s2)              // 两个引用指向了同一个切片
	assert.Equal(t, []int{1, 20, 3}, s1) // 赋值运算无法复制切片
	assert.Equal(t, []int{1, 20, 3}, s2) // 赋值运算无法复制切片

	s2 = append(s2, 4)
	assert.NotEqual(t, s1, s2)              // 两个引用指向了同一个切片
	assert.Equal(t, []int{1, 20, 3}, s1)    // 赋值运算无法复制切片
	assert.Equal(t, []int{1, 20, 3, 4}, s2) // 赋值运算无法复制切片

	// 切片拷贝
	// 通过 copy 函数可以将切片元素复制到目标切片中
	// 复制的元素个数以两个切片中长度较小的为准
	s1 = []int{1, 2, 3}
	s2 = make([]int, len(s1))
	copy(s2, s1)                        // 将 s1 的元素复制到 s2 中
	assert.Equal(t, []int{1, 2, 3}, s2) // copy 会复制切片的内容

	s2 = make([]int, 4)
	copy(s2, s1)
	assert.Equal(t, []int{1, 2, 3, 0}, s2) // 因为 s1 长度较小，所以会复制 s1 的全部元素，并保留 s2 的多余元素

	s2 = make([]int, 2)
	copy(s2, s1)
	assert.Equal(t, []int{1, 2}, s2) // 因为 s2 长度较小，所以会复制 s1 中和 s2 长度匹配的那部分，其余的不复制

	// 多维度切片
	var s3 [][]int        // 声明多维度切片类型变量
	s3 = make([][]int, 0) // 初始化变量
	s3 = append(s3, []int{1, 2, 3})
	assert.Equal(t, [][]int{{1, 2, 3}}, s3)

	s4 := [][]int{} // 声明并初始化一个空的多维度切片
	s4 = append(s4, []int{1, 2, 3}, []int{4, 5, 6})
	assert.Equal(t, [][]int{{1, 2, 3}, {4, 5, 6}}, s4)
}

// 测试切片 cap 的增长率
func TestSliceCapGrowUp(t *testing.T) {
	var slice []int
	assert.Equal(t, len(slice), cap(slice), 0)

	c := 1
	for i := 0; i < 20; i++ {
		// 当切片不为空时，每当 cap 和 len 相等，在添加元素时，会分配原有长度 2 倍的空间作为新的 cap
		if len(slice) > 0 && len(slice) == cap(slice) {
			c = len(slice) * 2
		}

		slice = append(slice, i)
		assert.Equal(t, len(slice), i+1)
		assert.Equal(t, cap(slice), c)
	}
}

// 定义结构体作为 map key
type UserKey struct {
	id   int
	name string
}

// 定义结构体作为 map value
type UserValue struct {
	gender   rune
	birthday string
	address  string
}

// 测试 map 数据结构
// map 即 hash map，是通过 hash 运算存储 key 和 value 键值对的结构
func TestMap(t *testing.T) {
	// 定义 map 类型变量
	var m1 map[string]int
	assert.Nil(t, m1)           // 变量此时为 nil
	assert.Equal(t, 0, len(m1)) // nil map 的长度为 0

	m1 = make(map[string]int, 100) // 通过 make 函数初始化 map, 第二个参数为 map 的初始容积，默认为 0
	m1["a"] = 1                    // 存放键值对
	assert.Equal(t, 1, m1["a"])    // 通过键获取值
	assert.Equal(t, 1, len(m1))    // 此时 map 长度为 1

	delete(m1, "a")             // 删除键值对
	assert.Equal(t, 0, m1["a"]) // 对于 value 为数值类型，返回 0 表示 key 不存在
	assert.Equal(t, 0, len(m1)) // 此时 map 长度为 0

	_, exist := m1["a"]
	assert.False(t, exist) // 判断 key 是否存在

	// 定义和初始化同时进行
	m2 := make(map[string]string)
	m2["a"] = "A"
	assert.Equal(t, "A", m2["a"]) // 通过键获取值

	delete(m2, "a")
	assert.Equal(t, "", m2["a"]) // 对于 value 为 string 类型，返回 空字符串 表示 key 不存在

	// 直接初始化
	m3 := map[string]interface{}{"a": 100, "b": "B", "c": []int{1, 2, 3}}
	assert.Equal(t, 100, m3["a"])
	assert.Equal(t, "B", m3["b"])
	assert.Equal(t, []int{1, 2, 3}, m3["c"])

	delete(m3, "a")
	assert.Equal(t, nil, m3["a"]) // 对于 value 为 其它 类型，返回 nil 表示 key 不存在

	// 遍历 map
	m4 := map[string]interface{}{"a": 100, "b": "B", "c": []int{1, 2, 3}}

	keys, values := make([]string, 0, len(m4)), make([]interface{}, 0, len(m4))
	for k, v := range m4 { // 遍历 key / value
		keys = append(keys, k)
		values = append(values, v)
	}
	assert.ElementsMatch(t, []string{"a", "b", "c"}, keys)
	assert.ElementsMatch(t, []interface{}{100, "B", []int{1, 2, 3}}, values)

	keys = make([]string, 0, len(m4))
	for k := range m4 { // 遍历所有 key
		keys = append(keys, k)
	}
	assert.ElementsMatch(t, []string{"a", "b", "c"}, keys)

	values = make([]interface{}, 0, len(m4))
	for _, v := range m4 { // 遍历所有 value
		values = append(values, v)
	}
	assert.ElementsMatch(t, []interface{}{100, "B", []int{1, 2, 3}}, values)

	// 使用其它类型的 key
	// 参考 builtin.UserKey, builtin.UserValue
	m5 := map[UserKey]*UserValue{} // 定义 map 使用结构体对象作为 key
	m5[UserKey{1, "Alvin"}] = &UserValue{gender: 'M', birthday: "1981-03", address: "ShanXi, Xi'an"}
	assert.Equal(t, UserValue{gender: 'M', birthday: "1981-03", address: "ShanXi, Xi'an"}, *m5[UserKey{1, "Alvin"}])
}

// 测试同步 sync.Map
// 同步 Map 用于异步场合，当多个任务同时访问一个 map 时，必须使用锁，否则会导致错误
// 如果需要降低锁对性能的影响，则需要使用 sync.Map 进行操作
func TestSyncMap(t *testing.T) {
	// 定义
	m := sync.Map{}

	// 定义一个等待组
	w := sync.WaitGroup{}

	w.Add(1)    // 表示等待数 加1
	go func() { // 开启一个任务（通过执行匿名函数）
		defer w.Done() // 表示任务结束后，等待数 减1

		for n := 0; n < 1000; n++ {
			m.Store(n, n+1) // 向 同步map 中添加元素
		}
	}()

	w.Add(1)
	go func() {
		defer w.Done()

		for n := 0; n < 1000; n++ {
			m.Store(fmt.Sprintf("%d", n), fmt.Sprintf("%d", n+1)) // 向 同步map 中添加元素
		}
	}()

	w.Wait() // 等待组的数值为 0 时返回

	// 通过 key 读取 value
	v, ok := m.Load("999")
	assert.True(t, ok) // ok 表示 key 是否存在
	assert.Equal(t, "1000", v)

	v, ok = m.Load(999)
	assert.True(t, ok)
	assert.Equal(t, 1000, v)

	// 通过 key 删除 key / value
	m.Delete("999")
	_, ok = m.Load("999")
	assert.False(t, ok)

	// 通过 key 读取并同时删除
	v, ok = m.LoadAndDelete(999)
	assert.True(t, ok) // ok 表示 key 是否存在
	assert.Equal(t, 1000, v)

	_, ok = m.Load(999)
	assert.False(t, ok)

	// 通过 key 读取 Value 否则存储 新Value
	_, ok = m.LoadOrStore(999, 1000)
	assert.False(t, ok) // ok 表示存储新 Value 前 key 是否存在

	v, ok = m.Load(999)
	assert.True(t, ok)
	assert.Equal(t, 1000, v)

	// 遍历 key / value
	keys := make([]interface{}, 0, 1000)
	values := make([]interface{}, 0, 1000)
	m.Range(func(k, v interface{}) bool { // 遍历需要通过传递一个函数参数完成
		keys = append(keys, k)
		values = append(values, v)
		return true // 返回遍历是否结束, 任意一个迭代返回 false，则整个遍历结束
	})
	assert.Contains(t, keys, 999)
	assert.Contains(t, keys, "0")

	assert.Contains(t, values, 1000)
	assert.Contains(t, values, "1")
}

// 将列表转化为切片
func toSlice(l *list.List) []interface{} {
	// 生成一个 cap 为列表长度的空切片
	slice := make([]interface{}, 0, l.Len())

	// 遍历列表，将列表元素依次加入切片中
	for iter := l.Front(); iter != nil; iter = iter.Next() {
		slice = append(slice, iter.Value)
	}
	return slice
}

// 将列表元素反转，得到新列表
func reverse(l *list.List) *list.List {
	rl := list.New()
	for iter := l.Back(); iter != nil; iter = iter.Prev() {
		rl.PushBack(iter.Value)
	}
	return rl
}

// go 语言的 list 实际上是双向链表，适合一些需要插入或删除中间节点的集合操作
// list 也能作为队列或栈来使用
func TestList(t *testing.T) {
	// 通过 New 创建一个空列表
	lst := list.New()
	assert.Equal(t, 0, lst.Len()) // 列表长度目前为 0

	// 添加元素
	// 返回值是一个 Element 指针，表示链表的节点
	elem := lst.PushBack(1) // 在列表末尾添加一个元素
	assert.Equal(t, []interface{}{1}, toSlice(lst))
	assert.Equal(t, elem.Value, 1) // 返回的 Element 即为刚添加元素的节点

	elem = lst.PushFront("Hello") // 在列表开头添加一个元素
	assert.Equal(t, []interface{}{"Hello", 1}, toSlice(lst))
	assert.Equal(t, elem.Value, "Hello")

	// 插入元素
	// 插入元素表示将新的元素添加在已有 Element 之前（或之后）
	elem = lst.Back().Prev() // 获取列表末尾元素的前一个元素节点
	assert.Equal(t, "Hello", elem.Value)

	lst.InsertAfter("OK", elem) // 在节点前插入
	assert.Equal(t, []interface{}{"Hello", "OK", 1}, toSlice(lst))

	elem = lst.Front().Next()
	assert.Equal(t, "OK", elem.Value)

	lst.InsertAfter("Bye", elem)
	assert.Equal(t, []interface{}{"Hello", "OK", "Bye", 1}, toSlice(lst))

	// 删除元素
	// 删除元素依赖被删除元素的节点对象，所以要先找到这个节点
	elem = lst.Front().Next().Next()
	assert.Equal(t, "Bye", elem.Value)

	value := lst.Remove(elem) // 删除节点，返回节点的 Value
	assert.Equal(t, "Bye", value)
	assert.Equal(t, []interface{}{"Hello", "OK", 1}, toSlice(lst))

	// 连接两个列表
	lst.PushBackList(reverse(lst)) // 在列表后连接列表
	assert.Equal(t, []interface{}{"Hello", "OK", 1, 1, "OK", "Hello"}, toSlice(lst))

	lst.PushFrontList(reverse(lst)) // 在列表前连接列表
	assert.Equal(t, []interface{}{"Hello", "OK", 1, 1, "OK", "Hello", "Hello", "OK", 1, 1, "OK", "Hello"}, toSlice(lst))

	// 重新初始化列表（清空）
	lst.Init()
	assert.Equal(t, 0, lst.Len())
}

// 测试 Set 集合
func TestSet(t *testing.T) {
	// 初始化并添加元素
	s1 := set.New(100)           // 初始化
	s1.Add(1, 2, 3, 4, 2)        // 批量添加元素
	assert.Equal(t, 4, s1.Len()) // 实际添加了 4 个元素，重复的 2 只存在 1 份

	// 判断集合是否包含指定值
	ok := s1.Contains(1)
	assert.True(t, ok)

	ok = s1.Contains(2, 3, 4) // 多值判断
	assert.True(t, ok)

	ok = s1.Contains(3, 4, 5) // 5 不在集合中，返回 false
	assert.False(t, ok)

	// 集合相等判断
	s2 := set.New(10) // 产生一个元素相同的集合
	s2.Add(1, 2, 3, 4)

	ok = s2.Equal(s1) // 两个集合元素是否相同
	assert.True(t, ok)

	s2.Add("Hello")   // 在其中一个集合中添加新元素
	ok = s2.Equal(s1) // 此时两个集合不再相同
	assert.False(t, ok)

	s2.Remove("Hello") // 移除之前添加的新元素
	ok = s2.Equal(s1)  // 此时两个集合恢复相同
	assert.True(t, ok)

	// 判断是否为子集
	ok = s2.IsSubset(s1) // 判断两个相同的集合是否互为子集
	assert.True(t, ok)
	ok = s1.IsSubset(s2)
	assert.True(t, ok)

	s1.Remove(2) // 删除集合元素，此时两个集合不再互为子集
	ok = s2.IsSubset(s1)
	assert.False(t, ok)
	ok = s1.IsSubset(s2)
	assert.True(t, ok)
}
