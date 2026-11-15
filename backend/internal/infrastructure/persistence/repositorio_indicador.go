package persistence

import (
	"context"
	"fmt"
	"strings"
	"time"

	"gorm.io/gorm"

	"escalafon/backend/internal/domain/indicador"
	"escalafon/backend/internal/domain/legajo"
	"escalafon/backend/internal/domain/tramite"
)

// RepositorioIndicador implementa indicador.Consulta (CU-25 y CU-26): solo lee.
type RepositorioIndicador struct {
	DB *gorm.DB
}

func NewRepositorioIndicador(dbase *gorm.DB) *RepositorioIndicador {
	return &RepositorioIndicador{DB: dbase}
}

// resolucion da, por trámite, el primer paso a Remitido (o a Archivado en datos anteriores al flujo actual).
const resolucion = "(SELECT tramite_id, min(fecha) AS resuelto_en FROM tramite_historial " +
	"WHERE estado IN ('REMITIDO','ARCHIVADO') GROUP BY tramite_id)"

// enPeriodo filtra columna por el periodo (desde inclusivo, hasta exclusivo).
func enPeriodo(q *gorm.DB, columna string, p indicador.Periodo) *gorm.DB {
	if p.Desde != nil {
		q = q.Where(columna+" >= ?", *p.Desde)
	}
	if p.Hasta != nil {
		q = q.Where(columna+" < ?", *p.Hasta)
	}
	return q
}

func porTipoCategoria(q *gorm.DB, f indicador.FiltroReporte) *gorm.DB {
	if f.Tipo != "" {
		q = q.Where("t.tipo = ?", string(f.Tipo))
	}
	if f.CategoriaID > 0 {
		q = q.Where("t.categoria_id = ?", f.CategoriaID)
	}
	return q
}

func (r *RepositorioIndicador) ContarPorEstado(ctx context.Context, p indicador.Periodo) ([]indicador.ConteoEstado, error) {
	var filas []struct {
		Estado   string
		Cantidad int64
	}
	err := enPeriodo(db(ctx, r.DB).Table("tramite t"), "t.creado_en", p).
		Select("t.estado, count(*) AS cantidad").Group("t.estado").Scan(&filas).Error
	if err != nil {
		return nil, err
	}
	out := make([]indicador.ConteoEstado, len(filas))
	for i, f := range filas {
		out[i] = indicador.ConteoEstado{Estado: tramite.Estado(f.Estado), Cantidad: f.Cantidad}
	}
	return out, nil
}

func (r *RepositorioIndicador) ListarResueltos(ctx context.Context, f indicador.FiltroReporte) ([]indicador.Resuelto, error) {
	q := db(ctx, r.DB).Table("tramite t").
		Joins("JOIN " + resolucion + " r ON r.tramite_id = t.id").
		Joins("JOIN docente d ON d.id = t.docente_id").
		Joins("LEFT JOIN categoria c ON c.id = t.categoria_id").
		Joins("LEFT JOIN usuario a ON a.id = t.analista_id")
	q = porTipoCategoria(enPeriodo(q, "r.resuelto_en", f.Periodo), f)
	var filas []struct {
		Codigo, Tipo, Nombres, Apellidos, Categoria string
		CategoriaID                                 int16
		Resultado, DCodigo, ANombres, AApellidos    *string
		CreadoEn, ResueltoEn                        time.Time
	}
	err := q.Select("t.codigo, t.tipo, t.resultado, COALESCE(t.categoria_id, 0) AS categoria_id, t.creado_en, r.resuelto_en, d.codigo AS d_codigo, " +
		"d.nombres, d.apellidos, COALESCE(c.nombre, '') AS categoria, a.nombres AS a_nombres, a.apellidos AS a_apellidos").
		Order("r.resuelto_en DESC, t.codigo").Scan(&filas).Error
	if err != nil {
		return nil, err
	}
	out := make([]indicador.Resuelto, len(filas))
	for i, f := range filas {
		out[i] = indicador.Resuelto{
			Codigo: f.Codigo, DocenteCodigo: val(f.DCodigo), DocenteNombre: nombreCompleto(f.Nombres, f.Apellidos),
			Tipo: tramite.Tipo(f.Tipo), CategoriaID: f.CategoriaID, Categoria: f.Categoria,
			Analista: nombreCompleto(val(f.ANombres), val(f.AApellidos)), Resultado: tramite.Resultado(val(f.Resultado)),
			RadicadoEn: f.CreadoEn, ResueltoEn: f.ResueltoEn,
		}
	}
	return out, nil
}

