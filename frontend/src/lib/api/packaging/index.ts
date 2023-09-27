
export interface KeyValue {
Key: string
Value: string
}


export interface HttpRequest {
Method: string
URL: string
Host: string
Path: string
QueryString: string
Scheme: string
Body: string
ID: string
LocalID: number
Headers: KeyValue[]
Query: KeyValue[]
Tags: string[]
}

