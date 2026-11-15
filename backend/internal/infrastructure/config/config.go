// Package config carga y valida las variables de entorno del backend.
package config

import (
	"fmt"
	"os"
	"strconv"
	"strings"
)

type Config struct {
	AppEnv            string
	Puerto            string
	URLFrontend       string
	ProxiesConfiables []string

	DatabaseURL       string
	DatabaseURLDirect string
	DatabaseURLTest   string

	RedisURL string

	JWTSecreto  string
	JWTEmisor   string
	SesionHoras int

	MaxIntentosLogin         int
	MinutosBloqueo           int
	MinutosTokenRecuperacion int
	HorasTokenActivacion     int

	SMTPHost      string
	SMTPPuerto    int
	SMTPUsuario   string
	SMTPClave     string
	SMTPRemitente string

	// S3 (opcional): sin S3_BUCKET la API arranca, pero los adjuntos de trámites no están disponibles.
	S3Bucket    string
	S3Region    string
	S3AccessKey string
	S3SecretKey string
	S3Endpoint  string // vacío = AWS; para MinIO local

	// Servicio de OCR y compresión (ia-service, CU-15). Sin IA_SERVICE_URL los documentos se suben pero no se procesan.
	IAServiceURL   string
	IAServiceToken string

	// Intérprete de lenguaje natural de los reportes (Groq). Sin GROQ_API_KEY la función responde 503.
	GroqKey    string
	GroqModelo string

	// Umbral de días de préstamo de una carpeta para las alertas de mora del tablero (CU-25).
	DiasMaximosPrestamo int

	SeedDirectorCorreo    string
	SeedDirectorNombres   string
	SeedDirectorApellidos string
}

// Cargar lee las variables de entorno y falla con un mensaje claro si falta alguna obligatoria
// o si JWT_SECRETO tiene menos de 32 bytes.
func Cargar() (*Config, error) {
	var faltantes []string
	req := func(nombre string) string {
		v := os.Getenv(nombre)
		if v == "" {
			faltantes = append(faltantes, nombre)
		}
		return v
	}
	reqInt := func(nombre string) int {
		v := req(nombre)
		n, err := strconv.Atoi(v)
		if v != "" && err != nil {
			faltantes = append(faltantes, nombre+" (debe ser numérico)")
		}
		return n
	}

	c := &Config{
		AppEnv:                   getEnvDefault("APP_ENV", "development"),
		Puerto:                   getEnvDefault("APP_PUERTO", "8080"),
		URLFrontend:              req("APP_URL_FRONTEND"),
		ProxiesConfiables:        splitCSV(os.Getenv("PROXIES_CONFIABLES")),
		DatabaseURL:              req("DATABASE_URL"),
		DatabaseURLDirect:        os.Getenv("DATABASE_URL_DIRECT"),
		DatabaseURLTest:          os.Getenv("DATABASE_URL_TEST"),
		RedisURL:                 req("REDIS_URL"),
		JWTSecreto:               req("JWT_SECRETO"),
		JWTEmisor:                req("JWT_EMISOR"),
		SesionHoras:              reqInt("SESION_HORAS"),
		MaxIntentosLogin:         reqInt("MAX_INTENTOS_LOGIN"),
		MinutosBloqueo:           reqInt("MINUTOS_BLOQUEO"),
		MinutosTokenRecuperacion: reqInt("MINUTOS_TOKEN_RECUPERACION"),
		HorasTokenActivacion:     reqInt("HORAS_TOKEN_ACTIVACION"),
		SMTPHost:                 req("SMTP_HOST"),
		SMTPPuerto:               reqInt("SMTP_PUERTO"),
		SMTPUsuario:              os.Getenv("SMTP_USUARIO"),
		SMTPClave:                os.Getenv("SMTP_CLAVE"),
		SMTPRemitente:            req("SMTP_REMITENTE"),
		S3Bucket:                 os.Getenv("S3_BUCKET"),
		S3Region:                 primero("AWS_REGION", "S3_REGION", "us-east-1"),
		S3AccessKey:              primero("AWS_ACCESS_KEY_ID", "S3_ACCESS_KEY", ""),
		S3SecretKey:              primero("AWS_SECRET_ACCESS_KEY", "S3_SECRET_KEY", ""),
		IAServiceURL:             os.Getenv("IA_SERVICE_URL"),
		IAServiceToken:           os.Getenv("IA_SERVICE_TOKEN"),
		GroqKey: os.Getenv("GROQ_API_KEY"),
		GroqModelo: os.Getenv("GROQ_MODEL"),
		S3Endpoint:               os.Getenv("S3_ENDPOINT"),
		SeedDirectorCorreo:       os.Getenv("SEED_DIRECTOR_CORREO"),
		SeedDirectorNombres:      os.Getenv("SEED_DIRECTOR_NOMBRES"),
		SeedDirectorApellidos:    os.Getenv("SEED_DIRECTOR_APELLIDOS"),
	}
	// Opcional: sin DIAS_MAXIMOS_PRESTAMO el tablero usa 15 días.
	if v := os.Getenv("DIAS_MAXIMOS_PRESTAMO"); v != "" {
		n, err := strconv.Atoi(v)
		if err != nil || n < 1 {
			faltantes = append(faltantes, "DIAS_MAXIMOS_PRESTAMO (debe ser un entero mayor que 0)")
		}
		c.DiasMaximosPrestamo = n
	}

	if len(faltantes) > 0 {
		return nil, fmt.Errorf("variables de entorno faltantes o inválidas: %s", strings.Join(faltantes, ", "))
	}
	if len(c.JWTSecreto) < 32 {
		return nil, fmt.Errorf("JWT_SECRETO debe tener al menos 32 bytes (tiene %d)", len(c.JWTSecreto))
	}
	return c, nil
}

func getEnvDefault(nombre, porDefecto string) string {
	if v := os.Getenv(nombre); v != "" {
		return v
	}
	return porDefecto
}

// primero devuelve la variable estándar de AWS o, si falta, el nombre antiguo S3_*.
func primero(nombre, alterno, porDefecto string) string {
	if v := os.Getenv(nombre); v != "" {
		return v
	}
	return getEnvDefault(alterno, porDefecto)
}

func splitCSV(v string) []string {
	if v == "" {
		return nil
	}
	partes := strings.Split(v, ",")
	for i := range partes {
		partes[i] = strings.TrimSpace(partes[i])
	}
	return partes
}
