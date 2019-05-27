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

// Remove duplicates from array
func RemoveDuplicates(arr []string) []string {
	// Get all elements, harmlessly overwrite if already exists
	encountered := map[string]bool{}
	for _, v := range arr {
		encountered[v] = true
	}

	// Iterate over keys and add to result
	var result []string
	for key := range encountered {
		result = append(result, key)
	}
	return result
}
