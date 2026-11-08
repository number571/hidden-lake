package dto

func MessageStateIsValid(pState IMessageState) bool {
	return pState != cBeginFlagState && pState < cEndFlagState
}

func Uint64InStateRange(pValue uint64) bool {
	return pValue != uint64(cBeginFlagState) && pValue < uint64(cEndFlagState)
}
