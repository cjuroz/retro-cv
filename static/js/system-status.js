(() => {
    const systemStatus = document.getElementById("system-status");

    if (!systemStatus) {
        return;
    }

    const statusText = systemStatus.querySelector(".status-text");
    const bootScreen = document.getElementById("boot-screen");

    let bootFinished = !bootScreen;
    let pageLoaded = document.readyState === "complete";
    let timerStarted = false;

    function updateStatus() {
        if (!bootFinished || !pageLoaded || timerStarted) {
            return;
        }

        timerStarted = true;

        window.setTimeout(() => {
            systemStatus.classList.remove("is-loading");
            systemStatus.classList.add("is-online");
            statusText.textContent = "SYSTEM ONLINE";
        }, 1200);
    }

    document.addEventListener("boot-finished", () => {
        bootFinished = true;
        updateStatus();
    }, { once: true });

    if (!pageLoaded) {
        window.addEventListener("load", () => {
            pageLoaded = true;
            updateStatus();
        }, { once: true });
    }

    updateStatus();
})();
