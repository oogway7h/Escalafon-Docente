FROM node:24-alpine
RUN npm install -g pnpm@11.24.0
WORKDIR /app
COPY package.json pnpm-lock.yaml pnpm-workspace.yaml ./
RUN pnpm install --frozen-lockfile
COPY . .
# La CSP se calcula al compilar: CSP_CONNECT_EXTRA (origen extra para subir documentos, p. ej. MinIO) va como build arg.
ARG CSP_CONNECT_EXTRA=""
ENV CSP_CONNECT_EXTRA=$CSP_CONNECT_EXTRA
# El build comprueba que exista la variable (no consulta la API); en ejecución la fija el compose.
ENV API_URL_INTERNA=http://backend:8080
RUN pnpm build && chown -R node:node /app
ENV NODE_ENV=production
USER node
EXPOSE 3000
CMD ["pnpm", "start", "-H", "0.0.0.0"]
