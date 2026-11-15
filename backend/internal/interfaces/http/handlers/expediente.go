package handlers

import (
	"github.com/gofiber/fiber/v2"

	"escalafon/backend/internal/application/tramites"
	"escalafon/backend/internal/domain/rol"
	"escalafon/backend/internal/domain/tramite"
	"escalafon/backend/internal/interfaces/http/apierror"
	"escalafon/backend/internal/interfaces/http/dto"
	"escalafon/backend/internal/interfaces/http/middleware"
)

// ExpedienteHandler expone /api/docentes/:id/documentos: el expediente digital de cada docente, que
// escanea el Personal de Archivo y consultan Archivo, analistas y Director (búsqueda por contenido).
type ExpedienteHandler struct {
	CasoPrepararSubida  *tramites.PrepararSubidaDocente
	CasoConfirmarSubida *tramites.ConfirmarSubida
	CasoListar          *tramites.ListarExpediente
	CasoBuscar          *tramites.BuscarEnExpediente
	CasoUrl             *tramites.UrlExpediente
	CasoBaja            *tramites.DarDeBajaExpediente
}

func (h *ExpedienteHandler) Listar(c *fiber.Ctx) error {
	actor, _ := middleware.DeContexto(c)
	docs, err := h.CasoListar.Ejecutar(c.Context(), c.Params("id"), actor.TienePermiso(rol.PermisoDocumentoRestituir))
	if err != nil {
		return apierror.Mapear(err)
	}
	out := make([]dto.DocumentoDTO, len(docs))
	for i, d := range docs {
		out[i] = aDocumentoDTO(d)
	}
	return c.JSON(out)
}

func (h *ExpedienteHandler) Buscar(c *fiber.Ctx) error {
	lista, err := h.CasoBuscar.Ejecutar(c.Context(), c.Params("id"), c.Query("q"))
	if err != nil {
		return apierror.Mapear(err)
	}
	out := dto.BusquedaExpedienteDTO{Coincidencias: make([]dto.CoincidenciaExpedienteDTO, len(lista))}
	for i, m := range lista {
		out.Coincidencias[i] = dto.CoincidenciaExpedienteDTO{
			DocumentoID: m.DocumentoID, TipoDocumento: string(m.TipoDocumento), NombreOriginal: m.NombreOriginal,
			Pagina: m.Pagina, Fragmento: m.Fragmento, Rango: m.Rango,
		}
	}
	return c.JSON(out)
}

func (h *ExpedienteHandler) Url(c *fiber.Ctx) error {
	actor, _ := middleware.DeContexto(c)
	url, err := h.CasoUrl.Ejecutar(c.Context(), c.Params("id"), c.Params("docId"), actor.TienePermiso(rol.PermisoDocumentoRestituir))
	if err != nil {
		return apierror.Mapear(err)
	}
	return c.JSON(fiber.Map{"url": url})
}

func (h *ExpedienteHandler) PrepararSubida(c *fiber.Ctx) error {
	var in dto.SubidaExpedienteEntrada
	if err := leerCuerpo(c, &in); err != nil {
		return err
	}
	r, err := h.CasoPrepararSubida.Ejecutar(c.Context(), tramites.SubidaDocenteEntrada{
		DocenteID: c.Params("id"), Nombre: in.Nombre, TipoMime: in.TipoMime, Tamano: in.Tamano, Tipo: tramite.TipoDocumento(in.TipoDocumento),
	})
	if err != nil {
		return apierror.Mapear(err)
	}
	return c.Status(fiber.StatusCreated).JSON(dto.SubidaPreparadaDTO{Documento: aDocumentoDTO(r.Documento), URL: r.URL})
}

func (h *ExpedienteHandler) ConfirmarSubida(c *fiber.Ctx) error {
	actor, _ := middleware.DeContexto(c)
	d, err := h.CasoConfirmarSubida.EjecutarDocente(c.Context(), actor.UsuarioID, c.Params("id"), c.Params("docId"))
	if err != nil {
		return apierror.Mapear(err)
	}
	return c.JSON(aDocumentoDTO(*d))
}

func (h *ExpedienteHandler) DarDeBaja(c *fiber.Ctx) error {
	var in dto.BajaDocumentoEntrada
	if err := leerCuerpo(c, &in); err != nil {
		return err
	}
	actor, _ := middleware.DeContexto(c)
	if err := h.CasoBaja.Ejecutar(c.Context(), actor.UsuarioID, c.Params("id"), c.Params("docId"), in.Motivo); err != nil {
		return apierror.Mapear(err)
	}
	return c.SendStatus(fiber.StatusNoContent)
}
