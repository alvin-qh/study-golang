package time

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

var (
	ZONE_LOCAL, _ = time.LoadLocation("Asia/Shanghai") // 获取东八区时区

	// Go 语言中格式化时间模板不是常见的 Y,M,S 等
	// 而是 Go 语言的诞生时间 2006-01-02 15:04:05.000 MST
	TIME_LAYOUT_UTC = "2006-01-02T15:04:05.000Z" // 以 Z 结尾的标准 UTC 时间格式
	TIME_UTC        = "2021-11-11T12:00:00.100Z"

	TIME_LAYOUT_OFF = "2006-01-02T15:04:05.000-07:00" // 带时区偏移量的时间格式
	TIME_OFF        = "2021-11-11T12:00:00.100+08:00"

	TIME_LAYOUT_ST = "2006-01-02 15:04:05 MST" // 带时间标准符合的时间格式

	TIME_LAYOUT_LOCAL = "2006-01-02 15:04:05.0000" // 设置本地时间格式
	TIME_LOCAL        = "2021-11-11 12:00:00.0000"
)

// 测试创建时间对象
//
// 通过 time.Date 创建时间对象
func TestTime_Date(t *testing.T) {
	// 创建 UTC 时区时间对象
	tm := time.Date(2012, 11, 11, 12, 0, 0, 0, time.UTC)
	assert.Equal(t, "2012-11-11 12:00:00 +0000 UTC", tm.String())

	// 创建东八区时区时间对象
	tm = time.Date(2012, 11, 11, 12, 0, 0, 0, ZONE_LOCAL)
	assert.Equal(t, "2012-11-11 12:00:00 +0800 CST", tm.String())
}

// 测试时区转换
//
// 时间对象的 Location 函数获取该时间的时区对象; In 函数可以在不同时区转换时间; UTC 函数将时间转换为 UTC 时区
func TestTime_Location(t *testing.T) {
	// UTC 时区时间
	tm := time.Date(2012, 11, 11, 12, 0, 0, 0, time.UTC)

	// 获取对应时区
	loc := tm.Location()
	assert.Equal(t, "UTC", loc.String())

	// 时区转换到东八区
	tm = tm.In(ZONE_LOCAL)
	assert.Equal(t, "2012-11-11 20:00:00 +0800 CST", tm.String())

	// 获取时区
	loc = tm.Location()
	assert.Equal(t, "Asia/Shanghai", loc.String())

	// 时区再次转换到 UTC
	tm = tm.UTC()
	assert.Equal(t, "2012-11-11 12:00:00 +0000 UTC", tm.String())
}

// 测试一个时间和时间差值相加
//
// 时间对象的 Add 函数用于求时间和一个 Duration 对象的结果
func TestTime_Add(t *testing.T) {
	tm := time.Date(2012, 11, 11, 12, 0, 0, 0, time.UTC)

	// 给一个时间加上时间差值
	tmAdd := tm.Add(8 * 24 * time.Hour)
	assert.Equal(t, time.Date(2012, 11, 19, 12, 0, 0, 0, time.UTC), tmAdd)

	// 给一个时间加上负的时间差值
	tmAdd = tm.Add(-8 * 24 * time.Hour)
	assert.Equal(t, time.Date(2012, 11, 3, 12, 0, 0, 0, time.UTC), tmAdd)
}

// 测试两个时间的时间差
//
// 时间对象的 Sub 函数用于求两个时间对象的差, 结果为 Duration 对象
func TestTime_Sub(t *testing.T) {
	tm1 := time.Date(2012, 11, 11, 12, 0, 0, 0, time.UTC)
	tm2 := time.Date(2012, 11, 11, 20, 0, 0, 0, time.UTC)

	d := tm2.Sub(tm1) // 求两个时间的时间差
	assert.Equal(t, float64(8), d.Hours())

	d = tm1.Sub(tm2) // 时间差可以为负数
	assert.Equal(t, float64(-8), d.Hours())
}

// 测试从字符串解析时间差值
func TestTime_ParseDuration(t *testing.T) {
	// 可以用一个字符串解析为时间差对象
	// 时间差为 1 小时 30 分 0 秒
	d, err := time.ParseDuration("1h30m")
	assert.Nil(t, err)
	assert.Equal(t, float64(1.5), d.Hours())
	assert.Equal(t, "1h30m0s", d.String())
}

func TestDuration_Round(t *testing.T) {
	// 定义时间差 2 小时 30 分
	d := time.Hour*2 + time.Minute*30

	// 将时间差按小时取整, 结果为 3 小时, 超过半小时向上取整, 反之向下取整
	d = d.Round(time.Hour)
	assert.Equal(t, "3h0m0s", d.String())

	// 求 3 小时时间差和 35 分钟的倍数, 结果为 2 小时 55 分钟
	d = d.Round(time.Minute * 35)
	assert.Equal(t, "2h55m0s", d.String())
}

