package main

import (
	storage "CloudStorage/pkg/grpc/storage/proto"
	"errors"
	"fmt"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"io"
	"log"
	"net"
	"os"
	"os/signal"
)

type StorageServiceServer struct {
	storage.UnimplementedStorageServiceServer
}

func NewStorageServiceServer() *StorageServiceServer {
	return &StorageServiceServer{}
}

func (S *StorageServiceServer) NewUploadFile(stream storage.StorageService_NewUploadFileServer) error {
	data := make([]byte, 0)

	for {
		req, err := stream.Recv()

		if errors.Is(err, io.EOF) {

			fmt.Println("data: ", len(data))
			fmt.Println("all received")

			return stream.SendAndClose(&storage.NewUploadFileResponse{
				Success: true,
			})
		}
		if err != nil {
			return err
		}

		data = append(data, req.GetData()...)
	}
}

func main() {
	port := 8080
	listener, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		log.Fatal("failed to listen: ", err)
		return
	}

	s := grpc.NewServer()

	storage.RegisterStorageServiceServer(s, NewStorageServiceServer())

	reflection.Register(s)

	go func() {
		log.Println("start server on port: ", port)
		err := s.Serve(listener)
		if err != nil {
			log.Fatal("failed to serve: ", err)
			return
		}
	}()

	quit := make(chan os.Signal, 1)

	signal.Notify(quit, os.Interrupt)
	<-quit
	log.Println("shutdown server")
	s.GracefulStop()

}
