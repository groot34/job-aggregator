package parsers

import (
	"fmt"
	"strings"
	"time"
)

func getIDFromURL(url string) string {
	url = strings.TrimSpace(url)
	if url == "" {
		return fmt.Sprint(time.Now().UnixNano())
	}

	url = strings.TrimSuffix(url, "/")
	url = strings.Split(url, "?")[0]
	url = strings.Split(url, "#")[0]
	parts := strings.FieldsFunc(url, func(r rune) bool {
		return r == '/' || r == '-' || r == '_' || r == '='
	})
	if len(parts) == 0 {
		return fmt.Sprint(time.Now().UnixNano())
	}
	return parts[len(parts)-1]
}
