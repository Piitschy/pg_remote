package utils

import (
	"fmt"
	"strings"
)

func Ext(flag string) string {
	if flag == "p" {
		return "sql"
	}
	if flag == "t" {
		return "tar"
	}
	if flag == "c" {
		return ""
	}
	return ""
}

func FormatFlag(ext string) (string, error) {
	if strings.HasPrefix(ext, ".") {
		ext = ext[1:]
	}
	if ext == "" || ext == "c" || ext == "custom" {
		return "c", nil
	}
	if ext == "sql" || ext == "p" {
		return "p", nil
	}
	if ext == "tar" || ext == "t" {
		return "t", nil
	}
	return "p", fmt.Errorf("Invalid format %s", ext)
}
