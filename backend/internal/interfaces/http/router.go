// Package http arma las rutas de la API (sección 7) sobre los handlers y middlewares.
package http

import (
	"github.com/gofiber/fiber/v2"

	"escalafon/backend/internal/application/ports"
	"escalafon/backend/internal/domain/rol"
	"escalafon/backend/internal/domain/usuario"
	"escalafon/backend/internal/interfaces/http/handlers"
	"escalafon/backend/internal/interfaces/http/middleware"
)

// Dependencias son las piezas ya construidas en cmd/api/main.go (regla de la sección 3: solo
// main.go instancia implementaciones concretas).
type Dependencias struct {
	Auth   *handlers.AuthHandler
	Perfil *handlers.PerfilHandler

	Usuarios  *handlers.UsuariosHandler
	Roles     *handlers.RolesHandler
	Auditoria *handlers.AuditoriaHandler // CU-06

	// CU-07 al CU-13
	Docentes   *handlers.DocentesHandler
	Categorias *handlers.CategoriasHandler
	Tramites   *handlers.TramitesHandler
	Expediente *handlers.ExpedienteHandler // expediente digital del docente
	Documentos *handlers.DocumentosHandler // CU-17 (bajas y restitución) y CU-18 (búsqueda global)

	Legajos *handlers.LegajosHandler // CU-19 al CU-24
	Notificaciones *handlers.NotificacionesHandler // CU-14
	Indicadores    *handlers.IndicadoresHandler    // CU-25 y CU-26

	Tokens        ports.TokenService
	Cache         ports.SessionCache
	UsuarioRepo   usuario.Repository
	LimitadorAuth fiber.Handler // limitador de tasa para /api/auth/*

	LimitadorPublico fiber.Handler // limitador de tasa para /api/publico/*
	// LimitadorSolicitudes acota la radicación desde el portal (sube archivos y crea trámites): mucho más estricto.
	LimitadorSolicitudes fiber.Handler
}

func Registrar(app *fiber.App, d Dependencias) {
	autenticar := middleware.Autenticar(d.Tokens, d.Cache, d.UsuarioRepo)

	api := app.Group("/api")

	// 7.1 Autenticación (pública, con límite de tasa).
	autGrupo := api.Group("/auth")
	if d.LimitadorAuth != nil {
		autGrupo.Use(d.LimitadorAuth)
	}
	autGrupo.Post("/login", d.Auth.Login)
	autGrupo.Post("/recuperacion", d.Auth.Recuperacion)
	autGrupo.Post("/restablecer", d.Auth.Restablecer)
	autGrupo.Post("/activar", d.Auth.Activar)

	// 7.2 Sesión y perfil (requiere autenticación).
	autGrupo.Post("/logout", autenticar, d.Auth.Logout)
	autGrupo.Get("/sesion", autenticar, d.Auth.Sesion)

	perfilGrupo := api.Group("/perfil", autenticar, middleware.RequierePermiso(rol.PermisoPerfilGestionar))
	perfilGrupo.Get("/", d.Perfil.Obtener)
	perfilGrupo.Put("/", d.Perfil.ActualizarContacto)
	perfilGrupo.Put("/contrasena", d.Perfil.CambiarContrasena)

	// 7.3 Usuarios y 7.4 Roles (HU-02): registradas si los handlers están construidos.
	if d.Usuarios != nil {
		usuariosGrupo := api.Group("/usuarios", autenticar, middleware.RequierePermiso(rol.PermisoUsuarioGestionar))
		usuariosGrupo.Get("/", d.Usuarios.Listar)
		usuariosGrupo.Get("/:id", d.Usuarios.Obtener)
		usuariosGrupo.Post("/", d.Usuarios.Crear)
		usuariosGrupo.Put("/:id", d.Usuarios.Editar)
		usuariosGrupo.Patch("/:id/rol", d.Usuarios.CambiarRol)
		usuariosGrupo.Patch("/:id/estado", d.Usuarios.CambiarEstado)
		usuariosGrupo.Post("/:id/reenviar-activacion", d.Usuarios.ReenviarActivacion)
	}
	if d.Roles != nil {
		api.Get("/permisos", autenticar, middleware.RequierePermiso(rol.PermisoRolGestionar), d.Roles.ListarPermisos)
		rolesGrupo := api.Group("/roles", autenticar, middleware.RequierePermiso(rol.PermisoRolGestionar))
		rolesGrupo.Get("/", d.Roles.Listar)
		rolesGrupo.Get("/:id", d.Roles.Obtener)
		rolesGrupo.Post("/", d.Roles.Guardar)
		rolesGrupo.Put("/:id", d.Roles.Guardar)
		rolesGrupo.Delete("/:id", d.Roles.Eliminar)
	}
	// 7.5 Auditoría (CU-06): solo lectura; no hay rutas que escriban en la bitácora.
	if d.Auditoria != nil {
		auditoriaGrupo := api.Group("/auditoria", autenticar, middleware.RequierePermiso(rol.PermisoAuditoriaConsultar))
		auditoriaGrupo.Get("/", d.Auditoria.Listar)
		auditoriaGrupo.Get("/filtros", d.Auditoria.Filtros)
	}

	registrarTramites(api, autenticar, d)
	registrarLegajos(api, autenticar, d)
	if d.Notificaciones != nil {
		api.Get("/notificaciones", autenticar, middleware.RequierePermiso(rol.PermisoNotificacionVer), d.Notificaciones.Listar)
	}
	// CU-25 tablero (solo lectura) y CU-26 reportes (la emisión queda en auditoría).
	if d.Indicadores != nil {
		api.Get("/tablero", autenticar, middleware.RequierePermiso(rol.PermisoTableroVer), d.Indicadores.Tablero)
		api.Post("/reportes", autenticar, middleware.RequierePermiso(rol.PermisoReporteGenerar), d.Indicadores.Reporte)
		if d.Indicadores.CasoInterpretar != nil {
			api.Post("/reportes/interpretar", autenticar, middleware.RequierePermiso(rol.PermisoReporteGenerar), d.Indicadores.InterpretarReporte)
		}
	}
}

