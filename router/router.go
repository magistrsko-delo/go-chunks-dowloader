package router

import (
	"github.com/gorilla/mux"
	"go-chunks-dowloader/controllers"
)

type ChunkDownloaderRouter struct {
	Router *mux.Router
}

func (chunkRouter *ChunkDownloaderRouter) RegisterHandlers()  {
	controller :=  &controllers.ChunkDownloaderController{}
	(*chunkRouter).Router.HandleFunc("/chunk/bucket/{bucketName}/storage/{storageName}", controller.DownloadChunk).Methods("GET")
}