package resthandler

import (
	"errors"
	"mime/multipart"
	"net/http"
)

var errNoSuchFile = errors.New("no such file")

func getFileFromMultipartform(w http.ResponseWriter, r *http.Request, key string) (multipart.File, *multipart.FileHeader, error) {
	err := r.ParseMultipartForm(1000)
	if err != nil {
		return nil, nil, err
	}

	file, fileHeader, err := r.FormFile(key)
	if err != nil {
		return nil, nil, err
	}

	return file, fileHeader, nil
}
