package FileMonitor

//
//import (
//	"fmt"
//	"os"
//	"path/filepath"
//	"strings"
//)
//
//func main() {
//	//path, err := filepath.Abs("/test/a.log")
//	path := filepath.ToSlash("f:/temp/a.log")
//	fmt.Println(path)
//	//fmt.Println(temp)
//	//if err == nil {
//	//	fmt.Println("ok")
//	//	fmt.Println(path)
//	//} else {
//	//	fmt.Println("wrong")
//	//	fmt.Println(err)
//	//}
//	fileInfos := make(map[string]BaseCompareItemInfo)
//	basePath := "f:\\etlproject\\"
//	filepath.Walk("f:\\etlproject\\", func(path string, info os.FileInfo, err error) (returnErr error) {
//		if path == basePath {
//			return nil
//		}
//		if info.IsDir() {
//			Scan(path, info, err, &fileInfos, basePath, path)
//		} else {
//			relativePath := strings.Replace(path, basePath, "", -1)
//			fileInfos[relativePath] = BaseCompareItemInfo{
//				RelativePath: relativePath,
//				Extention:    filepath.Ext(path),
//				FileName:     strings.Replace(filepath.Base(path), filepath.Ext(path), "", -1),
//			}
//		}
//		return nil
//	})
//	for _, item := range fileInfos {
//		fmt.Println(item)
//	}
//}
