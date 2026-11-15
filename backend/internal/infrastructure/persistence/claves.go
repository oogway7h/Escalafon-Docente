package persistence

import "strings"

var sinTildes = strings.NewReplacer(
	"á", "a", "é", "e", "í", "i", "ó", "o", "ú", "u", "ü", "u",
	"Á", "A", "É", "E", "Í", "I", "Ó", "O", "Ú", "U", "Ü", "U",
)

// Normalizar colapsa espacios y pasa a mayúsculas sin tildes (conserva la Ñ).
func Normalizar(s string) string {
	return strings.ToUpper(sinTildes.Replace(strings.Join(strings.Fields(s), " ")))
}

// ClaveCategoria es la clave con la que se reconoce una categoría aunque venga escrita distinto
// ("PROFESOR EMERITO", "Emérito" y "EMERITO" son la misma).
func ClaveCategoria(s string) string {
	return strings.TrimSpace(strings.TrimPrefix(Normalizar(s), "PROFESOR "))
}
