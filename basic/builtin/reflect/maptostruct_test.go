package reflect

import (
	"basic/builtin/reflect/maptostruct"
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

func TestMapToStruct(t *testing.T) {
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

	err := maptostruct.Decode(m, u)
	assert.NoError(t, err)

	assert.Equal(t, 100, u.Id)
	assert.Equal(t, "Alvin", u.Name)
	assert.Equal(t, GenderM, u.Gender)
	assert.Equal(t, time.Date(1981, 3, 17, 0, 0, 0, 0, time.UTC), u.Birthday)
	assert.Equal(t, []string{"Manager", "Engineer"}, u.Titles)
	assert.Equal(t, "China", u.Locale.Country)
	assert.Equal(t, "Chinese", u.Locale.Language)
}
