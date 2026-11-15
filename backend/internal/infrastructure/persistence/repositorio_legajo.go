package persistence

import (
	"context"
	"strings"
	"time"

	"gorm.io/gorm"

	"escalafon/backend/internal/domain/legajo"
)

type RepositorioLegajo struct {
	DB *gorm.DB
}

func NewRepositorioLegajo(dbase *gorm.DB) *RepositorioLegajo { return &RepositorioLegajo{DB: dbase} }

type filaLegajo struct {
	ID                          int32
	Codigo, Estado, DocenteID   string
	Gaveta, Anexos, Observacion *string
	Orden                       *int32
	DCodigo                     *string
	Nombres, Apellidos          string
}

func (f filaLegajo) dominio() legajo.Legajo {
	var orden int
	if f.Orden != nil {
		orden = int(*f.Orden)
	}
	return legajo.Legajo{
		ID: f.ID, Codigo: f.Codigo, DocenteID: f.DocenteID, DocenteCodigo: val(f.DCodigo),
		DocenteNombre: nombreCompleto(f.Nombres, f.Apellidos), Gaveta: val(f.Gaveta), Orden: orden,
		Anexos: val(f.Anexos), Observacion: val(f.Observacion), Estado: legajo.Estado(f.Estado),
	}
}

const selectLegajo = "l.id, l.codigo, l.estado, l.docente_id, l.gaveta, l.orden, l.anexos, l.observacion, " +
	"d.codigo AS d_codigo, d.nombres, d.apellidos"

func (r *RepositorioLegajo) consultaLegajos(ctx context.Context, f legajo.FiltroLegajos) *gorm.DB {
	q := db(ctx, r.DB).Table("legajo l").Joins("JOIN docente d ON d.id = l.docente_id")
	if f.Estado != "" {
		q = q.Where("l.estado = ?", string(f.Estado))
	}
	if b := strings.ToLower(strings.TrimSpace(f.Busqueda)); b != "" {
		// Códigos y cédula coinciden exactos; el nombre, por fragmento.
		if tieneDigitos(b) {
			q = q.Where("lower(l.codigo) = ? OR lower(d.codigo) = ? OR lower(d.cedula) = ?", b, b, b)
		} else {
			q = q.Where("lower(d.apellidos || ' ' || d.nombres) LIKE ?", "%"+escaparLike(b)+"%")
		}
	}
	return q
}

func (r *RepositorioLegajo) ListarLegajos(ctx context.Context, f legajo.FiltroLegajos, pagina, tamano int) ([]legajo.Legajo, int64, error) {
	var total int64
	if err := r.consultaLegajos(ctx, f).Count(&total).Error; err != nil {
		return nil, 0, err
	}
	var filas []filaLegajo
	err := r.consultaLegajos(ctx, f).Select(selectLegajo).Order("d.apellidos, d.nombres, l.gaveta, l.orden, l.id").
		Offset((pagina - 1) * tamano).Limit(tamano).Scan(&filas).Error
	if err != nil {
		return nil, 0, err
	}
	out := make([]legajo.Legajo, len(filas))
	for i, fila := range filas {
		out[i] = fila.dominio()
	}
	return out, total, nil
}

func (r *RepositorioLegajo) BuscarLegajo(ctx context.Context, id int32) (*legajo.Legajo, error) {
	var filas []filaLegajo
	err := db(ctx, r.DB).Table("legajo l").Joins("JOIN docente d ON d.id = l.docente_id").
		Select(selectLegajo).Where("l.id = ?", id).Scan(&filas).Error
	if err != nil {
		return nil, err
	}
	if len(filas) == 0 {
		return nil, legajo.ErrNoEncontrado
	}
	l := filas[0].dominio()
	return &l, nil
}

func errorDeLegajo(err error) error {
	switch {
	case esViolacionUnicidad(err):
		return legajo.ErrCodigoDuplicado
	case esViolacionFK(err):
		return legajo.ErrDocenteInvalido
	}
	return err
}

