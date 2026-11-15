package testdobles

import (
	"context"
	"strconv"
	"strings"
	"time"

	"escalafon/backend/internal/application/ports"
	"escalafon/backend/internal/domain/categoria"
	"escalafon/backend/internal/domain/docente"
	"escalafon/backend/internal/domain/tramite"
)

// --- docente.Repository -------------------------------------------------------

type DocenteRepoFalso struct{ Docentes map[string]*docente.Docente }

func NewDocenteRepoFalso() *DocenteRepoFalso {
	return &DocenteRepoFalso{Docentes: map[string]*docente.Docente{}}
}

func (r *DocenteRepoFalso) Sembrar(d docente.Docente) { r.Docentes[d.ID] = &d }

func (r *DocenteRepoFalso) Listar(_ context.Context, busqueda string, soloActivos bool) ([]docente.Docente, error) {
	var out []docente.Docente
	for _, d := range r.Docentes {
		if soloActivos && !d.Activo {
			continue
		}
		if busqueda == "" || strings.Contains(d.Cedula, busqueda) {
			out = append(out, *d)
		}
	}
	return out, nil
}

func (r *DocenteRepoFalso) ListarPagina(ctx context.Context, busqueda string, soloActivos bool, pagina, tamano int) ([]docente.Docente, int64, error) {
	todos, _ := r.Listar(ctx, busqueda, soloActivos)
	total := int64(len(todos))
	ini := min((pagina-1)*tamano, len(todos))
	return todos[ini:min(ini+tamano, len(todos))], total, nil
}

func (r *DocenteRepoFalso) Sugerir(_ context.Context, texto string, limite int) ([]docente.Docente, error) {
	t := strings.ToLower(texto)
	var out []docente.Docente
	for _, d := range r.Docentes {
		if d.Activo && (strings.HasPrefix(d.Codigo, t) || strings.Contains(strings.ToLower(d.Apellidos+" "+d.Nombres), t)) && len(out) < limite {
			out = append(out, *d)
		}
	}
	return out, nil
}

func (r *DocenteRepoFalso) BuscarPorID(_ context.Context, id string) (*docente.Docente, error) {
	if d, ok := r.Docentes[id]; ok {
		c := *d
		return &c, nil
	}
	return nil, docente.ErrNoEncontrado
}

func (r *DocenteRepoFalso) ListarFacultades(context.Context) ([]docente.Facultad, error) {
	return nil, nil
}

func (r *DocenteRepoFalso) BuscarPorIdentificador(_ context.Context, valor string) (*docente.Docente, error) {
	for _, d := range r.Docentes {
		if d.Codigo == valor || (d.Cedula != "" && d.Cedula == valor) {
			c := *d
			return &c, nil
		}
	}
	return nil, docente.ErrNoEncontrado
}

func (r *DocenteRepoFalso) Crear(_ context.Context, d *docente.Docente) error {
	for _, e := range r.Docentes {
		if e.Codigo == d.Codigo && d.Codigo != "" {
			return docente.ErrCodigoDuplicado
		}
	}
	d.ID, d.Activo = "doc-"+strconv.Itoa(len(r.Docentes)+1), true
	c := *d
	r.Docentes[d.ID] = &c
	return nil
}

func (r *DocenteRepoFalso) Actualizar(_ context.Context, d *docente.Docente) error {
	e, ok := r.Docentes[d.ID]
	if !ok {
		return docente.ErrNoEncontrado
	}
	e.Nombres, e.Apellidos, e.Correo, e.Telefono, e.CategoriaID = d.Nombres, d.Apellidos, d.Correo, d.Telefono, d.CategoriaID
	return nil
}

func (r *DocenteRepoFalso) CambiarActivo(_ context.Context, id string, activo bool) error {
	d, ok := r.Docentes[id]
	if !ok {
		return docente.ErrNoEncontrado
	}
	d.Activo = activo
	return nil
}

// --- categoria.Repository --------------------------------------------------------

type CategoriaRepoFalso struct {
	Categorias map[int16]*categoria.Categoria
}

func NewCategoriaRepoFalso() *CategoriaRepoFalso {
	return &CategoriaRepoFalso{Categorias: map[int16]*categoria.Categoria{}}
}

func (r *CategoriaRepoFalso) Sembrar(c categoria.Categoria) { r.Categorias[c.ID] = &c }

