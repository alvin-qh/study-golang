package fs

import (
	"fmt"
	"io/fs"
	"testing"

	"github.com/stretchr/testify/assert"
)

func mapDirEntriesToStrings(entries []fs.DirEntry) []string {
	var strings []string
	for _, entry := range entries {
		if entry.IsDir() {
			strings = append(strings, fmt.Sprintf("%v <dir>", entry.Name()))
		} else {
			strings = append(strings, fmt.Sprintf("%v <file>", entry.Name()))
		}

	}
	return strings
}

func TestReadEmbedDir(t *testing.T) {
	entires, err := STATIC_ASSETS.ReadDir(".")
	assert.Nil(t, err)
	assert.ElementsMatch(t, mapDirEntriesToStrings(entires), []string{"asset <dir>"})

	entires, err = STATIC_ASSETS.ReadDir("asset")
	assert.Nil(t, err)
	assert.ElementsMatch(t, mapDirEntriesToStrings(entires), []string{
		"01 <dir>",
		"02 <dir>",
		"static.txt <file>",
	})

	entires, err = STATIC_ASSETS.ReadDir("asset/01")
	assert.Nil(t, err)
	assert.ElementsMatch(t, mapDirEntriesToStrings(entires), []string{"static1.txt <file>"})

	entires, err = STATIC_ASSETS.ReadDir("asset/02")
	assert.Nil(t, err)
	assert.ElementsMatch(t, mapDirEntriesToStrings(entires), []string{"static2.txt <file>"})
}
