package persistence

import (
	"context"
	"errors"
	"testing"
	"time"

	"escalafon/backend/internal/domain/legajo"
)

// Recorre el ciclo real contra Postgres: código automático, búsqueda exacta, solicitud única,
// salida/entrada y la inmutabilidad de la bitácora.
func TestRepositorioLegajo_Ciclo(t *testing.T) {
	dbase := conectarPrueba(t)
	conTx(t, dbase, func(ctx context.Context) {
		rol := RolModel{Nombre: "rol-legajo-prueba"}
		usr := UsuarioModel{Nombres: "Ana", Apellidos: "Archivo", Correo: "legajo-prueba@ded.dev", HashContrasena: "x"}
		if err := db(ctx, dbase).Create(&rol).Error; err != nil {
			t.Fatal(err)
		}
		usr.RolID = rol.ID
		if err := db(ctx, dbase).Create(&usr).Error; err != nil {
			t.Fatal(err)
		}
		for _, c := range []string{"98765", "987650"} {
			c := c
			if err := db(ctx, dbase).Create(&DocenteModel{Codigo: &c, Nombres: "Doc" + c, Apellidos: "Prueba"}).Error; err != nil {
				t.Fatal(err)
			}
		}
		docs := NewRepositorioDocente(dbase)
		d, err := docs.BuscarPorIdentificador(ctx, "98765")
		if err != nil {
			t.Fatal(err)
		}

		// La búsqueda por código no trae códigos parecidos.
		lista, err := docs.Listar(ctx, "98765", false)
		if err != nil || len(lista) != 1 || lista[0].Codigo != "98765" {
			t.Fatalf("búsqueda exacta: %v %+v", err, lista)
		}

		repo := NewRepositorioLegajo(dbase)
		l := &legajo.Legajo{DocenteID: d.ID, Gaveta: "G1", Orden: 1}
		if err := repo.CrearLegajo(ctx, l); err != nil {
			t.Fatal(err)
		}
		got, err := repo.BuscarLegajo(ctx, l.ID)
		if err != nil || got.Codigo == "" || got.Estado != legajo.EnArchivo {
			t.Fatalf("el trigger debe asignar código: %v %+v", err, got)
		}

		s := &legajo.Solicitud{LegajoID: l.ID, SolicitanteID: usr.ID, Motivo: "revisión", Estado: legajo.Pendiente}
		if err := repo.CrearSolicitud(ctx, s); err != nil {
			t.Fatal(err)
		}
		// Un savepoint evita que el error de unicidad aborte la transacción de la prueba.
		sp := db(ctx, dbase).SavePoint("dup")
		if err := repo.CrearSolicitud(ctx, &legajo.Solicitud{LegajoID: l.ID, SolicitanteID: usr.ID, Motivo: "x", Estado: legajo.Pendiente}); !errors.Is(err, legajo.ErrSolicitudAbierta) {
			t.Fatalf("segunda solicitud abierta: %v", err)
		}
		sp.RollbackTo("dup")

		ahora := time.Now()
		if err := repo.ActualizarSolicitud(ctx, s.ID, legajo.Pendiente, legajo.Atendida, "", usr.ID, ahora); err != nil {
			t.Fatal(err)
		}
		if err := repo.ActualizarSolicitud(ctx, s.ID, legajo.Pendiente, legajo.Atendida, "", usr.ID, ahora); !errors.Is(err, legajo.ErrEstadoSolicitud) {
			t.Fatalf("transición repetida: %v", err)
		}
		if ok, err := repo.CambiarEstadoLegajo(ctx, l.ID, legajo.EnArchivo, legajo.Prestado); err != nil || !ok {
			t.Fatalf("prestar: %v %v", ok, err)
		}
		limite := ahora.AddDate(0, 0, -1)
		if err := repo.AnexarMovimiento(ctx, &legajo.Movimiento{
			LegajoID: l.ID, SolicitudID: s.ID, Tipo: legajo.Salida, CustodioID: usr.ID, FechaLimite: &limite,
			EstadoMaterial: legajo.Bueno, RegistradoPorID: usr.ID, Fecha: ahora,
		}); err != nil {
			t.Fatal(err)
		}
		if err := repo.ActualizarSolicitud(ctx, s.ID, legajo.Atendida, legajo.Despachada, "", usr.ID, ahora); err != nil {
			t.Fatal(err)
		}
		movs, total, err := repo.ListarMovimientos(ctx, legajo.FiltroMovimientos{Vencidos: true, Busqueda: "98765"}, 1, 10)
		if err != nil || total != 1 || !movs[0].Vencido {
			t.Fatalf("la salida con fecha pasada debe figurar vencida: %v %d %+v", err, total, movs)
		}
		if err := db(ctx, dbase).Exec("UPDATE movimiento_custodia SET observacion = 'x'").Error; err == nil {
			t.Fatal("la bitácora de custodia debe ser de solo anexión")
		}
	})
}
