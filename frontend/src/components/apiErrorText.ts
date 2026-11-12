import type { TFunction } from 'i18next';
import en from '../i18n/locales/en';
import type { ApiError } from '../models/ApiError';

type ApiErrorCode = keyof typeof en.apiErrors;

function isKnownCode(code: string): code is ApiErrorCode {
  return Object.prototype.hasOwnProperty.call(en.apiErrors, code);
}

// apiErrorText words an API error with t. An unknown code, or no error body
// at all, gets a generic text.
export function apiErrorText(
  apiError: ApiError | undefined,
  t: TFunction
): string {
  const code = apiError?.error ?? '';
  if (!isKnownCode(code)) {
    return t('apiErrors.unknown');
  }

  return t(`apiErrors.${code}`, {
    name: apiError?.name ?? t('apiErrors.somePlayer'),
  });
}
