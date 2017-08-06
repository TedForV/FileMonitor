package FileMonitor

import (
	"FileMonitor/fileToolkit"
)

type CompareItemInfo struct {
	fileToolkit.BaseCompareItemInfo
	IsNotMatched bool
	IsAdditional bool
	IsMissing    bool
	Message      string
}

type Monitor interface {
	Run() ([]CompareItemInfo, []error, bool)
}

type BaseMonitor struct {
	OriginalFolderPath       string   `json:"originalfolderpath"`
	MonitorFolderPath        string   `json:"monitorfolderpath"`
	UnMonitoredSubFolderPath []string `json:"unwatchedfolderpath"`
	cacheFilePath            string
	compareResult            []CompareItemInfo
	errors                   []error
	//MonitorPeriod      int    `json:"monitorperiod"`

}

func (bMonitor *BaseMonitor) Run() (infos []CompareItemInfo, errs []error, result bool) {
	result = true
	errs = make([]error, 10)
	originalFileInfos, err1 := bMonitor.scanOriginalFiles()
	monitorFileInfos, err2 := bMonitor.scanMonitorFiles()
	if err1 != nil {
		errs = append(errs, err1)
	}
	if err2 != nil {
		errs = append(errs, err2)
	}
	infos, err3 := bMonitor.compareFiles(originalFileInfos, monitorFileInfos)
	if err3 != nil {
		errs = append(errs, err3...)
	}
	if len(errs) > 0 {
		result = false
	}
}

func (bMonitor *BaseMonitor) GetCacheFilePath() string {
	return bMonitor.cacheFilePath
}

func (bMonitor *BaseMonitor) initialBaseMonitor() {
	bMonitor.compareResult = make([]CompareItemInfo, 50)
	bMonitor.errors = make([]error, 10)
}

func (bMonitor *BaseMonitor) scanOriginalFiles() (map[string]fileToolkit.BaseCompareItemInfo, error) {
	return fileToolkit.RecursiveScanFiles(bMonitor.OriginalFolderPath)
}

func (bMonitor *BaseMonitor) scanMonitorFiles() (map[string]fileToolkit.BaseCompareItemInfo, error) {
	return fileToolkit.RecursiveScanFiles(bMonitor.MonitorFolderPath)
}

func (bMonitor *BaseMonitor) compareFiles(originalInfos map[string]fileToolkit.BaseCompareItemInfo, monitorInfos map[string]fileToolkit.BaseCompareItemInfo) ([]CompareItemInfo, []error) {
	bMonitor.compareOriginalWithMonitor(originalInfos, monitorInfos)
	bMonitor.compareMonitorWithOriginal(originalInfos, monitorInfos)
	return bMonitor.compareResult, bMonitor.errors
}

func (bMonitor *BaseMonitor) compareOriginalWithMonitor(originalInfos map[string]fileToolkit.BaseCompareItemInfo, monitorInfos map[string]fileToolkit.BaseCompareItemInfo) {
	var (
		monitorItem fileToolkit.BaseCompareItemInfo
		ok          bool
	)
	for key, value := range originalInfos {
		monitorItem, ok = monitorInfos[key]
		if !ok {
			bMonitor.addMissingRecord(&value)
		} else {
			if result, err := fileToolkit.HasTheSameContent(value.AbsPath, monitorItem.AbsPath); !result {
				if err != nil {
					bMonitor.errors = append(bMonitor.errors, err)
				}
				bMonitor.addDifferRecord(&value)
			}
			monitorInfos[key] = fileToolkit.BaseCompareItemInfo{
				AbsPath:      monitorItem.AbsPath,
				RelativePath: monitorItem.RelativePath,
				FileName:     monitorItem.FileName,
				Extention:    monitorItem.Extention,
				IsCompared:   true,
			}
		}
	}
}

func (bMonitor *BaseMonitor) compareMonitorWithOriginal(originalInfos map[string]fileToolkit.BaseCompareItemInfo, monitorInfos map[string]fileToolkit.BaseCompareItemInfo) {
	for _, value := range monitorInfos {
		if value.IsCompared {
			continue
		}
		bMonitor.addAdditionalRecord(&value)

	}
}
func (bMonitor *BaseMonitor) addAdditionalRecord(monitorFileInfo *fileToolkit.BaseCompareItemInfo) {
	bMonitor.compareResult = append(bMonitor.compareResult, createCompareItemInfo(monitorFileInfo, false, true, true))
	//bMonitor.compareResult = append(bMonitor.compareResult, CompareItemInfo{
	//	BaseCompareItemInfo: fileToolkit.BaseCompareItemInfo{
	//		AbsPath:      "",
	//		RelativePath: monitorFileInfo.RelativePath,
	//		FileName:     monitorFileInfo.FileName,
	//		Extention:    monitorFileInfo.Extention,
	//		IsCompared:   true,
	//	},
	//	IsMissing:    false,
	//	IsAdditional: true,
	//	IsNotMatched: false,
	//	Message:      "",
	//})
}

func (bMonitor *BaseMonitor) addDifferRecord(originalFileInfo *fileToolkit.BaseCompareItemInfo) {
	bMonitor.compareResult = append(bMonitor.compareResult, createCompareItemInfo(originalFileInfo, false, false, true))
	//bMonitor.compareResult = append(bMonitor.compareResult, CompareItemInfo{
	//	BaseCompareItemInfo: fileToolkit.BaseCompareItemInfo{
	//		AbsPath:      "",
	//		RelativePath: originalFileInfo.RelativePath,
	//		FileName:     originalFileInfo.FileName,
	//		Extention:    originalFileInfo.Extention,
	//		IsCompared:   true,
	//	},
	//	IsMissing:    false,
	//	IsAdditional: false,
	//	IsNotMatched: true,
	//	Message:      "",
	//})
}

func (bMonitor *BaseMonitor) addMissingRecord(originalFileInfo *fileToolkit.BaseCompareItemInfo) {
	bMonitor.compareResult = append(bMonitor.compareResult, createCompareItemInfo(originalFileInfo, true, false, false))
	//bMonitor.compareResult = append(bMonitor.compareResult, CompareItemInfo{
	//	BaseCompareItemInfo: fileToolkit.BaseCompareItemInfo{
	//		AbsPath:      "",
	//		RelativePath: originalFileInfo.RelativePath,
	//		FileName:     originalFileInfo.FileName,
	//		Extention:    originalFileInfo.Extention,
	//		IsCompared:   true,
	//	},
	//	IsMissing:    true,
	//	IsAdditional: false,
	//	IsNotMatched: false,
	//	Message:      "",
	//})
}

func createCompareItemInfo(baseInfo *fileToolkit.BaseCompareItemInfo, isMissing bool, isAdditional bool, isNotMatched bool) CompareItemInfo {
	return CompareItemInfo{
		BaseCompareItemInfo: fileToolkit.BaseCompareItemInfo{
			AbsPath:      "",
			RelativePath: baseInfo.RelativePath,
			FileName:     baseInfo.FileName,
			Extention:    baseInfo.Extention,
			IsCompared:   true,
		},
		IsMissing:    isMissing,
		IsAdditional: isAdditional,
		IsNotMatched: isNotMatched,
		Message:      "",
	}
}
