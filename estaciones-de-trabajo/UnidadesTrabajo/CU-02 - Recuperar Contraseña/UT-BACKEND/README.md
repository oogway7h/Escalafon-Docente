# Unidades de trabajo — CU-02-BE

**Caso de uso:** Tokens de recuperación y activación de un solo uso, correo y revocación de sesiones.
**Proyecto:** `backend`
**Fuente CU:** [`CU-02 - Recuperar Contraseña`](../../../CasosDeUso/CU-02%20-%20Recuperar%20Contrase%C3%B1a.md)
**Recurso HTTP / contrato:** `POST /api/auth/recuperacion`, `/auth/restablecer`, `/auth/activar`
**Estado:** Hecho (documentación retroactiva del código existente)

## Antes de implementar cualquier UT

Leer [ARQUITECTURA_BACKEND.md](../../../.claude/rules/ARQUITECTURA_BACKEND.md) y [UNIDADES_DE_TRABAJO.md](../../../.claude/rules/UNIDADES_DE_TRABAJO.md). Estas fichas no sustituyen las reglas: si hay conflicto, ganan las reglas.

## Contrato de datos / recurso HTTP

| Método | Ruta | Éxito | Errores | Permiso |
|---|---|---|---|---|
| POST | /api/auth/recuperacion `{correo}` | 202 | 422 | público |
| POST | /api/auth/restablecer `{token, nueva}` | 204 | 400 `TOKEN_INVALIDO`, 422 `POLITICA_CONTRASENA` | público |
| POST | /api/auth/activar `{token, nueva}` | 204 | 400 `TOKEN_INVALIDO`, 409 `CUENTA_YA_ACTIVADA` | público |

## Orden de ejecución

| Código | Nombre | Depende de |
|--------|--------|------------|
| [UT-CU02-BE-01](./UT-CU02-BE-01.md) | Dominio de token y generador seguro | — |
| [UT-CU02-BE-02](./UT-CU02-BE-02.md) | SolicitarRecuperacion y RestablecerContrasena | 01 |
| [UT-CU02-BE-03](./UT-CU02-BE-03.md) | Handlers y rutas públicas | 02 |

**Cadena:** 01 → 02 → 03

## Cobertura del caso de uso

| Paso / CA | Descripción | UTs |
|-----------|-------------|-----|
| Flujo principal | Pedir enlace y restablecer | 01, 02, 03 |
| A1 | Respuesta neutra | 02 |
| A2, A3 | Token y política | 01, 02 |
| CA05 | Revocar sesiones | 02 |

## Follow-ups (no se implementan en estas UT)

- —

## Qué NO incluyen estas UT

- Cambio de contraseña con sesión iniciada (CU-03).
