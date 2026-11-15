package security

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
)

// GeneradorAleatorio implementa ports.GeneradorToken: 32 bytes de crypto/rand en base64url,
// con el hash SHA-256 en hex como lo que se guarda en la base de datos (sección 6).
type GeneradorAleatorio struct{}

func NewGeneradorAleatorio() *GeneradorAleatorio { return &GeneradorAleatorio{} }

func (GeneradorAleatorio) Generar() (plano string, hash string, err error) {
	buf := make([]byte, 32)
	if _, err := rand.Read(buf); err != nil {
		return "", "", err
	}
	plano = base64.RawURLEncoding.EncodeToString(buf)
	return plano, hashSHA256(plano), nil
}

func (GeneradorAleatorio) Hash(plano string) string {
	return hashSHA256(plano)
}

func hashSHA256(plano string) string {
	suma := sha256.Sum256([]byte(plano))
	return hex.EncodeToString(suma[:])
}
