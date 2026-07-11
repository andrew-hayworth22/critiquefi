package models

import "fmt"

// ValidationErrors maps fields to their validation error
type ValidationErrors map[string]string

func (ve ValidationErrors) Error() string {
	return "validation failed"
}

func (ve ValidationErrors) Add(field, message string) {
	ve[field] = message
}

func (ve ValidationErrors) Any() bool {
	return len(ve) > 0
}

func (ve ValidationErrors) Equals(other ValidationErrors) bool {
	for k, v := range ve {
		if other[k] != v {
			fmt.Printf("'%s' != '%s'\n", v, other[k])
			return false
		}
		delete(other, k)
	}
	return len(other) == 0
}
