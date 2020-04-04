package controllers

import (
	"github.com/gorilla/mux"
	"go-chunks-dowloader/services"
	"net/http"
	"strings"
)

type ChunkDownloaderController struct {
	ChunkDownloadService *services.ChunkMediaDownloaderService
}

func (chunkDownloaderController *ChunkDownloaderController) DownloadChunk(w http.ResponseWriter, r *http.Request)  {
	params := mux.Vars(r)

	storageNameSplit := strings.Split(params["storageName"], ".")

	if storageNameSplit[len(storageNameSplit) -1] != "ts" {
		http.Error(w, "chunk should be be MPEG2TS", http.StatusBadRequest)
		return
	}

	rsp, err := chunkDownloaderController.ChunkDownloadService.DownloadChunk(params["bucketName"], params["storageName"])
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "video/MP2T")  // video MPEG2 transport stream.
	w.Write(rsp)
}
