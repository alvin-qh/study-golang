package interfaces

// 定义接口, 将对象转为字符串
type ToString interface {
	String() string
}

// 定义接口, 用于比较两个对象
type Comparable interface {
	Compare(other interface{}) int
}

// 定义一组通过 typedef.Comparable 接口对对象进行比较的函数
func Eq(left, right Comparable) bool {
	return left.Compare(right) == 0
}

func Ne(left, right Comparable) bool {
	return left.Compare(right) != 0
}

func Gt(left, right Comparable) bool {
	return left.Compare(right) > 0
}

func Lt(left, right Comparable) bool {
	return left.Compare(right) < 0
}

func Ge(left, right Comparable) bool {
	return left.Compare(right) >= 0
}

func Le(left, right Comparable) bool {
	return left.Compare(right) <= 0
}
