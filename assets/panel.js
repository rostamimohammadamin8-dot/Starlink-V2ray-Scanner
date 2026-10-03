(function () {
  "use strict";

  var RTL_LANGUAGES = ["fa", "ar", "ur", "ps"];
  var BOOTSTRAP_CSS = "https://cdn.jsdelivr.net/npm/bootstrap@5.3.0/dist/css/";
  var translations = window.I18N || {};
  var languages = window.I18N_LANGUAGES || [["en", "English"]];

  var langSelect = document.getElementById("lang-select");
  var countrySelect = document.getElementById("country-select");
  var staticOnly = document.getElementById("static-only");
  var mainSub = document.getElementById("main-sub");
  var visibleCount = document.getElementById("visible-count");
  var emptyMessage = document.getElementById("empty-message");
  var cards = Array.prototype.slice.call(document.querySelectorAll(".node-card"));

  function detectLanguage() {
    var saved = localStorage.getItem("lang");
    if (saved && translations[saved]) {
      return saved;
    }
    var preferred = navigator.languages || [navigator.language || "en"];
    for (var i = 0; i < preferred.length; i++) {
      var code = String(preferred[i]).toLowerCase().split("-")[0];
      if (translations[code]) {
        return code;
      }
    }
    return "en";
  }

  var lang = detectLanguage();

  function t(key) {
    var table = translations[lang] || {};
    return table[key] || (translations.en || {})[key] || key;
  }

  function flag(code) {
    if (!/^[A-Z]{2}$/.test(code) || code === "XX") {
      return "";
    }
    return String.fromCodePoint(0x1f1e6 + code.charCodeAt(0) - 65, 0x1f1e6 + code.charCodeAt(1) - 65) + " ";
  }

  function countryName(code) {
    if (code === "XX") {
      return t("unknown");
    }
    try {
      return new Intl.DisplayNames([lang], { type: "region" }).of(code) || code;
    } catch (e) {
      return code;
    }
  }

  function applyLanguage() {
    var rtl = RTL_LANGUAGES.indexOf(lang) !== -1;
    document.documentElement.lang = lang;
    document.documentElement.dir = rtl ? "rtl" : "ltr";
    document.getElementById("bootstrap-css").href = BOOTSTRAP_CSS + (rtl ? "bootstrap.rtl.min.css" : "bootstrap.min.css");

    document.querySelectorAll("[data-i18n]").forEach(function (el) {
      el.textContent = t(el.dataset.i18n);
    });
    document.querySelectorAll("[data-i18n-title]").forEach(function (el) {
      el.title = t(el.dataset.i18nTitle);
    });
    document.querySelectorAll(".country-name").forEach(function (el) {
      el.textContent = flag(el.dataset.code) + countryName(el.dataset.code);
    });

    var options = Array.prototype.slice.call(countrySelect.options, 1);
    var selected = countrySelect.value;
    options.forEach(function (option) {
      option.textContent = flag(option.value) + countryName(option.value) + " (" + option.dataset.count + ")";
    });
    options.sort(function (a, b) {
      return countryName(a.value).localeCompare(countryName(b.value), lang);
    });
    options.forEach(function (option) {
      countrySelect.appendChild(option);
    });
    countrySelect.value = selected;
  }

  function applyFilter() {
    var country = countrySelect.value;
    var onlyStatic = staticOnly.checked;
    var shown = 0;
    cards.forEach(function (card) {
      var match = country ? card.dataset.country === country : card.dataset.top === "1";
      if (onlyStatic && card.dataset.static !== "1") {
        match = false;
      }
      card.classList.toggle("d-none", !match);
      if (match) {
        shown++;
        card.querySelector(".node-index").textContent = shown;
      }
    });
    visibleCount.textContent = shown;
    emptyMessage.classList.toggle("d-none", shown > 0);

    if (country) {
      mainSub.href = "subs/countries/" + country + ".txt";
    } else if (onlyStatic) {
      mainSub.href = "subs/static-ip.txt";
    } else {
      mainSub.href = "cleaned_configs.txt";
    }
  }

  languages.forEach(function (entry) {
    var option = document.createElement("option");
    option.value = entry[0];
    option.textContent = entry[1];
    langSelect.appendChild(option);
  });
  langSelect.value = lang;

  langSelect.addEventListener("change", function () {
    lang = langSelect.value;
    localStorage.setItem("lang", lang);
    applyLanguage();
  });
  countrySelect.addEventListener("change", applyFilter);
  staticOnly.addEventListener("change", applyFilter);

  document.getElementById("nodes").addEventListener("click", function (event) {
    var button = event.target.closest(".copy-btn");
    if (!button) {
      return;
    }
    navigator.clipboard.writeText(button.dataset.config).then(function () {
      alert(t("copied"));
    });
  });

  applyLanguage();
  applyFilter();
})();
