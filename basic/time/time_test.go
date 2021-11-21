package time

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

const (
	// Go 语言中格式化时间模板不是常见的 Y,M,S 等
	// 而是 Go 语言的诞生时间 2006-01-02 15:04:05.000
	TIME_LAYOUT_UTC = "2006-01-02T15:04:05.000Z"
	TIME_UTC        = "2021-11-11T12:00:00.100Z"

	TIME_LAYOUT_OFF = "2006-01-02T15:04:05.000-07:00"
	TIME_OFF        = "2021-11-11T12:00:00.100+08:00"
)

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
