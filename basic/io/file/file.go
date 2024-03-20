package file

import "os"

func CompareFileInfo(fi1 os.FileInfo, fi2 os.FileInfo) bool {
	return fi1.Name() == fi2.Name() &&
		fi1.Size() == fi2.Size() &&
		fi1.Mode() == fi2.Mode() &&
		fi1.ModTime() == fi2.ModTime() &&
		fi1.IsDir() == fi2.IsDir()
}