func (r *RepositorioLegajo) CrearLegajo(ctx context.Context, l *legajo.Legajo) error {
	orden := int32(l.Orden)
	m := LegajoModel{
		Codigo: l.Codigo, Estado: string(legajo.EnArchivo), DocenteID: l.DocenteID, Gaveta: ptr(l.Gaveta), Orden: &orden,
		Anexos: ptr(l.Anexos), Observacion: ptr(l.Observacion),
	}
	if err := db(ctx, r.DB).Create(&m).Error; err != nil {
		return errorDeLegajo(err)
	}
	l.ID = m.ID
	return nil
}

func (r *RepositorioLegajo) ActualizarLegajo(ctx context.Context, l *legajo.Legajo) error {
	cambios := map[string]any{"gaveta": ptr(l.Gaveta), "orden": l.Orden, "anexos": ptr(l.Anexos), "observacion": ptr(l.Observacion)}
	if l.Codigo != "" {
		cambios["codigo"] = l.Codigo
	}
	res := db(ctx, r.DB).Model(&LegajoModel{}).Where("id = ?", l.ID).Updates(cambios)
	if res.Error != nil {
		return errorDeLegajo(res.Error)
	}
	if res.RowsAffected == 0 {
		return legajo.ErrNoEncontrado
	}
	return nil
}

func (r *RepositorioLegajo) CambiarEstadoLegajo(ctx context.Context, id int32, desde, hacia legajo.Estado) (bool, error) {
	res := db(ctx, r.DB).Model(&LegajoModel{}).Where("id = ? AND estado = ?", id, string(desde)).Update("estado", string(hacia))
	return res.RowsAffected == 1, res.Error
}

// --- Solicitudes ------------------------------------------------------------------

func (r *RepositorioLegajo) CrearSolicitud(ctx context.Context, s *legajo.Solicitud) error {
	m := SolicitudCarpetaModel{LegajoID: s.LegajoID, SolicitanteID: s.SolicitanteID, TramiteID: ptr(s.TramiteID), Motivo: s.Motivo, Estado: string(s.Estado)}
	if err := db(ctx, r.DB).Create(&m).Error; err != nil {
		if esViolacionUnicidad(err) {
			return legajo.ErrSolicitudAbierta
		}
		return err
	}
	s.ID, s.CreadoEn = m.ID, m.CreadoEn
	return nil
}

type filaSolicitud struct {
	ID, Estado, Motivo, SolicitanteID, LegajoCodigo, Nombres, Apellidos, SNombres, SApellidos string
	LegajoID                                                                                  int32
	TramiteID, TramiteCodigo                                                                  *string
	Respuesta                                                                                 *string
	AtendidoEn                                                                                *time.Time
	AtendidoPor                                                                               *string
	CreadoEn                                                                                  time.Time
	FechaLimite                                                                               *time.Time
}

func (f filaSolicitud) dominio() legajo.Solicitud {
	return legajo.Solicitud{
		ID: f.ID, LegajoID: f.LegajoID, LegajoCodigo: f.LegajoCodigo, DocenteNombre: nombreCompleto(f.Nombres, f.Apellidos),
		TramiteID: val(f.TramiteID), TramiteCodigo: val(f.TramiteCodigo),
		SolicitanteID: f.SolicitanteID, Solicitante: nombreCompleto(f.SNombres, f.SApellidos), Motivo: f.Motivo,
		Estado: legajo.EstadoSolicitud(f.Estado), Respuesta: val(f.Respuesta), AtendidoPor: val(f.AtendidoPor),
		CreadoEn: f.CreadoEn, AtendidoEn: f.AtendidoEn, FechaLimite: f.FechaLimite,
	}
}

func (r *RepositorioLegajo) consultaSolicitudes(ctx context.Context) *gorm.DB {
	return db(ctx, r.DB).Table("solicitud_carpeta s").
		Joins("JOIN legajo l ON l.id = s.legajo_id").
		Joins("JOIN docente d ON d.id = l.docente_id").
		Joins("JOIN usuario u ON u.id = s.solicitante_id").
		Joins("LEFT JOIN tramite tr ON tr.id = s.tramite_id")
}

