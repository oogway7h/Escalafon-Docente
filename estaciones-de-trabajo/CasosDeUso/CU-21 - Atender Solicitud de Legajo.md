# Caso de uso

| Id | Nombre | Paquete | Prioridad | Estado |
| --- | --- | --- | --- | --- |
| CU-21 | Atender Solicitud de Legajo | P4 Archivo y custodia | Alta | Hecho (web) |

**Plataformas:** Web (Hecho) · Mobile (Hecho)
**Actor principal:** Personal de Archivo
**Actores secundarios:** —
**Objetivo:** Aprobar o rechazar las solicitudes de carpetas recibidas, con motivo al rechazar.

**Descripción:**
El Personal de Archivo ve en su bandeja las solicitudes Pendientes. Aprueba una (queda Aprobada, lista para despachar) o la rechaza escribiendo el motivo. El analista solicitante ve el nuevo estado y el motivo.

## Precondiciones

- Sesión con `LEGAJO_GESTIONAR`.
- Existe una solicitud Pendiente ([CU-20](CU-20%20-%20Solicitar%20Legajo%20F%C3%ADsico.md)).

## Flujo principal

1. Archivo abre «Solicitudes y custodia», pestaña Por atender.
2. Elige una solicitud.
3. Aprueba o rechaza con nota.
4. El sistema actualiza el estado, registra la respuesta y audita.

## Flujos alternativos y excepciones

- **A1:** rechazo sin nota → 422 `NOTA_REQUERIDA`.
- **A2:** solicitud que ya no está Pendiente → 409 `SOLICITUD_ESTADO_INVALIDO`.

## Reglas de negocio

- Rechazar libera la carpeta para otra solicitud.
- «Aprobar y entregar» (CU-22) hace los dos pasos juntos.

## Postcondiciones

- La solicitud está Aprobada o Rechazada con su respuesta.

## Criterios de aceptación

- **CA01:** Dada una solicitud Pendiente, aprobarla debe dejarla Aprobada.
- **CA02:** Dado un rechazo sin nota, el sistema debe impedirlo.
- **CA03:** El analista debe ver el motivo del rechazo.

## Unidades de trabajo

- [UT-BACKEND](../UnidadesTrabajo/CU-21%20-%20Atender%20Solicitud%20de%20Legajo/UT-BACKEND/README.md) — Hecho
- [UT-FRONTEND](../UnidadesTrabajo/CU-21%20-%20Atender%20Solicitud%20de%20Legajo/UT-FRONTEND/README.md) — Hecho
- [UT-MOBILE](../UnidadesTrabajo/CU-21%20-%20Atender%20Solicitud%20de%20Legajo/UT-MOBILE/README.md) — Hecho

## Desarrollador

Carlos Otsubo — especificación y revisión; implementación asistida por IA (ver [BITACORA_IA.md](../BITACORA_IA.md)).
