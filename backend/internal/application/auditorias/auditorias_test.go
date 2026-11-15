package auditorias_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"escalafon/backend/internal/application/auditorias"
	"escalafon/backend/internal/application/testdobles"
	"escalafon/backend/internal/domain/auditoria"
)

var base = time.Date(2026, 10, 1, 9, 0, 0, 0, time.UTC)

func ptr[T any](v T) *T { return &v }

// sembrar registra tres acciones de dos usuarios en días distintos y una sin usuario.
func sembrar(t *testing.T) *testdobles.AuditoriaRepoFalso {
	t.Helper()
	repo := &testdobles.AuditoriaRepoFalso{Autores: map[string]auditoria.Autor{
		"u-ana":  {ID: "u-ana", Nombres: "Ana", Apellidos: "Pérez", Correo: "ana@uagrm.edu.bo"},
		"u-luis": {ID: "u-luis", Nombres: "Luis", Apellidos: "Rojas", Correo: "luis@uagrm.edu.bo"},
	}}
	eventos := []auditoria.Evento{
		{UsuarioID: ptr("u-ana"), Accion: auditoria.AccionLogin, Entidad: auditoria.EntidadUsuario, Fecha: base},
		{UsuarioID: ptr("u-ana"), Accion: auditoria.AccionCrear, Entidad: auditoria.EntidadDocente, Fecha: base.AddDate(0, 0, 1)},
		{UsuarioID: ptr("u-luis"), Accion: auditoria.AccionLogin, Entidad: auditoria.EntidadUsuario, Fecha: base.AddDate(0, 0, 2)},
		{Accion: auditoria.AccionLoginFallido, Entidad: auditoria.EntidadUsuario, Fecha: base.AddDate(0, 0, 3)},
	}
	for i := range eventos {
		if err := repo.Registrar(context.Background(), &eventos[i]); err != nil {
			t.Fatalf("no se pudo sembrar: %v", err)
		}
	}
	return repo
}

func TestConsultar_SinFiltrosOrdenaDelMasReciente(t *testing.T) {
	caso := &auditorias.Consultar{Auditoria: sembrar(t)}
	res, err := caso.Ejecutar(context.Background(), auditoria.Filtro{})
	if err != nil {
		t.Fatalf("no esperaba error: %v", err)
	}
	if res.Total != 4 || len(res.Registros) != 4 {
		t.Fatalf("esperaba 4 registros, obtuve total=%d len=%d", res.Total, len(res.Registros))
	}
	if res.Registros[0].Accion != auditoria.AccionLoginFallido || res.Registros[3].Accion != auditoria.AccionLogin {
		t.Fatalf("el orden debería ser del más reciente al más antiguo: %+v", res.Registros)
	}
	if res.Registros[0].Autor != nil || res.Registros[1].Autor == nil || res.Registros[1].Autor.Nombres != "Luis" {
		t.Fatalf("el autor no se resolvió como se esperaba")
	}
	if res.Pagina != 1 || res.Tamano != 20 {
		t.Fatalf("esperaba la paginación por defecto (1, 20), obtuve (%d, %d)", res.Pagina, res.Tamano)
	}
}

func TestConsultar_FiltrosPorSeparadoYCombinados(t *testing.T) {
	caso := &auditorias.Consultar{Auditoria: sembrar(t)}
	casos := []struct {
		nombre string
		filtro auditoria.Filtro
		total  int64
	}{
		{"por usuario", auditoria.Filtro{UsuarioID: "u-ana"}, 2},
		{"por acción", auditoria.Filtro{Accion: auditoria.AccionLogin}, 2},
		{"por rango", auditoria.Filtro{Desde: ptr(base.AddDate(0, 0, 1)), Hasta: ptr(base.AddDate(0, 0, 3))}, 2},
		{"usuario y acción", auditoria.Filtro{UsuarioID: "u-ana", Accion: auditoria.AccionLogin}, 1},
		{"acción y rango", auditoria.Filtro{Accion: auditoria.AccionLogin, Desde: ptr(base.AddDate(0, 0, 1))}, 1},
		{"los tres", auditoria.Filtro{UsuarioID: "u-luis", Accion: auditoria.AccionLogin, Desde: ptr(base), Hasta: ptr(base.AddDate(0, 0, 3))}, 1},
		{"sin coincidencias", auditoria.Filtro{UsuarioID: "u-luis", Accion: auditoria.AccionCrear}, 0},
	}
	for _, c := range casos {
		res, err := caso.Ejecutar(context.Background(), c.filtro)
		if err != nil {
			t.Fatalf("%s: no esperaba error: %v", c.nombre, err)
		}
		if res.Total != c.total {
			t.Fatalf("%s: esperaba %d registros, obtuve %d", c.nombre, c.total, res.Total)
		}
	}
}

