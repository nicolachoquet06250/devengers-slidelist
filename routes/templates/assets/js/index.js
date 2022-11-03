window.addEventListener("load", () => {
    if ("serviceWorker" in navigator) {
      navigator.serviceWorker.register("/service-worker.js");
    }

    var beforeInstallPrompt = null;
    window.addEventListener("beforeinstallprompt", eventHandler, hideBanner);

    function eventHandler(e) {
        e.preventDefault();

        beforeInstallPrompt = e;

        document.querySelector('#installBtn').removeAttribute('disabled'); 
        document.querySelector('#installBtn').addEventListener('click', install);
        document.querySelector('#close').addEventListener('click', hideBanner);
        document.querySelector('#install-pwa-alert').classList.add('show');
    }

    function hideBanner() {
        document.querySelector('#install-pwa-alert').classList.remove('show');
        document.querySelector('#installBtn').removeEventListener('click', install);
    }

    async function install() {
        if (beforeInstallPrompt) {
            beforeInstallPrompt.prompt();
            
            // Find out whether the user confirmed the installation or not
            const { outcome } = await beforeInstallPrompt.userChoice;
            // The deferredPrompt can only be used once.
            beforeInstallPrompt = null;
            // Act on the user's choice
            hideBanner();
        }
    }
  });