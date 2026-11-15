package persistence

import (
	"context"
	"errors"
	"os"
	"testing"
	"time"

	"gorm.io/gorm"

	"escalafon/backend/internal/domain/auditoria"
	"escalafon/backend/internal/domain/rol"
	"escalafon/backend/internal/domain/token"
	"escalafon/backend/internal/domain/usuario"
)

// conectarPrueba abre una conexión al Postgres de docker compose (nunca a Neon, sección 4.1).
func conectarPrueba(t *testing.T) *gorm.DB {
	t.Helper()
	dsn := os.Getenv("DATABASE_URL_TEST")
	if dsn == "" {
		t.Skip("DATABASE_URL_TEST no está definida; omitiendo pruebas de integración")
	}
	dbase, err := Conectar(dsn)
	if err != nil {
		t.Fatalf("no se pudo conectar a la base de datos de pruebas: %v", err)
	}
	return dbase
}

var errRollbackPrueba = errors.New("rollback de prueba")

// conTx ejecuta fn dentro de una transacción que siempre se revierte, para que cada prueba
// parta de una base limpia sin necesidad de truncar tablas (la auditoría es de solo anexión).
func conTx(t *testing.T, dbase *gorm.DB, fn func(ctx context.Context)) {
	t.Helper()
	txm := NewTxManager(dbase)
	err := txm.WithinTx(context.Background(), func(ctx context.Context) error {
		fn(ctx)
		return errRollbackPrueba
	})
	if !errors.Is(err, errRollbackPrueba) {
		t.Fatalf("la transacción de prueba no se revirtió como se esperaba: %v", err)
	}
}

func crearPermisoPrueba(t *testing.T, ctx context.Context, dbase *gorm.DB, codigo string) int16 {
	t.Helper()
	m := PermisoModel{Codigo: codigo, Nombre: codigo, Modulo: "Prueba", CasoUso: "HU-99"}
	if err := db(ctx, dbase).Create(&m).Error; err != nil {
		t.Fatalf("no se pudo crear el permiso de apoyo %s: %v", codigo, err)
	}
	return m.ID
}

func TestRepositorioUsuario_CrearYBuscar(t *testing.T) {
	dbase := conectarPrueba(t)
	conTx(t, dbase, func(ctx context.Context) {
		roles := NewRepositorioRol(dbase)
		r := &rol.Rol{Nombre: "Director-prueba"}
		if err := roles.Guardar(ctx, r, nil); err != nil {
			t.Fatalf("no se pudo crear el rol de apoyo: %v", err)
		}

		repo := NewRepositorioUsuario(dbase)
		u := &usuario.Usuario{
			RolID: r.ID, Nombres: "Ana", Apellidos: "Pérez", Correo: "Ana.Perez@uagrm.edu.bo",
			HashContrasena: "hash", VersionSesion: 1, Activo: true,
		}
		if err := repo.Crear(ctx, u); err != nil {
			t.Fatalf("no se pudo crear el usuario: %v", err)
		}
		if u.ID == "" {
			t.Fatalf("esperaba que Crear asignara un ID")
		}

		// Correo único sin distinguir mayúsculas.
		encontrado, err := repo.BuscarPorCorreo(ctx, "ana.perez@uagrm.edu.bo")
		if err != nil {
			t.Fatalf("esperaba encontrar el usuario por correo en minúsculas: %v", err)
		}
		if encontrado.ID != u.ID {
			t.Fatalf("se encontró un usuario distinto")
		}

		if _, err := repo.BuscarPorID(ctx, "00000000-0000-0000-0000-000000000000"); !errors.Is(err, usuario.ErrNoEncontrado) {
			t.Fatalf("esperaba ErrNoEncontrado, obtuve %v", err)
		}
	})
}

func TestRepositorioUsuario_RechazaCorreoDuplicado(t *testing.T) {
	dbase := conectarPrueba(t)
	conTx(t, dbase, func(ctx context.Context) {
		roles := NewRepositorioRol(dbase)
		r := &rol.Rol{Nombre: "Director-prueba-2"}
		roles.Guardar(ctx, r, nil)

		repo := NewRepositorioUsuario(dbase)
		original := &usuario.Usuario{RolID: r.ID, Nombres: "Ana", Apellidos: "Pérez", Correo: "dup@uagrm.edu.bo", HashContrasena: "hash", VersionSesion: 1, Activo: true}
		if err := repo.Crear(ctx, original); err != nil {
			t.Fatalf("no se pudo crear el usuario original: %v", err)
		}

		duplicado := &usuario.Usuario{RolID: r.ID, Nombres: "Otra", Apellidos: "Persona", Correo: "DUP@UAGRM.EDU.BO", HashContrasena: "hash", VersionSesion: 1, Activo: true}
		if err := repo.Crear(ctx, duplicado); !errors.Is(err, usuario.ErrCorreoDuplicado) {
			t.Fatalf("esperaba ErrCorreoDuplicado, obtuve %v", err)
		}
	})
}

