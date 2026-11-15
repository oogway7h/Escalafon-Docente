# UT-CU11-BE-02 — Remisiones, regreso y archivo

| Campo | Valor |
|-------|-------|
| **Código** | `UT-CU11-BE-02` |
| **Nombre** | Remisiones, regreso y archivo |
| **Caso de uso** | CU-11-BE |
| **Proyecto** | `backend` |
| **Depende de** | UT-CU11-BE-01 |
| **Estado** | Hecho (documentada de forma retroactiva) |

## Antes de implementar

[ARQUITECTURA_BACKEND.md](../../../.claude/rules/ARQUITECTURA_BACKEND.md)

## Solicitud de negocio

Seguir el expediente fuera y de vuelta al Departamento.

## Objetivo

`RemitirAFacultad`, `RegistrarRegreso`, `RemitirARectorado`, `Archivar`.

## Alcance

- Archivar exige carpetas devueltas.
- Tests del flujo completo en `flujo_test.go`.

Archivos esperados (crear/tocar):

- `backend/internal/application/tramites/flujo.go`
- `backend/internal/application/tramites/tramites.go`
- `backend/internal/application/tramites/flujo_test.go`

## Fuera de alcance (no tocar)

- Notificaciones (CU-14).
- Cualquier archivo de `backend` que no figure en «Archivos esperados», y las otras capas del CU (tienen sus propias UT).

## Resultado esperado

`RemitirAFacultad`, `RegistrarRegreso`, `RemitirARectorado`, `Archivar`.

## Validación

1. `go test ./internal/application/tramites -count=1` en verde.
2. `make arch` pasa (domain y application sin dependencias externas).
