package validation

import (
	"fmt"
)

var required = "field ini tidak boleh dikosongkan"
var minString = "field ini tidak boleh kurang dari %d"
var maxString = "field ini tidak boleh lebih dari %d"
var passwordAndRePassword = "password dan re password tidak sesuai"
var emailMsg = "email harus menggunakan yourmail@gmail.com"

func maxMinString(s string, min, max int) string {
	switch {
	case len(s) < min:
		return fmt.Sprintf(minString, min)
	case len(s) > max:
		return fmt.Sprintf(maxString, max)
	}

	return ""
}
