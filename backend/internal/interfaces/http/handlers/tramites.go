package handlers

import (
	"io"
	"strconv"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"

	"escalafon/backend/internal/application/tramites"
	"escalafon/backend/internal/domain/rol"
	"escalafon/backend/internal/domain/tramite"
	"escalafon/backend/internal/interfaces/http/apierror"
	"escalafon/backend/internal/interfaces/http/dto"
	"escalafon/backend/internal/interfaces/http/middleware"
)

const maxLecturaArchivo = 8<<20 + 1 // un byte más que el máximo para detectar el exceso

// TramitesHandler expone /api/tramites (CU-08 al CU-11), /api/analistas y el portal público
// (CU-12, CU-13).
type TramitesHandler struct {
	CasoRadicar              *tramites.Radicar
	CasoRadicarPublico       *tramites.RadicarPublico
	CasoValidar              *tramites.ValidarSolicitud
	CasoListar               *tramites.Listar
	CasoListarPagina         *tramites.ListarPagina
	CasoObtener              *tramites.Obtener
	CasoListarAnalistas      *tramites.ListarAnalistas
	CasoUrlDocumento         *tramites.UrlDocumento
	CasoConsultarPublico     *tramites.ConsultarPublico
	CasoPublicoPorDocente    *tramites.ConsultarPublicoPorDocente
	CasoAsignar              *tramites.Asignar
	CasoIniciarRevision      *tramites.IniciarRevision
	CasoRegistrarObservacion *tramites.RegistrarObservacion
	CasoResolverObservacion  *tramites.ResolverObservacion
	CasoDictaminar           *tramites.Dictaminar
	CasoArchivar             *tramites.Archivar
	CasoRegistrarPeriodo     *tramites.RegistrarPeriodo
	CasoRemitirFacultad      *tramites.RemitirAFacultad
	CasoRegistrarRegreso     *tramites.RegistrarRegreso
	CasoRemitirRectorado     *tramites.RemitirARectorado
	CasoBajaDocumento        *tramites.DarDeBajaDocumento // CU-16
	CasoPrepararSubida       *tramites.PrepararSubida     // CU-15
	CasoConfirmarSubida      *tramites.ConfirmarSubida
	CasoCorregirTexto        *tramites.CorregirPaginas
	CasoTextoDocumento      *tramites.TextoDocumento
	CasoBuscarEnDocumento    *tramites.BuscarEnDocumento
}

const formatoFecha = "2006-01-02"

func aResumenDTO(r tramite.Resumen) dto.ResumenTramiteDTO {
	return dto.ResumenTramiteDTO{
		ID: r.ID, Codigo: r.Codigo, DocenteNombre: r.DocenteNombre, DocenteCodigo: r.DocenteCodigo,
		Categoria: r.Categoria, Tipo: string(r.Tipo), Estado: string(r.Estado), Origen: string(r.Origen),
		PendienteValidacion: r.PendienteValidacion, AnalistaNombre: r.AnalistaNombre, CreadoEn: r.CreadoEn,
	}
}

func aDocumentoDTO(doc tramite.Documento) dto.DocumentoDTO {
	return dto.DocumentoDTO{
		ID: doc.ID, TipoDocumento: string(doc.TipoDocumento), NombreOriginal: doc.NombreOriginal, TipoMime: doc.TipoMime, Tamano: doc.Tamano, RequisitoID: doc.RequisitoID,
		Activo: doc.Activo, MotivoBaja: doc.MotivoBaja, BajaEn: doc.BajaEn, EstadoOcr: string(doc.EstadoOCR), ErrorOcr: doc.ErrorOCR,
		TextoValidado: doc.TextoValidado, HashSha256: doc.HashSHA256, TamanoComprimido: doc.TamanoComprimido,
	}
}

