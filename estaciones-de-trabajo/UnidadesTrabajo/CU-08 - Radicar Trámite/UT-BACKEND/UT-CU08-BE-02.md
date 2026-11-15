# UT-CU08-BE-02 — Casos de uso Radicar, ValidarSolicitud y Listar

| Campo | Valor |
|-------|-------|
| **Código** | `UT-CU08-BE-02` |
| **Nombre** | Casos de uso Radicar, ValidarSolicitud y Listar |
| **Caso de uso** | CU-08-BE |
| **Proyecto** | `backend` |
| **Depende de** | UT-CU08-BE-01 |
| **Estado** | Hecho (documentada de forma retroactiva) |

## Antes de implementar

[ARQUITECTURA_BACKEND.md](../../../.claude/rules/ARQUITECTURA_BACKEND.md)

## Solicitud de negocio

Registrar el ingreso y listar bandejas.

## Objetivo

`Radicar`, `ValidarSolicitud`, `Listar`, `ListarPagina`, `Obtener`.

## Alcance

- Filtros `mios`, `sinAsignar`, `estado`, `buscar`.
- Tests con dobles.

Archivos esperados (crear/tocar):

- `backend/internal/application/tramites/tramites.go`
- `backend/internal/application/tramites/tramites_test.go`
- `backend/internal/infrastructure/persistence/repositorio_tramite.go`

## Fuera de alcance (no tocar)

- Asignación.
- Cualquier archivo de `backend` que no figure en «Archivos esperados», y las otras capas del CU (tienen sus propias UT).

## Resultado esperado

`Radicar`, `ValidarSolicitud`, `Listar`, `ListarPagina`, `Obtener`.

## Validación

1. `go test ./internal/application/tramites -count=1` en verde.
2. `make arch` pasa (domain y application sin dependencias externas).
