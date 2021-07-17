package serializer

import (
	"fmt"
	"regexp"
)

// findConfig find the regex group and if required and not found return an error, otherwise return the string found
func findConfig(source *string, regex *regexp.Regexp, isRequired bool) (value string, err error) {
	match := regex.FindStringSubmatch(*source)
	if match == nil && isRequired {
		return value, fmt.Errorf("configuration not found")
	}
	return match[1], nil
}
