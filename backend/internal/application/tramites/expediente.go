package tramites

import (
	"context"
	"errors"
	"strings"
	"time"

	"github.com/google/uuid"

	"escalafon/backend/internal/application/ports"
	"escalafon/backend/internal/domain/auditoria"
	"escalafon/backend/internal/domain/docente"
	"escalafon/backend/internal/domain/tramite"
)

// --- Expediente digital del docente ----------------------------------------------------------
//
// El Personal de Archivo escanea la carpeta física y sube cada documento al docente con su tipo
// (Currículum, Regularización, Titular A…). Usa el mismo flujo que CU-15 (URL firmada a S3, OCR y
// compresión en el ia-service); el texto queda buscable al terminar el OCR, de modo que el analista
// encuentra una resolución o un título sin pasar página por página.

func registrarDocente(ctx context.Context, a auditoria.Repository, actorID *string, accion, docenteID string, ahora time.Time, detalle map[string]any) {
	_ = a.Registrar(ctx, &auditoria.Evento{
		UsuarioID: actorID, Accion: accion, Entidad: auditoria.EntidadDocente, EntidadID: docenteID, Detalle: detalle, Fecha: ahora,
	})
}

type PrepararSubidaDocente struct {
	Docentes docente.Repository
	Tramites tramite.Repository
	Almacen  ports.ObjectStorage
}

type SubidaDocenteEntrada struct {
	DocenteID, Nombre, TipoMime string
	Tipo                        tramite.TipoDocumento
	Tamano                      int64
}

// Ejecutar valida el archivo declarado, lo registra como SUBIENDO y devuelve la URL de subida.
func (c *PrepararSubidaDocente) Ejecutar(ctx context.Context, in SubidaDocenteEntrada) (*SubidaPreparada, error) {
	doc, err := c.Docentes.BuscarPorID(ctx, in.DocenteID)
	if err != nil {
		return nil, err
	}
	if !in.Tipo.Valido() {
		return nil, tramite.ErrTipoDocumentoInvalido
	}
	nombre := strings.TrimSpace(in.Nombre)
	if nombre == "" || in.Tamano <= 0 {
		return nil, tramite.ErrArchivoInvalido
	}
	if !tiposPermitidos[in.TipoMime] {
		return nil, tramite.ErrTipoNoPermitido
	}
	if in.Tamano > limitePorTipo(in.TipoMime) {
		return nil, tramite.ErrArchivoGrande
	}
	if r := []rune(nombre); len(r) > 200 {
		nombre = string(r[len(r)-200:])
	}

	id := uuid.NewString()
	clave := claveDocumento(carpetaDocente(doc.Codigo, doc.ID), "expediente", id, nombre)
	url, err := c.Almacen.URLSubida(ctx, clave, in.TipoMime, time.Hour)
	if err != nil {
		return nil, err
	}
	d := tramite.Documento{
		ID: id, DocenteID: doc.ID, TipoDocumento: in.Tipo, NombreOriginal: nombre, TipoMime: in.TipoMime, Tamano: in.Tamano,
		Bucket: c.Almacen.Bucket(), ClaveObjeto: clave, Activo: true, EstadoOCR: tramite.OCRSubiendo,
	}
	if err := c.Tramites.CrearDocumentoDocente(ctx, doc.ID, &d); err != nil {
		return nil, err
	}
	return &SubidaPreparada{Documento: d, URL: url}, nil
}

// EjecutarDocente confirma la subida de un escaneo del expediente y lanza su OCR; también reintenta
// uno cuyo procesamiento falló.
func (c *ConfirmarSubida) EjecutarDocente(ctx context.Context, actorID, docenteID, docID string) (*tramite.Documento, error) {
	return c.confirmar(ctx, "", docenteID, docID, func(d *tramite.Documento, tam int64) {
		registrarDocente(ctx, c.Auditoria, &actorID, auditoria.AccionCargarDocumento, docenteID, c.Reloj.Ahora(),
			map[string]any{"documentoId": d.ID, "documento": d.NombreOriginal, "tipo": d.TipoDocumento, "tamano": tam})
	})
}

