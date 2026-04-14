import zhCN from "./locales/zh-CN.js";
import en from "./locales/en.js";

const locales = { "zh-CN": zhCN, en };

let currentLocale = "zh-CN";
let messages = zhCN;

export function setLocale(locale) {
  currentLocale = locale;
  messages = locales[locale] || locales["zh-CN"];
}

export function locale() {
  return currentLocale;
}

export function t(key, params) {
  const keys = key.split(".");
  let value = messages;
  for (const k of keys) {
    if (value && typeof value === "object" && k in value) {
      value = value[k];
    } else {
      return key;
    }
  }
  if (typeof value !== "string") {
    return key;
  }
  if (params) {
    return value.replace(/\{(\w+)\}/g, (_, name) =>
      params[name] !== undefined ? params[name] : `{${name}}`
    );
  }
  return value;
}