func TestRepositorioUsuario_Actualizar(t *testing.T) {
	dbase := conectarPrueba(t)
	conTx(t, dbase, func(ctx context.Context) {
		roles := NewRepositorioRol(dbase)
		r := &rol.Rol{Nombre: "Director-prueba-3"}
		roles.Guardar(ctx, r, nil)

		repo := NewRepositorioUsuario(dbase)
		u := &usuario.Usuario{RolID: r.ID, Nombres: "Ana", Apellidos: "Pérez", Correo: "act@uagrm.edu.bo", HashContrasena: "hash", VersionSesion: 1, Activo: true}
		repo.Crear(ctx, u)

		u.Telefono = "70000000"
		if err := repo.Actualizar(ctx, u); err != nil {
			t.Fatalf("no se pudo actualizar: %v", err)
		}
		releido, err := repo.BuscarPorID(ctx, u.ID)
		if err != nil || releido.Telefono != "70000000" {
			t.Fatalf("la actualización no se reflejó: %+v, err=%v", releido, err)
		}
	})
}

func TestRepositorioRol_GuardarYPermisos(t *testing.T) {
	dbase := conectarPrueba(t)
	conTx(t, dbase, func(ctx context.Context) {
		roles := NewRepositorioRol(dbase)
		p1 := crearPermisoPrueba(t, ctx, dbase, "PRUEBA_UNO")
		p2 := crearPermisoPrueba(t, ctx, dbase, "PRUEBA_DOS")

		r := &rol.Rol{Nombre: "Rol de prueba", Descripcion: "temporal"}
		if err := roles.Guardar(ctx, r, []int16{p1, p2}); err != nil {
			t.Fatalf("no se pudo guardar el rol: %v", err)
		}

		codigos, err := roles.PermisosPorRol(ctx, r.ID)
		if err != nil {
			t.Fatalf("no se pudo listar permisos por rol: %v", err)
		}
		if len(codigos) != 2 {
			t.Fatalf("esperaba 2 permisos, obtuve %v", codigos)
		}

		conPermisos, err := roles.BuscarPorID(ctx, r.ID)
		if err != nil || len(conPermisos.Permisos) != 2 {
			t.Fatalf("BuscarPorID debería traer los permisos: %+v, err=%v", conPermisos, err)
		}

		// Guardar de nuevo con un solo permiso debe reemplazar la asignación completa.
		if err := roles.Guardar(ctx, r, []int16{p1}); err != nil {
			t.Fatalf("no se pudo reemplazar los permisos: %v", err)
		}
		codigos, _ = roles.PermisosPorRol(ctx, r.ID)
		if len(codigos) != 1 {
			t.Fatalf("esperaba que la reasignación dejara 1 permiso, obtuve %v", codigos)
		}
	})
}

func TestRepositorioToken_CrearBuscarMarcarUsado(t *testing.T) {
	dbase := conectarPrueba(t)
	conTx(t, dbase, func(ctx context.Context) {
		roles := NewRepositorioRol(dbase)
		r := &rol.Rol{Nombre: "Rol para token"}
		roles.Guardar(ctx, r, nil)
		usuarios := NewRepositorioUsuario(dbase)
		u := &usuario.Usuario{RolID: r.ID, Nombres: "A", Apellidos: "B", Correo: "tok@uagrm.edu.bo", HashContrasena: "h", VersionSesion: 1, Activo: true}
		usuarios.Crear(ctx, u)

		tokens := NewRepositorioToken(dbase)
		tk := &token.TokenRecuperacion{UsuarioID: u.ID, Proposito: token.PropositoActivacion, TokenHash: "abc123hash", ExpiraEn: time.Now().Add(time.Hour)}
		if err := tokens.Crear(ctx, tk); err != nil {
			t.Fatalf("no se pudo crear el token: %v", err)
		}

		encontrado, err := tokens.BuscarPorHash(ctx, "abc123hash")
		if err != nil {
			t.Fatalf("no se pudo encontrar el token: %v", err)
		}
		if !encontrado.EsValido(time.Now()) {
			t.Fatalf("el token debería ser válido")
		}

		if err := tokens.MarcarUsado(ctx, encontrado.ID, time.Now()); err != nil {
			t.Fatalf("no se pudo marcar como usado: %v", err)
		}
		releido, _ := tokens.BuscarPorHash(ctx, "abc123hash")
		if releido.EsValido(time.Now()) {
			t.Fatalf("el token ya no debería ser válido tras marcarse usado")
		}
	})
}

