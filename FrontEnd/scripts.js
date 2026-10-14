document.addEventListener("DOMContentLoaded", function () {
    // Menu główne
    const przyciskMenu = document.getElementById("menu-toggle");
    const panelMenu = document.getElementById("panel-menu");

    if (przyciskMenu && panelMenu) {
        przyciskMenu.addEventListener("click", function () {
            const otwarte = panelMenu.classList.toggle("otwarte");
            przyciskMenu.setAttribute("aria-expanded", otwarte);
        });
    }

    // Panel Wyglądu
    const przyciskWygladu = document.getElementById("wyglad-toggle");
    const panelWygladu = document.getElementById("panel-wygladu");
    const przyciskZamknij = document.getElementById("zamknij-wyglad");

    if (przyciskWygladu && panelWygladu) {
        przyciskWygladu.addEventListener("click", function () {
            panelWygladu.classList.toggle("otwarty");
        });
    }

    if (przyciskZamknij && panelWygladu) {
        przyciskZamknij.addEventListener("click", function () {
            panelWygladu.classList.remove("otwarty");
        });
    }
});
