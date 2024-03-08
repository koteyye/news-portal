package resthandler

import (
	"mime/multipart"
	"net/http"
)

func getFileFromMultipartform(w http.ResponseWriter, r *http.Request, key string) (multipart.File, *multipart.FileHeader, error) {
	err := r.ParseMultipartForm(1000)
	if err != nil {
		return nil, nil, err
	}

	file, fileHeader, err := r.FormFile(key)
	if err != nil {
		return nil, nil, err
	}
	defer file.Close()

	return file, fileHeader, nil
}
