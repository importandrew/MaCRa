import i18n from 'i18next';
import { initReactI18next } from 'react-i18next';
import en from './locales/en';
import ru from './locales/ru';

// LANGUAGES are the languages the app speaks, each named in itself, as the
// language picker shows them. Adding one: a file in locales/ typed as
// Messages, an entry here and in resources below.
export const LANGUAGES = [
  { code: 'en', label: 'English' },
  { code: 'ru', label: 'Русский' },
] as const;

export type Language = (typeof LANGUAGES)[number]['code'];

const FALLBACK: Language = 'en';

// STORAGE_KEY keeps the chosen language on this device.
const STORAGE_KEY = 'language';

function isLanguage(code: string | null | undefined): code is Language {
  return LANGUAGES.some((lang) => lang.code === code);
}

// initialLanguage is the language chosen on this device, or else the first of
// the browser's languages the app speaks ("ru-RU" → "ru"), or English.
function initialLanguage(): Language {
  try {
    const stored = localStorage.getItem(STORAGE_KEY);
    if (isLanguage(stored)) return stored;
  } catch {
    // Storage blocked (private mode): fall through to the browser language.
  }

  for (const tag of navigator.languages ?? [navigator.language]) {
    const code = tag.toLowerCase().split('-')[0];
    if (isLanguage(code)) return code;
  }

  return FALLBACK;
}

// currentLanguage is the language the app shows now.
export function currentLanguage(): Language {
  return isLanguage(i18n.language) ? i18n.language : FALLBACK;
}

// setLanguage switches the app to lang and remembers it on this device.
export function setLanguage(lang: Language) {
  try {
    localStorage.setItem(STORAGE_KEY, lang);
  } catch {
    // Not remembered, but the switch still applies until the page reloads.
  }
  void i18n.changeLanguage(lang);
}

i18n.on('languageChanged', (lang) => {
  document.documentElement.lang = lang;
});

void i18n.use(initReactI18next).init({
  resources: {
    en: { translation: en },
    ru: { translation: ru },
  },
  lng: initialLanguage(),
  fallbackLng: FALLBACK,
  // React escapes text itself.
  interpolation: { escapeValue: false },
});

document.documentElement.lang = currentLanguage();

export default i18n;
