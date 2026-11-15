# Caso de uso

| Id | Nombre | Paquete | Prioridad | Estado |
| --- | --- | --- | --- | --- |
| CU-04 | Administrar Cuentas de Usuario | P1 Seguridad | Alta | Hecho (web) |

**Plataformas:** Web (Hecho) · Mobile (—)
**Actor principal:** Director
**Actores secundarios:** Servidor de correo
**Objetivo:** Registrar, editar, cambiar de rol y dar de baja lógica las cuentas del personal.

**Descripción:**
El Director lista las cuentas, crea una nueva con nombres, correo institucional y rol (el sistema envía el enlace de activación), edita sus datos, cambia su rol o la da de baja y la reactiva. El correo no puede repetirse. El Director no puede darse de baja ni cambiarse el rol a sí mismo. Toda operación queda en auditoría y el historial de actividad se conserva.

## Precondiciones

- Sesión con `USUARIO_GESTIONAR` ([CU-01](CU-01%20-%20Iniciar%20Sesi%C3%B3n.md)).
- Existen roles ([CU-05](CU-05%20-%20Gestionar%20Roles%20y%20Permisos.md)).

## Flujo principal

1. El Director abre Usuarios y ve el listado con búsqueda.
2. Elige «Nuevo», completa nombres, apellidos, correo y rol.
3. El sistema verifica que el correo no exista, crea la cuenta inactiva y envía el enlace de activación.
4. Para editar, cambia los datos y guarda.
5. Para dar de baja, confirma; la cuenta queda inactiva y sus sesiones se cierran.
6. Cada operación queda en auditoría.

## Flujos alternativos y excepciones

- **A1 (paso 3):** correo ya registrado (sin distinguir mayúsculas) → 409 `CORREO_DUPLICADO`.
- **A2 (paso 5):** el Director intenta darse de baja → 409 `AUTODESACTIVACION`.
- **A3:** el Director intenta cambiar su propio rol → 409 `AUTOCAMBIO_ROL`.
- **A4:** rol inexistente → 422 `ROL_INVALIDO`.
- **A5:** el usuario no activó su cuenta → «Reenviar activación» genera un enlace nuevo.

## Reglas de negocio

- La baja es lógica: no se borra la cuenta ni su historial.
- Correo único sin distinguir mayúsculas (índice manual `ux_usuario_correo`).
- Dar de baja o cambiar de rol invalida las sesiones del usuario.

## Postcondiciones

- La cuenta existe/cambió y la operación está en auditoría.

## Criterios de aceptación

- **CA01:** Dado un correo nuevo, el sistema debe crear la cuenta y enviar el enlace de activación.
- **CA02:** Dado un correo existente, el sistema debe responder `CORREO_DUPLICADO`.
- **CA03:** Dado que el Director se da de baja a sí mismo, el sistema debe impedirlo.
- **CA04:** Dado que una cuenta se da de baja, sus sesiones deben dejar de funcionar.

## Unidades de trabajo

- [UT-BACKEND](../UnidadesTrabajo/CU-04%20-%20Administrar%20Cuentas%20de%20Usuario/UT-BACKEND/README.md) — Hecho
- [UT-FRONTEND](../UnidadesTrabajo/CU-04%20-%20Administrar%20Cuentas%20de%20Usuario/UT-FRONTEND/README.md) — Hecho

## Desarrollador

Carlos Otsubo — especificación y revisión; implementación asistida por IA (ver [BITACORA_IA.md](../BITACORA_IA.md)).
