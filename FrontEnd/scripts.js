/* =========================================================
   1. USTAWIENIA WYGLĄDU STRONY
========================================================= */

function zastosujUstawieniaWygladu() {
    const trescStrony = document.querySelector(".tresc-strony");
    const ukladStrony = document.querySelector(".uklad-strony");

    if (!trescStrony || !ukladStrony) return;

    const rozmiarTekstu = localStorage.getItem("rozmiarTekstu") || "normalny";
    const szerokoscStrony = localStorage.getItem("szerokoscStrony") || "standardowa";
    const kolorStrony = localStorage.getItem("kolorStrony") || "jasny";

    // Rozmiar tekstu
    trescStrony.classList.remove("tekst-maly", "tekst-duzy");

    if (rozmiarTekstu === "maly") {
        trescStrony.classList.add("tekst-maly");
    } else if (rozmiarTekstu === "duzy") {
        trescStrony.classList.add("tekst-duzy");
    }

    // Szerokość strony
    ukladStrony.classList.remove("strona-szeroka");

    if (szerokoscStrony === "szeroka") {
        ukladStrony.classList.add("strona-szeroka");
    }

    // Jasny / ciemny motyw
    document.body.classList.toggle("tryb-ciemny", kolorStrony === "ciemny");

    // Synchronizacja przycisków wyboru
    document.querySelectorAll('input[name="rozmiar"]').forEach(radio => {
        radio.checked = radio.value === rozmiarTekstu;
    });

    document.querySelectorAll('input[name="szerokosc"]').forEach(radio => {
        radio.checked = radio.value === szerokoscStrony;
    });

    document.querySelectorAll('input[name="kolor"]').forEach(radio => {
        radio.checked = radio.value === kolorStrony;
    });
}


// Zapisywanie ustawień po zmianie
document.addEventListener("change", function (event) {
    const element = event.target;

    if (element.matches('input[name="rozmiar"]')) {
        localStorage.setItem("rozmiarTekstu", element.value);
        zastosujUstawieniaWygladu();
    }

    if (element.matches('input[name="szerokosc"]')) {
        localStorage.setItem("szerokoscStrony", element.value);
        zastosujUstawieniaWygladu();
    }

    if (element.matches('input[name="kolor"]')) {
        localStorage.setItem("kolorStrony", element.value);
        zastosujUstawieniaWygladu();
    }
});


/* =========================================================
   2. PANEL WYGLĄDU
========================================================= */

document.addEventListener("click", function (event) {
    const panelWygladu = document.getElementById("panel-wygladu");

    if (!panelWygladu) return;

    if (event.target.closest("#wyglad-toggle")) {
        panelWygladu.classList.toggle("otwarty");
    }

    if (event.target.closest("#zamknij-wyglad")) {
        panelWygladu.classList.remove("otwarty");
    }
});


/* =========================================================
   3. MENU NARZĘDZIA
========================================================= */

document.addEventListener("click", function (event) {
    const przycisk = event.target.closest(".narzedzia-button");
    const panel = document.getElementById("panel-narzedzi");

    if (!przycisk || !panel) return;

    panel.classList.toggle("otwarte");
});



/* =========================================================
   4. INFORMACJE O UŻYTKOWNIKU
========================================================= */

function pokazZalogowanegoUzytkownika(nick) {
    const przyciskRejestracji = document.getElementById("utworzenie-k");
    const przyciskLogowania = document.getElementById("logowanie-k");
    const nickUzytkownika = document.getElementById("nick-uzytkownika");

    if (przyciskRejestracji) {
        przyciskRejestracji.style.display = "none";
    }

    if (przyciskLogowania) {
        przyciskLogowania.style.display = "none";
    }

    if (nickUzytkownika) {
        nickUzytkownika.textContent = nick;
        nickUzytkownika.style.display = "block";
    }
}


/* =========================================================
   5. DANE ZDJĘĆ ARTYKUŁU
========================================================= */

const obrazyArtykulu = new Map();
let licznikObrazow = 0;


/* =========================================================
   6. POMOCNICZE FUNKCJE EDYTORA
========================================================= */

// Wstawianie tekstu w miejscu kursora
function wstawTekstDoEdytora(tekst) {
    const pole = document.getElementById("tresc-artykulu");

    if (!pole) return;

    const start = pole.selectionStart;
    const koniec = pole.selectionEnd;

    pole.setRangeText(tekst, start, koniec, "end");

    pole.focus();
    aktualizujPodglad();
}


// Zabezpieczenie tekstu przed interpretacją jako HTML
function zabezpieczHTML(tekst) {
    return String(tekst).replace(/[&<>"]/g, znak => {
        const znaki = {
            "&": "&amp;",
            "<": "&lt;",
            ">": "&gt;",
            '"': "&quot;"
        };

        return znaki[znak];
    });
}


