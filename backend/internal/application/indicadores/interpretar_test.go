package indicadores_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"escalafon/backend/internal/application/indicadores"
	"escalafon/backend/internal/application/ports"
	"escalafon/backend/internal/application/testdobles"
	"escalafon/backend/internal/domain/indicador"
)

type interpreteFalso struct {
	salida *ports.PeticionReporte
	err    error
}

func (i interpreteFalso) Interpretar(context.Context, string, time.Time, []ports.CategoriaRef) (*ports.PeticionReporte, error) {
	return i.salida, i.err
}

func interpretar(t *testing.T, i ports.InterpreteReportes, frase string) (*ports.PeticionReporte, error) {
	t.Helper()
	caso := &indicadores.InterpretarReporte{Interprete: i, Reloj: &testdobles.RelojFijo{Momento: ahora}}
	return caso.Ejecutar(context.Background(), frase, []ports.CategoriaRef{{ID: 3, Nombre: "Titular"}})
}

func TestInterpretar_SaneaLaSalidaDelModelo(t *testing.T) {
	p, err := interpretar(t, interpreteFalso{salida: &ports.PeticionReporte{
		Tipo: "TRAMITES", Formato: "DOCX", Desde: "2026-10-01", Hasta: "no-es-fecha", Estado: "OBSERVADO",
		TipoTramite: "INVENTADO", CategoriaID: 99, EstadoLegajo: "PRESTADO",
	}}, "trámites observados de octubre")
	if err != nil {
		t.Fatal(err)
	}
	if p.Formato != "PDF" || p.Desde != "2026-10-01" || p.Hasta != "" || p.Estado != "OBSERVADO" ||
		p.TipoTramite != "" || p.CategoriaID != 0 || p.EstadoLegajo != "" {
		t.Fatalf("salida sin sanear: %+v", p)
	}
}

func TestInterpretar_CategoriaConocida(t *testing.T) {
	p, err := interpretar(t, interpreteFalso{salida: &ports.PeticionReporte{Tipo: "TIEMPOS", Formato: "XLSX", CategoriaID: 3}}, "tiempos de titulares")
	if err != nil || p.CategoriaID != 3 || p.Formato != "XLSX" {
		t.Fatalf("p=%+v err=%v", p, err)
	}
}

func TestInterpretar_Errores(t *testing.T) {
	if _, err := interpretar(t, interpreteFalso{salida: &ports.PeticionReporte{Tipo: "OTRO"}}, "hola"); !errors.Is(err, indicador.ErrNoEntendido) {
		t.Fatalf("tipo inválido: %v", err)
	}
	if _, err := interpretar(t, interpreteFalso{}, "  "); !errors.Is(err, indicador.ErrNoEntendido) {
		t.Fatalf("frase vacía: %v", err)
	}
	caso := &indicadores.InterpretarReporte{Reloj: &testdobles.RelojFijo{Momento: ahora}}
	if _, err := caso.Ejecutar(context.Background(), "x", nil); !errors.Is(err, indicador.ErrInterpreteNoDisponible) {
		t.Fatalf("sin intérprete: %v", err)
	}
}
