package util

import (
	"net"
	"strconv"
	"strings"
)

// Validate a fields in a JSON request body
// Returns a string to be used as an error or empty if no error
func ValidateBody(body map[string]interface{}, keys []string, options map[string]map[string]string) (string, map[string]bool) {
	valid := make(map[string]bool)

	// Iterate over each value to check
	for _, key := range keys {
		// Set key as invalid to start
		valid[key] = false

		// Check if exists
		if exists := Exists(body, key); !exists && options[key]["required"] == "true" {
			return "field '" + key + "' is required", valid
		} else if !exists && options[key]["required"] == "false" {
			continue
		}

		// Check type based on options map
		switch options[key]["type"] {
		case "ipv4":
			if !Types.String(body[key]) {
				return "field '" + key + "' must be a string", valid
			} else if body[key].(string) == "" && options[key]["required"] == "true" {
				return "field '" + key + "' must be of length longer than 0", valid
			} else if ip := net.ParseIP(body[key].(string)); ip.To4().String() == "<nil>" {
				return "field '" + key + "' must be an IPv4 address", valid
			}

		case "ipv6":
			if !Types.String(body[key]) {
				return "field '" + key + "' must be a string", valid
			} else if body[key].(string) == "" && options[key]["required"] == "true" {
				return "field '" + key + "' must be of length longer than 0", valid
			} else if ip := net.ParseIP(body[key].(string)); ip.To4().String() != "<nil>" {
				return "field '" + key + "' must be an IPv4 address", valid
			}

		case "string":
			if !Types.String(body[key]) {
				return "field '" + key + "' must be a string", valid
			} else if body[key].(string) == "" && options[key]["required"] == "true" {
				return "field '" + key + "' must be of length longer than 0", valid
			} else if valuesString, ok := options[key]["oneOf"]; ok {
				values := strings.Split(valuesString, ",")
				if !StringInArray(body[key].(string), values) {
					return "field '" + key + "' must be one of " + valuesString, valid
				}
			}

		case "uint8":
			if !Types.Uint8(body[key]) {
				return "field '" + key + "' must be an integer between 0 and 255", valid
			} else if minString, ok := options[key]["min"]; ok {
				min, _ := strconv.ParseInt(minString, 10, 8)
				if uint8(min) > uint8(body[key].(float64)) {
					return "field '" + key + "' must be greater than " + minString, valid
				}
			} else if maxString, ok := options[key]["max"]; ok {
				max, _ := strconv.ParseInt(maxString, 10, 8)
				if uint8(max) < uint8(body[key].(float64)) {
					return "field '" + key + "' must be less than " + maxString, valid
				}
			}

		case "uint16":
			if !Types.Uint16(body[key]) {
				return "field '" + key + "' must be an integer between 0 and 65535", valid
			}

		case "uint32":
			if !Types.Uint32(body[key]) {
				return "field '" + key + "' must be an integer between 0 and 4294967296", valid
			}

		case "stringarray":
			if !Types.StringArray(body[key]) {
				return "field '" + key + "' must be an array of strings", valid
			} else if text, _ := ConvertArrayToString(body["text"].([]interface{})); len(text) < 1 {
				return "field '" + key + "' must be of at least length 1", valid
			}
		}

		// Set key as valid if
		valid[key] = true
	}

	return "", valid
}
