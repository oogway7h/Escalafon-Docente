# UT-CU06-BE-01 — Tabla de auditoría de solo anexión

| Campo | Valor |
|-------|-------|
| **Código** | `UT-CU06-BE-01` |
| **Nombre** | Tabla de auditoría de solo anexión |
| **Caso de uso** | CU-06-BE |
| **Proyecto** | `backend` |
| **Depende de** | — |
| **Estado** | Hecho (documentada de forma retroactiva) |

## Antes de implementar

[ARQUITECTURA_BACKEND.md](../../../.claude/rules/ARQUITECTURA_BACKEND.md)

## Solicitud de negocio

Una bitácora alterable no sirve como evidencia.

## Objetivo

Tabla `auditoria` con trigger que rechaza cambios e índice por fecha.

## Alcance

- Trigger en `ajustes_manuales.sql`; `idx_auditoria_fecha`.

Archivos esperados (crear/tocar):

- `backend/internal/domain/auditoria/auditoria.go`
- `backend/internal/infrastructure/persistence/repositorio_auditoria.go`
- `backend/migrations/20260922003354_ajustes_manuales.sql`

## Fuera de alcance (no tocar)

- Consulta (UT 02).
- Cualquier archivo de `backend` que no figure en «Archivos esperados», y las otras capas del CU (tienen sus propias UT).

## Resultado esperado

Tabla `auditoria` con trigger que rechaza cambios e índice por fecha.

## Validación

1. Un `UPDATE auditoria …` en la base falla con el error del trigger.
2. `make arch` pasa (domain y application sin dependencias externas).
