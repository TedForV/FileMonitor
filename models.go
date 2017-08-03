package FileMonitor

type BaseCompareItemInfo struct {
	RelativePath string
	FileName     string
	Extention    string
}

type CompareItemInfo struct {
	BaseCompareItemInfo
	IsDir        bool
	IsMatched    bool
	IsAdditional bool
	IsMissing    bool
	Message      string
	IsCompared   bool
}
