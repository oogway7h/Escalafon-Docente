package persistence

import (
	"context"
	"errors"
	"strings"
	"time"

	"gorm.io/gorm"

	"escalafon/backend/internal/domain/tramite"
)

type RepositorioTramite struct {
	DB *gorm.DB
}

func NewRepositorioTramite(dbase *gorm.DB) *RepositorioTramite { return &RepositorioTramite{DB: dbase} }

// delTramite limita la consulta a los documentos de un trámite; tramiteID "" alcanza los escaneos del
// expediente digital del docente (los que no tienen trámite).
func delTramite(tramiteID string) func(*gorm.DB) *gorm.DB {
	return func(q *gorm.DB) *gorm.DB {
		if tramiteID == "" {
			return q.Where("tramite_id IS NULL")
		}
		return q.Where("tramite_id = ?", tramiteID)
	}
}

func ptr(s string) *string {
	if s == "" {
		return nil
	}
	return &s
}

func val(s *string) string {
	if s == nil {
		return ""
	}
	return *s
}

func nombreCompleto(nombres, apellidos string) string {
	return strings.TrimSpace(nombres + " " + apellidos)
}

func (r *RepositorioTramite) Crear(ctx context.Context, t *tramite.Tramite) error {
	return db(ctx, r.DB).Transaction(func(tx *gorm.DB) error {
		m := TramiteModel{
			ID: t.ID, Codigo: t.Codigo, DocenteID: t.DocenteID,
			Tipo: string(t.Tipo), Estado: string(t.Estado), Origen: string(t.Origen),
			PendienteValidacion: t.PendienteValidacion, AnalistaID: t.AnalistaID, RadicadoPor: t.RadicadoPor,
		}
		if t.CategoriaID != 0 { // 0 = sin categoría (regularización)
			m.CategoriaID = &t.CategoriaID
		}
		if err := tx.Create(&m).Error; err != nil {
			if esViolacionUnicidad(err) {
				return tramite.ErrCodigoDuplicado
			}
			return err
		}
		t.CreadoEn = m.CreadoEn
		for _, d := range t.Documentos {
			dm := TramiteDocumentoModel{
				ID: d.ID, TramiteID: ptr(t.ID), RequisitoID: d.RequisitoID, NombreOriginal: d.NombreOriginal,
				TipoMime: d.TipoMime, Tamano: d.Tamano, HashSHA256: ptr(d.HashSHA256), Bucket: d.Bucket, ClaveObjeto: d.ClaveObjeto,
				Activo: true, EstadoOCR: string(estadoOCRInicial(d.EstadoOCR)),
			}
			if err := tx.Create(&dm).Error; err != nil {
				return err
			}
		}
		for _, h := range t.Historial {
			if err := anexarHistorial(tx, t.ID, &h); err != nil {
				return err
			}
		}
		return nil
	})
}

func anexarHistorial(tx *gorm.DB, tramiteID string, h *tramite.Historial) error {
	hm := TramiteHistorialModel{
		TramiteID: tramiteID, Estado: string(h.Estado), EstadoAnterior: ptr(string(h.EstadoAnterior)),
		UsuarioID: h.UsuarioID, Detalle: ptr(h.Detalle), Fecha: h.Fecha,
	}
	return tx.Create(&hm).Error
}

func (r *RepositorioTramite) BuscarPorID(ctx context.Context, id string) (*tramite.Tramite, error) {
	return r.buscar(ctx, "tramite.id = ?", id)
}

func (r *RepositorioTramite) BuscarPorCodigo(ctx context.Context, codigo string) (*tramite.Tramite, error) {
	return r.buscar(ctx, "tramite.codigo = ?", codigo)
}

