package resthandlers

import (
	"fmt"
	"io"
	"k/usecases"
	"k/userfile"
	"log"
	"net/http"
)

type UploadHandler struct {
	Usecases *usecases.Usecases
}

func (h UploadHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	log.Println("handleUpload")

	input, err := getUploadParamsFromRequest(r)
	if err != nil {
		http.Error(w, "Failed to read the file", http.StatusBadRequest)
	}

	if err := h.Usecases.Upload(input); err != nil {

	}

	w.WriteHeader(http.StatusOK)
	fmt.Fprint(w, "File uploaded successfully")
}

func getUploadParamsFromRequest(r *http.Request) (usecases.UploadParams, error) {
	file, header, err := r.FormFile("file")
	if err != nil {
		return usecases.UploadParams{}, err
	}
	defer file.Close()

	bb, err := io.ReadAll(file)
	if err != nil {
		return usecases.UploadParams{}, err
	}

	uploadedFile := userfile.New(header.Filename, bb)
	return usecases.UploadParams{
		UserID: 1,
		File:   uploadedFile,
	}, nil
}
