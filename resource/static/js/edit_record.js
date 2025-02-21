import * as formatter from './formatter.js'

window.onload = function () {
    Array.from(document.getElementsByClassName('data-enable-on-click')).forEach(e => {
        e.onclick = function() {
            Array.from(e.children).forEach(c => {
                c.disabled = false
                c.focus()
            })
        }
    })

    Array.from(document.getElementsByClassName('data-format-decimal')).forEach(e => {
        e.onchange = function() {
            e.value = formatter.formatDecimal(e.value)

            // we need to add here since when we create onchange event, it will replace the old one
            if (e.classList.contains('data-recalculate-profit')) {
                let id = e.id.replace('-bought', '').replace('-current', '')
                let boughtId = id + '-bought'
                let currentId = id + '-current'
                let profitId = id + '-profit'

                let boughtVal = Number(formatter.formatDecimal(document.getElementById(boughtId).value))
                let currentVal = Number(formatter.formatDecimal(document.getElementById(currentId).value))
                document.getElementById(profitId).innerText = formatter.formatPercent((currentVal - boughtVal) / boughtVal)
            }
        }
    })

    document.getElementById('change-date').onclick = function () {
        document.getElementById('date').showPicker()
    }

    document.getElementById('date').onchange = function () {
        document.getElementById('shown-date').innerText = document.getElementById('date').value
    }
}

window.onbeforeunload = function (e) {
    return "Are you sure to exit?"
}