func (r *RepositorioTramite) buscar(ctx context.Context, condicion string, valor any) (*tramite.Tramite, error) {
	var m TramiteModel
	if err := db(ctx, r.DB).Preload("Docente").Preload("Categoria").Preload("Analista").
		First(&m, condicion, valor).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, tramite.ErrNoEncontrado
		}
		return nil, err
	}
	t := &tramite.Tramite{
		ID: m.ID, Codigo: m.Codigo, DocenteID: m.DocenteID, DocenteCodigo: val(m.Docente.Codigo),
		DocenteNombre: nombreCompleto(m.Docente.Nombres, m.Docente.Apellidos),
		Tipo: tramite.Tipo(m.Tipo),
		Estado: tramite.Estado(m.Estado), Origen: tramite.Origen(m.Origen), PendienteValidacion: m.PendienteValidacion,
		AnalistaID: m.AnalistaID, RadicadoPor: m.RadicadoPor, Resultado: tramite.Resultado(val(m.Resultado)),
		Dictamen: val(m.Dictamen), DictaminadoEn: m.DictaminadoEn, ArchivadoEn: m.ArchivadoEn,
		UbicacionLegajo: val(m.UbicacionLegajo), CreadoEn: m.CreadoEn,
		ResolucionNumero: val(m.ResolucionNumero), ResolucionFecha: m.ResolucionFecha, RemisionNota: val(m.RemisionNota),
	}
	if m.CategoriaID != nil && m.Categoria != nil {
		t.CategoriaID, t.Categoria = *m.CategoriaID, m.Categoria.Nombre
	}
	if m.AnioDesde != nil {
		t.AnioDesde = *m.AnioDesde
	}
	if m.AnioHasta != nil {
		t.AnioHasta = *m.AnioHasta
	}
	var carpetas struct{ Ligadas, Abiertas int64 }
	if err := db(ctx, r.DB).Raw(`SELECT count(*) FILTER (WHERE estado <> 'RECHAZADA') AS ligadas,
		count(*) FILTER (WHERE estado IN ('PENDIENTE','ATENDIDA','DESPACHADA')) AS abiertas
		FROM solicitud_carpeta WHERE tramite_id = ?`, m.ID).Scan(&carpetas).Error; err != nil {
		return nil, err
	}
	t.CarpetasLigadas, t.CarpetasAbiertas = carpetas.Ligadas, carpetas.Abiertas
	if m.Analista != nil {
		t.AnalistaNombre = nombreCompleto(m.Analista.Nombres, m.Analista.Apellidos)
	}

	var hist []struct {
		TramiteHistorialModel
		Nombres   *string
		Apellidos *string
	}
	if err := db(ctx, r.DB).Table("tramite_historial h").
		Select("h.*, u.nombres, u.apellidos").
		Joins("LEFT JOIN usuario u ON u.id = h.usuario_id").
		Where("h.tramite_id = ?", m.ID).Order("h.fecha, h.id").Scan(&hist).Error; err != nil {
		return nil, err
	}
	for _, h := range hist {
		e := tramite.Historial{
			ID: h.ID, Estado: tramite.Estado(h.Estado), EstadoAnterior: tramite.Estado(val(h.EstadoAnterior)),
			UsuarioID: h.UsuarioID, Detalle: val(h.Detalle), Fecha: h.Fecha,
		}
		if h.Nombres != nil {
			e.UsuarioNombre = nombreCompleto(*h.Nombres, val(h.Apellidos))
		}
		t.Historial = append(t.Historial, e)
	}

	var obs []ObservacionModel
	if err := db(ctx, r.DB).Where("tramite_id = ?", m.ID).Order("creado_en").Find(&obs).Error; err != nil {
		return nil, err
	}
	for _, o := range obs {
		t.Observaciones = append(t.Observaciones, tramite.Observacion{
			ID: o.ID, Documento: o.Documento, Motivo: o.Motivo, Plazo: o.Plazo,
			Resuelta: o.Resuelta, ResueltaEn: o.ResueltaEn, CreadoEn: o.CreadoEn,
		})
	}

	var docs []TramiteDocumentoModel
	if err := db(ctx, r.DB).Where("tramite_id = ?", m.ID).Order("creado_en").Find(&docs).Error; err != nil {
		return nil, err
	}
	presentes := map[int32]bool{}
	for _, d := range docs {
		t.Documentos = append(t.Documentos, aDominioDocumento(&d))
		// Un documento dado de baja (CU-16) deja de cubrir el requisito de la categoría.
		if d.RequisitoID != nil && d.Activo {
			presentes[*d.RequisitoID] = true
		}
	}

	var reqs []CategoriaRequisitoModel // sin categoría (regularización) no hay requisitos: t.CategoriaID es 0
	if err := db(ctx, r.DB).Where("categoria_id = ?", t.CategoriaID).Order("id").Find(&reqs).Error; err != nil {
		return nil, err
	}
	for _, rq := range reqs {
		t.Requisitos = append(t.Requisitos, tramite.Requisito{ID: rq.ID, Nombre: rq.Nombre, Presente: presentes[rq.ID]})
	}
	return t, nil
}

