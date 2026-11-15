# UT-CU13-BE-03 — Rutas públicas de solicitud y categorías

| Campo | Valor |
|-------|-------|
| **Código** | `UT-CU13-BE-03` |
| **Nombre** | Rutas públicas de solicitud y categorías |
| **Caso de uso** | CU-13-BE |
| **Proyecto** | `backend` |
| **Depende de** | UT-CU13-BE-02 |
| **Estado** | Hecho (documentada de forma retroactiva) |

## Antes de implementar

[ARQUITECTURA_BACKEND.md](../../../.claude/rules/ARQUITECTURA_BACKEND.md)

## Solicitud de negocio

Exponer el formulario del portal.

## Objetivo

`POST /publico/solicitudes` (multipart) y `GET /publico/categorias`.

## Alcance

- Límite de cuerpo y limitador público.

Archivos esperados (crear/tocar):

- `backend/internal/interfaces/http/handlers/tramites.go`
- `backend/internal/interfaces/http/handlers/catalogo.go`
- `backend/internal/interfaces/http/router.go`

## Fuera de alcance (no tocar)

- Cualquier archivo de `backend` que no figure en «Archivos esperados», y las otras capas del CU (tienen sus propias UT).

## Resultado esperado

`POST /publico/solicitudes` (multipart) y `GET /publico/categorias`.

## Validación

1. `go test ./internal/interfaces/http -count=1` en verde.
2. `make arch` pasa (domain y application sin dependencias externas).
