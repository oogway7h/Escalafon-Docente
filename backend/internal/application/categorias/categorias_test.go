package categorias_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"escalafon/backend/internal/application/categorias"
	"escalafon/backend/internal/application/testdobles"
	"escalafon/backend/internal/domain/categoria"
)

func nuevoCaso(c categoria.Categoria) (*categorias.Eliminar, *testdobles.CategoriaRepoFalso) {
	repo := testdobles.NewCategoriaRepoFalso()
	repo.Sembrar(c)
	return &categorias.Eliminar{Categorias: repo, Auditoria: &testdobles.AuditoriaRepoFalso{}, Reloj: &testdobles.RelojFijo{Momento: time.Now()}}, repo
}

func TestEliminar_BorraUnaCategoriaSinUso(t *testing.T) {
	caso, repo := nuevoCaso(categoria.Categoria{ID: 1, Nombre: "Nueva", Activa: true})
	if err := caso.Ejecutar(context.Background(), "u1", 1); err != nil {
		t.Fatalf("eliminar: %v", err)
	}
	if _, ok := repo.Categorias[1]; ok {
		t.Fatal("la categoría debía borrarse del repositorio")
	}
}

func TestEliminar_RechazaUnaCategoriaEnUso(t *testing.T) {
	caso, repo := nuevoCaso(categoria.Categoria{ID: 1, Nombre: "Titular A", Activa: true, Docentes: 12, Tramites: 3})
	err := caso.Ejecutar(context.Background(), "u1", 1)
	var enUso *categoria.ErrEnUso
	if !errors.As(err, &enUso) || enUso.Docentes != 12 || enUso.Tramites != 3 {
		t.Fatalf("esperaba ErrEnUso con los conteos, obtuve %v", err)
	}
	if _, ok := repo.Categorias[1]; !ok {
		t.Fatal("una categoría en uso no debe borrarse")
	}
}

func TestEliminar_CategoriaInexistente(t *testing.T) {
	caso, _ := nuevoCaso(categoria.Categoria{ID: 1})
	if err := caso.Ejecutar(context.Background(), "u1", 99); !errors.Is(err, categoria.ErrNoEncontrado) {
		t.Fatalf("esperaba ErrNoEncontrado, obtuve %v", err)
	}
}
