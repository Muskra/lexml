package lexml

import (
    "bytes"
)

// IntEq function checks equality of two given Int values
func IntEq(orig int, given int) bool {
	return (orig == given)
}

// ItrEq function checks equality of two given String values
func StrEq(orig string, given string) bool {
	return bytes.Equal([]byte(orig), []byte(given))
}
