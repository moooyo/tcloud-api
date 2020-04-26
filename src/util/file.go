package util

import (
	"os"
	"path/filepath"
)

func RemoveFileOrDirectoryWithoutError(name string) {
	info, err := os.Stat(name)
	if err != nil {
		ERROR("%s", err)
		return
	}
	if info.IsDir() {
		err = os.RemoveAll(name)
	} else {
		err = os.Remove(name)
	}
	if err != nil {
		ERROR("%s", err)
	}
}
func CreateDirectoryWhenNotExist(path string) {
	_, err := os.Stat(path)
	if err != nil {
		err = os.Mkdir(path, 0777)
		if err != nil {
			ERROR("%s", err.Error())
		}
	}
}

func WorkSpaceInit() {
	basePath := GetConfig().File.Path
	template := filepath.Join(basePath, "tmp")
	download := filepath.Join(template, "download")
	CreateDirectoryWhenNotExist(template)
	CreateDirectoryWhenNotExist(download)

}

func WorkSpaceClean() {
	basePath := GetConfig().File.Path
	template := filepath.Join(basePath, "tmp")
	download := filepath.Join(template, "download")
	RemoveFileOrDirectoryWithoutError(template)
	RemoveFileOrDirectoryWithoutError(download)
}
