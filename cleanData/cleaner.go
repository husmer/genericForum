package cleanData

import "strings"

func CleanEmail(name string) string {
	return strings.ToLower(name)
}

func CleanName(name string) string {
	middleman := CleanEmail(name)
	return strings.ToUpper(string(name[0])) + middleman[1:]
}
