# UT-CU03-BE-01 — Casos de uso del perfil

| Campo | Valor |
|-------|-------|
| **Código** | `UT-CU03-BE-01` |
| **Nombre** | Casos de uso del perfil |
| **Caso de uso** | CU-03-BE |
| **Proyecto** | `backend` |
| **Depende de** | — |
| **Estado** | Hecho (documentada de forma retroactiva) |

## Antes de implementar

[ARQUITECTURA_BACKEND.md](../../../.claude/rules/ARQUITECTURA_BACKEND.md)

## Solicitud de negocio

El usuario mantiene sus datos sin pedirlo al Director.

## Objetivo

`ObtenerPerfil`, `ActualizarContacto`, `CambiarContrasena` con auditoría.

## Alcance

- Rechaza campos no editables.
- Verifica la contraseña actual con el hasher.
- Tests con dobles.

Archivos esperados (crear/tocar):

- `backend/internal/application/perfil/perfil.go`
- `backend/internal/application/perfil/perfil_test.go`

## Fuera de alcance (no tocar)

- HTTP (UT 02).
- Cualquier archivo de `backend` que no figure en «Archivos esperados», y las otras capas del CU (tienen sus propias UT).

## Resultado esperado

`ObtenerPerfil`, `ActualizarContacto`, `CambiarContrasena` con auditoría.

## Validación

1. `go test ./internal/application/perfil -count=1` en verde.
2. `make arch` pasa (domain y application sin dependencias externas).
