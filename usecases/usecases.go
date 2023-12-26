package usecases

import (
	"context"
	"errors"
	"k/storageserver"
	"k/userfile"
	"log"
	"sync"
)

type Usecases struct {
	storageserverspool storageserver.Pool
}

func New(pool storageserver.Pool) *Usecases {
	return &Usecases{
		storageserverspool: pool,
	}
}

type UploadParams struct {
	UserID int
	File   userfile.File
}

func (uc Usecases) Upload(p UploadParams) error {
	fileparts := p.File.Divide(6)

	srvrs := uc.storageserverspool.LeastCompleted()
	for i := 0; i < len(srvrs); i++ {
		storageServer := srvrs[i]
		part := fileparts[i]

		storageServer.ReserveSpace(len(part.Bytes()))
		response, err := storageServer.Client.SaveFileChunk(context.Background(), &storageserver.DataRequest{
			BinaryData: part.Bytes(),
			Filename:   p.File.Name(),
			UserId:     1,
		})
		if err != nil {
			log.Println("error saving file chunk to:", storageServer.Addr())
			return err
		}
		if response.GetErrorMessage() != "" {
			log.Println("error saving file chunk to:", storageServer.Addr())
			return err
		}
	}

	uc.storageserverspool.SetFileServers(p.UserID, p.File.Name(), srvrs)

	return nil
}

type DownloadParams struct {
	UserID   int
	FileName string
}

func (uc Usecases) Download(p DownloadParams) (userfile.File, error) {
	var (
		Err            error
		filepartsbytes = make([][]byte, 6)
		filesize       = 0
		mu             = &sync.Mutex{}
		wg             = &sync.WaitGroup{}
	)

	wg.Add(6)
	srvrs := uc.storageserverspool.GetFileServers(p.FileName)
	for i := 0; i < len(srvrs); i++ {
		go func(j int) {
			defer wg.Done()
			storageServer := srvrs[j]
			response, err := storageServer.Client.GetFileChunk(context.Background(), &storageserver.GetDataRequest{
				Filename: p.FileName,
				UserId:   1,
			})
			if err != nil {
				Err = err
			}
			if response.GetErrorMessage() != "" {
				Err = errors.New(response.GetErrorMessage())
			}
			filepartsbytes[j] = response.GetBinaryData()
			mu.Lock()
			filesize += len(filepartsbytes[j])
			mu.Unlock()
		}(i)
	}
	wg.Wait()
	filesbytes := make([]byte, 0)

	for i := 0; i < len(filepartsbytes); i++ {
		filesbytes = append(filesbytes, filepartsbytes[i]...)
	}

	if Err != nil {
		return userfile.File{}, Err
	}

	return userfile.New(p.FileName, filesbytes), nil
}
