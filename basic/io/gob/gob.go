package gob

import (
	"encoding/gob"
	"io"

	"github.com/google/uuid"
)

// 用于序列化和反序列化的结构体
type Product struct {
	Id     uuid.UUID
	Name   string
	Weight int
	Price  float64
}

// 生产一个结构体对象
func NewProduct(name string, weight int, price float64) (*Product, error) {
	id, err := uuid.NewUUID()
	if err != nil {
		return nil, err
	}
	return &Product{
		Id:     id,
		Name:   name,
		Weight: weight,
		Price:  price,
	}, nil
}

// 序列化当前结构体对象
func (p *Product) Serialize(w io.Writer) error {
	enc := gob.NewEncoder(w)
	return enc.Encode(p)
}

// 反序列化结构体对象
func DeserializeProduct(r io.Reader) (*Product, error) {
	dec := gob.NewDecoder(r)

	p := Product{}
	err := dec.Decode(&p)
	if err != nil {
		return nil, err
	}
	return &p, nil
}
