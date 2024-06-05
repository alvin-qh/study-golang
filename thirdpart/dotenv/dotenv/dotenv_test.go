package dotenv

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDotEnv_Load(t *testing.T) {
	keyA := os.Getenv("KEY_A")
	assert.Equal(t, "Test Key A", keyA)

	keyB := os.Getenv("KEY_B")
	assert.Equal(t, "Test Key B", keyB)
}