// ListarExpediente devuelve los escaneos del docente; las bajas solo las ve quien puede restituirlas.
type ListarExpediente struct {
	Docentes docente.Repository
	Tramites tramite.Repository
}

func (c *ListarExpediente) Ejecutar(ctx context.Context, docenteID string, verBajas bool) ([]tramite.Documento, error) {
	if _, err := c.Docentes.BuscarPorID(ctx, docenteID); err != nil {
		return nil, err
	}
	return c.Tramites.ListarDocumentosDocente(ctx, docenteID, verBajas)
}

// BuscarEnExpediente busca por contenido en todos los documentos procesados del docente.
type BuscarEnExpediente struct {
	Docentes docente.Repository
	Tramites tramite.Repository
}

func (c *BuscarEnExpediente) Ejecutar(ctx context.Context, docenteID, consulta string) ([]tramite.CoincidenciaDocente, error) {
	consulta = strings.TrimSpace(consulta)
	if n := len([]rune(consulta)); n < 2 || n > 100 {
		return nil, tramite.ErrConsultaInvalida
	}
	if _, err := c.Docentes.BuscarPorID(ctx, docenteID); err != nil {
		return nil, err
	}
	return c.Tramites.BuscarEnDocente(ctx, docenteID, consulta)
}

// cargarDeExpediente trae un escaneo y comprueba que sea del docente indicado.
func cargarDeExpediente(ctx context.Context, repo tramite.Repository, docenteID, docID string) (*tramite.Documento, error) {
	d, err := repo.BuscarDocumento(ctx, "", docID)
	if err != nil {
		if errors.Is(err, tramite.ErrNoEncontrado) {
			return nil, tramite.ErrDocumentoNoEncontrado
		}
		return nil, err
	}
	if d.DocenteID != docenteID {
		return nil, tramite.ErrDocumentoNoEncontrado
	}
	return d, nil
}

// UrlExpediente entrega el enlace firmado para ver el escaneo (la versión comprimida si existe).
type UrlExpediente struct {
	Tramites tramite.Repository
	Almacen  ports.ObjectStorage
}

func (c *UrlExpediente) Ejecutar(ctx context.Context, docenteID, docID string, verBajas bool) (string, error) {
	d, err := cargarDeExpediente(ctx, c.Tramites, docenteID, docID)
	if err != nil {
		return "", err
	}
	if !d.Activo && !verBajas {
		return "", tramite.ErrDocumentoDeBaja
	}
	return c.Almacen.URLFirmada(ctx, d.ClaveDescarga(), 30*time.Minute)
}

// DarDeBajaExpediente retira un escaneo subido por error (baja lógica, con motivo y auditoría).
type DarDeBajaExpediente struct {
	Tramites  tramite.Repository
	Auditoria auditoria.Repository
	Reloj     ports.Clock
}

func (c *DarDeBajaExpediente) Ejecutar(ctx context.Context, actorID, docenteID, docID, motivo string) error {
	motivo = strings.TrimSpace(motivo)
	if motivo == "" {
		return tramite.ErrMotivoBajaRequerido
	}
	d, err := cargarDeExpediente(ctx, c.Tramites, docenteID, docID)
	if err != nil {
		return err
	}
	if !d.Activo {
		return tramite.ErrDocumentoDeBaja
	}
	ahora := c.Reloj.Ahora()
	if err := c.Tramites.DarDeBajaDocumento(ctx, "", d.ID, motivo, actorID, ahora); err != nil {
		return err
	}
	registrarDocente(ctx, c.Auditoria, &actorID, auditoria.AccionBajaDocumento, docenteID, ahora,
		map[string]any{"documentoId": d.ID, "documento": d.NombreOriginal, "motivo": motivo})
	return nil
}
