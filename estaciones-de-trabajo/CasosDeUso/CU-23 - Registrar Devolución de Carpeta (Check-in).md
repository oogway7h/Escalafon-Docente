# Caso de uso

| Id | Nombre | Paquete | Prioridad | Estado |
| --- | --- | --- | --- | --- |
| CU-23 | Registrar Devolución de Carpeta (Check-in) | P4 Archivo y custodia | Alta | Hecho (web) |

**Plataformas:** Web (Hecho) · Mobile (Hecho)
**Actor principal:** Personal de Archivo
**Actores secundarios:** —
**Objetivo:** Registrar el retorno de una carpeta catalogando su estado material y responsabilizando al último custodio.

**Descripción:**
Cuando la carpeta vuelve, el Personal de Archivo registra la devolución indicando el estado material (Buen estado, Hojas sueltas o Deteriorado) y una observación, obligatoria si no vuelve en buen estado. El sistema marca el legajo En archivo, la solicitud Devuelta y crea un asiento de entrada en la bitácora asignando la responsabilidad al último custodio. La alerta de mora del tablero desaparece sola.

## Precondiciones

- Sesión con `CUSTODIA_REGISTRAR`.
- El legajo está Prestado ([CU-22](CU-22%20-%20Despachar%20Carpeta%20F%C3%ADsica%20%28Check-out%29.md)).

## Flujo principal

1. Archivo abre la pestaña Prestadas o Vencidas.
2. Elige la carpeta y «Registrar devolución».
3. Indica el estado material y la observación.
4. El sistema, en una transacción, marca el legajo En archivo, la solicitud Devuelta y crea el movimiento de entrada con el custodio.

## Flujos alternativos y excepciones

- **A1:** estado distinto de Buen estado sin observación → 422 `OBSERVACION_REQUERIDA`.
- **A2:** el legajo no está prestado → 409 `LEGAJO_NO_PRESTADO`.
- **A3:** estado material desconocido → 422 `ESTADO_MATERIAL_INVALIDO`.

## Reglas de negocio

- Desvío de `cu.md`: sin evidencia fotográfica (se quitó `evidencia_clave`).
- La bitácora es de solo anexión.
- Un trámite solo se archiva con todas sus carpetas devueltas (CU-11).

## Postcondiciones

- El legajo está En archivo y la entrada en la bitácora.

## Criterios de aceptación

- **CA01:** Dada una devolución en buen estado, el legajo debe quedar En archivo.
- **CA02:** Dada una devolución deteriorada sin observación, el sistema debe rechazarla.
- **CA03:** Tras la devolución, la alerta de mora debe desaparecer del tablero.

## Unidades de trabajo

- [UT-BACKEND](../UnidadesTrabajo/CU-23%20-%20Registrar%20Devoluci%C3%B3n%20de%20Carpeta%20%28Check-in%29/UT-BACKEND/README.md) — Hecho
- [UT-FRONTEND](../UnidadesTrabajo/CU-23%20-%20Registrar%20Devoluci%C3%B3n%20de%20Carpeta%20%28Check-in%29/UT-FRONTEND/README.md) — Hecho
- [UT-MOBILE](../UnidadesTrabajo/CU-23%20-%20Registrar%20Devoluci%C3%B3n%20de%20Carpeta%20%28Check-in%29/UT-MOBILE/README.md) — Hecho

## Desarrollador

Carlos Otsubo — especificación y revisión; implementación asistida por IA (ver [BITACORA_IA.md](../BITACORA_IA.md)).
