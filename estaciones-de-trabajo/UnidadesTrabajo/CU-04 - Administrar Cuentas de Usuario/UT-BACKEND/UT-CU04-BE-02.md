# UT-CU04-BE-02 — Casos de uso de usuarios

| Campo | Valor |
|-------|-------|
| **Código** | `UT-CU04-BE-02` |
| **Nombre** | Casos de uso de usuarios |
| **Caso de uso** | CU-04-BE |
| **Proyecto** | `backend` |
| **Depende de** | UT-CU04-BE-01 |
| **Estado** | Hecho (documentada de forma retroactiva) |

## Antes de implementar

[ARQUITECTURA_BACKEND.md](../../../.claude/rules/ARQUITECTURA_BACKEND.md)

## Solicitud de negocio

El Director gestiona al personal sin perder el historial.

## Objetivo

`Crear`, `Editar`, `CambiarRol`, `CambiarEstado`, `ReenviarActivacion` con auditoría.

## Alcance

- Autoprotección del Director.
- Subir `version_sesion` al dar de baja o cambiar rol.
- Tests con dobles.

Archivos esperados (crear/tocar):

- `backend/internal/application/usuarios/usuarios.go`
- `backend/internal/application/usuarios/usuarios_test.go`

## Fuera de alcance (no tocar)

- Roles (CU-05).
- Cualquier archivo de `backend` que no figure en «Archivos esperados», y las otras capas del CU (tienen sus propias UT).

## Resultado esperado

`Crear`, `Editar`, `CambiarRol`, `CambiarEstado`, `ReenviarActivacion` con auditoría.

## Validación

1. `go test ./internal/application/usuarios -count=1` en verde.
2. `make arch` pasa (domain y application sin dependencias externas).