// 对时间进行 Round 操作
//
// Round 函数将时间计算为指定时间差的整数倍时间
func TestTime_Round(t *testing.T) {
	tm := time.Date(2012, 11, 11, 11, 40, 0, 0, time.UTC)
	d := time.Hour * 6

	// 求能够和指定时间差形成整数倍的时间值
	tm = tm.Round(d)
	assert.Equal(t, "2012-11-11 12:00:00 +0000 UTC", tm.String())

	// 计算整数倍数为 62622 倍
	mul := float64(tm.UnixNano()) / float64(d)
	assert.Equal(t, 62622.0, float64(mul))
}

// 对时间进行 Truncate 操作
//
// Truncate 函数将时间计算为最接近的整点时间
func TestTime_Truncate(t *testing.T) {
	tm := time.Date(2012, 11, 11, 11, 40, 0, 0, time.UTC)
	d := time.Hour * 6

	tm = tm.Truncate(d)
	assert.Equal(t, "2012-11-11 06:00:00 +0000 UTC", tm.String())
}

// 时间格式化为字符串
func TestTimeFormat(t *testing.T) {
	tm := time.Date(2012, 11, 11, 12, 0, 0, 0, ZONE_LOCAL)

	s := tm.Format(TIME_LAYOUT_UTC)
	assert.Equal(t, "2012-11-11T12:00:00.000Z", s) // 格式化为标准 UTC 格式

	s = tm.Format(TIME_LAYOUT_OFF)
	assert.Equal(t, "2012-11-11T12:00:00.000+08:00", s) // 格式化为带时区偏移量的格式

	s = tm.Format(TIME_LAYOUT_ST)
	assert.Equal(t, "2012-11-11 12:00:00 CST", s) // 格式化为带标准时间标识的格式
}

// 将字符串转化为时间对象
func TestParseTimeInUTC(t *testing.T) {
	checkResult := func(tm *time.Time) {
		assert.Equal(t, 2021, int(tm.Year()))
		assert.Equal(t, 11, int(tm.Month()))
		assert.Equal(t, 11, int(tm.Day()))
		assert.Equal(t, 12, int(tm.Hour()))
		assert.Equal(t, 0, int(tm.Minute()))
		assert.Equal(t, 0, int(tm.Second()))
		assert.Equal(t, 100000000, int(tm.Nanosecond()))
	}

	tm, err := time.Parse(TIME_LAYOUT_UTC, TIME_UTC) // 将一个标准 UTC 格式的时间转换为时间对象
	assert.Nil(t, err)
	checkResult(&tm)

	tm, err = time.Parse(TIME_LAYOUT_OFF, TIME_OFF) // 将一个带时区偏移量的字符串转为时间对象c
	assert.Nil(t, err)
	checkResult(&tm)
}

// 将时间进行序列化
func TestTimeMarshal(t *testing.T) {
	tm1 := time.Date(2012, 11, 11, 12, 0, 0, 0, time.UTC)
	b, err := tm1.MarshalBinary() // 序列化为 []byte
	assert.Nil(t, err)

	var tm2 time.Time
	err = tm2.UnmarshalBinary(b) // 从 []byte 反序列化
	assert.Nil(t, err)

	assert.Equal(t, tm1, tm2)

	b, err = tm1.MarshalJSON() // 序列化为 JSON 可用的格式
	assert.Nil(t, err)
	assert.Equal(t, "\"2012-11-11T12:00:00Z\"", string(b))

	tm2 = time.Time{}
	tm2.UnmarshalJSON([]byte("\"2012-11-11T12:00:00Z\"")) // 从 JSON 可用格式反序列化
	assert.Equal(t, tm1, tm2)

	b, err = tm1.MarshalText() // 序列化为字符串
	assert.Nil(t, err)
	assert.Equal(t, "2012-11-11T12:00:00Z", string(b))

	tm2 = time.Time{}
	tm2.UnmarshalText([]byte("2012-11-11T12:00:00Z")) // 从字符串反序列化
	assert.Equal(t, tm1, tm2)
}

// 测试通过指定的时区对象解析时间字符串
func TestParseTimeWithGivenLocationObject(t *testing.T) {
	// 通过 Asia/Shanghai 作为时区对象解析时间字符串
	tmLoc, err := time.ParseInLocation(TIME_LAYOUT_LOCAL, TIME_LOCAL, ZONE_LOCAL)
	assert.Nil(t, err)

	// 解析 UTC 时间
	tmUtc, err := time.Parse(TIME_LAYOUT_UTC, TIME_UTC)
	assert.Nil(t, err)

	// 此时两个时间的小时数一致
	assert.Equal(t, 0, tmLoc.Hour()-tmUtc.Hour())

	// 将 Asia/Shanghai 时间转为 UTC 时间
	tmLoc = tmLoc.In(time.UTC)

	// 此时两个时间的小时数相差 8 小时
	assert.Equal(t, 8, tmUtc.Hour()-tmLoc.Hour())
}