func aDominioDocumento(d *TramiteDocumentoModel) tramite.Documento {
	doc := tramite.Documento{
		ID: d.ID, DocenteID: val(d.DocenteID), TipoDocumento: tramite.TipoDocumento(val(d.TipoDocumento)),
		RequisitoID: d.RequisitoID, NombreOriginal: d.NombreOriginal, TipoMime: d.TipoMime,
		Tamano: d.Tamano, HashSHA256: strings.TrimSpace(val(d.HashSHA256)), Bucket: d.Bucket, ClaveObjeto: d.ClaveObjeto,
		Activo: d.Activo, MotivoBaja: val(d.MotivoBaja), BajaEn: d.BajaEn, BajaPor: d.BajaPor,
		EstadoOCR: tramite.EstadoOCR(d.EstadoOCR), TextoOCR: val(d.TextoOCR), TextoValidado: d.TextoValidado,
		ErrorOCR: val(d.ErrorOCR), ClaveComprimida: val(d.ClaveComprimida),
	}
	if d.TamanoComprimido != nil {
		doc.TamanoComprimido = *d.TamanoComprimido
	}
	return doc
}

func estadoOCRInicial(e tramite.EstadoOCR) tramite.EstadoOCR {
	if e == "" {
		return tramite.OCRNoAplica
	}
	return e
}

// consultaListado arma el filtro común de Listar y ListarPagina.
func (r *RepositorioTramite) consultaListado(ctx context.Context, f tramite.Filtro) *gorm.DB {
	q := db(ctx, r.DB).Table("tramite t").
		Joins("JOIN docente d ON d.id = t.docente_id").
		Joins("LEFT JOIN categoria c ON c.id = t.categoria_id"). // las regularizaciones no tienen categoría
		Joins("LEFT JOIN usuario a ON a.id = t.analista_id")
	if f.Estado != "" {
		q = q.Where("t.estado = ?", string(f.Estado))
	}
	if f.AnalistaID != "" {
		q = q.Where("t.analista_id = ?", f.AnalistaID)
	}
	if f.SinAsignar {
		q = q.Where("t.analista_id IS NULL AND t.estado NOT IN ('REMITIDO','ARCHIVADO')")
	}
	if b := strings.ToLower(strings.TrimSpace(f.Busqueda)); b != "" {
		// Código de docente y cédula coinciden exactos; el código de trámite, por fragmento; el nombre solo si no hay dígitos.
		if tieneDigitos(b) {
			q = q.Where("lower(t.codigo) LIKE ? OR lower(d.codigo) = ? OR lower(d.cedula) = ?", "%"+escaparLike(b)+"%", b, b)
		} else {
			q = q.Where("lower(d.apellidos || ' ' || d.nombres) LIKE ?", "%"+escaparLike(b)+"%")
		}
	}
	return q
}

func (r *RepositorioTramite) Listar(ctx context.Context, f tramite.Filtro) ([]tramite.Resumen, error) {
	return resumenes(r.consultaListado(ctx, f))
}

func (r *RepositorioTramite) ListarPorDocente(ctx context.Context, docenteID string) ([]tramite.Resumen, error) {
	return resumenes(r.consultaListado(ctx, tramite.Filtro{}).Where("t.docente_id = ?", docenteID))
}

func (r *RepositorioTramite) ListarPagina(ctx context.Context, f tramite.Filtro, pagina, tamano int) ([]tramite.Resumen, int64, error) {
	var total int64
	if err := r.consultaListado(ctx, f).Count(&total).Error; err != nil {
		return nil, 0, err
	}
	out, err := resumenes(r.consultaListado(ctx, f).Offset((pagina - 1) * tamano).Limit(tamano))
	return out, total, err
}