func (r *CategoriaRepoFalso) Listar(_ context.Context, soloActivas bool) ([]categoria.Categoria, error) {
	var out []categoria.Categoria
	for _, c := range r.Categorias {
		if !soloActivas || c.Activa {
			out = append(out, *c)
		}
	}
	return out, nil
}

func (r *CategoriaRepoFalso) BuscarPorID(_ context.Context, id int16) (*categoria.Categoria, error) {
	if c, ok := r.Categorias[id]; ok {
		cp := *c
		return &cp, nil
	}
	return nil, categoria.ErrNoEncontrado
}

func (r *CategoriaRepoFalso) Guardar(_ context.Context, c *categoria.Categoria, requisitos []string) error {
	if c.ID == 0 {
		c.ID = int16(len(r.Categorias) + 1)
	}
	c.Requisitos = nil
	for i, n := range requisitos {
		c.Requisitos = append(c.Requisitos, categoria.Requisito{ID: int32(i + 1), Nombre: n})
	}
	cp := *c
	r.Categorias[c.ID] = &cp
	return nil
}

func (r *CategoriaRepoFalso) Eliminar(_ context.Context, id int16) error {
	if _, ok := r.Categorias[id]; !ok {
		return categoria.ErrNoEncontrado
	}
	delete(r.Categorias, id)
	return nil
}

func (r *CategoriaRepoFalso) CambiarActiva(_ context.Context, id int16, activa bool) error {
	c, ok := r.Categorias[id]
	if !ok {
		return categoria.ErrNoEncontrado
	}
	c.Activa = activa
	return nil
}

// --- tramite.Repository -------------------------------------------------------------

type TramiteRepoFalso struct {
	Paginas       map[string][]string // texto OCR por documento
	Tramites      map[string]*tramite.Tramite
	Observaciones map[string][]*tramite.Observacion // por trámite
	Analistas     []tramite.Analista
	contador      int
}

func NewTramiteRepoFalso() *TramiteRepoFalso {
	return &TramiteRepoFalso{Tramites: map[string]*tramite.Tramite{}, Observaciones: map[string][]*tramite.Observacion{}, Paginas: map[string][]string{}}
}

func (r *TramiteRepoFalso) copia(t *tramite.Tramite) *tramite.Tramite {
	c := *t
	c.Observaciones = nil
	for _, o := range r.Observaciones[t.ID] {
		c.Observaciones = append(c.Observaciones, *o)
	}
	return &c
}

func (r *TramiteRepoFalso) Crear(_ context.Context, t *tramite.Tramite) error {
	for _, e := range r.Tramites {
		if e.Codigo == t.Codigo {
			return tramite.ErrCodigoDuplicado
		}
	}
	c := *t
	r.Tramites[t.ID] = &c
	return nil
}

func (r *TramiteRepoFalso) BuscarPorID(_ context.Context, id string) (*tramite.Tramite, error) {
	if t, ok := r.Tramites[id]; ok {
		return r.copia(t), nil
	}
	return nil, tramite.ErrNoEncontrado
}

func (r *TramiteRepoFalso) BuscarPorCodigo(_ context.Context, codigo string) (*tramite.Tramite, error) {
	for _, t := range r.Tramites {
		if t.Codigo == codigo {
			return r.copia(t), nil
		}
	}
	return nil, tramite.ErrNoEncontrado
}

func (r *TramiteRepoFalso) Listar(_ context.Context, f tramite.Filtro) ([]tramite.Resumen, error) {
	var out []tramite.Resumen
	for id, t := range r.Tramites {
		if id == "" { // contenedor imaginario de los escaneos del expediente del docente
			continue
		}
		if f.Estado == "" || t.Estado == f.Estado {
			out = append(out, tramite.Resumen{ID: t.ID, Codigo: t.Codigo, Estado: t.Estado})
		}
	}
	return out, nil
}

func (r *TramiteRepoFalso) ListarPagina(ctx context.Context, f tramite.Filtro, pagina, tamano int) ([]tramite.Resumen, int64, error) {
	todos, _ := r.Listar(ctx, f)
	ini := min((pagina-1)*tamano, len(todos))
	return todos[ini:min(ini+tamano, len(todos))], int64(len(todos)), nil
}

func (r *TramiteRepoFalso) Guardar(_ context.Context, t *tramite.Tramite, h *tramite.Historial) error {
	e, ok := r.Tramites[t.ID]
	if !ok {
		return tramite.ErrNoEncontrado
	}
	hist := e.Historial
	*e = *t
	e.Historial = hist
	if h != nil {
		e.Historial = append(e.Historial, *h)
	}
	return nil
}

