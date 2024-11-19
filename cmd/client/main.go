package main

import (
	storage "GoStorage/pkg/grpc/storage/proto"
	"context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"log"
)

var client storage.StorageServiceClient
var MB = 1024 * 1024

func main() {
	log.Println("client start")
	address := "localhost:8080"

	// セキュリティを無効にして接続
	opt := grpc.WithTransportCredentials(insecure.NewCredentials())
	conn, err := grpc.NewClient(address, opt)
	if err != nil {
		log.Fatal("connection error: ", err)
		return
	}
	defer conn.Close()

	client = storage.NewStorageServiceClient(conn)

	data := []byte("Hello, World!")
	NewUploadFile(data)
}

func NewUploadFile(data []byte) {
	// streamを生成
	stream, err := client.NewUploadFile(context.Background())
	if err != nil {
		log.Fatal("upload file error: ", err)
		return
	}

	size := len(data)

	// リクエストを生成
	req := &storage.NewUploadFileRequest{
		Meta: &storage.FileMeta{
			Id:       "",
			Filename: "",
			Path:     "",
			Size:     uint64(size),
		},
		Data: data,
	}

	start := 0
	end := 0

	for (size - start) > 0 {
		start = end
		if (size - start) > MB {
			end = start + MB
		} else {
			end = size
		}

		req.Data = data[start:end]
		if err := stream.Send(req); err != nil {
			log.Fatal("send error: ", err)
			return
		}

	}

	res, err := stream.CloseAndRecv()
	if err != nil {
		log.Fatal("close and recv error: ", err)
		return
	}
	log.Println("upload file response: ", res.GetSuccess())
}
