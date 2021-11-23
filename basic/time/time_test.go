package time

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

var (
	TIME_LOCAL, _ = time.LoadLocation("Asia/Shanghai")

	// Go 语言中格式化时间模板不是常见的 Y,M,S 等
	// 而是 Go 语言的诞生时间 2006-01-02 15:04:05.000
	TIME_LAYOUT_UTC = "2006-01-02T15:04:05.000Z"
	TIME_UTC        = "2021-11-11T12:00:00.100Z"

	TIME_LAYOUT_OFF = "2006-01-02T15:04:05.000-07:00"
	TIME_OFF        = "2021-11-11T12:00:00.100+08:00Z"

	TIME_LAYOUT_ST = "2006-01-02 15:04:05 MST"
)

// 创建时间对象
func TestCreateTime(t *testing.T) {
	// 创建 UTC 时区时间对象
	tm := time.Date(2012, 11, 11, 12, 0, 0, 0, time.UTC)
	assert.Equal(t, "2012-11-11 12:00:00 +0000 UTC", tm.String())

	// 创建东八区时区时间对象
	tm = time.Date(2012, 11, 11, 12, 0, 0, 0, TIME_LOCAL)
	assert.Equal(t, "2012-11-11 12:00:00 +0800 CST", tm.String())
}

// 时区转换
func TestConvertLocalTime(t *testing.T) {
	tm := time.Date(2012, 11, 11, 12, 0, 0, 0, time.UTC) // UTC 时区时间

	loc := tm.Location() // 获取对应时区
	assert.Equal(t, "UTC", loc.String())

	tm = tm.In(TIME_LOCAL) // 时区转换到东八区
	assert.Equal(t, "2012-11-11 20:00:00 +0800 CST", tm.String())

	loc = tm.Location() // 获取时区
	assert.Equal(t, "Asia/Shanghai", loc.String())

	tm = tm.UTC() // 时区再次转换到 UTC
	assert.Equal(t, "2012-11-11 12:00:00 +0000 UTC", tm.String())
}

// 时间计算
func TestDuration(t *testing.T) {
	tm1 := time.Date(2012, 11, 11, 12, 0, 0, 0, time.UTC)
	tm2 := time.Date(2012, 11, 11, 20, 0, 0, 0, time.UTC)

	d := tm2.Sub(tm1)
	assert.Equal(t, float64(8), d.Hours())

	tm3 := tm1.Add(d)
	assert.Equal(t, tm2, tm3)

	d = tm1.Sub(tm2)
	assert.Equal(t, float64(-8), d.Hours())

	tm3 = tm2.Add(d)
	assert.Equal(t, tm3, tm1)

	d, err := time.ParseDuration("1h20m")
	assert.NoError(t, err)
	assert.Equal(t, int64(4800000000000), int64(d))
	assert.Equal(t, "1h20m0s", d.String())

	d = d.Round(time.Hour)
	assert.Equal(t, "1h0m0s", d.String())

	d = d.Round(time.Minute * 35)
	assert.Equal(t, "1h10m0s", d.String())

	mul := float64(d) / float64(time.Minute*35)
	assert.Equal(t, 2.0, mul)
}

func TestRoundAndTruncate(t *testing.T) {
	tm := time.Date(2012, 11, 11, 11, 40, 0, 0, time.UTC)
	d := time.Hour * 12

	tm = tm.Round(d)
	assert.Equal(t, "2012-11-11 12:00:00 +0000 UTC", tm.String())

	mul := float64(tm.UnixNano()) / float64(d)
	assert.Equal(t, 31311.0, float64(mul))

	tm = time.Date(2012, 11, 11, 11, 40, 0, 0, time.UTC)
	tm = tm.Truncate(d)
	assert.Equal(t, "2012-11-11 00:00:00 +0000 UTC", tm.String())
}

func TestTimeFormat(t *testing.T) {
	tm := time.Date(2012, 11, 11, 12, 0, 0, 0, TIME_LOCAL)

	s := tm.Format(TIME_LAYOUT_UTC)
	assert.Equal(t, "2012-11-11T12:00:00.000Z", s)

	s = tm.Format(TIME_LAYOUT_OFF)
	assert.Equal(t, "2012-11-11T12:00:00.000+08:00", s)

	s = tm.Format(TIME_LAYOUT_ST)
	assert.Equal(t, "2012-11-11 12:00:00 CST", s)
}

func TestParseTimeInUTC(t *testing.T) {
	tm, err := time.Parse(TIME_LAYOUT_UTC, TIME_UTC)
	assert.NoError(t, err)
	assert.Equal(t, "2021-11-11 12:00:00.1 +0000 UTC", tm.String())

	assert.Equal(t, 2021, int(tm.Year()))
	assert.Equal(t, 11, int(tm.Month()))
	assert.Equal(t, 11, int(tm.Day()))
	assert.Equal(t, 12, int(tm.Hour()))
	assert.Equal(t, 0, int(tm.Minute()))
	assert.Equal(t, 0, int(tm.Second()))
	assert.Equal(t, 100000000, int(tm.Nanosecond()))
}

func TestParseTimeInCurrentLocale(t *testing.T) {
	tm, err := time.Parse(TIME_LAYOUT_OFF, TIME_OFF)
	assert.NoError(t, err)
	assert.Equal(t, "2021-11-11 12:00:00.1 +0800 CST", tm.String())

	assert.Equal(t, 2021, int(tm.Year()))
	assert.Equal(t, 11, int(tm.Month()))
	assert.Equal(t, 11, int(tm.Day()))
	assert.Equal(t, 12, int(tm.Hour()))
	assert.Equal(t, 0, int(tm.Minute()))
	assert.Equal(t, 0, int(tm.Second()))
	assert.Equal(t, 100000000, int(tm.Nanosecond()))
}

func TestTimeMarshal(t *testing.T) {
	tm1 := time.Date(2012, 11, 11, 12, 0, 0, 0, time.UTC)
	b, err := tm1.MarshalBinary()
	assert.NoError(t, err)

	var tm2 time.Time
	err = tm2.UnmarshalBinary(b)
	assert.NoError(t, err)

	assert.Equal(t, tm1, tm2)

	b, err = tm1.MarshalJSON()
	assert.NoError(t, err)
	assert.Equal(t, "\"2012-11-11T12:00:00Z\"", string(b))

	tm2 = time.Time{}
	tm2.UnmarshalJSON([]byte("\"2012-11-11T12:00:00Z\""))
	assert.Equal(t, tm1, tm2)

    b, err = tm1.MarshalText()
    assert.NoError(t, err)
    assert.Equal(t, "2012-11-11T12:00:00Z", string(b))

    tm2 = time.Time{}
	tm2.UnmarshalText([]byte("2012-11-11T12:00:00Z"))
	assert.Equal(t, tm1, tm2)
}
