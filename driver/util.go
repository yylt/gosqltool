package driver

import (
	"fmt"
	"strings"
)

func addMark(dsn string) string {
	if strings.HasSuffix(dsn, "?") {
		return dsn
	} else {
		return fmt.Sprintf("%s?", dsn)
	}
}
