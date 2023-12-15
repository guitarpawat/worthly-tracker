import {ApiFetcher} from '../common/fetcher.js'
import {renderErrorInfo} from '../common/error.js'
import {fromRecordResponse, RecordsSummary} from '../model/record.js'
import * as formatter from '../common/formatter.js'
import {BigNumber} from '../bignumber.mjs'

let fetcher = new ApiFetcher()
let param = new URLSearchParams(window.location.search)

let getRecordOnLoad = async function () {
    let date = param.get('date')
    if(!date) {
        date = null
    }
    await fetcher.getRecordsByDate(null)
    await loadRecord(date)
    editBtnTarget()
}

let disableBtn = function (resp) {
    let btns = document.getElementsByClassName('disable-on-error')
    for (let i = 0; i < btns.length; i++) {
        btns[i].disabled = true
        btns[i].classList.add('disabled')
    }
}

let enableBtn = function (resp) {
    let btns = document.getElementsByClassName('disable-on-error')
    for (let i = 0; i < btns.length; i++) {
        btns[i].disabled = false
        btns[i].classList.remove('disabled')
    }
}

let renderResponse = function (resp) {
    let records = fromRecordResponse(resp.body.types)

    let table = document.getElementById('record-table')
    for(let i = 0; i < records.length; i++) {
        let tbody = document.createElement('tbody')
        createRecordHead(tbody, records[i].name)
        for(let j = 0; j < records[i].assets.length; j++) {
            createAssetRecord(tbody, records[i].assets[j])
        }
        table.appendChild(tbody)
    }

    createRecordFoot(table, records)
}

let createRecordHead = function(tbody, name) {
    let tr = document.createElement('tr')
    tr.setAttribute('class', 'table-light')

    let thName = document.createElement('th')
    thName.setAttribute('class', 'text-center border-end-0')
    thName.setAttribute('colspan', '8')
    thName.setAttribute('scope', 'row')
    thName.innerText = name
    tr.appendChild(thName)

    let tdNote = document.createElement('td')
    tdNote.setAttribute('class', 'text-start border-start-0')
    tr.appendChild(tdNote)

    tbody.appendChild(tr)
}

let createAssetRecord = function (tbody, record) {
    let tr = document.createElement('tr')

    let thName = document.createElement('th')
    thName.setAttribute('class', 'text-start')
    thName.setAttribute('scope', 'row')
    thName.innerText = record.name
    tr.appendChild(thName)

    let tdPlace = document.createElement('td')
    tdPlace.setAttribute('class', 'text-start')
    tdPlace.innerText = record.broker
    tr.appendChild(tdPlace)

    let tdBoughtValue = document.createElement('td')
    tdBoughtValue.setAttribute('class', 'text-end')
    if(record.boughtValue) {
        tdBoughtValue.innerText = formatter.formatCurrency(record.boughtValue)
    }
    tr.appendChild(tdBoughtValue)

    let tdCurrentValue = document.createElement('td')
    tdCurrentValue.setAttribute('class', 'text-end')
    tdCurrentValue.innerText = formatter.formatCurrency(record.currentValue)
    tr.appendChild(tdCurrentValue)

    let tdUnrealizedValue = document.createElement('td')
    tdUnrealizedValue.setAttribute('class', 'text-end')
    tdUnrealizedValue.innerText = formatter.formatCurrency(record.unrealizedValue)
    tr.appendChild(tdUnrealizedValue)

    let tdUnrealizedPercent = document.createElement('td')
    tdUnrealizedPercent.setAttribute('class', 'text-end')
    tdUnrealizedPercent.innerText = formatter.formatPercent(record.unrealizedPercent)
    tr.appendChild(tdUnrealizedPercent)

    let tdRealizedValue = document.createElement('td')
    tdRealizedValue.setAttribute('class', 'text-end')
    if(record.realizedValue) {
        tdRealizedValue.innerText = formatter.formatCurrency(record.realizedValue)
    }
    tr.appendChild(tdRealizedValue)

    let tdProfitPercent = document.createElement('td')
    tdProfitPercent.setAttribute('class', 'text-end')
    tdProfitPercent.innerText = formatter.formatPercent(record.profitPercent)
    tr.appendChild(tdProfitPercent)

    let tdNote = document.createElement('td')
    tdNote.setAttribute('class', 'text-start')
    tdNote.innerText = record.note
    tr.appendChild(tdNote)

    tbody.appendChild(tr)
}

