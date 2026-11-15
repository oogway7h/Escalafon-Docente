# Caso de uso

| Id | Nombre | Paquete | Prioridad | Estado |
| --- | --- | --- | --- | --- |
| CU-14 | Notificar Cambios de Estado al Docente | P2 Trámites | Media | Hecho (web) |

**Plataformas:** Web (Hecho) · Mobile (—)
**Actor principal:** Sistema
**Actores secundarios:** Docente, servidor de correo, Director (consulta los envíos)
**Objetivo:** Avisar por correo al docente cuando su trámite pasa a una etapa relevante, sin exponer datos sensibles.

**Descripción:**
Cuando un trámite pasa a Observado, En facultad, Remitido o Archivado, el sistema registra un aviso y envía un correo al docente con su código de seguimiento y el enlace al portal, sin datos sensibles. Si el envío falla, se reintenta cada minuto hasta tres veces y queda registrado; un fallo nunca revierte la transición. El Director puede consultar los envíos y su estado.

## Precondiciones

- El trámite cambia de estado ([CU-10](CU-10%20-%20Revisar%20Expediente%20y%20Registrar%20Observaciones.md), [CU-11](CU-11%20-%20Emitir%20Dictamen%20y%20Archivar%20Tr%C3%A1mite.md)).

## Flujo principal

1. Ocurre una transición a Observado, En facultad, Remitido o Archivado.
2. El sistema crea el aviso en estado PENDIENTE.
3. Envía el correo en segundo plano tras responder la transición.
4. Marca el aviso ENVIADO.
5. El Director consulta los envíos filtrando por estado.

## Flujos alternativos y excepciones

- **A1 (paso 3):** el envío falla → queda PENDIENTE con el error y un ticker lo reintenta cada minuto.
- **A2:** tres intentos fallidos → FALLIDO.
- **A3:** el docente no tiene correo → FALLIDO sin intentar.
- **A4:** nueva observación sobre un trámite ya Observado → no se envía otro aviso.

## Reglas de negocio

- El correo solo lleva código, estado y enlace al portal.
- La transición nunca se revierte por un fallo de correo.
- Una sola instancia de API (con varias, un aviso podría duplicarse).

## Postcondiciones

- Cada aviso queda registrado con su estado, intentos y error.

## Criterios de aceptación

- **CA01:** Dado que un trámite pasa a Observado, el sistema debe registrar y enviar un aviso al docente.
- **CA02:** Dado que el servidor de correo falla, la transición debe completarse y el aviso reintentarse.
- **CA03:** Dado un docente sin correo, el aviso debe quedar FALLIDO.
- **CA04:** El Director debe poder listar los avisos filtrando por estado.

## Unidades de trabajo

- [UT-BACKEND](../UnidadesTrabajo/CU-14%20-%20Notificar%20Cambios%20de%20Estado%20al%20Docente/UT-BACKEND/README.md) — Hecho
- [UT-FRONTEND](../UnidadesTrabajo/CU-14%20-%20Notificar%20Cambios%20de%20Estado%20al%20Docente/UT-FRONTEND/README.md) — Hecho

## Desarrollador

Carlos Otsubo — especificación y revisión; implementación asistida por IA (ver [BITACORA_IA.md](../BITACORA_IA.md)).
