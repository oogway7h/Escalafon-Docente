package reportes

import (
	"strings"
	"time"

	"github.com/xuri/excelize/v2"

	"escalafon/backend/internal/domain/indicador"
)

const (
	hojaDatos   = "Datos"
	hojaResumen = "Resumen"
	filaTabla   = 5 // la cabecera de la tabla va en la fila 5; arriba, título, emisión y filtros
)

// Excel compone un libro con dos hojas: «Datos», con la tabla filtrable (autofiltro y cabecera fija),
// y «Resumen», con la emisión, los filtros y los totales.
func (Exportador) Excel(r *indicador.Reporte) ([]byte, error) {
	f := excelize.NewFile()
	defer f.Close()
	if err := f.SetSheetName("Sheet1", hojaDatos); err != nil {
		return nil, err
	}
	if _, err := f.NewSheet(hojaResumen); err != nil {
		return nil, err
	}
	_ = f.SetDocProps(&excelize.DocProperties{Title: r.Titulo, Creator: r.EmitidoPor, Description: sistema})

	titulo, _ := f.NewStyle(&excelize.Style{Font: &excelize.Font{Bold: true, Size: 14, Color: "0F1E36"}})
	nota, _ := f.NewStyle(&excelize.Style{Font: &excelize.Font{Size: 9, Color: "5C6478"}})
	cabecera, _ := f.NewStyle(&excelize.Style{
		Font:      &excelize.Font{Bold: true, Color: "FFFFFF"},
		Fill:      excelize.Fill{Type: "pattern", Pattern: 1, Color: []string{"0F1E36"}},
		Alignment: &excelize.Alignment{Vertical: "center"},
	})
	formatoFecha, formatoDecimal := "dd/mm/yyyy", "0.0"
	estiloFecha, _ := f.NewStyle(&excelize.Style{CustomNumFmt: &formatoFecha, Alignment: &excelize.Alignment{Horizontal: "center"}})
	estiloDecimal, _ := f.NewStyle(&excelize.Style{CustomNumFmt: &formatoDecimal})
	negrita, _ := f.NewStyle(&excelize.Style{Font: &excelize.Font{Bold: true}})

	// --- Hoja Datos ---
	emision := "Emitido el " + r.EmitidoEn.Format("02/01/2006 15:04") + " por " + r.EmitidoPor + " · " + departamento
	filtros := "Sin filtros: incluye todos los registros."
	if len(r.Filtros) > 0 {
		filtros = "Filtros: " + strings.Join(r.Filtros, " · ")
	}
	for i, v := range []string{r.Titulo, emision, filtros} {
		celda, _ := excelize.CoordinatesToCellName(1, i+1)
		_ = f.SetCellValue(hojaDatos, celda, v)
		estilo := nota
		if i == 0 {
			estilo = titulo
		}
		_ = f.SetCellStyle(hojaDatos, celda, celda, estilo)
	}

	for c, col := range r.Columnas {
		celda, _ := excelize.CoordinatesToCellName(c+1, filaTabla)
		_ = f.SetCellValue(hojaDatos, celda, col.Titulo)
		nombreCol, _ := excelize.ColumnNumberToName(c + 1)
		_ = f.SetColWidth(hojaDatos, nombreCol, nombreCol, max(col.Ancho*1.1, float64(len([]rune(col.Titulo)))+3))
	}
	inicio, _ := excelize.CoordinatesToCellName(1, filaTabla)
	finCab, _ := excelize.CoordinatesToCellName(len(r.Columnas), filaTabla)
	_ = f.SetCellStyle(hojaDatos, inicio, finCab, cabecera)
	_ = f.SetRowHeight(hojaDatos, filaTabla, 20)

	for n, fila := range r.Filas {
		for c, v := range fila {
			if v == nil || c >= len(r.Columnas) {
				continue
			}
			if t, ok := v.(time.Time); ok {
				// Solo el día, sin zona: así Excel no lo corre a otra fecha.
				y, m, d := t.Date()
				v = time.Date(y, m, d, 0, 0, 0, 0, time.UTC)
			}
			celda, _ := excelize.CoordinatesToCellName(c+1, filaTabla+1+n)
			if err := f.SetCellValue(hojaDatos, celda, v); err != nil {
				return nil, err
			}
		}
	}
	ultima := filaTabla + len(r.Filas)
	for c, col := range r.Columnas {
		var estilo int
		switch col.Tipo {
		case indicador.ColFecha:
			estilo = estiloFecha
		case indicador.ColDecimal:
			estilo = estiloDecimal
		default:
			continue
		}
		desde, _ := excelize.CoordinatesToCellName(c+1, filaTabla+1)
		hasta, _ := excelize.CoordinatesToCellName(c+1, ultima)
		_ = f.SetCellStyle(hojaDatos, desde, hasta, estilo)
	}
	fin, _ := excelize.CoordinatesToCellName(len(r.Columnas), ultima)
	if err := f.AutoFilter(hojaDatos, inicio+":"+fin, nil); err != nil {
		return nil, err
	}
	primeraDato, _ := excelize.CoordinatesToCellName(1, filaTabla+1)
	_ = f.SetPanes(hojaDatos, &excelize.Panes{Freeze: true, YSplit: filaTabla, TopLeftCell: primeraDato, ActivePane: "bottomLeft"})

	// --- Hoja Resumen ---
	filasResumen := [][]any{
		{institucion}, {departamento}, {r.Titulo}, {},
		{"Emitido el", r.EmitidoEn.Format("02/01/2006 15:04")}, {"Responsable", r.EmitidoPor}, {"Registros", len(r.Filas)},
	}
	for _, filtro := range r.Filtros {
		filasResumen = append(filasResumen, []any{"Filtro", filtro})
	}
	filasResumen = append(filasResumen, []any{})
	for _, d := range r.Resumen {
		filasResumen = append(filasResumen, []any{d.Etiqueta, d.Valor})
	}
	for i, fila := range filasResumen {
		celda, _ := excelize.CoordinatesToCellName(1, i+1)
		_ = f.SetSheetRow(hojaResumen, celda, &fila)
		if len(fila) == 2 {
			_ = f.SetCellStyle(hojaResumen, celda, celda, negrita)
		}
	}
	_ = f.SetCellStyle(hojaResumen, "A3", "A3", titulo)
	_ = f.SetColWidth(hojaResumen, "A", "A", 34)
	_ = f.SetColWidth(hojaResumen, "B", "B", 60)

	f.SetActiveSheet(0)
	buf, err := f.WriteToBuffer()
	if err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}
