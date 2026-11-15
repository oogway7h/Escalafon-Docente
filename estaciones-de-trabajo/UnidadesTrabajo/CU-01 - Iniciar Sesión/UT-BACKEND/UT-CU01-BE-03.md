# UT-CU01-BE-03 — Caso de uso AutenticarUsuario

| Campo | Valor |
|-------|-------|
| **Código** | `UT-CU01-BE-03` |
| **Nombre** | Caso de uso AutenticarUsuario |
| **Caso de uso** | CU-01-BE |
| **Proyecto** | `backend` |
| **Depende de** | UT-CU01-BE-01, UT-CU01-BE-02 |
| **Estado** | Hecho (documentada de forma retroactiva) |

## Antes de implementar

[ARQUITECTURA_BACKEND.md](../../../.claude/rules/ARQUITECTURA_BACKEND.md)

## Solicitud de negocio

Centralizar la regla de login: validar, contar intentos, bloquear, emitir y auditar.

## Objetivo

`AutenticarUsuario.Ejecutar` devuelve token, vencimiento, usuario y permisos.

## Alcance

- Busca por correo (sin distinguir mayúsculas), verifica estado y bloqueo, compara hash.
- Fallo: registra intento y audita `LOGIN_FALLIDO` con IP.
- Éxito: reinicia intentos, emite token y audita `LOGIN`.
- Tests con dobles: éxito, contraseña errónea, tercer fallo bloquea, cuenta inactiva.

Archivos esperados (crear/tocar):

- `backend/internal/application/auth/auth.go`
- `backend/internal/application/auth/auth_test.go`

## Fuera de alcance (no tocar)

- Cerrar sesión (UT 04).
- Recuperación (CU-02).
- Cualquier archivo de `backend` que no figure en «Archivos esperados», y las otras capas del CU (tienen sus propias UT).

## Resultado esperado

`AutenticarUsuario.Ejecutar` devuelve token, vencimiento, usuario y permisos.

## Validación

1. `go test ./internal/application/auth -count=1` en verde.
2. `make arch` pasa (domain y application sin dependencias externas).
