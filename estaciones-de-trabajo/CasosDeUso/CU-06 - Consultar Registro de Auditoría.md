# Caso de uso

| Id | Nombre | Paquete | Prioridad | Estado |
| --- | --- | --- | --- | --- |
| CU-06 | Consultar Registro de Auditoría | P1 Seguridad | Media | Hecho (web) |

**Plataformas:** Web (Hecho) · Mobile (—)
**Actor principal:** Director
**Actores secundarios:** —
**Objetivo:** Consultar la bitácora de acciones filtrando por usuario, fechas o tipo de acción.

**Descripción:**
El Director abre la bitácora y ve las acciones de los usuarios del más reciente al más antiguo, paginadas. Puede combinar filtros por usuario, tipo de acción y rango de fechas. Los registros son de solo lectura para todos los roles: no existe ninguna operación que los modifique o borre.

## Precondiciones

- Sesión con `AUDITORIA_CONSULTAR`.

## Flujo principal

1. El Director abre Auditoría.
2. El sistema muestra las últimas acciones (fecha y hora de Bolivia, usuario, acción, detalle, IP).
3. El Director elige filtros y consulta.
4. El sistema devuelve la página filtrada.

## Flujos alternativos y excepciones

- **A1:** acción desconocida → 422 `ACCION_INVALIDA`.
- **A2:** «desde» posterior a «hasta» → 422 `RANGO_FECHAS_INVALIDO`.
- **A3:** sin resultados → lista vacía con mensaje.

## Reglas de negocio

- La tabla `auditoria` es de solo anexión: un trigger rechaza UPDATE, DELETE y TRUNCATE.
- Las fechas se interpretan en hora de Bolivia y ambos extremos son inclusivos.
- Nunca se registran contraseñas, tokens ni enlaces.

## Postcondiciones

- La consulta no modifica nada.

## Criterios de aceptación

- **CA01:** Dado un filtro por usuario y rango, el sistema debe devolver solo sus acciones en ese rango.
- **CA02:** Dado un rango invertido, el sistema debe responder `RANGO_FECHAS_INVALIDO`.
- **CA03:** Dado cualquier intento de modificar un registro en la base, el trigger debe rechazarlo.

## Unidades de trabajo

- [UT-BACKEND](../UnidadesTrabajo/CU-06%20-%20Consultar%20Registro%20de%20Auditor%C3%ADa/UT-BACKEND/README.md) — Hecho
- [UT-FRONTEND](../UnidadesTrabajo/CU-06%20-%20Consultar%20Registro%20de%20Auditor%C3%ADa/UT-FRONTEND/README.md) — Hecho

## Desarrollador

Carlos Otsubo — especificación y revisión; implementación asistida por IA (ver [BITACORA_IA.md](../BITACORA_IA.md)).
