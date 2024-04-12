import {HeaderModel} from '../model/header.js'
import {ApiFetcher} from '../common/fetcher.js'
import {renderErrorInfo} from '../common/error.js';

const pageNameAttr = 'page-name';

export class HeaderComponent extends HTMLElement {
    constructor() {
        super();

        let pageName = ''

        if(this.hasAttribute(pageNameAttr)){
            pageName = this.getAttribute(pageNameAttr);
        }

        (async()=>{
            await this.#getHeaderConfig(pageName)
        })();
    }

    async #getHeaderConfig(pageName) {
        let fetcher = new ApiFetcher()
        let resp = await fetcher.getHeader(pageName)
        if(!resp || resp.status !== 200) {
            renderErrorInfo(resp)
            return
        }

        // render nav items
        let navItems = HeaderModel.fromJson(resp.body)
        let navHtml = navItems.links.map(item => this.#renderTopLink(item)).join('\n');

        let innerHtml = `
<nav class="navbar navbar-expand-sm navbar-dark bg-secondary navbar-dar mb-3">
    <div class="container">
        <a class="navbar-brand" href="/">${navItems.title}</a>
        <button class="navbar-toggler" type="button" data-bs-toggle="collapse" data-bs-target="#navbar-menu" aria-controls="navbar-menu" aria-expanded="false" aria-label="Toggle navigation">
            <i class="navbar-toggler-icon"></i>
        </button>
        <div class="collapse navbar-collapse" id="navbar-menu">
            <ul class="navbar-nav ms-auto mb-2 mb-xs-0">
                ${navHtml}
            </ul>
        </div>
    </div>
</nav>
        `

        this.innerHTML = innerHtml
    }


    /**
     *
     * @param {TopLinkModel} link
     */
    #renderTopLink(link) {
        if(link.href === '#') {
            link.href = `#${link.href}" onclick="return false;`
        }
        if(link.childNodes && link.childNodes.length > 0) {
            let childHtml = link.childNodes.map(child => 
                `<a href="${child.href}" class="dropdown-item">${child.name}</a>`
            ).join('\n')

            let aClass = 'nav-link dropdown-toggle'
            if(link.highlight) {
                aClass += ' active'
            }

            return `
                <li class="nav-item dropdown">
                    <a href="${link.href}" class="${aClass}" role="button" data-bs-toggle="dropdown" aria-expanded="false">${link.name}</a>
                    <ul class="dropdown-menu dropdown-menu-end">
                        <li>
                            ${childHtml}
                        </li>
                    </ul>
                </li>
            `
        } else {
            let aClass = 'nav-link'
            if(link.highlight) {
                aClass += ' active'
            }

            return `
                <li class="nav-item">
                    <a href="${link.href}" class="${aClass}">${link.name}</a>
                </li>
            `
        }
    }
}