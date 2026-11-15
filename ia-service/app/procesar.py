"""OCR (Tesseract) y compresión (ocrmypdf) de un documento guardado en S3."""
import hashlib
import os
import subprocess
import tempfile
from functools import lru_cache

import boto3
import ocrmypdf
import pytesseract
from PIL import Image

TIPOS = {"application/pdf", "image/jpeg", "image/png"}
SALTO_DE_PAGINA = "\x0c"  # pdftotext cierra cada página con un form feed


@lru_cache
def s3():
    return boto3.client(
        "s3",
        region_name=os.getenv("AWS_REGION", "us-east-1"),
        aws_access_key_id=os.getenv("AWS_ACCESS_KEY_ID"),
        aws_secret_access_key=os.getenv("AWS_SECRET_ACCESS_KEY"),
        endpoint_url=os.getenv("S3_ENDPOINT") or None,
    )


def sha256(ruta: str) -> str:
    h = hashlib.sha256()
    with open(ruta, "rb") as f:
        for bloque in iter(lambda: f.read(1 << 20), b""):
            h.update(bloque)
    return h.hexdigest()


def texto_de_imagen(ruta: str) -> str:
    with Image.open(ruta) as img:
        # El tipo lo declara quien sube el archivo: aquí se comprueba contra el contenido real.
        if img.format not in ("JPEG", "PNG"):
            raise ValueError("el archivo no es una imagen JPG o PNG válida")
        return pytesseract.image_to_string(img, lang="spa").strip()


def dividir_paginas(salida: str) -> list[str]:
    """Parte la salida de pdftotext en páginas (la última termina con su propio salto)."""
    paginas = salida.split(SALTO_DE_PAGINA)
    if paginas and not paginas[-1].strip():
        paginas.pop()
    return [p.strip() for p in paginas]


def paginas_de_pdf(ruta: str) -> list[str]:
    """Texto de cada página del PDF ya procesado (el mismo que verá el visor), con pdftotext."""
    salida = subprocess.run(["pdftotext", "-enc", "UTF-8", ruta, "-"], capture_output=True, check=True, timeout=1800)
    return dividir_paginas(salida.stdout.decode("utf-8", errors="replace"))


def procesar_pdf(entrada: str, salida: str) -> list[str]:
    """OCR y compresión en una pasada. `salida` queda con la capa de texto buscable; devuelve el texto por página."""
    ocrmypdf.ocr(
        entrada,
        salida,
        language=["spa"],
        skip_text=True,  # páginas que ya traen texto no se rasterizan de nuevo
        optimize=3,
        jpeg_quality=int(os.getenv("JPEG_QUALITY", "75")),
        oversample=300,  # evita comprimir por debajo de 300 dpi: el texto se sigue leyendo bien
        output_type="pdf",
        jobs=int(os.getenv("OCR_JOBS", "2")),
        progress_bar=False,
    )
    return paginas_de_pdf(salida)


def procesar(bucket: str, clave: str, tipo_mime: str, clave_salida: str | None) -> dict:
    """Descarga el objeto, extrae el texto por página y, en PDF, sube una versión comprimida a `clave_salida`.

    Todo pasa por disco (no por memoria): los PDF escaneados llegan a cientos de MB.
    """
    if tipo_mime not in TIPOS:
        raise ValueError("tipo de archivo no permitido")
    with tempfile.TemporaryDirectory() as tmp:
        entrada = os.path.join(tmp, "entrada")
        s3().download_file(bucket, clave, entrada)
        tamano = os.path.getsize(entrada)
        resultado = {"hashSha256": sha256(entrada), "tamanoOriginal": tamano, "tamanoFinal": tamano, "claveComprimida": ""}

        if tipo_mime != "application/pdf":
            texto = texto_de_imagen(entrada)
            resultado["paginas"], resultado["texto"] = [texto], texto
            return resultado

        with open(entrada, "rb") as f:
            if f.read(5) != b"%PDF-":
                raise ValueError("el archivo no es un PDF válido")
        salida = os.path.join(tmp, "salida.pdf")
        paginas = procesar_pdf(entrada, salida)
        resultado["paginas"], resultado["texto"] = paginas, "\n\n".join(paginas)

        reducido = os.path.getsize(salida)
        if clave_salida and reducido < tamano:  # si no reduce, se conserva solo el original
            s3().upload_file(salida, bucket, clave_salida, ExtraArgs={"ContentType": "application/pdf"})
            resultado["tamanoFinal"], resultado["claveComprimida"] = reducido, clave_salida
        return resultado