// Sprawdzenie bezpiecznego adresu internetowego
function bezpiecznyAdres(adres) {
    try {
        const url = new URL(adres);

        return ["http:", "https:"].includes(url.protocol)
            ? url.href
            : null;
    } catch {
        return null;
    }
}


/* =========================================================
   7. FORMATOWANIE TEKSTU
========================================================= */

function formatujTekst(tekst) {
    let wynik = zabezpieczHTML(tekst);

    // Pogrubienie
    wynik = wynik.replace(
        /'''([\s\S]+?)'''/g,
        "<strong>$1</strong>"
    );

    // Kursywa
    wynik = wynik.replace(
        /''([\s\S]+?)''/g,
        "<em>$1</em>"
    );

    // Podkreślenie
    wynik = wynik.replace(
        /&lt;u&gt;([\s\S]+?)&lt;\/u&gt;/g,
        "<u>$1</u>"
    );

    // Linki: [https://adres.pl opis]
    wynik = wynik.replace(
        /\[(https?:\/\/[^\s\]]+)(?:\s+([^\]]+))?\]/g,
        function (_, adres, opis) {
            const bezpiecznyURL = bezpiecznyAdres(adres);

            if (!bezpiecznyURL) {
                return opis || adres;
            }

            return `<a href="${bezpiecznyURL}" target="_blank" rel="noopener noreferrer">${opis || adres}</a>`;
        }
    );

    return wynik;
}


/* =========================================================
   8. PRZYGOTOWANIE PODGLĄDU ARTYKUŁU
========================================================= */

