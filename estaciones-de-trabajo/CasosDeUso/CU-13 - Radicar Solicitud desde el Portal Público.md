# Caso de uso

| Id | Nombre | Paquete | Prioridad | Estado |
| --- | --- | --- | --- | --- |
| CU-13 | Radicar Solicitud desde el Portal Público | P2 Trámites | Alta | Hecho (web) |

**Plataformas:** Web (Hecho) · Mobile (Hecho)
**Actor principal:** Docente (sin cuenta)
**Actores secundarios:** Personal de Ventanilla (valida)
**Objetivo:** Iniciar una solicitud adjuntando los documentos que exige la categoría, sin ir a ventanilla.

**Descripción:**
El docente se identifica con su código o cédula, elige el tipo de trámite y la categoría; el sistema lista los documentos que esa categoría exige y el docente adjunta un archivo por cada uno (PDF, JPG o PNG). El sistema valida identidad, tipos y tamaños, guarda los archivos con su hash, genera el código de seguimiento y crea el trámite pendiente de validación en ventanilla.

## Precondiciones

- El docente está en el padrón y activo.
- La categoría está vigente y tiene requisitos ([CU-07](CU-07%20-%20Gestionar%20Docentes%20y%20Categor%C3%ADas.md)).
- El almacenamiento de archivos (S3) está configurado.

## Flujo principal

1. El docente abre «Nueva solicitud».
2. Escribe su código o cédula, elige tipo y categoría.
3. El sistema muestra los documentos exigidos.
4. El docente adjunta un archivo por documento y envía.
5. El sistema valida, guarda los archivos, genera el código y crea el trámite Ingresado pendiente de validación.
6. El sistema muestra el código y lleva al seguimiento.

## Flujos alternativos y excepciones

- **A1:** docente no encontrado o inactivo → 422 `DOCENTE_INVALIDO` / 409 `DOCENTE_INACTIVO`.
- **A2:** falta un documento exigido → 422 `DOCUMENTOS_FALTANTES`.
- **A3:** archivo que no es PDF/JPG/PNG por contenido → 422 `ARCHIVO_INVALIDO`.
- **A4:** archivo de más de 8 MB → 422 `ARCHIVO_GRANDE`.
- **A5:** almacenamiento no configurado → 503 `ALMACEN_NO_DISPONIBLE`.

## Reglas de negocio

- El tipo de archivo se valida por contenido, no por extensión.
- Cada archivo se guarda con su SHA-256 en `docentes/<código>/<trámite>/…`.
- La solicitud no puede asignarse hasta que ventanilla la valide (CU-08).

## Postcondiciones

- Existe un trámite Ingresado con origen PORTAL, pendiente de validación, con sus adjuntos.

## Criterios de aceptación

- **CA01:** Dada una solicitud completa, el sistema debe devolver un código de seguimiento.
- **CA02:** Dado que falta un documento exigido, el sistema debe rechazarla con `DOCUMENTOS_FALTANTES`.
- **CA03:** Dado un archivo de 9 MB, el sistema debe rechazarlo con `ARCHIVO_GRANDE`.
- **CA04:** Dada una solicitud radicada, ventanilla debe verla como pendiente de validación.

## Unidades de trabajo

- [UT-BACKEND](../UnidadesTrabajo/CU-13%20-%20Radicar%20Solicitud%20desde%20el%20Portal%20P%C3%BAblico/UT-BACKEND/README.md) — Hecho
- [UT-FRONTEND](../UnidadesTrabajo/CU-13%20-%20Radicar%20Solicitud%20desde%20el%20Portal%20P%C3%BAblico/UT-FRONTEND/README.md) — Hecho
- [UT-MOBILE](../UnidadesTrabajo/CU-13%20-%20Radicar%20Solicitud%20desde%20el%20Portal%20P%C3%BAblico/UT-MOBILE/README.md) — Hecho

## Desarrollador

Carlos Otsubo — especificación y revisión; implementación asistida por IA (ver [BITACORA_IA.md](../BITACORA_IA.md)).
