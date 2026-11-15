package persistence

import (
	"context"
	"math"
	"testing"
	"time"

	"github.com/google/uuid"

	"escalafon/backend/internal/domain/indicador"
	"escalafon/backend/internal/domain/legajo"
	"escalafon/backend/internal/domain/tramite"
)

// Comprueba contra Postgres las consultas del tablero y los reportes (CU-25, CU-26). Los trámites
// se fechan en 2001 para que el periodo los aísle de cualquier otro dato de la base de pruebas.
func TestRepositorioIndicador_Consultas(t *testing.T) {
	dbase := conectarPrueba(t)
	conTx(t, dbase, func(ctx context.Context) {
		crear := func(v any) {
			t.Helper()
			if err := db(ctx, dbase).Create(v).Error; err != nil {
				t.Fatal(err)
			}
		}
		permiso := PermisoModel{Codigo: "TRAMITE_REVISAR", Nombre: "Revisar trámites", Modulo: "Trámites", CasoUso: "HU-06"}
		if err := db(ctx, dbase).Where("codigo = ?", permiso.Codigo).FirstOrCreate(&permiso).Error; err != nil {
			t.Fatal(err)
		}
		rolm := RolModel{Nombre: "rol-indicadores-prueba"}
		crear(&rolm)
		if err := db(ctx, dbase).Exec("INSERT INTO rol_permiso (rol_id, permiso_id) VALUES (?, ?)", rolm.ID, permiso.ID).Error; err != nil {
			t.Fatal(err)
		}
		analista := UsuarioModel{RolID: rolm.ID, Nombres: "Ana", Apellidos: "Indicadores", Correo: "indicadores-prueba@ded.dev", HashContrasena: "x"}
		crear(&analista)
		codigo := "77701"
		docente := DocenteModel{Codigo: &codigo, Nombres: "Doc", Apellidos: "Indicadores"}
		crear(&docente)
		cat := CategoriaModel{Nombre: "Categoría de prueba CU-25"}
		crear(&cat)

		f := func(m time.Month, d int) time.Time { return time.Date(2001, m, d, 12, 0, 0, 0, time.UTC) }
		remitido := TramiteModel{
			ID: uuid.NewString(), Codigo: "ED-2001-INDREM", DocenteID: docente.ID, CategoriaID: &cat.ID, Tipo: string(tramite.Categorizacion),
			Estado: string(tramite.Remitido), Origen: string(tramite.OrigenStaff), AnalistaID: &analista.ID, CreadoEn: f(3, 1),
		}
		enRevision := TramiteModel{
			ID: uuid.NewString(), Codigo: "ED-2001-INDREV", DocenteID: docente.ID, CategoriaID: &cat.ID, Tipo: string(tramite.Regularizacion),
			Estado: string(tramite.EnRevision), Origen: string(tramite.OrigenPortal), AnalistaID: &analista.ID, CreadoEn: f(3, 5),
		}
		crear(&remitido)
		crear(&enRevision)
		crear(&TramiteHistorialModel{TramiteID: remitido.ID, Estado: string(tramite.Ingresado), Fecha: f(3, 1)})
		crear(&TramiteHistorialModel{TramiteID: remitido.ID, Estado: string(tramite.Remitido), Fecha: f(3, 21)})
		crear(&TramiteHistorialModel{TramiteID: remitido.ID, Estado: string(tramite.Archivado), Fecha: f(4, 30)}) // cuenta el primer paso

		leg := LegajoModel{Codigo: "LEG-IND-01", Estado: string(legajo.Prestado), DocenteID: docente.ID}
		crear(&leg)
		sol := SolicitudCarpetaModel{LegajoID: leg.ID, SolicitanteID: analista.ID, TramiteID: &enRevision.ID, Motivo: "revisión", Estado: string(legajo.Despachada)}
		crear(&sol)
		limite := time.Date(2001, 3, 10, 0, 0, 0, 0, time.UTC)
		crear(&MovimientoCustodiaModel{
			LegajoID: leg.ID, SolicitudID: sol.ID, Tipo: string(legajo.Salida), CustodioID: analista.ID, FechaLimite: &limite,
			EstadoMaterial: string(legajo.Bueno), RegistradoPor: analista.ID, Fecha: f(3, 6),
		})

		repo := NewRepositorioIndicador(dbase)
		desde, hasta := time.Date(2001, 1, 1, 0, 0, 0, 0, time.UTC), time.Date(2002, 1, 1, 0, 0, 0, 0, time.UTC)
		p := indicador.Periodo{Desde: &desde, Hasta: &hasta}

		conteos, err := repo.ContarPorEstado(ctx, p)
		if err != nil || len(conteos) != 2 {
			t.Fatalf("esperaba dos estados en 2001: %v %+v", err, conteos)
		}

		resueltos, err := repo.ListarResueltos(ctx, indicador.FiltroReporte{Periodo: p})
		if err != nil || len(resueltos) != 1 {
			t.Fatalf("esperaba un resuelto: %v %+v", err, resueltos)
		}
		if r := resueltos[0]; r.Codigo != remitido.Codigo || math.Abs(r.Dias()-20) > 1e-9 || r.Analista != "Ana Indicadores" || r.Categoria != cat.Nombre {
			t.Fatalf("resuelto inesperado (se resuelve al pasar a Remitido, 20 días): %+v", r)
		}
		if otros, _ := repo.ListarResueltos(ctx, indicador.FiltroReporte{Periodo: p, Tipo: tramite.Regularizacion}); len(otros) != 0 {
			t.Fatalf("el filtro por tipo excluye la categorización: %+v", otros)
		}

		analistas, err := repo.CargaAnalistas(ctx, p)
		if err != nil {
			t.Fatal(err)
		}
		var carga *indicador.CargaAnalista
		for i := range analistas {
			if analistas[i].ID == analista.ID {
				carga = &analistas[i]
			}
		}
		if carga == nil || carga.Activos != 1 || carga.EnFacultad != 0 || carga.Resueltos != 1 {
			t.Fatalf("carga inesperada: %+v", carga)
		}

		prestamos, err := repo.ListarPrestamos(ctx)
		if err != nil {
			t.Fatal(err)
		}
		var prestamo *indicador.Prestamo
		for i := range prestamos {
			if prestamos[i].LegajoCodigo == leg.Codigo {
				prestamo = &prestamos[i]
			}
		}
		if prestamo == nil || prestamo.TramiteCodigo != enRevision.Codigo || prestamo.Custodio != "Ana Indicadores" ||
			prestamo.FechaLimite == nil || prestamo.FechaLimite.Format("2006-01-02") != "2001-03-10" {
			t.Fatalf("préstamo inesperado: %+v", prestamo)
		}

		tramites, err := repo.ListarTramites(ctx, indicador.FiltroReporte{Periodo: p, Estado: tramite.EnRevision})
		if err != nil || len(tramites) != 1 || tramites[0].Codigo != enRevision.Codigo || tramites[0].Analista != "Ana Indicadores" {
			t.Fatalf("reporte de trámites inesperado: %v %+v", err, tramites)
		}

		inventario, err := repo.ListarInventario(ctx, indicador.FiltroReporte{EstadoLegajo: legajo.Prestado})
		if err != nil {
			t.Fatal(err)
		}
		var fila *indicador.FilaInventario
		for i := range inventario {
			if inventario[i].LegajoCodigo == leg.Codigo {
				fila = &inventario[i]
			}
			if inventario[i].Estado != legajo.Prestado {
				t.Fatalf("el filtro por estado deja solo los prestados: %+v", inventario[i])
			}
		}
		if fila == nil || fila.Custodio != "Ana Indicadores" || fila.Salida == nil || fila.DocenteCodigo != codigo {
			t.Fatalf("legajo del inventario inesperado: %+v", fila)
		}
	})
}