func (r *RepositorioIndicador) CargaAnalistas(ctx context.Context, p indicador.Periodo) ([]indicador.CargaAnalista, error) {
	cond, args := []string{"r.resuelto_en IS NOT NULL"}, []any{}
	if p.Desde != nil {
		cond, args = append(cond, "r.resuelto_en >= ?"), append(args, *p.Desde)
	}
	if p.Hasta != nil {
		cond, args = append(cond, "r.resuelto_en < ?"), append(args, *p.Hasta)
	}
	// Mismos analistas que la bandeja de asignación: usuarios activos con TRAMITE_REVISAR.
	sql := fmt.Sprintf(`SELECT u.id, u.nombres, u.apellidos,
		count(t.id) FILTER (WHERE t.estado IN ('ASIGNADO','EN_REVISION','OBSERVADO','DICTAMINADO','REVISION_FINAL')) AS activos,
		count(t.id) FILTER (WHERE t.estado = 'EN_FACULTAD') AS en_facultad,
		count(t.id) FILTER (WHERE %s) AS resueltos
		FROM usuario u
		JOIN rol_permiso rp ON rp.rol_id = u.rol_id
		JOIN permiso p ON p.id = rp.permiso_id AND p.codigo = 'TRAMITE_REVISAR'
		LEFT JOIN tramite t ON t.analista_id = u.id
		LEFT JOIN %s r ON r.tramite_id = t.id
		WHERE u.activo
		GROUP BY u.id, u.nombres, u.apellidos
		ORDER BY u.apellidos, u.nombres`, strings.Join(cond, " AND "), resolucion)
	var filas []struct {
		ID, Nombres, Apellidos         string
		Activos, EnFacultad, Resueltos int64
	}
	if err := db(ctx, r.DB).Raw(sql, args...).Scan(&filas).Error; err != nil {
		return nil, err
	}
	out := make([]indicador.CargaAnalista, len(filas))
	for i, f := range filas {
		out[i] = indicador.CargaAnalista{
			ID: f.ID, Nombre: nombreCompleto(f.Nombres, f.Apellidos), Activos: f.Activos, EnFacultad: f.EnFacultad, Resueltos: f.Resueltos,
		}
	}
	return out, nil
}

func (r *RepositorioIndicador) ListarPrestamos(ctx context.Context) ([]indicador.Prestamo, error) {
	var filas []struct {
		LegajoID                                               int32
		LegajoCodigo, Nombres, Apellidos, CNombres, CApellidos string
		TramiteCodigo                                          *string
		Salida                                                 time.Time
		FechaLimite                                            *time.Time
	}
	err := db(ctx, r.DB).Table("solicitud_carpeta s").
		Joins("JOIN legajo l ON l.id = s.legajo_id").
		Joins("JOIN docente d ON d.id = l.docente_id").
		Joins("JOIN movimiento_custodia m ON m.solicitud_id = s.id AND m.tipo = 'SALIDA'").
		Joins("JOIN usuario c ON c.id = m.custodio_id").
		Joins("LEFT JOIN tramite tr ON tr.id = s.tramite_id").
		Where("s.estado = 'DESPACHADA'").
		Select("l.id AS legajo_id, l.codigo AS legajo_codigo, d.nombres, d.apellidos, c.nombres AS c_nombres, " +
			"c.apellidos AS c_apellidos, tr.codigo AS tramite_codigo, m.fecha AS salida, m.fecha_limite").
		Order("m.fecha, l.codigo").Scan(&filas).Error
	if err != nil {
		return nil, err
	}
	out := make([]indicador.Prestamo, len(filas))
	for i, f := range filas {
		out[i] = indicador.Prestamo{
			LegajoID: f.LegajoID, LegajoCodigo: f.LegajoCodigo, DocenteNombre: nombreCompleto(f.Nombres, f.Apellidos),
			Custodio: nombreCompleto(f.CNombres, f.CApellidos), TramiteCodigo: val(f.TramiteCodigo),
			Salida: f.Salida, FechaLimite: f.FechaLimite,
		}
	}
	return out, nil
}

