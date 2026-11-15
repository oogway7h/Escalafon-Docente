// Comando importdocentes: carga una sola vez el padrón docente desde un CSV (columnas COD,
// APELLIDO Y NOMBRE, CATEGORIA, FACULTAD, ORDEN, GAVETA, NRO ANEXOS, OBSERVACION). Es idempotente:
// repetirlo no duplica nada. Con -dry-run no toca la base de datos y solo informa.
//
//	make import-docentes file=../archivos/docsAyuda/datos.csv
//	make import-docentes-dry file=../archivos/docsAyuda/datos.csv
package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"sort"
	"strings"

	"gorm.io/gorm"

	"escalafon/backend/internal/infrastructure/config"
	"escalafon/backend/internal/infrastructure/persistence"
)

// docenteImport es un docente ya consolidado: todas las filas de un mismo código y nombre.
type docenteImport struct {
	Codigo    string // vacío = sin código (conflicto): se corrige a mano desde el panel
	Apellidos string
	Nombres   string
	Categ     string
	Facultad  string
	Filas     []fila
}

func (d docenteImport) clave() string { return persistence.Normalizar(d.Apellidos + " " + d.Nombres) }

func main() {
	archivo := flag.String("file", "", "ruta del CSV")
	seco := flag.Bool("dry-run", false, "no escribe en la base de datos; solo informa")
	reporte := flag.String("reporte", "reporte-importacion.txt", "archivo donde se escribe el detalle (conflictos, omitidos, nombres ambiguos)")
	flag.Parse()
	if *archivo == "" {
		log.Fatal("indica el CSV: -file=ruta/datos.csv")
	}

	filas, omitidas, err := leerCSV(*archivo)
	if err != nil {
		log.Fatalf("no se pudo leer el CSV: %v", err)
	}
	docentes, conflictos, ambiguos := consolidar(filas)

	categorias, facultades, legajos := map[string]string{}, map[string]bool{}, 0
	for _, d := range docentes {
		if d.Categ != "" {
			categorias[persistence.ClaveCategoria(d.Categ)] = d.Categ
		}
		if d.Facultad != "" {
			facultades[persistence.Normalizar(d.Facultad)] = true
		}
		legajos += len(legajosDe(d))
	}

	detalle := armarReporte(omitidas, conflictos, ambiguos)
	if err := os.WriteFile(*reporte, []byte(detalle), 0o644); err != nil {
		log.Printf("no se pudo escribir %s: %v", *reporte, err)
	}
	log.Printf("CSV: %d filas con datos, %d omitidas", len(filas), len(omitidas))
	log.Printf("a importar: %d docentes, %d categorías, %d facultades, %d legajos", len(docentes), len(categorias), len(facultades), legajos)
	log.Printf("revisar: %d códigos en conflicto, %d nombres ambiguos (detalle en %s)", len(conflictos), len(ambiguos), *reporte)
	if *seco {
		log.Println("-dry-run: no se escribió nada en la base de datos")
		return
	}

	cfg, err := config.Cargar()
	if err != nil {
		log.Fatalf("configuración inválida: %v", err)
	}
	db, err := persistence.Conectar(cfg.DatabaseURL)
	if err != nil {
		log.Fatalf("no se pudo conectar a la base de datos: %v", err)
	}
	if err := db.Transaction(func(tx *gorm.DB) error { return importar(tx, docentes) }); err != nil {
		log.Fatalf("la importación se revirtió: %v", err)
	}
}

