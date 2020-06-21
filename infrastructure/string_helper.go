package infrastructure

import "strings"

type StringHelper interface {
	StringToURL(string) string
}

type stringHelper struct {
}

func ProvideStringHelper() StringHelper {
	return &stringHelper{}
}

func (helper *stringHelper) StringToURL(str string) string {
	lowerStr := strings.ToLower(str)
	return strings.Join(strings.Split(lowerStr, " "), "-")
}
