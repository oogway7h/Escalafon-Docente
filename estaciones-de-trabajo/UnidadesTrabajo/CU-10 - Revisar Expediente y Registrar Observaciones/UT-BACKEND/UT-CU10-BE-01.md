# UT-CU10-BE-01 — Casos de uso IniciarRevision, RegistrarObservacion y ResolverObservacion

| Campo | Valor |
|-------|-------|
| **Código** | `UT-CU10-BE-01` |
| **Nombre** | Casos de uso IniciarRevision, RegistrarObservacion y ResolverObservacion |
| **Caso de uso** | CU-10-BE |
| **Proyecto** | `backend` |
| **Depende de** | — |
| **Estado** | Hecho (documentada de forma retroactiva) |

## Antes de implementar

[ARQUITECTURA_BACKEND.md](../../../.claude/rules/ARQUITECTURA_BACKEND.md)

## Solicitud de negocio

El analista documenta lo que falta con plazo.

## Objetivo

Transiciones Asignado→En revisión⇄Observado con sus validaciones.

## Alcance

- Exige solicitud de carpeta ligada.
- Notifica (CU-14) solo al pasar a Observado.

Archivos esperados (crear/tocar):

- `backend/internal/application/tramites/tramites.go`
- `backend/internal/application/tramites/tramites_test.go`
- `backend/internal/domain/tramite/tramite.go`

## Fuera de alcance (no tocar)

- Dictamen.
- Cualquier archivo de `backend` que no figure en «Archivos esperados», y las otras capas del CU (tienen sus propias UT).

## Resultado esperado

Transiciones Asignado→En revisión⇄Observado con sus validaciones.

## Validación

1. `go test ./internal/application/tramites -count=1` en verde.
2. `make arch` pasa (domain y application sin dependencias externas).
