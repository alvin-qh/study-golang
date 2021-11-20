package archive

import (
	"io"
	"os"
	"path/filepath"
	"reflect"
)

var (
	FileList = []string{ // 待归档文件列表
		"common.go",
		"gzip_test.go",
		"gzip.go",
		"tar_test.go",
		"tar.go",
		"zip_test.go",
		"zip.go",
	}
)

// 确认归档前的文件和归档恢复后的文件数量和内容一致
func CheckUnarchiveFiles(unarchivePath string) (bool, error) {
	for _, srcFile := range FileList {
		distFile := filepath.Join(unarchivePath, srcFile)
		eq, err := compareTwoFiles(srcFile, distFile)
		if !eq || err != nil {
			return eq, err
		}
	}
	return true, nil
}

func compareTwoFiles(fa, fb string) (bool, error) {
	fileA, err := os.Open(fa)
	if err != nil {
		return false, err
	}

	fileB, err := os.Open(fb)
	if err != nil {
		return false, err
	}

	dataA, err := io.ReadAll(fileA)
	if err != nil {
		return false, err
	}

	dataB, err := io.ReadAll(fileB)
	if err != nil {
		return false, err
	}

	return reflect.DeepEqual(dataA, dataB), nil
}

func createDirIfNotExists(dir string) error {
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		if err = os.MkdirAll(dir, 0755); err != nil {
			return err
		}
	} else if err != nil {
		return err
	}
	return nil
}
