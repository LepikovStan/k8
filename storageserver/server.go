package storageserver

import (
	"context"
	"fmt"
	"google.golang.org/grpc"
	"k/repository"
	"log"
)

type DataServer struct {
	storage repository.FileStorage
	UnimplementedDataServiceServer
}

func getDataKeyFromRequest(userID int32, filename string) string {
	return fmt.Sprintf("%d_%s", userID, filename)
}

// SaveFileChunk saves a file chunk to the data server.
//
// ctx - the context.Context object used for request cancellation.
// req - the DataRequest object containing the file chunk and metadata.
// Returns a DataResponse object and an error.
func (s *DataServer) SaveFileChunk(ctx context.Context, req *DataRequest) (*DataResponse, error) {
	s.storage.Save(getDataKeyFromRequest(req.GetUserId(), req.GetFilename()), req.BinaryData)

	return &DataResponse{}, nil
}

func (s *DataServer) GetFileChunk(ctx context.Context, req *GetDataRequest) (*GetDataResponse, error) {
	filebytes, err := s.storage.Get(getDataKeyFromRequest(req.GetUserId(), req.GetFilename()))

	return &GetDataResponse{BinaryData: filebytes}, err
}

func NewDataServer(storage repository.FileStorage) *DataServer {
	return &DataServer{
		storage: storage,
	}
}

type Server struct {
	addr          string
	spaceReserved int
	Client        DataServiceClient
}

func (s *Server) ReserveSpace(bb int) {
	s.spaceReserved += bb
}

func (s Server) SpaceReserved() int {
	return s.spaceReserved
}

func (s Server) Addr() string {
	return s.addr
}

func NewServer(addr string, storage repository.FileStorage) (*Server, error) {
	conn, err := grpc.Dial(addr, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("Failed to connect: %v", err)
	}
	defer conn.Close()

	// Create a gRPC client
	client := NewDataServiceClient(conn)

	return &Server{
		addr:          addr,
		spaceReserved: 0,
		Client:        client,
	}, nil
}

//
//func (s Server) Host() string {
//	return s.host
//}
//
//func (s Server) Port() string {
//	return s.port
//}
//
//func (s Server) GetConnectionPath() string {
//	return fmt.Sprintf("%s:%s", s.host, s.port)
//}
//
//func (s *Server) Upload(bb []byte) (int, error) {
//	s.bytesReserved += len(bb)
//	defer func() { s.bytesLoaded += len(bb) }()
//
//	return len(bb), nil
//}
//
//func NewServer(path string) Server {
//	parts := strings.Split(path, ":")
//
//	return Server{
//		host: parts[0],
//		port: parts[1],
//	}
//}
//
//type Meta struct {
//	host string
//	port string
//}
