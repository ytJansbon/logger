package logger

import (
	"os"
	"path/filepath"
	"time"
)

//get current date
func getCurrentDate() *time.Time {
	tm, _ := time.Parse("2006-01-02", time.Now().Format("2006-01-02"))
	return &tm
}

//check if a directory or path is exist
func isPathExist(path string) bool {
	_, err := os.Stat(path)
	if err == nil {
		return true
	}
	if os.IsExist(err) {
		return true
	}
	return false
}

//join directory and file name into a path
func joinFilePath(dir string, file string) string {
	return filepath.Join(dir, file)
}

func getFileSize(path string) uint64 {
	info, err := os.Stat(path)
	if err != nil {
		return 0
	}
	return uint64(info.Size())
}
