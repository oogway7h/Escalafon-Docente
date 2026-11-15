// Package tramites implementa los casos de uso CU-08 al CU-13: radicación, asignación,
// revisión, dictamen, archivo y consulta pública de trámites.
package tramites

import (
	"crypto/subtle"
	"context"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"net/http"
	"regexp"
	"strings"
	"time"

	"github.com/google/uuid"

	"escalafon/backend/internal/application/ports"
	"escalafon/backend/internal/domain/auditoria"
	"escalafon/backend/internal/domain/categoria"
	"escalafon/backend/internal/domain/docente"
	"escalafon/backend/internal/domain/tramite"
)

const maxTamanoArchivo = 8 << 20 // 8 MB por adjunto

var tiposPermitidos = map[string]bool{"application/pdf": true, "image/jpeg": true, "image/png": true}

// registrar anota el evento en auditoría; un fallo de auditoría no deshace la operación.
func registrar(ctx context.Context, a auditoria.Repository, actorID *string, accion, id string, ahora time.Time, detalle map[string]any) {
	_ = a.Registrar(ctx, &auditoria.Evento{
		UsuarioID: actorID, Accion: accion, Entidad: auditoria.EntidadTramite, EntidadID: id, Detalle: detalle, Fecha: ahora,
	})
}

// Notificador avisa al docente de una transición (CU-14). Es opcional: nil = sin avisos.
type Notificador interface {
	Transicion(ctx context.Context, t *tramite.Tramite)
}

func notificar(ctx context.Context, n Notificador, t *tramite.Tramite) {
	if n != nil {
		n.Transicion(ctx, t)
	}
}

// crearConCodigo inserta el trámite reintentando si el código aleatorio ya existe.
func crearConCodigo(ctx context.Context, repo tramite.Repository, t *tramite.Tramite, ahora time.Time) error {
	for i := 0; i < 5; i++ {
		codigo, err := tramite.NuevoCodigo(ahora.Year())
		if err != nil {
			return err
		}
		t.Codigo = codigo
		err = repo.Crear(ctx, t)
		if !errors.Is(err, tramite.ErrCodigoDuplicado) {
			return err
		}
	}
	return tramite.ErrCodigoDuplicado
}

func categoriaActiva(ctx context.Context, repo categoria.Repository, id int16) (*categoria.Categoria, error) {
	cat, err := repo.BuscarPorID(ctx, id)
	if err != nil || !cat.Activa {
		return nil, tramite.ErrCategoriaInvalida
	}
	return cat, nil
}

// --- CU-08 Radicar trámite -------------------------------------------------

type Radicar struct {
	Tramites   tramite.Repository
	Docentes   docente.Repository
	Categorias categoria.Repository
	Auditoria  auditoria.Repository
	Reloj      ports.Clock
}

type RadicarEntrada struct {
	ActorID     string
	DocenteID   string
	Tipo        tramite.Tipo
	CategoriaID int16
}

func (c *Radicar) Ejecutar(ctx context.Context, in RadicarEntrada) (*tramite.Tramite, error) {
	if !in.Tipo.Valido() {
		return nil, tramite.ErrTransicionInvalida
	}
	d, err := c.Docentes.BuscarPorID(ctx, in.DocenteID)
	if err != nil {
		return nil, tramite.ErrDocenteNoEncontrado
	}
	if !d.Activo {
		return nil, tramite.ErrDocenteInactivo
	}
	// Una regularización no postula a ninguna categoría: se ignora la que llegue (CategoriaID 0 = sin categoría).
	var categoriaID int16
	if in.Tipo == tramite.Categorizacion {
		if _, err := categoriaActiva(ctx, c.Categorias, in.CategoriaID); err != nil {
			return nil, err
		}
		categoriaID = in.CategoriaID
	}
	ahora := c.Reloj.Ahora()
	t := &tramite.Tramite{
		ID: uuid.NewString(), DocenteID: d.ID, DocenteCodigo: d.Codigo, CategoriaID: categoriaID, Tipo: in.Tipo,
		Estado: tramite.Ingresado, Origen: tramite.OrigenStaff, RadicadoPor: &in.ActorID,
		Historial: []tramite.Historial{{Estado: tramite.Ingresado, UsuarioID: &in.ActorID, Detalle: "Trámite radicado en ventanilla", Fecha: ahora}},
	}
	if err := crearConCodigo(ctx, c.Tramites, t, ahora); err != nil {
		return nil, err
	}
	registrar(ctx, c.Auditoria, &in.ActorID, auditoria.AccionRadicar, t.ID, ahora, map[string]any{"codigo": t.Codigo})
	return c.Tramites.BuscarPorID(ctx, t.ID)
}

