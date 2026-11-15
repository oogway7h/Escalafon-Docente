# Caso de uso

| Id | Nombre | Paquete | Prioridad | Estado |
| --- | --- | --- | --- | --- |
| CU-11 | Emitir Dictamen y Archivar Trámite | P2 Trámites | Alta | Hecho (web) |

**Plataformas:** Web (Hecho) · Mobile (—)
**Actor principal:** Analista asignado
**Actores secundarios:** Personal de Archivo (archiva), Personal de Ventanilla (registra el regreso)
**Objetivo:** Emitir el dictamen, remitir el expediente por el tramo de facultad y Rectorado y archivarlo como solo lectura.

**Descripción:**
Sin observaciones pendientes y con los años a evaluar registrados, el analista emite el dictamen (favorable o desfavorable, con texto): el trámite pasa a Dictaminado. Lo remite a la facultad (En facultad, tramo sin seguimiento). Cuando vuelve con la resolución (número y fecha), el analista o ventanilla registra el regreso y el trámite pasa a Revisión final con el mismo analista, que lo remite al Rectorado (Remitido). Finalmente el Personal de Archivo lo archiva indicando la ubicación del legajo, solo si todas sus carpetas volvieron; queda Archivado y de solo lectura.

## Precondiciones

- Trámite En revisión del analista asignado ([CU-10](CU-10%20-%20Revisar%20Expediente%20y%20Registrar%20Observaciones.md)).

## Flujo principal

1. El analista registra los años a evaluar.
2. Emite el dictamen con resultado y texto; el trámite pasa a Dictaminado.
3. Remite a la facultad con una nota; pasa a En facultad.
4. Al regresar, el analista o ventanilla registra número y fecha de resolución; pasa a Revisión final.
5. El analista remite al Rectorado; pasa a Remitido.
6. El Personal de Archivo archiva con la ubicación del legajo; pasa a Archivado.

## Flujos alternativos y excepciones

- **A1 (paso 2):** observaciones pendientes → 409 `OBSERVACIONES_PENDIENTES`.
- **A2 (paso 2):** sin años a evaluar → 409 `PERIODO_REQUERIDO`; rango inválido → 422 `PERIODO_INVALIDO`.
- **A3 (paso 4):** sin número o fecha → 422 `RESOLUCION_REQUERIDA`.
- **A4 (paso 6):** carpetas sin devolver → 409 `CARPETAS_ABIERTAS`.
- **A5:** modificar un trámite remitido o archivado → 409 `TRAMITE_REMITIDO` / `TRAMITE_ARCHIVADO`.

## Reglas de negocio

- Desvío de `cu.md`: el analista dictamina y remite; **archiva el Personal de Archivo** (`TRAMITE_ARCHIVAR`).
- El dictamen es interno: el docente ve el Resultado recién desde Remitido.
- En facultad es un tramo sin seguimiento (Tribunal, Decano, Consejo Facultativo).
- Archivado es de solo lectura.

## Postcondiciones

- El trámite está Archivado con fecha y ubicación del legajo.

## Criterios de aceptación

- **CA01:** Dado un trámite con observaciones pendientes, el sistema debe impedir dictaminar.
- **CA02:** Dado un dictamen válido, el trámite debe pasar a Dictaminado.
- **CA03:** Dado el regreso con resolución, el trámite debe volver al mismo analista en Revisión final.
- **CA04:** Dado un trámite con carpetas prestadas, el sistema debe impedir archivarlo.
- **CA05:** Dado un trámite archivado, el sistema debe rechazar cualquier modificación.

## Unidades de trabajo

- [UT-BACKEND](../UnidadesTrabajo/CU-11%20-%20Emitir%20Dictamen%20y%20Archivar%20Tr%C3%A1mite/UT-BACKEND/README.md) — Hecho
- [UT-FRONTEND](../UnidadesTrabajo/CU-11%20-%20Emitir%20Dictamen%20y%20Archivar%20Tr%C3%A1mite/UT-FRONTEND/README.md) — Hecho

## Desarrollador

Carlos Otsubo — especificación y revisión; implementación asistida por IA (ver [BITACORA_IA.md](../BITACORA_IA.md)).
