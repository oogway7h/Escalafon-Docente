# UT-CU15-BE-04 — Expediente digital del docente y rutas

| Campo | Valor |
|-------|-------|
| **Código** | `UT-CU15-BE-04` |
| **Nombre** | Expediente digital del docente y rutas |
| **Caso de uso** | CU-15-BE |
| **Proyecto** | `backend` |
| **Depende de** | UT-CU15-BE-03 |
| **Estado** | Hecho (documentada de forma retroactiva) |

## Antes de implementar

[ARQUITECTURA_BACKEND.md](../../../.claude/rules/ARQUITECTURA_BACKEND.md)

## Solicitud de negocio

Archivo digitaliza las carpetas físicas por docente.

## Objetivo

`PrepararSubidaDocente`, `ListarExpediente`, `UrlExpediente` y rutas de documentos.

## Alcance

- Tipos CURRICULUM … EMERITO.

Archivos esperados (crear/tocar):

- `backend/internal/application/tramites/expediente.go`
- `backend/internal/interfaces/http/handlers/documentos.go`
- `backend/internal/interfaces/http/handlers/expediente.go`
- `backend/internal/interfaces/http/router.go`

## Fuera de alcance (no tocar)

- Búsqueda.
- Cualquier archivo de `backend` que no figure en «Archivos esperados», y las otras capas del CU (tienen sus propias UT).

## Resultado esperado

`PrepararSubidaDocente`, `ListarExpediente`, `UrlExpediente` y rutas de documentos.

## Validación

1. `go test ./internal/application/tramites -count=1` en verde.
2. `go test ./internal/interfaces/http -count=1` en verde.
3. `make arch` pasa (domain y application sin dependencias externas).