func aDetalleDTO(t *tramite.Tramite) dto.TramiteDetalleDTO {
	d := dto.TramiteDetalleDTO{
		ResumenTramiteDTO: dto.ResumenTramiteDTO{
			ID: t.ID, Codigo: t.Codigo, DocenteNombre: t.DocenteNombre, DocenteCodigo: t.DocenteCodigo,
			Categoria: t.Categoria, Tipo: string(t.Tipo), Estado: string(t.Estado), Origen: string(t.Origen),
			PendienteValidacion: t.PendienteValidacion, AnalistaNombre: t.AnalistaNombre, CreadoEn: t.CreadoEn,
		},
		DocenteID: t.DocenteID, CategoriaID: t.CategoriaID, AnalistaID: t.AnalistaID, Resultado: string(t.Resultado),
		Dictamen: t.Dictamen, DictaminadoEn: t.DictaminadoEn, ArchivadoEn: t.ArchivadoEn, UbicacionLegajo: t.UbicacionLegajo,
		AnioDesde: t.AnioDesde, AnioHasta: t.AnioHasta, ResolucionNumero: t.ResolucionNumero, RemisionNota: t.RemisionNota,
		CarpetasLigadas: t.CarpetasLigadas, CarpetasAbiertas: t.CarpetasAbiertas,
		Historial: []dto.HistorialDTO{}, Observaciones: []dto.ObservacionDTO{}, Documentos: []dto.DocumentoDTO{}, Requisitos: []dto.RequisitoDTO{},
	}
	if t.ResolucionFecha != nil {
		f := t.ResolucionFecha.Format(formatoFecha)
		d.ResolucionFecha = &f
	}
	for _, h := range t.Historial {
		d.Historial = append(d.Historial, dto.HistorialDTO{
			Estado: string(h.Estado), EstadoAnterior: string(h.EstadoAnterior), Usuario: h.UsuarioNombre, Detalle: h.Detalle, Fecha: h.Fecha,
		})
	}
	for _, o := range t.Observaciones {
		d.Observaciones = append(d.Observaciones, dto.ObservacionDTO{
			ID: o.ID, Documento: o.Documento, Motivo: o.Motivo, Plazo: o.Plazo.Format(formatoFecha), Resuelta: o.Resuelta, ResueltaEn: o.ResueltaEn,
		})
	}
	for _, doc := range t.Documentos {
		d.Documentos = append(d.Documentos, aDocumentoDTO(doc))
	}
	for _, r := range t.Requisitos {
		presente := r.Presente
		d.Requisitos = append(d.Requisitos, dto.RequisitoDTO{ID: r.ID, Nombre: r.Nombre, Presente: &presente})
	}
	return d
}

func (h *TramitesHandler) Listar(c *fiber.Ctx) error {
	actor, _ := middleware.DeContexto(c)
	f := tramite.Filtro{Estado: tramite.Estado(c.Query("estado")), Busqueda: c.Query("buscar"), SinAsignar: c.Query("sinAsignar") == "true"}
	if c.Query("mios") == "true" {
		f.AnalistaID = actor.UsuarioID
	}
	// Con ?pagina=N devuelve {datos, total, pagina, tamano} (?tamano=, 10 por defecto); sin él, el arreglo completo.
	if c.Query("pagina") != "" {
		pagina, _ := strconv.Atoi(c.Query("pagina"))
		tamano, _ := strconv.Atoi(c.Query("tamano", "10"))
		pagina, tamano = max(pagina, 1), min(max(tamano, 1), 100)
		lista, total, err := h.CasoListarPagina.Ejecutar(c.Context(), f, pagina, tamano)
		if err != nil {
			return apierror.Mapear(err)
		}
		out := dto.ListaTramitesDTO{Datos: make([]dto.ResumenTramiteDTO, len(lista)), Total: total, Pagina: pagina, Tamano: tamano}
		for i, r := range lista {
			out.Datos[i] = aResumenDTO(r)
		}
		return c.JSON(out)
	}
	lista, err := h.CasoListar.Ejecutar(c.Context(), f)
	if err != nil {
		return apierror.Mapear(err)
	}
	out := make([]dto.ResumenTramiteDTO, len(lista))
	for i, r := range lista {
		out[i] = aResumenDTO(r)
	}
	return c.JSON(out)
}

func (h *TramitesHandler) Obtener(c *fiber.Ctx) error {
	t, err := h.CasoObtener.Ejecutar(c.Context(), c.Params("id"))
	if err != nil {
		return apierror.Mapear(err)
	}
	return c.JSON(aDetalleDTO(t))
}