// registrarLegajos arma las rutas de inventario (CU-19), solicitudes (CU-20, CU-21) y custodia (CU-22 al CU-24).
func registrarLegajos(api fiber.Router, autenticar fiber.Handler, d Dependencias) {
	if d.Legajos == nil {
		return
	}
	req, alguno := middleware.RequierePermiso, middleware.RequiereAlguno
	h := d.Legajos

	g := api.Group("/legajos", autenticar)
	g.Get("/", alguno(rol.PermisoLegajoGestionar, rol.PermisoCarpetaSolicitar, rol.PermisoCustodiaRegistrar), h.Listar)
	g.Get("/:id", alguno(rol.PermisoLegajoGestionar, rol.PermisoCarpetaSolicitar, rol.PermisoCustodiaRegistrar), h.Obtener)
	g.Post("/", req(rol.PermisoLegajoGestionar), h.Guardar)
	g.Put("/:id", req(rol.PermisoLegajoGestionar), h.Guardar)
	g.Post("/:id/solicitudes", req(rol.PermisoCarpetaSolicitar), h.Solicitar)
	g.Post("/:id/devolver", req(rol.PermisoCustodiaRegistrar), h.Devolver)

	s := api.Group("/solicitudes-carpeta", autenticar)
	s.Get("/", alguno(rol.PermisoLegajoGestionar, rol.PermisoCarpetaSolicitar, rol.PermisoCustodiaRegistrar), h.ListarSolicitudes)
	s.Post("/:id/atender", req(rol.PermisoLegajoGestionar), h.Atender)
	s.Post("/:id/despachar", req(rol.PermisoCustodiaRegistrar), h.Despachar)
	s.Post("/:id/entregar", req(rol.PermisoLegajoGestionar), req(rol.PermisoCustodiaRegistrar), h.Entregar) // aprobar y despachar

	api.Get("/custodia/bitacora", autenticar, req(rol.PermisoBitacoraConsultar), h.Bitacora)
}

