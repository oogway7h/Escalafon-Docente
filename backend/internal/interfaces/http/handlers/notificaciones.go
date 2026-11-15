package handlers

import (
	"github.com/gofiber/fiber/v2"

	"escalafon/backend/internal/application/notificaciones"
	"escalafon/backend/internal/domain/notificacion"
	"escalafon/backend/internal/interfaces/http/apierror"
	"escalafon/backend/internal/interfaces/http/dto"
)

// NotificacionesHandler expone el registro de avisos de cambio de estado (CU-14).
type NotificacionesHandler struct {
	CasoListar *notificaciones.Listar
}

func (h *NotificacionesHandler) Listar(c *fiber.Ctx) error {
	pagina, tamano := paginacion(c)
	f := notificacion.Filtro{Estado: notificacion.Estado(c.Query("estado"))}
	lista, total, err := h.CasoListar.Ejecutar(c.Context(), f, pagina, tamano)
	if err != nil {
		return apierror.Mapear(err)
	}
	out := dto.ListaNotificacionesDTO{Datos: make([]dto.NotificacionDTO, len(lista)), Total: total, Pagina: pagina, Tamano: tamano}
	for i, e := range lista {
		out.Datos[i] = dto.NotificacionDTO{
			ID: e.ID, TramiteCodigo: e.TramiteCodigo, EstadoTramite: e.EstadoTramite, Destino: e.Destino,
			Estado: string(e.Estado), Intentos: e.Intentos, Error: e.Error, CreadoEn: e.CreadoEn, EnviadoEn: e.EnviadoEn,
		}
	}
	return c.JSON(out)
}