func (h *TramitesHandler) Radicar(c *fiber.Ctx) error {
	var in dto.RadicarEntrada
	if err := leerCuerpo(c, &in); err != nil {
		return err
	}
	actor, _ := middleware.DeContexto(c)
	t, err := h.CasoRadicar.Ejecutar(c.Context(), tramites.RadicarEntrada{
		ActorID: actor.UsuarioID, DocenteID: in.DocenteID, Tipo: tramite.Tipo(in.Tipo), CategoriaID: in.CategoriaID,
	})
	if err != nil {
		return apierror.Mapear(err)
	}
	return c.Status(fiber.StatusCreated).JSON(aDetalleDTO(t))
}

func (h *TramitesHandler) Validar(c *fiber.Ctx) error {
	actor, _ := middleware.DeContexto(c)
	if err := h.CasoValidar.Ejecutar(c.Context(), actor.UsuarioID, c.Params("id")); err != nil {
		return apierror.Mapear(err)
	}
	return c.SendStatus(fiber.StatusNoContent)
}

func (h *TramitesHandler) ListarAnalistas(c *fiber.Ctx) error {
	lista, err := h.CasoListarAnalistas.Ejecutar(c.Context())
	if err != nil {
		return apierror.Mapear(err)
	}
	out := make([]dto.AnalistaDTO, len(lista))
	for i, a := range lista {
		out[i] = dto.AnalistaDTO{ID: a.ID, Nombre: a.Nombre, Carga: a.Carga}
	}
	return c.JSON(out)
}

// Asignar atiende tanto la asignación inicial como la reasignación (esta exige motivo).
func (h *TramitesHandler) Asignar(c *fiber.Ctx) error {
	var in dto.AsignarEntrada
	if err := leerCuerpo(c, &in); err != nil {
		return err
	}
	actor, _ := middleware.DeContexto(c)
	t, err := h.CasoAsignar.Ejecutar(c.Context(), tramites.AsignarEntrada{
		ActorID: actor.UsuarioID, TramiteID: c.Params("id"), AnalistaID: in.AnalistaID, Motivo: in.Motivo,
	})
	if err != nil {
		return apierror.Mapear(err)
	}
	return c.JSON(aDetalleDTO(t))
}

func (h *TramitesHandler) IniciarRevision(c *fiber.Ctx) error {
	actor, _ := middleware.DeContexto(c)
	if err := h.CasoIniciarRevision.Ejecutar(c.Context(), actor.UsuarioID, c.Params("id")); err != nil {
		return apierror.Mapear(err)
	}
	return c.SendStatus(fiber.StatusNoContent)
}

func (h *TramitesHandler) RegistrarObservacion(c *fiber.Ctx) error {
	var in dto.ObservacionEntrada
	if err := leerCuerpo(c, &in); err != nil {
		return err
	}
	plazo, err := time.Parse(formatoFecha, in.Plazo)
	if err != nil {
		return apierror.Validacion(map[string]string{"Plazo": "Debe ser una fecha AAAA-MM-DD"})
	}
	actor, _ := middleware.DeContexto(c)
	if err := h.CasoRegistrarObservacion.Ejecutar(c.Context(), tramites.ObservacionEntrada{
		ActorID: actor.UsuarioID, TramiteID: c.Params("id"), Documento: in.Documento, Motivo: in.Motivo, Plazo: plazo,
	}); err != nil {
		return apierror.Mapear(err)
	}
	return c.SendStatus(fiber.StatusCreated)
}

func (h *TramitesHandler) ResolverObservacion(c *fiber.Ctx) error {
	actor, _ := middleware.DeContexto(c)
	if err := h.CasoResolverObservacion.Ejecutar(c.Context(), actor.UsuarioID, c.Params("id"), c.Params("oid")); err != nil {
		return apierror.Mapear(err)
	}
	return c.SendStatus(fiber.StatusNoContent)
}

func (h *TramitesHandler) Dictaminar(c *fiber.Ctx) error {
	var in dto.DictamenEntrada
	if err := leerCuerpo(c, &in); err != nil {
		return err
	}
	actor, _ := middleware.DeContexto(c)
	if err := h.CasoDictaminar.Ejecutar(c.Context(), tramites.DictamenEntrada{
		ActorID: actor.UsuarioID, TramiteID: c.Params("id"), Resultado: tramite.Resultado(in.Resultado), Texto: in.Texto,
	}); err != nil {
		return apierror.Mapear(err)
	}
	return c.SendStatus(fiber.StatusNoContent)
}

