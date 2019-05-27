package util

var Types = types{}
type types struct {}

// Check if value is a string
func (t types) String(value interface{}) bool {
	_, ok := value.(string)
	return ok
}

// Check if value is a uint8
func (t types) Uint8(value interface{}) bool {
	// Check initial type
	if _, ok := value.(float64); !ok {
		return false
	}
	// Check if exceeds uint8 range
	return int(uint8(value.(float64))) == int(value.(float64))
}

// Check if value is a uint16
func (t types) Uint16(value interface{}) bool {
	// Check initial type
	if _, ok := value.(float64); !ok {
		return false
	}
	// Check if exceeds uint16 range
	return int(uint16(value.(float64))) == int(value.(float64))
}

// Check if value is a uint32
func (t types) Uint32(value interface{}) bool {
	// Check initial type
	if _, ok := value.(float64); !ok {
		return false
	}
	// Check if exceeds uint32 range
	return int(uint32(value.(float64))) == int(value.(float64))
}

// Check if value is an array of strings
func (t types) StringArray(value interface{}) bool {
	// Check if array
	if _, ok := value.([]interface{}); !ok {
		return false
	}

	// Check each value
	for _, v := range value.([]interface{}) {
		if _, ok := v.(string); !ok {
			return false
		}
	}

	return true
}
