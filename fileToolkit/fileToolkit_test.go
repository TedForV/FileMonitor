package fileToolkit

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestRecursiveScanFiles(t *testing.T) {
	unWatched := []string{"F:\\ETLProject\\AccountHistory"}
	data, err := RecursiveScanFiles(
		"F:\\ETLProject\\", &unWatched)
	if err != nil {
		for _, item := range err {
			t.Error(item.Error())
		}
		return
	}
	assert.Equal(t, 17, len(data))
	for _, item := range data {
		t.Log(item)
	}
}

func TestHasTheSameContent_Same(t *testing.T) {
	data, err := HasTheSameContent("F:\\ETLProject\\AccountHistory\\AccountHistoryProcedure.kjb", "F:\\ETLProject\\AccountHistory\\AccountHistoryProcedure.kjb")
	if err != nil {
		t.Error(err.Error())
	}
	if !data {
		t.Error("it should be the same.")
	}
}

func TestHasTheSameContent_NotSame(t *testing.T) {
	data, err := HasTheSameContent("F:\\ETLProject\\AccountHistory\\AccountHistoryProcedure.kjb", "F:\\ETLProject\\AccountHistory\\InitialProcedure.ktr")
	if err != nil {
		t.Log(err.Error())
	}
	if data {
		t.Error("it should be not the same.")
	}
}

func TestHasTheSameContent_NotSame_withWrongPath(t *testing.T) {
	data, err := HasTheSameContent("F:\\ETLProject\\AccountHistory\\test.kjb", "F:\\ETLProject\\AccountHistory\\InitialProcedure.ktr")
	if err != nil {
		t.Log(err.Error())
	}
	if data {
		t.Error("it should be not the same.")
	}
}

func TestHasTheSameContent_NotSame_withNullPath(t *testing.T) {
	data, err := HasTheSameContent("", "F:\\ETLProject\\AccountHistory\\InitialProcedure.ktr")
	if err != nil {
		t.Log(err.Error())
	}
	if data {
		t.Error("it should be not the same.")
	}
}

func TestHasTheSameContent_NotSame_withBothNullPath(t *testing.T) {
	data, err := HasTheSameContent("", "")
	if err != nil {
		t.Log(err.Error())
	}
	if data {
		t.Error("it should be not the same.")
	}
}

func TestIsExistedDir_IsDir(t *testing.T) {
	if !IsExistedDir("F:\\ETLProject\\") {
		t.Error("it should be a dir")
	}
}

func TestIsExistedDir_IsWrongDir(t *testing.T) {
	if IsExistedDir("F:\\ETLProject1\\") {
		t.Error("it should not be a dir")
	}
}

func TestIsExistedDir_IsFile(t *testing.T) {
	if IsExistedDir("") {
		t.Error("it should be a nil")
	}
}
