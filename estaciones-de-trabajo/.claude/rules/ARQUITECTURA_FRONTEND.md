---
globs: ../frontend/**/*
---

# Arquitectura Frontend (Next.js 16)

Repo `frontend/`: Next.js 16 App Router, React 19, Tailwind 4, zod, pnpm. **Next 16 trae cambios
incompatibles**: leer `frontend/node_modules/next/dist/docs/` antes de tocar su API. Detalle vivo en
`frontend/docs/estado-del-frontend.md`.

## Estructura

| Carpeta | Qué va |
|---|---|
| `app/(auth)` | Login, recuperar, restablecer, activar |
| `app/(publico)` | Portal sin sesión: `/seguimiento`, `/solicitud` |
| `app/panel/<módulo>` | Área autenticada: `page.tsx` (Server Component), `actions.ts` (server actions), formularios cliente junto a la página, `loading.tsx` propio |
| `components/common` · `components/layout` | UI reutilizable · estructura (`PanelShell`, `Sidebar`) |
| `lib/` | `api*.ts` (cliente de la API, `server-only`), `esquemas.ts` (zod), `formularios.ts` (mensajes de error de la API), `navegacion.ts` (menú por permisos) |

## Convenciones que una UT de frontend debe respetar

- El navegador **nunca** habla con la API: Server Components y server actions llaman con `api()` y el
  token va en cookie `httpOnly`. Excepción: la subida directa a S3 por URL firmada.
- Formularios: `useActionState` + zod; errores de la API traducidos en `MENSAJES_API`.
- Cada página exige su permiso (`exigirPermiso`); el menú sale de `lib/navegacion.ts`.
- Buscadores con `next/form` o `CargaListado.tsx`, nunca `<form method="get">` plano.
- Listados con `loading.tsx` → `ListadoEsqueleto`. No hay `loading.tsx` en `app/panel/` ni en la raíz.
- Contraste WCAG AA: un par de colores nuevo se agrega a `scripts/contraste.mjs`.

Validación estándar de una UT de frontend: `pnpm lint` + `pnpm build` + paso manual en el navegador.
