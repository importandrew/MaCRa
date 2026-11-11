import 'i18next';
import type en from './locales/en';

// Type the keys of t(): a key missing from the English texts is a build error.
declare module 'i18next' {
  interface CustomTypeOptions {
    defaultNS: 'translation';
    resources: {
      translation: typeof en;
    };
  }
}