// --- CU-13 Radicar solicitud desde el portal público --------------------------

type Archivo struct {
	RequisitoID int32
	Nombre      string
	Datos       []byte
}

type RadicarPublico struct {
	Tramites   tramite.Repository
	Docentes   docente.Repository
	Categorias categoria.Repository
	Almacen    ports.ObjectStorage
	Auditoria  auditoria.Repository
	Reloj      ports.Clock
}

type RadicarPublicoEntrada struct {
	Codigo      string // código universitario del docente
	Cedula      string // debe coincidir con la del padrón: impide radicar a nombre de otro
	Tipo        tramite.Tipo
	CategoriaID int16     // solo en categorización
	Archivos    []Archivo // categorización: uno por requisito; regularización: de 1 a 5 con RequisitoID 0
}

const (
	maxAdjuntosLibres       = 5 // regularización: documentos sin requisito asociado
	maxPendientesPorDocente = 2 // solicitudes del portal sin validar en ventanilla
)

var caracteresNoSeguros = regexp.MustCompile(`[^A-Za-z0-9._-]+`)

// docenteVerificado identifica al docente del portal por su código Y su cédula. La cédula es lo que impide que
// cualquiera recorra los códigos de docente (cortos y consecutivos). Código desconocido, cédula que no coincide
// y docente sin cédula registrada devuelven el mismo error, para que el portal no sirva para averiguar quién existe.
func docenteVerificado(ctx context.Context, repo docente.Repository, codigo, cedula string) (*docente.Docente, error) {
	codigo, cedula = strings.ToLower(strings.TrimSpace(codigo)), strings.ToLower(strings.TrimSpace(cedula))
	if codigo == "" || cedula == "" || len(codigo) > 20 || len(cedula) > 20 {
		return nil, tramite.ErrDocenteNoVerificado
	}
	d, err := repo.BuscarPorIdentificador(ctx, codigo)
	if err != nil {
		if errors.Is(err, docente.ErrNoEncontrado) {
			return nil, tramite.ErrDocenteNoVerificado
		}
		return nil, err
	}
	// BuscarPorIdentificador también acepta una cédula como "código": aquí el código debe ser el del docente.
	cedulaReal := strings.ToLower(strings.TrimSpace(d.Cedula))
	if strings.ToLower(d.Codigo) != codigo || cedulaReal == "" || subtle.ConstantTimeCompare([]byte(cedulaReal), []byte(cedula)) != 1 {
		return nil, tramite.ErrDocenteNoVerificado
	}
	return d, nil
}

