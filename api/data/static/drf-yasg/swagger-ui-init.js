"use strict";
var currentPath = window.location.protocol + "//" + window.location.host + window.location.pathname;
var specURL = currentPath + '?format=openapi';

function patchSwaggerUi() {
    var authWrapper = document.querySelector('.auth-wrapper');
    var authorizeButton = document.querySelector('.auth-wrapper .authorize');
    var djangoSessionAuth = document.querySelector('#django-session-auth');
    if (!djangoSessionAuth) {
        console.log("WARNING: session auth disabled");
        return;
    }

    if (document.querySelector('.auth-wrapper #django-session-auth')) {
        console.log("WARNING: session auth already patched; skipping patchSwaggerUi()");
        return;
    }

    authWrapper.insertBefore(djangoSessionAuth, authorizeButton);
    djangoSessionAuth.classList.remove("hidden");

    var divider = document.createElement("div");
    divider.classList.add("divider");
    authWrapper.insertBefore(divider, authorizeButton);
}

function initSwaggerUi() {
    if (window.ui) {
        console.log("WARNING: skipping initSwaggerUi() because window.ui is already defined");
        return;
    }
    var swaggerConfig = {
        url: specURL,
        dom_id: '#swagger-ui',
        displayRequestDuration: true,
        presets: [
            SwaggerUIBundle.presets.apis,
            SwaggerUIStandalonePreset
        ],
        plugins: [
            SwaggerUIBundle.plugins.DownloadUrl
        ],
        layout: "StandaloneLayout",
        filter: true,
        requestInterceptor: function (request) {
            var headers = request.headers || {};
            var csrftoken = document.querySelector("[name=csrfmiddlewaretoken]");
            if (csrftoken) {
                headers["X-CSRFToken"] = csrftoken.value;
            }
            return request;
        }
    };

    var swaggerSettings = JSON.parse(document.getElementById('swagger-settings').innerHTML);
    if (!('oauth2RedirectUrl' in swaggerSettings)) {
        var oauth2RedirectUrl = document.getElementById('oauth2-redirect-url');
        if (oauth2RedirectUrl) {
            swaggerSettings['oauth2RedirectUrl'] = oauth2RedirectUrl.href;
            oauth2RedirectUrl.parentNode.removeChild(oauth2RedirectUrl);
        }
    }

    console.log('swaggerSettings', swaggerSettings);
    for (var p in swaggerSettings) {
        if (swaggerSettings.hasOwnProperty(p)) {
            swaggerConfig[p] = swaggerSettings[p];
        }
    }
    window.ui = SwaggerUIBundle(swaggerConfig);

    var oauth2Config = JSON.parse(document.getElementById('oauth2-config').innerHTML);
    console.log('oauth2Config', oauth2Config);
    window.ui.initOAuth(oauth2Config);
}

window.onload = function () {
    initSwaggerUi();
};

if (document.querySelector('.auth-wrapper .authorize')) {
    patchSwaggerUi();
}
else {
    insertionQ('.auth-wrapper .authorize').every(patchSwaggerUi);
}
