// Comando seed: crea los datos iniciales del módulo de seguridad (sección 4.5). Es idempotente:
// puede ejecutarse varias veces sin duplicar roles, permisos ni el primer Director.
package main

import (
	"fmt"
	"log"
	"time"

	"gorm.io/gorm/clause"

	"escalafon/backend/internal/infrastructure/config"
	"escalafon/backend/internal/infrastructure/mail"
	"escalafon/backend/internal/infrastructure/persistence"
	"escalafon/backend/internal/infrastructure/security"
)

const (
	nombreDirector   = "Director"
	nombreAnalista   = "Analista"
	nombreArchivo    = "Personal de Archivo"
	nombreVentanilla = "Personal de Ventanilla"
)

type permisoSeed struct {
	Codigo, Nombre, Modulo, CasoUso string
}

var catalogoPermisos = []permisoSeed{
	{"PERFIL_GESTIONAR", "Gestionar mi perfil", "Seguridad", "HU-01"},
	{"USUARIO_GESTIONAR", "Gestionar usuarios", "Seguridad", "HU-02"},
	{"ROL_GESTIONAR", "Gestionar roles", "Seguridad", "HU-02"},
	{"AUDITORIA_CONSULTAR", "Consultar auditoría", "Seguridad", "HU-03"},
	{"DOCENTE_GESTIONAR", "Gestionar docentes", "Dominio", "HU-04"},
	{"CATEGORIA_GESTIONAR", "Gestionar categorías", "Dominio", "HU-04"},
	{"TRAMITE_RADICAR", "Radicar trámites", "Trámites", "HU-05"},
	{"TRAMITE_ASIGNAR", "Asignar trámites", "Trámites", "HU-05"},
	{"TRAMITE_REVISAR", "Revisar trámites", "Trámites", "HU-06"},
	{"TRAMITE_ARCHIVAR", "Archivar trámites", "Trámites", "HU-06"},
	{"DOCUMENTO_GESTIONAR", "Gestionar documentos", "Documental", "HU-07"},
	{"DOCUMENTO_RESTITUIR", "Restituir documentos", "Documental", "HU-07"},
	{"DOCUMENTO_BUSCAR", "Buscar documentos", "Documental", "HU-08"},
	{"DOCUMENTO_ESCANEAR", "Cargar escaneos al expediente digital", "Documental", "HU-07"},
	{"LEGAJO_GESTIONAR", "Gestionar legajos", "Custodia", "HU-09"},
	{"CARPETA_SOLICITAR", "Solicitar carpetas", "Custodia", "HU-09"},
	{"CUSTODIA_REGISTRAR", "Registrar custodia", "Custodia", "HU-10"},
	{"BITACORA_CONSULTAR", "Consultar bitácora", "Custodia", "HU-10"},
	{"NOTIFICACION_CONSULTAR", "Consultar notificaciones", "Trámites", "HU-14"},
	{"TABLERO_VER", "Ver tablero", "Indicadores", "HU-11"},
	{"REPORTE_GENERAR", "Generar reportes", "Indicadores", "HU-12"},
}

var permisosPorRol = map[string][]string{
	nombreDirector: nil, // se completa con TODOS los códigos del catálogo
	nombreAnalista: {
		"PERFIL_GESTIONAR", "DOCENTE_GESTIONAR", "TRAMITE_REVISAR",
		"DOCUMENTO_GESTIONAR", "DOCUMENTO_BUSCAR", "CARPETA_SOLICITAR",
	},
	// Ventanilla (secretaría): radica trámites, valida las solicitudes del portal y registra el regreso de la facultad.
	// Consulta el padrón por TRAMITE_RADICAR; no modifica docentes.
	nombreVentanilla: {"PERFIL_GESTIONAR", "TRAMITE_RADICAR"},
	nombreArchivo: {
		"PERFIL_GESTIONAR", "LEGAJO_GESTIONAR", "CUSTODIA_REGISTRAR", "BITACORA_CONSULTAR", "TRAMITE_ARCHIVAR", "REPORTE_GENERAR",
			"DOCUMENTO_ESCANEAR", "DOCUMENTO_BUSCAR",
	},
}