// Ejecutar verifica al docente (código + cédula) y los adjuntos, sube los archivos y crea el trámite en Ingresado,
// pendiente de validación. El OCR no se lanza aquí sino al validar en ventanilla: una solicitud sin validar no
// debe poder ocupar el servicio de OCR.
// ponytail: si falla la inserción después de subir, los objetos quedan huérfanos en el bucket;
// limpiar con una política de ciclo de vida del bucket si molesta.
func (c *RadicarPublico) Ejecutar(ctx context.Context, in RadicarPublicoEntrada) (string, error) {
	if !in.Tipo.Valido() {
		in.Tipo = tramite.Categorizacion
	}
	d, err := docenteVerificado(ctx, c.Docentes, in.Codigo, in.Cedula)
	if err != nil {
		return "", err
	}
	if !d.Activo { // mismo error que un docente desconocido
		return "", tramite.ErrDocenteNoVerificado
	}
	previos, err := c.Tramites.ListarPorDocente(ctx, d.ID)
	if err != nil {
		return "", err
	}
	pendientes := 0
	for _, t := range previos {
		if t.PendienteValidacion {
			pendientes++
		}
	}
	if pendientes >= maxPendientesPorDocente {
		return "", tramite.ErrDemasiadasPendientes
	}

	type adjunto struct {
		requisito *int32
		archivo   Archivo
	}
	var adjuntos []adjunto
	var categoriaID int16
	if in.Tipo == tramite.Categorizacion {
		cat, err := categoriaActiva(ctx, c.Categorias, in.CategoriaID)
		if err != nil {
			return "", err
		}
		categoriaID = cat.ID
		porRequisito := map[int32]Archivo{}
		for _, a := range in.Archivos {
			porRequisito[a.RequisitoID] = a
		}
		var faltan []string
		for _, r := range cat.Requisitos {
			a, ok := porRequisito[r.ID]
			if !ok {
				faltan = append(faltan, r.Nombre)
				continue
			}
			reqID := r.ID
			adjuntos = append(adjuntos, adjunto{&reqID, a})
		}
		if len(faltan) > 0 {
			return "", &tramite.ErrDocumentosFaltantes{Faltan: faltan}
		}
	} else { // regularización: sin categoría, adjuntos libres
		for _, a := range in.Archivos {
			if a.RequisitoID == 0 {
				adjuntos = append(adjuntos, adjunto{nil, a})
			}
		}
		if len(adjuntos) < 1 || len(adjuntos) > maxAdjuntosLibres {
			return "", tramite.ErrAdjuntosLibres
		}
	}

	id := uuid.NewString()
	var docs []tramite.Documento
	for _, adj := range adjuntos {
		a := adj.archivo
		if len(a.Datos) == 0 || len(a.Datos) > maxTamanoArchivo {
			return "", tramite.ErrArchivoInvalido
		}
		mime := http.DetectContentType(a.Datos)
		if !tiposPermitidos[mime] {
			return "", tramite.ErrArchivoInvalido
		}
		docID := uuid.NewString()
		clave := claveDocumento(carpetaDocente(d.Codigo, d.ID), id, docID, a.Nombre)
		if err := c.Almacen.Guardar(ctx, clave, mime, a.Datos); err != nil {
			return "", err
		}
		hash := sha256.Sum256(a.Datos)
		docs = append(docs, tramite.Documento{
			ID: docID, RequisitoID: adj.requisito, NombreOriginal: a.Nombre, TipoMime: mime, Tamano: int64(len(a.Datos)),
			HashSHA256: hex.EncodeToString(hash[:]), Bucket: c.Almacen.Bucket(), ClaveObjeto: clave, Activo: true,
			EstadoOCR: tramite.OCRNoAplica, // pasa a PENDIENTE al validarse en ventanilla
		})
	}

	ahora := c.Reloj.Ahora()
	t := &tramite.Tramite{
		ID: id, DocenteID: d.ID, CategoriaID: categoriaID, Tipo: in.Tipo, Estado: tramite.Ingresado,
		Origen: tramite.OrigenPortal, PendienteValidacion: true, Documentos: docs,
		Historial: []tramite.Historial{{Estado: tramite.Ingresado, Detalle: "Solicitud radicada desde el portal público", Fecha: ahora}},
	}
	if err := crearConCodigo(ctx, c.Tramites, t, ahora); err != nil {
		return "", err
	}
	registrar(ctx, c.Auditoria, nil, auditoria.AccionRadicarPortal, t.ID, ahora, map[string]any{"codigo": t.Codigo})
	return t.Codigo, nil
}

// ValidarSolicitud confirma en ventanilla una solicitud llegada por el portal.
type ValidarSolicitud struct {
	Tramites  tramite.Repository
	Auditoria auditoria.Repository
	Reloj     ports.Clock
	OCR       *ProcesoOCR // opcional: si está, los adjuntos del portal se procesan con OCR al validar (nil = sin OCR)
}

