package main

import (
	"os"
	"path/filepath"
	"testing"
)

func TestParsearNombre(t *testing.T) {
	casos := []struct {
		entrada, apellidos, nombres string
		ambiguo                     bool
	}{
		{"ABAL COLUMBA JOHNNY", "ABAL COLUMBA", "JOHNNY", true},
		{"ABREGO ASCARRAGA CARLOS EDUARDO", "ABREGO ASCARRAGA", "CARLOS EDUARDO", false},
		{"ABUJDER BALDERAS MARIA DEL CARMEN", "ABUJDER BALDERAS", "MARIA DEL CARMEN", false},
		{"PEREZ DE LA TORRE JUAN CARLOS", "PEREZ DE LA TORRE", "JUAN CARLOS", false},
		{"DE LA ROCA SOLIZ ANA", "DE LA ROCA SOLIZ", "ANA", true},
		{"AÑEZ ARIAS", "AÑEZ", "ARIAS", false},
		{"SOLO", "SOLO", "", false},
		{"", "", "", false},
	}
	for _, c := range casos {
		ap, no, amb := parsearNombre(c.entrada)
		if ap != c.apellidos || no != c.nombres || amb != c.ambiguo {
			t.Errorf("parsearNombre(%q) = %q, %q, %v; esperaba %q, %q, %v", c.entrada, ap, no, amb, c.apellidos, c.nombres, c.ambiguo)
		}
	}
}

func TestTitulo(t *testing.T) {
	casos := map[string]string{
		"PROFESOR TITULAR A":           "Profesor Titular A",
		"JEFE DE LABORATORIO II":       "Jefe de Laboratorio II",
		"DIRECTOR POST-GRADO FACULTAT": "Director Post-Grado Facultat",
		"MARIA DEL CARMEN":             "Maria del Carmen",
		"AÑEZ":                         "Añez",
	}
	for in, esperado := range casos {
		if got := titulo(in); got != esperado {
			t.Errorf("titulo(%q) = %q; esperaba %q", in, got, esperado)
		}
	}
}

func fl(linea int, codigo, nombre, categ string) fila {
	return fila{Linea: linea, Codigo: codigo, Nombre: nombre, Categ: categ, Gaveta: "A1"}
}

func TestConsolidar(t *testing.T) {
	docentes, conflictos, _ := consolidar([]fila{
		fl(2, "100", "CAZAS GONZALES MILTON PABEL", "PROFESOR TITULAR B"),
		fl(3, "100", "CAZAS GONZALES MILTON PABEL", "PROFESOR TITULAR C"), // misma persona, otra carpeta
		fl(4, "200", "CUELLAR TELLEZ VICENTE REMBERTO", "RECTOR"),
		fl(5, "200", "DAJER ULLOA ERIKA", "PROFESOR TITULAR C"), // código de otra persona
		fl(6, "300", "ACOSTA CLAROS EDY", ""),
	})
	if len(docentes) != 4 {
		t.Fatalf("esperaba 4 docentes, obtuve %d", len(docentes))
	}
	if d := docentes[0]; d.Codigo != "100" || len(d.Filas) != 2 || d.Categ != "PROFESOR TITULAR C" {
		t.Errorf("la persona repetida debe unirse con la mayor categoría: %+v", d)
	}
	if docentes[1].Codigo != "" || docentes[2].Codigo != "" || len(conflictos) != 1 {
		t.Errorf("los códigos en conflicto deben quedar sin código y reportarse: %+v / %v", docentes, conflictos)
	}
}

func TestLegajosDe_IgnoraOkYFilasSinUbicacion(t *testing.T) {
	d := docenteImport{Filas: []fila{
		{Gaveta: "C7", Anexos: "Carp. 1 y 2", Obs: "ok"},
		{Gaveta: "", Anexos: "", Obs: "x"},
		{Gaveta: "C8", Obs: "NO ESTA EN EL DRIVE DE TITULOS"},
	}}
	l := legajosDe(d)
	if len(l) != 2 || l[0].Obs != "" || l[1].Obs == "" {
		t.Fatalf("legajos inesperados: %+v", l)
	}
}

func TestLeerCSV(t *testing.T) {
	ruta := filepath.Join(t.TempDir(), "d.csv")
	contenido := string(rune(0xFEFF)) + ",COD,APELLIDO Y NOMBRE,CATEGORIA ,FACULTAD,ORDEN,GAVETA,NRO ANEXOS ,OBSERVACION,X\n" +
		"1,5948,ABAL COLUMBA JOHNNY,PROFESOR TITULAR B,CARRERA DE DERECHO,1,A1,Carp. 1 y 2 ,OK,#N/A\n" +
		"2,,,,,5,,,,\n" +
		"3,77,SIN CATEGORIA UNO,,,2,A2,,,\n"
	if err := os.WriteFile(ruta, []byte(contenido), 0o600); err != nil {
		t.Fatal(err)
	}
	filas, omitidas, err := leerCSV(ruta)
	if err != nil {
		t.Fatal(err)
	}
	if len(filas) != 2 || len(omitidas) != 1 || filas[0].Anexos != "Carp. 1 y 2" || filas[0].Orden == nil || *filas[0].Orden != 1 {
		t.Fatalf("lectura inesperada: %+v / %v", filas, omitidas)
	}
}
