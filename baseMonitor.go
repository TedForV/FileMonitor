package FileMonitor

import (
	"crypto/md5"
	"encoding/hex"
	"hash"
	"os"
	"path/filepath"
	"strings"
)

type Monitor interface {
	LoadSetting(jsonFilePath string) error
	Run() []CompareItemInfo
	ErrorProcedure(err error, additionInfo string)
	Alert(msg string)
}

type BaseMonitor struct {
	OriginalFolderPath string `json:"originalfolderpath"`
	MonitorFolderPath  string `json:"monitorfolderpath"`
	MonitorPeriod      int    `json:"monitorperiod"`
}

var md5Obj hash.Hash

func Init() {
	md5Obj = md5.New()
}

func Scan(path string, info os.FileInfo, err error, fileInfos *map[string]BaseCompareItemInfo, basePath string, currentPath string) error {
	if path == currentPath {
		return nil
	}
	if info.IsDir() {
		filepath.Walk(path, func(path string, info os.FileInfo, err error) (returnErr error) {
			return Scan(path, info, err, fileInfos, basePath, path)
		})
	} else {
		relativePath := strings.Replace(path, basePath, "", -1)
		(*fileInfos)[relativePath] = BaseCompareItemInfo{
			RelativePath: relativePath,
			Extention:    filepath.Ext(path),
			FileName:     strings.Replace(filepath.Base(path), filepath.Ext(path), "", -1),
		}
	}
	return nil
}

//
//func (bMonitor *BaseMonitor) FileComparer(path string, file os.FileInfo) (bool, string) {
//	if file.IsDir() {
//		return true, ""
//	}
//	f, err := os.Open(path)
//	if err != nil {
//		return false, err.Error()
//	}
//
//}

func getMD5(data []byte) string {
	md5Obj.Reset()
	md5Obj.Write(data)
	return hex.EncodeToString(md5Obj.Sum(nil))
}
