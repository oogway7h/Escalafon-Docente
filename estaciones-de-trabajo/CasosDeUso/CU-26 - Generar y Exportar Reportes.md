# Caso de uso

| Id | Nombre | Paquete | Prioridad | Estado |
| --- | --- | --- | --- | --- |
| CU-26 | Generar y Exportar Reportes | P5 Indicadores | Media | Hecho (web) |

**Plataformas:** Web (Hecho) · Mobile (—)
**Actor principal:** Director (y Personal de Archivo por seed)
**Actores secundarios:** —
**Objetivo:** Generar reportes de trámites, inventario y tiempos en PDF con membrete o Excel filtrable.

**Descripción:**
El usuario elige el tipo de reporte (trámites, inventario de legajos o tiempos de atención), sus filtros y el formato. El sistema genera un PDF con membrete institucional, fecha de emisión, responsable, filtros, resumen y paginación, o un Excel con una hoja de datos con autofiltro y una de resumen. Cada emisión queda en auditoría; un conjunto vacío no genera archivo.

## Precondiciones

- Sesión con `REPORTE_GENERAR`.

## Flujo principal

1. El usuario abre Reportes.
2. Elige tipo, filtros y formato.
3. El sistema consulta, arma el reporte independiente del formato y lo compone en PDF o Excel.
4. Devuelve el archivo para descargar y audita `GENERAR_REPORTE`.

## Flujos alternativos y excepciones

- **A1:** sin filas → 422 `REPORTE_VACIO` (no se audita).
- **A2:** rango inválido → 422 `RANGO_FECHAS_INVALIDO`.
- **A3:** tipo o filtro desconocido → 422 `TIPO_REPORTE_INVALIDO` / `FILTRO_INVALIDO`.

## Reglas de negocio

- PDF A4 apaisado con logo y membrete en cada página y «Página X de N».
- Excel con fechas y números con su tipo.
- La auditoría guarda tipo, formato, filas y filtros.

## Postcondiciones

- El archivo se descargó y la emisión está en auditoría.

## Criterios de aceptación

- **CA01:** Dado un filtro con datos, el sistema debe devolver un PDF con membrete y paginación.
- **CA02:** Dado un filtro sin datos, el sistema debe responder `REPORTE_VACIO`.
- **CA03:** Cada emisión debe aparecer en la auditoría como «Generar reporte».

## Unidades de trabajo

- [UT-BACKEND](../UnidadesTrabajo/CU-26%20-%20Generar%20y%20Exportar%20Reportes/UT-BACKEND/README.md) — Hecho
- [UT-FRONTEND](../UnidadesTrabajo/CU-26%20-%20Generar%20y%20Exportar%20Reportes/UT-FRONTEND/README.md) — Hecho

## Desarrollador

Carlos Otsubo — especificación y revisión; implementación asistida por IA (ver [BITACORA_IA.md](../BITACORA_IA.md)).
