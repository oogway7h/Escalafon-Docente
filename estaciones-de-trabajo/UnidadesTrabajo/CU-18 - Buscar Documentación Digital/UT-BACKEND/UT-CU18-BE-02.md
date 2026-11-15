# UT-CU18-BE-02 — Rutas de búsqueda

| Campo | Valor |
|-------|-------|
| **Código** | `UT-CU18-BE-02` |
| **Nombre** | Rutas de búsqueda |
| **Caso de uso** | CU-18-BE |
| **Proyecto** | `backend` |
| **Depende de** | UT-CU18-BE-01 |
| **Estado** | Hecho (documentada de forma retroactiva) |

## Antes de implementar

[ARQUITECTURA_BACKEND.md](../../../.claude/rules/ARQUITECTURA_BACKEND.md)

## Solicitud de negocio

Exponer la búsqueda.

## Objetivo

Rutas de búsqueda global, por expediente y por documento.

## Alcance

- `ListaResultadosDTO`.

Archivos esperados (crear/tocar):

- `backend/internal/interfaces/http/handlers/documentos.go`
- `backend/internal/interfaces/http/handlers/expediente.go`
- `backend/internal/interfaces/http/router.go`

## Fuera de alcance (no tocar)

- Cualquier archivo de `backend` que no figure en «Archivos esperados», y las otras capas del CU (tienen sus propias UT).

## Resultado esperado

Rutas de búsqueda global, por expediente y por documento.

## Validación

1. `go test ./internal/interfaces/http -count=1` en verde.
2. `make arch` pasa (domain y application sin dependencias externas).