func TestRepositorioToken_HashInexistente(t *testing.T) {
	dbase := conectarPrueba(t)
	conTx(t, dbase, func(ctx context.Context) {
		tokens := NewRepositorioToken(dbase)
		if _, err := tokens.BuscarPorHash(ctx, "no-existe"); !errors.Is(err, token.ErrInvalido) {
			t.Fatalf("esperaba ErrInvalido, obtuve %v", err)
		}
	})
}

func TestRepositorioAuditoria_Registrar(t *testing.T) {
	dbase := conectarPrueba(t)
	conTx(t, dbase, func(ctx context.Context) {
		repo := NewRepositorioAuditoria(dbase)
		e := &auditoria.Evento{
			Accion: "LOGIN", Entidad: "usuario", EntidadID: "algún-id",
			Detalle: map[string]any{"ip": "127.0.0.1"}, Fecha: time.Now(),
		}
		if err := repo.Registrar(ctx, e); err != nil {
			t.Fatalf("no se pudo registrar el evento de auditoría: %v", err)
		}
		if e.ID == 0 {
			t.Fatalf("esperaba que Registrar asignara un ID")
		}
	})
}

// sembrarAuditoria crea dos usuarios y cuatro eventos en días consecutivos (uno sin usuario).
func sembrarAuditoria(t *testing.T, ctx context.Context, dbase *gorm.DB) (ana, luis string, base time.Time) {
	t.Helper()
	r := &rol.Rol{Nombre: "Rol para auditoría"}
	if err := NewRepositorioRol(dbase).Guardar(ctx, r, nil); err != nil {
		t.Fatalf("no se pudo crear el rol de apoyo: %v", err)
	}
	usuarios := NewRepositorioUsuario(dbase)
	ua := &usuario.Usuario{RolID: r.ID, Nombres: "Ana", Apellidos: "Pérez", Correo: "ana.aud@uagrm.edu.bo", HashContrasena: "h", VersionSesion: 1, Activo: true}
	ul := &usuario.Usuario{RolID: r.ID, Nombres: "Luis", Apellidos: "Rojas", Correo: "luis.aud@uagrm.edu.bo", HashContrasena: "h", VersionSesion: 1, Activo: true}
	for _, u := range []*usuario.Usuario{ua, ul} {
		if err := usuarios.Crear(ctx, u); err != nil {
			t.Fatalf("no se pudo crear el usuario de apoyo: %v", err)
		}
	}

	// Fechas futuras para no mezclarse con otros eventos que existan en la base de pruebas.
	base = time.Date(2099, 1, 10, 12, 0, 0, 0, time.UTC)
	repo := NewRepositorioAuditoria(dbase)
	for _, e := range []auditoria.Evento{
		{UsuarioID: &ua.ID, Accion: auditoria.AccionLogin, Entidad: auditoria.EntidadUsuario, EntidadID: ua.ID, DireccionIP: "10.0.0.1", Fecha: base},
		{UsuarioID: &ua.ID, Accion: auditoria.AccionEditar, Entidad: auditoria.EntidadUsuario, EntidadID: ul.ID,
			Detalle: map[string]any{"anterior": map[string]any{"telefono": "1"}, "nuevo": map[string]any{"telefono": "2"}}, Fecha: base.AddDate(0, 0, 1)},
		{UsuarioID: &ul.ID, Accion: auditoria.AccionLogin, Entidad: auditoria.EntidadUsuario, EntidadID: ul.ID, Fecha: base.AddDate(0, 0, 2)},
		{Accion: auditoria.AccionLoginFallido, Entidad: auditoria.EntidadUsuario, Fecha: base.AddDate(0, 0, 3)},
	} {
		if err := repo.Registrar(ctx, &e); err != nil {
			t.Fatalf("no se pudo registrar el evento: %v", err)
		}
	}
	return ua.ID, ul.ID, base
}

