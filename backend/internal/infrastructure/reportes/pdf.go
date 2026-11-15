// Package reportes compone los archivos de CU-26 a partir de un indicador.Reporte ya calculado:
// PDF con membrete institucional (go-pdf/fpdf) y Excel tabular con filtros (excelize).
package reportes

import (
	"bytes"
	_ "embed"
	"fmt"
	"strings"
	"time"

	"github.com/go-pdf/fpdf"

	"escalafon/backend/internal/domain/indicador"
)

//go:embed logo-ded.png
var logo []byte

const (
	institucion  = "UNIVERSIDAD AUTÓNOMA GABRIEL RENÉ MORENO"
	departamento = "Departamento de Escalafón Docente"
	sistema      = "Sistema de Escalafón Docente — UAGRM"
)

// Colores del panel (app/globals.css): azul de la barra lateral, rojo institucional y gris de fondo.
var (
	azul  = [3]int{15, 30, 54}
	rojo  = [3]int{200, 16, 46}
	gris  = [3]int{92, 100, 120}
	fondo = [3]int{245, 246, 248}
	linea = [3]int{225, 228, 234}
)

// Exportador implementa ports.ExportadorReportes.
type Exportador struct{}

func NewExportador() *Exportador { return &Exportador{} }

// texto convierte una celda a su representación en el reporte (las fechas ya vienen en hora de Bolivia).
func texto(v any) string {
	switch x := v.(type) {
	case nil:
		return ""
	case string:
		return x
	case time.Time:
		return x.Format("02/01/2006")
	case float64:
		return strings.Replace(fmt.Sprintf("%.1f", x), ".", ",", 1)
	default:
		return fmt.Sprint(x)
	}
}

func alineacion(tipo indicador.TipoColumna) string {
	switch tipo {
	case indicador.ColEntero, indicador.ColDecimal:
		return "R"
	case indicador.ColFecha:
		return "C"
	}
	return "L"
}