func main() {
	cfg, err := config.Cargar()
	if err != nil {
		log.Fatalf("configuración inválida: %v", err)
	}

	db, err := persistence.Conectar(cfg.DatabaseURL)
	if err != nil {
		log.Fatalf("no se pudo conectar a la base de datos: %v", err)
	}

	// 1. Roles de sistema.
	roles := map[string]*persistence.RolModel{}
	for _, nombre := range []string{nombreDirector, nombreAnalista, nombreArchivo, nombreVentanilla} {
		r := &persistence.RolModel{Nombre: nombre, EsSistema: true, Activo: true}
		if err := db.Where(persistence.RolModel{Nombre: nombre}).
			Attrs(persistence.RolModel{EsSistema: true, Activo: true}).
			FirstOrCreate(r).Error; err != nil {
			log.Fatalf("no se pudo crear el rol %s: %v", nombre, err)
		}
		roles[nombre] = r
		log.Printf("rol %-20s id=%d", nombre, r.ID)
	}

	// 2. Catálogo de permisos.
	permisoIDPorCodigo := map[string]int16{}
	for _, p := range catalogoPermisos {
		m := &persistence.PermisoModel{Codigo: p.Codigo, Nombre: p.Nombre, Modulo: p.Modulo, CasoUso: p.CasoUso}
		if err := db.Where(persistence.PermisoModel{Codigo: p.Codigo}).
			Attrs(persistence.PermisoModel{Nombre: p.Nombre, Modulo: p.Modulo, CasoUso: p.CasoUso}).
			FirstOrCreate(m).Error; err != nil {
			log.Fatalf("no se pudo crear el permiso %s: %v", p.Codigo, err)
		}
		permisoIDPorCodigo[p.Codigo] = m.ID
	}
	log.Printf("catálogo de permisos: %d entradas", len(catalogoPermisos))

	// 3. Asignación de permisos por rol (Director recibe todos los permisos del catálogo).
	todos := make([]string, 0, len(catalogoPermisos))
	for _, p := range catalogoPermisos {
		todos = append(todos, p.Codigo)
	}
	permisosPorRol[nombreDirector] = todos

	for nombreRol, codigos := range permisosPorRol {
		r := roles[nombreRol]
		for _, codigo := range codigos {
			pid, ok := permisoIDPorCodigo[codigo]
			if !ok {
				log.Fatalf("permiso desconocido en la asignación: %s", codigo)
			}
			rp := persistence.RolPermisoModel{RolID: r.ID, PermisoID: pid}
			if err := db.Clauses(clause.OnConflict{DoNothing: true}).Create(&rp).Error; err != nil {
				log.Fatalf("no se pudo asignar %s a %s: %v", codigo, nombreRol, err)
			}
		}
		log.Printf("rol %-20s -> %d permisos asignados", nombreRol, len(codigos))
	}

	// Radicar es de Ventanilla: se retira del Analista, que lo tenía antes de existir ese rol.
	if err := db.Exec(`DELETE FROM rol_permiso WHERE rol_id = ? AND permiso_id = ?`,
		roles[nombreAnalista].ID, permisoIDPorCodigo["TRAMITE_RADICAR"]).Error; err != nil {
		log.Fatalf("no se pudo retirar TRAMITE_RADICAR al Analista: %v", err)
	}

	// 3b. Categorías del escalafón con sus puntajes y documentos exigidos (idempotente). Se reconocen por
	// clave normalizada, así "Titular A" y "Profesor Titular A" son la misma. Los puntajes son los del art. 18
	// del Reglamento del Escalafón y los documentos son los de "Requisitos para categorización docente" del
	// Departamento: son los mismos para todas las categorías (solo cambia el puntaje). Cada ejecución deja
	// exactamente esa lista (quita los documentos que ya no se exigen). Los docentes reales se cargan con
	// `make import-docentes`.
	requisitosCategorizacion := []string{
		"Carta «Solicitud de categorización» al Jefe del Dpto. de Escalafón Docente, con la facultad",
		"Certificado de Trabajo original de RR.HH. – UAGRM, actualizado, con facultad y carrera",
		"Dos fotos 4x4 (cualquier fondo)",
		"Hoja de vida (sin respaldo)",
		"Fotocopia simple de la cédula de identidad",
		"Fotocopia legalizada de títulos (académico, provisión nacional, diplomado, maestría)",
		"Certificado de evaluación de desempeño docente (D.E.A.) de las gestiones a evaluar",
		"Méritos: cursos, libros, textos guía, publicaciones y proyectos visados por Decano o Director",
	}
	var existentes []persistence.CategoriaModel
	if err := db.Find(&existentes).Error; err != nil {
		log.Fatalf("no se pudieron leer las categorías: %v", err)
	}
	porClave := map[string]persistence.CategoriaModel{}
	for _, c := range existentes {
		porClave[persistence.ClaveCategoria(c.Nombre)] = c
	}
	p := func(v int16) *int16 { return &v }
	for _, c := range []struct {
		Nombre       string
		Desde, Hasta *int16 // nil = sin puntaje; Hasta nil con Desde = "o más"
	}{
		{"Profesor Adjunto", p(30), p(89)}, {"Profesor Asociado", p(90), p(149)}, {"Profesor Titular A", p(150), p(299)},
		{"Profesor Titular B", p(300), p(449)}, {"Profesor Titular C", p(450), nil},
		{"Profesor Emérito", nil, nil}, {"Profesor Interino", nil, nil},
	} {
		cat, ok := porClave[persistence.ClaveCategoria(c.Nombre)]
		if !ok {
			cat = persistence.CategoriaModel{Nombre: c.Nombre, Activa: true}
			if err := db.Create(&cat).Error; err != nil {
				log.Fatalf("no se pudo crear la categoría %s: %v", c.Nombre, err)
			}
		}
		// La antigüedad ya no es un dato de la categoría: cada docente tiene las gestiones que le corresponde evaluar.
		if err := db.Model(&persistence.CategoriaModel{}).Where("id = ?", cat.ID).
			Updates(map[string]any{"anios_antiguedad": 0, "puntaje_desde": c.Desde, "puntaje_hasta": c.Hasta}).Error; err != nil {
			log.Fatalf("no se pudo actualizar el puntaje de %s: %v", c.Nombre, err)
		}
		if err := db.Where("categoria_id = ? AND nombre NOT IN ?", cat.ID, requisitosCategorizacion).
			Delete(&persistence.CategoriaRequisitoModel{}).Error; err != nil {
			log.Fatalf("no se pudieron retirar los requisitos antiguos de %s: %v", c.Nombre, err)
		}
		for _, n := range requisitosCategorizacion {
			rq := persistence.CategoriaRequisitoModel{CategoriaID: cat.ID, Nombre: n}
			if err := db.Where(rq).FirstOrCreate(&rq).Error; err != nil {
				log.Fatalf("no se pudo crear el requisito %s: %v", n, err)
			}
		}
	}
	log.Println("categorías del escalafón cargadas")

	// 4. Usuarios de prueba (uno por rol) con una contraseña fija y pública (está en este archivo): SOLO en
	// desarrollo. En cualquier otro entorno no se crean, porque serían cuentas abiertas a cualquiera (el
	// Director con todos los permisos incluido). No pasan por activación: quedan activos de inmediato.
	if cfg.AppEnv != "development" {
		log.Printf("APP_ENV=%q: no se crean los usuarios de prueba", cfg.AppEnv)
	} else {
		usuariosPrueba := []struct{ Correo, Nombres, Apellidos, Rol string }{
			{"admin@ded.dev", "Admin", "Prueba", nombreDirector},
			{"analista@ded.dev", "Analista", "Prueba", nombreAnalista},
			{"archivo@ded.dev", "Archivo", "Prueba", nombreArchivo},
			{"ventanilla@ded.dev", "Ventanilla", "Prueba", nombreVentanilla},
		}
		hasherPrueba := security.NewBcryptHasher()
		hashPrueba, err := hasherPrueba.Hash("123123")
		if err != nil {
			log.Fatalf("no se pudo hashear la contraseña de prueba: %v", err)
		}
		for _, up := range usuariosPrueba {
			var existente persistence.UsuarioModel
			if err := db.Where("lower(correo) = lower(?)", up.Correo).First(&existente).Error; err == nil {
				log.Printf("usuario de prueba %s ya existe: no se recrea", up.Correo)
				continue
			}
			u := persistence.UsuarioModel{
				RolID: roles[up.Rol].ID, Nombres: up.Nombres, Apellidos: up.Apellidos,
				Correo: up.Correo, HashContrasena: hashPrueba, VersionSesion: 1, Activo: true,
			}
			if err := db.Create(&u).Error; err != nil {
				log.Fatalf("no se pudo crear el usuario de prueba %s: %v", up.Correo, err)
			}
			log.Printf("usuario de prueba creado: %-20s rol=%s", up.Correo, up.Rol)
		}
	}

	// 5. Primer Director, con contraseña aleatoria inutilizable y token de activación.
	if cfg.SeedDirectorCorreo == "" {
		log.Println("SEED_DIRECTOR_CORREO vacío: se omite la creación del primer Director")
		return
	}
	var existente persistence.UsuarioModel
	err = db.Where("lower(correo) = lower(?)", cfg.SeedDirectorCorreo).First(&existente).Error
	if err == nil {
		log.Printf("el primer Director (%s) ya existe: no se recrea", cfg.SeedDirectorCorreo)
		return
	}

	hasher := security.NewBcryptHasher()
	generador := security.NewGeneradorAleatorio()

	claveAleatoria, _, err := generador.Generar()
	if err != nil {
		log.Fatalf("no se pudo generar la contraseña aleatoria: %v", err)
	}
	hashClave, err := hasher.Hash(claveAleatoria)
	if err != nil {
		log.Fatalf("no se pudo hashear la contraseña: %v", err)
	}

	director := persistence.UsuarioModel{
		RolID:          roles[nombreDirector].ID,
		Nombres:        cfg.SeedDirectorNombres,
		Apellidos:      cfg.SeedDirectorApellidos,
		Correo:         cfg.SeedDirectorCorreo,
		HashContrasena: hashClave,
		VersionSesion:  1,
		Activo:         true,
	}
	if err := db.Create(&director).Error; err != nil {
		log.Fatalf("no se pudo crear el primer Director: %v", err)
	}

	tokenPlano, tokenHash, err := generador.Generar()
	if err != nil {
		log.Fatalf("no se pudo generar el token de activación: %v", err)
	}
	activacion := persistence.TokenRecuperacionModel{
		UsuarioID: director.ID,
		Proposito: "ACTIVACION",
		TokenHash: tokenHash,
		ExpiraEn:  time.Now().Add(time.Duration(cfg.HorasTokenActivacion) * time.Hour),
	}
	if err := db.Create(&activacion).Error; err != nil {
		log.Fatalf("no se pudo crear el token de activación: %v", err)
	}

	enlace := fmt.Sprintf("%s/activar?token=%s", cfg.URLFrontend, tokenPlano)
	fmt.Println("========================================")
	fmt.Println("Primer Director creado:", cfg.SeedDirectorCorreo)
	fmt.Println("Enlace de activación:", enlace)
	fmt.Println("========================================")

	mailer := mail.NewMailerSMTP(cfg.SMTPHost, cfg.SMTPPuerto, cfg.SMTPUsuario, cfg.SMTPClave, cfg.SMTPRemitente)
	if err := mailer.EnviarActivacion(cfg.SeedDirectorCorreo, cfg.SeedDirectorNombres, enlace); err != nil {
		log.Printf("advertencia: no se pudo enviar el correo de activación (¿Mailpit está corriendo?): %v", err)
	}
}
