# UT-CU13-BE-02 — Caso de uso RadicarPublico

| Campo | Valor |
|-------|-------|
| **Código** | `UT-CU13-BE-02` |
| **Nombre** | Caso de uso RadicarPublico |
| **Caso de uso** | CU-13-BE |
| **Proyecto** | `backend` |
| **Depende de** | UT-CU13-BE-01 |
| **Estado** | Hecho (documentada de forma retroactiva) |

## Antes de implementar

[ARQUITECTURA_BACKEND.md](../../../.claude/rules/ARQUITECTURA_BACKEND.md)

## Solicitud de negocio

El docente radica desde su casa.

## Objetivo

`RadicarPublico` valida identidad, requisitos, tipo y tamaño y guarda con hash.

## Alcance

- Validación por contenido.
- Tests con almacenamiento falso.

Archivos esperados (crear/tocar):

- `backend/internal/application/tramites/tramites.go`
- `backend/internal/application/tramites/tramites_test.go`

## Fuera de alcance (no tocar)

- Validación en ventanilla (CU-08).
- Cualquier archivo de `backend` que no figure en «Archivos esperados», y las otras capas del CU (tienen sus propias UT).

## Resultado esperado

`RadicarPublico` valida identidad, requisitos, tipo y tamaño y guarda con hash.

## Validación

1. `go test ./internal/application/tramites -count=1` en verde.
2. `make arch` pasa (domain y application sin dependencias externas).
