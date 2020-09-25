package tools

import (
	"strings"
)

// GetExtension from filename
// Ex: bayu.jpg return jpg
func GetExtension(filename string) string {
	names := strings.Split(filename, ".")
	return names[len(names)-1]
}
