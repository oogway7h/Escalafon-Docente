package tramites

import (
	"context"
	"strings"

	"escalafon/backend/internal/domain/tramite"
)

// BuscarDocumentos es la búsqueda global de CU-18: contenido (incluidos números de resolución),
// nombre, código o cédula del docente, en todos los documentos activos. Los dados de baja nunca
// salen; el acceso lo controla el permiso DOCUMENTO_BUSCAR en la ruta.
type BuscarDocumentos struct {
	Tramites tramite.Repository
}

func (c *BuscarDocumentos) Ejecutar(ctx context.Context, f tramite.FiltroBusqueda, pagina, tamano int) ([]tramite.ResultadoBusqueda, int64, error) {
	f.Consulta = strings.TrimSpace(f.Consulta)
	if n := len([]rune(f.Consulta)); n < 2 || n > 100 {
		return nil, 0, tramite.ErrConsultaInvalida
	}
	if !f.Campo.Valido() {
		return nil, 0, tramite.ErrConsultaInvalida
	}
	if f.Tipo != "" && !f.Tipo.Valido() {
		return nil, 0, tramite.ErrTipoDocumentoInvalido
	}
	return c.Tramites.BuscarDocumentos(ctx, f, pagina, tamano)
}
