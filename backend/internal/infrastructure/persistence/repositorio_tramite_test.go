package persistence

import (
	"context"
	"errors"
	"strings"
	"testing"
	"time"

	"github.com/google/uuid"

	"escalafon/backend/internal/domain/tramite"
)

// Comprueba contra Postgres la baja lógica de un adjunto (CU-16): el objeto se conserva con su
// motivo, el requisito de la categoría queda descubierto y la baja no se puede repetir.
func TestRepositorioTramite_DarDeBajaDocumento(t *testing.T) {
	dbase := conectarPrueba(t)
	conTx(t, dbase, func(ctx context.Context) {
		rolm := RolModel{Nombre: "rol-documento-prueba"}
		if err := db(ctx, dbase).Create(&rolm).Error; err != nil {
			t.Fatal(err)
		}
		usr := UsuarioModel{RolID: rolm.ID, Nombres: "Ana", Apellidos: "Analista", Correo: "documento-prueba@ded.dev", HashContrasena: "x"}
		if err := db(ctx, dbase).Create(&usr).Error; err != nil {
			t.Fatal(err)
		}
		codigo := "55501"
		docente := DocenteModel{Codigo: &codigo, Nombres: "Doc", Apellidos: "Prueba"}
		if err := db(ctx, dbase).Create(&docente).Error; err != nil {
			t.Fatal(err)
		}
		cat := CategoriaModel{Nombre: "Categoría de prueba CU-16"}
		if err := db(ctx, dbase).Create(&cat).Error; err != nil {
			t.Fatal(err)
		}
		req := CategoriaRequisitoModel{CategoriaID: cat.ID, Nombre: "Título"}
		if err := db(ctx, dbase).Create(&req).Error; err != nil {
			t.Fatal(err)
		}

		repo := NewRepositorioTramite(dbase)
		docID := uuid.NewString()
		tr := &tramite.Tramite{
			ID: uuid.NewString(), Codigo: "ED-2026-DOCBAJ", DocenteID: docente.ID, CategoriaID: cat.ID,
			Tipo: tramite.Categorizacion, Estado: tramite.Ingresado, Origen: tramite.OrigenPortal,
			Documentos: []tramite.Documento{{
				ID: docID, RequisitoID: &req.ID, NombreOriginal: "titulo.pdf", TipoMime: "application/pdf",
				Tamano: 10, HashSHA256: strings.Repeat("a", 64), Bucket: "bucket-prueba",
				ClaveObjeto: "tramites/prueba/titulo.pdf", Activo: true,
			}},
			Historial: []tramite.Historial{{Estado: tramite.Ingresado, Detalle: "Radicado", Fecha: time.Now()}},
		}
		if err := repo.Crear(ctx, tr); err != nil {
			t.Fatal(err)
		}

		antes, err := repo.BuscarPorID(ctx, tr.ID)
		if err != nil || !antes.Documentos[0].Activo || !antes.Requisitos[0].Presente {
			t.Fatalf("el adjunto nace activo y cubre el requisito: %v %+v", err, antes.Documentos)
		}

		ahora := time.Now()
		if err := repo.DarDeBajaDocumento(ctx, tr.ID, docID, "Cargado por error", usr.ID, ahora); err != nil {
			t.Fatal(err)
		}
		despues, err := repo.BuscarPorID(ctx, tr.ID)
		if err != nil {
			t.Fatal(err)
		}
		d := despues.Documentos[0]
		if d.Activo || d.MotivoBaja != "Cargado por error" || d.BajaEn == nil || d.ClaveObjeto == "" {
			t.Fatalf("la baja es lógica y conserva el objeto con su motivo: %+v", d)
		}
		if d.BajaPor == nil || *d.BajaPor != usr.ID {
			t.Fatalf("la baja debe quedar atribuida a quien la hizo: %+v", d.BajaPor)
		}
		if despues.Requisitos[0].Presente {
			t.Fatal("un documento dado de baja deja de cubrir el requisito de la categoría")
		}

		if err := repo.DarDeBajaDocumento(ctx, tr.ID, docID, "Otra vez", usr.ID, ahora); !errors.Is(err, tramite.ErrDocumentoDeBaja) {
			t.Fatalf("esperaba ErrDocumentoDeBaja al repetir la baja, obtuve %v", err)
		}
	})
}