let createRecordFoot = function (table, types) {
    let result = RecordsSummary.fromResp(types)

    let tfoot = document.createElement('tfoot')
    tfoot.setAttribute('class', 'table-light table-group-divider')

    let trInvestment = document.createElement('tr')

    let thName = document.createElement('th')
    thName.setAttribute('class', 'text-center')
    thName.setAttribute('scope', 'rowgroup')
    thName.setAttribute('rowspan', '2')
    thName.innerText = 'Total'
    trInvestment.appendChild(thName)

    let thPlace = document.createElement('th')
    thPlace.setAttribute('class', 'text-center')
    thPlace.setAttribute('scope', 'row')
    thPlace.innerText = 'Investment'
    trInvestment.appendChild(thPlace)

    let tdBoughtValue = document.createElement('td')
    tdBoughtValue.setAttribute('class', 'text-end')
    tdBoughtValue.innerText = formatter.formatCurrency(result.boughtValue)
    trInvestment.appendChild(tdBoughtValue)

    let tdCurrentValue = document.createElement('td')
    tdCurrentValue.setAttribute('class', 'text-end')
    tdCurrentValue.innerText = formatter.formatCurrency(result.currentValue)
    trInvestment.appendChild(tdCurrentValue)

    let tdUnrealizedValue = document.createElement('td')
    tdUnrealizedValue.setAttribute('class', 'text-end')
    tdUnrealizedValue.innerText = formatter.formatCurrency(result.unrealizedValue)
    trInvestment.appendChild(tdUnrealizedValue)

    let tdUnrealizedPercent = document.createElement('td')
    tdUnrealizedPercent.setAttribute('class', 'text-end')
    tdUnrealizedPercent.innerText = formatter.formatPercent(result.unrealizedPercent)
    trInvestment.appendChild(tdUnrealizedPercent)

    let tdRealizedValue = document.createElement('td')
    tdRealizedValue.setAttribute('class', 'text-end')
    tdRealizedValue.innerText = formatter.formatCurrency(result.realizedValue)
    trInvestment.appendChild(tdRealizedValue)

    let tdProfitPercent = document.createElement('td')
    tdProfitPercent.setAttribute('class', 'text-end')
    tdProfitPercent.innerText = formatter.formatPercent(result.profitPercent)
    trInvestment.appendChild(tdProfitPercent)

    let tdNote = document.createElement('td')
    tdNote.setAttribute('class', 'text-start')
    tdNote.setAttribute('rowspan', '2')
    tdNote.innerText = `Net Worth: ${formatter.formatCurrency(result.netWorth)}\nCash/NW: ${formatter.formatPercent(result.cash/result.netWorth)}`
    trInvestment.appendChild(tdNote)

    tfoot.appendChild(trInvestment)

    let trCash = document.createElement('tr')

    let thPlace2 = document.createElement('th')
    thPlace2.setAttribute('class', 'text-center')
    thPlace2.setAttribute('scope', 'row')
    thPlace2.innerText = 'Cash'
    trCash.appendChild(thPlace2)

    let tdBoughtValue2 = document.createElement('td')
    tdBoughtValue2.setAttribute('class', 'text-end')
    trCash.appendChild(tdBoughtValue2)

    let tdCurrentValue2 = document.createElement('td')
    tdCurrentValue2.setAttribute('class', 'text-end')
    tdCurrentValue2.innerText = formatter.formatCurrency(result.cash)
    trCash.appendChild(tdCurrentValue2)

    let tdUnrealizedValue2 = document.createElement('td')
    tdUnrealizedValue2.setAttribute('class', 'text-end')
    trCash.appendChild(tdUnrealizedValue2)

    let tdUnrealizedPercent2 = document.createElement('td')
    tdUnrealizedPercent2.setAttribute('class', 'text-end')
    trCash.appendChild(tdUnrealizedPercent2)

    let tdRealizedValue2 = document.createElement('td')
    tdRealizedValue2.setAttribute('class', 'text-end')
    trCash.appendChild(tdRealizedValue2)

    let tdProfitPercent2 = document.createElement('td')
    tdProfitPercent2.setAttribute('class', 'text-end')
    trCash.appendChild(tdProfitPercent2)

    tfoot.appendChild(trCash)

    table.appendChild(tfoot)
}

