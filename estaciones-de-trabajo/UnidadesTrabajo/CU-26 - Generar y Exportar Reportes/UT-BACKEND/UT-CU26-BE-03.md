# UT-CU26-BE-03 — Ruta de reportes

| Campo | Valor |
|-------|-------|
| **Código** | `UT-CU26-BE-03` |
| **Nombre** | Ruta de reportes |
| **Caso de uso** | CU-26-BE |
| **Proyecto** | `backend` |
| **Depende de** | UT-CU26-BE-02 |
| **Estado** | Hecho (documentada de forma retroactiva) |

## Antes de implementar

[ARQUITECTURA_BACKEND.md](../../../.claude/rules/ARQUITECTURA_BACKEND.md)

## Solicitud de negocio

Descargar el archivo.

## Objetivo

`POST /reportes`.

## Alcance

- Cabeceras de descarga.

Archivos esperados (crear/tocar):

- `backend/internal/interfaces/http/handlers/indicadores.go`
- `backend/internal/interfaces/http/router.go`

## Fuera de alcance (no tocar)

- Cualquier archivo de `backend` que no figure en «Archivos esperados», y las otras capas del CU (tienen sus propias UT).

## Resultado esperado

`POST /reportes`.

## Validación

1. `go test ./internal/interfaces/http -count=1` en verde.
2. `make arch` pasa (domain y application sin dependencias externas).