// CU-17 y CU-18 contra Postgres: la búsqueda global halla por contenido y por docente, ignora los
// documentos de baja, y la bandeja de bajas y la restitución funcionan sobre el mismo documento.
func TestRepositorioTramite_BusquedaGlobalYRestitucion(t *testing.T) {
	dbase := conectarPrueba(t)
	conTx(t, dbase, func(ctx context.Context) {
		rolm := RolModel{Nombre: "rol-restituir-prueba"}
		if err := db(ctx, dbase).Create(&rolm).Error; err != nil {
			t.Fatal(err)
		}
		usr := UsuarioModel{RolID: rolm.ID, Nombres: "Luis", Apellidos: "Vaca", Correo: "restituir-prueba@ded.dev", HashContrasena: "x"}
		if err := db(ctx, dbase).Create(&usr).Error; err != nil {
			t.Fatal(err)
		}
		fac := FacultadModel{ID: 9901, Nombre: "Facultad de prueba CU-18"}
		if err := db(ctx, dbase).Create(&fac).Error; err != nil {
			t.Fatal(err)
		}
		codigo, cedula := "55502", "9876543"
		docente := DocenteModel{Codigo: &codigo, Cedula: &cedula, Nombres: "Rosa", Apellidos: "Terrazas Prueba", FacultadID: &fac.ID}
		if err := db(ctx, dbase).Create(&docente).Error; err != nil {
			t.Fatal(err)
		}
		repo := NewRepositorioTramite(dbase)
		docID := uuid.NewString()
		d := tramite.Documento{ID: docID, TipoDocumento: tramite.TipoTitularA, NombreOriginal: "resolucion.pdf", TipoMime: "application/pdf",
			Tamano: 10, Bucket: "b", ClaveObjeto: "docentes/55502/resolucion.pdf", Activo: true, EstadoOCR: tramite.OCRPendiente}
		if err := repo.CrearDocumentoDocente(ctx, docente.ID, &d); err != nil {
			t.Fatal(err)
		}
		if err := repo.GuardarPaginas(ctx, "", docID, []string{"Portada", "Se aprueba la resolución de ascenso a la categoría Titular"}); err != nil {
			t.Fatal(err)
		}

		buscar := func(f tramite.FiltroBusqueda) []tramite.ResultadoBusqueda {
			t.Helper()
			r, _, err := repo.BuscarDocumentos(ctx, f, 1, 10)
			if err != nil {
				t.Fatal(err)
			}
			return r
		}
		if r := buscar(tramite.FiltroBusqueda{Consulta: "resolucion ascenso"}); len(r) != 1 || len(r[0].Paginas) != 1 || r[0].Paginas[0] != 2 || !strings.Contains(r[0].Fragmento, "⟦") || r[0].DocenteCodigo != codigo {
			t.Fatalf("por contenido, sin tildes, en la página 2 y con resaltado: %+v", r)
		}
		if r := buscar(tramite.FiltroBusqueda{Consulta: "terrazas", Campo: tramite.BuscarDocente}); len(r) != 1 || len(r[0].Paginas) != 0 {
			t.Fatalf("por nombre del docente: %+v", r)
		}
		if r := buscar(tramite.FiltroBusqueda{Consulta: cedula}); len(r) != 1 {
			t.Fatalf("por cédula exacta: %+v", r)
		}
		if r := buscar(tramite.FiltroBusqueda{Consulta: "terrazas", Campo: tramite.BuscarContenido}); len(r) != 0 {
			t.Fatalf("el nombre no está en el contenido: %+v", r)
		}
		if r := buscar(tramite.FiltroBusqueda{Consulta: "titular", FacultadID: fac.ID, Recientes: true}); len(r) != 1 {
			t.Fatalf("filtro por la facultad del docente y orden por recientes: %+v", r)
		}
		if r := buscar(tramite.FiltroBusqueda{Consulta: "titular", FacultadID: fac.ID + 1}); len(r) != 0 {
			t.Fatalf("otra facultad excluye: %+v", r)
		}
		if r := buscar(tramite.FiltroBusqueda{Consulta: "titular", Tipo: tramite.TipoCurriculum}); len(r) != 0 {
			t.Fatalf("el filtro de tipo excluye: %+v", r)
		}

		if err := repo.DarDeBajaDocumento(ctx, "", docID, "Duplicado", usr.ID, time.Now()); err != nil {
			t.Fatal(err)
		}
		if r := buscar(tramite.FiltroBusqueda{Consulta: "titular"}); len(r) != 0 {
			t.Fatalf("un documento de baja no se busca: %+v", r)
		}
		bajas, total, err := repo.ListarDocumentosDeBaja(ctx, "terrazas", 1, 10)
		if err != nil || total != 1 || bajas[0].ID != docID || bajas[0].BajaPorNombre != "Luis Vaca" || bajas[0].MotivoBaja != "Duplicado" || bajas[0].DocenteID != docente.ID {
			t.Fatalf("bandeja de bajas: %v %d %+v", err, total, bajas)
		}
		if err := repo.RestituirDocumento(ctx, docID); err != nil {
			t.Fatal(err)
		}
		if err := repo.RestituirDocumento(ctx, docID); !errors.Is(err, tramite.ErrDocumentoActivo) {
			t.Fatalf("esperaba ErrDocumentoActivo, obtuve %v", err)
		}
		if err := repo.RestituirDocumento(ctx, uuid.NewString()); !errors.Is(err, tramite.ErrDocumentoNoEncontrado) {
			t.Fatalf("esperaba ErrDocumentoNoEncontrado, obtuve %v", err)
		}
		if r := buscar(tramite.FiltroBusqueda{Consulta: "titular"}); len(r) != 1 {
			t.Fatalf("restituido vuelve a la búsqueda: %+v", r)
		}
	})
}
