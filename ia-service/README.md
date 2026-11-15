# ia-service

Servicio FastAPI de OCR (Tesseract, español) y compresión de PDF (ocrmypdf) para CU-15. El backend (Go) le pide un trabajo y consulta su avance; **el archivo no viaja por HTTP**: el servicio lo descarga de S3 en streaming, trabaja sobre disco temporal y sube la versión comprimida junto al original.

## Ejecutar

```bash
cp .env.example .env        # token compartido con backend/.env (IA_SERVICE_TOKEN) y credenciales de S3
docker compose up -d --build ia-service   # desde la raíz del repo; escucha en http://localhost:8000
```

Sin Docker: instalar Tesseract (con el idioma `spa`), Ghostscript y pngquant, y `pip install -r requirements.txt && uvicorn app.main:app`.

## API (cabecera `X-Token` obligatoria salvo en `/salud`)

| Método | Ruta | Descripción |
| --- | --- | --- |
| `POST` | `/trabajos` | `{bucket, clave, tipoMime, claveSalida?}` → 202 `{id}`. Solo acepta el bucket configurado y claves `docentes/…`. |
| `GET` | `/trabajos/{id}` | `{estado: EN_COLA\|PROCESANDO\|LISTO\|ERROR, texto, hashSha256, tamanoOriginal, tamanoFinal, claveComprimida, error}` |
| `GET` | `/salud` | Estado del servicio. |

PDF: `ocrmypdf -l spa --skip-text -O3 --oversample 300` (JPEG a calidad `JPEG_QUALITY`, 75 por defecto). La versión comprimida solo se sube si pesa menos que el original. Imágenes: solo OCR.

## Límites conocidos

- Los trabajos viven en memoria: si el servicio reinicia, el backend vuelve a lanzar los pendientes al arrancar.
- Un PDF de 400 MB puede tardar decenas de minutos; `MAX_TRABAJOS` (1) y `OCR_JOBS` (2) acotan la CPU y la RAM. El disco temporal necesita varias veces el tamaño del PDF.
- Para comprimir aún más, instalar `jbig2enc` en la imagen y usar `--jbig2-lossy` (texto en blanco y negro).

## CORS del bucket (hace falta una vez, en la consola de AWS)

El navegador sube con `PUT` directo a S3, así que el bucket necesita (Permisos → CORS):

```json
[{"AllowedOrigins": ["http://localhost:3000"], "AllowedMethods": ["PUT"], "AllowedHeaders": ["Content-Type"], "ExposeHeaders": [], "MaxAgeSeconds": 3000}]
```

Agrega el origen de producción del frontend cuando exista. El usuario IAM `ded-backend` no tiene `s3:PutBucketCORS`.

## Pruebas

`pip install -r requirements-dev.txt && pytest` (la prueba de OCR se omite si `tesseract` no está instalado).
