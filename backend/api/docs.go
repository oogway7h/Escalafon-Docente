// Package api embebe la especificación OpenAPI (sección 7) para servirla desde /api/docs.
package api

import _ "embed"

//go:embed openapi.yaml
var OpenAPI []byte
