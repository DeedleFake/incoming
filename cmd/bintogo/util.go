package main

import (
	"path/filepath"
	"strings"
)

func TrimExt(path string) string {
	return strings.TrimSuffix(path, filepath.Ext(path))
}
