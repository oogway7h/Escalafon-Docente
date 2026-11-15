# UT-CU10-FE-03 — Acciones de revisión y observaciones

| Campo | Valor |
|-------|-------|
| **Código** | `UT-CU10-FE-03` |
| **Nombre** | Acciones de revisión y observaciones |
| **Caso de uso** | CU-10-FE |
| **Proyecto** | `frontend` |
| **Depende de** | UT-CU10-FE-02 |
| **Estado** | Hecho (documentada de forma retroactiva) |

## Antes de implementar

[ARQUITECTURA_FRONTEND.md](../../../.claude/rules/ARQUITECTURA_FRONTEND.md)

## Solicitud de negocio

Registrar observaciones sin salir del detalle.

## Objetivo

Formularios de iniciar revisión, observar y resolver.

## Alcance

- Sin carpeta ofrece «Solicitar carpeta».

Archivos esperados (crear/tocar):

- `frontend/app/panel/tramites/[id]/FormulariosTramite.tsx`
- `frontend/app/panel/tramites/actions.ts`

## Fuera de alcance (no tocar)

- Dictamen.
- Cualquier archivo de `frontend` que no figure en «Archivos esperados», y las otras capas del CU (tienen sus propias UT).

## Resultado esperado

Formularios de iniciar revisión, observar y resolver.

## Validación

1. Con `pnpm dev` y la API levantada, registrar una observación cambia el estado a Observado.
2. `pnpm lint` y `pnpm build` en `frontend/` sin errores.
