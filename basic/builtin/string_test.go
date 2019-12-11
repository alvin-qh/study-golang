package builtin

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestString(t *testing.T) {
	s := "Hello, 大家好"
	assert.Equal(t, rune(s[1]), 'e')
	assert.Equal(t, string(s[1]), "e")

	as := []rune(s)
	assert.Equal(t, as[1], int32(s[1]))
	assert.NotEqual(t, as[8], int32(s[8]))
	assert.Equal(t, as[8], '家')
}

func TestFromInt(t *testing.T) {
	assert.Equal(t, string(FromInt(123)), "123")
}

func TestFromFloat(t *testing.T) {
	assert.Equal(t, string(FromFloat(123.123)), "123.123")
}

func TestString_Len(t *testing.T) {
	s := String("Hello, 大家好")
	assert.Equal(t, s.Len(), 10)
}

func TestString_CharAt(t *testing.T) {
	s := String("Hello, 大家好")
	assert.Equal(t, s.CharAt(1), Char('e'))
	assert.Equal(t, s.CharAt(8), Char('家'))
}

func TestString_ToInt(t *testing.T) {
	n, e := String("123").ToInt()
	assert.Equal(t, n, 123)
	assert.Equal(t, e, nil)

	n, e = String("ABC").ToInt()
	assert.Equal(t, n, 0)
	assert.NotEqual(t, e, nil)
}

func TestString_ToFloat(t *testing.T) {
	n, e := String("123.123").ToFloat()
	assert.Equal(t, n, float32(123.123))
	assert.Equal(t, e, nil)

	n, e = String("ABC").ToFloat()
	assert.Equal(t, n, float32(0))
	assert.NotEqual(t, e, nil)
}

func TestString_Compare(t *testing.T) {
	assert.Equal(t, String("abc").Compare("abc"), 0)
	assert.Equal(t, String("abc").Compare("def"), -1)
	assert.Equal(t, String("def").Compare("abc"), 1)
}

func TestString_Contains(t *testing.T) {
	assert.True(t, String("abcde").Contains("cde"))
	assert.False(t, String("abcde").Contains("cdf"))
}

func TestString_Index(t *testing.T) {
	assert.Equal(t, String("abcde").Index("cde"), 2)
	assert.Equal(t, String("abcde").LastIndex("cde"), 2)
}

func TestString_Count(t *testing.T) {
	s := String("abababc")
	assert.Equal(t, s.Count("ab"), 3)
	assert.Equal(t, s.Count(""), s.Len()+1)
}

func TestRepeat(t *testing.T) {
	s := Repeat("abc", 3)
	assert.Equal(t, s.Len(), 9)
	assert.Equal(t, s.Count("abc"), 3)
}

func TestString_Replace(t *testing.T) {
	s := String("Hello")

	s = s.Replace("l", "L", 0) // do nothing replace
	assert.Equal(t, string(s), "Hello")

	s = s.Replace("l", "L", 1)
	assert.Equal(t, string(s), "HeLlo")

	s = String("Hello")
	s = s.Replace("l", "L", 2)
	assert.Equal(t, string(s), "HeLLo")

	s = String("Hello")
	s = s.Replace("ello", "ELLO", -1)
	assert.Equal(t, string(s), "HELLO")

	s = String("Hello")
	s.ReplaceSelf("ello", "ELLO", -1)
	assert.Equal(t, string(s), "HELLO")
}

func TestString_Trim(t *testing.T) {
	s := String("   Hello   ")
	assert.Equal(t, string(s.Trim()), "Hello")
}

func TestString_TrimWith(t *testing.T) {
	s := String("*^%Hello^&*")
	assert.Equal(t, string(s.TrimWith("*^%&")), "Hello")
}

func TestString_StartWith_And_EndWith(t *testing.T) {
	s := String("http://www.google.com")
	assert.True(t, s.StartWith("http://"))
	assert.True(t, s.EndWith(".com"))

	s = s.TrimStart("http://").TrimEnd(".com")
	assert.Equal(t, string(s), "www.google")
}

func TestString_Split(t *testing.T) {
	s := String("www.google.com")
	rs := s.Split(".")
	assert.Equal(t, len(rs), 3)
	assert.Equal(t, string(rs[0]), "www")
	assert.Equal(t, string(rs[1]), "google")
	assert.Equal(t, string(rs[2]), "com")
}

func TestJoin(t *testing.T) {
	parts := make([]string, 26)
	for i := 'A'; i <= 'Z'; i++ {
		parts[i-'A'] = string(i)
	}

	join := Join(",", parts...)
	assert.Equal(t, string(join), "A,B,C,D,E,F,G,H,I,J,K,L,M,N,O,P,Q,R,S,T,U,V,W,X,Y,Z")
}

func TestStringBuilder(t *testing.T) {
	builder := new(StringBuilder)

	builder.Append("Hello")
	builder.Append(" World")

	assert.Equal(t, builder.Size(), 11)
	assert.Equal(t, builder.ToString(), "Hello World")

	builder.AppendInt(123)
	assert.Equal(t, builder.Size(), 14)
	assert.Equal(t, builder.ToString(), "Hello World123")

	builder.AppendFloat(0.1234567)
	assert.Equal(t, builder.Size(), 23)
	assert.Equal(t, builder.ToString(), "Hello World1230.1234567")

	builder.Clear()
	assert.Equal(t, builder.Size(), 0)
	assert.Equal(t, builder.ToString(), "")
}