func (c *ValidarSolicitud) Ejecutar(ctx context.Context, actorID, id string) error {
	t, err := c.Tramites.BuscarPorID(ctx, id)
	if err != nil {
		return err
	}
	if t.Estado.Cerrado() {
		return t.Estado.ErrCerrado()
	}
	if !t.PendienteValidacion {
		return tramite.ErrYaValidado
	}
	ahora := c.Reloj.Ahora()
	t.PendienteValidacion = false
	h := &tramite.Historial{Estado: t.Estado, EstadoAnterior: t.Estado, UsuarioID: &actorID, Detalle: "Solicitud validada en ventanilla", Fecha: ahora}
	if err := c.Tramites.Guardar(ctx, t, h); err != nil {
		return err
	}
	registrar(ctx, c.Auditoria, &actorID, auditoria.AccionValidarSolicitud, t.ID, ahora, nil)
	if c.OCR != nil { // el texto de los adjuntos queda indexado para el buscador del visor
		for _, d := range t.Documentos {
			if !d.Activo || d.EstadoOCR != tramite.OCRNoAplica {
				continue
			}
			d.EstadoOCR = tramite.OCRPendiente
			if err := c.Tramites.GuardarDocumento(ctx, t.ID, &d); err != nil {
				return err
			}
			c.OCR.Lanzar(t.ID, d)
		}
	}
	return nil
}

// --- Consultas internas ---------------------------------------------------------

type Listar struct{ Tramites tramite.Repository }

func (c *Listar) Ejecutar(ctx context.Context, f tramite.Filtro) ([]tramite.Resumen, error) {
	return c.Tramites.Listar(ctx, f)
}

// ListarPagina devuelve una página de trámites (tamaño acotado a 1..100) y el total de coincidencias.
type ListarPagina struct{ Tramites tramite.Repository }

func (c *ListarPagina) Ejecutar(ctx context.Context, f tramite.Filtro, pagina, tamano int) ([]tramite.Resumen, int64, error) {
	if pagina < 1 {
		pagina = 1
	}
	if tamano < 1 || tamano > 100 {
		tamano = 10
	}
	return c.Tramites.ListarPagina(ctx, f, pagina, tamano)
}

type Obtener struct{ Tramites tramite.Repository }

func (c *Obtener) Ejecutar(ctx context.Context, id string) (*tramite.Tramite, error) {
	return c.Tramites.BuscarPorID(ctx, id)
}

type ListarAnalistas struct{ Tramites tramite.Repository }

func (c *ListarAnalistas) Ejecutar(ctx context.Context) ([]tramite.Analista, error) {
	return c.Tramites.ListarAnalistas(ctx)
}

// UrlDocumento devuelve un enlace temporal para descargar un adjunto del trámite (la versión comprimida
// si existe). Un documento dado de baja solo lo descarga quien puede restituirlo (verBajas).
type UrlDocumento struct {
	Tramites tramite.Repository
	Almacen  ports.ObjectStorage
}

func (c *UrlDocumento) Ejecutar(ctx context.Context, tramiteID, docID string, verBajas bool) (string, error) {
	d, err := c.Tramites.BuscarDocumento(ctx, tramiteID, docID)
	if err != nil {
		return "", err
	}
	if !d.Activo && !verBajas {
		return "", tramite.ErrDocumentoDeBaja
	}
	return c.Almacen.URLFirmada(ctx, d.ClaveDescarga(), 30*time.Minute) // el visor pide el PDF por tramos durante la lectura
}

// --- CU-12 Consulta pública -----------------------------------------------------

type EtapaPublica struct {
	Estado tramite.Estado
	Fecha  time.Time
}

type ObservacionPublica struct {
	Documento string
	Motivo    string
	Plazo     time.Time
	Resuelta  bool
}

// VistaPublica es lo único que ve el docente: ningún dato personal ni interno del trámite.
type VistaPublica struct {
	Codigo              string
	Tipo                tramite.Tipo
	Categoria           string
	Estado              tramite.Estado
	RadicadoEn          time.Time
	PendienteValidacion bool
	Resultado           tramite.Resultado
	Etapas              []EtapaPublica
	Observaciones       []ObservacionPublica
}

