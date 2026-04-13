package ingestor

import (
	"context"
	"os"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"

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

// AuthInterceptor validates the API-KEY from the metadata
func AuthInterceptor(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return nil, status.Errorf(codes.Unauthenticated, "No metadata")
	}

	// We read the key from env variables from k8s
	apiKey := os.Getenv("APP_API_KEY")
	values := md["api-key"]

	if len(values) == 0 || values[0] != apiKey {
		return nil, status.Errorf(codes.Unauthenticated, "Invalid API-KEY")
	}
	return handler(ctx, req)
}
