package rules

import "unicode"

const EnglishMessageDiagnostic = "log message must contain English text only"

func CheckEnglish(message string) bool {
	for _, r := range message {
		if !unicode.IsLetter(r) {
			continue
		}

		if !isASCIILetter(r) {
			return false
		}
	}

	return true
}

func isASCIILetter(r rune) bool {
	return ('a' <= r && r <= 'z') || ('A' <= r && r <= 'Z')
}
