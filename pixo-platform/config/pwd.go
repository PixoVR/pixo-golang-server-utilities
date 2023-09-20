package config

import (
	"path/filepath"
	"runtime"
)

func GetProjectRoot(differential ...string) (root string) {
	_, b, _, _ := runtime.Caller(0)
	dir := filepath.Dir(b)

	if len(differential) > 0 {
		diff := differential[0]
		root = filepath.Join(dir, diff)
	}

	return root
}
