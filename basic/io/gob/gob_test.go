package gob

import (
	"bytes"
	"encoding/gob"
	"reflect"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

// 测试对简单类型数据进行编码和解码
//
// 通过 `gob.NewEncoder(io.Writer)` 可以创建一个 `gob.Encoder` 类型的编码器实例, 用于将一个实例序列化为字节流
//
// 通过 `gob.NewDecoder(io.Reader)` 可以创建一个 `gob.Decoder` 类型的解码器实例, 用于将一个字节流反序列化为实例
//
// 序列化时可以指定实例的值或实例的指针, 反序列化时需要指定实例的指针
func TestGob_SimpleData(t *testing.T) {
	// 用于序列化的 Buffer 对象
	buf := bytes.NewBuffer(make([]byte, 0))

	// 生成编码器实例
	enc := gob.NewEncoder(buf)

	// 被序列化的整数实例
	en := 100

	// 序列化一个整数, 使用实例的值
	err := enc.Encode(en)
	assert.Nil(t, err)

	// 被序列化的字符串实例
	es := "Hello World"

	// 序列化一个字符串, 使用实例的指针
	err = enc.Encode(&es)
	assert.Nil(t, err)

	// 生成解码器对象
	dec := gob.NewDecoder(buf)

	// 用于接受反序列化整数的变量
	var an int

	// 反序列化整数
	err = dec.Decode(&an)
	assert.Nil(t, err)
	assert.Equal(t, en, an)

	// 用于接收反序列化字符串的变量
	var as string

	// 反序列化字符串
	err = dec.Decode(&as)
	assert.Nil(t, err)
	assert.Equal(t, es, as)
}

// 用于序列化和反序列化的结构体
type Product struct {
	Id     uuid.UUID
	Name   string
	Weight int
	Price  float64
}

// 测试结构体实例的序列化和反序列化
//
// 和序列化/反序列化简单类型类似, 可以对结构体实例的值或指针进行序列化, 并在之后通过变量指针反序列化到该变量中
func TestGob_Struct(t *testing.T) {
	ep := Product{
		Id:     uuid.New(),
		Name:   "Apple",
		Weight: 100,
		Price:  2.56,
	}

	// 用于序列化的 Buffer 对象
	buf := bytes.NewBuffer(make([]byte, 0))

	// 创建编码器
	enc := gob.NewEncoder(buf)

	// 通过结构体实例将结构体进行编码
	err := enc.Encode(ep)
	assert.Nil(t, err)

	// 通过结构体实例指针进行编码
	err = enc.Encode(&ep)
	assert.Nil(t, err)

	// 创建解码器
	dec := gob.NewDecoder(buf)

	// 用于接收解码结果的结构体实例
	var ap1, ap2 Product

	// 解码结构体
	err = dec.Decode(&ap1)
	assert.Nil(t, err)
	assert.Equal(t, ep, ap1)

	// 解码结构体
	err = dec.Decode(&ap2)
	assert.Nil(t, err)
	assert.Equal(t, ep, ap2)
}

// 包含未知类型字段的结构体
//
// `Product` 字段为 `interface{}` 类型, 不能直接进行序列化
type Order struct {
	Id      string
	Product interface{}
}

// 测试包含 `interface{}` 类型字段的结构体实例的序列化
//
// 如果结构体中包含了 `interface{}` 类型字段, 则其实例在序列化时会发生错误, 因为编码器并不知道该字段对应的实际类型
//
// 此时可通过 `gob.Register` 方法对类型进行注册, 从而让编码器可以正确识别类型
func TestGob_InterfaceType(t *testing.T) {
	// 未注册类型前, 对包含 `interface{}` 类型字段的结构体实例进行编码, 返回类型未注册错误
	t.Run("without register", func(t *testing.T) {
		buf := bytes.NewBuffer(make([]byte, 0))

		enc := gob.NewEncoder(buf)

		// 对 `Order` 结构体实例进行序列化
		err := enc.Encode(&Order{
			Id: "001",
			Product: Product{
				Id:     uuid.New(),
				Name:   "Apple",
				Price:  2.25,
				Weight: 100,
			},
		})

		// 序列化失败, 返回错误
		assert.EqualError(t, err, "gob: type not registered for interface: gob.Product")
	})

	// 注册了 `interface{}` 将表示的类型后, 对应的结构体实例方可正确序列化
	t.Run("with register", func(t *testing.T) {
		// 注册类型
		gob.Register(Product{})

		eo := Order{
			Id: "001",
			Product: Product{
				Id:     uuid.New(),
				Name:   "Apple",
				Price:  2.25,
				Weight: 100,
			},
		}

		buf := bytes.NewBuffer(make([]byte, 0))

		enc := gob.NewEncoder(buf)

		// 进行序列化, 此时序列化成功
		err := enc.Encode(&eo)
		assert.Nil(t, err)

		dec := gob.NewDecoder(buf)

		var ao Order

		// 进行反序列化
		err = dec.Decode(&ao)
		assert.Nil(t, err)
		assert.Equal(t, eo, ao)
	})
}

// 测试对 `reflect.Value` 实例进行序列化和反序列化
func TestGob_ReflectValue(t *testing.T) {
	buf := bytes.NewBuffer(make([]byte, 0))

	enc := gob.NewEncoder(buf)

	ep := Product{
		Id:     uuid.New(),
		Name:   "Apple",
		Price:  2.25,
		Weight: 100,
	}

	// 将一个 `reflect.Value` 类型实例进行序列化
	err := enc.EncodeValue(reflect.ValueOf(ep))
	assert.Nil(t, err)

	dec := gob.NewDecoder(buf)

	var ap Product

	// 将字节流反序列化到变量指针的 `reflect.Value` 类型实例中
	// 相当于反序列化到 `reflect.Value` 包装的变量中
	err = dec.DecodeValue(reflect.ValueOf(&ap))
	assert.Nil(t, err)
	assert.Equal(t, ep, ap)
}