const selectSolicitud = "s.tramite_id, tr.codigo AS tramite_codigo, s.id, s.estado, s.motivo, s.solicitante_id, s.legajo_id, s.respuesta, s.atendido_por, s.atendido_en, s.creado_en, " +
	"l.codigo AS legajo_codigo, d.nombres, d.apellidos, u.nombres AS s_nombres, u.apellidos AS s_apellidos, " +
	"(SELECT m.fecha_limite FROM movimiento_custodia m WHERE m.solicitud_id = s.id AND m.tipo = 'SALIDA' LIMIT 1) AS fecha_limite"

func (r *RepositorioLegajo) BuscarSolicitud(ctx context.Context, id string) (*legajo.Solicitud, error) {
	var filas []filaSolicitud
	if err := r.consultaSolicitudes(ctx).Select(selectSolicitud).Where("s.id = ?", id).Scan(&filas).Error; err != nil {
		return nil, err
	}
	if len(filas) == 0 {
		return nil, legajo.ErrSolicitudNoExiste
	}
	s := filas[0].dominio()
	return &s, nil
}

func (r *RepositorioLegajo) ListarSolicitudes(ctx context.Context, f legajo.FiltroSolicitudes, pagina, tamano int) ([]legajo.Solicitud, int64, error) {
	filtrar := func(q *gorm.DB) *gorm.DB {
		if f.Estado != "" {
			q = q.Where("s.estado = ?", string(f.Estado))
		}
		if f.SolicitanteID != "" {
			q = q.Where("s.solicitante_id = ?", f.SolicitanteID)
		}
		return q
	}
	var total int64
	if err := filtrar(r.consultaSolicitudes(ctx)).Count(&total).Error; err != nil {
		return nil, 0, err
	}
	var filas []filaSolicitud
	err := filtrar(r.consultaSolicitudes(ctx)).Select(selectSolicitud).Order("s.creado_en DESC, s.id").
		Offset((pagina - 1) * tamano).Limit(tamano).Scan(&filas).Error
	if err != nil {
		return nil, 0, err
	}
	out := make([]legajo.Solicitud, len(filas))
	for i, fila := range filas {
		out[i] = fila.dominio()
	}
	return out, total, nil
}

func (r *RepositorioLegajo) HaySolicitudAbierta(ctx context.Context, legajoID int32) (bool, error) {
	var n int64
	err := db(ctx, r.DB).Model(&SolicitudCarpetaModel{}).
		Where("legajo_id = ? AND estado IN ('PENDIENTE','ATENDIDA','DESPACHADA')", legajoID).Count(&n).Error
	return n > 0, err
}

func (r *RepositorioLegajo) ActualizarSolicitud(ctx context.Context, id string, desde, hacia legajo.EstadoSolicitud, respuesta, usuarioID string, en time.Time) error {
	res := db(ctx, r.DB).Model(&SolicitudCarpetaModel{}).Where("id = ? AND estado = ?", id, string(desde)).
		Updates(map[string]any{"estado": string(hacia), "respuesta": ptr(respuesta), "atendido_por": usuarioID, "atendido_en": en})
	if res.Error != nil {
		return res.Error
	}
	if res.RowsAffected == 0 {
		var n int64
		if err := db(ctx, r.DB).Model(&SolicitudCarpetaModel{}).Where("id = ?", id).Count(&n).Error; err != nil {
			return err
		}
		if n == 0 {
			return legajo.ErrSolicitudNoExiste
		}
		return legajo.ErrEstadoSolicitud
	}
	return nil
}

func (r *RepositorioLegajo) SolicitudDespachada(ctx context.Context, legajoID int32) (*legajo.Solicitud, error) {
	var filas []filaSolicitud
	err := r.consultaSolicitudes(ctx).Select(selectSolicitud).
		Where("s.legajo_id = ? AND s.estado = 'DESPACHADA'", legajoID).Scan(&filas).Error
	if err != nil {
		return nil, err
	}
	if len(filas) == 0 {
		return nil, legajo.ErrNoPrestado
	}
	s := filas[0].dominio()
	return &s, nil
}

// --- Movimientos ------------------------------------------------------------------

