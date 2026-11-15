# Caso de uso

| Id | Nombre | Paquete | Prioridad | Estado |
| --- | --- | --- | --- | --- |
| CU-20 | Solicitar Legajo Físico | P4 Archivo y custodia | Alta | Hecho (web) |

**Plataformas:** Web (Hecho) · Mobile (Hecho)
**Actor principal:** Analista
**Actores secundarios:** —
**Objetivo:** Pedir digitalmente al archivo la carpeta física de un docente para un trámite y seguir su estado.

**Descripción:**
El analista elige uno de sus trámites y la carpeta del docente de ese trámite, escribe el motivo y envía la solicitud, que queda Pendiente. Ve sus solicitudes con su estado (Pendiente, Aprobada, Rechazada con motivo, Prestada, Devuelta). Solo puede haber una solicitud abierta por carpeta y la carpeta debe estar en archivo.

## Precondiciones

- Sesión con `CARPETA_SOLICITAR`.
- El trámite es del analista y no está cerrado ([CU-09](CU-09%20-%20Asignar%20o%20Reasignar%20Tr%C3%A1mite.md)).
- La carpeta es del mismo docente ([CU-19](CU-19%20-%20Gestionar%20Inventario%20de%20Legajos.md)).

## Flujo principal

1. El analista abre «Solicitar carpeta» y elige el trámite.
2. El sistema lista las carpetas del docente de ese trámite.
3. Elige una, escribe el motivo y envía.
4. El sistema crea la solicitud Pendiente y audita.
5. El analista sigue el estado en «Mis solicitudes».

## Flujos alternativos y excepciones

- **A1:** carpeta prestada → 409 `LEGAJO_NO_DISPONIBLE`.
- **A2:** ya hay una solicitud abierta → 409 `SOLICITUD_ABIERTA`.
- **A3:** trámite ajeno, cerrado o de otro docente → 422 `TRAMITE_INVALIDO`.
- **A4:** sin trámite o sin motivo → 422 `TRAMITE_REQUERIDO` / `VALIDACION`.

## Reglas de negocio

- Una solicitud abierta por legajo (índice único parcial `ux_solicitud_abierta`).
- El analista solo ve sus solicitudes.
- La solicitud ligada habilita iniciar la revisión (CU-10).

## Postcondiciones

- Existe una solicitud Pendiente visible para el archivo y el analista.

## Criterios de aceptación

- **CA01:** Dada una carpeta en archivo, la solicitud debe quedar Pendiente.
- **CA02:** Dada una carpeta con solicitud abierta, el sistema debe rechazar otra.
- **CA03:** El analista debe ver solo sus solicitudes con su estado y, si fue rechazada, el motivo.

## Unidades de trabajo

- [UT-BACKEND](../UnidadesTrabajo/CU-20%20-%20Solicitar%20Legajo%20F%C3%ADsico/UT-BACKEND/README.md) — Hecho
- [UT-FRONTEND](../UnidadesTrabajo/CU-20%20-%20Solicitar%20Legajo%20F%C3%ADsico/UT-FRONTEND/README.md) — Hecho
- [UT-MOBILE](../UnidadesTrabajo/CU-20%20-%20Solicitar%20Legajo%20F%C3%ADsico/UT-MOBILE/README.md) — Hecho

## Desarrollador

Carlos Otsubo — especificación y revisión; implementación asistida por IA (ver [BITACORA_IA.md](../BITACORA_IA.md)).
