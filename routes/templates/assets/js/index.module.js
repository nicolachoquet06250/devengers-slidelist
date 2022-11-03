/**
 * @typedef {object} Owner
 * @property {string} displayName
 * @property {string} photoLink
 */

/**
 * @typedef {object} File
 * @property {string} id
 * @property {string} name
 * @property {string} webContentLink
 * @property {string} webViewLink
 * @property {string} iconLink
 * @property {boolean} hasThumbnail
 * @property {string} thumbnailLink
 * @property {boolean} trashed
 * @property {string} createdTime
 * @property {string} modifiedTime
 * @property {Array<Owner>} owners
 * @property {Owner} lastModifyingUser
 * @property {boolean} shared
 * @property {string} originalFilename
 * @property {boolean} size
 */

/**
 * @typedef {object} GoogleDriveAPIResponse
 * @property {string} kind
 * @property {Array<File>} files
 */

Object.prototype.getDeepProperty = function(prop) {
    const [firstKey, ...lastKeys] = prop.split('.');

    if (this[firstKey]) {
        if (lastKeys.length >= 1) {
            return this[firstKey].getDeepProperty(lastKeys.join('.'))
        }
        return this[firstKey];
    }

    return null;
};

/**
 * @param {string} format
 * @return {string}
 */
Date.prototype.format = function (format) {
    const assoc = {
        'DD': () => {
            const d = this.getDate()

            return `${d < 10 ? '0' : ''}${d}`;
        },
        'MM': () => {
            const m = this.getMonth() + 1

            return `${m < 10 ? '0' : ''}${m}`;
        },
        'D': this.getDate.bind(this),
        'M': () => this.getMonth() + 1,
        'YYYY': this.getFullYear.bind(this),
        'YY': () => {
            const y = this.getFullYear();

            return y.toString(10).substring(2, 4);
        },
        'HH': () => {
            const h = this.getHours();

            return `${h < 10 ? '0' : ''}${h}`
        },
        'II': () => {
            const m = this.getMinutes();

            return `${m < 10 ? '0' : ''}${m}`
        },
        'SS': () => {
            const s = this.getSeconds();

            return `${s < 10 ? '0' : ''}${s}`
        },
        'H': this.getHours.bind(this),
        'I': this.getMinutes.bind(this),
        'S': this.getSeconds.bind(this),
    };

    const regex = new RegExp(`(${Object.keys(assoc).join('|')})([\\:\\\/\\|\\(\\)\\[\\]\\#\\ ]+)?`, 'gm');

    let m;

    /**
     * @type { {format: string, separator: ?string}[] }
     */
    const formattedArray = [];

    while ((m = regex.exec(format)) !== null) {
        if (m.index === regex.lastIndex) {
            regex.lastIndex++;
        }

        formattedArray.push({
            format: m[1],
            separator: m[2]
        })
    }

    return formattedArray.map(({ format, separator }) => `${assoc[format]()}${separator ?? ''}`).join('');
}

const app = document.querySelector('#app');
let removeLoginLink = () => {};

function openWin(e) {
    e.preventDefault();

    const url = e.target.href;

    const win = window.open(url, '', 'width=400,height=400,top=200,left=200');

    (function detectWinClose() {
        if (win.closed) {
            app.innerHTML = '';

            getDevengersElements();

            return;
        }

        setTimeout(detectWinClose, 100);
    })();
}

function getAsyncAccessToken() {
    const notReallyConnected =
        (localStorage.getItem('loggedIn') === '0' || localStorage.getItem('loggedIn') === null) &&
        localStorage.getItem('googleOAuthEndTimestamp') !== null &&
        (Date.now() / 1000) < parseInt(localStorage.getItem('googleOAuthEndTimestamp'))

    if (notReallyConnected) {
        localStorage.setItem('loggedIn', '1');
    }

    const askAccessToken =
        (
            localStorage.getItem('googleOAuthEndTimestamp') === null ||
            (Date.now() / 1000) >= parseInt(localStorage.getItem('googleOAuthEndTimestamp'))
        ) &&
        localStorage.getItem('googleOAuthToken') !== null

    if (askAccessToken) {
        localStorage.setItem('loggedIn', '0');

        const { ClientID: client_id, ClientSecret: client_secret, RedirectURI: redirect_uri } = window.params;

        return fetch('https://oauth2.googleapis.com/token', {
            method: 'post',
            headers: {
                Accept: 'applicaiton/json',
                'Content-Type': 'application/json',
            },
            body: JSON.stringify({
                grant_type: 'authorization_code',
                code: localStorage.getItem('googleOAuthToken'),
                client_id, client_secret, redirect_uri
            })
        })
            .then(async r => {
                if (!r.ok) {
                    throw new Error(JSON.stringify(await r.json()))
                }
                return r.json()
            })
            .then(json => {
                const {access_token, expires_in, refresh_token, token_type} = json;

                localStorage.setItem('googleOAuthAccessToken', access_token);
                if (refresh_token) {
                    localStorage.setItem('googleOAuthToken', refresh_token);
                }
                localStorage.setItem('googleOAuthExpiresIn', expires_in);
                localStorage.setItem('googleOAuthTokenType', token_type);
                localStorage.setItem('googleOAuthEndTimestamp', (Date.now() / 1000) + expires_in);

                localStorage.setItem('loggedIn', '1');

                return access_token;
            })
    }

    return Promise.resolve(localStorage.getItem('googleOAuthAccessToken'));
}