func (r *RepositorioLegajo) AnexarMovimiento(ctx context.Context, m *legajo.Movimiento) error {
	mm := MovimientoCustodiaModel{
		LegajoID: m.LegajoID, SolicitudID: m.SolicitudID, Tipo: string(m.Tipo), CustodioID: m.CustodioID, FechaLimite: m.FechaLimite,
		EstadoMaterial: string(m.EstadoMaterial), Observacion: ptr(m.Observacion), RegistradoPor: m.RegistradoPorID, Fecha: m.Fecha,
	}
	if err := db(ctx, r.DB).Create(&mm).Error; err != nil {
		return err
	}
	m.ID = mm.ID
	return nil
}

func (r *RepositorioLegajo) ListarMovimientos(ctx context.Context, f legajo.FiltroMovimientos, pagina, tamano int) ([]legajo.Movimiento, int64, error) {
	consulta := func() *gorm.DB {
		q := db(ctx, r.DB).Table("movimiento_custodia m").
			Joins("JOIN legajo l ON l.id = m.legajo_id").
			Joins("JOIN docente d ON d.id = l.docente_id").
			Joins("JOIN usuario c ON c.id = m.custodio_id").
			Joins("JOIN usuario r ON r.id = m.registrado_por").
			Joins("JOIN solicitud_carpeta s ON s.id = m.solicitud_id")
		if f.Tipo != "" {
			q = q.Where("m.tipo = ?", string(f.Tipo))
		}
		if f.Vencidos {
			q = q.Where("m.tipo = 'SALIDA' AND s.estado = 'DESPACHADA' AND m.fecha_limite < CURRENT_DATE")
		}
		if b := strings.ToLower(strings.TrimSpace(f.Busqueda)); b != "" {
			if tieneDigitos(b) {
				q = q.Where("lower(l.codigo) = ? OR lower(d.codigo) = ? OR lower(d.cedula) = ?", b, b, b)
			} else {
				q = q.Where("lower(d.apellidos || ' ' || d.nombres) LIKE ?", "%"+escaparLike(b)+"%")
			}
		}
		return q
	}
	var total int64
	if err := consulta().Count(&total).Error; err != nil {
		return nil, 0, err
	}
	var filas []struct {
		ID                                                           int64
		LegajoID                                                     int32
		SolicitudID, Tipo, CustodioID, EstadoMaterial, RegistradoPor string
		LegajoCodigo, Nombres, Apellidos, CNombres, CApellidos       string
		RNombres, RApellidos                                         string
		FechaLimite                                                  *time.Time
		Observacion                                                  *string
		Fecha                                                        time.Time
		Vencido                                                      bool
	}
	err := consulta().Select("m.id, m.legajo_id, m.solicitud_id, m.tipo, m.custodio_id, m.estado_material, m.registrado_por, " +
		"m.fecha_limite, m.observacion, m.fecha, l.codigo AS legajo_codigo, d.nombres, d.apellidos, " +
		"c.nombres AS c_nombres, c.apellidos AS c_apellidos, r.nombres AS r_nombres, r.apellidos AS r_apellidos, " +
		"(m.tipo = 'SALIDA' AND s.estado = 'DESPACHADA' AND m.fecha_limite < CURRENT_DATE) AS vencido").
		Order("m.fecha DESC, m.id DESC").Offset((pagina - 1) * tamano).Limit(tamano).Scan(&filas).Error
	if err != nil {
		return nil, 0, err
	}
	out := make([]legajo.Movimiento, len(filas))
	for i, f := range filas {
		out[i] = legajo.Movimiento{
			ID: f.ID, LegajoID: f.LegajoID, LegajoCodigo: f.LegajoCodigo, DocenteNombre: nombreCompleto(f.Nombres, f.Apellidos),
			SolicitudID: f.SolicitudID, Tipo: legajo.TipoMovimiento(f.Tipo), CustodioID: f.CustodioID,
			Custodio: nombreCompleto(f.CNombres, f.CApellidos), FechaLimite: f.FechaLimite,
			EstadoMaterial: legajo.EstadoMaterial(f.EstadoMaterial), Observacion: val(f.Observacion),
			RegistradoPor: nombreCompleto(f.RNombres, f.RApellidos), RegistradoPorID: f.RegistradoPor, Fecha: f.Fecha, Vencido: f.Vencido,
		}
	}
	return out, total, nil
}
