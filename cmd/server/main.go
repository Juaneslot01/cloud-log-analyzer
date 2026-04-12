package main

import (
	"log"
	"net"

	pb "github.com/juaneslot01/cloud-log-analyzer/api/proto"
	"github.com/juaneslot01/cloud-log-analyzer/internal/ingestor"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func main() {
	// 1. Crear el canal (Buffer de 100 para manejar rafagas)
	logQueue := make(chan *pb.LogRequest, 100)

	// 2. Instanciar nuestro servidor interno
	ingestorSrv := ingestor.NewServer(logQueue)

	// 3. Arrancar el Worker en segundo plano (Goroutine)
	go ingestorSrv.RunLogWorker()

	// 4. Configurar el socket TCP
	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("Error al abrir puerto: %v", err)
	}

	// 5. Configurar gRPC y registrar servicios
	grpcServer := grpc.NewServer()
	pb.RegisterLogServiceServer(grpcServer, ingestorSrv)

	// Habilitar reflexion para grpcurl
	reflection.Register(grpcServer)

	log.Printf("Servidor gRPC y Worker iniciados en %v", lis.Addr())

	// 6. Empezar a servir
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("Error al servir gRPC: %v", err)
	}
}
