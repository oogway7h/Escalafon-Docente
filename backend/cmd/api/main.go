// Comando api: compone las dependencias e inicia el servidor Fiber (sección 3, regla 5:
// este es el único lugar donde se instancian implementaciones concretas).
package main

import (
	"context"
	"log"
	"log/slog"
	"os"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/compress"
	"github.com/gofiber/fiber/v2/middleware/helmet"
	"github.com/gofiber/fiber/v2/middleware/limiter"
	fiberlog "github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/requestid"
	redisStorage "github.com/gofiber/storage/redis/v3"

	auditoriasApp "escalafon/backend/internal/application/auditorias"
	"escalafon/backend/internal/application/auth"
	"escalafon/backend/internal/application/ports"
	categoriasApp "escalafon/backend/internal/application/categorias"
	docentesApp "escalafon/backend/internal/application/docentes"
	indicadoresApp "escalafon/backend/internal/application/indicadores"
	"escalafon/backend/internal/application/perfil"
	rolesApp "escalafon/backend/internal/application/roles"
	legajosApp "escalafon/backend/internal/application/legajos"
	notificacionesApp "escalafon/backend/internal/application/notificaciones"
	tramitesApp "escalafon/backend/internal/application/tramites"
	usuariosApp "escalafon/backend/internal/application/usuarios"
	"escalafon/backend/internal/domain/token"
	"escalafon/backend/internal/infrastructure/cache"
	"escalafon/backend/internal/infrastructure/config"
	"escalafon/backend/internal/infrastructure/iaservice"
	"escalafon/backend/internal/infrastructure/mail"
	"escalafon/backend/internal/infrastructure/persistence"
	"escalafon/backend/internal/infrastructure/llm"
	"escalafon/backend/internal/infrastructure/reportes"
	"escalafon/backend/internal/infrastructure/security"
	"escalafon/backend/internal/infrastructure/storage"
	httpapi "escalafon/backend/internal/interfaces/http"
	"escalafon/backend/internal/interfaces/http/apierror"
	"escalafon/backend/internal/interfaces/http/handlers"
	"escalafon/backend/internal/interfaces/http/middleware"
)

