package basic

var Int int = 10
var PInt *int = &Int
var PPInt **int = &PInt

var Ints = []int{1, 2, 3}
var PInts = &Ints

func VarPointer(a int, b *int) {
	a = 100
	*b = 100
}

func ExchangeByPointer(a, b *int) {
	*a, *b = *b, *a
}

func ArrayPointer(a []int, b *[]int) {
	a[0] = 100
	(*b)[0] = 100
}

const (
	Male   = 1
	Female
)

type Person struct {
	Name   string
	Age    int
	Gender int
}

func StructPointer(p1 Person, p2 *Person)  {
	p1.Age = 20
	p2.Age = 20
}