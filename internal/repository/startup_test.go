package repository

import (
	"os"
	"path/filepath"
)

const tmpDirPath = "./tmp"

func createTmpDir() {
	_ = os.Mkdir(tmpDirPath, 0777)
}

func removeTmpDir() {
	// _= os.RemoveAll(tmpDirPath)
}

func pathInTmpDir(path string) string {
	return filepath.Join(tmpDirPath, path)
}
