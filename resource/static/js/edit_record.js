import * as formatter from './formatter'

function enableField(htmlId) {
    document.getElementById(htmlId).disabled = false
}

export function formatDecimalInput(htmlId) {
    document.getElementById(htmlId).innerText = formatter.formatDecimal(document.getElementById(htmlId).innerText)
}

export function calculatePercent(assetId) {
    let boughtValue = Number(formatter.formatDecimal(document.getElementById(assetId+"-bought").innerText))
    let currentValue = Number(formatter.formatDecimal(document.getElementById(assetId+"-current").innerText))
    document.getElementById(assetId+"-profit").innerText = formatter.formatPercent((currentValue - boughtValue) / boughtValue)
}

window.onbeforeunload = function (e) {
    return "Are you sure to exit?"
}