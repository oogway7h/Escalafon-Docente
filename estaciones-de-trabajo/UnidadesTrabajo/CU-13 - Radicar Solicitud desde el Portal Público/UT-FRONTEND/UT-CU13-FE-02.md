# UT-CU13-FE-02 — Requisitos por categoría en la portada

| Campo | Valor |
|-------|-------|
| **Código** | `UT-CU13-FE-02` |
| **Nombre** | Requisitos por categoría en la portada |
| **Caso de uso** | CU-13-FE |
| **Proyecto** | `frontend` |
| **Depende de** | UT-CU13-FE-01 |
| **Estado** | Hecho (documentada de forma retroactiva) |

## Antes de implementar

[ARQUITECTURA_FRONTEND.md](../../../.claude/rules/ARQUITECTURA_FRONTEND.md)

## Solicitud de negocio

El docente sabe qué preparar antes de empezar.

## Objetivo

Tarjetas de dos caras en `CategoriasRequisitos`.

## Alcance

- Lee `/publico/categorias`; si falla, la sección se omite.

Archivos esperados (crear/tocar):

- `frontend/components/publico/CategoriasRequisitos.tsx`
- `frontend/app/page.tsx`

## Fuera de alcance (no tocar)

- Cualquier archivo de `frontend` que no figure en «Archivos esperados», y las otras capas del CU (tienen sus propias UT).

## Resultado esperado

Tarjetas de dos caras en `CategoriasRequisitos`.

## Validación

1. Con `pnpm dev` y la API levantada, tocar una tarjeta muestra sus documentos y «Iniciar mi solicitud».
2. `pnpm lint` y `pnpm build` en `frontend/` sin errores.
