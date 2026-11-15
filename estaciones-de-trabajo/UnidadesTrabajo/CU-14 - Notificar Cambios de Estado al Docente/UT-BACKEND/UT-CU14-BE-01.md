# UT-CU14-BE-01 — Dominio y tabla de envíos

| Campo | Valor |
|-------|-------|
| **Código** | `UT-CU14-BE-01` |
| **Nombre** | Dominio y tabla de envíos |
| **Caso de uso** | CU-14-BE |
| **Proyecto** | `backend` |
| **Depende de** | — |
| **Estado** | Hecho (documentada de forma retroactiva) |

## Antes de implementar

[ARQUITECTURA_BACKEND.md](../../../.claude/rules/ARQUITECTURA_BACKEND.md)

## Solicitud de negocio

Cada aviso debe quedar trazable.

## Objetivo

`domain/notificacion` y tabla `notificacion_envio`.

## Alcance

- Estados PENDIENTE/ENVIADO/FALLIDO, intentos y error.

Archivos esperados (crear/tocar):

- `backend/internal/domain/notificacion/notificacion.go`
- `backend/internal/infrastructure/persistence/repositorio_notificacion.go`
- `backend/migrations/20261018100000_notificacion_envio.sql`

## Fuera de alcance (no tocar)

- Envío (UT 02).
- Cualquier archivo de `backend` que no figure en «Archivos esperados», y las otras capas del CU (tienen sus propias UT).

## Resultado esperado

`domain/notificacion` y tabla `notificacion_envio`.

## Validación

1. `go test ./internal/application/notificaciones -count=1` en verde.
2. `make arch` pasa (domain y application sin dependencias externas).
