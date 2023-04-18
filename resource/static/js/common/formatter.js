let currencyFormatter = new Intl.NumberFormat('en-US', {
    style: 'currency',
    currency: 'THB',
    currencySign: 'accounting',
})

let percentageFormatter = new Intl.NumberFormat('en-US', {
    style: 'percent',
    minimumFractionDigits: 2,
    maximumFractionDigits: 2,
})

let fixedDecimalFormatter = new Intl.NumberFormat('en-US', {
    style: 'decimal',
    minimumFractionDigits: 2,
    maximumFractionDigits: 2,
})

let frontendDateFormatter = new Intl.DateTimeFormat('en-GB', {
    day: "2-digit",
    month: "2-digit",
    year: "numeric"
})

let backendDateFormatter = new Intl.DateTimeFormat('sv-SE', {
    day: "2-digit",
    month: "2-digit",
    year: "numeric"
})

export function formatCurrency(val) {
    if(val !== 0 && !val) {
        return ''
    }
    if(val < 0) {
        return `${currencyFormatter.format(val)}`
    } else {
        return `${currencyFormatter.format(val)}\u00A0`
    }
}

export function formatPercent(val) {
    if(val !== 0 && !val) {
        return ''
    }
    return percentageFormatter.format(val)
}

export function formatDecimal(val) {
    if(val !== 0 && !val) {
        return ''
    }
    return fixedDecimalFormatter.format(val).replaceAll(',', '')
}

export function formatDateAsFrontend(date) {
    let d
    if(date instanceof Date) {
        d = date
    } else {
        d = new Date(date)
    }
    return frontendDateFormatter.format(d)
}

export function formatDateAsBackend(date) {
    let d
    if(date instanceof Date) {
        d = date
    } else {
        d = new Date(date)
    }
    return backendDateFormatter.format(d)
}