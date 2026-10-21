document.addEventListener("DOMContentLoaded", () => {
    const openButton = document.getElementById("open-terminal");
    const terminal = document.getElementById("terminal");

    if (!openButton || !terminal) {
        return;
    }

    openButton.addEventListener("click", (event) => {
        event.preventDefault();

        terminal.classList.remove("terminal-collapsed");
        terminal.classList.add("terminal-expanded");

        openButton.setAttribute("aria-expanded", "true");
        openButton.textContent = "TERMINAL OPEN ✓";

        history.replaceState(null, "", "#terminal");

        terminal.scrollIntoView({
            behavior: window.matchMedia(
                "(prefers-reduced-motion: reduce)"
            ).matches ? "auto" : "smooth",
            block: "start"
        });
    });
});