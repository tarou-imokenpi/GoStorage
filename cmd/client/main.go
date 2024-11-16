package main

import (
	storage "CloudStorage/pkg/grpc/storage/proto"
	"context"
	"fmt"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"log"
)

var clinet storage.StorageServiceClient

func main() {
	address := "localhost:8080"

	// セキュリティを無効にして接続
	opt := grpc.WithTransportCredentials(insecure.NewCredentials())
	conn, err := grpc.NewClient(address, opt)
	if err != nil {
		log.Fatal("connection error: ", err)
		return
	}
	defer conn.Close()

	clinet := storage.NewStorageServiceClient(conn)

}

func NewUploadFile(data []byte) error {
	// streamを生成
	stream, err := clinet.NewUploadFile(context.Background())
	if err != nil {
		log.Fatal("upload file error: ", err)
		return err
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

	sendCount := size / (1024 * 1024)
	if size%(1024*1024) == 0 {
		if err := stream.Send(req); err != nil {
			log.Fatal("send error: ", err)
			return err
		}

	} else {
		// 1MBずつ送信
		// 1MB未満の場合は1回で送信
		for i := 0; i < sendCount; i++ {
			req.Data = data[i*1024*1024 : (i+1)*1024*1024]
			err := stream.Send(req)
			if err != nil {
				log.Fatal("send error: ", err)
				return err
			}
		}
	}

	res, err := stream.CloseAndRecv()
	if err != nil {
		log.Fatal("close and recv error: ", err)
		return err
	}

	fmt.Println("upload file response: ", res.GetSuccess())
	return nil
}
