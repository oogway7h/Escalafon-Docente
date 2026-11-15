// Package auditorias implementa CU-06: Consultar registro de auditoría (solo lectura).
package auditorias

import (
	"context"

	"escalafon/backend/internal/domain/auditoria"
)

const (
	tamanoPorDefecto = 20
	tamanoMaximo     = 100
)

type Resultado struct {
	Registros []auditoria.Registro
	Total     int64
	Pagina    int
	Tamano    int
}

// Consultar lista la bitácora con filtros combinables, del evento más reciente al más antiguo.
type Consultar struct{ Auditoria auditoria.Consulta }

func (c *Consultar) Ejecutar(ctx context.Context, f auditoria.Filtro) (*Resultado, error) {
	if f.Accion != "" && !auditoria.AccionValida(f.Accion) {
		return nil, auditoria.ErrAccionInvalida
	}
	if f.Desde != nil && f.Hasta != nil && !f.Desde.Before(*f.Hasta) {
		return nil, auditoria.ErrRangoFechas
	}
	if f.Pagina < 1 {
		f.Pagina = 1
	}
	if f.Tamano < 1 || f.Tamano > tamanoMaximo {
		f.Tamano = tamanoPorDefecto
	}
	lista, total, err := c.Auditoria.Listar(ctx, f)
	if err != nil {
		return nil, err
	}
	return &Resultado{Registros: lista, Total: total, Pagina: f.Pagina, Tamano: f.Tamano}, nil
}

// Opciones son los valores posibles de los filtros de la pantalla de auditoría.
type Opciones struct {
	Acciones []string
	Usuarios []auditoria.Autor
}

type ListarOpciones struct{ Auditoria auditoria.Consulta }

func (c *ListarOpciones) Ejecutar(ctx context.Context) (*Opciones, error) {
	autores, err := c.Auditoria.ListarAutores(ctx)
	if err != nil {
		return nil, err
	}
	return &Opciones{Acciones: auditoria.Acciones, Usuarios: autores}, nil
}
