// Package iaservice implementa ports.ProcesadorDocumentos sobre el servicio FastAPI de OCR y compresión.
package iaservice

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"escalafon/backend/internal/application/ports"
	"escalafon/backend/internal/domain/tramite"
)

type Cliente struct {
	url, token string
	http       *http.Client
}

func New(url, token string) *Cliente {
	return &Cliente{url: strings.TrimRight(url, "/"), token: token, http: &http.Client{Timeout: 30 * time.Second}}
}

func (c *Cliente) hacer(ctx context.Context, metodo, ruta string, cuerpo, destino any) error {
	var r *bytes.Reader
	if cuerpo != nil {
		b, _ := json.Marshal(cuerpo)
		r = bytes.NewReader(b)
	} else {
		r = bytes.NewReader(nil)
	}
	req, err := http.NewRequestWithContext(ctx, metodo, c.url+ruta, r)
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-Token", c.token)
	resp, err := c.http.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode >= 300 {
		return fmt.Errorf("ia-service respondió %d", resp.StatusCode)
	}
	return json.NewDecoder(resp.Body).Decode(destino)
}

func (c *Cliente) Iniciar(ctx context.Context, s ports.SolicitudProceso) (string, error) {
	var out struct {
		ID string `json:"id"`
	}
	err := c.hacer(ctx, http.MethodPost, "/trabajos", map[string]string{
		"bucket": s.Bucket, "clave": s.Clave, "tipoMime": s.TipoMime, "claveSalida": s.ClaveSalida,
	}, &out)
	return out.ID, err
}

func (c *Cliente) Estado(ctx context.Context, id string) (*ports.ResultadoProceso, error) {
	var out struct {
		Estado          string `json:"estado"`
		Texto           string   `json:"texto"`
		Paginas         []string `json:"paginas"`
		Hash            string `json:"hashSha256"`
		TamanoOriginal  int64  `json:"tamanoOriginal"`
		TamanoFinal     int64  `json:"tamanoFinal"`
		ClaveComprimida string `json:"claveComprimida"`
		Error           string `json:"error"`
	}
	if err := c.hacer(ctx, http.MethodGet, "/trabajos/"+id, nil, &out); err != nil {
		return nil, err
	}
	return &ports.ResultadoProceso{
		Estado: out.Estado, Texto: out.Texto, Paginas: out.Paginas, HashSHA256: out.Hash, TamanoOriginal: out.TamanoOriginal,
		TamanoFinal: out.TamanoFinal, ClaveComprimida: out.ClaveComprimida, Error: out.Error,
	}, nil
}

// NoConfigurado se usa sin IA_SERVICE_URL: los documentos se suben pero no se procesan.
type NoConfigurado struct{}

func (NoConfigurado) Iniciar(context.Context, ports.SolicitudProceso) (string, error) {
	return "", tramite.ErrProcesadorNoDisponible
}
func (NoConfigurado) Estado(context.Context, string) (*ports.ResultadoProceso, error) {
	return nil, tramite.ErrProcesadorNoDisponible
}