func (r *TramiteRepoFalso) AgregarObservacion(_ context.Context, tramiteID string, o *tramite.Observacion, _ string) error {
	r.contador++
	o.ID = "obs-" + strconv.Itoa(r.contador)
	c := *o
	r.Observaciones[tramiteID] = append(r.Observaciones[tramiteID], &c)
	return nil
}

func (r *TramiteRepoFalso) BuscarObservacion(_ context.Context, tramiteID, id string) (*tramite.Observacion, error) {
	for _, o := range r.Observaciones[tramiteID] {
		if o.ID == id {
			c := *o
			return &c, nil
		}
	}
	return nil, tramite.ErrObservacionNoEncontrada
}

func (r *TramiteRepoFalso) ResolverObservacion(_ context.Context, id string, en time.Time) error {
	for _, obs := range r.Observaciones {
		for _, o := range obs {
			if o.ID == id {
				o.Resuelta, o.ResueltaEn = true, &en
			}
		}
	}
	return nil
}

func (r *TramiteRepoFalso) ContarObservacionesPendientes(_ context.Context, tramiteID string) (int64, error) {
	var n int64
	for _, o := range r.Observaciones[tramiteID] {
		if !o.Resuelta {
			n++
		}
	}
	return n, nil
}

func (r *TramiteRepoFalso) BuscarDocumento(_ context.Context, tramiteID, id string) (*tramite.Documento, error) {
	if t, ok := r.Tramites[tramiteID]; ok {
		for _, d := range t.Documentos {
			if d.ID == id {
				c := d
				return &c, nil
			}
		}
	}
	return nil, tramite.ErrNoEncontrado
}

func (r *TramiteRepoFalso) DarDeBajaDocumento(_ context.Context, tramiteID, id, motivo, actorID string, en time.Time) error {
	t, ok := r.Tramites[tramiteID]
	if !ok {
		return tramite.ErrNoEncontrado
	}
	for i := range t.Documentos {
		d := &t.Documentos[i]
		if d.ID != id {
			continue
		}
		if !d.Activo {
			return tramite.ErrDocumentoDeBaja
		}
		d.Activo, d.MotivoBaja, d.BajaEn, d.BajaPor = false, motivo, &en, &actorID
		return nil
	}
	return tramite.ErrNoEncontrado
}

func (r *TramiteRepoFalso) CrearDocumento(_ context.Context, tramiteID string, d *tramite.Documento) error {
	t, ok := r.Tramites[tramiteID]
	if !ok {
		return tramite.ErrNoEncontrado
	}
	c := *d
	c.Activo = true
	t.Documentos = append(t.Documentos, c)
	return nil
}

func (r *TramiteRepoFalso) GuardarDocumento(_ context.Context, tramiteID string, d *tramite.Documento) error {
	if t, ok := r.Tramites[tramiteID]; ok {
		for i := range t.Documentos {
			if t.Documentos[i].ID == d.ID {
				c := *d
				c.Activo, c.MotivoBaja, c.BajaEn, c.BajaPor = t.Documentos[i].Activo, t.Documentos[i].MotivoBaja, t.Documentos[i].BajaEn, t.Documentos[i].BajaPor
				t.Documentos[i] = c
				return nil
			}
		}
	}
	return tramite.ErrDocumentoNoEncontrado
}

func (r *TramiteRepoFalso) tieneDocumento(tramiteID, docID string) bool {
	if t, ok := r.Tramites[tramiteID]; ok {
		for _, d := range t.Documentos {
			if d.ID == docID {
				return true
			}
		}
	}
	return false
}

func (r *TramiteRepoFalso) GuardarPaginas(_ context.Context, tramiteID, docID string, paginas []string) error {
	if !r.tieneDocumento(tramiteID, docID) {
		return tramite.ErrDocumentoNoEncontrado
	}
	r.Paginas[docID] = append([]string(nil), paginas...)
	t := r.Tramites[tramiteID]
	for i := range t.Documentos {
		if t.Documentos[i].ID == docID {
			t.Documentos[i].TextoOCR = strings.Join(paginas, "\n\n")
		}
	}
	return nil
}

func (r *TramiteRepoFalso) ListarPaginas(_ context.Context, tramiteID, docID string) ([]string, error) {
	if !r.tieneDocumento(tramiteID, docID) {
		return nil, tramite.ErrDocumentoNoEncontrado
	}
	return append([]string(nil), r.Paginas[docID]...), nil
}

