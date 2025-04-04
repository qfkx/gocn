package util

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

// GetParentDirectory 获取父文件夹
func GetParentDirectory(dirctory string) string {
	return SubStr(dirctory, 0, strings.LastIndex(dirctory, "/"))
}

// GetCurrentDirectory 获取当前文件夹：windows
func GetCurrentDirectory() string {
	dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		log.Fatal(err)
	}
	return strings.Replace(dir, "\\", "/", -1)
}

// JoinX 获取文件夹，如果不存在就创建。创建后合并指定的文件。
func JoinX(folderName string, filename string) string {
	file, _ := exec.LookPath(os.Args[0])
	path, _ := filepath.Abs(file)
	basePath := filepath.Dir(path)
	folderPath := filepath.Join(basePath, folderName)
	return JoinX(folderPath, filename)
}

// JoinAbsX 获取文件夹，如果不存在就创建。创建后合并指定的文件。
func JoinAbsX(folderPath string, filename string) string {
	if _, err := os.Stat(folderPath); os.IsNotExist(err) {
		os.Mkdir(folderPath, 0777)
		os.Chmod(folderPath, 0777)
	}
	return filepath.Join(folderPath, filename)
}

// IsFileExist1 判断文件是否存在，如果存在就删除
func IsFileExist1(filename string) bool {
	_, err := os.Stat(filename)
	if err != nil {
		if os.IsNotExist(err) {
			return false
		}
	}
	del := os.Remove(filename)
	if del != nil {
		return true
	}
	return false
}

// IsFileExist2 严格版本：判断文件是否存在，如果存在 & 大小不同，删除。
func IsFileExist2(filename string, filesize int64) bool {
	info, err := os.Stat(filename)
	if os.IsNotExist(err) {
		return false
	}
	if info.Size() == filesize {
		//fmt.Println("文件已存在", info.Name())
		return true
	}
	del := os.Remove(filename)
	if del != nil {
		fmt.Println(del)
	}
	return false
}

// IsFileExist3 版本：如果存在 & < limitSize，删除。
func IsFileExist3(filename string, limitSize int64) bool {
	info, err := os.Stat(filename)
	if os.IsNotExist(err) {
		return false
	}
	if info.Size() > limitSize {
		return true
	}
	del := os.Remove(filename)
	if del != nil {
		fmt.Println(del)
	}
	return false
}
