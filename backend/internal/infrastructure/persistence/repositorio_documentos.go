package persistence

import (
	"context"
	"strconv"
	"strings"

	"escalafon/backend/internal/domain/tramite"
)

// Consultas transversales sobre tramite_documento: abarcan a la vez los adjuntos de trámites y los
// escaneos del expediente del docente (CU-17 y CU-18). El docente sale del documento o de su trámite.
const fromDocumentos = `
	FROM tramite_documento d
	LEFT JOIN tramite t ON t.id = d.tramite_id
	JOIN docente k ON k.id = COALESCE(d.docente_id, t.docente_id)`

// condicionDocente filtra por docente con la misma regla que el padrón: con dígitos, código o cédula
// exactos; sin dígitos, un fragmento del nombre en cualquier orden (apellidos nombres o al revés).
func condicionDocente(b string) (string, map[string]any) {
	b = strings.ToLower(strings.TrimSpace(b))
	if tieneDigitos(b) {
		return "(lower(k.codigo) = @b OR lower(k.cedula) = @b)", map[string]any{"b": b}
	}
	return "(lower(k.apellidos || ' ' || k.nombres) LIKE @like OR lower(k.nombres || ' ' || k.apellidos) LIKE @like)",
		map[string]any{"like": "%" + escaparLike(b) + "%"}
}

type filaDeBaja struct {
	ID, TramiteID, TramiteCodigo, TramiteEstado string
	DocenteID                                   string
	Nombres, Apellidos, DocenteCodigo           string
	BajaNombres, BajaApellidos                  string
	Total                                       int64
}

func (r *RepositorioTramite) ListarDocumentosDeBaja(ctx context.Context, busqueda string, pagina, tamano int) ([]tramite.DocumentoDeBaja, int64, error) {
	where, args := "NOT d.activo", map[string]any{}
	if b := strings.TrimSpace(busqueda); b != "" {
		cond, a := condicionDocente(b)
		args = a
		args["t"], args["n"] = strings.ToLower(b), "%"+escaparLike(strings.ToLower(b))+"%"
		where += " AND (" + cond + " OR lower(t.codigo) = @t OR lower(d.nombre_original) LIKE @n)"
	}
	args["lim"], args["off"] = tamano, (pagina-1)*tamano
	return r.documentosDeBaja(ctx, where+" ORDER BY d.baja_en DESC NULLS LAST, d.creado_en DESC LIMIT @lim OFFSET @off", args)
}

func (r *RepositorioTramite) ObtenerDocumentoDeBaja(ctx context.Context, id string) (*tramite.DocumentoDeBaja, error) {
	lista, _, err := r.documentosDeBaja(ctx, "d.id = @id", map[string]any{"id": id})
	if err != nil {
		return nil, err
	}
	if len(lista) == 0 {
		return nil, tramite.ErrDocumentoNoEncontrado
	}
	return &lista[0], nil
}

// documentosDeBaja trae las filas de contexto con where (y su orden/límite) y luego los modelos completos.
func (r *RepositorioTramite) documentosDeBaja(ctx context.Context, where string, args map[string]any) ([]tramite.DocumentoDeBaja, int64, error) {
	var filas []filaDeBaja
	err := db(ctx, r.DB).Raw(`
		SELECT d.id, COALESCE(d.tramite_id::text, '') AS tramite_id, COALESCE(t.codigo, '') AS tramite_codigo,
		       COALESCE(t.estado, '') AS tramite_estado, k.id AS docente_id, k.nombres, k.apellidos, COALESCE(k.codigo, '') AS docente_codigo,
		       COALESCE(u.nombres, '') AS baja_nombres, COALESCE(u.apellidos, '') AS baja_apellidos,
		       count(*) OVER () AS total`+fromDocumentos+`
		LEFT JOIN usuario u ON u.id = d.baja_por
		WHERE `+where, args).Scan(&filas).Error
	if err != nil || len(filas) == 0 {
		return nil, 0, err
	}
	ids := make([]string, len(filas))
	for i, f := range filas {
		ids[i] = f.ID
	}
	var ms []TramiteDocumentoModel
	if err := db(ctx, r.DB).Where("id IN ?", ids).Find(&ms).Error; err != nil {
		return nil, 0, err
	}
	porID := make(map[string]*TramiteDocumentoModel, len(ms))
	for i := range ms {
		porID[ms[i].ID] = &ms[i]
	}
	out := make([]tramite.DocumentoDeBaja, 0, len(filas))
	for _, f := range filas {
		m, ok := porID[f.ID]
		if !ok {
			continue
		}
		doc := aDominioDocumento(m)
		doc.DocenteID = f.DocenteID
		out = append(out, tramite.DocumentoDeBaja{
			Documento: doc, TramiteID: f.TramiteID, TramiteCodigo: f.TramiteCodigo, TramiteEstado: tramite.Estado(f.TramiteEstado),
			DocenteNombre: nombreCompleto(f.Nombres, f.Apellidos), DocenteCodigo: f.DocenteCodigo,
			BajaPorNombre: nombreCompleto(f.BajaNombres, f.BajaApellidos),
		})
	}
	return out, filas[0].Total, nil
}

