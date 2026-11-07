package dto

func MessageStateIsValid(pState IMessageState) bool {
	return pState != cBeginFlagState && pState < cEndFlagState
}
