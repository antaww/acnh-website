package main

import "strings"

func TrimURLPrefix(url string) string {
	return strings.TrimPrefix(url, "/")
}
