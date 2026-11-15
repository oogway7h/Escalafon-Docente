# Caso de uso

| Id | Nombre | Paquete | Prioridad | Estado |
| --- | --- | --- | --- | --- |
| CU-17 | Restituir Documento Dado de Baja | P3 Gestión documental | Baja | Hecho (web) |

**Plataformas:** Web (Hecho) · Mobile (—)
**Actor principal:** Director
**Actores secundarios:** —
**Objetivo:** Reactivar un documento dado de baja y devolverlo a las búsquedas.

**Descripción:**
El Director ve la bandeja de documentos dados de baja (docente, trámite, motivo, quién y cuándo), elige uno y lo restituye escribiendo un motivo. El documento vuelve a estar activo y a aparecer en las búsquedas sin reprocesar el OCR. La operación queda en auditoría con ambos motivos. Es exclusiva del Director.

## Precondiciones

- Sesión con `DOCUMENTO_RESTITUIR`.
- Existe un documento dado de baja ([CU-16](CU-16%20-%20Dar%20de%20Baja%20Documento.md)).

## Flujo principal

1. El Director abre «Documentos dados de baja».
2. Elige un documento y ve su detalle.
3. Escribe el motivo de la restitución y confirma.
4. El sistema reactiva el documento con un UPDATE condicional y audita `RESTITUIR_DOCUMENTO`.

## Flujos alternativos y excepciones

- **A1:** motivo de menos de 5 caracteres → 422 `MOTIVO_RESTITUCION_REQUERIDO`.
- **A2:** el documento ya está activo → 409 `DOCUMENTO_ACTIVO`.
- **A3:** trámite remitido o archivado → 409 `TRAMITE_REMITIDO` / `TRAMITE_ARCHIVADO`.

## Reglas de negocio

- Solo el Director.
- No se reprocesa el OCR: las búsquedas filtran por `activo`.

## Postcondiciones

- El documento está activo y vuelve a las búsquedas.

## Criterios de aceptación

- **CA01:** Dado un documento de baja y un motivo válido, debe volver a aparecer en las búsquedas.
- **CA02:** Dado un usuario sin `DOCUMENTO_RESTITUIR`, la ruta debe responder 403.
- **CA03:** Dado un trámite archivado, el sistema debe impedir restituir.

## Unidades de trabajo

- [UT-BACKEND](../UnidadesTrabajo/CU-17%20-%20Restituir%20Documento%20Dado%20de%20Baja/UT-BACKEND/README.md) — Hecho
- [UT-FRONTEND](../UnidadesTrabajo/CU-17%20-%20Restituir%20Documento%20Dado%20de%20Baja/UT-FRONTEND/README.md) — Hecho

## Desarrollador

Carlos Otsubo — especificación y revisión; implementación asistida por IA (ver [BITACORA_IA.md](../BITACORA_IA.md)).