// ResumenPublico es lo único que el portal muestra de cada trámite al buscar por código de docente: sin nombres,
// analistas ni observaciones. El detalle se abre con el código de seguimiento.
type ResumenPublico struct {
	Codigo              string
	Tipo                tramite.Tipo
	Categoria           string
	Estado              tramite.Estado
	RadicadoEn          time.Time
	PendienteValidacion bool
}

// ConsultarPublicoPorDocente ayuda al docente que perdió su código de seguimiento: con su código de docente Y su
// cédula lista sus trámites. La cédula es lo que impide que cualquiera recorra códigos de docente (cortos y
// consecutivos). Un código desconocido, una cédula que no coincide o un docente sin cédula registrada responden
// igual que quien no tiene trámites, para que el portal no sirva para averiguar qué docentes existen.
type ConsultarPublicoPorDocente struct {
	Docentes docente.Repository
	Tramites tramite.Repository
}

func (c *ConsultarPublicoPorDocente) Ejecutar(ctx context.Context, codigo, cedula string) ([]ResumenPublico, error) {
	d, err := docenteVerificado(ctx, c.Docentes, codigo, cedula)
	if errors.Is(err, tramite.ErrDocenteNoVerificado) {
		return []ResumenPublico{}, nil
	}
	if err != nil {
		return nil, err
	}
	lista, err := c.Tramites.ListarPorDocente(ctx, d.ID)
	if err != nil {
		return nil, err
	}
	out := make([]ResumenPublico, len(lista))
	for i, t := range lista {
		out[i] = ResumenPublico{Codigo: t.Codigo, Tipo: t.Tipo, Categoria: t.Categoria, Estado: t.Estado, RadicadoEn: t.CreadoEn, PendienteValidacion: t.PendienteValidacion}
	}
	return out, nil
}

type ConsultarPublico struct{ Tramites tramite.Repository }

func (c *ConsultarPublico) Ejecutar(ctx context.Context, codigo string) (*VistaPublica, error) {
	t, err := c.Tramites.BuscarPorCodigo(ctx, strings.ToUpper(strings.TrimSpace(codigo)))
	if err != nil {
		return nil, err
	}
	v := &VistaPublica{
		Codigo: t.Codigo, Tipo: t.Tipo, Categoria: t.Categoria, Estado: t.Estado,
		RadicadoEn: t.CreadoEn, PendienteValidacion: t.PendienteValidacion,
	}
	if t.Estado == tramite.Remitido || t.Estado == tramite.Archivado { // antes es un informe interno
		v.Resultado = t.Resultado
	}
	for _, h := range t.Historial {
		if h.Estado != h.EstadoAnterior { // solo cambios de fase, no reasignaciones ni validaciones
			v.Etapas = append(v.Etapas, EtapaPublica{Estado: h.Estado, Fecha: h.Fecha})
		}
	}
	for _, o := range t.Observaciones {
		v.Observaciones = append(v.Observaciones, ObservacionPublica{Documento: o.Documento, Motivo: o.Motivo, Plazo: o.Plazo, Resuelta: o.Resuelta})
	}
	return v, nil
}

// --- CU-09 Asignar o reasignar -----------------------------------------------------

type Asignar struct {
	Tramites  tramite.Repository
	Auditoria auditoria.Repository
	Reloj     ports.Clock
}

type AsignarEntrada struct {
	ActorID    string
	TramiteID  string
	AnalistaID string
	Motivo     string // obligatorio al reasignar
}