func resumenes(q *gorm.DB) ([]tramite.Resumen, error) {
	var filas []struct {
		ID, Codigo, Tipo, Estado, Origen, Nombres, Apellidos, Categoria string
		PendienteValidacion                                             bool
		CreadoEn                                                        time.Time
		ANombres, AApellidos, DCodigo                                   *string
	}
	err := q.Select("t.id, t.codigo, t.tipo, t.estado, t.origen, t.pendiente_validacion, t.creado_en, " +
		"d.codigo AS d_codigo, d.nombres, d.apellidos, COALESCE(c.nombre, '') AS categoria, a.nombres AS a_nombres, a.apellidos AS a_apellidos").
		Order("t.creado_en DESC, t.id").Scan(&filas).Error
	if err != nil {
		return nil, err
	}
	out := make([]tramite.Resumen, len(filas))
	for i, f := range filas {
		out[i] = tramite.Resumen{
			ID: f.ID, Codigo: f.Codigo, DocenteCodigo: val(f.DCodigo), DocenteNombre: nombreCompleto(f.Nombres, f.Apellidos),
			Categoria: f.Categoria, Tipo: tramite.Tipo(f.Tipo), Estado: tramite.Estado(f.Estado), Origen: tramite.Origen(f.Origen),
			PendienteValidacion: f.PendienteValidacion, CreadoEn: f.CreadoEn,
		}
		if f.ANombres != nil {
			out[i].AnalistaNombre = nombreCompleto(*f.ANombres, val(f.AApellidos))
		}
	}
	return out, nil
}

func anio(a int16) *int16 {
	if a == 0 {
		return nil
	}
	return &a
}

func (r *RepositorioTramite) Guardar(ctx context.Context, t *tramite.Tramite, h *tramite.Historial) error {
	return db(ctx, r.DB).Transaction(func(tx *gorm.DB) error {
		res := tx.Model(&TramiteModel{}).Where("id = ?", t.ID).Updates(map[string]any{
			"estado": string(t.Estado), "analista_id": t.AnalistaID, "pendiente_validacion": t.PendienteValidacion,
			"resultado": ptr(string(t.Resultado)), "dictamen": ptr(t.Dictamen), "dictaminado_en": t.DictaminadoEn,
			"archivado_en": t.ArchivadoEn, "ubicacion_legajo": ptr(t.UbicacionLegajo), "actualizado_en": gorm.Expr("now()"),
			"anio_desde": anio(t.AnioDesde), "anio_hasta": anio(t.AnioHasta), "resolucion_numero": ptr(t.ResolucionNumero),
			"resolucion_fecha": t.ResolucionFecha, "remision_nota": ptr(t.RemisionNota),
		})
		if res.Error != nil {
			return res.Error
		}
		if res.RowsAffected == 0 {
			return tramite.ErrNoEncontrado
		}
		if h != nil {
			return anexarHistorial(tx, t.ID, h)
		}
		return nil
	})
}

func (r *RepositorioTramite) AgregarObservacion(ctx context.Context, tramiteID string, o *tramite.Observacion, creadaPor string) error {
	m := ObservacionModel{TramiteID: tramiteID, Documento: o.Documento, Motivo: o.Motivo, Plazo: o.Plazo, CreadaPor: creadaPor}
	if err := db(ctx, r.DB).Create(&m).Error; err != nil {
		return err
	}
	o.ID, o.CreadoEn = m.ID, m.CreadoEn
	return nil
}

func (r *RepositorioTramite) BuscarObservacion(ctx context.Context, tramiteID, id string) (*tramite.Observacion, error) {
	var m ObservacionModel
	if err := db(ctx, r.DB).First(&m, "id = ? AND tramite_id = ?", id, tramiteID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, tramite.ErrObservacionNoEncontrada
		}
		return nil, err
	}
	return &tramite.Observacion{
		ID: m.ID, Documento: m.Documento, Motivo: m.Motivo, Plazo: m.Plazo,
		Resuelta: m.Resuelta, ResueltaEn: m.ResueltaEn, CreadoEn: m.CreadoEn,
	}, nil
}

func (r *RepositorioTramite) ResolverObservacion(ctx context.Context, id string, en time.Time) error {
	return db(ctx, r.DB).Model(&ObservacionModel{}).Where("id = ?", id).
		Updates(map[string]any{"resuelta": true, "resuelta_en": en}).Error
}

