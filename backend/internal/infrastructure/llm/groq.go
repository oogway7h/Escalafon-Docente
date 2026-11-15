// Package llm implementa ports.InterpreteReportes sobre la API de Groq (compatible con OpenAI).
package llm

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"escalafon/backend/internal/application/ports"
	"escalafon/backend/internal/domain/indicador"
)

const (
	urlAPI        = "https://api.groq.com/openai/v1/chat/completions"
	modeloDefecto = "openai/gpt-oss-120b"
)

type Groq struct {
	clave, modelo string
	http          *http.Client
}

func NewGroq(clave, modelo string) *Groq {
	if modelo == "" {
		modelo = modeloDefecto
	}
	return &Groq{clave: clave, modelo: modelo, http: &http.Client{Timeout: 20 * time.Second}}
}

// Herramienta obligatoria: el modelo solo puede responder con estos campos (salida estructurada).
var parametros = map[string]any{
	"type": "object",
	"properties": map[string]any{
		"tipo":         map[string]any{"type": "string", "enum": []string{"TRAMITES", "INVENTARIO", "TIEMPOS"}, "description": "TRAMITES: trámites radicados y su estado. INVENTARIO: legajos físicos y préstamos. TIEMPOS: tiempos de resolución."},
		"formato":      map[string]any{"type": "string", "enum": []string{"PDF", "XLSX"}, "description": "PDF salvo que pida Excel/planilla/hoja de cálculo."},
		"desde":        map[string]any{"type": "string", "description": "AAAA-MM-DD inclusive; vacío si no hay."},
		"hasta":        map[string]any{"type": "string", "description": "AAAA-MM-DD inclusive; vacío si no hay."},
		"estado":       map[string]any{"type": "string", "enum": []string{"", "INGRESADO", "ASIGNADO", "EN_REVISION", "OBSERVADO", "DICTAMINADO", "EN_FACULTAD", "REVISION_FINAL", "REMITIDO", "ARCHIVADO"}, "description": "Solo para TRAMITES. EN_FACULTAD = Fuera de Escalafón."},
		"tipoTramite":  map[string]any{"type": "string", "enum": []string{"", "CATEGORIZACION", "REGULARIZACION"}},
		"categoriaId":  map[string]any{"type": "integer", "description": "Id de la lista de categorías; 0 si no menciona ninguna."},
		"estadoLegajo": map[string]any{"type": "string", "enum": []string{"", "EN_ARCHIVO", "PRESTADO"}, "description": "Solo para INVENTARIO."},
		"explicacion":  map[string]any{"type": "string", "description": "Una frase en español que resume qué reporte se armó."},
	},
	"required": []string{"tipo", "formato", "explicacion"},
}

func (g *Groq) Interpretar(ctx context.Context, frase string, hoy time.Time, categorias []ports.CategoriaRef) (*ports.PeticionReporte, error) {
	var lista strings.Builder
	for _, c := range categorias {
		fmt.Fprintf(&lista, "- %d: %s\n", c.ID, c.Nombre)
	}
	sistema := fmt.Sprintf("Eres el asistente de reportes del Departamento de Escalafón Docente. Traduce la petición a los filtros de un reporte "+
		"llamando a definir_reporte. Hoy es %s (hora de Bolivia); resuelve fechas relativas (\"este mes\", \"octubre\", \"último trimestre\") a partir de hoy. "+
		"Si la petición no pide un reporte, elige TRAMITES sin filtros y explícalo. No inventes filtros que no se mencionen.\nCategorías:\n%s",
		hoy.Format("2006-01-02"), lista.String())

	cuerpo, _ := json.Marshal(map[string]any{
		"model": g.modelo, "temperature": 0, "max_tokens": 500,
		"messages": []any{
			map[string]any{"role": "system", "content": sistema},
			map[string]any{"role": "user", "content": frase},
		},
		"tools": []any{map[string]any{"type": "function", "function": map[string]any{
			"name": "definir_reporte", "description": "Define los filtros del reporte que pide el usuario.", "parameters": parametros,
		}}},
		"tool_choice": map[string]any{"type": "function", "function": map[string]any{"name": "definir_reporte"}},
	})
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, urlAPI, bytes.NewReader(cuerpo))
	if err != nil {
		return nil, err
	}
	req.Header.Set("content-type", "application/json")
	req.Header.Set("authorization", "Bearer "+g.clave)
	resp, err := g.http.Do(req)
	if err != nil {
		return nil, indicador.ErrInterpreteNoDisponible
	}
	defer resp.Body.Close()
	datos, _ := io.ReadAll(io.LimitReader(resp.Body, 1<<20))
	if resp.StatusCode != http.StatusOK {
		return nil, indicador.ErrInterpreteNoDisponible
	}

	var sal struct {
		Choices []struct {
			Message struct {
				ToolCalls []struct {
					Function struct {
						Arguments string `json:"arguments"`
					} `json:"function"`
				} `json:"tool_calls"`
			} `json:"message"`
		} `json:"choices"`
	}
	if err := json.Unmarshal(datos, &sal); err != nil || len(sal.Choices) == 0 || len(sal.Choices[0].Message.ToolCalls) == 0 {
		return nil, indicador.ErrNoEntendido
	}
	var i struct {
		Tipo         string `json:"tipo"`
		Formato      string `json:"formato"`
		Desde        string `json:"desde"`
		Hasta        string `json:"hasta"`
		Estado       string `json:"estado"`
		TipoTramite  string `json:"tipoTramite"`
		CategoriaID  int16  `json:"categoriaId"`
		EstadoLegajo string `json:"estadoLegajo"`
		Explicacion  string `json:"explicacion"`
	}
	if err := json.Unmarshal([]byte(sal.Choices[0].Message.ToolCalls[0].Function.Arguments), &i); err != nil {
		return nil, indicador.ErrNoEntendido
	}
	return &ports.PeticionReporte{
		Tipo: i.Tipo, Formato: i.Formato, Desde: i.Desde, Hasta: i.Hasta, Estado: i.Estado,
		TipoTramite: i.TipoTramite, CategoriaID: i.CategoriaID, EstadoLegajo: i.EstadoLegajo, Explicacion: i.Explicacion,
	}, nil
}
