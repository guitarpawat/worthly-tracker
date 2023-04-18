export function renderErrorInfo(resp) {
    let info = document.createElement('p')
    info.setAttribute('class', 'alert alert-danger')
    info.setAttribute('role', 'alert')
    if(!!resp.body && !!resp.body.message) {
        info.textContent = `Error with http status ${resp.status}: ${resp.body.message}`
    } else {
        info.textContent = `Error with http status ${resp.status}`
    }

    document.getElementById('error-info').appendChild(info)
}