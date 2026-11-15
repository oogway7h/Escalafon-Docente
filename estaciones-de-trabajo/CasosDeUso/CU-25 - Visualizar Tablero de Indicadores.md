# Caso de uso

| Id | Nombre | Paquete | Prioridad | Estado |
| --- | --- | --- | --- | --- |
| CU-25 | Visualizar Tablero de Indicadores | P5 Indicadores | Media | Hecho (web) |

**Plataformas:** Web (Hecho) · Mobile (Hecho)
**Actor principal:** Director
**Actores secundarios:** —
**Objetivo:** Ver indicadores consolidados y alertas de carpetas retenidas fuera de plazo.

**Descripción:**
El Director abre el tablero y ve, para un periodo opcional: trámites radicados, en curso, resueltos y tiempo promedio de resolución; la distribución por estado; la carga vigente por analista; el tiempo de resolución por tipo y categoría; y las alertas de carpetas en mora (legajo, docente, custodio, salida, límite, días de mora). Las alertas desaparecen solas al registrarse la devolución. El tablero se actualiza periódicamente.

## Precondiciones

- Sesión con `TABLERO_VER`.

## Flujo principal

1. El Director abre el tablero.
2. Elige un periodo (o todo el historial).
3. El sistema calcula los indicadores sobre trámites, historial y custodia.
4. Muestra indicadores, gráficos de barras y la tabla de alertas.
5. Se refresca cada 5 minutos.

## Flujos alternativos y excepciones

- **A1:** «desde» posterior a «hasta» → 422 `RANGO_FECHAS_INVALIDO`.
- **A2:** fecha con formato inválido → 422 `VALIDACION`.
- **A3:** sin carpetas en mora → la tabla de alertas no aparece.

## Reglas de negocio

- Resolución de un trámite = su primer paso a Remitido.
- Mora: pasó la fecha límite o lleva más de `DIAS_MAXIMOS_PRESTAMO` (15) días fuera.
- La carga y las alertas son siempre las de hoy, sin importar el periodo.
- Solo lectura: no hay tablas propias.

## Postcondiciones

- La consulta no modifica nada.

## Criterios de aceptación

- **CA01:** Dado un periodo, los radicados deben contar solo los trámites de ese rango.
- **CA02:** Dada una carpeta con fecha límite vencida, debe aparecer en alertas con sus días de mora.
- **CA03:** Dada la devolución de esa carpeta, la alerta debe desaparecer en la siguiente consulta.

## Unidades de trabajo

- [UT-BACKEND](../UnidadesTrabajo/CU-25%20-%20Visualizar%20Tablero%20de%20Indicadores/UT-BACKEND/README.md) — Hecho
- [UT-FRONTEND](../UnidadesTrabajo/CU-25%20-%20Visualizar%20Tablero%20de%20Indicadores/UT-FRONTEND/README.md) — Hecho
- [UT-MOBILE](../UnidadesTrabajo/CU-25%20-%20Visualizar%20Tablero%20de%20Indicadores/UT-MOBILE/README.md) — Hecho

## Desarrollador

Carlos Otsubo — especificación y revisión; implementación asistida por IA (ver [BITACORA_IA.md](../BITACORA_IA.md)).
