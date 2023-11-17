export function renderErrorInfo(resp) {
    let info = document.createElement('div')
    info.setAttribute('class', 'alert alert-danger alert-dismissible fade show')
    info.setAttribute('role', 'alert')
    if(!!resp.body && !!resp.body.message) {
        info.textContent = `Error with http status ${resp.status}: ${resp.body.message}`
    } else {
        info.textContent = `Error with http status ${resp.status}`
    }

    info.insertAdjacentHTML('beforeend', '<button type="button" class="btn-close" data-bs-dismiss="alert"></button>')

    document.getElementById('error-info').innerHTML = ''
    document.getElementById('error-info').appendChild(info)
}