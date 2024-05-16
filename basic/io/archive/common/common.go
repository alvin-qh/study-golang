package common

import (
	"io"
	"os"
	"path/filepath"
	"reflect"
)

// 确认归档前的文件和归档恢复后的文件数量和内容一致
func CheckUnarchiveFiles(unarchivePath string, fileList []string) (bool, error) {
	for _, srcFile := range fileList {
		distFile := filepath.Join(unarchivePath, srcFile)
		eq, err := CompareTwoFiles(srcFile, distFile)
		if !eq || err != nil {
			return eq, err
		}
	}
	return true, nil
}

// 比较两个文件内容是否一致
func CompareTwoFiles(fa, fb string) (bool, error) {
	fileA, err := os.Open(fa)
	if err != nil {
		return false, err
	}
	defer fileA.Close()

	fileB, err := os.Open(fb)
	if err != nil {
		return false, err
	}
	defer fileB.Close()

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

// 创建目录
//
// 如果目录不存在, 则创建该目录
func CreateDirIfNotExists(dir string) error {
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		if err = os.MkdirAll(dir, 0755); err != nil {
			return err
		}
	} else if err != nil {
		return err
	}
	return nil
}
