# UT-CU08-BE-01 — Agregado Trámite y máquina de estados

| Campo | Valor |
|-------|-------|
| **Código** | `UT-CU08-BE-01` |
| **Nombre** | Agregado Trámite y máquina de estados |
| **Caso de uso** | CU-08-BE |
| **Proyecto** | `backend` |
| **Depende de** | — |
| **Estado** | Hecho (documentada de forma retroactiva) |

## Antes de implementar

[ARQUITECTURA_BACKEND.md](../../../.claude/rules/ARQUITECTURA_BACKEND.md)

## Solicitud de negocio

Todas las transiciones del flujo real dependen de este modelo.

## Objetivo

`domain/tramite` con estados, tipos, historial y transiciones permitidas.

## Alcance

- Estados INGRESADO … ARCHIVADO.
- `Transicion` valida origen y destino.
- Generador de código.

Archivos esperados (crear/tocar):

- `backend/internal/domain/tramite/tramite.go`
- `backend/migrations/20261005021952_tramites.sql`
- `backend/migrations/20261020140000_flujo_escalafon.sql`

## Fuera de alcance (no tocar)

- Casos de uso.
- Cualquier archivo de `backend` que no figure en «Archivos esperados», y las otras capas del CU (tienen sus propias UT).

## Resultado esperado

`domain/tramite` con estados, tipos, historial y transiciones permitidas.

## Validación

1. `go test ./internal/application/tramites -run Radicar -count=1` en verde.
2. `make arch` pasa (domain y application sin dependencias externas).
