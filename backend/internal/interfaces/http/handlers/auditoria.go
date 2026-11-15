package handlers

import (
	"strconv"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"

	"escalafon/backend/internal/application/auditorias"
	"escalafon/backend/internal/domain/auditoria"
	"escalafon/backend/internal/interfaces/http/apierror"
	"escalafon/backend/internal/interfaces/http/dto"
)

// zonaBolivia interpreta las fechas de los filtros como días de Bolivia (UTC-4, sin horario de verano).
var zonaBolivia = time.FixedZone("BOT", -4*60*60)

// AuditoriaHandler expone /api/auditoria (CU-06). Solo tiene lecturas: la bitácora no se edita.
type AuditoriaHandler struct {
	CasoConsultar *auditorias.Consultar
	CasoOpciones  *auditorias.ListarOpciones
}

func aAutorDTO(a auditoria.Autor) dto.AutorAuditoriaDTO {
	return dto.AutorAuditoriaDTO{ID: a.ID, Nombres: a.Nombres, Apellidos: a.Apellidos, Correo: a.Correo}
}

// leerFiltroAuditoria valida los parámetros de consulta: usuarioId (UUID), accion, desde y hasta
// (AAAA-MM-DD, ambos inclusive), pagina y tamano (enteros).
func leerFiltroAuditoria(c *fiber.Ctx) (auditoria.Filtro, error) {
	var f auditoria.Filtro
	campos := map[string]string{}

	if v := strings.TrimSpace(c.Query("usuarioId")); v != "" {
		id, err := uuid.Parse(v)
		if err != nil {
			campos["usuarioId"] = "Debe ser un identificador de usuario válido"
		} else {
			f.UsuarioID = id.String()
		}
	}
	f.Accion = strings.ToUpper(strings.TrimSpace(c.Query("accion")))

	fecha := func(nombre string) *time.Time {
		v := strings.TrimSpace(c.Query(nombre))
		if v == "" {
			return nil
		}
		t, err := time.ParseInLocation(formatoFecha, v, zonaBolivia)
		if err != nil {
			campos[nombre] = "Debe ser una fecha AAAA-MM-DD"
			return nil
		}
		return &t
	}
	f.Desde = fecha("desde")
	if hasta := fecha("hasta"); hasta != nil {
		fin := hasta.AddDate(0, 0, 1) // el día "hasta" completo
		f.Hasta = &fin
	}

	entero := func(nombre string) int {
		v := strings.TrimSpace(c.Query(nombre))
		if v == "" {
			return 0
		}
		n, err := strconv.Atoi(v)
		if err != nil {
			campos[nombre] = "Debe ser un número entero"
		}
		return n
	}
	f.Pagina, f.Tamano = entero("pagina"), entero("tamano")

	if len(campos) > 0 {
		return f, apierror.Validacion(campos)
	}
	return f, nil
}

func (h *AuditoriaHandler) Listar(c *fiber.Ctx) error {
	f, err := leerFiltroAuditoria(c)
	if err != nil {
		return err
	}
	res, err := h.CasoConsultar.Ejecutar(c.Context(), f)
	if err != nil {
		return apierror.Mapear(err)
	}
	out := dto.ListaAuditoriaDTO{Datos: make([]dto.RegistroAuditoriaDTO, len(res.Registros)), Total: res.Total, Pagina: res.Pagina, Tamano: res.Tamano}
	for i, r := range res.Registros {
		reg := dto.RegistroAuditoriaDTO{
			ID: r.ID, Fecha: r.Fecha, Accion: r.Accion, Entidad: r.Entidad, EntidadID: r.EntidadID,
			Detalle: r.Detalle, DireccionIP: r.DireccionIP,
		}
		if r.Autor != nil {
			autor := aAutorDTO(*r.Autor)
			reg.Usuario = &autor
		}
		out.Datos[i] = reg
	}
	return c.JSON(out)
}

func (h *AuditoriaHandler) Filtros(c *fiber.Ctx) error {
	op, err := h.CasoOpciones.Ejecutar(c.Context())
	if err != nil {
		return apierror.Mapear(err)
	}
	out := dto.FiltrosAuditoriaDTO{Acciones: op.Acciones, Usuarios: make([]dto.AutorAuditoriaDTO, len(op.Usuarios))}
	for i, a := range op.Usuarios {
		out.Usuarios[i] = aAutorDTO(a)
	}
	return c.JSON(out)
}
