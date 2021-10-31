package typedef

type Comparable interface {
	Compare(other interface{}) int
}