function getDevengersElements() {
    const { ApiKey } = window.params;

    return getAsyncAccessToken()
        .then(access_token => {
            if (localStorage.getItem('loggedIn') === '1') {
                return fetch('https://www.googleapis.com/drive/v3/files?q=%271ZKwJPKXIKXI5YdSwV3mICWsC2eUq-kj0%27%20in%20parents&fields=*&key='+ApiKey, {
                    method: 'get',
                    headers: {
                        Authorization: `${localStorage.getItem('googleOAuthTokenType')} ${access_token}`,
                        Accept: 'application/json'
                    }
                })
            }
        })
        .then(async r => {
            if (!r.ok) {
                throw new Error(JSON.stringify(await r.json()))
            }
            return r.json()
        })
        .then(
            /**
             * @param {GoogleDriveAPIResponse} json
             */
            json => {
                removeLoginLink();
                createDOMPptxList(json.files);
            }
        )
}

/**
 * @param {Array<File>} files
 */
function createDOMPptxList(files) {
    const container = document.createElement('div');
    container.classList.add('container');

    files = files.filter(f => !f.trashed)

    for (const file of files) {
        /**
         * @const
         * @type {HTMLElement}
         */
        const pptxTpl = document.querySelector('#pptx-card').cloneNode(true).content.firstElementChild;

        Array.from(pptxTpl.querySelectorAll('[data]')).map(e => {
            e.getAttributeNames().filter(a => a.startsWith('data-')).map(a => {
                const prop = a
                    .replace('data-', '')
                    .split('-')
                    .map((p, i) => (i === 0 ? p.substring(0, 1) : p.substring(0, 1).toUpperCase()) + p.substring(1, p.length).toLowerCase())
                    .join('');

                e[prop] = file.getDeepProperty(e.getAttribute(a));

                if (e.hasAttribute('date-format')) {
                    e[prop] = new Date(file.getDeepProperty(e.getAttribute(a))).format(e.getAttribute('date-format').toUpperCase())
                }

                if (e.hasAttribute('format-data') && tools !== null && tools[e.getAttribute('format-data').toString()]) {
                    const func = e.getAttribute('format-data').toString();
                    e[prop] = tools[func](e[prop]);
                }
            })
        });

        container.appendChild(pptxTpl);
    }

    app.innerHTML = '';
    app.appendChild(container);
}

const tools = {
    formatSize(size) {
        size = parseInt(size);

        return Math.round(size / 1000) + 'Ko';
    }
};

function createLoginLink() {
    const { LinkUrl, LinkLabel, Referer } = window.params;

    const a = document.createElement('a');
    a.href = `${LinkUrl}&referrer=${Referer}`.replace(/"/g, '');
    a.target = '_blank';
    a.innerText = LinkLabel.replace(/"/g, '');
    a.classList.add('login');

    const disableClick = e => e.preventDefault();

    const init = () => {
        console.log('link init');
        a.removeAttribute('disabled', '');
        a.removeEventListener('click', disableClick);
        a.addEventListener('click', openWin);
    };
    const destroy = () => {
        console.log('link destroyed');
        a.setAttribute('disabled', '');
        a.removeEventListener('click', openWin);
        a.addEventListener('click', disableClick);
    };

    window.addEventListener('offline', destroy);
    window.addEventListener('online', init);

    app.innerHTML = '';
    app.appendChild(a);

    init();

    return destroy;
}

export default function () {
    getDevengersElements().catch(err => {
        console.log(err)

        removeLoginLink = createLoginLink();
    })
};
