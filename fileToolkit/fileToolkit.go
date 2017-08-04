package fileToolkit

import (
	"crypto/md5"
	"encoding/hex"
	"errors"
	"hash"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"sync"
)

type BaseCompareItemInfo struct {
	AbsPath      string
	RelativePath string
	FileName     string
	Extention    string
	IsCompared   bool
}

var md5ObjPool *sync.Pool

func init() {
	md5ObjPool = &sync.Pool{
		New: func() interface{} {
			return md5.New()
		},
	}
}

//Scan the all files's info in the folderPath and all sub folders recursively
func RecursiveScanFiles(folderPath string) (map[string]BaseCompareItemInfo, error) {
	if isDir := IsExistedDir(folderPath); !isDir {
		return nil, errors.New("the folderPath(" + folderPath + ") is not a correct path.")
	}
	fileInfos := make(map[string]BaseCompareItemInfo)
	filepath.Walk(folderPath, func(path string, info os.FileInfo, err error) (returnErr error) {
		if path == folderPath {
			return nil
		}
		if info.IsDir() {
			scan(path, info, err, &fileInfos, folderPath)
		} else {
			relativePath := strings.Replace(path, folderPath, "", -1)
			fileInfos[relativePath] = BaseCompareItemInfo{
				AbsPath:      path,
				RelativePath: relativePath,
				Extention:    filepath.Ext(path),
				FileName:     strings.Replace(filepath.Base(path), filepath.Ext(path), "", -1),
			}
		}
		return nil
	})
	return fileInfos, nil
}

//Compare the two files by calculating the MD5 string
func HasTheSameContent(filePath1 string, filePath2 string) (bool, error) {
	var (
		wg           sync.WaitGroup
		err1, err2   error
		data1, data2 string
	)
	wg.Add(2)
	go func() {
		defer wg.Done()
		data1, err1 = getFileMD5(filePath1)

	}()
	go func() {
		defer wg.Done()
		data2, err2 = getFileMD5(filePath2)
	}()

	wg.Wait()
	if err1 == nil && err2 == nil {
		return data1 == data2, nil
	}
	msg := ""
	if err1 != nil {
		msg = err1.Error() + "(" + filePath1 + ");"
	}
	if err2 != nil {
		msg = err2.Error() + "(" + filePath2 + ")"
	}
	return false, errors.New(msg)

}

//Check the path is directory or not
func IsExistedDir(path string) bool {
	path = strings.TrimSpace(path)
	if path == "" {
		return false
	}
	fileInfo, err := os.Stat(path)
	if err != nil {
		return false
	}
	return fileInfo.IsDir()
}

func getFileMD5(filePath string) (string, error) {
	data, err := ioutil.ReadFile(filePath)
	if err != nil {
		return "", err
	}
	return getMD5(data), nil
}

func getMD5(data []byte) string {
	md5Obj := md5ObjPool.Get().(hash.Hash)
	defer md5ObjPool.Put(md5Obj)

	md5Obj.Reset()
	md5Obj.Write(data)
	return hex.EncodeToString(md5Obj.Sum(nil))
}

func scan(currentPath string, info os.FileInfo, err error, fileInfos *map[string]BaseCompareItemInfo, basePath string) error {
	if info.IsDir() {
		filepath.Walk(currentPath, func(path string, info os.FileInfo, err error) (returnErr error) {
			if path == currentPath {
				return nil
			}
			return scan(path, info, err, fileInfos, basePath)
		})
	} else {
		relativePath := strings.Replace(currentPath, basePath, "", -1)
		(*fileInfos)[relativePath] = BaseCompareItemInfo{
			RelativePath: relativePath,
			Extention:    filepath.Ext(currentPath),
			FileName:     strings.Replace(filepath.Base(currentPath), filepath.Ext(currentPath), "", -1),
			IsCompared:   false,
		}
	}
	return nil
}
