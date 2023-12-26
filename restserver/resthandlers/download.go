package resthandlers

import (
	"fmt"
	"k/usecases"
	"net/http"
)

type DownloadHandler struct {
	Usecases *usecases.Usecases
}

func (h DownloadHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	params := getParamsFromRequest(r)
	file, err := h.Usecases.Download(params)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprint(w, "File download error", err.Error())
	}

	w.WriteHeader(http.StatusOK)
	w.Write(file.Bytes())
}

func getParamsFromRequest(r *http.Request) usecases.DownloadParams {
	return usecases.DownloadParams{}
}
