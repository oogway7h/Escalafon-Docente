package handlers

import (
	"github.com/gofiber/fiber/v2"

	"escalafon/backend/internal/application/tramites"
	"escalafon/backend/internal/domain/tramite"
	"escalafon/backend/internal/interfaces/http/apierror"
	"escalafon/backend/internal/interfaces/http/dto"
	"escalafon/backend/internal/interfaces/http/middleware"
)

// DocumentosHandler expone /api/documentos: la búsqueda global (CU-18) y la bandeja de documentos
// dados de baja con su restitución (CU-17), sobre adjuntos de trámites y escaneos de expedientes.
type DocumentosHandler struct {
	CasoBuscar    *tramites.BuscarDocumentos
	CasoBajas     *tramites.ListarDocumentosDeBaja
	CasoRestituir *tramites.RestituirDocumento
}

func (h *DocumentosHandler) Buscar(c *fiber.Ctx) error {
	pagina, tamano := paginacion(c)
	f := tramite.FiltroBusqueda{
		Consulta: c.Query("q"), Campo: tramite.CampoBusqueda(c.Query("campo")), Tipo: tramite.TipoDocumento(c.Query("tipo")),
		FacultadID: int16(c.QueryInt("facultad")), Recientes: c.Query("orden") == "recientes",
	}
	lista, total, err := h.CasoBuscar.Ejecutar(c.Context(), f, pagina, tamano)
	if err != nil {
		return apierror.Mapear(err)
	}
	out := dto.ListaResultadosDTO{Datos: make([]dto.ResultadoBusquedaDTO, len(lista)), Total: total, Pagina: pagina, Tamano: tamano}
	for i, r := range lista {
		out.Datos[i] = dto.ResultadoBusquedaDTO{
			DocumentoID: r.DocumentoID, TramiteID: r.TramiteID, TramiteCodigo: r.TramiteCodigo, TipoDocumento: string(r.TipoDocumento),
			NombreOriginal: r.NombreOriginal, DocenteID: r.DocenteID, DocenteNombre: r.DocenteNombre, DocenteCodigo: r.DocenteCodigo,
			Paginas: r.Paginas, Fragmento: r.Fragmento,
		}
	}
	return c.JSON(out)
}

func (h *DocumentosHandler) ListarBajas(c *fiber.Ctx) error {
	pagina, tamano := paginacion(c)
	lista, total, err := h.CasoBajas.Ejecutar(c.Context(), c.Query("buscar"), pagina, tamano)
	if err != nil {
		return apierror.Mapear(err)
	}
	out := dto.ListaBajasDTO{Datos: make([]dto.DocumentoDeBajaDTO, len(lista)), Total: total, Pagina: pagina, Tamano: tamano}
	for i, d := range lista {
		out.Datos[i] = dto.DocumentoDeBajaDTO{
			DocumentoDTO: aDocumentoDTO(d.Documento), TramiteID: d.TramiteID, TramiteCodigo: d.TramiteCodigo, TramiteEstado: string(d.TramiteEstado),
			DocenteID: d.DocenteID, DocenteNombre: d.DocenteNombre, DocenteCodigo: d.DocenteCodigo, BajaPorNombre: d.BajaPorNombre,
		}
	}
	return c.JSON(out)
}

func (h *DocumentosHandler) Restituir(c *fiber.Ctx) error {
	var in dto.RestituirEntrada
	if err := leerCuerpo(c, &in); err != nil {
		return err
	}
	actor, _ := middleware.DeContexto(c)
	if err := h.CasoRestituir.Ejecutar(c.Context(), actor.UsuarioID, c.Params("docId"), in.Motivo); err != nil {
		return apierror.Mapear(err)
	}
	return c.SendStatus(fiber.StatusNoContent)
}
