# Unidades de trabajo — CU-19-MOB

**Caso de uso:** Buscar una carpeta física y ver su ubicación y estado (consulta).
**Proyecto:** `ded_Mobile`
**Fuente CU:** [`CU-19 - Gestionar Inventario de Legajos`](../../../CasosDeUso/CU-19%20-%20Gestionar%20Inventario%20de%20Legajos.md)
**Recurso HTTP / contrato:** `GET /api/legajos?buscar=&estado=&pagina=&tamano=`
**Estado:** Hecho

## Antes de implementar cualquier UT

Leer [ARQUITECTURA_MOBILE.md](../../../.claude/rules/ARQUITECTURA_MOBILE.md) (y `mobile/*`) y [UNIDADES_DE_TRABAJO.md](../../../.claude/rules/UNIDADES_DE_TRABAJO.md). Estas fichas no sustituyen las reglas: si hay conflicto, ganan las reglas.

## Decisiones de diseño

1. En el móvil el inventario es de **consulta**: dar de alta o mover carpetas se hace en la web.
2. Con dígitos la búsqueda es exacta (código de legajo, código o cédula del docente), igual que la API.

## Contrato de datos / recurso HTTP

| Método | Ruta | Éxito | Errores | Permiso |
|---|---|---|---|---|
| GET | /legajos?buscar=&estado=&pagina=1&tamano=20 | 200 `{datos:[{id,codigo,docenteId,docenteCodigo,docenteNombre,gaveta,orden,anexos,observacion,estado}],total,pagina,tamano}` | 403 | `LEGAJO_GESTIONAR`, `CARPETA_SOLICITAR` o `CUSTODIA_REGISTRAR` |

## Orden de ejecución

| Código | Nombre | Depende de |
|--------|--------|------------|
| [UT-CU19-MOB-01](./UT-CU19-MOB-01.md) | Modelo `Legajo` y `LegajosService` | — |
| [UT-CU19-MOB-02](./UT-CU19-MOB-02.md) | Pantalla de legajos | 01 |

**Cadena:** 01 → 02

## Cobertura del caso de uso

| Paso / CA | Descripción | UTs |
|-----------|-------------|-----|
| CA03 | Búsqueda exacta | 01, 02 |

## Follow-ups (no se implementan en estas UT)

- Escanear la etiqueta del legajo con la cámara (cuando se impriman etiquetas).

## Qué NO incluyen estas UT

- Alta y edición de legajos (web).
