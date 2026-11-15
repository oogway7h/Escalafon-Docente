// Package dto define los objetos de entrada/salida HTTP y su validación (sección 7).
package dto

import (
	"fmt"

	"github.com/go-playground/validator/v10"
)

var validador = validator.New()

// Validar devuelve un mapa campo -> mensaje si v no cumple sus etiquetas `validate`, o nil si es válido.
func Validar(v any) map[string]string {
	err := validador.Struct(v)
	if err == nil {
		return nil
	}
	errsValidacion, ok := err.(validator.ValidationErrors)
	if !ok {
		return map[string]string{"_": "solicitud inválida"}
	}
	campos := map[string]string{}
	for _, fe := range errsValidacion {
		campos[fe.Field()] = mensajeParaTag(fe)
	}
	return campos
}

func mensajeParaTag(fe validator.FieldError) string {
	switch fe.Tag() {
	case "required":
		return "Este campo es obligatorio"
	case "email":
		return "Debe ser un correo electrónico válido"
	case "min":
		return fmt.Sprintf("Debe tener al menos %s caracteres", fe.Param())
	default:
		return "Valor inválido"
	}
}
