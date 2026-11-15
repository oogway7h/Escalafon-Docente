package tramites

import (
	"context"
	"strings"

	"escalafon/backend/internal/application/ports"
	"escalafon/backend/internal/domain/auditoria"
	"escalafon/backend/internal/domain/tramite"
)

// --- CU-17: restituir un documento dado de baja ------------------------------------------------
//
// Operación exclusiva del Director (DOCUMENTO_RESTITUIR). La baja de CU-16 es lógica, así que
// restituir solo vuelve a marcar el documento como activo: las búsquedas filtran por `activo`, de
// modo que reaparece en ellas sin reprocesar el OCR.

// ListarDocumentosDeBaja devuelve, paginados, los documentos retirados (de trámites y de expedientes).
type ListarDocumentosDeBaja struct {
	Tramites tramite.Repository
}

func (c *ListarDocumentosDeBaja) Ejecutar(ctx context.Context, busqueda string, pagina, tamano int) ([]tramite.DocumentoDeBaja, int64, error) {
	return c.Tramites.ListarDocumentosDeBaja(ctx, strings.TrimSpace(busqueda), pagina, tamano)
}

type RestituirDocumento struct {
	Tramites  tramite.Repository
	Auditoria auditoria.Repository
	Reloj     ports.Clock
}

// Ejecutar reactiva el documento. Exige motivo (queda en auditoría) y no toca expedientes de
// trámites remitidos o archivados, que son de solo lectura.
func (c *RestituirDocumento) Ejecutar(ctx context.Context, actorID, docID, motivo string) error {
	motivo = strings.TrimSpace(motivo)
	if motivo == "" {
		return tramite.ErrMotivoRestitucion
	}
	d, err := c.Tramites.ObtenerDocumentoDeBaja(ctx, docID)
	if err != nil {
		return err
	}
	if d.Activo {
		return tramite.ErrDocumentoActivo
	}
	if d.TramiteID != "" && d.TramiteEstado.Cerrado() {
		return d.TramiteEstado.ErrCerrado()
	}
	if err := c.Tramites.RestituirDocumento(ctx, d.ID); err != nil {
		return err
	}
	entidad, entidadID := auditoria.EntidadTramite, d.TramiteID
	if d.TramiteID == "" {
		entidad, entidadID = auditoria.EntidadDocente, d.DocenteID
	}
	_ = c.Auditoria.Registrar(ctx, &auditoria.Evento{
		UsuarioID: &actorID, Accion: auditoria.AccionRestituirDocumento, Entidad: entidad, EntidadID: entidadID, Fecha: c.Reloj.Ahora(),
		Detalle: map[string]any{"documentoId": d.ID, "documento": d.NombreOriginal, "motivo": motivo, "motivoBaja": d.MotivoBaja},
	})
	return nil
}
