# UT-CU12-BE-02 — Rutas públicas con limitador

| Campo | Valor |
|-------|-------|
| **Código** | `UT-CU12-BE-02` |
| **Nombre** | Rutas públicas con limitador |
| **Caso de uso** | CU-12-BE |
| **Proyecto** | `backend` |
| **Depende de** | UT-CU12-BE-01 |
| **Estado** | Hecho (documentada de forma retroactiva) |

## Antes de implementar

[ARQUITECTURA_BACKEND.md](../../../.claude/rules/ARQUITECTURA_BACKEND.md)

## Solicitud de negocio

Evitar enumeración de trámites.

## Objetivo

Grupo `/publico` con 30 peticiones/min por IP.

## Alcance

- `VistaPublicaDTO`, `ConsultaDocenteEntrada`.

Archivos esperados (crear/tocar):

- `backend/internal/interfaces/http/handlers/tramites.go`
- `backend/internal/interfaces/http/dto/tramites.go`
- `backend/internal/interfaces/http/router.go`
- `backend/cmd/api/main.go`

## Fuera de alcance (no tocar)

- Cualquier archivo de `backend` que no figure en «Archivos esperados», y las otras capas del CU (tienen sus propias UT).

## Resultado esperado

Grupo `/publico` con 30 peticiones/min por IP.

## Validación

1. `go test ./internal/interfaces/http -count=1` en verde.
2. `make arch` pasa (domain y application sin dependencias externas).
