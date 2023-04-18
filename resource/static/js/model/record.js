export class AssetRecord {
    id;
    assetId;
    name;
    broker;
    defaultIncrement;
    boughtValue;
    currentValue;
    realizedValue;
    note;
    unrealizedValue;
    unrealizedPercent;
    profitPercent;

    constructor(id, assetId, name, broker, defaultIncrement, boughtValue, currentValue, realizedValue, note, isCash) {
        this.id = id;
        this.assetId = assetId;
        this.name = name;
        this.broker = broker;
        this.defaultIncrement = Number(defaultIncrement);
        this.boughtValue = Number(boughtValue);
        this.currentValue = Number(currentValue);
        this.realizedValue = Number(realizedValue);

        if(isCash) {
            this.boughtValue = 0
            this.realizedValue = 0
        } else {
            this.unrealizedValue = currentValue - boughtValue
            this.unrealizedPercent = this.unrealizedValue / boughtValue
            this.profitPercent = (this.unrealizedValue + this.realizedValue) / boughtValue
        }

        if(!note) {
            this.note = ''
        } else {
            this.note = note
        }
    }

    static fromResp(asset, isCash) {
        return new AssetRecord(asset.id, asset.assetId, asset.name, asset.broker, asset.defaultIncrement,
            asset.boughtValue, asset.currentValue, asset.realizedValue, asset.note, isCash)
    }
}

export class AssetTypeRecord {
    id;
    name;
    isCash;
    isLiability;
    assets;

    constructor(id, name, isCash, isLiability, assets) {
        this.id = id;
        this.name = name;
        this.isCash = isCash;
        this.isLiability = isLiability;
        this.assets = assets;
    }

    static fromResp(type) {
        let assets = []
        for(let i = 0; i < type.assets.length; i++) {
            assets.push(AssetRecord.fromResp(type.assets[i], type.isCash))
        }
        return new AssetTypeRecord(type.id, type.name, type.isCash, type.isLiability, assets)
    }
}

export function fromRecordResponse(resp) {
    let types = []
    for(let i = 0; i < resp.length; i++) {
        types.push(AssetTypeRecord.fromResp(resp[i]))
    }

    return types
}

export class RecordsSummary {
    boughtValue = 0;
    currentValue = 0;
    unrealizedValue = 0;
    unrealizedPercent = 0;
    realizedValue = 0;
    profitPercent = 0;
    netWorth = 0;
    cash = 0;

    static fromResp(types) {
        let res = new RecordsSummary()

        for(let i = 0; i < types.length; i++) {
            for(let j = 0; j < types[i].assets.length; j++) {
                if(types[i].isCash) {
                    res.cash += Number(types[i].assets[j].currentValue)
                } else {
                    res.boughtValue += Number(types[i].assets[j].boughtValue)
                    res.currentValue += Number(types[i].assets[j].currentValue)
                    res.unrealizedValue += Number(types[i].assets[j].unrealizedValue)
                    res.realizedValue += Number(types[i].assets[j].realizedValue)
                }

                res.netWorth += Number(types[i].assets[j].currentValue)
            }
        }

        res.unrealizedValue = res.currentValue - res.boughtValue
        res.unrealizedPercent = res.unrealizedValue / res.boughtValue
        res.profitPercent = (res.unrealizedValue + Number(res.realizedValue)) / res.boughtValue

        return res
    }
}