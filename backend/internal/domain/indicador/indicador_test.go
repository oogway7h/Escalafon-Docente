package indicador_test

import (
	"math"
	"testing"
	"time"

	"escalafon/backend/internal/domain/indicador"
	"escalafon/backend/internal/domain/tramite"
)

var bolivia = time.FixedZone("BOT", -4*60*60)

func fechaSinHora(y int, m time.Month, d int) *time.Time {
	t := time.Date(y, m, d, 0, 0, 0, 0, time.UTC) // así llegan las columnas date desde Postgres
	return &t
}

func TestEvaluarMora(t *testing.T) {
	hoy := time.Date(2026, 10, 25, 10, 0, 0, 0, bolivia)
	casos := []struct {
		nombre       string
		salida       time.Time
		limite       *time.Time
		enMora       bool
		diasMora     int
		motivo       indicador.MotivoMora
		diasPrestado int
	}{
		{"dentro del plazo", time.Date(2026, 10, 20, 9, 0, 0, 0, bolivia), fechaSinHora(2026, 10, 27), false, 0, "", 5},
		{"vence hoy: todavía no está en mora", time.Date(2026, 10, 20, 9, 0, 0, 0, bolivia), fechaSinHora(2026, 10, 25), false, 0, "", 5},
		{"fecha límite vencida", time.Date(2026, 10, 18, 9, 0, 0, 0, bolivia), fechaSinHora(2026, 10, 21), true, 4, indicador.MoraFechaLimite, 7},
		{"supera el máximo aunque la fecha límite sea lejana", time.Date(2026, 10, 1, 9, 0, 0, 0, bolivia), fechaSinHora(2026, 11, 30), true, 9, indicador.MoraExcedeMaximo, 24},
		{"gana el motivo con más días", time.Date(2026, 10, 1, 9, 0, 0, 0, bolivia), fechaSinHora(2026, 10, 5), true, 20, indicador.MoraFechaLimite, 24},
		{"sin fecha límite solo cuenta el máximo", time.Date(2026, 10, 5, 9, 0, 0, 0, bolivia), nil, true, 5, indicador.MoraExcedeMaximo, 20},
		// 23:30 del 9/10 en Bolivia ya es el 10/10 en UTC: el día de salida es el boliviano.
		{"día de salida en hora de Bolivia", time.Date(2026, 10, 9, 23, 30, 0, 0, bolivia).UTC(), nil, true, 1, indicador.MoraExcedeMaximo, 16},
	}
	for _, c := range casos {
		a, enMora := indicador.EvaluarMora(indicador.Prestamo{Salida: c.salida, FechaLimite: c.limite}, hoy, 15)
		if enMora != c.enMora || a.DiasMora != c.diasMora || a.Motivo != c.motivo || a.DiasPrestado != c.diasPrestado {
			t.Fatalf("%s: obtuve enMora=%v diasMora=%d motivo=%q diasPrestado=%d", c.nombre, enMora, a.DiasMora, a.Motivo, a.DiasPrestado)
		}
	}
}

func TestOrdenarAlertas_MasMoraPrimero(t *testing.T) {
	alertas := []indicador.AlertaMora{
		{Prestamo: indicador.Prestamo{LegajoCodigo: "LEG-2"}, DiasMora: 2},
		{Prestamo: indicador.Prestamo{LegajoCodigo: "LEG-9"}, DiasMora: 9},
		{Prestamo: indicador.Prestamo{LegajoCodigo: "LEG-1"}, DiasMora: 2},
	}
	indicador.OrdenarAlertas(alertas)
	if alertas[0].LegajoCodigo != "LEG-9" || alertas[1].LegajoCodigo != "LEG-1" || alertas[2].LegajoCodigo != "LEG-2" {
		t.Fatalf("orden inesperado: %+v", alertas)
	}
}

func TestAgruparTiempos(t *testing.T) {
	base := time.Date(2026, 1, 1, 8, 0, 0, 0, time.UTC)
	resuelto := func(tipo tramite.Tipo, categoria string, dias int) indicador.Resuelto {
		return indicador.Resuelto{Tipo: tipo, Categoria: categoria, RadicadoEn: base, ResueltoEn: base.AddDate(0, 0, dias)}
	}
	rs := []indicador.Resuelto{
		resuelto(tramite.Regularizacion, "Adjunto", 30),
		resuelto(tramite.Categorizacion, "Titular B", 10),
		resuelto(tramite.Categorizacion, "Titular A", 20),
		resuelto(tramite.Categorizacion, "Titular A", 10),
	}
	g := indicador.AgruparTiempos(rs)
	if len(g) != 3 {
		t.Fatalf("esperaba 3 grupos, obtuve %+v", g)
	}
	if g[0].Tipo != tramite.Categorizacion || g[0].Categoria != "Titular A" || g[2].Tipo != tramite.Regularizacion {
		t.Fatalf("los grupos se ordenan por tipo y categoría: %+v", g)
	}
	if g[0].Cantidad != 2 || g[0].PromedioDias != 15 || g[0].MinimoDias != 10 || g[0].MaximoDias != 20 {
		t.Fatalf("estadísticas de Titular A inesperadas: %+v", g[0])
	}
	if p := indicador.PromedioDias(rs); math.Abs(p-17.5) > 1e-9 {
		t.Fatalf("promedio general: esperaba 17.5, obtuve %v", p)
	}
	if indicador.PromedioDias(nil) != 0 || len(indicador.AgruparTiempos(nil)) != 0 {
		t.Fatal("sin resueltos no hay promedio ni grupos")
	}
}

func TestPeriodo_Valido(t *testing.T) {
	a, b := time.Now(), time.Now().Add(time.Hour)
	if !(indicador.Periodo{}).Valido() || !(indicador.Periodo{Desde: &a, Hasta: &b}).Valido() || (indicador.Periodo{Desde: &b, Hasta: &a}).Valido() {
		t.Fatal("solo un rango invertido es inválido")
	}
}
