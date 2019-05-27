package util

import "fmt"

// Check if a value exists within a map
func Exists(m map[string]interface{}, key string) bool {
	_, ok := m[key]
	return ok
}

// Convert []interface to []string)
func ConvertArrayToString(iarr []interface{}) ([]string, error) {
	strings := make([]string, len(iarr))

	for _, v := range iarr {
		if s, ok := v.(string); !ok {
			return []string{}, fmt.Errorf("not all values are strings")
		} else {
			strings = append(strings, s)
		}
	}

	return strings, nil
}
