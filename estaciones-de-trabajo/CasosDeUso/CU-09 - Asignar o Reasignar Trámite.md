# Caso de uso

| Id | Nombre | Paquete | Prioridad | Estado |
| --- | --- | --- | --- | --- |
| CU-09 | Asignar o Reasignar Trámite | P2 Trámites | Alta | Hecho (web) |

**Plataformas:** Web (Hecho) · Mobile (Hecho)
**Actor principal:** Director
**Actores secundarios:** —
**Objetivo:** Distribuir los trámites ingresados entre los analistas según su carga vigente.

**Descripción:**
El Director ve los trámites sin asignar y la carga vigente de cada analista. Elige un analista y asigna: el trámite pasa a Asignado. Para reasignar un trámite ya asignado debe escribir el motivo, que queda permanente en el historial; la reasignación no cambia el estado. No se puede reasignar un trámite remitido o archivado.

## Precondiciones

- Sesión con `TRAMITE_ASIGNAR`.
- Existe al menos un usuario activo con `TRAMITE_REVISAR`.
- El trámite fue radicado ([CU-08](CU-08%20-%20Radicar%20Tr%C3%A1mite.md)) y, si vino del portal, validado.

## Flujo principal

1. El Director abre la bandeja de asignación.
2. El sistema lista los trámites en Ingresado y los analistas con su carga vigente.
3. El Director elige el analista y confirma.
4. El sistema pasa el trámite a Asignado, registra la etapa y audita.

## Flujos alternativos y excepciones

- **A1 (reasignar):** sin motivo → 422 `MOTIVO_REQUERIDO`.
- **A2:** usuario sin `TRAMITE_REVISAR` o inactivo → 422 `ANALISTA_INVALIDO`.
- **A3:** trámite remitido o archivado → 409 `TRAMITE_REMITIDO` / `TRAMITE_ARCHIVADO`.
- **A4:** solicitud del portal sin validar → 409 `PENDIENTE_VALIDACION`.

## Reglas de negocio

- Carga vigente = trámites en Asignado, En revisión, Observado, Dictaminado y Revisión final (no En facultad).
- La reasignación no cambia el estado y su motivo queda en el historial.
- Se puede reasignar en cualquier estado salvo Remitido y Archivado.

## Postcondiciones

- El trámite tiene analista y la etapa está en el historial.

## Criterios de aceptación

- **CA01:** Dado un trámite Ingresado, asignarlo debe dejarlo en Asignado con el analista elegido.
- **CA02:** Dado un trámite asignado, reasignarlo sin motivo debe rechazarse con `MOTIVO_REQUERIDO`.
- **CA03:** Dada una reasignación válida, el estado no debe cambiar y el motivo debe verse en el historial.
- **CA04:** La bandeja debe mostrar la carga vigente de cada analista.

## Unidades de trabajo

- [UT-BACKEND](../UnidadesTrabajo/CU-09%20-%20Asignar%20o%20Reasignar%20Tr%C3%A1mite/UT-BACKEND/README.md) — Hecho
- [UT-FRONTEND](../UnidadesTrabajo/CU-09%20-%20Asignar%20o%20Reasignar%20Tr%C3%A1mite/UT-FRONTEND/README.md) — Hecho
- [UT-MOBILE](../UnidadesTrabajo/CU-09%20-%20Asignar%20o%20Reasignar%20Tr%C3%A1mite/UT-MOBILE/README.md) — Hecho

## Desarrollador

Carlos Otsubo — especificación y revisión; implementación asistida por IA (ver [BITACORA_IA.md](../BITACORA_IA.md)).
