package maptostruct

import (
	"reflect"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

type Gender string

const (
	GenderM = Gender("M")
	GenderF = Gender("F")
)

type user struct {
	Id       int       `json:"id"`
	Name     string    `json:"name"`
	Gender   Gender    `json:"gender"`
	Birthday time.Time `json:"birthday"`
	Titles   []string  `json:"titles"`
	Locale   struct {
		Country  string `json:"country"`
		Language string `json:"language"`
	} `json:"locale"`
}

func TestFindTag(t *testing.T) {
	f, _ := reflect.TypeOf(new(user)).Elem().FieldByName("Id")

	mts := New("json")
	assert.Equal(t, "id", mts.findTag(&f))

	mts = New("unknown")
	assert.Equal(t, "id", mts.findTag(&f))
}

func TestDecodeByInvalidTarget(t *testing.T) {
	mts := New("json")

	v := 123
	m := map[string]any{}
	err := mts.Decode(m, &v)
	assert.EqualError(t, err, "\"target\" argument must be a struct pointer")
}

func TestDecodeStruct(t *testing.T) {
	m := map[string]any{
		"id":       100,
		"name":     "Alvin",
		"gender":   "M",
		"birthday": "1981-03-17",
		"titles":   []string{"Manager", "Engineer"},
		"locale": map[string]any{
			"country":  "China",
			"language": "Chinese",
		},
	}

	u := new(user)

	mts := New("json")

	err := mts.Decode(m, u)
	assert.NoError(t, err)

	assert.Equal(t, 100, u.Id)
	assert.Equal(t, "Alvin", u.Name)
	assert.Equal(t, GenderM, u.Gender)
	assert.Equal(t, time.Date(1981, 3, 17, 0, 0, 0, 0, time.UTC), u.Birthday)
	assert.Equal(t, []string{"Manager", "Engineer"}, u.Titles)
	assert.Equal(t, "China", u.Locale.Country)
	assert.Equal(t, "Chinese", u.Locale.Language)
}

func TestDecodeSlice(t *testing.T) {
	m := []any{
		map[string]any{
			"id":       100,
			"name":     "Alvin",
			"gender":   "M",
			"birthday": "1981-03-17",
			"titles":   []string{"Manager", "Engineer"},
			"locale": map[string]any{
				"country":  "China",
				"language": "Chinese",
			},
		},
		map[string]any{
			"id":       101,
			"name":     "Emma",
			"gender":   "F",
			"birthday": "1985-03-29",
			"titles":   []string{"Manager", "Engineer"},
			"locale": map[string]any{
				"country":  "China",
				"language": "Chinese",
			},
		},
	}

	mts := New("json")

	var us []user

	err := mts.Decode(m, &us)
	assert.NoError(t, err)
	assert.Len(t, us, 2)

	assert.Equal(t, 100, us[0].Id)
	assert.Equal(t, "Alvin", us[0].Name)
	assert.Equal(t, GenderM, us[0].Gender)
	assert.Equal(t, time.Date(1981, 3, 17, 0, 0, 0, 0, time.UTC), us[0].Birthday)
	assert.Equal(t, []string{"Manager", "Engineer"}, us[0].Titles)
	assert.Equal(t, "China", us[0].Locale.Country)
	assert.Equal(t, "Chinese", us[0].Locale.Language)

	assert.Equal(t, 101, us[1].Id)
	assert.Equal(t, "Emma", us[1].Name)
	assert.Equal(t, GenderF, us[1].Gender)
	assert.Equal(t, time.Date(1985, 3, 29, 0, 0, 0, 0, time.UTC), us[1].Birthday)
	assert.Equal(t, []string{"Manager", "Engineer"}, us[1].Titles)
	assert.Equal(t, "China", us[1].Locale.Country)
	assert.Equal(t, "Chinese", us[1].Locale.Language)
}