func (r *RepositorioIndicador) ListarTramites(ctx context.Context, f indicador.FiltroReporte) ([]indicador.FilaTramite, error) {
	q := db(ctx, r.DB).Table("tramite t").
		Joins("JOIN docente d ON d.id = t.docente_id").
		Joins("LEFT JOIN categoria c ON c.id = t.categoria_id").
		Joins("LEFT JOIN usuario a ON a.id = t.analista_id")
	q = porTipoCategoria(enPeriodo(q, "t.creado_en", f.Periodo), f)
	if f.Estado != "" {
		q = q.Where("t.estado = ?", string(f.Estado))
	}
	var filas []struct {
		Codigo, Tipo, Estado, Origen, Nombres, Apellidos, Categoria string
		Resultado, DCodigo, ANombres, AApellidos                    *string
		CreadoEn                                                    time.Time
	}
	err := q.Select("t.codigo, t.tipo, t.estado, t.origen, t.resultado, t.creado_en, d.codigo AS d_codigo, d.nombres, " +
		"d.apellidos, COALESCE(c.nombre, '') AS categoria, a.nombres AS a_nombres, a.apellidos AS a_apellidos").
		Order("t.creado_en DESC, t.codigo").Scan(&filas).Error
	if err != nil {
		return nil, err
	}
	out := make([]indicador.FilaTramite, len(filas))
	for i, f := range filas {
		out[i] = indicador.FilaTramite{
			Codigo: f.Codigo, DocenteCodigo: val(f.DCodigo), DocenteNombre: nombreCompleto(f.Nombres, f.Apellidos),
			Tipo: tramite.Tipo(f.Tipo), Categoria: f.Categoria, Estado: tramite.Estado(f.Estado), Origen: tramite.Origen(f.Origen),
			Analista: nombreCompleto(val(f.ANombres), val(f.AApellidos)), Resultado: tramite.Resultado(val(f.Resultado)),
			RadicadoEn: f.CreadoEn,
		}
	}
	return out, nil
}

func (r *RepositorioIndicador) ListarInventario(ctx context.Context, f indicador.FiltroReporte) ([]indicador.FilaInventario, error) {
	// ux_solicitud_abierta garantiza a lo sumo una solicitud DESPACHADA (y su salida) por legajo.
	q := db(ctx, r.DB).Table("legajo l").
		Joins("JOIN docente d ON d.id = l.docente_id").
		Joins("LEFT JOIN solicitud_carpeta s ON s.legajo_id = l.id AND s.estado = 'DESPACHADA'").
		Joins("LEFT JOIN movimiento_custodia m ON m.solicitud_id = s.id AND m.tipo = 'SALIDA'").
		Joins("LEFT JOIN usuario c ON c.id = m.custodio_id")
	if f.EstadoLegajo != "" {
		q = q.Where("l.estado = ?", string(f.EstadoLegajo))
	}
	var filas []struct {
		LegajoCodigo, Estado, Nombres, Apellidos      string
		Gaveta, Anexos, DCodigo, CNombres, CApellidos *string
		Orden                                         *int32
		Salida, FechaLimite                           *time.Time
	}
	err := q.Select("l.codigo AS legajo_codigo, l.estado, l.gaveta, l.orden, l.anexos, d.codigo AS d_codigo, d.nombres, " +
		"d.apellidos, c.nombres AS c_nombres, c.apellidos AS c_apellidos, m.fecha AS salida, m.fecha_limite").
		// Orden físico y natural de las gavetas: A1, A2, …, A10 (no A1, A10, A2).
		Order("substring(l.gaveta from '^[^0-9]*') NULLS LAST, substring(l.gaveta from '[0-9]+')::int NULLS LAST, " +
			"l.gaveta, l.orden NULLS LAST, l.codigo").Scan(&filas).Error
	if err != nil {
		return nil, err
	}
	out := make([]indicador.FilaInventario, len(filas))
	for i, f := range filas {
		var orden int
		if f.Orden != nil {
			orden = int(*f.Orden)
		}
		out[i] = indicador.FilaInventario{
			LegajoCodigo: f.LegajoCodigo, DocenteCodigo: val(f.DCodigo), DocenteNombre: nombreCompleto(f.Nombres, f.Apellidos),
			Gaveta: val(f.Gaveta), Orden: orden, Anexos: val(f.Anexos), Estado: legajo.Estado(f.Estado),
			Custodio: nombreCompleto(val(f.CNombres), val(f.CApellidos)), Salida: f.Salida, FechaLimite: f.FechaLimite,
		}
	}
	return out, nil
}
