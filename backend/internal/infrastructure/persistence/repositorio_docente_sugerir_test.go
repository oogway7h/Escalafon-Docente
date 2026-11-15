package persistence

import (
	"context"
	"testing"
)

// Sugerencias del buscador: el código coincide por prefijo y el nombre sin tildes y en cualquier orden.
func TestRepositorioDocente_Sugerir(t *testing.T) {
	dbase := conectarPrueba(t)
	conTx(t, dbase, func(ctx context.Context) {
		codigo := "77123"
		d := DocenteModel{Codigo: &codigo, Nombres: "Ana María", Apellidos: "Gutiérrez Sugerencia"}
		if err := db(ctx, dbase).Create(&d).Error; err != nil {
			t.Fatal(err)
		}
		repo := NewRepositorioDocente(dbase)
		for _, q := range []string{"771", "gutierrez sug", "ana maria gutierrez"} {
			lista, err := repo.Sugerir(ctx, q, 8)
			if err != nil {
				t.Fatal(err)
			}
			hallado := false
			for _, x := range lista {
				hallado = hallado || x.ID == d.ID
			}
			if !hallado {
				t.Fatalf("«%s» debía sugerir al docente: %+v", q, lista)
			}
		}
	})
}
