# UT-CU01-BE-02 — Puertos de seguridad y sus implementaciones

| Campo | Valor |
|-------|-------|
| **Código** | `UT-CU01-BE-02` |
| **Nombre** | Puertos de seguridad y sus implementaciones |
| **Caso de uso** | CU-01-BE |
| **Proyecto** | `backend` |
| **Depende de** | UT-CU01-BE-01 |
| **Estado** | Hecho (documentada de forma retroactiva) |

## Antes de implementar

[ARQUITECTURA_BACKEND.md](../../../.claude/rules/ARQUITECTURA_BACKEND.md)

## Solicitud de negocio

Hashing y firma de tokens deben poder cambiarse y probarse sin tocar los casos de uso.

## Objetivo

Puertos `PasswordHasher`, `TokenService`, `SessionCache` con bcrypt, JWT y Redis.

## Alcance

- Puertos en `application/ports`.
- bcrypt con costo configurable; JWT HS256 con `jti`, emisor y vencimiento `SESION_HORAS`.
- Caché de sesión en Redis con caída a Postgres.

Archivos esperados (crear/tocar):

- `backend/internal/application/ports/ports.go`
- `backend/internal/infrastructure/security/bcrypt.go`
- `backend/internal/infrastructure/security/jwt.go`
- `backend/internal/infrastructure/cache/sesion.go`

## Fuera de alcance (no tocar)

- El caso de uso de login (UT 03).
- Cualquier archivo de `backend` que no figure en «Archivos esperados», y las otras capas del CU (tienen sus propias UT).

## Resultado esperado

Puertos `PasswordHasher`, `TokenService`, `SessionCache` con bcrypt, JWT y Redis.

## Validación

1. `go test ./internal/infrastructure/security ./internal/infrastructure/cache -count=1` en verde.
2. `make arch` pasa (domain y application sin dependencias externas).
