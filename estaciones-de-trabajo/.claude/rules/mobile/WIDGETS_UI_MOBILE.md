---
globs: ../ded_Mobile/**/*
---

# Widgets y UI — Mobile

## Compartidos vs. de una feature

- Un widget de una sola feature vive en `features/{feature}/widgets/`.
- Pasa a `lib/shared/widgets/` cuando **una segunda feature** lo usa de verdad, no antes.

## Tema

`lib/core/theme/app_theme.dart`: Material 3, claro y oscuro (sigue al sistema), espejo manual de
`frontend/app/globals.css`:

| Token web | Claro | Oscuro | Uso |
|---|---|---|---|
| `--color-accent` | `#C8102E` | `#D42740` | relleno de marca (botones), texto blanco encima |
| `--color-accent-ink` | `#C8102E` | `#FF8293` | acento como texto |
| `--color-bg` | `#F5F6F8` | `#0B1220` | fondo |
| `--color-surface` | `#FFFFFF` | `#131C2E` | tarjetas |
| `--color-ink` | `#0F1B2E` | — | texto |
| `--color-sidebar` | `#0F1E36` | — | AppBar |
| `--color-success` | `#15744C` | `#3FB981` | estados OK |
| `--color-grafico` | `#2C5AA0` | `#6A96D8` | barras del tablero |

Un cambio de paleta se hace en los dos lados. En pantallas no se escribe `Color(0x…)`: se usa
`Theme.of(context).colorScheme`.

## Piezas de uso común

`Card`, `ListTile`, `FilledButton`, `Chip` salen del tema. Estados vacíos y errores con un widget
simple (`EstadoVacio`, `ErrorReintentar`) en `shared/widgets/` cuando lo usen dos features. Etiquetas
de estado del trámite y del legajo con los mismos textos que `frontend/lib/tramites.ts` y
`frontend/lib/legajos.ts`.

## Accesibilidad

Objetivos táctiles de 48 dp, `Semantics`/`tooltip` en botones de solo icono, el texto escala con el
sistema (sin tamaños fijos que corten).
