# Caso de uso

| Id | Nombre | Paquete | Prioridad | Estado |
| --- | --- | --- | --- | --- |
| CU-12 | Consultar Estado de Trámite (Portal Público) | P2 Trámites | Alta | Hecho (web) |

**Plataformas:** Web (Hecho) · Mobile (Hecho)
**Actor principal:** Docente (sin cuenta)
**Actores secundarios:** —
**Objetivo:** Conocer la fase actual, el historial y las observaciones de su trámite sin iniciar sesión.

**Descripción:**
El docente escribe su código de seguimiento y ve la fase actual, la línea de etapas y las observaciones (documento, motivo, plazo). Si no recuerda el código, puede escribir su código de docente junto con su cédula y elegir entre sus trámites. El sistema nunca expone datos personales ni de otros docentes, y un error no revela si falló el código o la cédula.

## Precondiciones

- El trámite fue radicado ([CU-08](CU-08%20-%20Radicar%20Tr%C3%A1mite.md) o [CU-13](CU-13%20-%20Radicar%20Solicitud%20desde%20el%20Portal%20P%C3%BAblico.md)).

## Flujo principal

1. El docente abre el seguimiento.
2. Escribe el código de seguimiento.
3. El sistema muestra código, tipo, categoría, estado con su descripción, etapas y observaciones.
4. Alternativa: escribe su código de docente y su cédula; el sistema lista sus trámites y el docente elige uno.

## Flujos alternativos y excepciones

- **A1:** código inexistente → 404 con mensaje genérico.
- **A2:** código de docente y cédula que no coinciden, o docente sin cédula → lista vacía (no revela si el docente existe).
- **A3:** más de 30 consultas por minuto desde la misma IP → 429.
- **A4:** trámite En facultad → se explica que esa etapa no tiene seguimiento.

## Reglas de negocio

- Solo se devuelven código, tipo, categoría, estado, etapas y observaciones; nunca nombres, cédula ni analista.
- La cédula viaja en el cuerpo de la petición, nunca en la URL.
- El Resultado se muestra recién desde Remitido.

## Postcondiciones

- La consulta no modifica nada.

## Criterios de aceptación

- **CA01:** Dado un código válido, el sistema debe mostrar el estado y las etapas sin datos personales.
- **CA02:** Dado un código inexistente, el sistema debe responder un mensaje genérico.
- **CA03:** Dado un código de docente con cédula incorrecta, el sistema debe responder lista vacía.
- **CA04:** Dado un trámite Observado, el docente debe ver documento, motivo y plazo de cada observación.

## Unidades de trabajo

- [UT-BACKEND](../UnidadesTrabajo/CU-12%20-%20Consultar%20Estado%20de%20Tr%C3%A1mite%20%28Portal%20P%C3%BAblico%29/UT-BACKEND/README.md) — Hecho
- [UT-FRONTEND](../UnidadesTrabajo/CU-12%20-%20Consultar%20Estado%20de%20Tr%C3%A1mite%20%28Portal%20P%C3%BAblico%29/UT-FRONTEND/README.md) — Hecho
- [UT-MOBILE](../UnidadesTrabajo/CU-12%20-%20Consultar%20Estado%20de%20Tr%C3%A1mite%20%28Portal%20P%C3%BAblico%29/UT-MOBILE/README.md) — Hecho

## Desarrollador

Carlos Otsubo — especificación y revisión; implementación asistida por IA (ver [BITACORA_IA.md](../BITACORA_IA.md)).
