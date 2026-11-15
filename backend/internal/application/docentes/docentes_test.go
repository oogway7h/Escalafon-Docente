package docentes_test

import (
	"context"
	"testing"
	"time"

	"escalafon/backend/internal/application/docentes"
	"escalafon/backend/internal/application/testdobles"
	"escalafon/backend/internal/domain/docente"
)

func TestGuardarConservaCedulaAsignada(t *testing.T) {
	repo := testdobles.NewDocenteRepoFalso()
	repo.Sembrar(docente.Docente{ID: "d1", Codigo: "100", Cedula: "555", Nombres: "Ana", Apellidos: "Paz", Activo: true})
	caso := &docentes.Guardar{
		Docentes: repo, Categorias: testdobles.NewCategoriaRepoFalso(), Auditoria: &testdobles.AuditoriaRepoFalso{},
		Reloj: &testdobles.RelojFijo{Momento: time.Now()},
	}
	if _, err := caso.Ejecutar(context.Background(), docentes.Entrada{ActorID: "u", ID: "d1", Cedula: "999", Nombres: "Ana María", Apellidos: "Paz"}); err != nil {
		t.Fatal(err)
	}
	got, _ := repo.BuscarPorID(context.Background(), "d1")
	if got.Cedula != "555" {
		t.Errorf("la cédula asignada no debe cambiar, quedó %q", got.Cedula)
	}
}

func TestListarPaginaAcotaParametros(t *testing.T) {
	repo := testdobles.NewDocenteRepoFalso()
	repo.Sembrar(docente.Docente{ID: "d1", Nombres: "Ana", Apellidos: "Paz", Activo: true})
	_, total, err := (&docentes.ListarPagina{Docentes: repo}).Ejecutar(context.Background(), "", false, -3, 5000)
	if err != nil || total != 1 {
		t.Fatalf("página y tamaño inválidos deben acotarse: %v %d", err, total)
	}
}
