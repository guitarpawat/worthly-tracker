import {HeaderComponent} from '../component/header.js';

const headerTag = 'worthly-header'

export function loadAllCustomTags() {
    if(document.getElementsByTagName(headerTag)) {
        customElements.define(headerTag, HeaderComponent);
    }
}