func (h *TramitesHandler) Archivar(c *fiber.Ctx) error {
	var in dto.ArchivarEntrada
	if err := leerCuerpo(c, &in); err != nil {
		return err
	}
	actor, _ := middleware.DeContexto(c)
	if err := h.CasoArchivar.Ejecutar(c.Context(), actor.UsuarioID, c.Params("id"), in.Ubicacion); err != nil {
		return apierror.Mapear(err)
	}
	return c.SendStatus(fiber.StatusNoContent)
}

func (h *TramitesHandler) RegistrarPeriodo(c *fiber.Ctx) error {
	var in dto.PeriodoEntrada
	if err := leerCuerpo(c, &in); err != nil {
		return err
	}
	actor, _ := middleware.DeContexto(c)
	if err := h.CasoRegistrarPeriodo.Ejecutar(c.Context(), actor.UsuarioID, c.Params("id"), in.AnioDesde, in.AnioHasta); err != nil {
		return apierror.Mapear(err)
	}
	return c.SendStatus(fiber.StatusNoContent)
}

// RemitirFacultad registra la salida del expediente: desde ahí el trámite no se puede rastrear.
func (h *TramitesHandler) RemitirFacultad(c *fiber.Ctx) error {
	actor, _ := middleware.DeContexto(c)
	if err := h.CasoRemitirFacultad.Ejecutar(c.Context(), actor.UsuarioID, c.Params("id")); err != nil {
		return apierror.Mapear(err)
	}
	return c.SendStatus(fiber.StatusNoContent)
}

// RegistrarRegreso lo puede hacer el analista asignado o quien radica trámites (ventanilla).
func (h *TramitesHandler) RegistrarRegreso(c *fiber.Ctx) error {
	var in dto.RegresoEntrada
	if err := leerCuerpo(c, &in); err != nil {
		return err
	}
	fecha, err := time.Parse(formatoFecha, in.ResolucionFecha)
	if err != nil {
		return apierror.Validacion(map[string]string{"ResolucionFecha": "Debe ser una fecha AAAA-MM-DD"})
	}
	actor, _ := middleware.DeContexto(c)
	if err := h.CasoRegistrarRegreso.Ejecutar(c.Context(), tramites.RegresoEntrada{
		ActorID: actor.UsuarioID, EsVentanilla: actor.TienePermiso(rol.PermisoTramiteRadicar), TramiteID: c.Params("id"),
		ResolucionNumero: in.ResolucionNumero, ResolucionFecha: fecha,
	}); err != nil {
		return apierror.Mapear(err)
	}
	return c.SendStatus(fiber.StatusNoContent)
}

func (h *TramitesHandler) RemitirRectorado(c *fiber.Ctx) error {
	var in dto.RemisionEntrada
	if err := leerCuerpo(c, &in); err != nil {
		return err
	}
	actor, _ := middleware.DeContexto(c)
	if err := h.CasoRemitirRectorado.Ejecutar(c.Context(), actor.UsuarioID, c.Params("id"), in.Nota); err != nil {
		return apierror.Mapear(err)
	}
	return c.SendStatus(fiber.StatusNoContent)
}

// DarDeBajaDocumento retira del expediente un adjunto cargado por error (CU-16): la baja es lógica
// y exige un motivo, que queda en el documento y en la bitácora.
func (h *TramitesHandler) DarDeBajaDocumento(c *fiber.Ctx) error {
	var in dto.BajaDocumentoEntrada
	if err := leerCuerpo(c, &in); err != nil {
		return err
	}
	actor, _ := middleware.DeContexto(c)
	if err := h.CasoBajaDocumento.Ejecutar(c.Context(), tramites.BajaDocumentoEntrada{
		ActorID: actor.UsuarioID, TramiteID: c.Params("id"), DocumentoID: c.Params("docId"), Motivo: in.Motivo,
	}); err != nil {
		return apierror.Mapear(err)
	}
	return c.SendStatus(fiber.StatusNoContent)
}

// UrlDocumento responde con un enlace temporal de descarga del adjunto.
func (h *TramitesHandler) UrlDocumento(c *fiber.Ctx) error {
	actor, _ := middleware.DeContexto(c)
	url, err := h.CasoUrlDocumento.Ejecutar(c.Context(), c.Params("id"), c.Params("docId"), actor.TienePermiso(rol.PermisoDocumentoRestituir))
	if err != nil {
		return apierror.Mapear(err)
	}
	return c.JSON(fiber.Map{"url": url})
}

