package ingestor

import (
	"context"

	pb "github.com/juaneslot01/cloud-log-analyzer/api/proto"
)

// server implementa el servicio LogServiceServer
type Server struct {
	pb.UnimplementedLogServiceServer
	LogQueue chan *pb.LogRequest
}

// NewServer es el constructor para nuestro servidor
func NewServer(queue chan *pb.LogRequest) *Server {
	return &Server{
		LogQueue: queue,
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
