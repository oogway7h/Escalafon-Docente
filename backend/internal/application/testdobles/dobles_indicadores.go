package testdobles

import (
	"context"

	"escalafon/backend/internal/domain/indicador"
)

// IndicadorRepoFalso devuelve lo sembrado y aplica los mismos filtros que el repositorio real.
type IndicadorRepoFalso struct {
	Conteos    []indicador.ConteoEstado
	Resueltos  []indicador.Resuelto
	Analistas  []indicador.CargaAnalista
	Prestamos  []indicador.Prestamo
	Tramites   []indicador.FilaTramite
	Inventario []indicador.FilaInventario
}

func (r *IndicadorRepoFalso) ContarPorEstado(context.Context, indicador.Periodo) ([]indicador.ConteoEstado, error) {
	return r.Conteos, nil
}

func (r *IndicadorRepoFalso) ListarResueltos(_ context.Context, f indicador.FiltroReporte) ([]indicador.Resuelto, error) {
	var out []indicador.Resuelto
	for _, x := range r.Resueltos {
		if (f.Desde != nil && x.ResueltoEn.Before(*f.Desde)) || (f.Hasta != nil && !x.ResueltoEn.Before(*f.Hasta)) {
			continue
		}
		if (f.Tipo != "" && x.Tipo != f.Tipo) || (f.CategoriaID > 0 && x.CategoriaID != f.CategoriaID) {
			continue
		}
		out = append(out, x)
	}
	return out, nil
}

func (r *IndicadorRepoFalso) CargaAnalistas(context.Context, indicador.Periodo) ([]indicador.CargaAnalista, error) {
	return r.Analistas, nil
}

func (r *IndicadorRepoFalso) ListarPrestamos(context.Context) ([]indicador.Prestamo, error) {
	return r.Prestamos, nil
}

func (r *IndicadorRepoFalso) ListarTramites(_ context.Context, f indicador.FiltroReporte) ([]indicador.FilaTramite, error) {
	var out []indicador.FilaTramite
	for _, x := range r.Tramites {
		if (f.Desde != nil && x.RadicadoEn.Before(*f.Desde)) || (f.Hasta != nil && !x.RadicadoEn.Before(*f.Hasta)) {
			continue
		}
		if (f.Estado != "" && x.Estado != f.Estado) || (f.Tipo != "" && x.Tipo != f.Tipo) {
			continue
		}
		out = append(out, x)
	}
	return out, nil
}

func (r *IndicadorRepoFalso) ListarInventario(_ context.Context, f indicador.FiltroReporte) ([]indicador.FilaInventario, error) {
	var out []indicador.FilaInventario
	for _, x := range r.Inventario {
		if f.EstadoLegajo == "" || x.Estado == f.EstadoLegajo {
			out = append(out, x)
		}
	}
	return out, nil
}

// ExportadorFalso guarda el último reporte recibido y devuelve un contenido reconocible.
type ExportadorFalso struct {
	Ultimo *indicador.Reporte
}

func (e *ExportadorFalso) PDF(r *indicador.Reporte) ([]byte, error) {
	e.Ultimo = r
	return []byte("%PDF-falso"), nil
}

func (e *ExportadorFalso) Excel(r *indicador.Reporte) ([]byte, error) {
	e.Ultimo = r
	return []byte("PK-falso"), nil
}
