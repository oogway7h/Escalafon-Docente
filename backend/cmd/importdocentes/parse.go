package main

import (
	"encoding/csv"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"

	"escalafon/backend/internal/infrastructure/persistence"
)

// fila es una línea útil del CSV.
type fila struct {
	Linea    int // número de línea en el archivo (1 = encabezado)
	Codigo   string
	Nombre   string
	Categ    string
	Facultad string
	Orden    *int32
	Gaveta   string
	Anexos   string
	Obs      string
}

var particulas = map[string]bool{"DE": true, "DEL": true, "LA": true, "LAS": true, "LOS": true, "SAN": true, "SANTA": true, "DA": true, "DOS": true, "VAN": true, "VON": true}

// parsearNombre separa "APELLIDOS NOMBRES" (el CSV no marca la frontera): 2 apellidos y el resto
// nombres; las partículas (DE, DEL, LA…) se unen al apellido que las sigue. ambiguo es true
// cuando el número de palabras no permite un corte confiable (3 o 5 o más).
func parsearNombre(completo string) (apellidos, nombres string, ambiguo bool) {
	tokens := strings.Fields(completo)
	var unidades []string
	var pendiente []string
	for _, t := range tokens {
		pendiente = append(pendiente, t)
		if !particulas[strings.ToUpper(t)] {
			unidades = append(unidades, strings.Join(pendiente, " "))
			pendiente = nil
		}
	}
	if len(pendiente) > 0 { // terminó en partícula: se conserva como unidad propia
		unidades = append(unidades, strings.Join(pendiente, " "))
	}
	switch n := len(unidades); {
	case n == 0:
		return "", "", false
	case n == 1:
		return unidades[0], "", false
	case n == 2:
		return unidades[0], unidades[1], false
	default:
		return unidades[0] + " " + unidades[1], strings.Join(unidades[2:], " "), n != 4
	}
}

var minusculas = map[string]bool{"de": true, "del": true, "la": true, "las": true, "los": true, "y": true, "e": true, "en": true, "da": true, "dos": true, "van": true, "von": true}

// titulo pasa "PROFESOR TITULAR A" a "Profesor Titular A" (letras sueltas y números romanos en mayúscula).
func titulo(s string) string {
	palabras := strings.Fields(strings.ToLower(s))
	for i, p := range palabras {
		partes := strings.Split(p, "-")
		for j, parte := range partes {
			r := []rune(parte)
			switch {
			case len(r) == 0:
			case len(r) == 1 && i > 0:
				partes[j] = strings.ToUpper(parte)
			case parte == "ii" || parte == "iii" || parte == "iv":
				partes[j] = strings.ToUpper(parte)
			case minusculas[parte] && i > 0:
			default:
				partes[j] = strings.ToUpper(string(r[0])) + string(r[1:])
			}
		}
		palabras[i] = strings.Join(partes, "-")
	}
	return strings.Join(palabras, " ")
}

// jerarquia ordena las categorías de escalafón; las demás (cargos) valen 0.
var jerarquia = map[string]int{"INTERINO": 1, "ADJUNTO": 2, "ASOCIADO": 3, "TITULAR A": 4, "TITULAR B": 5, "TITULAR C": 6, "EMERITO": 7}

func esEscalafon(clave string) bool { _, ok := jerarquia[clave]; return ok }

// esObservacionOk indica si la observación es solo el "ok" de revisión (no se guarda).
func esObservacionOk(s string) bool {
	n := strings.ToLower(strings.TrimSpace(s))
	return n == "" || n == "ok" || n == "oks"
}

// leerCSV devuelve las filas con datos y las líneas omitidas con su motivo.
func leerCSV(ruta string) (filas []fila, omitidas []string, err error) {
	f, err := os.Open(ruta)
	if err != nil {
		return nil, nil, err
	}
	defer f.Close()

	lector := csv.NewReader(f)
	lector.FieldsPerRecord = -1
	lector.LazyQuotes = true

	encabezado, err := lector.Read()
	if err != nil {
		return nil, nil, fmt.Errorf("no se pudo leer el encabezado: %w", err)
	}
	col := map[string]int{}
	for i, h := range encabezado {
		col[persistence.Normalizar(strings.TrimPrefix(h, string(rune(0xFEFF))))] = i
	}
	for _, requerida := range []string{"COD", "APELLIDO Y NOMBRE", "CATEGORIA", "FACULTAD", "ORDEN", "GAVETA", "NRO ANEXOS", "OBSERVACION"} {
		if _, ok := col[requerida]; !ok {
			return nil, nil, fmt.Errorf("falta la columna %q en el encabezado", requerida)
		}
	}

	linea := 1
	for {
		reg, err := lector.Read()
		if err == io.EOF {
			break
		}
		linea++
		if err != nil {
			omitidas = append(omitidas, fmt.Sprintf("línea %d: no se pudo leer (%v)", linea, err))
			continue
		}
		campo := func(nombre string) string {
			if i := col[nombre]; i < len(reg) {
				return strings.Join(strings.Fields(reg[i]), " ")
			}
			return ""
		}
		fl := fila{
			Linea: linea, Codigo: campo("COD"), Nombre: campo("APELLIDO Y NOMBRE"), Categ: campo("CATEGORIA"),
			Facultad: campo("FACULTAD"), Gaveta: campo("GAVETA"), Anexos: campo("NRO ANEXOS"), Obs: campo("OBSERVACION"),
		}
		if n, err := strconv.Atoi(campo("ORDEN")); err == nil {
			o := int32(n)
			fl.Orden = &o
		}
		switch {
		case fl.Codigo == "" && fl.Nombre == "":
			omitidas = append(omitidas, fmt.Sprintf("línea %d: fila vacía", linea))
		case fl.Codigo == "":
			omitidas = append(omitidas, fmt.Sprintf("línea %d: %s no tiene código", linea, fl.Nombre))
		case fl.Nombre == "":
			omitidas = append(omitidas, fmt.Sprintf("línea %d: el código %s no tiene nombre", linea, fl.Codigo))
		default:
			filas = append(filas, fl)
		}
	}
	return filas, omitidas, nil
}
