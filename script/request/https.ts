export declare const HttpsRequest:
    (method, serverURL, body: string, ...headers: string[])
        => {
        Status: string
        StatusCode: number
        Headers: Array<[string, string]>
        Body: string
    }
