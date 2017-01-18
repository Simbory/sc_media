package handlers

import (
	"path"
	"strings"
)

func getDir(filePath string) string {
	return path.Dir(strings.Replace(filePath, "\\", "/", -1))
}
