import {loadAllCustomTags} from './custom_tags.js';

// BEGIN EXPORT ALL DEPENDENCIES
export {BigNumber} from 'https://cdn.jsdelivr.net/npm/bignumber.js@9/bignumber.mjs'
// END EXPORT ALL DEPENDENCIES

export function init() {
    dynamicallyLoadScript('https://cdn.jsdelivr.net/npm/bootstrap@5/dist/js/bootstrap.bundle.min.js')
    loadAllCustomTags()
}

function dynamicallyLoadScript(url) {
    let script = document.createElement("script");
    script.src = url;
    document.head.appendChild(script);
}
