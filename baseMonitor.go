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
}

func (bMonitor *BaseMonitor) Run() (
	infos []CompareItemInfo,
	errs []error, result bool) {
	result = true
	errs = make([]error, 0)
	originalFileInfos, err1 := bMonitor.scanOriginalFiles()
	monitorFileInfos, err2 := bMonitor.scanMonitorFiles()
	if err1 != nil {
		errs = append(errs, err1...)
	}
	if err2 != nil {
		errs = append(errs, err2...)
	}
	infos, err3 := bMonitor.compareFiles(originalFileInfos, monitorFileInfos)
	if err3 != nil {
		errs = append(errs, err3...)
	}
	if len(infos) > 0 {
		result = false
	}
	return infos, errs, result
}

func (bMonitor *BaseMonitor) GetCacheFilePath() string {
	return bMonitor.cacheFilePath
}

func (bMonitor *BaseMonitor) initialBaseMonitor() {
	bMonitor.compareResult = make([]CompareItemInfo, 50)
	bMonitor.errors = make([]error, 10)
}

func (bMonitor *BaseMonitor) scanOriginalFiles() (map[string]fileToolkit.BaseCompareItemInfo, []error) {
	return fileToolkit.RecursiveScanFiles(bMonitor.OriginalFolderPath, &bMonitor.UnMonitoredSubFolderPath)
}

func (bMonitor *BaseMonitor) scanMonitorFiles() (map[string]fileToolkit.BaseCompareItemInfo, []error) {
	return fileToolkit.RecursiveScanFiles(bMonitor.MonitorFolderPath, &bMonitor.UnMonitoredSubFolderPath)
}

func (bMonitor *BaseMonitor) compareFiles(
	originalInfos map[string]fileToolkit.BaseCompareItemInfo,
	monitorInfos map[string]fileToolkit.BaseCompareItemInfo) ([]CompareItemInfo, []error) {
	bMonitor.compareOriginalWithMonitor(originalInfos, monitorInfos)
	bMonitor.compareMonitorWithOriginal(originalInfos, monitorInfos)
	return bMonitor.compareResult, bMonitor.errors
}

func (bMonitor *BaseMonitor) compareOriginalWithMonitor(
	originalInfos map[string]fileToolkit.BaseCompareItemInfo,
	monitorInfos map[string]fileToolkit.BaseCompareItemInfo) {
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

func (bMonitor *BaseMonitor) compareMonitorWithOriginal(
	originalInfos map[string]fileToolkit.BaseCompareItemInfo,
	monitorInfos map[string]fileToolkit.BaseCompareItemInfo) {
	for _, value := range monitorInfos {
		if value.IsCompared {
			continue
		}
		bMonitor.addAdditionalRecord(&value)

	}
}
func (bMonitor *BaseMonitor) addAdditionalRecord(
	monitorFileInfo *fileToolkit.BaseCompareItemInfo) {
	bMonitor.compareResult = append(
		bMonitor.compareResult,
		createCompareItemInfo(
			monitorFileInfo,
			false,
			true,
			true))
}

func (bMonitor *BaseMonitor) addDifferRecord(
	originalFileInfo *fileToolkit.BaseCompareItemInfo) {
	bMonitor.compareResult = append(
		bMonitor.compareResult,
		createCompareItemInfo(
			originalFileInfo,
			false,
			false,
			true))
}

func (bMonitor *BaseMonitor) addMissingRecord(originalFileInfo *fileToolkit.BaseCompareItemInfo) {
	bMonitor.compareResult = append(
		bMonitor.compareResult,
		createCompareItemInfo(
			originalFileInfo,
			true,
			false,
			false))
}

func createCompareItemInfo(
	baseInfo *fileToolkit.BaseCompareItemInfo,
	isMissing bool,
	isAdditional bool,
	isNotMatched bool) CompareItemInfo {
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
