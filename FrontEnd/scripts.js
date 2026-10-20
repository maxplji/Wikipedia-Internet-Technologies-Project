/* PANEL WYGLĄDU STRONY */

    // Otwieranie i zamykanie panelu
    document.addEventListener("click", function (event) {
        if (event.target.closest("#wyglad-toggle")) {
            document.getElementById("panel-wygladu")
                ?.classList.toggle("otwarty");
        }

        if (event.target.closest("#zamknij-wyglad")) {
            document.getElementById("panel-wygladu")
                ?.classList.remove("otwarty");
        }
    });


    /* Zapisywanie i odtwarzanie ustawień wyglądu */
    function zastosujUstawieniaWygladu() {
        const tresc = document.querySelector(".tresc-strony");
        const uklad = document.querySelector(".uklad-strony");

        const rozmiar = localStorage.getItem("rozmiarTekstu") || "standardowy";
        const szerokosc = localStorage.getItem("szerokoscStrony") || "standardowa";
        const kolor = localStorage.getItem("kolorStrony") || "jasny";

        if (tresc) {
            tresc.classList.remove("tekst-maly", "tekst-duzy");

            if (rozmiar === "maly") tresc.classList.add("tekst-maly");
            if (rozmiar === "duzy") tresc.classList.add("tekst-duzy");
        }

        if (uklad) {
            uklad.classList.toggle("strona-szeroka", szerokosc === "szeroka");
        }

        document.body.classList.toggle("tryb-ciemny", kolor === "ciemny");

        document.querySelectorAll('input[name="rozmiar"]').forEach(input => {
            input.checked = input.value === rozmiar;
        });

        document.querySelectorAll('input[name="szerokosc"]').forEach(input => {
            input.checked = input.value === szerokosc;
        });

        document.querySelectorAll('input[name="kolor"]').forEach(input => {
            input.checked = input.value === kolor;
        });
    }

    /* Przywróć ustawienia po załadowaniu dokumentu */
    document.addEventListener("DOMContentLoaded", zastosujUstawieniaWygladu);

    /* Zapisuj każdą zmianę ustawień */
    document.addEventListener("change", function (event) {
        const opcja = event.target;

        if (opcja.name === "rozmiar") {
            localStorage.setItem("rozmiarTekstu", opcja.value);
        }

        if (opcja.name === "szerokosc") {
            localStorage.setItem("szerokoscStrony", opcja.value);
        }

        if (opcja.name === "kolor") {
            localStorage.setItem("kolorStrony", opcja.value);
        }

        zastosujUstawieniaWygladu();
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