func TestConsultar_Paginacion(t *testing.T) {
	caso := &auditorias.Consultar{Auditoria: sembrar(t)}
	res, err := caso.Ejecutar(context.Background(), auditoria.Filtro{Pagina: 2, Tamano: 3})
	if err != nil {
		t.Fatalf("no esperaba error: %v", err)
	}
	if res.Total != 4 || len(res.Registros) != 1 || res.Pagina != 2 || res.Tamano != 3 {
		t.Fatalf("página inesperada: total=%d len=%d pagina=%d tamano=%d", res.Total, len(res.Registros), res.Pagina, res.Tamano)
	}

	res, _ = caso.Ejecutar(context.Background(), auditoria.Filtro{Pagina: -1, Tamano: 500})
	if res.Pagina != 1 || res.Tamano != 20 {
		t.Fatalf("esperaba normalizar a (1, 20), obtuve (%d, %d)", res.Pagina, res.Tamano)
	}
}

func TestConsultar_RechazaAccionDesconocida(t *testing.T) {
	caso := &auditorias.Consultar{Auditoria: sembrar(t)}
	if _, err := caso.Ejecutar(context.Background(), auditoria.Filtro{Accion: "BORRAR_TODO"}); !errors.Is(err, auditoria.ErrAccionInvalida) {
		t.Fatalf("esperaba ErrAccionInvalida, obtuve %v", err)
	}
}

func TestConsultar_RechazaRangoInvertido(t *testing.T) {
	caso := &auditorias.Consultar{Auditoria: sembrar(t)}
	_, err := caso.Ejecutar(context.Background(), auditoria.Filtro{Desde: ptr(base.AddDate(0, 0, 2)), Hasta: ptr(base)})
	if !errors.Is(err, auditoria.ErrRangoFechas) {
		t.Fatalf("esperaba ErrRangoFechas, obtuve %v", err)
	}
}

func TestConsultar_NoModificaLaBitacora(t *testing.T) {
	repo := sembrar(t)
	antes := append([]auditoria.Evento(nil), repo.Eventos...)
	caso := &auditorias.Consultar{Auditoria: repo}
	if _, err := caso.Ejecutar(context.Background(), auditoria.Filtro{UsuarioID: "u-ana"}); err != nil {
		t.Fatalf("no esperaba error: %v", err)
	}
	if len(repo.Eventos) != len(antes) {
		t.Fatalf("la consulta no debe agregar ni quitar eventos")
	}
	for i := range antes {
		if repo.Eventos[i].ID != antes[i].ID || repo.Eventos[i].Accion != antes[i].Accion || !repo.Eventos[i].Fecha.Equal(antes[i].Fecha) {
			t.Fatalf("la consulta no debe modificar eventos")
		}
	}
}

func TestListarOpciones_AccionesYAutores(t *testing.T) {
	caso := &auditorias.ListarOpciones{Auditoria: sembrar(t)}
	op, err := caso.Ejecutar(context.Background())
	if err != nil {
		t.Fatalf("no esperaba error: %v", err)
	}
	if len(op.Acciones) != len(auditoria.Acciones) {
		t.Fatalf("esperaba el catálogo completo de acciones")
	}
	if len(op.Usuarios) != 2 || op.Usuarios[0].Nombres != "Ana" {
		t.Fatalf("esperaba los dos autores ordenados por apellido, obtuve %+v", op.Usuarios)
	}
}
