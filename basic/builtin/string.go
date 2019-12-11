package builtin

import (
	"bytes"
	"strconv"
	"strings"
)

type String string
type Char rune

func FromInt(n int) String {
	//return String(strconv.Itoa(n))
	return String(strconv.FormatInt(int64(n), 10))
}

func FromFloat(f float32) String {
	return String(strconv.FormatFloat(float64(f), 'f', -1, 32))
}

func (i String) Len() int {
	return len([]rune(i))
}

func (i String) CharAt(index int) Char {
	return Char([]rune(i)[index])
}

func (i String) ToInt() (int, error) {
	//n, e := strconv.Atoi(string(i))
	n, e := strconv.ParseInt(string(i), 10, 32)
	return int(n), e
}

func (i String) ToFloat() (float32, error) {
	f, e := strconv.ParseFloat(string(i), 32)
	return float32(f), e
}

func (i String) Compare(other string) int {
	return strings.Compare(string(i), other)
}

func (i String) Contains(s string) bool {
	return strings.Contains(string(i), s)
}

func (i String) Index(s string) int {
	return strings.Index(string(i), s)
}

func (i String) LastIndex(s string) int {
	return strings.LastIndex(string(i), s)
}

func (i String) Count(s string) int {
	return strings.Count(string(i), s)
}

func Repeat(s string, count int) String {
	return String(strings.Repeat(s, count))
}

func (i String) Replace(old string, new string, count int) String {
	return String(strings.Replace(string(i), old, new, count))
}

func (i *String) ReplaceSelf(old string, new string, count int) {
	*i = String(strings.Replace(string(*i), old, new, count))
}

func (i String) Trim() String {
	return String(strings.TrimSpace(string(i)))
}

func (i String) TrimWith(charset string) String {
	s := strings.TrimLeft(string(i), charset)
	return String(strings.TrimRight(s, charset))
}

func (i String) StartWith(prefix string) bool {
	return strings.HasPrefix(string(i), prefix)
}

func (i String) TrimStart(prefix string) String {
	return String(strings.TrimPrefix(string(i), prefix))
}

func (i String) EndWith(suffix string) bool {
	return strings.HasSuffix(string(i), suffix)
}

func (i String) TrimEnd(suffix string) String {
	return String(strings.TrimSuffix(string(i), suffix))
}

func (i String) Split(sep string) []String {
	split := strings.Split(string(i), sep)
	rs := make([]String, len(split))
	for n, s := range strings.Split(string(i), sep) {
		rs[n] = String(s)
	}
	return rs
}

func (i String) S() string {
	return string(i)
}

func Join(sep string, part ...string) String {
	return String(strings.Join(part, sep))
}

type StringBuilder struct {
	buffer bytes.Buffer
}

func (i *StringBuilder) Append(s string) *StringBuilder {
	i.buffer.WriteString(s)
	return i
}

func (i *StringBuilder) AppendInt(n int) *StringBuilder {
	i.buffer.WriteString(FromInt(n).S())
	return i
}

func (i *StringBuilder) AppendFloat(n float32) *StringBuilder {
	i.buffer.WriteString(FromFloat(n).S())
	return i
}

func (i *StringBuilder) ToString() string {
	return i.buffer.String()
}

func (i *StringBuilder) Size() int {
	return i.buffer.Len()
}

func (i *StringBuilder) Clear() {
	i.buffer.Reset()
}
