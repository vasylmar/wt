package utils

import "strings"

func GetWorktreeName(s string) string {
	parts := strings.Split(s, "/")
	return parts[len(parts)-1]
}
