package FileMonitor

import (
	"encoding/json"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestJsonStringToModel(t *testing.T) {
	var bMonitor BaseMonitor
	str := `{"originalfolderpath":"originalfolderpath","monitorfolderpath":"monitorfolderpath","unwatchedfolderpath":["unwatchedfolderpath1","unwatchedfolderpath2","unwatchedfolderpath3"]}`
	if err := json.Unmarshal([]byte(str), &bMonitor); err != nil {
		t.Error(err.Error())
		return
	}
	assert := assert.New(t)
	assert.Equal(bMonitor.OriginalFolderPath, "originalfolderpath", "this should be the same.")
	assert.Equal(bMonitor.MonitorFolderPath, "monitorfolderpath", "this should be the same.")
	assert.Equal(len(bMonitor.UnMonitoredSubFolderPath), 3, "array is 3")
	assert.Equal(bMonitor.UnMonitoredSubFolderPath[0], "unwatchedfolderpath1", "this should be the same.")
}

func TestBaseMonitor_Run(t *testing.T) {
	bMonitor, err := createTestBaseMonitor(`{"originalfolderpath":"F:\\ETLProject\\","monitorfolderpath":"F:\\Temp\\","unwatchedfolderpath":[]}`)
	if err != nil {
		t.Error(err.Error())
		return
	}
	compareInfo, errs, result := bMonitor.Run()
	if errs != nil && len(errs) != 0 {
		t.Log("errors :")
		for _, item := range errs {
			t.Error(item.Error())
		}
	}
	if !result {
		for _, item := range compareInfo {
			t.Error(item)
		}
	}
}

func createTestBaseMonitor(jsonStr string) (BaseMonitor, error) {
	var bMonitor BaseMonitor
	err := json.Unmarshal([]byte(jsonStr), &bMonitor)
	return bMonitor, err

}
