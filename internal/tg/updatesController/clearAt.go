package updatesController

import "strings"

func clearAt(s string) string {
	return strings.TrimPrefix(s, "@")
}
