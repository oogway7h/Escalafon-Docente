---
globs: ../ded_Mobile/**/*
---

# Formularios — Mobile

- `Form` + `TextFormField` + `validator` nativos. Sin librerías de formularios.
- Mensajes en español, escritos en el `validator`.
- **Misma regla que la web** para el mismo campo (`frontend/lib/esquemas.ts`): p. ej. motivo de
  reasignación y nota de rechazo obligatorios, observación obligatoria si la carpeta no vuelve en buen
  estado. Si cambia una regla, se cambia en los dos lados en el mismo CU.
- El backend es la fuente de verdad: el error que devuelve se muestra con `mensajeDeError(e)`.
- Selectores: `DropdownButtonFormField` y `showDatePicker` del framework (en móvil son controles
  propios de Flutter, no nativos del navegador).
