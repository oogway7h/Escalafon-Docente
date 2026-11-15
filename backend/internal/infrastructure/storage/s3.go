// Package storage implementa ports.ObjectStorage sobre S3 (o un servicio compatible como MinIO).
package storage

import (
	"bytes"
	"context"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"

	"escalafon/backend/internal/domain/tramite"
)

type S3 struct {
	cliente *s3.Client
	firma   *s3.PresignClient
	bucket  string
}

type ConfigS3 struct {
	Bucket, Region, AccessKey, SecretKey string
	Endpoint                             string // vacío = AWS; con valor (MinIO) usa path-style
}

func NewS3(c ConfigS3) *S3 {
	opts := s3.Options{
		Region:       c.Region,
		Credentials:  credentials.NewStaticCredentialsProvider(c.AccessKey, c.SecretKey, ""),
		UsePathStyle: c.Endpoint != "",
	}
	if c.Endpoint != "" {
		opts.BaseEndpoint = aws.String(c.Endpoint)
	}
	cliente := s3.New(opts)
	return &S3{cliente: cliente, firma: s3.NewPresignClient(cliente), bucket: c.Bucket}
}

func (s *S3) Bucket() string { return s.bucket }

func (s *S3) Guardar(ctx context.Context, clave, tipoMime string, datos []byte) error {
	_, err := s.cliente.PutObject(ctx, &s3.PutObjectInput{
		Bucket: &s.bucket, Key: &clave, Body: bytes.NewReader(datos), ContentType: &tipoMime,
	})
	return err
}

func (s *S3) URLFirmada(ctx context.Context, clave string, vigencia time.Duration) (string, error) {
	r, err := s.firma.PresignGetObject(ctx, &s3.GetObjectInput{Bucket: &s.bucket, Key: &clave}, s3.WithPresignExpires(vigencia))
	if err != nil {
		return "", err
	}
	return r.URL, nil
}

func (s *S3) URLSubida(ctx context.Context, clave, tipoMime string, vigencia time.Duration) (string, error) {
	r, err := s.firma.PresignPutObject(ctx, &s3.PutObjectInput{Bucket: &s.bucket, Key: &clave, ContentType: &tipoMime}, s3.WithPresignExpires(vigencia))
	if err != nil {
		return "", err
	}
	return r.URL, nil
}

func (s *S3) Tamano(ctx context.Context, clave string) (int64, error) {
	r, err := s.cliente.HeadObject(ctx, &s3.HeadObjectInput{Bucket: &s.bucket, Key: &clave})
	if err != nil {
		return 0, err
	}
	return aws.ToInt64(r.ContentLength), nil
}

// NoConfigurado se usa cuando no hay variables S3: la API arranca, pero subir o descargar
// adjuntos responde 503.
type NoConfigurado struct{}

func (NoConfigurado) Bucket() string { return "" }
func (NoConfigurado) Guardar(context.Context, string, string, []byte) error {
	return tramite.ErrAlmacenNoDisponible
}
func (NoConfigurado) URLFirmada(context.Context, string, time.Duration) (string, error) {
	return "", tramite.ErrAlmacenNoDisponible
}
func (NoConfigurado) URLSubida(context.Context, string, string, time.Duration) (string, error) {
	return "", tramite.ErrAlmacenNoDisponible
}
func (NoConfigurado) Tamano(context.Context, string) (int64, error) {
	return 0, tramite.ErrAlmacenNoDisponible
}