// RestituirDocumento reactiva el documento (CU-17). El UPDATE es condicional: si ya estaba activo
// no afecta filas y devuelve ErrDocumentoActivo.
func (r *RepositorioTramite) RestituirDocumento(ctx context.Context, id string) error {
	res := db(ctx, r.DB).Model(&TramiteDocumentoModel{}).Where("id = ? AND NOT activo", id).
		Updates(map[string]any{"activo": true, "motivo_baja": nil, "baja_en": nil, "baja_por": nil})
	if res.Error != nil {
		return res.Error
	}
	if res.RowsAffected == 0 {
		var n int64
		if err := db(ctx, r.DB).Model(&TramiteDocumentoModel{}).Where("id = ?", id).Count(&n).Error; err != nil {
			return err
		}
		if n == 0 {
			return tramite.ErrDocumentoNoEncontrado
		}
		return tramite.ErrDocumentoActivo
	}
	return nil
}

// BuscarDocumentos es la búsqueda global de CU-18: texto completo por página (índice es_unaccent,
// sin tildes ni terminaciones) y/o docente, solo sobre documentos activos. Agrupa por documento,
// lo más relevante primero, con las páginas halladas y el mejor fragmento.
func (r *RepositorioTramite) BuscarDocumentos(ctx context.Context, f tramite.FiltroBusqueda, pagina, tamano int) ([]tramite.ResultadoBusqueda, int64, error) {
	condDoc, args := condicionDocente(f.Consulta)
	args["q"], args["tipo"], args["fac"] = f.Consulta, string(f.Tipo), f.FacultadID
	args["lim"], args["off"] = tamano, (pagina-1)*tamano

	orden := "rango DESC, apellidos, nombres, nombre_original"
	if f.Recientes {
		orden = "creado_en DESC, nombre_original"
	}
	conContenido, having := "true", "count(pag.pagina) > 0 OR "+condDoc
	switch f.Campo {
	case tramite.BuscarContenido:
		having = "count(pag.pagina) > 0"
	case tramite.BuscarDocente:
		conContenido, having = "false", condDoc
	}

	var filas []struct {
		ID, TramiteID, TramiteCodigo, TipoDocumento, NombreOriginal string
		DocenteID, Nombres, Apellidos, DocenteCodigo                string
		Paginas, Fragmento                                          string
		Rango                                                       float64
		Total                                                       int64
	}
	err := db(ctx, r.DB).Raw(`
		WITH q AS (SELECT plainto_tsquery('es_unaccent', @q) AS q),
		pag AS (
			SELECT p.documento_id, p.pagina, p.texto, ts_rank(to_tsvector('es_unaccent', p.texto), q.q) AS rango
			FROM documento_pagina p, q
			WHERE `+conContenido+` AND to_tsvector('es_unaccent', p.texto) @@ q.q
		),
		docs AS (
			SELECT d.id, COALESCE(d.tramite_id::text, '') AS tramite_id, COALESCE(t.codigo, '') AS tramite_codigo,
			       COALESCE(d.tipo_documento, '') AS tipo_documento, d.nombre_original,
			       k.id AS docente_id, k.nombres, k.apellidos, COALESCE(k.codigo, '') AS docente_codigo,
			       COALESCE(max(pag.rango), 0) AS rango, d.creado_en,
			       COALESCE(string_agg(pag.pagina::text, ',' ORDER BY pag.pagina), '') AS paginas`+fromDocumentos+`
			LEFT JOIN pag ON pag.documento_id = d.id
			WHERE d.activo AND (@tipo = '' OR d.tipo_documento = @tipo) AND (@fac = 0 OR k.facultad_id = @fac)
			GROUP BY d.id, t.id, k.id
			HAVING `+having+`
		)
		SELECT docs.*, count(*) OVER () AS total,
		       COALESCE((SELECT ts_headline('es_unaccent', p.texto, q.q, 'StartSel=⟦, StopSel=⟧, MinWords=8, MaxWords=30, MaxFragments=1, ShortWord=2')
		                 FROM pag p, q WHERE p.documento_id = docs.id ORDER BY p.rango DESC, p.pagina LIMIT 1), '') AS fragmento
		FROM docs
		ORDER BY `+orden+`
		LIMIT @lim OFFSET @off`, args).Scan(&filas).Error
	if err != nil || len(filas) == 0 {
		return nil, 0, err
	}
	out := make([]tramite.ResultadoBusqueda, len(filas))
	for i, x := range filas {
		out[i] = tramite.ResultadoBusqueda{
			DocumentoID: x.ID, TramiteID: x.TramiteID, TramiteCodigo: x.TramiteCodigo, TipoDocumento: tramite.TipoDocumento(x.TipoDocumento),
			NombreOriginal: x.NombreOriginal, DocenteID: x.DocenteID, DocenteNombre: nombreCompleto(x.Nombres, x.Apellidos),
			DocenteCodigo: x.DocenteCodigo, Paginas: listaPaginas(x.Paginas), Fragmento: x.Fragmento, Rango: x.Rango,
		}
	}
	return out, filas[0].Total, nil
}

func listaPaginas(s string) []int {
	out := []int{}
	for _, p := range strings.Split(s, ",") {
		if n, err := strconv.Atoi(p); err == nil {
			out = append(out, n)
		}
	}
	return out
}
