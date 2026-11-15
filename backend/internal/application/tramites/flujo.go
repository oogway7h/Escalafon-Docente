package tramites

import (
	"context"
	"fmt"
	"strings"
	"time"

	"escalafon/backend/internal/application/ports"
	"escalafon/backend/internal/domain/auditoria"
	"escalafon/backend/internal/domain/tramite"
)

// Etapas del trámite que siguen al dictamen del analista: salida a la facultad (tramo sin seguimiento),
// revisión de la resolución recibida y remisión al Rectorado. El archivo lo hace Archivo (Archivar).

// RegistrarPeriodo guarda los años a evaluar, que se verifican en el sistema durante la revisión.
type RegistrarPeriodo struct {
	Tramites  tramite.Repository
	Auditoria auditoria.Repository
	Reloj     ports.Clock
}

func (c *RegistrarPeriodo) Ejecutar(ctx context.Context, actorID, id string, desde, hasta int16) error {
	t, err := cargarDelAnalista(ctx, c.Tramites, actorID, id)
	if err != nil {
		return err
	}
	if t.Estado != tramite.EnRevision && t.Estado != tramite.Observado {
		return tramite.ErrTransicionInvalida
	}
	ahora := c.Reloj.Ahora()
	if desde < 1950 || hasta < desde || int(hasta) > ahora.Year() {
		return tramite.ErrPeriodoInvalido
	}
	t.AnioDesde, t.AnioHasta = desde, hasta
	h := &tramite.Historial{Estado: t.Estado, EstadoAnterior: t.Estado, UsuarioID: &actorID, Fecha: ahora,
		Detalle: fmt.Sprintf("Años a evaluar: %d a %d", desde, hasta)}
	if err := c.Tramites.Guardar(ctx, t, h); err != nil {
		return err
	}
	registrar(ctx, c.Auditoria, &actorID, auditoria.AccionRegistrarPeriodo, t.ID, ahora, map[string]any{"desde": desde, "hasta": hasta})
	return nil
}

// RemitirAFacultad registra la salida del expediente. Desde aquí el trámite no se puede rastrear
// hasta que regrese con la resolución.
type RemitirAFacultad struct {
	Tramites    tramite.Repository
	Auditoria   auditoria.Repository
	Reloj       ports.Clock
	Notificador Notificador
}

func (c *RemitirAFacultad) Ejecutar(ctx context.Context, actorID, id string) error {
	t, err := cargarDelAnalista(ctx, c.Tramites, actorID, id)
	if err != nil {
		return err
	}
	if err := tramite.Validar(t.Estado, tramite.EnFacultad); err != nil {
		return err
	}
	ahora := c.Reloj.Ahora()
	h := &tramite.Historial{Estado: tramite.EnFacultad, EstadoAnterior: t.Estado, UsuarioID: &actorID, Fecha: ahora,
		Detalle: "Remitido a la facultad (Tribunal, Decano y Consejo Facultativo): sin seguimiento hasta su regreso"}
	t.Estado = tramite.EnFacultad
	if err := c.Tramites.Guardar(ctx, t, h); err != nil {
		return err
	}
	registrar(ctx, c.Auditoria, &actorID, auditoria.AccionRemitirFacultad, t.ID, ahora, nil)
	notificar(ctx, c.Notificador, t)
	return nil
}

// RegistrarRegreso recibe la carpeta con la resolución de categorización. Lo registra el analista
// asignado o el personal de ventanilla; el trámite vuelve al mismo analista para la revisión final.
type RegistrarRegreso struct {
	Tramites  tramite.Repository
	Auditoria auditoria.Repository
	Reloj     ports.Clock
}

type RegresoEntrada struct {
	ActorID          string
	EsVentanilla     bool // el actor puede radicar trámites (ventanilla)
	TramiteID        string
	ResolucionNumero string
	ResolucionFecha  time.Time
}

func (c *RegistrarRegreso) Ejecutar(ctx context.Context, in RegresoEntrada) error {
	t, err := c.Tramites.BuscarPorID(ctx, in.TramiteID)
	if err != nil {
		return err
	}
	if t.Estado.Cerrado() {
		return t.Estado.ErrCerrado()
	}
	if !in.EsVentanilla && (t.AnalistaID == nil || *t.AnalistaID != in.ActorID) {
		return tramite.ErrNoAsignado
	}
	if err := tramite.Validar(t.Estado, tramite.RevisionFinal); err != nil {
		return err
	}
	numero := strings.TrimSpace(in.ResolucionNumero)
	ahora := c.Reloj.Ahora()
	if numero == "" || in.ResolucionFecha.IsZero() || in.ResolucionFecha.After(ahora) {
		return tramite.ErrResolucionRequerida
	}
	h := &tramite.Historial{Estado: tramite.RevisionFinal, EstadoAnterior: t.Estado, UsuarioID: &in.ActorID, Fecha: ahora,
		Detalle: fmt.Sprintf("Regresó de la facultad con la resolución %s (%s)", numero, in.ResolucionFecha.Format("02/01/2006"))}
	t.Estado, t.ResolucionNumero, t.ResolucionFecha = tramite.RevisionFinal, numero, &in.ResolucionFecha
	if err := c.Tramites.Guardar(ctx, t, h); err != nil {
		return err
	}
	registrar(ctx, c.Auditoria, &in.ActorID, auditoria.AccionRegresoFacultad, t.ID, ahora, map[string]any{"resolucion": numero})
	return nil
}

// RemitirARectorado cierra la revisión final: el expediente sale hacia el Rectorado (con los
// ejemplares de la categorización) y queda a la espera de que Archivo archive la carpeta.
type RemitirARectorado struct {
	Tramites    tramite.Repository
	Auditoria   auditoria.Repository
	Reloj       ports.Clock
	Notificador Notificador
}

func (c *RemitirARectorado) Ejecutar(ctx context.Context, actorID, id, nota string) error {
	t, err := cargarDelAnalista(ctx, c.Tramites, actorID, id)
	if err != nil {
		return err
	}
	if err := tramite.Validar(t.Estado, tramite.Remitido); err != nil {
		return err
	}
	nota = strings.TrimSpace(nota)
	if nota == "" {
		return tramite.ErrNotaRequerida
	}
	ahora := c.Reloj.Ahora()
	h := &tramite.Historial{Estado: tramite.Remitido, EstadoAnterior: t.Estado, UsuarioID: &actorID, Fecha: ahora,
		Detalle: "Remitido al Rectorado: " + nota}
	t.Estado, t.RemisionNota = tramite.Remitido, nota
	if err := c.Tramites.Guardar(ctx, t, h); err != nil {
		return err
	}
	registrar(ctx, c.Auditoria, &actorID, auditoria.AccionRemitirRectorado, t.ID, ahora, map[string]any{"nota": nota})
	notificar(ctx, c.Notificador, t)
	return nil
}
