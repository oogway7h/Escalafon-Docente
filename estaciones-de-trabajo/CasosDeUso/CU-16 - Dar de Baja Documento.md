# Caso de uso

| Id | Nombre | Paquete | Prioridad | Estado |
| --- | --- | --- | --- | --- |
| CU-16 | Dar de Baja Documento | P3 Gestión documental | Media | Hecho (web) |

**Plataformas:** Web (Hecho) · Mobile (—)
**Actor principal:** Analista asignado / Personal de Archivo
**Actores secundarios:** —
**Objetivo:** Retirar de las búsquedas un documento cargado por error, con motivo, sin borrarlo.

**Descripción:**
El analista elige un documento vigente del trámite y escribe el motivo. El sistema lo marca inactivo con fecha, autor y motivo: deja de cubrir su requisito y desaparece de las búsquedas, pero el archivo se conserva y la operación queda en auditoría para que el Director pueda restituirlo.

## Precondiciones

- Documento vigente cargado ([CU-15](CU-15%20-%20Cargar%20y%20Procesar%20Documento%20con%20OCR.md)).
- Ser el analista asignado (o Archivo para escaneos) y el trámite no estar archivado.

## Flujo principal

1. El usuario elige «Dar de baja» sobre un documento.
2. Escribe el motivo (mínimo 5 caracteres).
3. El sistema marca el documento inactivo con motivo, fecha y autor, y audita `BAJA_DOCUMENTO`.
4. El documento se muestra tachado con su motivo.

## Flujos alternativos y excepciones

- **A1:** motivo vacío o corto → 422 `MOTIVO_BAJA_REQUERIDO`.
- **A2:** documento ya de baja → 409 `DOCUMENTO_DE_BAJA`.
- **A3:** trámite archivado → 409 `TRAMITE_ARCHIVADO`.
- **A4:** no es el asignado → 409 `NO_ASIGNADO`.

## Reglas de negocio

- Baja lógica: el objeto sigue en el almacenamiento.
- Un documento de baja solo lo descarga quien tiene `DOCUMENTO_RESTITUIR`.
- Si el trámite ya tiene dictamen, la pantalla advierte que hay que justificar la baja.

## Postcondiciones

- El documento está inactivo, fuera de las búsquedas, y la baja auditada.

## Criterios de aceptación

- **CA01:** Dado un motivo válido, el documento debe dejar de aparecer en las búsquedas.
- **CA02:** Dado un documento de baja, su requisito debe volver a figurar como faltante.
- **CA03:** Dado un trámite archivado, el sistema debe impedir la baja.

## Unidades de trabajo

- [UT-BACKEND](../UnidadesTrabajo/CU-16%20-%20Dar%20de%20Baja%20Documento/UT-BACKEND/README.md) — Hecho
- [UT-FRONTEND](../UnidadesTrabajo/CU-16%20-%20Dar%20de%20Baja%20Documento/UT-FRONTEND/README.md) — Hecho

## Desarrollador

Carlos Otsubo — especificación y revisión; implementación asistida por IA (ver [BITACORA_IA.md](../BITACORA_IA.md)).
