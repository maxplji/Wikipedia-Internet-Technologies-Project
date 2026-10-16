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

// Zalogowany użytkownik (znika login oraz tworzenie konta -> pojawia się użytkownik)
function pokazZalogowanegoUzytkownika(nick) {
    const utworzKonto = document.getElementById("utworzenie-k");
    const zaloguj = document.getElementById("logowanie-k");
    const nickElement = document.getElementById("nick-uzytkownika");

    if (!utworzKonto || !zaloguj || !nickElement) {
        return;
    }

    utworzKonto.hidden = true;
    zaloguj.hidden = true;

    nickElement.textContent = nick;
    nickElement.hidden = false;
}
