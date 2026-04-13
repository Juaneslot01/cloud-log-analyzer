package main

import (
	"context"
	"log"
	"net"
	"os"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	pb "github.com/juaneslot01/cloud-log-analyzer/api/proto"
	"github.com/juaneslot01/cloud-log-analyzer/internal/ingestor"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func main() {
	// 1. Cargar configuracion de AWS (lee credenciales de /.aws/credentials)
	cfg, err := config.LoadDefaultConfig(context.TODO(), config.WithRegion("us-east-2"))
	if err != nil {
		log.Fatalf("No se pudo cargar la config de AWS: %v", err)
	}

	s3Client := s3.NewFromConfig(cfg)
	bucketName := os.Getenv("S3_BUCKET_NAME")
	if bucketName == "" {
		bucketName = "juanes-logs-bucket"
	}

	// 2. Setup del servidor y canal
	logQueue := make(chan *pb.LogRequest, 100)
	ingestorSrv := ingestor.NewServer(logQueue, s3Client, bucketName)

	go ingestorSrv.RunLogWorker()

	// 3. gRPC Setup
	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	// Register interceptor
	grpcServer := grpc.NewServer(
		grpc.UnaryInterceptor(ingestor.AuthInterceptor),
	)

	pb.RegisterLogServiceServer(grpcServer, ingestorSrv)
	reflection.Register(grpcServer)

	log.Printf("Server with Auth initialize on %v", lis.Addr())
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
