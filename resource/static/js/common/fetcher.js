class ApiConfig {
    root = "http://localhost:8080/api"
    getRecord = "/records/:date"
    getRecordDraft = "/records/draft"
    postRecord = "/records/"
}

export class ApiFetcher {
    config = new ApiConfig()

    async #get(url) {
        let response = await window.fetch(url, {
            method: 'GET',
        })
        let json = await response.json()
        return {
            body: json,
            status: response.status
        }
    }

    async #post(url, body) {
        let response = await window.fetch(url, {
            method: 'POST',
            body: body,
            headers: {
                "Content-Type": "application/json",
            },
        })
        let json = await response.json()
        return {
            body: json,
            status: response.status
        }
    }

    async getRecordsByDate(date) {
        if (!date) {
            date = ""
        }

        let url = this.config.root + this.config.getRecord.replace(":date", date)
        return this.#get(url)
    }

    async getRecordDraft() {
        let url = this.config.root + this.config.getRecordDraft
        return this.#get(url)
    }

    async postRecord(records, date) {
        let body = {
            "date": date,
            "assets": []
        }

        for(let i = 0; i < records.length; i++) {
            body.assets.push({
                "id": records[i].id,
                "assetId": records[i].assetId,
                "boughtValue": records[i].boughtValue,
                "currentValue": records[i].currentValue,
                "realizedValue": records[i].realizedValue,
                "note": records[i].note,
            })
        }

        let url = this.config.root + this.config.postRecord
        return await this.#post(url, JSON.stringify(body))
    }
}