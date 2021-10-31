package typedef

// 定义接口，用于比较两个对象
type Comparable interface {
	Compare(other interface{}) int
}
