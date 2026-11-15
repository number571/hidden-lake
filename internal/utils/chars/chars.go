package chars

import "unicode"

func HasNotGraphicCharacters(pS string) bool {
	for _, c := range pS {
		if !unicode.IsPrint(c) {
			return true
		}
	}
	return false
}
