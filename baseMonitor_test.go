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
