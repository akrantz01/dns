package util

import (
	"fmt"
	"github.com/akrantz01/krantz.dev/dns/db"
)

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

// Check if a record does not exist
func RecordDoesNotExist(r db.Record) bool {
	// This is a REALLY big hack, but its the best way I could think of
	return fmt.Sprintf("%v", r) == "<nil>"
}

// Check if value in array
func StringInArray(val string, list []string) bool {
	for _, el := range list {
		if el == val {
			return true
		}
	}
	return false
}
