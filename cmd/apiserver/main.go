package main

import (
	"fmt"
	"k/cmd/config"
	"k/restserver/resthandlers"
	"k/storageserver"
	"k/usecases"
	"log"
	"net/http"
)

type APIServer struct {
	storageServersPool storageserver.Pool
}

func main() {
	cfg := config.Config{
		StorageServersPath: []string{
			"localhost:50051",
			"localhost:50052",
			"localhost:50053",
			"localhost:50054",
			"localhost:50055",
			"localhost:50056",
		},
	}

	storageServersPool := storageserver.NewPool(nil)
	for i := 0; i < len(cfg.StorageServersPath); i++ {
		dataServer, err := storageserver.NewServer(cfg.StorageServersPath[i], nil)
		if err != nil {
			log.Fatal(err)
		}
		storageServersPool.Append(dataServer)
	}

	uc := usecases.New(storageServersPool)

	http.HandleFunc("/upload", resthandlers.UploadHandler{Usecases: uc}.ServeHTTP)
	http.HandleFunc("/download", resthandlers.DownloadHandler{Usecases: uc}.ServeHTTP)
	//http.HandleFunc("/download", fileServer.handleDownload)

	fmt.Println("Server is running on :8080")
	http.ListenAndServe(":8080", nil)
}
