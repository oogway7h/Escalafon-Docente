# Unidades de trabajo — CU-07-BE

**Caso de uso:** Padrón docente, facultades y catálogo de categorías con requisitos.
**Proyecto:** `backend`
**Fuente CU:** [`CU-07 - Gestionar Docentes y Categorías`](../../../CasosDeUso/CU-07%20-%20Gestionar%20Docentes%20y%20Categor%C3%ADas.md)
**Recurso HTTP / contrato:** `/api/docentes`, `/api/facultades`, `/api/categorias`
**Estado:** Hecho (documentación retroactiva del código existente)

## Antes de implementar cualquier UT

Leer [ARQUITECTURA_BACKEND.md](../../../.claude/rules/ARQUITECTURA_BACKEND.md) y [UNIDADES_DE_TRABAJO.md](../../../.claude/rules/UNIDADES_DE_TRABAJO.md). Estas fichas no sustituyen las reglas: si hay conflicto, ganan las reglas.

## Decisiones de diseño

1. Las categorías se comparan por clave normalizada (`persistence.ClaveCategoria`: sin tildes ni «Profesor»).
2. GORM ignora `false` en columnas `default:true`: las bajas se fuerzan con un `Update` aparte.

## Contrato de datos / recurso HTTP

| Método | Ruta | Éxito | Errores | Permiso |
|---|---|---|---|---|
| GET | /api/docentes?buscar=&pagina=&tamano= | 200 `{datos,total,pagina,tamano}` | — | lectura: `DOCENTE_GESTIONAR`, `TRAMITE_RADICAR`, `LEGAJO_GESTIONAR`… |
| POST/PUT | /api/docentes[/:id] | 201/200 | 409 `CODIGO_DUPLICADO`, `CEDULA_DUPLICADA` | `DOCENTE_GESTIONAR` |
| PATCH | /api/docentes/:id/estado | 200 | — | idem |
| GET/POST/PUT | /api/categorias[/:id] | 200/201 | 409 `CATEGORIA_DUPLICADA` | `CATEGORIA_GESTIONAR` |
| DELETE | /api/categorias/:id | 204 | 409 `CATEGORIA_EN_USO` | idem |

## Orden de ejecución

| Código | Nombre | Depende de |
|--------|--------|------------|
| [UT-CU07-BE-01](./UT-CU07-BE-01.md) | Dominio y persistencia de docente, facultad y categoría | — |
| [UT-CU07-BE-02](./UT-CU07-BE-02.md) | Casos de uso de docentes | 01 |
| [UT-CU07-BE-03](./UT-CU07-BE-03.md) | Casos de uso de categorías | 01 |
| [UT-CU07-BE-04](./UT-CU07-BE-04.md) | Handlers y rutas de catálogo | 02, 03 |

**Cadena:** 01 → 02 → 03 → 04

## Cobertura del caso de uso

| Paso / CA | Descripción | UTs |
|-----------|-------------|-----|
| Flujo principal | Docentes y categorías | 01-04 |
| A1-A3 | Unicidad e inmutabilidad | 02 |
| A4, A5 | Categorías | 03 |

## Follow-ups (no se implementan en estas UT)

- —

## Qué NO incluyen estas UT

- Carga masiva del padrón (`cmd/importdocentes`, fuera de CU).
- Legajos físicos (CU-19).
