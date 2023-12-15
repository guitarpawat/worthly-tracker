import * as formatter from '../common/formatter.js'
import {ApiFetcher} from '../common/fetcher.js'
import {renderErrorInfo} from '../common/error.js'
import {fromRecordResponse} from '../model/record.js'

let fetcher = new ApiFetcher()
let param = new URLSearchParams(window.location.search)
let drafts

let changeDate = function() {
    let newDate = document.getElementById('date').value
    document.getElementById('shown-date').innerText = formatter.formatDateAsFrontend(newDate)
}

let getDraft = async function() {
    let resp = await fetcher.getRecordDraft()
    if(resp.status === 200) {
        renderDraft(resp.body, true)
        document.getElementById('date').value = formatter.formatDateAsBackend(Date.now())
        changeDate()
    } else {
        renderErrorInfo(resp)
    }
}

let getRecord = async function() {
    let date = param.get('date')
    if(!date) {
        date = null
    }
    let resp = await fetcher.getRecordsByDate(date)
    if(resp.status === 200) {
        renderDraft(resp.body.types, false)
        document.getElementById('date').value = resp.body.date.current
        changeDate()
    } else {
        renderErrorInfo(resp)
    }
}

let renderDraft = function (resp, autoIncrement) {
    drafts = fromRecordResponse(resp)

    let table = document.getElementById('record-table')
    for(let i = 0; i < drafts.length; i++) {
        renderTypes(drafts[i], table, autoIncrement)
    }
}

let renderTypes = function (type, table, autoIncrement) {
    let tbody = document.createElement('tbody')
    table.appendChild(tbody)

    let tr = document.createElement('tr')
    tr.setAttribute('class', 'table-light')
    tbody.appendChild(tr)

    let th = document.createElement('th')
    th.setAttribute('class', 'text-center fs-6')
    th.setAttribute('colspan', '8')
    th.innerText = type.name
    tr.appendChild(th)

    for(let i = 0; i < type.assets.length; i++) {
        renderAssets(type.assets[i], type.isCash, tbody, autoIncrement)
    }
}