func (r *RepositorioTramite) ContarObservacionesPendientes(ctx context.Context, tramiteID string) (int64, error) {
	var n int64
	err := db(ctx, r.DB).Model(&ObservacionModel{}).Where("tramite_id = ? AND NOT resuelta", tramiteID).Count(&n).Error
	return n, err
}

func (r *RepositorioTramite) BuscarDocumento(ctx context.Context, tramiteID, id string) (*tramite.Documento, error) {
	var m TramiteDocumentoModel
	if err := db(ctx, r.DB).Scopes(delTramite(tramiteID)).First(&m, "id = ?", id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, tramite.ErrNoEncontrado
		}
		return nil, err
	}
	d := aDominioDocumento(&m)
	return &d, nil
}

// DarDeBajaDocumento aplica la baja lógica del adjunto (CU-16). El UPDATE es condicional: si el
// documento ya no estaba activo no afecta filas y devuelve ErrDocumentoDeBaja.
func (r *RepositorioTramite) DarDeBajaDocumento(ctx context.Context, tramiteID, id, motivo, actorID string, en time.Time) error {
	res := db(ctx, r.DB).Model(&TramiteDocumentoModel{}).Scopes(delTramite(tramiteID)).
		Where("id = ? AND activo", id).
		Updates(map[string]any{"activo": false, "motivo_baja": motivo, "baja_en": en, "baja_por": actorID})
	if res.Error != nil {
		return res.Error
	}
	if res.RowsAffected == 0 {
		return tramite.ErrDocumentoDeBaja
	}
	return nil
}

// CrearDocumento agrega un adjunto a un trámite ya radicado (CU-15).
func (r *RepositorioTramite) CrearDocumento(ctx context.Context, tramiteID string, d *tramite.Documento) error {
	m := TramiteDocumentoModel{
		ID: d.ID, TramiteID: ptr(tramiteID), RequisitoID: d.RequisitoID, NombreOriginal: d.NombreOriginal, TipoMime: d.TipoMime,
		Tamano: d.Tamano, HashSHA256: ptr(d.HashSHA256), Bucket: d.Bucket, ClaveObjeto: d.ClaveObjeto,
		Activo: true, EstadoOCR: string(estadoOCRInicial(d.EstadoOCR)),
	}
	return db(ctx, r.DB).Create(&m).Error
}

// CrearDocumentoDocente agrega un escaneo al expediente digital del docente (sin trámite).
func (r *RepositorioTramite) CrearDocumentoDocente(ctx context.Context, docenteID string, d *tramite.Documento) error {
	m := TramiteDocumentoModel{
		ID: d.ID, DocenteID: &docenteID, TipoDocumento: ptr(string(d.TipoDocumento)), NombreOriginal: d.NombreOriginal,
		TipoMime: d.TipoMime, Tamano: d.Tamano, HashSHA256: ptr(d.HashSHA256), Bucket: d.Bucket, ClaveObjeto: d.ClaveObjeto,
		Activo: true, EstadoOCR: string(estadoOCRInicial(d.EstadoOCR)),
	}
	return db(ctx, r.DB).Create(&m).Error
}

func (r *RepositorioTramite) ListarDocumentosDocente(ctx context.Context, docenteID string, incluirBajas bool) ([]tramite.Documento, error) {
	var ms []TramiteDocumentoModel
	q := db(ctx, r.DB).Where("docente_id = ? AND tramite_id IS NULL", docenteID)
	if !incluirBajas {
		q = q.Where("activo")
	}
	if err := q.Order("tipo_documento, creado_en").Find(&ms).Error; err != nil {
		return nil, err
	}
	out := make([]tramite.Documento, len(ms))
	for i := range ms {
		out[i] = aDominioDocumento(&ms[i])
	}
	return out, nil
}

