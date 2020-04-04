package services

import (
	"errors"
	"go-chunks-dowloader/Models"
	"io/ioutil"
	"log"
	"net/http"
)

type ChunkMediaDownloaderService struct {
}

func (chunkMediaDownloaderService *ChunkMediaDownloaderService) DownloadChunk(bucketName string, storageName string) ([]byte, error)  {
	resp, err := http.Get(Models.GetEnvStruct().AwsStorageUrl + "v1/awsStorage/media/" + bucketName + "/" + storageName)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if  resp.StatusCode >= 300 || resp.StatusCode < 200  {
		return nil, errors.New("chunk not found: " + resp.Status + " path: " + resp.Request.URL.Path)
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Printf("Error reading body: %v", err)
		return nil, err
	}

	return body, nil

}