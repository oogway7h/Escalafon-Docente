# Caso de uso

| Id | Nombre | Paquete | Prioridad | Estado |
| --- | --- | --- | --- | --- |
| CU-10 | Revisar Expediente y Registrar Observaciones | P2 Trámites | Alta | Hecho (web) |

**Plataformas:** Web (Hecho) · Mobile (Hecho)
**Actor principal:** Analista asignado
**Actores secundarios:** —
**Objetivo:** Revisar el expediente contra el catálogo normativo y registrar observaciones con documento, motivo y plazo.

**Descripción:**
El analista ve su bandeja de trámites asignados. Para iniciar la revisión debe haber solicitado la carpeta física del docente ligada al trámite; entonces el trámite pasa a En revisión. Contrasta la documentación con los requisitos de la categoría y, si falta algo, registra una observación con el documento afectado, el motivo y el plazo: el trámite pasa a Observado. Al resolver la última observación pendiente vuelve a En revisión.

## Precondiciones

- Sesión con `TRAMITE_REVISAR` y ser el analista asignado ([CU-09](CU-09%20-%20Asignar%20o%20Reasignar%20Tr%C3%A1mite.md)).
- Hay una solicitud de carpeta ligada al trámite que no fue rechazada ([CU-20](CU-20%20-%20Solicitar%20Legajo%20F%C3%ADsico.md)).

## Flujo principal

1. El analista abre «Mi bandeja» y elige un trámite.
2. El sistema muestra datos, fases, requisitos de la categoría, documentos e historial.
3. El analista inicia la revisión; el trámite pasa a En revisión.
4. Registra una observación (documento, motivo, plazo); el trámite pasa a Observado.
5. Cuando el docente subsana, el analista marca la observación como resuelta; sin pendientes, vuelve a En revisión.

## Flujos alternativos y excepciones

- **A1 (paso 3):** sin carpeta solicitada → 409 `CARPETA_NO_SOLICITADA`.
- **A2:** no es el analista asignado → 403/409 `NO_ASIGNADO`.
- **A3 (paso 4):** motivo de menos de 5 caracteres o plazo inválido → 422 `VALIDACION` / `PLAZO_INVALIDO`.
- **A4:** transición no permitida desde el estado actual → 409 `TRANSICION_INVALIDA`.

## Reglas de negocio

- «Observado» es un desvío de «En revisión».
- Registrar una observación más en un trámite ya observado no vuelve a notificar al docente (CU-14).
- Remitido y Archivado son de solo lectura para el analista.

## Postcondiciones

- El trámite está En revisión u Observado con sus observaciones registradas.

## Criterios de aceptación

- **CA01:** Dado un trámite asignado sin carpeta solicitada, el sistema debe impedir iniciar la revisión.
- **CA02:** Dada una observación válida, el trámite debe pasar a Observado y la observación verse en el portal público.
- **CA03:** Dado que se resuelve la última observación, el trámite debe volver a En revisión.
- **CA04:** Dado un analista que no es el asignado, el sistema debe rechazar la revisión.

## Unidades de trabajo

- [UT-BACKEND](../UnidadesTrabajo/CU-10%20-%20Revisar%20Expediente%20y%20Registrar%20Observaciones/UT-BACKEND/README.md) — Hecho
- [UT-FRONTEND](../UnidadesTrabajo/CU-10%20-%20Revisar%20Expediente%20y%20Registrar%20Observaciones/UT-FRONTEND/README.md) — Hecho
- [UT-MOBILE](../UnidadesTrabajo/CU-10%20-%20Revisar%20Expediente%20y%20Registrar%20Observaciones/UT-MOBILE/README.md) — Hecho

## Desarrollador

Carlos Otsubo — especificación y revisión; implementación asistida por IA (ver [BITACORA_IA.md](../BITACORA_IA.md)).
