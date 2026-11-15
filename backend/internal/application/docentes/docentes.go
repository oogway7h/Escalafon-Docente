// Package docentes implementa los casos de uso de CU-07 (padrón de docentes).
package docentes

import (
	"context"
	"strings"

	"escalafon/backend/internal/application/ports"
	"escalafon/backend/internal/domain/auditoria"
	"escalafon/backend/internal/domain/categoria"
	"escalafon/backend/internal/domain/docente"
)

type Listar struct{ Docentes docente.Repository }

func (c *Listar) Ejecutar(ctx context.Context, busqueda string, soloActivos bool) ([]docente.Docente, error) {
	return c.Docentes.Listar(ctx, busqueda, soloActivos)
}

// ListarPagina devuelve una página de docentes (tamaño acotado a 1..100) y el total de coincidencias.
// Sugerir alimenta el buscador con sugerencias mientras se escribe (como la barra de un navegador).
type Sugerir struct{ Docentes docente.Repository }

func (c *Sugerir) Ejecutar(ctx context.Context, texto string) ([]docente.Docente, error) {
	if len([]rune(strings.TrimSpace(texto))) < 2 {
		return []docente.Docente{}, nil
	}
	return c.Docentes.Sugerir(ctx, texto, 8)
}

type ListarPagina struct{ Docentes docente.Repository }

func (c *ListarPagina) Ejecutar(ctx context.Context, busqueda string, soloActivos bool, pagina, tamano int) ([]docente.Docente, int64, error) {
	if pagina < 1 {
		pagina = 1
	}
	if tamano < 1 || tamano > 100 {
		tamano = 10
	}
	return c.Docentes.ListarPagina(ctx, busqueda, soloActivos, pagina, tamano)
}

type Obtener struct{ Docentes docente.Repository }

func (c *Obtener) Ejecutar(ctx context.Context, id string) (*docente.Docente, error) {
	return c.Docentes.BuscarPorID(ctx, id)
}

type ListarFacultades struct{ Docentes docente.Repository }

func (c *ListarFacultades) Ejecutar(ctx context.Context) ([]docente.Facultad, error) {
	return c.Docentes.ListarFacultades(ctx)
}

type Entrada struct {
	ActorID     string
	ID          string // vacío al crear
	Codigo      string // código universitario del docente
	Cedula      string // opcional; no se puede cambiar una vez asignada
	Nombres     string
	Apellidos   string
	Correo      string
	Telefono    string
	CategoriaID *int16
	FacultadID  *int16
}

// Guardar crea o edita un docente. La cédula, una vez asignada, no cambia.
type Guardar struct {
	Docentes   docente.Repository
	Categorias categoria.Repository
	Auditoria  auditoria.Repository
	Reloj      ports.Clock
}

func (c *Guardar) Ejecutar(ctx context.Context, in Entrada) (*docente.Docente, error) {
	if in.CategoriaID != nil {
		cat, err := c.Categorias.BuscarPorID(ctx, *in.CategoriaID)
		if err != nil || !cat.Activa {
			return nil, docente.ErrCategoriaInvalida
		}
	}
	d := &docente.Docente{
		ID: in.ID, Codigo: strings.TrimSpace(in.Codigo), Cedula: strings.TrimSpace(in.Cedula), Nombres: strings.TrimSpace(in.Nombres),
		Apellidos: strings.TrimSpace(in.Apellidos), Correo: strings.TrimSpace(in.Correo),
		Telefono: strings.TrimSpace(in.Telefono), CategoriaID: in.CategoriaID, FacultadID: in.FacultadID,
	}
	accion := auditoria.AccionCrear
	if in.ID == "" {
		if err := c.Docentes.Crear(ctx, d); err != nil {
			return nil, err
		}
	} else {
		existente, err := c.Docentes.BuscarPorID(ctx, in.ID)
		if err != nil {
			return nil, err
		}
		if existente.Cedula != "" {
			d.Cedula = existente.Cedula
		}
		if d.Codigo == "" {
			d.Codigo = existente.Codigo
		}
		if err := c.Docentes.Actualizar(ctx, d); err != nil {
			return nil, err
		}
		accion = auditoria.AccionEditar
	}
	_ = c.Auditoria.Registrar(ctx, &auditoria.Evento{
		UsuarioID: &in.ActorID, Accion: accion, Entidad: auditoria.EntidadDocente, EntidadID: d.ID, Fecha: c.Reloj.Ahora(),
	})
	return c.Docentes.BuscarPorID(ctx, d.ID)
}

// CambiarEstado da de baja (o reactiva) lógicamente a un docente.
type CambiarEstado struct {
	Docentes  docente.Repository
	Auditoria auditoria.Repository
	Reloj     ports.Clock
}

func (c *CambiarEstado) Ejecutar(ctx context.Context, actorID, id string, activo bool) error {
	if err := c.Docentes.CambiarActivo(ctx, id, activo); err != nil {
		return err
	}
	accion := auditoria.AccionBaja
	if activo {
		accion = auditoria.AccionReactivar
	}
	_ = c.Auditoria.Registrar(ctx, &auditoria.Evento{
		UsuarioID: &actorID, Accion: accion, Entidad: auditoria.EntidadDocente, EntidadID: id, Fecha: c.Reloj.Ahora(),
	})
	return nil
}