// BuscarEnDocente busca en las páginas de todos los documentos activos y procesados del docente:
// evita revisar el expediente página por página. Lo más relevante primero.
func (r *RepositorioTramite) BuscarEnDocente(ctx context.Context, docenteID, consulta string) ([]tramite.CoincidenciaDocente, error) {
	var filas []struct {
		DocumentoID    string
		TipoDocumento  string
		NombreOriginal string
		Pagina         int
		Fragmento      string
		Rango          float64
	}
	err := db(ctx, r.DB).Raw(`
		SELECT d.id AS documento_id, d.tipo_documento, d.nombre_original, p.pagina,
		       ts_headline('es_unaccent', p.texto, q.q, 'StartSel=⟦, StopSel=⟧, MinWords=6, MaxWords=26, MaxFragments=1, ShortWord=2') AS fragmento,
		       ts_rank(to_tsvector('es_unaccent', p.texto), q.q) AS rango
		FROM documento_pagina p
		JOIN tramite_documento d ON d.id = p.documento_id,
		     (SELECT plainto_tsquery('es_unaccent', ?) AS q) q
		WHERE d.docente_id = ? AND d.tramite_id IS NULL AND d.activo AND d.estado_ocr = 'PROCESADO'
		  AND to_tsvector('es_unaccent', p.texto) @@ q.q
		ORDER BY rango DESC, d.tipo_documento, p.pagina
		LIMIT 100`, consulta, docenteID).Scan(&filas).Error
	if err != nil {
		return nil, err
	}
	out := make([]tramite.CoincidenciaDocente, len(filas))
	for i, f := range filas {
		out[i] = tramite.CoincidenciaDocente{
			DocumentoID: f.DocumentoID, TipoDocumento: tramite.TipoDocumento(f.TipoDocumento), NombreOriginal: f.NombreOriginal,
			Coincidencia: tramite.Coincidencia{Pagina: f.Pagina, Fragmento: f.Fragmento, Rango: f.Rango},
		}
	}
	return out, nil
}

// GuardarDocumento actualiza los campos que cambian durante la carga y el procesamiento (CU-15).
func (r *RepositorioTramite) GuardarDocumento(ctx context.Context, tramiteID string, d *tramite.Documento) error {
	var comprimido *int64
	if d.ClaveComprimida != "" {
		comprimido = &d.TamanoComprimido
	}
	res := db(ctx, r.DB).Model(&TramiteDocumentoModel{}).Scopes(delTramite(tramiteID)).Where("id = ?", d.ID).
		Updates(map[string]any{
			"tamano": d.Tamano, "hash_sha256": ptr(d.HashSHA256), "estado_ocr": string(d.EstadoOCR), "texto_ocr": ptr(d.TextoOCR),
			"texto_validado": d.TextoValidado, "error_ocr": ptr(d.ErrorOCR), "clave_comprimida": ptr(d.ClaveComprimida),
			"tamano_comprimido": comprimido,
		})
	if res.Error != nil {
		return res.Error
	}
	if res.RowsAffected == 0 {
		return tramite.ErrDocumentoNoEncontrado
	}
	return nil
}

// documentoDelTramite comprueba que el documento pertenezca al trámite.
func (r *RepositorioTramite) documentoDelTramite(ctx context.Context, tramiteID, docID string) error {
	var n int64
	if err := db(ctx, r.DB).Model(&TramiteDocumentoModel{}).Scopes(delTramite(tramiteID)).Where("id = ?", docID).Count(&n).Error; err != nil {
		return err
	}
	if n == 0 {
		return tramite.ErrDocumentoNoEncontrado
	}
	return nil
}

func (r *RepositorioTramite) GuardarPaginas(ctx context.Context, tramiteID, docID string, paginas []string) error {
	if err := r.documentoDelTramite(ctx, tramiteID, docID); err != nil {
		return err
	}
	return db(ctx, r.DB).Transaction(func(tx *gorm.DB) error {
		if err := tx.Where("documento_id = ?", docID).Delete(&DocumentoPaginaModel{}).Error; err != nil {
			return err
		}
		filas := make([]DocumentoPaginaModel, len(paginas))
		for i, p := range paginas {
			filas[i] = DocumentoPaginaModel{DocumentoID: docID, Pagina: i + 1, Texto: p}
		}
		if len(filas) > 0 {
			if err := tx.CreateInBatches(filas, 200).Error; err != nil {
				return err
			}
		}
		return tx.Model(&TramiteDocumentoModel{}).Where("id = ?", docID).Update("texto_ocr", strings.Join(paginas, "\n\n")).Error
	})
}

