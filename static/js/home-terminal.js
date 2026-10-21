document.addEventListener("DOMContentLoaded", () => {
    const openButton = document.getElementById("open-terminal");
    const terminal = document.getElementById("terminal");

    if (!openButton || !terminal) {
        return;
    }

    openButton.addEventListener("click", (event) => {
        event.preventDefault();

        // Reveal the terminal.
        terminal.hidden = false;
        openButton.setAttribute("aria-expanded", "true");

        // Update the URL without reloading the page.
        history.replaceState(null, "", "#terminal");

        // Scroll to the terminal after it becomes visible.
        terminal.scrollIntoView({
            behavior: window.matchMedia(
                "(prefers-reduced-motion: reduce)"
            ).matches ? "auto" : "smooth",
            block: "start"
        });
    });
});