let renderDate = function(resp) {
    let dates = resp.body.date
    let dateSelector = document.getElementById('date-selector')

    if(dates.prev) {
        for(let i = dates.prev.length - 1; i >= 0; i--) {
            createDateOption(dateSelector, dates.prev[i])
        }
        document.getElementById('prev-btn').value = dates.prev[0]
    } else {
        document.getElementById('prev-btn').removeAttribute('value')
    }

    createDateOption(dateSelector, dates.current, true)

    if(dates.next) {
        for(let i = 0; i < dates.next.length; i++) {
            createDateOption(dateSelector, dates.next[i])
        }
        document.getElementById('next-btn').value = dates.next[0]
    } else {
        document.getElementById('next-btn').removeAttribute('value')
    }
}

let createDateOption = function (dateSelector, date, isSelected = false) {
    let option = document.createElement('option')
    option.setAttribute('value', date)
    option.selected = isSelected
    option.innerText = formatter.formatDateAsFrontend(date)
    dateSelector.appendChild(option)
}

let loadRecord = async function (date) {
    let record = await fetcher.getRecordsByDate(date)
    let offset = await fetcher.getBoughtValueOffsetByDate(record.body.date.current)
    processOffset(record, offset)

    if(record.status === 200) {
        renderResponse(record)
        renderDate(record)
        enableBtn()
    } else {
        renderErrorInfo(record)
    }
}

let processOffset = function (record, offset) {
    if(offset && offset.body) {
        record.body.types.forEach(assetType => {
            assetType.assets.forEach(asset => {
                let assetOffset = offset.body.find(o => o.assetId === asset.assetId)
                if (assetOffset) {
                    asset.boughtValue = new BigNumber(asset.boughtValue).plus(new BigNumber(assetOffset.offsetPrice)).toString()
                }
            })
        })
    }
}

let clearData = function () {
    document.getElementById('date-selector').innerHTML = ''
    let table = document.getElementById('record-table')
    while (table.lastChild.tagName !== 'THEAD') {
        table.lastChild.remove()
    }
}

let editBtnTarget = function () {
    document.getElementById('edit-btn').href = '/edit' + window.location.search
}

let changeRecord = async function(date) {
    clearData()
    disableBtn()
    await loadRecord(date)
    if(date) {
        param.set('date', date)
        window.history.pushState(date, '', '?' + param.toString())
    } else {
        window.history.pushState(date, '', '/')
    }
    editBtnTarget()
}

window.onload = getRecordOnLoad
document.getElementById('date-selector').onchange = async function () {await changeRecord(document.getElementById('date-selector').value)}
document.getElementById('prev-btn').onclick = async function () {
    let date = document.getElementById('prev-btn').value
    if(!date) return
    await changeRecord(date)
}
document.getElementById('next-btn').onclick = async function () {
    let date = document.getElementById('next-btn').value
    if(!date) return
    await changeRecord(date)
}

document.getElementById('latest-btn').onclick = async function () {
    await changeRecord(null)
}