// PrepararSubida registra el documento (CU-15) y entrega la URL firmada para subirlo directo al almacén.
func (h *TramitesHandler) PrepararSubida(c *fiber.Ctx) error {
	var in dto.SubidaEntrada
	if err := leerCuerpo(c, &in); err != nil {
		return err
	}
	actor, _ := middleware.DeContexto(c)
	r, err := h.CasoPrepararSubida.Ejecutar(c.Context(), tramites.SubidaEntrada{
		ActorID: actor.UsuarioID, TramiteID: c.Params("id"), Nombre: in.Nombre, TipoMime: in.TipoMime, Tamano: in.Tamano, RequisitoID: in.RequisitoID,
	})
	if err != nil {
		return apierror.Mapear(err)
	}
	return c.Status(fiber.StatusCreated).JSON(dto.SubidaPreparadaDTO{Documento: aDocumentoDTO(r.Documento), URL: r.URL})
}

// ConfirmarSubida verifica el objeto subido y lanza el OCR; también reintenta un documento en ERROR.
func (h *TramitesHandler) ConfirmarSubida(c *fiber.Ctx) error {
	actor, _ := middleware.DeContexto(c)
	d, err := h.CasoConfirmarSubida.Ejecutar(c.Context(), actor.UsuarioID, c.Params("id"), c.Params("docId"))
	if err != nil {
		return apierror.Mapear(err)
	}
	return c.JSON(aDocumentoDTO(*d))
}

func (h *TramitesHandler) TextoDocumento(c *fiber.Ctx) error {
	actor, _ := middleware.DeContexto(c)
	paginas, err := h.CasoTextoDocumento.Ejecutar(c.Context(), c.Params("id"), c.Params("docId"), actor.TienePermiso(rol.PermisoDocumentoRestituir))
	if err != nil {
		return apierror.Mapear(err)
	}
	return c.JSON(dto.TextoDocumentoDTO{Paginas: paginas})
}

// BuscarEnDocumento busca en el texto OCR indexado y responde las páginas donde aparece lo buscado.
func (h *TramitesHandler) BuscarEnDocumento(c *fiber.Ctx) error {
	actor, _ := middleware.DeContexto(c)
	lista, err := h.CasoBuscarEnDocumento.Ejecutar(c.Context(), c.Params("id"), c.Params("docId"), c.Query("q"), actor.TienePermiso(rol.PermisoDocumentoRestituir))
	if err != nil {
		return apierror.Mapear(err)
	}
	out := dto.BusquedaDocumentoDTO{Coincidencias: make([]dto.CoincidenciaDTO, len(lista))}
	for i, m := range lista {
		out.Coincidencias[i] = dto.CoincidenciaDTO{Pagina: m.Pagina, Fragmento: m.Fragmento, Rango: m.Rango}
	}
	return c.JSON(out)
}

func (h *TramitesHandler) CorregirTexto(c *fiber.Ctx) error {
	var in dto.TextoDocumentoDTO
	if err := leerCuerpo(c, &in); err != nil {
		return err
	}
	actor, _ := middleware.DeContexto(c)
	if err := h.CasoCorregirTexto.Ejecutar(c.Context(), actor.UsuarioID, c.Params("id"), c.Params("docId"), in.Paginas); err != nil {
		return apierror.Mapear(err)
	}
	return c.SendStatus(fiber.StatusNoContent)
}

// --- Portal público ------------------------------------------------------------------

// ConsultarPublicoPorDocente lista los trámites de un docente a partir de su código y su cédula (portal público).
// Va por POST para que la cédula no viaje en la URL.
func (h *TramitesHandler) ConsultarPublicoPorDocente(c *fiber.Ctx) error {
	var in dto.ConsultaDocenteEntrada
	if err := leerCuerpo(c, &in); err != nil {
		return err
	}
	lista, err := h.CasoPublicoPorDocente.Ejecutar(c.Context(), in.Codigo, in.Cedula)
	if err != nil {
		return apierror.Mapear(err)
	}
	out := make([]dto.ResumenPublicoDTO, len(lista))
	for i, t := range lista {
		out[i] = dto.ResumenPublicoDTO{Codigo: t.Codigo, Tipo: string(t.Tipo), Categoria: t.Categoria, Estado: string(t.Estado), RadicadoEn: t.RadicadoEn, PendienteValidacion: t.PendienteValidacion}
	}
	return c.JSON(out)
}