// consolidar agrupa las filas por código. Mismo código y mismo nombre = un docente con varias
// carpetas (se queda la categoría de mayor jerarquía). Mismo código con nombres distintos =
// conflicto: se importan todos sin código para corregirlos a mano.
func consolidar(filas []fila) (docentes []docenteImport, conflictos, ambiguos []string) {
	porCodigo := map[string][]fila{}
	var orden []string
	for _, f := range filas {
		if _, ok := porCodigo[f.Codigo]; !ok {
			orden = append(orden, f.Codigo)
		}
		porCodigo[f.Codigo] = append(porCodigo[f.Codigo], f)
	}

	nuevo := func(codigo string, grupo []fila) docenteImport {
		ap, no, amb := parsearNombre(grupo[0].Nombre)
		d := docenteImport{Codigo: codigo, Apellidos: titulo(ap), Nombres: titulo(no), Filas: grupo}
		for _, f := range grupo {
			if mejor(f.Categ, d.Categ) {
				d.Categ = f.Categ
			}
			if d.Facultad == "" {
				d.Facultad = f.Facultad
			}
		}
		if amb {
			ambiguos = append(ambiguos, fmt.Sprintf("%s  →  apellidos «%s» · nombres «%s»", grupo[0].Nombre, d.Apellidos, d.Nombres))
		}
		return d
	}

	for _, codigo := range orden {
		grupo := porCodigo[codigo]
		porNombre := map[string][]fila{}
		var nombres []string
		for _, f := range grupo {
			k := persistence.Normalizar(f.Nombre)
			if _, ok := porNombre[k]; !ok {
				nombres = append(nombres, k)
			}
			porNombre[k] = append(porNombre[k], f)
		}
		if len(nombres) == 1 {
			docentes = append(docentes, nuevo(codigo, grupo))
			continue
		}
		var quienes []string
		for _, k := range nombres {
			docentes = append(docentes, nuevo("", porNombre[k]))
			quienes = append(quienes, fmt.Sprintf("%s (línea %d)", porNombre[k][0].Nombre, porNombre[k][0].Linea))
		}
		conflictos = append(conflictos, fmt.Sprintf("código %s usado por personas distintas: %s", codigo, strings.Join(quienes, " / ")))
	}
	return docentes, conflictos, ambiguos
}

// mejor indica si la categoría a tiene mayor jerarquía que b (una categoría vacía pierde siempre).
func mejor(a, b string) bool {
	if a == "" {
		return false
	}
	if b == "" {
		return true
	}
	return jerarquia[persistence.ClaveCategoria(a)] > jerarquia[persistence.ClaveCategoria(b)]
}

type legajoImport struct {
	Gaveta, Anexos, Obs string
	Orden               *int32
}

// legajosDe devuelve una carpeta por fila que tenga gaveta o anexos.
func legajosDe(d docenteImport) []legajoImport {
	var out []legajoImport
	for _, f := range d.Filas {
		if f.Gaveta == "" && f.Anexos == "" {
			continue
		}
		obs := f.Obs
		if esObservacionOk(obs) {
			obs = ""
		}
		out = append(out, legajoImport{Gaveta: f.Gaveta, Anexos: f.Anexos, Obs: obs, Orden: f.Orden})
	}
	return out
}

func armarReporte(omitidas, conflictos, ambiguos []string) string {
	var b strings.Builder
	seccion := func(titulo string, items []string) {
		fmt.Fprintf(&b, "== %s (%d) ==\n", titulo, len(items))
		for _, i := range items {
			fmt.Fprintln(&b, i)
		}
		b.WriteString("\n")
	}
	seccion("Filas omitidas", omitidas)
	seccion("Códigos en conflicto (importados sin código: asignarlo desde el panel)", conflictos)
	seccion("Nombres con corte ambiguo entre apellidos y nombres (revisar)", ambiguos)
	return b.String()
}

func ptr(s string) *string {
	if s == "" {
		return nil
	}
	return &s
}

