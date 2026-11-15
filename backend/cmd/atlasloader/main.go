package main

import (
	"fmt"
	"io"
	"os"

	"ariga.io/atlas-provider-gorm/gormschema"

	"escalafon/backend/internal/infrastructure/persistence"
)

func main() {
	stmts, err := gormschema.New("postgres").Load(
		&persistence.RolModel{},
		&persistence.PermisoModel{},
		&persistence.RolPermisoModel{},
		&persistence.UsuarioModel{},
		&persistence.TokenRecuperacionModel{},
		&persistence.AuditoriaModel{},
		&persistence.CategoriaModel{},
		&persistence.CategoriaRequisitoModel{},
		&persistence.DocenteModel{},
		&persistence.TramiteModel{},
		&persistence.TramiteHistorialModel{},
		&persistence.TramiteDocumentoModel{},
		&persistence.ObservacionModel{},
		&persistence.FacultadModel{},
		&persistence.LegajoModel{},
		&persistence.SolicitudCarpetaModel{},
		&persistence.MovimientoCustodiaModel{},
		&persistence.NotificacionEnvioModel{},
	)
	if err != nil {
		fmt.Fprintln(os.Stderr, "falló la carga del esquema gorm:", err)
		os.Exit(1)
	}
	io.WriteString(os.Stdout, stmts)
}
