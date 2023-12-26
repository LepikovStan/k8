package main

import (
	"fmt"
	"google.golang.org/grpc"
	"k/cmd/config"
	"k/infrastructure"
	pb "k/storageserver"
	"log"
	"net"
	"os"
	"os/signal"
	"sync"
	"syscall"
)

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

	redisStorage, err := infrastructure.NewRedisFileStorage("localhost:6379")
	if err != nil {
		log.Fatal(err)
	}

	wg := &sync.WaitGroup{}
	wg.Add(len(cfg.StorageServersPath))
	for i := 0; i < len(cfg.StorageServersPath); i++ {
		// Create a gRPC server
		server := grpc.NewServer()

		// Register the data server
		pb.RegisterDataServiceServer(server, pb.NewDataServer(redisStorage))

		// Start the server
		listener, err := net.Listen("tcp", cfg.StorageServersPath[i])
		if err != nil {
			log.Fatalf("Failed to listen: %v", err)
		}

		go func(j int, s *grpc.Server, l net.Listener) {
			log.Println("Starting gRPC server on " + cfg.StorageServersPath[j])
			if err := s.Serve(l); err != nil {
				log.Fatalf("Failed to serve: %v", err)
			}
		}(i, server, listener)

		// Graceful shutdown on interrupt signal
		go func(s *grpc.Server, l net.Listener) {
			waitForShutdown(l, s)
			wg.Done()
		}(server, listener)
	}
	wg.Wait()

}

func waitForShutdown(listener net.Listener, server *grpc.Server) {
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)

	sig := <-sigCh
	fmt.Printf("Received signal %v. Shutting down...\n", sig)

	// Create a wait group to wait for the gRPC server to stop
	var wg sync.WaitGroup
	wg.Add(1)

	go func() {
		defer wg.Done()

		// Stop the gRPC server
		server.GracefulStop()

		// Close the listener
		listener.Close()
	}()

	// Wait for the goroutine to finish
	wg.Wait()

	fmt.Println("Server gracefully stopped.")
}
