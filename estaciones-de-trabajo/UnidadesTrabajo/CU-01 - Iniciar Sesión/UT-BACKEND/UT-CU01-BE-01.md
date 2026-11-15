# UT-CU01-BE-01 — Entidad Usuario y reglas de bloqueo

| Campo | Valor |
|-------|-------|
| **Código** | `UT-CU01-BE-01` |
| **Nombre** | Entidad Usuario y reglas de bloqueo |
| **Caso de uso** | CU-01-BE |
| **Proyecto** | `backend` |
| **Depende de** | — |
| **Estado** | Hecho (documentada de forma retroactiva) |

## Antes de implementar

[ARQUITECTURA_BACKEND.md](../../../.claude/rules/ARQUITECTURA_BACKEND.md)

## Solicitud de negocio

El bloqueo por intentos es una regla de negocio y debe vivir en el dominio, no en el handler.

## Objetivo

`domain/usuario` modela el usuario, su estado y el registro de intentos fallidos con bloqueo temporal.

## Alcance

- `Usuario` con estado activo, intentos fallidos y `BloqueadoHasta`.
- `RegistrarFallo(max, minutos, ahora)` y `Bloqueado(ahora)`.
- Errores de dominio `ErrCredencialesInvalidas`, `ErrCuentaBloqueada`.

Archivos esperados (crear/tocar):

- `backend/internal/domain/usuario/usuario.go`
- `backend/internal/domain/usuario/errores.go`

## Fuera de alcance (no tocar)

- Persistencia GORM (UT 02).
- HTTP y DTOs (UT 05).
- Cualquier archivo de `backend` que no figure en «Archivos esperados», y las otras capas del CU (tienen sus propias UT).

## Resultado esperado

`domain/usuario` modela el usuario, su estado y el registro de intentos fallidos con bloqueo temporal.

## Validación

1. `go test ./internal/application/auth -count=1` en verde.
2. `make arch` pasa (domain y application sin dependencias externas).
