import i18n from "i18next";
import { initReactI18next } from "react-i18next";
import type { I18nProvider } from "@refinedev/core";
import zhCN from "./locales/zh-CN.json";
import en from "./locales/en.json";

// Supported languages: zh-CN (default) + en.
// The backend never translates: it returns { msgKey, msgParams } (see
// backend/core/i18n/keys.go) and the frontend renders them as
// t(`backend.${msgKey}`) from the dictionaries below.
export const SUPPORTED_LANGS = [
  { value: "zh-CN", label: "中文" },
  { value: "en", label: "EN" },
] as const;

export const LANG_KEY = "lang";
export const DEFAULT_LANG = "zh-CN";

export function getLang(): string {
  const saved = localStorage.getItem(LANG_KEY);
  if (saved === "en" || saved === "zh-CN") return saved;
  return DEFAULT_LANG;
}

i18n.use(initReactI18next).init({
  resources: {
    "zh-CN": { translation: zhCN },
    en: { translation: en },
  },
  lng: getLang(),
  fallbackLng: DEFAULT_LANG,
  interpolation: { escapeValue: false },
});

export function setLang(lang: string) {
  const next = lang === "en" ? "en" : "zh-CN";
  localStorage.setItem(LANG_KEY, next);
  void i18n.changeLanguage(next);
}

// tBackend renders a backend msgKey: tBackend("auth.invalidCredentials").
// Falls back to the backend-provided msg (default-language text) so unknown
// keys never render blank.
export function tBackend(msgKey?: string, params?: Record<string, any>, fallback?: string): string {
  if (!msgKey) return fallback || (i18n.t("common.requestFailed") as string);
  const text = i18n.t(`backend.${msgKey}`, { ...(params ?? {}), defaultValue: "" }) as string;
  return text || fallback || msgKey;
}

// backendErrorMessage turns a thrown unwrap() error into display text.
export function backendErrorMessage(err: any): string {
  if (err?.i18nKey) return tBackend(err.i18nKey, err.i18nParams, err.message);
  return err?.message || i18n.t("common.requestFailed");
}

// refineI18nProvider wires react-i18next into refine (delete confirmations,
// notifications, useTranslate). Our own UI uses useTranslation() directly.
export const refineI18nProvider: I18nProvider = {
  translate: (key: string, options?: any, defaultMessage?: string) =>
    i18n.t(key, { ...(options ?? {}), defaultValue: defaultMessage ?? key }) as string,
  changeLocale: (locale: string) => {
    setLang(locale);
    return Promise.resolve();
  },
  getLocale: () => i18n.language,
};

export default i18n;