func TestRepositorioAuditoria_ListarFiltraYOrdena(t *testing.T) {
	dbase := conectarPrueba(t)
	conTx(t, dbase, func(ctx context.Context) {
		ana, luis, base := sembrarAuditoria(t, ctx, dbase)
		repo := NewRepositorioAuditoria(dbase)
		desde := base
		fin := base.AddDate(0, 0, 4)
		todo := auditoria.Filtro{Desde: &desde, Hasta: &fin, Pagina: 1, Tamano: 20}

		lista, total, err := repo.Listar(ctx, todo)
		if err != nil {
			t.Fatalf("no se pudo listar: %v", err)
		}
		if total != 4 || len(lista) != 4 {
			t.Fatalf("esperaba 4 eventos, obtuve total=%d len=%d", total, len(lista))
		}
		if lista[0].Accion != auditoria.AccionLoginFallido || lista[0].Autor != nil || lista[3].Accion != auditoria.AccionLogin {
			t.Fatalf("el orden debería ser del más reciente al más antiguo: %+v", lista)
		}
		if lista[3].Autor == nil || lista[3].Autor.Nombres != "Ana" || lista[3].DireccionIP != "10.0.0.1" {
			t.Fatalf("el autor o la IP no se cargaron: %+v", lista[3])
		}
		if nuevo, ok := lista[2].Detalle["nuevo"].(map[string]any); !ok || nuevo["telefono"] != "2" {
			t.Fatalf("el detalle JSON no se leyó: %+v", lista[2].Detalle)
		}

		casos := []struct {
			nombre string
			ajuste func(f *auditoria.Filtro)
			total  int64
		}{
			{"usuario", func(f *auditoria.Filtro) { f.UsuarioID = ana }, 2},
			{"acción", func(f *auditoria.Filtro) { f.Accion = auditoria.AccionLogin }, 2},
			{"rango", func(f *auditoria.Filtro) {
				d, h := base.AddDate(0, 0, 1), base.AddDate(0, 0, 3)
				f.Desde, f.Hasta = &d, &h
			}, 2},
			{"usuario y acción", func(f *auditoria.Filtro) { f.UsuarioID, f.Accion = luis, auditoria.AccionLogin }, 1},
			{"sin coincidencias", func(f *auditoria.Filtro) { f.UsuarioID, f.Accion = luis, auditoria.AccionEditar }, 0},
		}
		for _, c := range casos {
			f := todo
			c.ajuste(&f)
			if _, total, err := repo.Listar(ctx, f); err != nil || total != c.total {
				t.Fatalf("%s: esperaba %d, obtuve %d (err=%v)", c.nombre, c.total, total, err)
			}
		}

		pagina := todo
		pagina.Pagina, pagina.Tamano = 2, 3
		if lista, total, _ := repo.Listar(ctx, pagina); total != 4 || len(lista) != 1 {
			t.Fatalf("la segunda página de 3 debería tener 1 evento: total=%d len=%d", total, len(lista))
		}

		autores, err := repo.ListarAutores(ctx)
		if err != nil {
			t.Fatalf("no se pudo listar los autores: %v", err)
		}
		encontrados := 0
		for _, a := range autores {
			if a.ID == ana || a.ID == luis {
				encontrados++
			}
		}
		if encontrados != 2 {
			t.Fatalf("esperaba a Ana y Luis entre los autores: %+v", autores)
		}
	})
}

// La tabla auditoria es de solo anexión: el trigger rechaza UPDATE y DELETE aun con SQL directo.
func TestRepositorioAuditoria_RechazaModificarYEliminar(t *testing.T) {
	dbase := conectarPrueba(t)
	conTx(t, dbase, func(ctx context.Context) {
		repo := NewRepositorioAuditoria(dbase)
		e := &auditoria.Evento{Accion: auditoria.AccionLogin, Entidad: auditoria.EntidadUsuario, Fecha: time.Now()}
		if err := repo.Registrar(ctx, e); err != nil {
			t.Fatalf("no se pudo registrar el evento: %v", err)
		}
		tx := db(ctx, dbase)
		for _, sentencia := range []string{
			"UPDATE auditoria SET accion = 'OTRA' WHERE id = ?",
			"DELETE FROM auditoria WHERE id = ?",
		} {
			tx.SavePoint("antes")
			if err := tx.Exec(sentencia, e.ID).Error; err == nil {
				t.Fatalf("esperaba que el trigger rechazara: %s", sentencia)
			}
			tx.RollbackTo("antes")
		}
		var accion string
		if err := tx.Raw("SELECT accion FROM auditoria WHERE id = ?", e.ID).Scan(&accion).Error; err != nil || accion != auditoria.AccionLogin {
			t.Fatalf("el evento debería seguir intacto: accion=%q err=%v", accion, err)
		}
	})
}
