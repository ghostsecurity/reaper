
interface HttpRequest {
    ID: number
    Method: string
    Scheme: string
    URL: string
    Host: string
    Path: string
    QueryString: string
    Raw: string
    Response: HttpResponse | null
    Headers: { [key: string]: Array<string> }
    Query: { [key: string]: Array<string> }
}

interface HttpResponse {
    ID: number
    Raw: string
    StatusCode: number
    Request: HttpRequest | null
    Headers: { [key: string]: Array<string> }
}

export type {
    HttpRequest,
    HttpResponse
}