func main() {
	if os.Getenv("APP_ENV") == "production" {
		slog.SetDefault(slog.New(slog.NewJSONHandler(os.Stdout, nil)))
	}

	cfg, err := config.Cargar()
	if err != nil {
		log.Fatalf("configuración inválida: %v", err)
	}

	db, err := persistence.Conectar(cfg.DatabaseURL)
	if err != nil {
		log.Fatalf("no se pudo conectar a la base de datos: %v", err)
	}

	sessionCache, err := cache.NewSessionCacheRedis(cfg.RedisURL)
	if err != nil {
		log.Fatalf("no se pudo crear el cliente de Redis: %v", err)
	}

	// --- Adaptadores de infraestructura ---
	usuarioRepo := persistence.NewRepositorioUsuario(db)
	rolRepo := persistence.NewRepositorioRol(db)
	tokenRepo := persistence.NewRepositorioToken(db)
	auditoria := persistence.NewRepositorioAuditoria(db)
	docenteRepo := persistence.NewRepositorioDocente(db)
	categoriaRepo := persistence.NewRepositorioCategoria(db)
	tramiteRepo := persistence.NewRepositorioTramite(db)
	legajoRepo := persistence.NewRepositorioLegajo(db)
	notificacionRepo := persistence.NewRepositorioNotificacion(db)
	indicadorRepo := persistence.NewRepositorioIndicador(db)
	tx := persistence.NewTxManager(db)

	hasher := security.NewBcryptHasher()
	jwtServicio := security.NewServicioJWT(cfg.JWTSecreto, cfg.JWTEmisor)
	generador := security.NewGeneradorAleatorio()
	mailer := mail.NewMailerSMTP(cfg.SMTPHost, cfg.SMTPPuerto, cfg.SMTPUsuario, cfg.SMTPClave, cfg.SMTPRemitente)
	reloj := relojReal{}

	var almacen ports.ObjectStorage = storage.NoConfigurado{}
	if cfg.S3Bucket != "" {
		almacen = storage.NewS3(storage.ConfigS3{
			Bucket: cfg.S3Bucket, Region: cfg.S3Region, AccessKey: cfg.S3AccessKey, SecretKey: cfg.S3SecretKey, Endpoint: cfg.S3Endpoint,
		})
	} else {
		slog.Warn("S3_BUCKET no configurado: los adjuntos de trámites estarán deshabilitados")
	}

	// --- Casos de uso de HU-01 ---
	autenticarCaso := &auth.AutenticarUsuario{
		Usuarios: usuarioRepo, Roles: rolRepo, Auditoria: auditoria, Hasher: hasher, Tokens: jwtServicio,
		Reloj: reloj, MaxIntentos: cfg.MaxIntentosLogin, MinutosBloqueo: cfg.MinutosBloqueo, SesionHoras: cfg.SesionHoras,
	}
	cerrarSesionCaso := &auth.CerrarSesion{Cache: sessionCache}
	obtenerSesionCaso := &auth.ObtenerSesion{Usuarios: usuarioRepo, Roles: rolRepo}
	aceptarTerminosCaso := &auth.AceptarTerminos{Usuarios: usuarioRepo, Auditoria: auditoria, Reloj: reloj}
	solicitarRecuperacionCaso := &auth.SolicitarRecuperacion{
		Usuarios: usuarioRepo, Tokens: tokenRepo, Generador: generador, Mailer: mailer, Reloj: reloj,
		MinutosVigencia: cfg.MinutosTokenRecuperacion, URLFrontend: cfg.URLFrontend, Asincrono: true,
	}
	restablecerRecuperacionCaso := &auth.RestablecerContrasena{
		Usuarios: usuarioRepo, Tokens: tokenRepo, Auditoria: auditoria, Hasher: hasher, Generador: generador,
		Cache: sessionCache, Tx: tx, Reloj: reloj, PropositoEsperado: token.PropositoRecuperacion,
	}
	restablecerActivacionCaso := &auth.RestablecerContrasena{
		Usuarios: usuarioRepo, Tokens: tokenRepo, Auditoria: auditoria, Hasher: hasher, Generador: generador,
		Cache: sessionCache, Tx: tx, Reloj: reloj, PropositoEsperado: token.PropositoActivacion,
	}
	fotoRepo := persistence.NewRepositorioFoto(db)
	obtenerPerfilCaso := &perfil.ObtenerPerfil{Usuarios: usuarioRepo}
	actualizarContactoCaso := &perfil.ActualizarContacto{Usuarios: usuarioRepo, Auditoria: auditoria, Reloj: reloj}
	cambiarContrasenaCaso := &perfil.CambiarContrasena{
		Usuarios: usuarioRepo, Roles: rolRepo, Auditoria: auditoria, Hasher: hasher, Tokens: jwtServicio,
		Cache: sessionCache, Tx: tx, Reloj: reloj, SesionHoras: cfg.SesionHoras,
	}

	authHandler := &handlers.AuthHandler{
		Autenticar: autenticarCaso, CerrarSesionCaso: cerrarSesionCaso, ObtenerSesionCaso: obtenerSesionCaso, AceptarTerminosCaso: aceptarTerminosCaso,
		SolicitarRecuperacion: solicitarRecuperacionCaso, RestablecerRecuperacion: restablecerRecuperacionCaso,
		RestablecerActivacion: restablecerActivacionCaso,
	}
	perfilHandler := &handlers.PerfilHandler{
		CasoObtener: obtenerPerfilCaso, CasoActualizarContacto: actualizarContactoCaso, CasoCambiarContrasena: cambiarContrasenaCaso,
		CasoObtenerFoto: &perfil.ObtenerFoto{Fotos: fotoRepo, Almacen: almacen},
		CasoSubirFoto:   &perfil.SubirFoto{Fotos: fotoRepo, Almacen: almacen, Auditoria: auditoria, Reloj: reloj},
		CasoQuitarFoto:  &perfil.QuitarFoto{Fotos: fotoRepo, Almacen: almacen, Auditoria: auditoria, Reloj: reloj},
	}

	// --- Casos de uso de HU-02 ---
	usuariosHandler := &handlers.UsuariosHandler{
		CasoListar:  &usuariosApp.Listar{Usuarios: usuarioRepo},
		CasoObtener: &usuariosApp.Obtener{Usuarios: usuarioRepo},
		CasoCrear: &usuariosApp.Crear{
			Usuarios: usuarioRepo, Roles: rolRepo, Tokens: tokenRepo, Auditoria: auditoria, Hasher: hasher,
			Generador: generador, Mailer: mailer, Tx: tx, Reloj: reloj,
			HorasTokenActivacion: cfg.HorasTokenActivacion, URLFrontend: cfg.URLFrontend,
		},
		CasoEditar: &usuariosApp.Editar{Usuarios: usuarioRepo, Auditoria: auditoria, Reloj: reloj},
		CasoCambiarRol: &usuariosApp.CambiarRol{
			Usuarios: usuarioRepo, Roles: rolRepo, Auditoria: auditoria, Cache: sessionCache, Tx: tx, Reloj: reloj,
		},
		CasoCambiarEstado: &usuariosApp.CambiarEstado{
			Usuarios: usuarioRepo, Auditoria: auditoria, Cache: sessionCache, Tx: tx, Reloj: reloj,
		},
		CasoReenviarActivacion: &usuariosApp.ReenviarActivacion{
			Usuarios: usuarioRepo, Tokens: tokenRepo, Auditoria: auditoria, Generador: generador, Mailer: mailer, Reloj: reloj,
			HorasTokenActivacion: cfg.HorasTokenActivacion, URLFrontend: cfg.URLFrontend,
		},
	}
	rolesHandler := &handlers.RolesHandler{
		CasoListar:         &rolesApp.Listar{Roles: rolRepo},
		CasoObtener:        &rolesApp.Obtener{Roles: rolRepo},
		CasoListarPermisos: &rolesApp.ListarPermisos{Roles: rolRepo},
		CasoGuardar:        &rolesApp.Guardar{Roles: rolRepo, Usuarios: usuarioRepo, Auditoria: auditoria, Cache: sessionCache, Tx: tx, Reloj: reloj},
		CasoEliminar:       &rolesApp.Eliminar{Roles: rolRepo, Auditoria: auditoria, Reloj: reloj},
	}
	// CU-06: consulta de la bitácora (el mismo repositorio, a través del puerto de solo lectura).
	auditoriaHandler := &handlers.AuditoriaHandler{
		CasoConsultar: &auditoriasApp.Consultar{Auditoria: auditoria},
		CasoOpciones:  &auditoriasApp.ListarOpciones{Auditoria: auditoria},
	}

	// --- CU-07 al CU-13: docentes, categorías y trámites ---
	docentesHandler := &handlers.DocentesHandler{
		CasoListar:        &docentesApp.Listar{Docentes: docenteRepo},
		CasoListarPagina:  &docentesApp.ListarPagina{Docentes: docenteRepo},
		CasoSugerir:       &docentesApp.Sugerir{Docentes: docenteRepo},
		CasoObtener:       &docentesApp.Obtener{Docentes: docenteRepo},
		CasoGuardar:       &docentesApp.Guardar{Docentes: docenteRepo, Categorias: categoriaRepo, Auditoria: auditoria, Reloj: reloj},
		CasoCambiarEstado: &docentesApp.CambiarEstado{Docentes: docenteRepo, Auditoria: auditoria, Reloj: reloj},
		CasoFacultades:    &docentesApp.ListarFacultades{Docentes: docenteRepo},
	}
	categoriasHandler := &handlers.CategoriasHandler{
		CasoListar:        &categoriasApp.Listar{Categorias: categoriaRepo},
		CasoObtener:       &categoriasApp.Obtener{Categorias: categoriaRepo},
		CasoGuardar:       &categoriasApp.Guardar{Categorias: categoriaRepo, Auditoria: auditoria, Reloj: reloj},
		CasoCambiarEstado: &categoriasApp.CambiarEstado{Categorias: categoriaRepo, Auditoria: auditoria, Reloj: reloj},
		CasoEliminar:      &categoriasApp.Eliminar{Categorias: categoriaRepo, Auditoria: auditoria, Reloj: reloj},
	}
	// CU-14: aviso por correo al docente en Observado, En facultad, Remitido y Archivado; los fallos se reintentan.
	notificador := &notificacionesApp.Notificar{
		Docentes: docenteRepo, Envios: notificacionRepo, Mailer: mailer, Reloj: reloj, URLFrontend: cfg.URLFrontend, Asincrono: true,
	}
	go func() {
		for range time.Tick(time.Minute) {
			if err := notificador.ReintentarPendientes(context.Background()); err != nil {
				slog.Warn("reintento de notificaciones falló", "error", err)
			}
		}
	}()
	notificacionesHandler := &handlers.NotificacionesHandler{CasoListar: &notificacionesApp.Listar{Envios: notificacionRepo}}
	// CU-15: OCR y compresión en el ia-service; sin IA_SERVICE_URL los documentos se suben pero no se procesan.
	var procesador ports.ProcesadorDocumentos = iaservice.NoConfigurado{}
	if cfg.IAServiceURL != "" {
		procesador = iaservice.New(cfg.IAServiceURL, cfg.IAServiceToken)
	} else {
		slog.Warn("IA_SERVICE_URL no configurado: el OCR de documentos estará deshabilitado")
	}
	procesoOCR := &tramitesApp.ProcesoOCR{Tramites: tramiteRepo, Procesador: procesador, Asincrono: true}
	confirmarSubida := &tramitesApp.ConfirmarSubida{Tramites: tramiteRepo, Almacen: almacen, OCR: procesoOCR, Auditoria: auditoria, Reloj: reloj}
	if err := procesoOCR.Reanudar(context.Background()); err != nil {
		slog.Warn("no se pudieron reanudar los documentos en proceso", "error", err)
	}
	documentosHandler := &handlers.DocumentosHandler{
		CasoBuscar:    &tramitesApp.BuscarDocumentos{Tramites: tramiteRepo},
		CasoBajas:     &tramitesApp.ListarDocumentosDeBaja{Tramites: tramiteRepo},
		CasoRestituir: &tramitesApp.RestituirDocumento{Tramites: tramiteRepo, Auditoria: auditoria, Reloj: reloj},
	}
	expedienteHandler := &handlers.ExpedienteHandler{
		CasoPrepararSubida:  &tramitesApp.PrepararSubidaDocente{Docentes: docenteRepo, Tramites: tramiteRepo, Almacen: almacen},
		CasoConfirmarSubida: confirmarSubida,
		CasoListar:          &tramitesApp.ListarExpediente{Docentes: docenteRepo, Tramites: tramiteRepo},
		CasoBuscar:          &tramitesApp.BuscarEnExpediente{Docentes: docenteRepo, Tramites: tramiteRepo},
		CasoUrl:             &tramitesApp.UrlExpediente{Tramites: tramiteRepo, Almacen: almacen},
		CasoBaja:            &tramitesApp.DarDeBajaExpediente{Tramites: tramiteRepo, Auditoria: auditoria, Reloj: reloj},
	}
	tramitesHandler := &handlers.TramitesHandler{
		CasoRadicar: &tramitesApp.Radicar{Tramites: tramiteRepo, Docentes: docenteRepo, Categorias: categoriaRepo, Auditoria: auditoria, Reloj: reloj},
		CasoRadicarPublico: &tramitesApp.RadicarPublico{
			Tramites: tramiteRepo, Docentes: docenteRepo, Categorias: categoriaRepo, Almacen: almacen, Auditoria: auditoria, Reloj: reloj,
		},
		// El OCR de los adjuntos del portal se lanza al validar en ventanilla, no al radicar.
		CasoValidar:              &tramitesApp.ValidarSolicitud{Tramites: tramiteRepo, Auditoria: auditoria, Reloj: reloj, OCR: procesoOCR},
		CasoListar:               &tramitesApp.Listar{Tramites: tramiteRepo},
		CasoListarPagina:         &tramitesApp.ListarPagina{Tramites: tramiteRepo},
		CasoObtener:              &tramitesApp.Obtener{Tramites: tramiteRepo},
		CasoListarAnalistas:      &tramitesApp.ListarAnalistas{Tramites: tramiteRepo},
		CasoUrlDocumento:         &tramitesApp.UrlDocumento{Tramites: tramiteRepo, Almacen: almacen},
		CasoConsultarPublico:     &tramitesApp.ConsultarPublico{Tramites: tramiteRepo},
		CasoPublicoPorDocente:    &tramitesApp.ConsultarPublicoPorDocente{Docentes: docenteRepo, Tramites: tramiteRepo},
		CasoAsignar:              &tramitesApp.Asignar{Tramites: tramiteRepo, Auditoria: auditoria, Reloj: reloj},
		CasoIniciarRevision:      &tramitesApp.IniciarRevision{Tramites: tramiteRepo, Auditoria: auditoria, Reloj: reloj},
		CasoRegistrarObservacion: &tramitesApp.RegistrarObservacion{Tramites: tramiteRepo, Auditoria: auditoria, Reloj: reloj, Notificador: notificador},
		CasoResolverObservacion:  &tramitesApp.ResolverObservacion{Tramites: tramiteRepo, Auditoria: auditoria, Reloj: reloj},
		CasoDictaminar:           &tramitesApp.Dictaminar{Tramites: tramiteRepo, Auditoria: auditoria, Reloj: reloj},
		CasoRegistrarPeriodo:     &tramitesApp.RegistrarPeriodo{Tramites: tramiteRepo, Auditoria: auditoria, Reloj: reloj},
		CasoRemitirFacultad:      &tramitesApp.RemitirAFacultad{Tramites: tramiteRepo, Auditoria: auditoria, Reloj: reloj, Notificador: notificador},
		CasoRegistrarRegreso:     &tramitesApp.RegistrarRegreso{Tramites: tramiteRepo, Auditoria: auditoria, Reloj: reloj},
		CasoRemitirRectorado:     &tramitesApp.RemitirARectorado{Tramites: tramiteRepo, Auditoria: auditoria, Reloj: reloj, Notificador: notificador},
		CasoArchivar:             &tramitesApp.Archivar{Tramites: tramiteRepo, Auditoria: auditoria, Reloj: reloj, Notificador: notificador},
		CasoBajaDocumento:        &tramitesApp.DarDeBajaDocumento{Tramites: tramiteRepo, Auditoria: auditoria, Reloj: reloj},
		CasoPrepararSubida:       &tramitesApp.PrepararSubida{Tramites: tramiteRepo, Almacen: almacen, Reloj: reloj},
		CasoConfirmarSubida:      confirmarSubida,
		CasoCorregirTexto:        &tramitesApp.CorregirPaginas{Tramites: tramiteRepo, Auditoria: auditoria, Tx: tx, Reloj: reloj},
		CasoTextoDocumento:       &tramitesApp.TextoDocumento{Tramites: tramiteRepo},
		CasoBuscarEnDocumento:    &tramitesApp.BuscarEnDocumento{Tramites: tramiteRepo},
	}

	// --- CU-19 al CU-24: legajos y custodia ---
	atender := &legajosApp.Atender{Legajos: legajoRepo, Auditoria: auditoria, Reloj: reloj}
	despachar := &legajosApp.Despachar{Legajos: legajoRepo, Auditoria: auditoria, Tx: tx, Reloj: reloj}
	legajosHandler := &handlers.LegajosHandler{
		CasoListar:            &legajosApp.ListarLegajos{Legajos: legajoRepo},
		CasoObtener:           &legajosApp.ObtenerLegajo{Legajos: legajoRepo},
		CasoGuardar:           &legajosApp.GuardarLegajo{Legajos: legajoRepo, Auditoria: auditoria, Reloj: reloj},
		CasoSolicitar:         &legajosApp.Solicitar{Legajos: legajoRepo, Tramites: tramiteRepo, Auditoria: auditoria, Reloj: reloj},
		CasoListarSolicitudes: &legajosApp.ListarSolicitudes{Legajos: legajoRepo},
		CasoAtender:           atender,
		CasoDespachar:         despachar,
		CasoEntregar:          &legajosApp.Entregar{Legajos: legajoRepo, Atender: atender, Despachar: despachar},
		CasoDevolver:          &legajosApp.Devolver{Legajos: legajoRepo, Auditoria: auditoria, Tx: tx, Reloj: reloj},
		CasoListarMovimientos: &legajosApp.ListarMovimientos{Legajos: legajoRepo},
	}

	// --- CU-25 y CU-26: tablero de indicadores y reportes ---
	indicadoresHandler := &handlers.IndicadoresHandler{
		CasoTablero: &indicadoresApp.ConsultarTablero{Indicadores: indicadorRepo, Reloj: reloj, DiasMaximosPrestamo: cfg.DiasMaximosPrestamo},
		CasoReporte: &indicadoresApp.GenerarReporte{
			Indicadores: indicadorRepo, Usuarios: usuarioRepo, Exportador: reportes.NewExportador(), Auditoria: auditoria,
			Reloj: reloj, DiasMaximosPrestamo: cfg.DiasMaximosPrestamo,
		},
	}
	// Reportes en lenguaje natural: sin GROQ_API_KEY el formulario sigue igual y la ruta no existe.
	if cfg.GroqKey != "" {
		indicadoresHandler.CasoInterpretar = &indicadoresApp.InterpretarReporte{Interprete: llm.NewGroq(cfg.GroqKey, cfg.GroqModelo), Reloj: reloj}
	} else {
		slog.Warn("GROQ_API_KEY no configurado: los reportes en lenguaje natural estarán deshabilitados")
	}

	// --- Límite de tasa en /api/auth/* (10 req/min por IP, sección 6) ---
	limitStorage := redisStorage.New(redisStorage.Config{URL: cfg.RedisURL})
	limitadorAuth := limiter.New(limiter.Config{
		Max: 10, Expiration: time.Minute, Storage: limitStorage,
		KeyGenerator: func(c *fiber.Ctx) string { return c.IP() },
		LimitReached: func(c *fiber.Ctx) error {
			return apierror.Nuevo(429, "LIMITE_TASA", "Demasiadas solicitudes, intenta más tarde")
		},
	})

	// Portal público (CU-12, CU-13): 30 req/min por IP, con su propio contador en Redis.
	limitadorPublico := limiter.New(limiter.Config{
		Max: 30, Expiration: time.Minute, Storage: limitStorage,
		KeyGenerator: func(c *fiber.Ctx) string { return "publico:" + c.IP() },
		LimitReached: func(c *fiber.Ctx) error {
			return apierror.Nuevo(429, "LIMITE_TASA", "Demasiadas solicitudes, intenta más tarde")
		},
	})

	// Radicar desde el portal sube archivos y crea trámites: 5 por hora por IP.
	limitadorSolicitudes := limiter.New(limiter.Config{
		Max: 5, Expiration: time.Hour, Storage: limitStorage,
		KeyGenerator: func(c *fiber.Ctx) string { return "solicitud:" + c.IP() },
		LimitReached: func(c *fiber.Ctx) error {
			return apierror.Nuevo(429, "LIMITE_TASA", "Demasiadas solicitudes, intenta más tarde")
		},
	})

	app := fiber.New(fiber.Config{
		BodyLimit:               40 << 20, // solicitudes del portal con hasta varios adjuntos de 8 MB
		ErrorHandler:            middleware.ManejadorErrores,
		ProxyHeader:             fiber.HeaderXForwardedFor,
		EnableTrustedProxyCheck: true,
		TrustedProxies:          cfg.ProxiesConfiables,
	})
	app.Use(requestid.New())
	app.Use(helmet.New())   // nosniff, X-Frame-Options, etc.
	app.Use(compress.New()) // el texto OCR de un documento pesa varios MB
	app.Use(fiberlog.New(fiberlog.Config{Format: "${time} ${status} ${method} ${path} reqid=${locals:requestid}\n"}))

	if cfg.AppEnv == "development" {
		httpapi.RegistrarDocs(app)
	}

	httpapi.Registrar(app, httpapi.Dependencias{
		Auth: authHandler, Perfil: perfilHandler,
		Usuarios: usuariosHandler, Roles: rolesHandler, Auditoria: auditoriaHandler,
		Docentes: docentesHandler, Categorias: categoriasHandler, Tramites: tramitesHandler, Expediente: expedienteHandler, Documentos: documentosHandler, Legajos: legajosHandler, Notificaciones: notificacionesHandler,
		Indicadores: indicadoresHandler, LimitadorPublico: limitadorPublico, LimitadorSolicitudes: limitadorSolicitudes,
		Tokens: jwtServicio, Cache: sessionCache, UsuarioRepo: usuarioRepo, LimitadorAuth: limitadorAuth,
	})

	log.Fatal(app.Listen(":" + cfg.Puerto))
}

type relojReal struct{}

func (relojReal) Ahora() time.Time { return time.Now() }
