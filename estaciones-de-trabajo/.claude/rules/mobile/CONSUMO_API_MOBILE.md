---
globs: ../ded_Mobile/**/*
---

# Consumo de API — Mobile

## Una sola URL: la API Go

La app habla con `API_URL` (por defecto `http://10.0.2.2:8080/api`, el host visto desde el emulador de
Android; en un teléfono real, la IP de la PC en la red local). Se cambia con
`flutter run --dart-define=API_URL=http://192.168.0.10:8080/api`. No hay `.env` en la app.

## Contrato de la API (distinto al de DriveHub)

- **Éxito:** el JSON viene directo, sin sobre (`{token, expiraEn, permisos, usuario}`,
  `{datos,total,pagina,tamano}`, …).
- **Error:** `{codigo, mensaje, detalles?}` con códigos estables (`CREDENCIALES_INVALIDAS`,
  `LEGAJO_NO_DISPONIBLE`, …). 401 = sesión vencida o revocada.
- **Sesión:** `Authorization: Bearer <token>`. El login devuelve el token en el cuerpo; no hay refresh
  token: al vencer (`expiraEn`, horas de `SESION_HORAS`) se vuelve a iniciar sesión.

## `ApiClient` (en `lib/core/network/api_client.dart`)

- Una instancia de `Dio` con `baseUrl = API_URL`, timeouts y un interceptor que agrega el Bearer
  leyendo un callback `Future<String?> Function() token` (core no importa `features/auth`).
- Toda `DioException` con respuesta JSON se convierte en `ApiException(status, codigo, mensaje)`; sin
  respuesta (servidor caído), `ApiException(0, 'SIN_CONEXION', 'No se pudo conectar con el servidor')`.
- Un 401 avisa por otro callback (`onNoAutorizado`) para que la sesión se borre y el router lleve al login.
- `mensajeDeError(Object e)` devuelve el texto en español para mostrar (usa `mensaje` del backend).

## Servicios y DTOs por feature

- `features/{feature}/services/*_service.dart` reciben el `Dio` por constructor (inyectable en tests) y
  solo conocen los endpoints de su feature.
- `features/{feature}/models/`: clases inmutables (`final` + `const`), `fromJson`/`toJson` escritos a
  mano, campos en camelCase igual que el JSON. Sin `json_serializable`.
- Solo el service construye y parsea DTOs; providers y pantallas reciben objetos tipados.

## Archivos (CU-13 y CU-15)

- CU-13 (portal): `multipart/form-data` a `/publico/solicitudes`, un archivo `doc_<requisitoId>` por
  requisito (PDF/JPG/PNG hasta 8 MB).
- CU-15 (analista): `POST …/documentos/subida` → URL firmada → `PUT` **con un `Dio` sin baseUrl ni
  Bearer** (S3 rechaza cabeceras extra) → `POST …/:docId/confirmar`.
