interface HttpRequest {
  ID: string
  LocalID: number
  Method: string
  Scheme: string
  URL: string
  Host: string
  Path: string
  QueryString: string
  Body: string
  Response: HttpResponse | null
  Headers: { [key: string]: Array<string> }
  Query: { [key: string]: Array<string> }
  Tags: Array<string>
}

interface HttpResponse {
  ID: string
  LocalID: number
  Body: string
  StatusCode: number
  Request: HttpRequest | null
  Headers: { [key: string]: Array<string> }
  Tags: Array<string>
  BodySize: number
}

export function MethodClass(req: HttpRequest): string {
  switch (req.Method) {
  case 'GET':
    return 'bg-frost-3'
  case 'OPTIONS':
    return 'bg-frost-1'
  case 'PUT':
    return 'bg-frost-4'
  case 'HEAD':
    return 'bg-frost-2'
  case 'DELETE':
    return 'bg-aurora-4'
  case 'CONNECT':
    return 'bg-aurora-3'
  case 'POST':
    return 'bg-aurora-2'
  case 'TRACE':
    return 'bg-aurora-1'
  case 'PATCH':
    return 'bg-aurora-4'
  default:
    return ''
  }
}

export function StatusClass(req: HttpRequest): string {
  if (req.Response === null) {
    return ''
  }
  const status = req.Response?.StatusCode
  if (status === undefined) {
    return 'bg-polar-night-1'
  }
  if (status >= 200 && status < 300) {
    return 'bg-aurora-4'
  }
  if (status >= 300 && status < 400) {
    return 'bg-aurora-3'
  }
  if (status >= 400 && status < 500) {
    return 'bg-aurora-2'
  }
  if (status >= 500 && status < 600) {
    return 'bg-aurora-1'
  }
  return ''
}

export {}

export type { HttpRequest, HttpResponse }