func (c *Asignar) Ejecutar(ctx context.Context, in AsignarEntrada) (*tramite.Tramite, error) {
	t, err := c.Tramites.BuscarPorID(ctx, in.TramiteID)
	if err != nil {
		return nil, err
	}
	if t.Estado.Cerrado() {
		return nil, t.Estado.ErrCerrado()
	}
	if t.PendienteValidacion {
		return nil, tramite.ErrPendienteValidacion
	}
	analistas, err := c.Tramites.ListarAnalistas(ctx)
	if err != nil {
		return nil, err
	}
	var nombre string
	for _, a := range analistas {
		if a.ID == in.AnalistaID {
			nombre = a.Nombre
		}
	}
	if nombre == "" {
		return nil, tramite.ErrAnalistaInvalido
	}

	ahora := c.Reloj.Ahora()
	motivo := strings.TrimSpace(in.Motivo)
	var h tramite.Historial
	accion := auditoria.AccionAsignarTramite
	switch {
	case t.Estado == tramite.Ingresado:
		if err := tramite.Validar(t.Estado, tramite.Asignado); err != nil {
			return nil, err
		}
		h = tramite.Historial{Estado: tramite.Asignado, EstadoAnterior: t.Estado, Detalle: "Asignado a " + nombre}
		t.Estado = tramite.Asignado
	default: // también En facultad: el analista de reemplazo queda listo para el regreso
		if motivo == "" {
			return nil, tramite.ErrMotivoRequerido
		}
		h = tramite.Historial{Estado: t.Estado, EstadoAnterior: t.Estado, Detalle: fmt.Sprintf("Reasignado a %s. Motivo: %s", nombre, motivo)}
		accion = auditoria.AccionReasignarTramite
	}
	h.UsuarioID, h.Fecha = &in.ActorID, ahora
	t.AnalistaID = &in.AnalistaID
	if err := c.Tramites.Guardar(ctx, t, &h); err != nil {
		return nil, err
	}
	registrar(ctx, c.Auditoria, &in.ActorID, accion, t.ID, ahora, map[string]any{"analistaId": in.AnalistaID, "motivo": motivo})
	return c.Tramites.BuscarPorID(ctx, t.ID)
}

// --- CU-10 Revisar expediente y registrar observaciones ------------------------------

// cargarDelAnalista trae el trámite y exige que el actor sea su analista asignado.
func cargarDelAnalista(ctx context.Context, repo tramite.Repository, actorID, id string) (*tramite.Tramite, error) {
	t, err := repo.BuscarPorID(ctx, id)
	if err != nil {
		return nil, err
	}
	if t.Estado.Cerrado() {
		return nil, t.Estado.ErrCerrado()
	}
	if t.AnalistaID == nil || *t.AnalistaID != actorID {
		return nil, tramite.ErrNoAsignado
	}
	return t, nil
}

type IniciarRevision struct {
	Tramites  tramite.Repository
	Auditoria auditoria.Repository
	Reloj     ports.Clock
}

func (c *IniciarRevision) Ejecutar(ctx context.Context, actorID, id string) error {
	t, err := cargarDelAnalista(ctx, c.Tramites, actorID, id)
	if err != nil {
		return err
	}
	if err := tramite.Validar(t.Estado, tramite.EnRevision); err != nil {
		return err
	}
	if t.Estado == tramite.Asignado && t.CarpetasLigadas == 0 { // la revisión parte de la carpeta física
		return tramite.ErrCarpetaNoSolicitada
	}
	ahora := c.Reloj.Ahora()
	h := &tramite.Historial{Estado: tramite.EnRevision, EstadoAnterior: t.Estado, UsuarioID: &actorID, Detalle: "Revisión iniciada", Fecha: ahora}
	t.Estado = tramite.EnRevision
	if err := c.Tramites.Guardar(ctx, t, h); err != nil {
		return err
	}
	registrar(ctx, c.Auditoria, &actorID, auditoria.AccionIniciarRevision, t.ID, ahora, nil)
	return nil
}

type RegistrarObservacion struct {
	Tramites    tramite.Repository
	Auditoria   auditoria.Repository
	Reloj       ports.Clock
	Notificador Notificador
}

type ObservacionEntrada struct {
	ActorID   string
	TramiteID string
	Documento string
	Motivo    string
	Plazo     time.Time
}

