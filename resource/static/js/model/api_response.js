export class ApiResponse {
    body
    status

    constructor(json, status) {
        this.body = json
        this.status = status
    }
}