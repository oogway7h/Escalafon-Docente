# UT-CU13-FE-01 — Formulario de solicitud

| Campo | Valor |
|-------|-------|
| **Código** | `UT-CU13-FE-01` |
| **Nombre** | Formulario de solicitud |
| **Caso de uso** | CU-13-FE |
| **Proyecto** | `frontend` |
| **Depende de** | — |
| **Estado** | Hecho (documentada de forma retroactiva) |

## Antes de implementar

[ARQUITECTURA_FRONTEND.md](../../../.claude/rules/ARQUITECTURA_FRONTEND.md)

## Solicitud de negocio

Radicar sin ir al Departamento.

## Objetivo

`/solicitud` con `SolicitudForm`.

## Alcance

- Un campo de archivo por requisito.
- Redirige a `/seguimiento?codigo=…&nueva=1`.

Archivos esperados (crear/tocar):

- `frontend/app/(publico)/solicitud/page.tsx`
- `frontend/app/(publico)/solicitud/SolicitudForm.tsx`
- `frontend/app/(publico)/solicitud/actions.ts`
- `frontend/next.config.ts`

## Fuera de alcance (no tocar)

- Seguimiento.
- Cualquier archivo de `frontend` que no figure en «Archivos esperados», y las otras capas del CU (tienen sus propias UT).

## Resultado esperado

`/solicitud` con `SolicitudForm`.

## Validación

1. Con `pnpm dev` y la API levantada, una solicitud completa termina en el seguimiento con su código nuevo.
2. `pnpm lint` y `pnpm build` en `frontend/` sin errores.