func (c *RegistrarObservacion) Ejecutar(ctx context.Context, in ObservacionEntrada) error {
	t, err := cargarDelAnalista(ctx, c.Tramites, in.ActorID, in.TramiteID)
	if err != nil {
		return err
	}
	if t.Estado != tramite.Observado { // ya observado: se agregan más observaciones sin cambiar de fase
		if err := tramite.Validar(t.Estado, tramite.Observado); err != nil {
			return err
		}
	}
	ahora := c.Reloj.Ahora()
	hoy := time.Date(ahora.Year(), ahora.Month(), ahora.Day(), 0, 0, 0, 0, in.Plazo.Location())
	if in.Plazo.Before(hoy) {
		return tramite.ErrPlazoInvalido
	}
	obs := &tramite.Observacion{Documento: strings.TrimSpace(in.Documento), Motivo: strings.TrimSpace(in.Motivo), Plazo: in.Plazo}
	if err := c.Tramites.AgregarObservacion(ctx, t.ID, obs, in.ActorID); err != nil {
		return err
	}
	h := &tramite.Historial{
		Estado: tramite.Observado, EstadoAnterior: t.Estado, UsuarioID: &in.ActorID, Fecha: ahora,
		Detalle: fmt.Sprintf("Observación sobre «%s»: %s", obs.Documento, obs.Motivo),
	}
	cambia := t.Estado != tramite.Observado
	t.Estado = tramite.Observado
	if err := c.Tramites.Guardar(ctx, t, h); err != nil {
		return err
	}
	registrar(ctx, c.Auditoria, &in.ActorID, auditoria.AccionObservar, t.ID, ahora, map[string]any{"documento": obs.Documento})
	if cambia { // más observaciones sobre un trámite ya observado no repiten el aviso
		notificar(ctx, c.Notificador, t)
	}
	return nil
}

type ResolverObservacion struct {
	Tramites  tramite.Repository
	Auditoria auditoria.Repository
	Reloj     ports.Clock
}

// Ejecutar marca la observación como resuelta; al cerrar la última, el trámite vuelve a En Revisión.
func (c *ResolverObservacion) Ejecutar(ctx context.Context, actorID, tramiteID, observacionID string) error {
	t, err := cargarDelAnalista(ctx, c.Tramites, actorID, tramiteID)
	if err != nil {
		return err
	}
	o, err := c.Tramites.BuscarObservacion(ctx, t.ID, observacionID)
	if err != nil {
		return err
	}
	if o.Resuelta {
		return nil
	}
	ahora := c.Reloj.Ahora()
	if err := c.Tramites.ResolverObservacion(ctx, o.ID, ahora); err != nil {
		return err
	}
	pendientes, err := c.Tramites.ContarObservacionesPendientes(ctx, t.ID)
	if err != nil {
		return err
	}
	var h *tramite.Historial
	if pendientes == 0 && t.Estado == tramite.Observado {
		h = &tramite.Historial{Estado: tramite.EnRevision, EstadoAnterior: t.Estado, UsuarioID: &actorID, Detalle: "Observaciones resueltas", Fecha: ahora}
		t.Estado = tramite.EnRevision
	}
	if err := c.Tramites.Guardar(ctx, t, h); err != nil {
		return err
	}
	registrar(ctx, c.Auditoria, &actorID, auditoria.AccionResolverObservacion, t.ID, ahora, map[string]any{"observacionId": o.ID})
	return nil
}

// --- CU-11 Emitir dictamen y archivar ------------------------------------------------------

type Dictaminar struct {
	Tramites  tramite.Repository
	Auditoria auditoria.Repository
	Reloj     ports.Clock
}

type DictamenEntrada struct {
	ActorID   string
	TramiteID string
	Resultado tramite.Resultado
	Texto     string
}

func (c *Dictaminar) Ejecutar(ctx context.Context, in DictamenEntrada) error {
	t, err := cargarDelAnalista(ctx, c.Tramites, in.ActorID, in.TramiteID)
	if err != nil {
		return err
	}
	if t.Estado == tramite.Observado {
		return tramite.ErrObservacionesPendientes
	}
	if err := tramite.Validar(t.Estado, tramite.Dictaminado); err != nil {
		return err
	}
	pendientes, err := c.Tramites.ContarObservacionesPendientes(ctx, t.ID)
	if err != nil {
		return err
	}
	if pendientes > 0 {
		return tramite.ErrObservacionesPendientes
	}
	if t.AnioDesde == 0 {
		return tramite.ErrPeriodoRequerido
	}
	ahora := c.Reloj.Ahora()
	h := &tramite.Historial{Estado: tramite.Dictaminado, EstadoAnterior: t.Estado, UsuarioID: &in.ActorID, Detalle: "Dictamen emitido: " + string(in.Resultado), Fecha: ahora}
	t.Estado, t.Resultado, t.Dictamen, t.DictaminadoEn = tramite.Dictaminado, in.Resultado, strings.TrimSpace(in.Texto), &ahora
	if err := c.Tramites.Guardar(ctx, t, h); err != nil {
		return err
	}
	registrar(ctx, c.Auditoria, &in.ActorID, auditoria.AccionDictaminar, t.ID, ahora, map[string]any{"resultado": string(in.Resultado)})
	return nil
}

