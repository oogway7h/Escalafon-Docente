# Caso de uso

| Id | Nombre | Paquete | Prioridad | Estado |
| --- | --- | --- | --- | --- |
| CU-24 | Consultar Bitácora de Custodia | P4 Archivo y custodia | Media | Hecho (web) |

**Plataformas:** Web (Hecho) · Mobile (—)
**Actor principal:** Director (con `BITACORA_CONSULTAR`)
**Actores secundarios:** —
**Objetivo:** Consultar el historial inmutable de movimientos de las carpetas.

**Descripción:**
El Director ve los movimientos de salida y entrada de las carpetas en orden cronológico, filtra por texto (legajo, docente o custodio), tipo de movimiento o solo vencidos, y obtiene el encadenamiento completo de la cadena de custodia de un legajo. Cada salida indica si está vencida.

## Precondiciones

- Sesión con `BITACORA_CONSULTAR`.

## Flujo principal

1. El Director abre la bitácora.
2. Filtra por texto, tipo o vencidos.
3. El sistema devuelve los movimientos paginados con custodio, fechas, estado material y quién registró.

## Flujos alternativos y excepciones

- **A1:** tipo desconocido → 422 `FILTRO_INVALIDO`.
- **A2:** sin resultados → lista vacía.

## Reglas de negocio

- La bitácora no se puede modificar (trigger).
- La fecha límite se compara por fecha (Bolivia es UTC-4).

## Postcondiciones

- La consulta no modifica nada.

## Criterios de aceptación

- **CA01:** Dado un código de legajo, el sistema debe devolver su cadena de custodia en orden.
- **CA02:** Dado el filtro vencidos, solo deben salir préstamos fuera de plazo.

## Unidades de trabajo

- [UT-BACKEND](../UnidadesTrabajo/CU-24%20-%20Consultar%20Bit%C3%A1cora%20de%20Custodia/UT-BACKEND/README.md) — Hecho
- [UT-FRONTEND](../UnidadesTrabajo/CU-24%20-%20Consultar%20Bit%C3%A1cora%20de%20Custodia/UT-FRONTEND/README.md) — Hecho

## Desarrollador

Carlos Otsubo — especificación y revisión; implementación asistida por IA (ver [BITACORA_IA.md](../BITACORA_IA.md)).