// PDF compone el reporte en A4 apaisado: membrete en cada página, título, filtros y resumen al
// inicio, tabla con la cabecera repetida en cada página y pie con la numeración.
func (Exportador) PDF(r *indicador.Reporte) ([]byte, error) {
	const margen, altoFila, altoCabecera = 12.0, 6.2, 7.5
	pdf := fpdf.New("L", "mm", "A4", "")
	pdf.SetMargins(margen, 30, margen)
	pdf.SetAutoPageBreak(false, margen)
	pdf.AliasNbPages("")
	pdf.SetTitle(r.Titulo, true)
	pdf.SetAuthor(r.EmitidoPor, true)
	pdf.SetCreator(sistema, true)
	tr := pdf.UnicodeTranslatorFromDescriptor("") // UTF-8 → cp1252 de las fuentes base

	ancho, alto := pdf.GetPageSize()
	util := ancho - 2*margen
	emitido := r.EmitidoEn.Format("02/01/2006 15:04")
	pdf.RegisterImageOptionsReader("logo", fpdf.ImageOptions{ImageType: "PNG"}, bytes.NewReader(logo))

	pdf.SetHeaderFunc(func() {
		pdf.ImageOptions("logo", margen, 8, 0, 15, false, fpdf.ImageOptions{ImageType: "PNG"}, 0, "")
		pdf.SetXY(margen+25, 9)
		pdf.SetTextColor(azul[0], azul[1], azul[2])
		pdf.SetFont("Helvetica", "B", 11)
		pdf.CellFormat(150, 5.5, tr(institucion), "", 2, "L", false, 0, "")
		pdf.SetFont("Helvetica", "", 9.5)
		pdf.SetTextColor(gris[0], gris[1], gris[2])
		pdf.CellFormat(150, 5, tr(departamento), "", 2, "L", false, 0, "")
		pdf.SetXY(ancho-margen-90, 9)
		pdf.SetFont("Helvetica", "", 8.5)
		pdf.CellFormat(90, 5, tr("Emitido: "+emitido), "", 2, "R", false, 0, "")
		pdf.CellFormat(90, 5, tr("Responsable: "+r.EmitidoPor), "", 2, "R", false, 0, "")
		pdf.SetDrawColor(rojo[0], rojo[1], rojo[2])
		pdf.SetLineWidth(0.6)
		pdf.Line(margen, 25, ancho-margen, 25)
		pdf.SetXY(margen, 30)
	})
	pdf.SetFooterFunc(func() {
		pdf.SetY(alto - 9)
		pdf.SetFont("Helvetica", "", 7.5)
		pdf.SetTextColor(gris[0], gris[1], gris[2])
		pdf.CellFormat(util/2, 4, tr(sistema+" · "+r.Titulo), "", 0, "L", false, 0, "")
		pdf.CellFormat(util/2, 4, tr(fmt.Sprintf("Página %d de {nb}", pdf.PageNo())), "", 0, "R", false, 0, "")
	})

	pdf.AddPage()
	pdf.SetTextColor(azul[0], azul[1], azul[2])
	pdf.SetFont("Helvetica", "B", 15)
	pdf.CellFormat(util, 8, tr(r.Titulo), "", 1, "L", false, 0, "")
	pdf.SetFont("Helvetica", "", 9)
	pdf.SetTextColor(gris[0], gris[1], gris[2])
	if len(r.Filtros) > 0 {
		pdf.MultiCell(util, 4.6, tr("Filtros: "+strings.Join(r.Filtros, " · ")), "", "L", false)
	} else {
		pdf.CellFormat(util, 4.6, tr("Sin filtros: incluye todos los registros."), "", 1, "L", false, 0, "")
	}
	pdf.Ln(2.5)

	// Resumen en una grilla de tres columnas: etiqueta pequeña y valor en negrita.
	const porFila, altoDato = 3, 10.5
	anchoDato := util / porFila
	for i, d := range r.Resumen {
		x := margen + float64(i%porFila)*anchoDato
		if i%porFila == 0 {
			if i > 0 {
				pdf.Ln(altoDato)
			}
			if pdf.GetY()+altoDato > alto-margen-10 {
				pdf.AddPage()
			}
		}
		y := pdf.GetY()
		pdf.SetFillColor(fondo[0], fondo[1], fondo[2])
		pdf.Rect(x, y, anchoDato-2, altoDato-1.5, "F")
		pdf.SetXY(x+2.5, y+1)
		pdf.SetFont("Helvetica", "", 7.5)
		pdf.SetTextColor(gris[0], gris[1], gris[2])
		pdf.CellFormat(anchoDato-6, 3.8, tr(d.Etiqueta), "", 2, "L", false, 0, "")
		pdf.SetFont("Helvetica", "B", 9.5)
		pdf.SetTextColor(azul[0], azul[1], azul[2])
		pdf.CellFormat(anchoDato-6, 4.6, tr(recortar(pdf, tr, d.Valor, anchoDato-6)), "", 0, "L", false, 0, "")
		pdf.SetXY(margen, y)
	}
	if len(r.Resumen) > 0 {
		pdf.Ln(altoDato + 2)
	}

	// Tabla: anchos proporcionales al peso de cada columna.
	var pesos float64
	for _, c := range r.Columnas {
		pesos += c.Ancho
	}
	anchos := make([]float64, len(r.Columnas))
	for i, c := range r.Columnas {
		anchos[i] = util * c.Ancho / pesos
	}
	cabecera := func() {
		pdf.SetFont("Helvetica", "B", 8)
		pdf.SetFillColor(azul[0], azul[1], azul[2])
		pdf.SetTextColor(255, 255, 255)
		for i, c := range r.Columnas {
			pdf.CellFormat(anchos[i], altoCabecera, tr(c.Titulo), "", 0, alineacion(c.Tipo), true, 0, "")
		}
		pdf.Ln(-1)
		pdf.SetFont("Helvetica", "", 7.8)
		pdf.SetTextColor(azul[0], azul[1], azul[2])
		pdf.SetDrawColor(linea[0], linea[1], linea[2]) // el membrete deja el rojo y una línea gruesa
		pdf.SetLineWidth(0.2)
	}
	cabecera()
	for n, fila := range r.Filas {
		if pdf.GetY()+altoFila > alto-margen-6 {
			pdf.AddPage()
			cabecera()
		}
		pdf.SetFillColor(fondo[0], fondo[1], fondo[2])
		for i, c := range r.Columnas {
			var v any
			if i < len(fila) {
				v = fila[i]
			}
			celda := recortar(pdf, tr, texto(v), anchos[i]-2)
			pdf.CellFormat(anchos[i], altoFila, tr(celda), "B", 0, alineacion(c.Tipo), n%2 == 1, 0, "")
		}
		pdf.Ln(-1)
	}
	pdf.Ln(2)
	pdf.SetFont("Helvetica", "", 8)
	pdf.SetTextColor(gris[0], gris[1], gris[2])
	pdf.CellFormat(util, 4.5, tr(fmt.Sprintf("%d registros. Fechas en hora de Bolivia.", len(r.Filas))), "", 1, "L", false, 0, "")

	var buf bytes.Buffer
	if err := pdf.Output(&buf); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

// recortar acorta s con "…" para que quepa en el ancho dado con la fuente actual (las celdas no se parten).
func recortar(pdf *fpdf.Fpdf, tr func(string) string, s string, ancho float64) string {
	if pdf.GetStringWidth(tr(s)) <= ancho {
		return s
	}
	runas := []rune(s)
	for len(runas) > 0 && pdf.GetStringWidth(tr(string(runas)+"…")) > ancho {
		runas = runas[:len(runas)-1]
	}
	return strings.TrimSpace(string(runas)) + "…"
}
