# Unidades de trabajo — CU-19-BE

**Caso de uso:** Inventario de legajos con código automático y búsqueda.
**Proyecto:** `backend`
**Fuente CU:** [`CU-19 - Gestionar Inventario de Legajos`](../../../CasosDeUso/CU-19%20-%20Gestionar%20Inventario%20de%20Legajos.md)
**Recurso HTTP / contrato:** `/api/legajos`
**Estado:** Hecho (documentación retroactiva del código existente)

## Antes de implementar cualquier UT

Leer [ARQUITECTURA_BACKEND.md](../../../.claude/rules/ARQUITECTURA_BACKEND.md) y [UNIDADES_DE_TRABAJO.md](../../../.claude/rules/UNIDADES_DE_TRABAJO.md). Estas fichas no sustituyen las reglas: si hay conflicto, ganan las reglas.

## Contrato de datos / recurso HTTP

| Método | Ruta | Éxito | Errores | Permiso |
|---|---|---|---|---|
| GET | /api/legajos?buscar=&estado=&pagina=&tamano= | 200 `{datos:[{id,codigo,docenteId,docenteCodigo,docenteNombre,gaveta,orden,anexos,observacion,estado}],…}` | — | `LEGAJO_GESTIONAR`, `CARPETA_SOLICITAR` o `CUSTODIA_REGISTRAR` |
| GET | /api/legajos/:id | 200 | 404 | idem |
| POST/PUT | /api/legajos[/:id] | 201/200 | 409 `LEGAJO_CODIGO_DUPLICADO` | `LEGAJO_GESTIONAR` |

## Orden de ejecución

| Código | Nombre | Depende de |
|--------|--------|------------|
| [UT-CU19-BE-01](./UT-CU19-BE-01.md) | Dominio legajo, migración y trigger de código | — |
| [UT-CU19-BE-02](./UT-CU19-BE-02.md) | Casos de uso ListarLegajos, ObtenerLegajo y GuardarLegajo | 01 |
| [UT-CU19-BE-03](./UT-CU19-BE-03.md) | Rutas de legajos | 02 |

**Cadena:** 01 → 02 → 03

## Cobertura del caso de uso

| Paso / CA | Descripción | UTs |
|-----------|-------------|-----|
| Flujo principal | Inventario | 01-03 |
| A1-A3 | Validación | 02 |

## Follow-ups (no se implementan en estas UT)

- —

## Qué NO incluyen estas UT

- Solicitudes y custodia (CU-20 a CU-24).
