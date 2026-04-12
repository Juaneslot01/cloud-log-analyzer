package ingestor

import (
	"log"
)

// RunLogWorker escucha el canal de forma infinita
func (s *Server) RunLogWorker() {
	log.Println("Worker: Iniciado y esperando logs...")

	for logEntry := range s.LogQueue {
		// Aqui es donde ira la logica de AWS S3 / DynamoDB
		log.Printf(" [WORKER] Procesando: %s - %s: %s",
			logEntry.Level,
			logEntry.ServiceName,
			logEntry.Message,
		)
	}
}
