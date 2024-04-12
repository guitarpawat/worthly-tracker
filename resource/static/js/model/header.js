export class HeaderModel {
    title = ''
    /** @type {TopLinkModel[]} */
    links = []

    constructor(title, links) {
        this.title = title
        this.links = links
    }

    static fromJson(obj) {
        let links = []
        if(obj.links) {
            obj.links.forEach(link => {
                links.push(TopLinkModel.fromJson(link));
            })
        }

        return new HeaderModel(
            obj.title,
            links,
        )
    }
}

class LinkModel {
    name = ''
    href = ''

    constructor(name, href) {
        this.name = name;
        this.href = href;
    }

    static fromJson(obj) {
        return new LinkModel(
            obj.name,
            obj.href,
        )
    }
}

class TopLinkModel extends LinkModel {
    highlight = false
    /** @type {LinkModel[]} */
    childNodes;

    constructor(name, href, highlight = false, childNodes = []) {
        super(name, href);
        this.highlight = highlight;
        this.childNodes = childNodes;
    }

    static fromJson(obj) {
        let childNodes = []
        if(obj.childNodes) {
            obj.childNodes.forEach(child => {
                childNodes.push(LinkModel.fromJson(child));
            })
        }

        return new TopLinkModel(
            obj.name,
            obj.href,
            obj.highlight,
            childNodes,
        )
    }
}