// BuscarEnDocumento del doble: coincide si la página contiene todas las palabras, sin mayúsculas.
func (r *TramiteRepoFalso) BuscarEnDocumento(_ context.Context, tramiteID, docID, consulta string) ([]tramite.Coincidencia, error) {
	if !r.tieneDocumento(tramiteID, docID) {
		return nil, tramite.ErrDocumentoNoEncontrado
	}
	var out []tramite.Coincidencia
	for i, p := range r.Paginas[docID] {
		minus, ok := strings.ToLower(p), true
		for _, w := range strings.Fields(strings.ToLower(consulta)) {
			ok = ok && strings.Contains(minus, w)
		}
		if ok {
			out = append(out, tramite.Coincidencia{Pagina: i + 1, Fragmento: p, Rango: 1})
		}
	}
	return out, nil
}

func (r *TramiteRepoFalso) ListarDocumentosEnProceso(context.Context) ([]tramite.DocumentoEnProceso, error) {
	var out []tramite.DocumentoEnProceso
	for _, t := range r.Tramites {
		for _, d := range t.Documentos {
			if d.Activo && (d.EstadoOCR == tramite.OCRPendiente || d.EstadoOCR == tramite.OCRProcesando) {
				out = append(out, tramite.DocumentoEnProceso{TramiteID: t.ID, Documento: d})
			}
		}
	}
	return out, nil
}

// Los escaneos del expediente del docente viven en un trámite imaginario de id "": así BuscarDocumento,
// GuardarDocumento, GuardarPaginas… funcionan con tramiteID "" igual que en la base real.
func (r *TramiteRepoFalso) CrearDocumentoDocente(ctx context.Context, docenteID string, d *tramite.Documento) error {
	if _, ok := r.Tramites[""]; !ok {
		r.Tramites[""] = &tramite.Tramite{}
	}
	c := *d
	c.DocenteID = docenteID
	return r.CrearDocumento(ctx, "", &c)
}

func (r *TramiteRepoFalso) ListarDocumentosDocente(_ context.Context, docenteID string, incluirBajas bool) ([]tramite.Documento, error) {
	var out []tramite.Documento
	if t, ok := r.Tramites[""]; ok {
		for _, d := range t.Documentos {
			if d.DocenteID == docenteID && (d.Activo || incluirBajas) {
				out = append(out, d)
			}
		}
	}
	return out, nil
}

func (r *TramiteRepoFalso) BuscarEnDocente(ctx context.Context, docenteID, consulta string) ([]tramite.CoincidenciaDocente, error) {
	docs, _ := r.ListarDocumentosDocente(ctx, docenteID, false)
	var out []tramite.CoincidenciaDocente
	for _, d := range docs {
		if d.EstadoOCR != tramite.OCRProcesado {
			continue
		}
		cs, _ := r.BuscarEnDocumento(ctx, "", d.ID, consulta)
		for _, c := range cs {
			out = append(out, tramite.CoincidenciaDocente{DocumentoID: d.ID, TipoDocumento: d.TipoDocumento, NombreOriginal: d.NombreOriginal, Coincidencia: c})
		}
	}
	return out, nil
}

func (r *TramiteRepoFalso) ListarPorDocente(_ context.Context, docenteID string) ([]tramite.Resumen, error) {
	var out []tramite.Resumen
	for id, t := range r.Tramites {
		if id != "" && t.DocenteID == docenteID {
			out = append(out, tramite.Resumen{ID: t.ID, Codigo: t.Codigo, Tipo: t.Tipo, Estado: t.Estado, Categoria: t.Categoria, PendienteValidacion: t.PendienteValidacion, CreadoEn: t.CreadoEn})
		}
	}
	return out, nil
}

func (r *TramiteRepoFalso) ListarAnalistas(context.Context) ([]tramite.Analista, error) {
	return r.Analistas, nil
}

// --- ports.ObjectStorage -------------------------------------------------------------

type AlmacenFalso struct{ Guardados map[string][]byte }

func NewAlmacenFalso() *AlmacenFalso { return &AlmacenFalso{Guardados: map[string][]byte{}} }

func (*AlmacenFalso) URLSubida(_ context.Context, clave, _ string, _ time.Duration) (string, error) {
	return "https://almacen.test/subir/" + clave, nil
}