// registrarTramites arma las rutas de CU-07 (docentes, categorías), CU-08 al CU-11 (trámites),
// CU-16 (baja de documentos) y el portal público de CU-12 y CU-13.
func registrarTramites(api fiber.Router, autenticar fiber.Handler, d Dependencias) {
	req, alguno := middleware.RequierePermiso, middleware.RequiereAlguno

	if d.Docentes != nil {
		g := api.Group("/docentes", autenticar)
		lectura := alguno(rol.PermisoDocenteGestionar, rol.PermisoTramiteRadicar, rol.PermisoLegajoGestionar, rol.PermisoDocumentoBuscar, rol.PermisoDocumentoEscanear)
		g.Get("/", lectura, d.Docentes.Listar)
		g.Get("/sugerencias", lectura, d.Docentes.Sugerir)
		api.Get("/facultades", autenticar, lectura, d.Docentes.ListarFacultades)
		g.Get("/:id", lectura, d.Docentes.Obtener)
		escritura := req(rol.PermisoDocenteGestionar)
		g.Post("/", escritura, d.Docentes.Guardar)
		g.Put("/:id", escritura, d.Docentes.Guardar)
		g.Patch("/:id/estado", escritura, d.Docentes.CambiarEstado)

		// Expediente digital del docente: Archivo sube los escaneos; Archivo, analistas y Director los consultan.
		if d.Expediente != nil {
			consultar := alguno(rol.PermisoDocumentoBuscar, rol.PermisoDocumentoEscanear)
			g.Get("/:id/documentos", consultar, d.Expediente.Listar)
			g.Get("/:id/documentos/buscar", req(rol.PermisoDocumentoBuscar), d.Expediente.Buscar)
			g.Get("/:id/documentos/:docId", consultar, d.Expediente.Url)
			escanear := req(rol.PermisoDocumentoEscanear)
			g.Post("/:id/documentos/subida", escanear, d.Expediente.PrepararSubida)
			g.Post("/:id/documentos/:docId/confirmar", escanear, d.Expediente.ConfirmarSubida)
			g.Patch("/:id/documentos/:docId/estado", escanear, d.Expediente.DarDeBaja)
		}
	}

	if d.Documentos != nil {
		g := api.Group("/documentos", autenticar)
		g.Get("/buscar", req(rol.PermisoDocumentoBuscar), d.Documentos.Buscar)
		g.Get("/bajas", req(rol.PermisoDocumentoRestituir), d.Documentos.ListarBajas)
		g.Post("/:docId/restituir", req(rol.PermisoDocumentoRestituir), d.Documentos.Restituir)
	}

	if d.Categorias != nil {
		g := api.Group("/categorias", autenticar)
		lectura := alguno(rol.PermisoCategoriaGestionar, rol.PermisoDocenteGestionar, rol.PermisoTramiteRadicar)
		g.Get("/", lectura, d.Categorias.Listar)
		g.Get("/:id", lectura, d.Categorias.Obtener)
		escritura := req(rol.PermisoCategoriaGestionar)
		g.Post("/", escritura, d.Categorias.Guardar)
		g.Put("/:id", escritura, d.Categorias.Guardar)
		g.Patch("/:id/estado", escritura, d.Categorias.CambiarEstado)
		g.Delete("/:id", escritura, d.Categorias.Eliminar)
	}

	if d.Tramites != nil {
		api.Get("/analistas", autenticar, req(rol.PermisoTramiteAsignar), d.Tramites.ListarAnalistas)

		g := api.Group("/tramites", autenticar)
		lectura := alguno(rol.PermisoTramiteRadicar, rol.PermisoTramiteAsignar, rol.PermisoTramiteRevisar, rol.PermisoTramiteArchivar)
		g.Get("/", lectura, d.Tramites.Listar)
		g.Get("/:id", lectura, d.Tramites.Obtener)
		g.Get("/:id/documentos/:docId", lectura, d.Tramites.UrlDocumento)
		g.Get("/:id/documentos/:docId/texto", lectura, d.Tramites.TextoDocumento)
		g.Get("/:id/documentos/:docId/buscar", lectura, d.Tramites.BuscarEnDocumento)

		radicar := req(rol.PermisoTramiteRadicar)
		g.Post("/", radicar, d.Tramites.Radicar)
		g.Post("/:id/validar", radicar, d.Tramites.Validar)
		g.Post("/:id/asignar", req(rol.PermisoTramiteAsignar), d.Tramites.Asignar)

		// CU-16: baja lógica de un adjunto cargado por error (el objeto se conserva en el almacén).
		g.Patch("/:id/documentos/:docId/estado", req(rol.PermisoDocumentoGestionar), d.Tramites.DarDeBajaDocumento)

		// CU-15: subida directa a S3 con URL firmada, confirmación (lanza el OCR) y corrección del texto extraído.
		gestionarDoc := req(rol.PermisoDocumentoGestionar)
		g.Post("/:id/documentos/subida", gestionarDoc, d.Tramites.PrepararSubida)
		g.Post("/:id/documentos/:docId/confirmar", gestionarDoc, d.Tramites.ConfirmarSubida)
		g.Patch("/:id/documentos/:docId/texto", gestionarDoc, d.Tramites.CorregirTexto)
		revisar := req(rol.PermisoTramiteRevisar)
		g.Post("/:id/revision", revisar, d.Tramites.IniciarRevision)
		g.Post("/:id/observaciones", revisar, d.Tramites.RegistrarObservacion)
		g.Post("/:id/observaciones/:oid/resolver", revisar, d.Tramites.ResolverObservacion)
		g.Post("/:id/dictamen", revisar, d.Tramites.Dictaminar)
		g.Put("/:id/periodo", revisar, d.Tramites.RegistrarPeriodo)
		g.Post("/:id/remitir-facultad", revisar, d.Tramites.RemitirFacultad)
		g.Post("/:id/regreso", alguno(rol.PermisoTramiteRevisar, rol.PermisoTramiteRadicar), d.Tramites.RegistrarRegreso)
		g.Post("/:id/remitir-rectorado", revisar, d.Tramites.RemitirRectorado)
		// El archivo físico lo hace el Personal de Archivo, no el analista.
		g.Post("/:id/archivar", req(rol.PermisoTramiteArchivar), d.Tramites.Archivar)

		pub := api.Group("/publico")
		if d.LimitadorPublico != nil {
			pub.Use(d.LimitadorPublico)
		}
		pub.Get("/tramites/:codigo", d.Tramites.ConsultarPublico)
		pub.Post("/consulta-docente", d.Tramites.ConsultarPublicoPorDocente)
		if d.LimitadorSolicitudes != nil {
			pub.Post("/solicitudes", d.LimitadorSolicitudes, d.Tramites.RadicarPublico)
		} else {
			pub.Post("/solicitudes", d.Tramites.RadicarPublico)
		}
		if d.Categorias != nil {
			pub.Get("/categorias", d.Categorias.ListarPublicas)
		}
	}
}
