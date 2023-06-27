package serial

import (
	"encoding/xml"
	"fmt"
)

type User struct {
	XMLName xml.Name `json:"-" xml:"user"`     // 用于定义 XML 根节点名称, json 中忽略 (json:"-" 表示该字段不出现在 json 中)
	Id      int64    `json:"id" xml:"id,attr"` // "attr" 表示在 XML 中, "id" 字段的值在 根节点属性上表示, 而不是使用 XML 节点
	Name    string   `json:"name" xml:"name"`
	Email   string   `json:"email,omitempty" xml:"email,omitempty"` // "omitempty" 表示如果为空, 则不出现在 json 或 XML 中
	Phone   []string `json:"phone,omitempty" xml:"phones>tel"`      // "phones>tel" 表示在 XML 中, 切片类型字段位于 "phones" 节点下, 每一项是一个 "tel" 节点
}

func NewUser(id int64, name, email string, phone []string) *User {
	return &User{Id: id, Name: name, Email: email, Phone: phone}
}

func (u *User) AddPhone(phone string) {
	u.Phone = append(u.Phone, phone)
}

// 定义 convert 函数, 对 interface{} 类型值进行类型转换
func ConvertInterfaceToStringSlice(obj interface{}) ([]string, error) {
	is, ok := obj.([]interface{})
	if !ok {
		return nil, fmt.Errorf("invalid type")
	}

	os := make([]string, len(is))
	for i, v := range is {
		os[i] = v.(string)
	}
	return os, nil
}
