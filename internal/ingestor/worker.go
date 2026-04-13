package ingestor

import (
	"bytes"
	"context"
	"fmt"
	"log"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

// RunLogWorker escucha el canal de forma infinita
func (s *Server) RunLogWorker() {
	log.Println("Worker: Iniciado y esperando logs...")

	for logEntry := range s.LogQueue {

		now := time.Now()
		key := fmt.Sprintf("logs/%s/%s/%d.log",
			logEntry.ServiceName,
			now.Format("2006-01-02"),
			now.UnixNano(),
		)

		// Contenido del archivo sera el mensaje del log
		content := fmt.Sprintf("Level: %s\nMessage: %s\nTime: %s",
			logEntry.Level,
			logEntry.Message,
			logEntry.Timestamp,
		)

		// Subir a S3
		_, err := s.S3Client.PutObject(context.TODO(), &s3.PutObjectInput{
			Bucket: aws.String(s.Bucket),
			Key:    aws.String(key),
			Body:   bytes.NewReader([]byte(content)),
		})

		if err != nil {
			log.Printf("Error subiendo a S3: %s", err)
		} else {
			log.Printf("Log guardado en S3: %s", key)
		}
	}
}
