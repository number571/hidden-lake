package msgdata

import (
	"errors"
	"io"
	"net/http"
	"strings"

	"github.com/number571/hidden-lake/internal/utils/chars"
)

func GetMessage(rawMsgBytes []byte, timestamp string) (SMessage, error) {
	switch {
	case isText(rawMsgBytes):
		textdata := unwrapText(rawMsgBytes)
		if textdata == "" {
			return SMessage{}, ErrMessageNull
		}
		return SMessage{
			FTimestamp: timestamp,
			FTextData:  textdata,
		}, nil
	case isFile(rawMsgBytes):
		filename, filedata := unwrapFile(rawMsgBytes)
		if filename == "" || filedata == "" {
			return SMessage{}, ErrUnwrapFile
		}
		return SMessage{
			FTimestamp: timestamp,
			FFileName:  filename,
			FFileData:  filedata,
		}, nil
	default:
		return SMessage{}, ErrUnknownMessageType
	}
}

func GetMessageBytes(pR *http.Request) ([]byte, error) {
	switch pR.FormValue("method") {
	case http.MethodPost:
		if pR.FormValue("ping") != "" {
			return nil, nil
		}
		strMsg := strings.TrimSpace(pR.FormValue("input_message"))
		if strMsg == "" {
			return nil, ErrMessageNull
		}
		if chars.HasNotGraphicCharacters(strMsg) {
			return nil, ErrHasNotWritableChars
		}
		return wrapText(strMsg), nil
	case http.MethodPut:
		filename, fileBytes, err := getUploadFile(pR)
		if err != nil {
			return nil, errors.Join(ErrUploadFile, err)
		}
		return wrapFile(filename, fileBytes), nil
	default:
		panic("got not supported method")
	}
}

func getUploadFile(pR *http.Request) (string, []byte, error) {
	// Get handler for filename, size and headers
	file, handler, err := pR.FormFile("input_file")
	if err != nil {
		return "", nil, errors.Join(ErrGetFormFile, err)
	}
	defer file.Close()

	if handler.Size == 0 {
		return "", nil, ErrReadFileSize
	}

	fileBytes, err := io.ReadAll(file)
	if err != nil {
		return "", nil, errors.Join(ErrReadFile, err)
	}

	return handler.Filename, fileBytes, nil
}