// Archivar lo hace el Personal de Archivo cuando el trámite ya fue remitido al Rectorado y las
// carpetas solicitadas volvieron al archivo.
type Archivar struct {
	Tramites    tramite.Repository
	Auditoria   auditoria.Repository
	Reloj       ports.Clock
	Notificador Notificador
}

func (c *Archivar) Ejecutar(ctx context.Context, actorID, id, ubicacion string) error {
	t, err := c.Tramites.BuscarPorID(ctx, id)
	if err != nil {
		return err
	}
	if err := tramite.Validar(t.Estado, tramite.Archivado); err != nil {
		return err
	}
	if t.CarpetasAbiertas > 0 {
		return tramite.ErrCarpetasAbiertas
	}
	ahora := c.Reloj.Ahora()
	ubicacion = strings.TrimSpace(ubicacion)
	h := &tramite.Historial{Estado: tramite.Archivado, EstadoAnterior: t.Estado, UsuarioID: &actorID, Detalle: "Archivado en: " + ubicacion, Fecha: ahora}
	t.Estado, t.ArchivadoEn, t.UbicacionLegajo = tramite.Archivado, &ahora, ubicacion
	if err := c.Tramites.Guardar(ctx, t, h); err != nil {
		return err
	}
	registrar(ctx, c.Auditoria, &actorID, auditoria.AccionArchivar, t.ID, ahora, map[string]any{"ubicacion": ubicacion})
	notificar(ctx, c.Notificador, t)
	return nil
}

// --- CU-16 Dar de baja un documento ---------------------------------------------------------

type DarDeBajaDocumento struct {
	Tramites  tramite.Repository
	Auditoria auditoria.Repository
	Reloj     ports.Clock
}

type BajaDocumentoEntrada struct {
	ActorID     string
	TramiteID   string
	DocumentoID string
	Motivo      string
}

// Ejecutar retira de las búsquedas un documento cargado por error: la baja es lógica, el objeto
// sigue en el almacén y el motivo queda en el documento y en la bitácora (el Director lo restituye
// en CU-17). Solo el analista asignado puede darla y un expediente archivado es de solo lectura.
func (c *DarDeBajaDocumento) Ejecutar(ctx context.Context, in BajaDocumentoEntrada) error {
	motivo := strings.TrimSpace(in.Motivo)
	if motivo == "" {
		return tramite.ErrMotivoBajaRequerido
	}
	t, err := cargarDelAnalista(ctx, c.Tramites, in.ActorID, in.TramiteID)
	if err != nil {
		return err
	}
	d, err := c.Tramites.BuscarDocumento(ctx, t.ID, in.DocumentoID)
	if err != nil {
		// El trámite ya se cargó, así que un "no encontrado" aquí solo puede ser del documento.
		if errors.Is(err, tramite.ErrNoEncontrado) {
			return tramite.ErrDocumentoNoEncontrado
		}
		return err
	}
	if !d.Activo {
		return tramite.ErrDocumentoDeBaja
	}
	ahora := c.Reloj.Ahora()
	if err := c.Tramites.DarDeBajaDocumento(ctx, t.ID, d.ID, motivo, in.ActorID, ahora); err != nil {
		return err
	}
	registrar(ctx, c.Auditoria, &in.ActorID, auditoria.AccionBajaDocumento, t.ID, ahora,
		map[string]any{"documentoId": d.ID, "documento": d.NombreOriginal, "motivo": motivo})
	return nil
}