func importar(tx *gorm.DB, docentes []docenteImport) error {
	// Facultades: se crean las que faltan, en lote.
	var existentesFac []persistence.FacultadModel
	if err := tx.Find(&existentesFac).Error; err != nil {
		return err
	}
	facultadID := map[string]int16{}
	for _, f := range existentesFac {
		facultadID[persistence.Normalizar(f.Nombre)] = f.ID
	}
	var nuevasFac []persistence.FacultadModel
	vistas := map[string]bool{}
	for _, d := range docentes {
		k := persistence.Normalizar(d.Facultad)
		if d.Facultad == "" || vistas[k] {
			continue
		}
		vistas[k] = true
		if _, ok := facultadID[k]; !ok {
			nuevasFac = append(nuevasFac, persistence.FacultadModel{Nombre: d.Facultad})
		}
	}
	if len(nuevasFac) > 0 {
		if err := tx.CreateInBatches(&nuevasFac, 100).Error; err != nil {
			return err
		}
		for _, f := range nuevasFac {
			facultadID[persistence.Normalizar(f.Nombre)] = f.ID
		}
	}

	// Categorías: se reconocen por clave; las de escalafón nuevas quedan activas y los cargos inactivos.
	var existentesCat []persistence.CategoriaModel
	if err := tx.Find(&existentesCat).Error; err != nil {
		return err
	}
	categoriaID := map[string]int16{}
	for _, c := range existentesCat {
		categoriaID[persistence.ClaveCategoria(c.Nombre)] = c.ID
	}
	resolverCategoria := func(texto string) (*int16, error) {
		if texto == "" {
			return nil, nil
		}
		clave := persistence.ClaveCategoria(texto)
		if id, ok := categoriaID[clave]; ok {
			return &id, nil
		}
		activa := esEscalafon(clave)
		c := persistence.CategoriaModel{Nombre: titulo(texto), Activa: activa}
		if err := tx.Create(&c).Error; err != nil {
			return nil, err
		}
		if !activa { // GORM omite false en columnas con default:true (y devuelve el default): se fuerza aparte
			if err := tx.Model(&c).Update("activa", false).Error; err != nil {
				return nil, err
			}
		}
		categoriaID[clave] = c.ID
		return &c.ID, nil
	}

	// Docentes existentes, por código o (si no tienen código) por nombre.
	var existentes []persistence.DocenteModel
	if err := tx.Find(&existentes).Error; err != nil {
		return err
	}
	porCodigo := map[string]*persistence.DocenteModel{}
	sinCodigo := map[string]*persistence.DocenteModel{}
	for i := range existentes {
		e := &existentes[i]
		if e.Codigo != nil {
			porCodigo[*e.Codigo] = e
		} else {
			sinCodigo[persistence.Normalizar(e.Apellidos+" "+e.Nombres)] = e
		}
	}

	type pendiente struct {
		modelo *persistence.DocenteModel
		datos  docenteImport
		nuevo  bool
	}
	var todos []pendiente
	var nuevos []persistence.DocenteModel
	actualizados := 0
	for _, d := range docentes {
		catID, err := resolverCategoria(d.Categ)
		if err != nil {
			return err
		}
		var facID *int16
		if d.Facultad != "" {
			id := facultadID[persistence.Normalizar(d.Facultad)]
			facID = &id
		}
		var e *persistence.DocenteModel
		if d.Codigo != "" {
			e = porCodigo[d.Codigo]
		} else {
			e = sinCodigo[d.clave()]
		}
		if e == nil {
			nuevos = append(nuevos, persistence.DocenteModel{
				Codigo: ptr(d.Codigo), Nombres: d.Nombres, Apellidos: d.Apellidos, CategoriaID: catID, FacultadID: facID, Activo: true,
			})
			todos = append(todos, pendiente{datos: d, nuevo: true})
			continue
		}
		cambios := map[string]any{}
		if e.Nombres != d.Nombres {
			cambios["nombres"] = d.Nombres
		}
		if e.Apellidos != d.Apellidos {
			cambios["apellidos"] = d.Apellidos
		}
		if catID != nil && (e.CategoriaID == nil || *e.CategoriaID != *catID) {
			cambios["categoria_id"] = *catID
		}
		if facID != nil && (e.FacultadID == nil || *e.FacultadID != *facID) {
			cambios["facultad_id"] = *facID
		}
		if len(cambios) > 0 {
			if err := tx.Model(&persistence.DocenteModel{}).Where("id = ?", e.ID).Updates(cambios).Error; err != nil {
				return err
			}
			actualizados++
		}
		todos = append(todos, pendiente{modelo: e, datos: d})
	}
	if len(nuevos) > 0 {
		if err := tx.CreateInBatches(&nuevos, 200).Error; err != nil {
			return err
		}
	}

	// Legajos: el CSV manda, así que se reemplazan los de los docentes importados.
	var ids []string
	var legajos []persistence.LegajoModel
	iNuevo := 0
	for _, p := range todos {
		var id string
		if p.nuevo {
			id = nuevos[iNuevo].ID
			iNuevo++
		} else {
			id = p.modelo.ID
			ids = append(ids, id)
		}
		for _, l := range legajosDe(p.datos) {
			legajos = append(legajos, persistence.LegajoModel{
				DocenteID: id, Gaveta: ptr(l.Gaveta), Orden: l.Orden, Anexos: ptr(l.Anexos), Observacion: ptr(l.Obs),
			})
		}
	}
	for i := 0; i < len(ids); i += 500 {
		fin := min(i+500, len(ids))
		if err := tx.Where("docente_id IN ?", ids[i:fin]).Delete(&persistence.LegajoModel{}).Error; err != nil {
			return err
		}
	}
	if len(legajos) > 0 {
		if err := tx.CreateInBatches(&legajos, 500).Error; err != nil {
			return err
		}
	}

	nombres := make([]string, 0, len(categoriaID))
	for k := range categoriaID {
		nombres = append(nombres, k)
	}
	sort.Strings(nombres)
	log.Printf("importación terminada: %d docentes nuevos, %d actualizados, %d legajos, %d facultades nuevas, %d categorías en total",
		len(nuevos), actualizados, len(legajos), len(nuevasFac), len(nombres))
	return nil
}
