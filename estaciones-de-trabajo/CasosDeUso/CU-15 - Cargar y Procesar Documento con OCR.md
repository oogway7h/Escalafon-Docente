# Caso de uso

| Id | Nombre | Paquete | Prioridad | Estado |
| --- | --- | --- | --- | --- |
| CU-15 | Cargar y Procesar Documento con OCR | P3 Gestión documental | Alta | Hecho (web) |

**Plataformas:** Web (Hecho) · Mobile (Hecho)
**Actor principal:** Analista asignado / Personal de Archivo
**Actores secundarios:** Servicio de OCR (`ia-service`), almacenamiento S3
**Objetivo:** Cargar y clasificar documentos, calcular su hash y extraer su texto por OCR para que sean buscables.

**Descripción:**
El analista elige un archivo (PDF, JPG o PNG) para el trámite y, si corresponde, el requisito que cubre. El archivo va directo del dispositivo al almacenamiento; al confirmar, el sistema verifica que llegó y lo procesa en segundo plano: calcula el SHA-256, extrae el texto por página y comprime el PDF. El usuario ve el estado (Pendiente, Procesando, Procesado o Error) y puede corregir el texto de cada página antes de darlo por validado. El Personal de Archivo hace lo mismo con los escaneos del expediente digital del docente, clasificados por tipo.

## Precondiciones

- Analista asignado al trámite no archivado, con `DOCUMENTO_GESTIONAR` (o Archivo con `DOCUMENTO_ESCANEAR` para el expediente).
- S3 e `ia-service` configurados.

## Flujo principal

1. El usuario elige el archivo y el requisito o tipo de documento.
2. El sistema valida tipo y tamaño y entrega una URL firmada de subida.
3. El dispositivo sube el archivo directo al almacenamiento.
4. El usuario confirma; el sistema verifica el objeto y lanza el procesamiento.
5. `ia-service` calcula el hash, extrae el texto por página y comprime el PDF.
6. El sistema guarda las páginas y marca el documento Procesado.
7. El usuario revisa y, si hace falta, corrige el texto por página.

## Flujos alternativos y excepciones

- **A1 (paso 2):** tipo distinto de PDF/JPG/PNG → 422 `TIPO_NO_PERMITIDO`.
- **A2 (paso 2):** imagen de más de 20 MB o PDF de más de 400 MB → 422 `ARCHIVO_GRANDE`.
- **A3 (paso 4):** el archivo no llegó → 409 `SUBIDA_INCOMPLETA`.
- **A4 (paso 5):** falla el OCR → estado Error con el motivo; confirmar de nuevo lo reintenta.
- **A5:** no es el analista asignado → 403/409 `NO_ASIGNADO`.
- **A6 (paso 7):** distinto número de páginas → 422 `PAGINAS_INVALIDAS`.

## Reglas de negocio

- Desvío de `cu.md`: 20 MB rige para imágenes; los PDF admiten hasta 400 MB.
- El archivo no pasa por la API: va por URL firmada.
- Estados OCR: SUBIENDO → PENDIENTE → PROCESANDO → PROCESADO | ERROR.
- Al arrancar, la API reanuda los trabajos pendientes.

## Postcondiciones

- El documento está almacenado con su hash y su texto por página, buscable (CU-18).

## Criterios de aceptación

- **CA01:** Dada una foto JPG de 3 MB, el sistema debe aceptarla y procesarla hasta Procesado.
- **CA02:** Dado un archivo de 25 MB JPG, el sistema debe rechazarlo con `ARCHIVO_GRANDE`.
- **CA03:** Dado que se confirma sin haber subido, el sistema debe responder `SUBIDA_INCOMPLETA`.
- **CA04:** Dado un OCR fallido, confirmar de nuevo debe reintentarlo.

## Unidades de trabajo

- [UT-BACKEND](../UnidadesTrabajo/CU-15%20-%20Cargar%20y%20Procesar%20Documento%20con%20OCR/UT-BACKEND/README.md) — Hecho
- [UT-FRONTEND](../UnidadesTrabajo/CU-15%20-%20Cargar%20y%20Procesar%20Documento%20con%20OCR/UT-FRONTEND/README.md) — Hecho
- [UT-MOBILE](../UnidadesTrabajo/CU-15%20-%20Cargar%20y%20Procesar%20Documento%20con%20OCR/UT-MOBILE/README.md) — Hecho

## Desarrollador

Carlos Otsubo — especificación y revisión; implementación asistida por IA (ver [BITACORA_IA.md](../BITACORA_IA.md)).
