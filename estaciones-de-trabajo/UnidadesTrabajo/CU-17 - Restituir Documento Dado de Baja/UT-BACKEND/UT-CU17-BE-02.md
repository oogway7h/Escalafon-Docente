# UT-CU17-BE-02 — Rutas de restitución

| Campo | Valor |
|-------|-------|
| **Código** | `UT-CU17-BE-02` |
| **Nombre** | Rutas de restitución |
| **Caso de uso** | CU-17-BE |
| **Proyecto** | `backend` |
| **Depende de** | UT-CU17-BE-01 |
| **Estado** | Hecho (documentada de forma retroactiva) |

## Antes de implementar

[ARQUITECTURA_BACKEND.md](../../../.claude/rules/ARQUITECTURA_BACKEND.md)

## Solicitud de negocio

Exponer al Director.

## Objetivo

Rutas `/documentos/bajas` y `/restituir`.

## Alcance

- `RestituirEntrada`, `ListaBajasDTO`.

Archivos esperados (crear/tocar):

- `backend/internal/interfaces/http/handlers/documentos.go`
- `backend/internal/interfaces/http/router.go`

## Fuera de alcance (no tocar)

- Cualquier archivo de `backend` que no figure en «Archivos esperados», y las otras capas del CU (tienen sus propias UT).

## Resultado esperado

Rutas `/documentos/bajas` y `/restituir`.

## Validación

1. `go test ./internal/interfaces/http -count=1` en verde.
2. `make arch` pasa (domain y application sin dependencias externas).
