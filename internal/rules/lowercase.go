package rules

import (
	"unicode"
	"unicode/utf8"
)

const LowercaseMessageDiagnostic = "log message must start with a lowercase letter"

func CheckLowercase(message string) bool {
	if message == "" {
		return true
	}

	r, _ := utf8.DecodeRuneInString(message)
	return unicode.IsLower(r)
}
