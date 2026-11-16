package persistence

import (
	"context"
	"errors"

	"gorm.io/gorm"

	"escalafon/backend/internal/domain/categoria"
)

type RepositorioCategoria struct {
	DB *gorm.DB
}

func NewRepositorioCategoria(dbase *gorm.DB) *RepositorioCategoria {
	return &RepositorioCategoria{DB: dbase}
}

func aDominioCategoria(m *CategoriaModel, reqs []CategoriaRequisitoModel) categoria.Categoria {
	c := categoria.Categoria{ID: m.ID, Nombre: m.Nombre, AniosAntiguedad: m.AniosAntiguedad, PuntajeDesde: m.PuntajeDesde, PuntajeHasta: m.PuntajeHasta, Activa: m.Activa}
	for _, r := range reqs {
		c.Requisitos = append(c.Requisitos, categoria.Requisito{ID: r.ID, Nombre: r.Nombre})
	}
	return c
}

func (r *RepositorioCategoria) requisitosDe(ctx context.Context, ids ...int16) (map[int16][]CategoriaRequisitoModel, error) {
	var reqs []CategoriaRequisitoModel
	q := db(ctx, r.DB).Order("id")
	if len(ids) > 0 {
		q = q.Where("categoria_id IN ?", ids)
	}
	if err := q.Find(&reqs).Error; err != nil {
		return nil, err
	}
	out := map[int16][]CategoriaRequisitoModel{}
	for _, rq := range reqs {
		out[rq.CategoriaID] = append(out[rq.CategoriaID], rq)
	}
	return out, nil
}

// usos cuenta, por categoría, los docentes y trámites que la referencian.
func (r *RepositorioCategoria) usos(ctx context.Context) (docentes, tramites map[int16]int64, err error) {
	contar := func(tabla string) (map[int16]int64, error) {
		var filas []struct {
			CategoriaID int16
			N           int64
		}
		if err := db(ctx, r.DB).Table(tabla).Select("categoria_id, count(*) AS n").
			Where("categoria_id IS NOT NULL").Group("categoria_id").Scan(&filas).Error; err != nil {
			return nil, err
		}
		out := map[int16]int64{}
		for _, f := range filas {
			out[f.CategoriaID] = f.N
		}
		return out, nil
	}
	if docentes, err = contar("docente"); err != nil {
		return nil, nil, err
	}
	tramites, err = contar("tramite")
	return docentes, tramites, err
}

func (r *RepositorioCategoria) Listar(ctx context.Context, soloActivas bool) ([]categoria.Categoria, error) {
	var modelos []CategoriaModel
	q := db(ctx, r.DB).Order("id")
	if soloActivas {
		q = q.Where("activa")
	}
	if err := q.Find(&modelos).Error; err != nil {
		return nil, err
	}
	reqs, err := r.requisitosDe(ctx)
	if err != nil {
		return nil, err
	}
	docentes, tramites, err := r.usos(ctx)
	if err != nil {
		return nil, err
	}
	out := make([]categoria.Categoria, len(modelos))
	for i := range modelos {
		out[i] = aDominioCategoria(&modelos[i], reqs[modelos[i].ID])
		out[i].Docentes, out[i].Tramites = docentes[modelos[i].ID], tramites[modelos[i].ID]
	}
	return out, nil
}

func (r *RepositorioCategoria) BuscarPorID(ctx context.Context, id int16) (*categoria.Categoria, error) {
	var m CategoriaModel
	if err := db(ctx, r.DB).First(&m, "id = ?", id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, categoria.ErrNoEncontrado
		}
		return nil, err
	}
	reqs, err := r.requisitosDe(ctx, id)
	if err != nil {
		return nil, err
	}
	c := aDominioCategoria(&m, reqs[id])
	docentes, tramites, err := r.usos(ctx)
	if err != nil {
		return nil, err
	}
	c.Docentes, c.Tramites = docentes[id], tramites[id]
	return &c, nil
}

func (r *RepositorioCategoria) Guardar(ctx context.Context, c *categoria.Categoria, requisitos []string) error {
	return db(ctx, r.DB).Transaction(func(tx *gorm.DB) error {
		m := CategoriaModel{ID: c.ID, Nombre: c.Nombre, AniosAntiguedad: c.AniosAntiguedad, PuntajeDesde: c.PuntajeDesde, PuntajeHasta: c.PuntajeHasta, Activa: c.Activa}
		if c.ID == 0 {
			if err := tx.Create(&m).Error; err != nil {
				if esViolacionUnicidad(err) {
					return categoria.ErrDuplicada
				}
				return err
			}
			c.ID = m.ID
		} else {
			res := tx.Model(&CategoriaModel{}).Where("id = ?", c.ID).
				Updates(map[string]any{"nombre": c.Nombre, "anios_antiguedad": c.AniosAntiguedad, "puntaje_desde": c.PuntajeDesde, "puntaje_hasta": c.PuntajeHasta})
			if res.Error != nil {
				if esViolacionUnicidad(res.Error) {
					return categoria.ErrDuplicada
				}
				return res.Error
			}
			if res.RowsAffected == 0 {
				return categoria.ErrNoEncontrado
			}
		}
		// Conserva los requisitos existentes por nombre (los trámites los referencian) y
		// solo crea los nuevos o elimina los retirados.
		var existentes []CategoriaRequisitoModel
		if err := tx.Where("categoria_id = ?", c.ID).Find(&existentes).Error; err != nil {
			return err
		}
		deseados := map[string]bool{}
		for _, n := range requisitos {
			deseados[n] = true
		}
		tenidos := map[string]bool{}
		for _, e := range existentes {
			tenidos[e.Nombre] = true
			if !deseados[e.Nombre] {
				if err := tx.Delete(&e).Error; err != nil {
					return err
				}
			}
		}
		for n := range deseados {
			if !tenidos[n] {
				if err := tx.Create(&CategoriaRequisitoModel{CategoriaID: c.ID, Nombre: n}).Error; err != nil {
					return err
				}
			}
		}
		return nil
	})
}

func (r *RepositorioCategoria) Eliminar(ctx context.Context, id int16) error {
	res := db(ctx, r.DB).Delete(&CategoriaModel{}, "id = ?", id)
	if res.Error != nil {
		if esViolacionFK(res.Error) { // alguien la referenció entre la lectura y el borrado
			return &categoria.ErrEnUso{}
		}
		return res.Error
	}
	if res.RowsAffected == 0 {
		return categoria.ErrNoEncontrado
	}
	return nil
}

func (r *RepositorioCategoria) CambiarActiva(ctx context.Context, id int16, activa bool) error {
	res := db(ctx, r.DB).Model(&CategoriaModel{}).Where("id = ?", id).Update("activa", activa)
	if res.Error != nil {
		return res.Error
	}
	if res.RowsAffected == 0 {
		return categoria.ErrNoEncontrado
	}
	return nil
}
