// Package categorias implementa los casos de uso de CU-07 (catálogo normativo de categorías).
package categorias

import (
	"context"
	"strconv"
	"strings"

	"escalafon/backend/internal/application/ports"
	"escalafon/backend/internal/domain/auditoria"
	"escalafon/backend/internal/domain/categoria"
)

type Listar struct{ Categorias categoria.Repository }

func (c *Listar) Ejecutar(ctx context.Context, soloActivas bool) ([]categoria.Categoria, error) {
	return c.Categorias.Listar(ctx, soloActivas)
}

type Obtener struct{ Categorias categoria.Repository }

func (c *Obtener) Ejecutar(ctx context.Context, id int16) (*categoria.Categoria, error) {
	return c.Categorias.BuscarPorID(ctx, id)
}

type Entrada struct {
	ActorID         string
	ID              int16 // 0 = crear
	Nombre          string
	AniosAntiguedad int16
	Requisitos      []string
}

// Guardar crea o actualiza una categoría y reemplaza sus documentos exigidos.
type Guardar struct {
	Categorias categoria.Repository
	Auditoria  auditoria.Repository
	Reloj      ports.Clock
}

func (c *Guardar) Ejecutar(ctx context.Context, in Entrada) (*categoria.Categoria, error) {
	reqs := make([]string, 0, len(in.Requisitos))
	vistos := map[string]bool{}
	for _, r := range in.Requisitos {
		r = strings.TrimSpace(r)
		if r != "" && !vistos[r] {
			vistos[r] = true
			reqs = append(reqs, r)
		}
	}
	cat := &categoria.Categoria{ID: in.ID, Nombre: strings.TrimSpace(in.Nombre), AniosAntiguedad: in.AniosAntiguedad, Activa: true}
	if err := c.Categorias.Guardar(ctx, cat, reqs); err != nil {
		return nil, err
	}
	accion := auditoria.AccionCrear
	if in.ID != 0 {
		accion = auditoria.AccionEditar
	}
	_ = c.Auditoria.Registrar(ctx, &auditoria.Evento{
		UsuarioID: &in.ActorID, Accion: accion, Entidad: auditoria.EntidadCategoria,
		EntidadID: strconv.Itoa(int(cat.ID)), Fecha: c.Reloj.Ahora(),
	})
	return c.Categorias.BuscarPorID(ctx, cat.ID)
}

// Eliminar borra definitivamente una categoría. Si docentes o trámites la usan, se rechaza:
// en ese caso lo correcto es darla de baja.
type Eliminar struct {
	Categorias categoria.Repository
	Auditoria  auditoria.Repository
	Reloj      ports.Clock
}

func (c *Eliminar) Ejecutar(ctx context.Context, actorID string, id int16) error {
	cat, err := c.Categorias.BuscarPorID(ctx, id)
	if err != nil {
		return err
	}
	if cat.Docentes+cat.Tramites > 0 {
		return &categoria.ErrEnUso{Docentes: cat.Docentes, Tramites: cat.Tramites}
	}
	if err := c.Categorias.Eliminar(ctx, id); err != nil {
		return err
	}
	_ = c.Auditoria.Registrar(ctx, &auditoria.Evento{
		UsuarioID: &actorID, Accion: auditoria.AccionEliminarCategoria, Entidad: auditoria.EntidadCategoria,
		EntidadID: strconv.Itoa(int(id)), Detalle: map[string]any{"nombre": cat.Nombre}, Fecha: c.Reloj.Ahora(),
	})
	return nil
}

type CambiarEstado struct {
	Categorias categoria.Repository
	Auditoria  auditoria.Repository
	Reloj      ports.Clock
}

func (c *CambiarEstado) Ejecutar(ctx context.Context, actorID string, id int16, activa bool) error {
	if err := c.Categorias.CambiarActiva(ctx, id, activa); err != nil {
		return err
	}
	accion := auditoria.AccionBaja
	if activa {
		accion = auditoria.AccionReactivar
	}
	_ = c.Auditoria.Registrar(ctx, &auditoria.Evento{
		UsuarioID: &actorID, Accion: accion, Entidad: auditoria.EntidadCategoria,
		EntidadID: strconv.Itoa(int(id)), Fecha: c.Reloj.Ahora(),
	})
	return nil
}
