import type { NextConfig } from "next";

const produccion = process.env.NODE_ENV === "production";

// CSP: Next inyecta scripts y estilos en línea (hidratación, tema), por eso 'unsafe-inline'; 'unsafe-eval' solo en
// desarrollo (recarga en caliente). connect-src admite S3 porque el navegador sube los documentos directo al bucket
// con una URL firmada; CSP_CONNECT_EXTRA agrega otro origen (p. ej. MinIO local).
const csp = [
  "default-src 'self'",
  `script-src 'self' 'unsafe-inline'${produccion ? "" : " 'unsafe-eval'"}`,
  "style-src 'self' 'unsafe-inline'",
  "img-src 'self' data: blob:",
  "font-src 'self' data:",
  `connect-src 'self' https://*.amazonaws.com ${process.env.CSP_CONNECT_EXTRA ?? ""}`.trim(),
  "worker-src 'self' blob:",
  "object-src 'none'",
  "base-uri 'self'",
  "form-action 'self'",
  "frame-ancestors 'none'",
].join("; ");

const cabeceras = [
  { key: "Content-Security-Policy", value: csp },
  { key: "X-Frame-Options", value: "DENY" },
  { key: "X-Content-Type-Options", value: "nosniff" },
  { key: "Referrer-Policy", value: "strict-origin-when-cross-origin" },
  { key: "Permissions-Policy", value: "camera=(), microphone=(), geolocation=()" },
  ...(produccion ? [{ key: "Strict-Transport-Security", value: "max-age=31536000; includeSubDomains" }] : []),
];

const nextConfig: NextConfig = {
  experimental: {
    // La solicitud del portal público envía hasta varios adjuntos de 8 MB (el backend los limita).
    serverActions: { bodySizeLimit: "40mb" },
  },
  async headers() {
    return [{ source: "/:path*", headers: cabeceras }];
  },
};

export default nextConfig;
