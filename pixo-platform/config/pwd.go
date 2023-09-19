package config

import (
	"path/filepath"
	"runtime"
)

func GetProjectRoot() (root string) {
	_, b, _, _ := runtime.Caller(0)
	root = filepath.Join(filepath.Dir(b), "../..")
	return root
}
