// Package notificaciones implementa CU-14: avisar por correo al docente cuando su trámite pasa a
// Observado, Dictaminado o Archivado. Los envíos fallidos se registran y se reintentan.
package notificaciones

import (
	"context"
	"log/slog"
	"net/url"
	"strings"

	"escalafon/backend/internal/application/ports"
	"escalafon/backend/internal/domain/docente"
	"escalafon/backend/internal/domain/notificacion"
	"escalafon/backend/internal/domain/tramite"
)

// Notificar crea el aviso de una transición y lo envía. Nunca devuelve error: un fallo de correo
// no debe deshacer la transición del trámite.
type Notificar struct {
	Docentes    docente.Repository
	Envios      notificacion.Repository
	Mailer      ports.Mailer
	Reloj       ports.Clock
	URLFrontend string
	// Asincrono envía en segundo plano para no retener la petición HTTP mientras responde el SMTP.
	Asincrono bool
}

// Transicion avisa de que t acaba de pasar a su estado actual; ignora los estados sin aviso.
func (n *Notificar) Transicion(ctx context.Context, t *tramite.Tramite) {
	switch t.Estado {
	case tramite.Observado, tramite.EnFacultad, tramite.Remitido, tramite.Archivado:
	default:
		return
	}
	d, err := n.Docentes.BuscarPorID(ctx, t.DocenteID)
	if err != nil {
		slog.Warn("notificación: no se pudo cargar el docente", "tramite", t.Codigo, "error", err)
		return
	}
	e := &notificacion.Envio{
		TramiteID: t.ID, TramiteCodigo: t.Codigo, EstadoTramite: string(t.Estado),
		Destino: strings.TrimSpace(d.Correo), Nombres: d.Nombres, Estado: notificacion.Pendiente, CreadoEn: n.Reloj.Ahora(),
	}
	if e.Destino == "" { // sin correo no hay nada que reintentar
		e.Estado, e.Error = notificacion.Fallido, "el docente no tiene correo registrado"
	}
	if err := n.Envios.Crear(ctx, e); err != nil {
		slog.Warn("notificación: no se pudo registrar el aviso", "tramite", t.Codigo, "error", err)
		return
	}
	if e.Estado != notificacion.Pendiente {
		return
	}
	if n.Asincrono {
		go n.enviar(context.WithoutCancel(ctx), e)
		return
	}
	n.enviar(ctx, e)
}

func (n *Notificar) enviar(ctx context.Context, e *notificacion.Envio) {
	enlace := strings.TrimRight(n.URLFrontend, "/") + "/seguimiento?codigo=" + url.QueryEscape(e.TramiteCodigo)
	err := n.Mailer.EnviarCambioEstado(e.Destino, e.Nombres, e.TramiteCodigo, e.EstadoTramite, enlace)
	e.Intentos++
	if err == nil {
		ahora := n.Reloj.Ahora()
		e.Estado, e.EnviadoEn, e.Error = notificacion.Enviado, &ahora, ""
	} else {
		e.Error = err.Error()
		if e.Intentos >= notificacion.MaxIntentos {
			e.Estado = notificacion.Fallido
		}
	}
	if err := n.Envios.Actualizar(ctx, e); err != nil {
		slog.Warn("notificación: no se pudo actualizar el aviso", "tramite", e.TramiteCodigo, "error", err)
	}
}

// ReintentarPendientes reenvía los avisos que aún tienen intentos disponibles.
func (n *Notificar) ReintentarPendientes(ctx context.Context) error {
	pendientes, err := n.Envios.ListarPendientes(ctx, 50)
	if err != nil {
		return err
	}
	for i := range pendientes {
		n.enviar(ctx, &pendientes[i])
	}
	return nil
}

// Listar devuelve una página de avisos para la consulta del Director.
type Listar struct{ Envios notificacion.Repository }

func (c *Listar) Ejecutar(ctx context.Context, f notificacion.Filtro, pagina, tamano int) ([]notificacion.Envio, int64, error) {
	if pagina < 1 {
		pagina = 1
	}
	if tamano < 1 || tamano > 100 {
		tamano = 20
	}
	return c.Envios.Listar(ctx, f, pagina, tamano)
}
