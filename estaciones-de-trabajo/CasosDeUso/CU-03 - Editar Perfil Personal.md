# Caso de uso

| Id | Nombre | Paquete | Prioridad | Estado |
| --- | --- | --- | --- | --- |
| CU-03 | Editar Perfil Personal | P1 Seguridad | Media | Hecho (web) |

**Plataformas:** Web (Hecho) · Mobile (—)
**Actor principal:** Usuario interno
**Actores secundarios:** —
**Objetivo:** Actualizar sus datos de contacto y su contraseña sin poder tocar su rol ni sus permisos.

**Descripción:**
Desde su perfil, el usuario ve sus datos y puede cambiar teléfono y datos de contacto, o cambiar su contraseña confirmando la actual. El correo, el rol y los permisos no se pueden editar desde aquí. Si la contraseña actual no coincide o la nueva no cumple la política, el sistema lo informa.

## Precondiciones

- El usuario tiene sesión iniciada ([CU-01](CU-01%20-%20Iniciar%20Sesi%C3%B3n.md)) y el permiso `PERFIL_GESTIONAR`.

## Flujo principal

1. El usuario abre su perfil.
2. El sistema muestra nombres, correo, rol y contacto.
3. El usuario edita su contacto y guarda; el sistema lo actualiza y audita.
4. Para la contraseña, escribe la actual y la nueva; el sistema verifica la actual, aplica la política y guarda.

## Flujos alternativos y excepciones

- **A1 (paso 4):** contraseña actual incorrecta → 422 `CONTRASENA_ACTUAL_INCORRECTA`.
- **A2 (paso 4):** contraseña nueva débil → 422 `POLITICA_CONTRASENA`.
- **A3 (paso 3):** intento de cambiar correo, rol o permisos → 422 `CAMPO_NO_EDITABLE`.

## Reglas de negocio

- Nadie se cambia su propio rol ni permisos (eso es CU-04/CU-05).
- El cambio de contraseña exige la actual.

## Postcondiciones

- Los datos de contacto o la contraseña quedan actualizados y auditados.

## Criterios de aceptación

- **CA01:** Dado que el usuario cambia su teléfono, el sistema debe guardarlo y auditar `ACTUALIZAR_CONTACTO`.
- **CA02:** Dado que la contraseña actual no coincide, el sistema debe rechazar el cambio.
- **CA03:** Dado que se envía un rol en el cuerpo, el sistema debe rechazarlo con `CAMPO_NO_EDITABLE`.

## Unidades de trabajo

- [UT-BACKEND](../UnidadesTrabajo/CU-03%20-%20Editar%20Perfil%20Personal/UT-BACKEND/README.md) — Hecho
- [UT-FRONTEND](../UnidadesTrabajo/CU-03%20-%20Editar%20Perfil%20Personal/UT-FRONTEND/README.md) — Hecho

## Desarrollador

Carlos Otsubo — especificación y revisión; implementación asistida por IA (ver [BITACORA_IA.md](../BITACORA_IA.md)).
