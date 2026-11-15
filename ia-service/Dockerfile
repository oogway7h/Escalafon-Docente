FROM python:3.12-slim

# Tesseract (con español) + Ghostscript/pngquant/unpaper para que ocrmypdf pueda optimizar los PDF.
RUN apt-get update && apt-get install -y --no-install-recommends \
      tesseract-ocr tesseract-ocr-spa ghostscript pngquant unpaper poppler-utils \
    && rm -rf /var/lib/apt/lists/*

WORKDIR /srv
COPY requirements.txt .
RUN pip install --no-cache-dir -r requirements.txt
COPY app ./app

# Sin privilegios: el servicio procesa PDF e imágenes que no son de confianza (Ghostscript, Tesseract).
RUN useradd --create-home --uid 10001 ocr
USER ocr

EXPOSE 8000
CMD ["uvicorn", "app.main:app", "--host", "0.0.0.0", "--port", "8000"]
