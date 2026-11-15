# UT-CU12-BE-01 — Casos de uso ConsultarPublico y ConsultarPublicoPorDocente

| Campo | Valor |
|-------|-------|
| **Código** | `UT-CU12-BE-01` |
| **Nombre** | Casos de uso ConsultarPublico y ConsultarPublicoPorDocente |
| **Caso de uso** | CU-12-BE |
| **Proyecto** | `backend` |
| **Depende de** | — |
| **Estado** | Hecho (documentada de forma retroactiva) |

## Antes de implementar

[ARQUITECTURA_BACKEND.md](../../../.claude/rules/ARQUITECTURA_BACKEND.md)

## Solicitud de negocio

El docente sigue su trámite sin ir a ventanilla.

## Objetivo

Vistas públicas mínimas sin datos personales.

## Alcance

- `VistaPublica`, `ResumenPublico`.
- Lista vacía ante cédula incorrecta.

Archivos esperados (crear/tocar):

- `backend/internal/application/tramites/tramites.go`
- `backend/internal/application/tramites/tramites_test.go`

## Fuera de alcance (no tocar)

- Radicación pública.
- Cualquier archivo de `backend` que no figure en «Archivos esperados», y las otras capas del CU (tienen sus propias UT).

## Resultado esperado

Vistas públicas mínimas sin datos personales.

## Validación

1. `go test ./internal/application/tramites -count=1` en verde.
2. `make arch` pasa (domain y application sin dependencias externas).