func (h *TramitesHandler) ConsultarPublico(c *fiber.Ctx) error {
	v, err := h.CasoConsultarPublico.Ejecutar(c.Context(), c.Params("codigo"))
	if err != nil {
		return apierror.Mapear(err)
	}
	out := dto.VistaPublicaDTO{
		Codigo: v.Codigo, Tipo: string(v.Tipo), Categoria: v.Categoria, Estado: string(v.Estado), RadicadoEn: v.RadicadoEn,
		PendienteValidacion: v.PendienteValidacion, Resultado: string(v.Resultado),
		Etapas: []dto.EtapaPublicaDTO{}, Observaciones: []dto.ObservacionPublicaDTO{},
	}
	for _, e := range v.Etapas {
		out.Etapas = append(out.Etapas, dto.EtapaPublicaDTO{Estado: string(e.Estado), Fecha: e.Fecha})
	}
	for _, o := range v.Observaciones {
		out.Observaciones = append(out.Observaciones, dto.ObservacionPublicaDTO{
			Documento: o.Documento, Motivo: o.Motivo, Plazo: o.Plazo.Format(formatoFecha), Resuelta: o.Resuelta,
		})
	}
	return c.JSON(out)
}

// RadicarPublico recibe multipart/form-data: codigo y cedula del docente, tipo y, en categorización, categoriaId y un
// archivo por requisito en el campo doc_<id del requisito>; en regularización, de 1 a 5 archivos en el campo adjunto.
func (h *TramitesHandler) RadicarPublico(c *fiber.Ctx) error {
	form, err := c.MultipartForm()
	if err != nil {
		return apierror.Nuevo(400, "CUERPO_INVALIDO", "La solicitud debe enviarse como formulario multipart")
	}
	valor := func(k string) string {
		if v := form.Value[k]; len(v) > 0 {
			return strings.TrimSpace(v[0])
		}
		return ""
	}
	tipo := tramite.Tipo(valor("tipo"))
	campos := map[string]string{}
	codigoDocente, cedula := valor("codigo"), valor("cedula")
	if codigoDocente == "" {
		campos["codigo"] = "Este campo es obligatorio"
	}
	if cedula == "" {
		campos["cedula"] = "Este campo es obligatorio"
	}
	catID := 0
	if tipo != tramite.Regularizacion { // la regularización no postula a ninguna categoría
		if catID, err = strconv.Atoi(valor("categoriaId")); err != nil {
			campos["categoriaId"] = "Selecciona una categoría"
		}
	}
	if len(campos) > 0 {
		return apierror.Validacion(campos)
	}

	var archivos []tramites.Archivo
	for clave, headers := range form.File {
		reqID := 0 // adjunto libre
		if clave != "adjunto" {
			reqID, err = strconv.Atoi(strings.TrimPrefix(clave, "doc_"))
			if err != nil || reqID <= 0 || !strings.HasPrefix(clave, "doc_") {
				return apierror.Nuevo(422, "ARCHIVO_INVALIDO", "Adjunto con nombre de campo inválido")
			}
			headers = headers[:min(len(headers), 1)] // un archivo por requisito
		}
		for _, fh := range headers {
			f, err := fh.Open()
			if err != nil {
				return apierror.Nuevo(422, "ARCHIVO_INVALIDO", "No se pudo leer un adjunto")
			}
			datos, err := io.ReadAll(io.LimitReader(f, maxLecturaArchivo))
			f.Close()
			if err != nil {
				return apierror.Nuevo(422, "ARCHIVO_INVALIDO", "No se pudo leer un adjunto")
			}
			archivos = append(archivos, tramites.Archivo{RequisitoID: int32(reqID), Nombre: fh.Filename, Datos: datos})
		}
	}

	codigo, err := h.CasoRadicarPublico.Ejecutar(c.Context(), tramites.RadicarPublicoEntrada{
		Codigo: codigoDocente, Cedula: cedula, Tipo: tipo, CategoriaID: int16(catID), Archivos: archivos,
	})
	if err != nil {
		return apierror.Mapear(err)
	}
	return c.Status(fiber.StatusCreated).JSON(dto.SolicitudRadicadaDTO{Codigo: codigo})
}
