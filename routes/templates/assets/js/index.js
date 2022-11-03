window.addEventListener("load", () => {
    if ("serviceWorker" in navigator) {
      navigator.serviceWorker.register("/service-worker.js");
    }

    var beforeInstallPrompt = null;
    window.addEventListener("beforeinstallprompt", eventHandler, errorHandler);

    function eventHandler(e) {
        e.preventDefault();

        beforeInstallPrompt = e;

        document.querySelector('#installBtn').removeAttribute('disabled'); 
        document.querySelector('#installBtn').addEventListener('click', install);
        document.querySelector('#close').addEventListener('click', errorHandler);
        document.querySelector('#install-pwa-alert').classList.add('show');
    }

    function errorHandler(e) {
        console.log('error: ' + e);

        document.querySelector('#install-pwa-alert').classList.remove('show');
        document.querySelector('#installBtn').removeEventListener('click', install);
    }

    function install() {
        if (beforeInstallPrompt) {
            beforeInstallPrompt.prompt();
        }
    }
  });