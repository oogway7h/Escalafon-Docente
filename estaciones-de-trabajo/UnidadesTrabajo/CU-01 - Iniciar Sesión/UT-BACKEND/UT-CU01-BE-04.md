# UT-CU01-BE-04 — CerrarSesion y ObtenerSesion + middleware Autenticar

| Campo | Valor |
|-------|-------|
| **Código** | `UT-CU01-BE-04` |
| **Nombre** | CerrarSesion y ObtenerSesion + middleware Autenticar |
| **Caso de uso** | CU-01-BE |
| **Proyecto** | `backend` |
| **Depende de** | UT-CU01-BE-03 |
| **Estado** | Hecho (documentada de forma retroactiva) |

## Antes de implementar

[ARQUITECTURA_BACKEND.md](../../../.claude/rules/ARQUITECTURA_BACKEND.md)

## Solicitud de negocio

Una sesión cerrada no debe servir aunque el token no venza; cada petición debe saber quién es el actor.

## Objetivo

Revocación por `jti` y middleware que valida el token y carga actor y permisos.

## Alcance

- `CerrarSesion` guarda el `jti` revocado hasta su vencimiento.
- `ObtenerSesion` devuelve usuario, rol y permisos.
- `middleware.Autenticar` valida firma, vencimiento, revocación y versión de sesión.
- `RequierePermiso` y `RequiereAlguno`.

Archivos esperados (crear/tocar):

- `backend/internal/application/auth/auth.go`
- `backend/internal/interfaces/http/middleware/autenticar.go`
- `backend/internal/interfaces/http/middleware/permiso.go`
- `backend/internal/interfaces/http/middleware/autenticar_test.go`

## Fuera de alcance (no tocar)

- Rutas de otros módulos.
- Cualquier archivo de `backend` que no figure en «Archivos esperados», y las otras capas del CU (tienen sus propias UT).

## Resultado esperado

Revocación por `jti` y middleware que valida el token y carga actor y permisos.

## Validación

1. `go test ./internal/interfaces/http/middleware -count=1` en verde.
2. `make arch` pasa (domain y application sin dependencias externas).