function przygotujPodglad(tekst) {
    const linie = tekst.split("\n");

    let wynik = "";
    let listaOtwarta = false;

    function zamknijListe() {
        if (listaOtwarta) {
            wynik += "</ul>";
            listaOtwarta = false;
        }
    }

    linie.forEach(linia => {
        const tresc = linia.trim();

        // Pusta linia
        if (!tresc) {
            zamknijListe();
            return;
        }

        // Znacznik zdjęcia
        const znacznikZdjecia = tresc.match(/^\[\[OBRAZ:(\d+)\]\]$/);

        if (znacznikZdjecia) {
            zamknijListe();

            const id = Number(znacznikZdjecia[1]);
            const obraz = obrazyArtykulu.get(id);

            if (obraz) {
                const polozenie = obraz.polozenie || "srodek";
                const opis = zabezpieczHTML(obraz.opis || "");
                const podpis = obraz.podpis
                    ? `<figcaption>${zabezpieczHTML(obraz.podpis)}</figcaption>`
                    : "";

                wynik += `
                    <figure class="zdjecie-artykulu polozenie-${polozenie}">
                        <img src="${zabezpieczHTML(obraz.adres)}"
                             alt="${opis}">
                        ${podpis}
                    </figure>
                `;
            }

            return;
        }

        // Nagłówki w formacie == tekst ==
        const naglowekWiki = tresc.match(/^==\s*(.*?)\s*==$/);

        if (naglowekWiki) {
            zamknijListe();
            wynik += `<h2>${formatujTekst(naglowekWiki[1])}</h2>`;
            return;
        }

        // Nagłówki w formacie Markdown
        const naglowekMarkdown = tresc.match(/^(#{1,3})\s+(.+)$/);

        if (naglowekMarkdown) {
            zamknijListe();

            const poziom = naglowekMarkdown[1].length + 1;

            wynik += `<h${poziom}>${formatujTekst(naglowekMarkdown[2])}</h${poziom}>`;
            return;
        }

        // Element listy
        if (/^[*-]\s+/.test(tresc)) {
            if (!listaOtwarta) {
                wynik += "<ul>";
                listaOtwarta = true;
            }

            const elementListy = tresc.replace(/^[*-]\s+/, "");

            wynik += `<li>${formatujTekst(elementListy)}</li>`;
            return;
        }

        // Zwykły akapit
        zamknijListe();
        wynik += `<p>${formatujTekst(tresc)}</p>`;
    });

    zamknijListe();

    return wynik;
}


/* =========================================================
   9. AUTOMATYCZNA AKTUALIZACJA PODGLĄDU
========================================================= */

function aktualizujPodglad() {
    const poleTytulu = document.getElementById("tytul-artykulu");
    const poleTresci = document.getElementById("tresc-artykulu");

    const podgladTytulu = document.getElementById("podglad-tytulu");
    const podgladTresci = document.getElementById("podglad-tresci");

    if (!poleTytulu || !poleTresci || !podgladTytulu || !podgladTresci) {
        return;
    }

    podgladTytulu.textContent =
        poleTytulu.value.trim() || "Tytuł artykułu";

    podgladTresci.innerHTML =
        przygotujPodglad(poleTresci.value);
}


// Aktualizacja podczas pisania
document.addEventListener("input", function (event) {
    if (
        event.target.matches("#tytul-artykulu") ||
        event.target.matches("#tresc-artykulu")
    ) {
        aktualizujPodglad();
    }
});


/* =========================================================
   10. PRZYCISKI FORMATOWANIA
========================================================= */

document.addEventListener("click", function (event) {
    const przycisk = event.target.closest("[data-format]");

    if (!przycisk) return;

    const pole = document.getElementById("tresc-artykulu");

    if (!pole) return;

    const format = przycisk.dataset.format;

    const zaznaczonyTekst = pole.value.substring(
        pole.selectionStart,
        pole.selectionEnd
    );

    let tekstDoWstawienia = "";

    switch (format) {
        case "bold":
            tekstDoWstawienia = `'''${zaznaczonyTekst || "pogrubiony tekst"}'''`;
            break;

        case "italic":
            tekstDoWstawienia = `''${zaznaczonyTekst || "pochylony tekst"}''`;
            break;

        case "underline":
            tekstDoWstawienia = `<u>${zaznaczonyTekst || "podkreślony tekst"}</u>`;
            break;

        case "heading":
            tekstDoWstawienia = `== ${zaznaczonyTekst || "Nagłówek"} ==`;
            break;

        case "list":
            tekstDoWstawienia = `* ${zaznaczonyTekst || "Element listy"}`;
            break;

        case "link": {
            const adres = prompt("Podaj adres strony:");

            if (!adres) return;

            const bezpiecznyURL = bezpiecznyAdres(adres);

            if (!bezpiecznyURL) {
                alert("Podany adres jest nieprawidłowy.");
                return;
            }

            const opis = zaznaczonyTekst || bezpiecznyURL;

            tekstDoWstawienia = `[${bezpiecznyURL} ${opis}]`;
            break;
        }

       case "image": {
            const panelZdjecia = document.getElementById("panel-dodawania-zdjecia");

            if (!panelZdjecia) {
                console.error("Nie znaleziono panelu dodawania zdjęcia!");
                return;
            }

            panelZdjecia.hidden = false;
            panelZdjecia.classList.add("otwarty");

            break;
        }

        default:
            return;
    }

    wstawTekstDoEdytora(tekstDoWstawienia);
});


/* =========================================================
   11. PANEL DODAWANIA ZDJĘĆ
========================================================= */

document.addEventListener("click", function (event) {
    const panel = document.querySelector(".panel-dodawania-zdjecia");

    if (!panel) return;

    // Anulowanie dodawania zdjęcia
    if (event.target.closest("#anuluj-zdjecie")) {
    panel.classList.remove("otwarty");
    panel.hidden = true;
}

    // Dodawanie zdjęcia
    if (event.target.closest("#wstaw-zdjecie")) {
        dodajZdjecie();
    }
});


/* =========================================================
   12. DODAWANIE ZDJĘCIA DO ARTYKUŁU
========================================================= */

function dodajZdjecie() {
    const polePliku = document.getElementById("plik-zdjecia");
    const poleURL = document.getElementById("url-zdjecia");
    const poleOpisu = document.getElementById("opis-zdjecia");
    const polePodpisu = document.getElementById("podpis-zdjecia");
    const polePolozenia = document.getElementById("polozenie-zdjecia");

    if (!polePliku || !poleURL) return;

    const plik = polePliku.files[0];
    const adresURL = poleURL.value.trim();

    let adresZdjecia = "";

    // Zdjęcie z komputera
    if (plik) {
        if (!plik.type.startsWith("image/")) {
            alert("Wybrany plik nie jest zdjęciem.");
            return;
        }

        adresZdjecia = URL.createObjectURL(plik);
    }

    // Zdjęcie z adresu internetowego
    else if (adresURL) {
        adresZdjecia = bezpiecznyAdres(adresURL);

        if (!adresZdjecia) {
            alert("Podany adres zdjęcia jest nieprawidłowy.");
            return;
        }
    }

    else {
        alert("Wybierz zdjęcie z komputera lub podaj jego adres URL.");
        return;
    }

    const id = ++licznikObrazow;

    const polozenie = polePolozenia?.value || "srodek";

    obrazyArtykulu.set(id, {
        adres: adresZdjecia,
        opis: poleOpisu?.value.trim() || "",
        podpis: polePodpisu?.value.trim() || "",
        polozenie: polozenie
    });

    // Wstawienie znacznika zdjęcia do treści
    wstawTekstDoEdytora(`\n[[OBRAZ:${id}]]\n`);

    // Wyczyszczenie formularza
    polePliku.value = "";
    poleURL.value = "";

    if (poleOpisu) poleOpisu.value = "";
    if (polePodpisu) polePodpisu.value = "";

    if (polePolozenia) {
        polePolozenia.value = "srodek";
    }

    // Zamknięcie panelu
    const panel = document.getElementById("panel-dodawania-zdjecia");

    if (panel) {
        panel.classList.remove("otwarty");
        panel.hidden = true;
    }
}


/* =========================================================
   13. INICJALIZACJA STRONY
========================================================= */

document.addEventListener("DOMContentLoaded", function () {
    zastosujUstawieniaWygladu();
    aktualizujPodglad();
});