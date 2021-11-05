package types

// User 结构体，所有成员字段为 public
type User struct {
	Id     int    `primaryKey:"true" null:"false"`
	Name   string `default:"Alvin"`
	Gender rune
}