func (r *RepositorioTramite) ListarPaginas(ctx context.Context, tramiteID, docID string) ([]string, error) {
	if err := r.documentoDelTramite(ctx, tramiteID, docID); err != nil {
		return nil, err
	}
	var filas []DocumentoPaginaModel
	if err := db(ctx, r.DB).Where("documento_id = ?", docID).Order("pagina").Find(&filas).Error; err != nil {
		return nil, err
	}
	out := make([]string, len(filas))
	for i, f := range filas {
		out[i] = f.Texto
	}
	return out, nil
}

// BuscarEnDocumento usa el índice de texto completo por página (español: ignora tildes, mayúsculas y
// terminaciones). ts_headline marca las palabras halladas con ⟦ ⟧ para resaltarlas en el visor.
func (r *RepositorioTramite) BuscarEnDocumento(ctx context.Context, tramiteID, docID, consulta string) ([]tramite.Coincidencia, error) {
	if err := r.documentoDelTramite(ctx, tramiteID, docID); err != nil {
		return nil, err
	}
	var filas []struct {
		Pagina    int
		Fragmento string
		Rango     float64
	}
	err := db(ctx, r.DB).Raw(`
		SELECT p.pagina,
		       ts_headline('es_unaccent', p.texto, q.q, 'StartSel=⟦, StopSel=⟧, MinWords=6, MaxWords=26, MaxFragments=1, ShortWord=2') AS fragmento,
		       ts_rank(to_tsvector('es_unaccent', p.texto), q.q) AS rango
		FROM documento_pagina p, (SELECT plainto_tsquery('es_unaccent', ?) AS q) q
		WHERE p.documento_id = ? AND to_tsvector('es_unaccent', p.texto) @@ q.q
		ORDER BY p.pagina`, consulta, docID).Scan(&filas).Error
	if err != nil {
		return nil, err
	}
	out := make([]tramite.Coincidencia, len(filas))
	for i, f := range filas {
		out[i] = tramite.Coincidencia{Pagina: f.Pagina, Fragmento: f.Fragmento, Rango: f.Rango}
	}
	return out, nil
}

func (r *RepositorioTramite) ListarDocumentosEnProceso(ctx context.Context) ([]tramite.DocumentoEnProceso, error) {
	var ms []TramiteDocumentoModel
	if err := db(ctx, r.DB).Where("estado_ocr IN ('PENDIENTE','PROCESANDO') AND activo").Find(&ms).Error; err != nil {
		return nil, err
	}
	out := make([]tramite.DocumentoEnProceso, len(ms))
	for i := range ms {
		out[i] = tramite.DocumentoEnProceso{TramiteID: val(ms[i].TramiteID), Documento: aDominioDocumento(&ms[i])}
	}
	return out, nil
}

// analistas devuelve los usuarios activos con el permiso TRAMITE_REVISAR.
func (r *RepositorioTramite) analistas(ctx context.Context) *gorm.DB {
	return db(ctx, r.DB).Table("usuario u").
		Joins("JOIN rol_permiso rp ON rp.rol_id = u.rol_id").
		Joins("JOIN permiso p ON p.id = rp.permiso_id AND p.codigo = 'TRAMITE_REVISAR'").
		Where("u.activo")
}

func (r *RepositorioTramite) ListarAnalistas(ctx context.Context) ([]tramite.Analista, error) {
	var filas []struct {
		ID, Nombres, Apellidos string
		Carga                  int64
	}
	err := r.analistas(ctx).
		Select("u.id, u.nombres, u.apellidos, (SELECT count(*) FROM tramite t WHERE t.analista_id = u.id AND t.estado IN ('ASIGNADO','EN_REVISION','OBSERVADO','DICTAMINADO','REVISION_FINAL')) AS carga").
		Order("u.apellidos, u.nombres").Scan(&filas).Error
	if err != nil {
		return nil, err
	}
	out := make([]tramite.Analista, len(filas))
	for i, f := range filas {
		out[i] = tramite.Analista{ID: f.ID, Nombre: nombreCompleto(f.Nombres, f.Apellidos), Carga: f.Carga}
	}
	return out, nil
}