let renderAssets = function (asset, isCash, tbody, autoIncrement) {
    if(window.location.pathname.startsWith('/add')) {
        asset.id = null
    }

    let tr = document.createElement('tr')
    tbody.appendChild(tr)

    let th = document.createElement('th')
    th.setAttribute('class', 'text-start fs-6')
    th.setAttribute('scope', 'row')
    th.innerText = asset.name
    tr.appendChild(th)

    let tdBroker = document.createElement('td')
    tdBroker.setAttribute('class', 'text-start fs-6')
    tdBroker.innerText = asset.broker
    tr.appendChild(tdBroker)

    let tdAutoIncrement = document.createElement('td')
    tdAutoIncrement.setAttribute('class', 'text-center fs-6')
    if(asset.defaultIncrement) {
        tdAutoIncrement.innerText = '✔'
    }
    tr.appendChild(tdAutoIncrement)

    let tdBoughtValue = document.createElement('td')
    tr.appendChild(tdBoughtValue)
    let inputBoughtValue = document.createElement('input')
    inputBoughtValue.setAttribute('type', 'number')
    inputBoughtValue.setAttribute('step', '0.01')
    inputBoughtValue.setAttribute('class', 'fs-6 text-end')
    inputBoughtValue.disabled = true

    if(isCash) {
        inputBoughtValue.value = null
    } else {
        if(!asset.defaultIncrement || !autoIncrement) {
            asset.defaultIncrement = 0
        }
        if(!asset.boughtValue) {
            asset.boughtValue = 0
        }

        asset.boughtValue = asset.boughtValue + asset.defaultIncrement
        inputBoughtValue.value = formatter.formatDecimal(asset.boughtValue)
    }

    tdBoughtValue.appendChild(inputBoughtValue)

    let tdCurrentValue = document.createElement('td')
    tr.appendChild(tdCurrentValue)
    let inputCurrentValue = document.createElement('input')
    inputCurrentValue.setAttribute('type', 'number')
    inputCurrentValue.setAttribute('step', '0.01')
    inputCurrentValue.setAttribute('class', 'fs-6 text-end')
    if(!asset.currentValue) {
        asset.currentValue = "0"
    }
    inputCurrentValue.value = formatter.formatDecimal(asset.currentValue)
    tdCurrentValue.appendChild(inputCurrentValue)

    let tdUrPercent = document.createElement('td')
    tdUrPercent.setAttribute('class', 'fs-6 text-end')
    if(!isCash) {
        calUrPercent(tdUrPercent, inputBoughtValue, inputCurrentValue)
    }
    tr.appendChild(tdUrPercent)

    let tdRealizedValue = document.createElement('td')
    tr.appendChild(tdRealizedValue)
    let inputRealizedValue = document.createElement('input')
    inputRealizedValue.setAttribute('type', 'number')
    inputRealizedValue.setAttribute('step', '0.01')
    inputRealizedValue.setAttribute('class', 'fs-6 text-end')
    if(isCash) {
        inputRealizedValue.value = null
    } else if(asset.realizedValue) {
        inputRealizedValue.value = formatter.formatDecimal(asset.realizedValue)
    } else {
        asset.realizedValue = "0"
        inputRealizedValue.value = formatter.formatDecimal(0)
    }
    inputRealizedValue.disabled = true
    tdRealizedValue.appendChild(inputRealizedValue)

    let tdNote = document.createElement('td')
    tr.appendChild(tdNote)
    let inputNote = document.createElement('input')
    inputNote.setAttribute('type', 'text')
    inputNote.setAttribute('class', 'fs-6 text-start')
    inputNote.value = asset.note
    inputNote.disabled = true
    tdNote.onclick = function () {
        inputNote.disabled = false
        inputNote.focus()
    }
    inputNote.onchange = function () {asset.note = inputNote.value}
    tdNote.appendChild(inputNote)

    if(isCash) {
        inputCurrentValue.onchange = function () {
            inputCurrentValue.value = formatter.formatDecimal(inputCurrentValue.value)
            asset.currentValue = Number(inputCurrentValue.value)
        }
    } else {
        tdBoughtValue.onclick = function () {
            inputBoughtValue.disabled = false
            inputBoughtValue.focus()
        }
        inputBoughtValue.onchange = function () {
            inputBoughtValue.value = formatter.formatDecimal(inputBoughtValue.value)
            asset.boughtValue = Number(inputBoughtValue.value)
            calUrPercent(tdUrPercent, inputBoughtValue, inputCurrentValue)
        }
        inputCurrentValue.onchange = function () {
            inputCurrentValue.value = formatter.formatDecimal(inputCurrentValue.value)
            asset.currentValue = Number(inputCurrentValue.value)
            calUrPercent(tdUrPercent, inputBoughtValue, inputCurrentValue)
        }
        tdRealizedValue.onclick = function () {
            inputRealizedValue.disabled = false
            inputRealizedValue.focus()
        }
        inputRealizedValue.onchange = function () {
            inputRealizedValue.value = formatter.formatDecimal(inputRealizedValue.value)
            asset.realizedValue = Number(inputRealizedValue.value)
        }
    }
}

let calUrPercent = function (urPercent, bought, current) {
    urPercent.innerText = formatter.formatPercent((current.value - bought.value) / bought.value)
}

let postRecord = async function () {
    let records = []
    for(let i = 0; i < drafts.length; i++) {
        for(let j = 0; j < drafts[i].assets.length; j++) {
            records.push(drafts[i].assets[j])
        }
    }
    let resp = await fetcher.postRecord(records, document.getElementById('date').value)
    if(resp.status === 200) {
        window.onbeforeunload = null
        window.location.href = '/?date=' + document.getElementById('date').value
    } else {
        renderErrorInfo(resp)
    }
}

document.getElementById('change-date').onclick = function () {
    document.getElementById('date').showPicker()
}

document.getElementById('date').onchange = changeDate

document.getElementById('save-btn').onclick = await postRecord

window.onload = async function () {
    let path = window.location.pathname
    if(path.startsWith('/add')) {
        document.getElementById('record-action').innerText = 'Add Record'
        await getDraft()
        document.getElementById('save-btn').classList.remove('disabled')
    } else if(path.startsWith('/edit')) {
        document.getElementById('record-action').innerText = 'Edit Record'
        await getRecord()
        document.getElementById('save-btn').classList.remove('disabled')
    } else {
        document.getElementById('record-action').innerText = 'Something must be wrong with me'
    }
}

let onExit = function (e) {
    return "Are you sure to exit?"
}

window.onbeforeunload = onExit