func (a *AlmacenFalso) Tamano(_ context.Context, clave string) (int64, error) {
	d, ok := a.Guardados[clave]
	if !ok {
		return 0, tramite.ErrSubidaIncompleta
	}
	return int64(len(d)), nil
}

// ProcesadorFalso responde de inmediato con el resultado configurado (o con Err).
type ProcesadorFalso struct {
	Resultado ports.ResultadoProceso
	Err       error
	Recibidos []ports.SolicitudProceso
}

func (p *ProcesadorFalso) Iniciar(_ context.Context, s ports.SolicitudProceso) (string, error) {
	p.Recibidos = append(p.Recibidos, s)
	return "job-1", p.Err
}

func (p *ProcesadorFalso) Estado(context.Context, string) (*ports.ResultadoProceso, error) {
	r := p.Resultado
	return &r, nil
}

func (*AlmacenFalso) Bucket() string { return "bucket-falso" }

func (a *AlmacenFalso) Guardar(_ context.Context, clave, _ string, datos []byte) error {
	a.Guardados[clave] = datos
	return nil
}

func (*AlmacenFalso) URLFirmada(_ context.Context, clave string, _ time.Duration) (string, error) {
	return "https://almacen.test/" + clave, nil
}

// --- CU-17 y CU-18 ---

func (r *TramiteRepoFalso) buscarEnTodos(id string) (string, *tramite.Tramite, *tramite.Documento) {
	for tid, t := range r.Tramites {
		for i := range t.Documentos {
			if t.Documentos[i].ID == id {
				return tid, t, &t.Documentos[i]
			}
		}
	}
	return "", nil, nil
}

func (r *TramiteRepoFalso) deBaja(tid string, t *tramite.Tramite, d *tramite.Documento) tramite.DocumentoDeBaja {
	x := tramite.DocumentoDeBaja{Documento: *d, TramiteID: tid}
	if tid != "" {
		x.TramiteCodigo, x.TramiteEstado, x.DocenteNombre, x.DocenteCodigo = t.Codigo, t.Estado, t.DocenteNombre, t.DocenteCodigo
		x.DocenteID = t.DocenteID
	}
	return x
}

func (r *TramiteRepoFalso) ListarDocumentosDeBaja(_ context.Context, _ string, _, _ int) ([]tramite.DocumentoDeBaja, int64, error) {
	var out []tramite.DocumentoDeBaja
	for tid, t := range r.Tramites {
		for i := range t.Documentos {
			if !t.Documentos[i].Activo {
				out = append(out, r.deBaja(tid, t, &t.Documentos[i]))
			}
		}
	}
	return out, int64(len(out)), nil
}

func (r *TramiteRepoFalso) ObtenerDocumentoDeBaja(_ context.Context, id string) (*tramite.DocumentoDeBaja, error) {
	tid, t, d := r.buscarEnTodos(id)
	if d == nil {
		return nil, tramite.ErrDocumentoNoEncontrado
	}
	x := r.deBaja(tid, t, d)
	return &x, nil
}

func (r *TramiteRepoFalso) RestituirDocumento(_ context.Context, id string) error {
	_, _, d := r.buscarEnTodos(id)
	if d == nil {
		return tramite.ErrDocumentoNoEncontrado
	}
	if d.Activo {
		return tramite.ErrDocumentoActivo
	}
	d.Activo, d.MotivoBaja, d.BajaEn, d.BajaPor = true, "", nil, nil
	return nil
}

func (r *TramiteRepoFalso) BuscarDocumentos(ctx context.Context, f tramite.FiltroBusqueda, _, _ int) ([]tramite.ResultadoBusqueda, int64, error) {
	var out []tramite.ResultadoBusqueda
	for tid, t := range r.Tramites {
		for _, d := range t.Documentos {
			if !d.Activo || (f.Tipo != "" && d.TipoDocumento != f.Tipo) {
				continue
			}
			cs, _ := r.BuscarEnDocumento(ctx, tid, d.ID, f.Consulta)
			if len(cs) == 0 {
				continue
			}
			res := tramite.ResultadoBusqueda{DocumentoID: d.ID, TramiteID: tid, NombreOriginal: d.NombreOriginal, TipoDocumento: d.TipoDocumento, Fragmento: cs[0].Fragmento}
			for _, c := range cs {
				res.Paginas = append(res.Paginas, c.Pagina)
			}
			out = append(out, res)
		}
	}
	return out, int64(len(out)), nil
}
