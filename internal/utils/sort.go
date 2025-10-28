package utils

import "strings"

// NormalizeAndValidateSort trims direction, validates fields against allowed set, and returns
// the original sort tokens that are valid and a slice of invalid field names (without the '-').
func NormalizeAndValidateSort(inputs []string, allowed map[string]struct{}) (valid []string, invalid []string) {
	if len(inputs) == 0 {
		return nil, nil
	}

	valid = make([]string, 0, len(inputs))
	for _, token := range inputs {
		field := strings.TrimPrefix(token, "-")
		if _, ok := allowed[field]; ok {
			valid = append(valid, token)
		} else {
			invalid = append(invalid, field)
		}
	}
	return
}
