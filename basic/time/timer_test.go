package time

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestSleep(t *testing.T) {
	tm := time.Now()
	time.Sleep(3*time.Second + 2*time.Millisecond)

	d := time.Since(tm)
	assert.Equal(t, 3, int(d.Seconds()))
}

func TestTimeAfter(t *testing.T) {
	tm1 := time.Now()

	c := time.After(3 * time.Second)

	tm2 := <-c
	d := tm2.Sub(tm1)
    assert.Equal(t, 3, int(d.Seconds()))
}
