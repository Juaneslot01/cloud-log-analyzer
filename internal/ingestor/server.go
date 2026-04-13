package ingestor

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/service/s3"
	pb "github.com/juaneslot01/cloud-log-analyzer/api/proto"
)

// server implementa el servicio LogServiceServer
type Server struct {
	pb.UnimplementedLogServiceServer
	LogQueue chan *pb.LogRequest
	S3Client *s3.Client
	Bucket   string
}

// NewServer es el constructor para nuestro servidor
func NewServer(queue chan *pb.LogRequest, s3Client *s3.Client, bucket string) *Server {
	return &Server{
		LogQueue: queue,
		S3Client: s3Client,
		Bucket:   bucket,
	}
}

// SendLog recibe el log y lo encola para procesamiento asincrono
func (s *Server) SendLog(ctx context.Context, in *pb.LogRequest) (*pb.LogResponse, error) {
	// Encolamos el log. Si el canal esta lleno, esto bloqueara
	s.LogQueue <- in

	return &pb.LogResponse{
		Success: true,
		Message: "Log encolado correctamente",
	}, nil
}
