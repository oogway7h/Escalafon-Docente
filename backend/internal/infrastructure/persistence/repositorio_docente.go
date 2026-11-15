package persistence

import (
	"context"
	"errors"
	"strings"

	"gorm.io/gorm"

	"escalafon/backend/internal/domain/docente"
)

type RepositorioDocente struct {
	DB *gorm.DB
}

func NewRepositorioDocente(dbase *gorm.DB) *RepositorioDocente { return &RepositorioDocente{DB: dbase} }

func aDominioDocente(m *DocenteModel) docente.Docente {
	d := docente.Docente{
		ID: m.ID, Codigo: val(m.Codigo), Cedula: val(m.Cedula), Nombres: m.Nombres, Apellidos: m.Apellidos,
		Correo: m.Correo, Telefono: m.Telefono, CategoriaID: m.CategoriaID, FacultadID: m.FacultadID,
		Activo: m.Activo, CreadoEn: m.CreadoEn,
	}
	if m.Categoria != nil {
		d.Categoria = m.Categoria.Nombre
	}
	if m.Facultad != nil {
		d.Facultad = m.Facultad.Nombre
	}
	return d
}

// escaparLike neutraliza los comodines de LIKE en el texto del usuario: sin esto "%" o "_" se interpretan como
// patrón (no es inyección, pero permite búsquedas costosas y resultados inesperados).
var escaparLike = strings.NewReplacer(`\`, `\\`, `%`, `\%`, `_`, `\_`).Replace

// tieneDigitos distingue una búsqueda por código o cédula de una por nombre.
func tieneDigitos(s string) bool { return strings.ContainsAny(s, "0123456789") }

// cargarLegajos adjunta las carpetas físicas de los docentes dados.
func (r *RepositorioDocente) cargarLegajos(ctx context.Context, docentes []docente.Docente) error {
	if len(docentes) == 0 {
		return nil
	}
	ids := make([]string, len(docentes))
	idx := map[string]int{}
	for i, d := range docentes {
		ids[i], idx[d.ID] = d.ID, i
	}
	var legajos []LegajoModel
	if err := db(ctx, r.DB).Where("docente_id IN ?", ids).Order("id").Find(&legajos).Error; err != nil {
		return err
	}
	for _, l := range legajos {
		var orden int
		if l.Orden != nil {
			orden = int(*l.Orden)
		}
		i := idx[l.DocenteID]
		docentes[i].Legajos = append(docentes[i].Legajos, docente.Legajo{
			Gaveta: val(l.Gaveta), Orden: orden, Anexos: val(l.Anexos), Observacion: val(l.Observacion),
		})
	}
	return nil
}

func (r *RepositorioDocente) Sugerir(ctx context.Context, texto string, limite int) ([]docente.Docente, error) {
	b := strings.ToLower(strings.TrimSpace(texto))
	q := db(ctx, r.DB).Model(&DocenteModel{}).Preload("Categoria").Where("activo")
	if tieneDigitos(b) {
		q = q.Where("lower(codigo) LIKE ? OR lower(cedula) LIKE ?", escaparLike(b)+"%", escaparLike(b)+"%").Order("length(codigo), codigo")
	} else {
		like := "%" + escaparLike(b) + "%"
		q = q.Where("unaccent(lower(apellidos || ' ' || nombres)) LIKE unaccent(?) OR unaccent(lower(nombres || ' ' || apellidos)) LIKE unaccent(?)", like, like).
			Order("apellidos, nombres")
	}
	var ms []DocenteModel
	if err := q.Limit(limite).Find(&ms).Error; err != nil {
		return nil, err
	}
	out := make([]docente.Docente, len(ms))
	for i := range ms {
		out[i] = aDominioDocente(&ms[i])
	}
	return out, nil
}

// filtrarDocentes aplica el filtro de activos y la búsqueda: código y cédula coinciden exactos
// (buscar "123" no trae "1234"); sin dígitos, el nombre coincide por fragmento.
func filtrarDocentes(q *gorm.DB, busqueda string, soloActivos bool) *gorm.DB {
	if soloActivos {
		q = q.Where("activo")
	}
	if b := strings.TrimSpace(busqueda); b != "" {
		b = strings.ToLower(b)
		if tieneDigitos(b) { // parece un código o una cédula: coincidencia exacta, sin parecidos
			q = q.Where("lower(codigo) = ? OR lower(cedula) = ?", b, b)
		} else {
			q = q.Where("lower(apellidos || ' ' || nombres) LIKE ?", "%"+escaparLike(b)+"%")
		}
	}
	return q
}

func (r *RepositorioDocente) hidratar(ctx context.Context, modelos []DocenteModel) ([]docente.Docente, error) {
	out := make([]docente.Docente, len(modelos))
	for i := range modelos {
		out[i] = aDominioDocente(&modelos[i])
	}
	if err := r.cargarLegajos(ctx, out); err != nil {
		return nil, err
	}
	return out, nil
}

func (r *RepositorioDocente) Listar(ctx context.Context, busqueda string, soloActivos bool) ([]docente.Docente, error) {
	var modelos []DocenteModel
	q := filtrarDocentes(db(ctx, r.DB).Preload("Categoria").Preload("Facultad").Order("apellidos, nombres"), busqueda, soloActivos)
	if err := q.Find(&modelos).Error; err != nil {
		return nil, err
	}
	return r.hidratar(ctx, modelos)
}

func (r *RepositorioDocente) ListarPagina(ctx context.Context, busqueda string, soloActivos bool, pagina, tamano int) ([]docente.Docente, int64, error) {
	var total int64
	if err := filtrarDocentes(db(ctx, r.DB).Model(&DocenteModel{}), busqueda, soloActivos).Count(&total).Error; err != nil {
		return nil, 0, err
	}
	var modelos []DocenteModel
	q := filtrarDocentes(db(ctx, r.DB).Preload("Categoria").Preload("Facultad").Order("apellidos, nombres, id"), busqueda, soloActivos).
		Offset((pagina - 1) * tamano).Limit(tamano)
	if err := q.Find(&modelos).Error; err != nil {
		return nil, 0, err
	}
	out, err := r.hidratar(ctx, modelos)
	return out, total, err
}

func (r *RepositorioDocente) buscar(ctx context.Context, condicion string, args ...any) (*docente.Docente, error) {
	var m DocenteModel
	if err := db(ctx, r.DB).Preload("Categoria").Preload("Facultad").First(&m, append([]any{condicion}, args...)...).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, docente.ErrNoEncontrado
		}
		return nil, err
	}
	lista := []docente.Docente{aDominioDocente(&m)}
	if err := r.cargarLegajos(ctx, lista); err != nil {
		return nil, err
	}
	return &lista[0], nil
}

func (r *RepositorioDocente) BuscarPorID(ctx context.Context, id string) (*docente.Docente, error) {
	return r.buscar(ctx, "id = ?", id)
}

func (r *RepositorioDocente) BuscarPorIdentificador(ctx context.Context, valor string) (*docente.Docente, error) {
	return r.buscar(ctx, "codigo = ? OR cedula = ?", valor, valor)
}

// errorDeEscritura traduce las violaciones de unicidad y de clave foránea de docente.
func errorDeEscritura(err error) error {
	switch {
	case esViolacionUnicidad(err) && strings.Contains(err.Error(), "codigo"):
		return docente.ErrCodigoDuplicado
	case esViolacionUnicidad(err):
		return docente.ErrCedulaDuplicada
	case esViolacionFK(err) && strings.Contains(err.Error(), "facultad"):
		return docente.ErrFacultadInvalida
	case esViolacionFK(err):
		return docente.ErrCategoriaInvalida
	}
	return err
}

func (r *RepositorioDocente) Crear(ctx context.Context, d *docente.Docente) error {
	m := DocenteModel{
		Codigo: ptr(d.Codigo), Cedula: ptr(d.Cedula), Nombres: d.Nombres, Apellidos: d.Apellidos, Correo: d.Correo,
		Telefono: d.Telefono, CategoriaID: d.CategoriaID, FacultadID: d.FacultadID, Activo: true,
	}
	if err := db(ctx, r.DB).Create(&m).Error; err != nil {
		return errorDeEscritura(err)
	}
	d.ID, d.Activo, d.CreadoEn = m.ID, true, m.CreadoEn
	return nil
}

func (r *RepositorioDocente) Actualizar(ctx context.Context, d *docente.Docente) error {
	res := db(ctx, r.DB).Model(&DocenteModel{}).Where("id = ?", d.ID).Updates(map[string]any{
		"codigo": ptr(d.Codigo), "cedula": ptr(d.Cedula),
		"nombres": d.Nombres, "apellidos": d.Apellidos, "correo": d.Correo, "telefono": d.Telefono,
		"categoria_id": d.CategoriaID, "facultad_id": d.FacultadID, "actualizado_en": gorm.Expr("now()"),
	})
	if res.Error != nil {
		return errorDeEscritura(res.Error)
	}
	if res.RowsAffected == 0 {
		return docente.ErrNoEncontrado
	}
	return nil
}

func (r *RepositorioDocente) ListarFacultades(ctx context.Context) ([]docente.Facultad, error) {
	var modelos []FacultadModel
	if err := db(ctx, r.DB).Order("nombre").Find(&modelos).Error; err != nil {
		return nil, err
	}
	out := make([]docente.Facultad, len(modelos))
	for i, m := range modelos {
		out[i] = docente.Facultad{ID: m.ID, Nombre: m.Nombre}
	}
	return out, nil
}

func (r *RepositorioDocente) CambiarActivo(ctx context.Context, id string, activo bool) error {
	res := db(ctx, r.DB).Model(&DocenteModel{}).Where("id = ?", id).
		Updates(map[string]any{"activo": activo, "actualizado_en": gorm.Expr("now()")})
	if res.Error != nil {
		return res.Error
	}
	if res.RowsAffected == 0 {
		return docente.ErrNoEncontrado
	}
	return nil
}
