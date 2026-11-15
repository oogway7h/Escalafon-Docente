package reportes

import (
	"bytes"
	"fmt"
	"testing"
	"time"

	"github.com/xuri/excelize/v2"

	"escalafon/backend/internal/domain/indicador"
)

var bolivia = time.FixedZone("BOT", -4*60*60)

// reporteDePrueba tiene tildes, eñes, celdas vacías y filas suficientes para ocupar varias páginas.
func reporteDePrueba(filas int) *indicador.Reporte {
	r := &indicador.Reporte{
		Tipo: indicador.ReporteTiempos, Titulo: "Tiempos de atención", EmitidoPor: "Rosa Peñaranda", EmitidoEn: time.Date(2026, 10, 25, 23, 30, 0, 0, bolivia),
		Filtros: []string{"Resueltos del 01/01/2026 al 25/10/2026", "Tipo: Categorización"},
		Resumen: []indicador.Dato{{Etiqueta: "Trámites resueltos", Valor: fmt.Sprint(filas)}, {Etiqueta: "Tiempo promedio", Valor: "12,5 días"}},
		Columnas: []indicador.Columna{
			{Titulo: "Código", Tipo: indicador.ColTexto, Ancho: 15}, {Titulo: "Docente", Tipo: indicador.ColTexto, Ancho: 26},
			{Titulo: "Orden", Tipo: indicador.ColEntero, Ancho: 6}, {Titulo: "Radicado", Tipo: indicador.ColFecha, Ancho: 10},
			{Titulo: "Días", Tipo: indicador.ColDecimal, Ancho: 7},
		},
	}
	for i := 0; i < filas; i++ {
		r.Filas = append(r.Filas, []any{
			fmt.Sprintf("ED-2026-%06d", i), "José María Ñúñez de la Peña con un nombre bastante largo", i,
			time.Date(2026, 3, 9, 23, 0, 0, 0, bolivia), 12.25, // 23:00 de Bolivia: en UTC ya es el 10/03
		})
	}
	r.Filas[0][2] = nil
	return r
}

func TestPDF_VariasPaginas(t *testing.T) {
	datos, err := NewExportador().PDF(reporteDePrueba(120))
	if err != nil {
		t.Fatalf("no esperaba error: %v", err)
	}
	if !bytes.HasPrefix(datos, []byte("%PDF-")) {
		t.Fatalf("el archivo no es un PDF")
	}
	if n := bytes.Count(datos, []byte("/Type /Page")) - bytes.Count(datos, []byte("/Type /Pages")); n < 3 {
		t.Fatalf("120 filas deberían ocupar varias páginas, obtuve %d", n)
	}
}

func TestExcel_TablaFiltrableConTiposCorrectos(t *testing.T) {
	datos, err := NewExportador().Excel(reporteDePrueba(3))
	if err != nil {
		t.Fatalf("no esperaba error: %v", err)
	}
	f, err := excelize.OpenReader(bytes.NewReader(datos))
	if err != nil {
		t.Fatalf("el libro no se puede abrir: %v", err)
	}
	defer f.Close()
	if hojas := f.GetSheetList(); len(hojas) != 2 || hojas[0] != hojaDatos || hojas[1] != hojaResumen {
		t.Fatalf("hojas inesperadas: %v", hojas)
	}
	if v, _ := f.GetCellValue(hojaDatos, "A1"); v != "Tiempos de atención" {
		t.Fatalf("título: %q", v)
	}
	if v, _ := f.GetCellValue(hojaDatos, "B5"); v != "Docente" {
		t.Fatalf("la cabecera de la tabla va en la fila 5: %q", v)
	}
	if v, _ := f.GetCellValue(hojaDatos, "D6"); v != "09/03/2026" {
		t.Fatalf("las fechas conservan el día de Bolivia: %q", v)
	}
	if v, _ := f.GetCellValue(hojaDatos, "E6", excelize.Options{RawCellValue: true}); v != "12.25" {
		t.Fatalf("los decimales se guardan como número: %q", v)
	}
	if v, _ := f.GetCellValue(hojaDatos, "C6"); v != "" {
		t.Fatalf("una celda nil queda vacía: %q", v)
	}
	if v, _ := f.GetCellValue(hojaResumen, "B6"); v != "Rosa Peñaranda" {
		t.Fatalf("el responsable figura en el resumen: %q", v)
	}
}
