package porcelain

import (
	"os"
	"strings"
)

func forceSlashSeparators(name string) string {
	return strings.Replace(name, os.PathSeparator, "/", -1)
}
