export interface KeyValue {
  key: string;
  value: string;
}

export interface HttpResponse {
  body: string;
  status_code: number;
  id: string;
  local_id: number;
  headers: KeyValue[];
  tags: string[];
  body_size: number;
  cookies: KeyValue[];
  request: HttpRequest|null;
}

export interface HttpRequest {
  method: string;
  url: string;
  host: string;
  path: string;
  query_string: string;
  scheme: string;
  body: string;
  id: string;
  local_id: number;
  headers: KeyValue[];
  query: KeyValue[];
  tags: string[];
  response: HttpResponse|null;
}
