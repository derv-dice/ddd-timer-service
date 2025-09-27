package calendar_drawer

import (
	"os"
	"path/filepath"
	"time"
)

const tmpDirPath = "./tmp"

var (
	testDate1 = time.Date(2024, time.November, 11, 0, 0, 0, 0, time.Local)
	testDate2 = time.Date(2025, time.November, 10, 0, 0, 0, 0, time.Local)
)

func createTmpDir() {
	_ = os.Mkdir(tmpDirPath, 0777)
}

func removeTmpDir() {
	//_ = os.RemoveAll(tmpDirPath)
}

func pathInTmpDir(path string) string {
	return filepath.Join(tmpDirPath, path)
}
