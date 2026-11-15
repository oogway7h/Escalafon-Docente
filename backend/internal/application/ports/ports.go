package ports

import (
	"context"
	"time"

	"escalafon/backend/internal/domain/indicador"
)

// Clock aísla el tiempo para que los casos de uso sean deterministas en pruebas.
type Clock interface {
	Ahora() time.Time
}

// PasswordHasher aísla bcrypt de la capa de aplicación.
type PasswordHasher interface {
	Hash(plano string) (string, error)
	Comparar(hash, plano string) bool
}

// ClaimsSesion son los datos que viajan dentro del JWT (sección 6).
type ClaimsSesion struct {
	UsuarioID string
	RolID     int16
	Permisos  []string
	Version   int32
	JTI       string
	Expira    time.Time
}

// TokenService emite y verifica JWT de sesión.
type TokenService interface {
	Emitir(claims ClaimsSesion, vigencia time.Duration) (token string, err error)
	Verificar(token string) (*ClaimsSesion, error)
}

// GeneradorToken produce tokens aleatorios de un solo uso (recuperación/activación).
type GeneradorToken interface {
	// Generar produce el token en claro (base64url) y su hash SHA-256 en hex.
	Generar() (plano string, hash string, err error)
	// Hash calcula el SHA-256 en hex de un token en claro, para buscarlo en el repositorio.
	Hash(plano string) string
}

// Mailer envía los correos transaccionales de seguridad.
type Mailer interface {
	EnviarActivacion(destino, nombres, enlace string) error
	EnviarRecuperacion(destino, nombres, enlace string) error
	// EnviarCambioEstado avisa al docente del nuevo estado de su trámite (solo código y enlace al portal).
	EnviarCambioEstado(destino, nombres, codigo, estado, enlace string) error
}

// SessionCache es la caché de sesión sobre Redis (sección 5), con caída a BD si no está disponible.
type SessionCache interface {
	ObtenerSesion(ctx context.Context, usuarioID string) (version int32, activo bool, encontrado bool, err error)
	GuardarSesion(ctx context.Context, usuarioID string, version int32, activo bool) error
	InvalidarSesion(ctx context.Context, usuarioID string) error
	RevocarJTI(ctx context.Context, jti string, vigencia time.Duration) error
	EstaRevocado(ctx context.Context, jti string) (bool, error)
}

// TxManager ejecuta fn dentro de una transacción de base de datos.
type TxManager interface {
	WithinTx(ctx context.Context, fn func(ctx context.Context) error) error
}

// ObjectStorage guarda los adjuntos de los trámites (S3 o compatible).
type ObjectStorage interface {
	Bucket() string
	Guardar(ctx context.Context, clave, tipoMime string, datos []byte) error
	// URLFirmada devuelve un enlace temporal de descarga.
	URLFirmada(ctx context.Context, clave string, vigencia time.Duration) (string, error)
	// URLSubida devuelve un enlace temporal para que el navegador suba el objeto directamente (PUT),
	// sin pasar por la API: los PDF escaneados llegan a cientos de MB.
	URLSubida(ctx context.Context, clave, tipoMime string, vigencia time.Duration) (string, error)
	// Tamano devuelve el tamaño en bytes del objeto, o un error si no existe.
	Tamano(ctx context.Context, clave string) (int64, error)
}

// SolicitudProceso pide al servicio de IA que extraiga el texto (OCR) y comprima un objeto del almacén.
type SolicitudProceso struct {
	Bucket, Clave, TipoMime string
	ClaveSalida             string // dónde subir la versión comprimida (solo PDF)
}

// ResultadoProceso es el estado de un trabajo del servicio de IA: EN_COLA, PROCESANDO, LISTO o ERROR.
type ResultadoProceso struct {
	Estado          string
	Texto           string
	Paginas         []string // texto de cada página (una sola para una imagen)
	HashSHA256      string // del archivo original
	TamanoOriginal  int64
	TamanoFinal     int64
	ClaveComprimida string // vacío si la compresión no redujo el tamaño
	Error           string
}

// ProcesadorDocumentos es el servicio de OCR y compresión (ia-service, Tesseract + ocrmypdf).
type ProcesadorDocumentos interface {
	Iniciar(ctx context.Context, s SolicitudProceso) (id string, err error)
	Estado(ctx context.Context, id string) (*ResultadoProceso, error)
}

// ExportadorReportes compone el archivo de un reporte ya calculado (CU-26): PDF con membrete
// institucional o Excel tabular con filtros.
type ExportadorReportes interface {
	PDF(r *indicador.Reporte) ([]byte, error)
	Excel(r *indicador.Reporte) ([]byte, error)
}

// CategoriaRef es una categoría que el intérprete puede nombrar en un filtro.
type CategoriaRef struct {
	ID     int16
	Nombre string
}

// PeticionReporte es lo que un modelo de lenguaje entendió de una frase: los mismos campos que el
// formulario de reportes. Es una propuesta sin validar; InterpretarReporte la sanea.
type PeticionReporte struct {
	Tipo, Formato, Desde, Hasta, Estado, TipoTramite, EstadoLegajo string
	CategoriaID                                                    int16
	Explicacion                                                    string
}

// InterpreteReportes traduce una frase en lenguaje natural a los filtros de un reporte. Solo recibe
// la frase, la fecha de hoy y los nombres de las categorías: nunca datos de docentes ni de trámites.
type InterpreteReportes interface {
	Interpretar(ctx context.Context, frase string, hoy time.Time, categorias []CategoriaRef) (*PeticionReporte, error)
}
