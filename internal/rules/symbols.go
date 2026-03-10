package rules

import "unicode"

const SymbolsMessageDiagnostic = "log message must not contain special symbols or emoji"

func CheckSymbols(message string) bool {
	for _, r := range message {
		if unicode.IsLetter(r) || unicode.IsDigit(r) || unicode.IsSpace(r) {
			continue
		}

		return false
	}

	return true
}
