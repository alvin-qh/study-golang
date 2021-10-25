package builtin

import (
	"errors"
	"fmt"
	"math/rand"
	"time"
	"unicode/utf8"
)

func RandomString(length int) string {
	bytes := []byte("0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")
	result := make([]byte, 0, length)

	rand.Seed(time.Now().UnixNano() + int64(rand.Intn(100)))

	for i := 0; i < length; i++ {
		result = append(result, bytes[rand.Intn(len(bytes))])
	}
	return string(result)
}

type User struct {
	id   string
	name string
}

func New(id, name string) (user *User, err error) {
	if len(name) == 0 {
		err = errors.New("invalid name")
		return nil, err
	}
	return &User{id: id, name: name}, nil
}

func SafeNew(id, name string) *User {
	if user, err := New(id, name); err != nil {
		panic(err)
	} else {
		return user
	}
}

var (
	ErrIdRequired = errors.New("id required")
	ErrIdLength   = errors.New("id not enough length")
)

func checkRequired(user *User) error {
	if len(user.id) == 0 {
		return fmt.Errorf("id caused error: %w", ErrIdRequired)
	}
	return nil
}

func checkLength(user *User) error {
	if utf8.RuneCountInString(user.id) < 10 {
		return fmt.Errorf("id %s caused error: %w", user.id, ErrIdLength)
	}
	return nil
}

func (user *User) IsValid() error {
	if err := checkRequired(user); err != nil {
		return err
	}

	if err := checkLength(user); err != nil {
		return err
	}

	return nil
}
