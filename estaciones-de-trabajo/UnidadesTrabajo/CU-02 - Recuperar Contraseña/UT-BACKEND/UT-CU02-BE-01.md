# UT-CU02-BE-01 — Dominio de token y generador seguro

| Campo | Valor |
|-------|-------|
| **Código** | `UT-CU02-BE-01` |
| **Nombre** | Dominio de token y generador seguro |
| **Caso de uso** | CU-02-BE |
| **Proyecto** | `backend` |
| **Depende de** | — |
| **Estado** | Hecho (documentada de forma retroactiva) |

## Antes de implementar

[ARQUITECTURA_BACKEND.md](../../../.claude/rules/ARQUITECTURA_BACKEND.md)

## Solicitud de negocio

Un enlace adivinable comprometería cualquier cuenta.

## Objetivo

`domain/token` + `GeneradorToken` con `crypto/rand` y hash SHA-256.

## Alcance

- Token con propósito (RECUPERACION, ACTIVACION), vencimiento y uso.
- Repositorio que guarda solo el hash.
- Índice manual de token pendiente (`ajustes_manuales.sql`).

Archivos esperados (crear/tocar):

- `backend/internal/domain/token/token.go`
- `backend/internal/infrastructure/security/tokens.go`
- `backend/internal/infrastructure/persistence/repositorio_token.go`

## Fuera de alcance (no tocar)

- Envío de correo (UT 02).
- Cualquier archivo de `backend` que no figure en «Archivos esperados», y las otras capas del CU (tienen sus propias UT).

## Resultado esperado

`domain/token` + `GeneradorToken` con `crypto/rand` y hash SHA-256.

## Validación

1. `go test ./internal/infrastructure/security -count=1` en verde.
2. `make arch` pasa (domain y application sin dependencias externas).
