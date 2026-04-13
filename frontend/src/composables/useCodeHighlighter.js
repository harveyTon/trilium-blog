import { bundledLanguages, getSingletonHighlighter } from "shiki";

const THEMES = {
  light: "github-light",
  dark: "github-dark",
};

let highlighterPromise;
const loadedLanguages = new Set(["plaintext"]);

function escapeHtml(value) {
  return value
    .replaceAll("&", "&amp;")
    .replaceAll("<", "&lt;")
    .replaceAll(">", "&gt;")
    .replaceAll('"', "&quot;")
    .replaceAll("'", "&#39;");
}

function normalizeLanguage(language) {
  if (!language) return "plaintext";
  const normalized = String(language).trim().toLowerCase();
  if (!normalized) return "plaintext";
  if (normalized in bundledLanguages) {
    return normalized;
  }
  if (normalized === "text") {
    return "plaintext";
  }
  return "plaintext";
}

async function getHighlighter() {
  if (!highlighterPromise) {
    highlighterPromise = getSingletonHighlighter({
      themes: [THEMES.light, THEMES.dark],
      langs: ["plaintext"],
    });
  }

  return highlighterPromise;
}

export async function highlightCodeBlock({ code, language, dark }) {
  const rawCode = typeof code === "string" ? code : "";
  const resolvedLanguage = normalizeLanguage(language);

  try {
    const highlighter = await getHighlighter();
    if (!loadedLanguages.has(resolvedLanguage) && resolvedLanguage in bundledLanguages) {
      await highlighter.loadLanguage(resolvedLanguage);
      loadedLanguages.add(resolvedLanguage);
    }
    const html = highlighter.codeToHtml(rawCode, {
      lang: resolvedLanguage,
      theme: dark ? THEMES.dark : THEMES.light,
    });

    return {
      html,
      resolvedLanguage,
    };
  } catch (error) {
    console.error("Failed to highlight code block:", error);
    return {
      html: `<pre class="shiki shiki-fallback"><code>${escapeHtml(rawCode)}</code></pre>`,
      resolvedLanguage: "plaintext",
    };
  }
}
