# UT-CU05-BE-02 — Casos de uso de roles

| Campo | Valor |
|-------|-------|
| **Código** | `UT-CU05-BE-02` |
| **Nombre** | Casos de uso de roles |
| **Caso de uso** | CU-05-BE |
| **Proyecto** | `backend` |
| **Depende de** | UT-CU05-BE-01 |
| **Estado** | Hecho (documentada de forma retroactiva) |

## Antes de implementar

[ARQUITECTURA_BACKEND.md](../../../.claude/rules/ARQUITECTURA_BACKEND.md)

## Solicitud de negocio

Evitar que un cambio de roles deje al sistema sin administrador.

## Objetivo

`Listar`, `Obtener`, `ListarPermisos`, `Guardar`, `Eliminar`.

## Alcance

- Validaciones A1-A4 y auditoría.
- Tests con dobles.

Archivos esperados (crear/tocar):

- `backend/internal/application/roles/roles.go`
- `backend/internal/application/roles/roles_test.go`
- `backend/internal/infrastructure/persistence/repositorio_rol.go`

## Fuera de alcance (no tocar)

- HTTP (UT 03).
- Cualquier archivo de `backend` que no figure en «Archivos esperados», y las otras capas del CU (tienen sus propias UT).

## Resultado esperado

`Listar`, `Obtener`, `ListarPermisos`, `Guardar`, `Eliminar`.

## Validación

1. `go test ./internal/application/roles -count=1` en verde.
2. `make arch` pasa (domain y application sin dependencias externas).
