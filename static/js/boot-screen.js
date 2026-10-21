document.addEventListener("DOMContentLoaded", () => {
    const screen = document.getElementById("boot-screen");
    const log = document.getElementById("boot-log");
    const progress = document.getElementById("boot-progress-bar");
    const complete = document.getElementById("boot-complete");

    if (!screen || !log || !progress || !complete) {
        return;
    }

    const sessionKey = "retro-cv-booted";
    const reducedMotion = window.matchMedia(
        "(prefers-reduced-motion: reduce)"
    ).matches;

    // Skip the boot screen if it already ran in this session.
    try {
        if (sessionStorage.getItem(sessionKey) === "true") {
            screen.remove();
            document.dispatchEvent(new Event("boot-finished"));
            return;
        }
    } catch {
        // If sessionStorage is unavailable, run the boot sequence.
    }

    // Mark this session so navigation and refreshes skip the boot.
    try {
        sessionStorage.setItem(sessionKey, "true");
    } catch {
        // The animation can still run without sessionStorage.
    }

    if (reducedMotion) {
        screen.remove();
        document.dispatchEvent(new Event("boot-finished"));
        return;
    }

    const messages = [
        "CJUROZ PERSONAL COMPUTER BIOS v0.1.0",
        "COPYRIGHT (C) 2026 CJUROZ SYSTEMS",
        "",
        "POST: POWER-ON SELF TEST .... PASS",
        "CPU: GO RUNTIME ............. OK",
        "MEMORY: PROFILE DATA ........ LOADED",
        "STORAGE: STATIC ASSETS ...... OK",
        "NETWORK: HTTP SERVER ....... READY",
        "",
        "INITIALIZING USER INTERFACE...",
        "ALL SYSTEMS NOMINAL"
    ];

    let index = 0;

    function showNextMessage() {
        if (index >= messages.length) {
            progress.style.width = "100%";
            complete.classList.add("visible");

            window.setTimeout(() => {
                screen.classList.add("boot-finished");

                // Tell the status indicator the boot sequence has finished.
                document.dispatchEvent(new Event("boot-finished"));

                // Remove the overlay after its fade-out.
                window.setTimeout(() => screen.remove(), 550);
            }, 650);

            return;
        }

        const line = document.createElement("div");
        line.textContent = messages[index];
        log.appendChild(line);

        index++;
        progress.style.width = `${(index / messages.length) * 100}%`;

        window.setTimeout(showNextMessage, 450);
    }

    window.setTimeout(showNextMessage, 250);
});
