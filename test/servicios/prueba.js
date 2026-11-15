import http from 'k6/http';
import { check, sleep } from 'k6';
import { Counter } from 'k6/metrics';

export const options = {
  vus: 50,
  duration: '1m',
  thresholds: { http_req_duration: ['p(95)<1000'] },
};

const limitadas = new Counter('respuestas_429');

export default function () {
  // Consulta pública de un trámite por código de seguimiento (el endpoint que usa /seguimiento).
  const res = http.get('http://host.docker.internal:8081/api/publico/tramites/ED-2026-2SN2WX');
  if (res.status === 429) limitadas.add(1);
  check(res, { 'respuesta 200': (r) => r.status === 200 });
  sleep(1);
}
