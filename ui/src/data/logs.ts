export interface Log {
  id: string
  timestamp: string
  message: string
  method: string
  host: string
  url: string
  status: number
  size: number
  latency: string
  protocol: string
  userAgent?: string
  remoteIp: string
  referer: string
  serverIp: string
  responseSize: number
  severity: string
  spanId: string
  headers: {
    'Content-Type': string
    'Authorization': string
  }
}

export const logs: Log[] = [
  {
    id: '1f0f2c02-e299-40de-9b1d-86ef9e42126b',
    timestamp: '2024-09-22T09:00:00',
    message: 'service unavailable',
    method: 'GET',
    host: 'api.dev.ghostsecurity.com',
    url: '/v2/auth/login',
    status: 401,
    size: 367,
    latency: '0.002371508s',
    protocol: 'HTTP/1.1',
    remoteIp: '162.158.77.81',
    referer: 'https://app.dev.ghostsecurity.com/',
    serverIp: '172.253.62.121',
    responseSize: 144,
    severity: 'INFO',
    spanId: 'dc20888fe895a433',
    headers: {
      'Content-Type': 'application/json',
      'Authorization': 'Bearer 1234567890',
    },
  },
  {
    id: '1f0f2c02-e299-40de-9b1d-86ef9e42126c',
    timestamp: '2024-09-22T09:00:00',
    message: 'service unavailable',
    method: 'POST',
    host: 'api.dev.ghostsecurity.com',
    url: '/v2/auth/login',
    status: 401,
    size: 367,
    latency: '0.002371508s',
    protocol: 'HTTP/1.1',
    remoteIp: '162.158.77.81',
    referer: 'https://app.dev.ghostsecurity.com/',
    serverIp: '172.253.62.121',
    responseSize: 144,
    severity: 'WARN',
    spanId: 'dc20888fe895a433',
    headers: {
      'Content-Type': 'application/json',
      'Authorization': 'Bearer 1234567890',
    },
  },
  {
    id: '1f0f2c02-e299-40de-9b1d-86ef9e42126d',
    timestamp: '2024-09-22T09:00:00',
    message: 'service unavailable',
    method: 'POST',
    host: 'api.dev.ghostsecurity.com',
    url: '/v2/auth/login',
    status: 401,
    size: 367,
    latency: '0.002371508s',
    protocol: 'HTTP/1.1',
    userAgent: 'Mozilla/5.0 (compatible;Cloudflare-Healthchecks/1.0;+https://www.cloudflare.com/; healthcheck-id: 80beee5bf3635830)',
    remoteIp: '162.158.77.81',
    referer: 'https://app.dev.ghostsecurity.com/',
    serverIp: '172.253.62.121',
    responseSize: 144,
    severity: 'WARN',
    spanId: 'dc20888fe895a433',
    headers: {
      'Content-Type': 'application/json',
      'Authorization': 'Bearer 1234567890',
    },
  },
  {
    id: '1f0f2c02-e299-40de-9b1d-86ef9e42126e',
    timestamp: '2024-09-22T09:00:00',
    message: 'service unavailable',
    method: 'GET',
    host: 'api.dev.ghostsecurity.com',
    url: '/v2/auth/login',
    status: 401,
    size: 367,
    latency: '0.002371508s',
    protocol: 'HTTP/1.1',
    userAgent: 'Mozilla/5.0 (compatible;Cloudflare-Healthchecks/1.0;+https://www.cloudflare.com/; healthcheck-id: 80beee5bf3635830)',
    remoteIp: '162.158.77.81',
    referer: 'https://app.dev.ghostsecurity.com/',
    serverIp: '172.253.62.121',
    responseSize: 144,
    severity: 'WARN',
    spanId: 'dc20888fe895a433',
    headers: {
      'Content-Type': 'application/json',
      'Authorization': 'Bearer 1234567890',
    },
  },
  {
    id: '1f0f2c02-e299-40de-9b1d-86ef9e42126f',
    timestamp: '2024-09-22T09:00:00',
    message: 'service unavailable',
    method: 'GET',
    host: 'api.dev.ghostsecurity.com',
    url: '/v2/auth/login',
    status: 401,
    size: 367,
    latency: '0.002371508s',
    protocol: 'HTTP/1.1',
    userAgent: 'Mozilla/5.0 (compatible;Cloudflare-Healthchecks/1.0;+https://www.cloudflare.com/; healthcheck-id: 80beee5bf3635830)',
    remoteIp: '162.158.77.81',
    referer: 'https://app.dev.ghostsecurity.com/',
    serverIp: '172.253.62.121',
    responseSize: 144,
    severity: 'WARN',
    spanId: 'dc20888fe895a433',
    headers: {
      'Content-Type': 'application/json',
      'Authorization': 'Bearer 1234567890',
    },
  },
  {
    id: '1f0f2c02-e299-40de-9b1d-86ef9e42126g',
    timestamp: '2024-09-22T09:00:00',
    message: 'service unavailable',
    method: 'GET',
    host: 'api.dev.ghostsecurity.com',
    url: '/v2/auth/login',
    status: 401,
    size: 367,
    latency: '0.002371508s',
    protocol: 'HTTP/1.1',
    userAgent: 'Mozilla/5.0 (compatible;Cloudflare-Healthchecks/1.0;+https://www.cloudflare.com/; healthcheck-id: 80beee5bf3635830)',
    remoteIp: '162.158.77.81',
    referer: 'https://app.dev.ghostsecurity.com/',
    serverIp: '172.253.62.121',
    responseSize: 144,
    severity: 'WARN',
    spanId: 'dc20888fe895a433',
    headers: {
      'Content-Type': 'application/json',
      'Authorization': 'Bearer 1234567890',
    },
  },

]
