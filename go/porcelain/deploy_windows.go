package porcelain

import (
	"fmt"
	"os"
	"strings"
)

func forceSlashSeparators(name string) string {
	return strings.Replace(name, fmt.Sprintf("%c", os.PathSeparator), "/", -1)
}
